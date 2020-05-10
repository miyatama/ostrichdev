package ostrich

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Ostrich struct {
	Repository    string
	FromBranch    string
	OstrichBranch string
	CommitId      string
	FileAccessor  FileAccesserInterface
}

func (o *Ostrich) Run() error {

	// remove working directory
	repositoryName, err := o.getRepositoryName(o.Repository)
	if err != nil {
		return err
	}

	git := o.getGitCommand()

	// get from branch
	if err := os.RemoveAll(repositoryName); err != nil {
		return err
	}
	if err := git.Clone(o.Repository); err != nil {
		return err
	}
	if err := o.chDir(repositoryName); err != nil {
		return err
	}
	defer o.chDir("..")
	o.showGitVersion(git)
	if err := o.checkout(o.FromBranch, git); err != nil {
		return err
	}
	if err := git.Pull(o.FromBranch); err != nil {
		return err
	}
	if err := git.Fetch(); err != nil {
		return err
	}

	// apply commit to ostrich branch
	commitTexts, err := git.Show(o.CommitId)
	if err != nil {
		return err
	}
	commit, err := o.parseCommit(commitTexts)
	if err != nil {
		return err
	}

	if err := git.Checkout(o.OstrichBranch); err != nil {
		return err
	}
	if err := git.Reset(o.FromBranch); err != nil {
		return err
	}
	if err := o.applyCommit(commit, git); err != nil {
		return err
	}

	// commit and push to ostrich branch
	if err := git.Commit(commit.Message); err != nil {
		return err
	}
	if err := git.Push(o.OstrichBranch); err != nil {
		return err
	}
	return nil
}

func (o *Ostrich) getGitCommand() GitCommand {
	return GitCommand{
		executor: &CommandExecutor{},
	}
}

func (o *Ostrich) parseCommit(commitTexts []string) (Commit, error) {

	if len(commitTexts) < 5 {
		return Commit{},
			fmt.Errorf(
				"invalid commit texts.text line count is %d",
				len(commitTexts))
	}

	getAuthor := func(text string) (string, error) {
		terms := strings.Split(text, " ")
		if len(terms) < 2 {
			return "", fmt.Errorf("can not detect author %s.", text)
		}
		return terms[1], nil
	}
	getCommitDate := func(text string) (time.Time, error) {
		terms := strings.Split(text, " ")
		if len(terms) < 7 {
			return time.Now(), fmt.Errorf("can not detect date %s", text)
		}
		// for japanese
		text = strings.Replace(text, "Date:", "", 1)
		text = strings.Trim(text, " ")
		format := "Mon Jan 2 15:04:05 2006 -0700"
		date, err := time.Parse(format, text)
		if err != nil {
			return time.Now(), err
		}
		return date, nil
	}

	// get author, commit date and message
	author := ""
	commitDate := time.Now()
	message := ""
	err := errors.New("")
	for _, text := range commitTexts {
		if strings.HasPrefix(text, "commit ") {
			continue
		}
		if strings.HasPrefix(text, "Author") {
			author, err = getAuthor(text)
			if err != nil {
				return Commit{}, err
			}
			continue
		}
		if strings.HasPrefix(text, "Date:") {
			commitDate, err = getCommitDate(text)
			if err != nil {
				return Commit{}, err
			}
			continue
		}
		if strings.HasPrefix(text, "diff") {
			break
		}

		// text is commit message
		buff := strings.Trim(text, " ")
		if len(buff) > 1 {
			message = message + ", " + strings.Trim(text, " ")
		}

	}
	if len(message) > 0 {
		message = strings.Replace(message, ", ", "", 1)
	}

	ostrichFileInfo, err := o.parseOstrichFiles(commitTexts)
	if err != nil {
		return Commit{}, err
	}

	return Commit{
		Message:          message,
		Author:           author,
		CommitDate:       commitDate,
		OstrichFileInfos: ostrichFileInfo,
	}, nil
}

func (o *Ostrich) parseOstrichFiles(texts []string) ([]OstrichFileInfo, error) {
	heading := func(texts []string) (int, error) {
		for i, text := range texts {
			if strings.HasPrefix(text, "diff") {
				return i, nil
			}
		}
		return -1, errors.New("can not detect diff heading")
	}
	head, err := heading(texts)
	if err != nil {
		return []OstrichFileInfo{}, err
	}
	result := []OstrichFileInfo{}
	for {
		i, err := heading(texts[head+1:])
		if err != nil {
			o.outputDebug(fmt.Sprintf(
				"block is %d to %d, error is %s\n",
				head,
				len(texts),
				err.Error()))
			ostrichFileInfo, err := o.parseOstrichFile(texts[head:len(texts)])
			if err != nil {
				return []OstrichFileInfo{}, err
			}
			result = append(result, ostrichFileInfo)
			break
		} else {
			o.outputDebug(fmt.Sprintf(
				"block is %d to %d, error is %s\n",
				head,
				head+i,
				err.Error()))
			ostrichFileInfo, err := o.parseOstrichFile(texts[head : head+i])
			if err != nil {
				return []OstrichFileInfo{}, err
			}
			result = append(result, ostrichFileInfo)
			head = head + i + 1
		}
	}

	return result, nil
}

func (o *Ostrich) parseOstrichFile(texts []string) (OstrichFileInfo, error) {
	// 6 is diff, index, -file, +file, @@, diff
	if len(texts) < 6 {
		return OstrichFileInfo{}, fmt.Errorf("invalid ostrich file info texts length %d", len(texts))
	}
	for i := 0; i < 6; i++ {
		o.outputDebug(fmt.Sprintf("ostricch file texts: %s", texts[i]))
	}

	// get filename and infotype
	buff := strings.Split(texts[0], " ")
	filename := buff[2]
	filename = "." + filename[1:len(filename)]

	infoType := OstrichFileInfoTypeModFile
	if strings.HasPrefix(texts[1], "new file mode") {
		infoType = OstrichFileInfoTypeNewFile
	}
	if strings.HasPrefix(texts[1], "deleted file mode") {
		infoType = OstrichFileInfoTypeDelFile
	}
	o.outputDebug(fmt.Sprintf("ostrich file info - filename: %s", filename))
	o.outputDebug(fmt.Sprintf("ostrich file info - info type: %d", infoType))

	if infoType == OstrichFileInfoTypeDelFile {
		return OstrichFileInfo{
			Filename:          filename,
			InfoType:          infoType,
			OstrichMergeInfos: []OstrichMergeInfo{},
		}, nil

	}

	ostrichMergeInfos, err := o.parseOstrichMerges(texts)
	if err != nil {
		return OstrichFileInfo{}, err
	}

	// generate ostrich merge infos
	return OstrichFileInfo{
		Filename:          filename,
		InfoType:          infoType,
		OstrichMergeInfos: ostrichMergeInfos,
	}, nil
}

func (o *Ostrich) parseOstrichMerges(texts []string) ([]OstrichMergeInfo, error) {
	heading := func(texts []string) (int, error) {
		for i, text := range texts {
			if strings.HasPrefix(text, "@@") {
				return i, nil
			}
		}
		return -1, errors.New("can not detect diff heading")
	}
	head, err := heading(texts)
	if err != nil {
		return []OstrichMergeInfo{}, err
	}
	result := []OstrichMergeInfo{}
	for {
		i, err := heading(texts[head+1:])
		if err != nil {
			o.outputDebug(fmt.Sprintf(
				"merge block is %d to %d, error is %s\n",
				head,
				len(texts),
				err.Error()))
			mergeInfos, err := o.parseOstrichMerge(texts[head:len(texts)])
			if err != nil {
				return []OstrichMergeInfo{}, err
			}
			result = append(result, mergeInfos...)
			break
		} else {
			o.outputDebug(fmt.Sprintf(
				"merge block is %d to %d, error is %s\n",
				head,
				head+i,
				err.Error()))
			mergeInfos, err := o.parseOstrichMerge(texts[head : head+i])
			if err != nil {
				return []OstrichMergeInfo{}, err
			}
			result = append(result, mergeInfos...)
			head = head + i + 1
		}
	}
	return result, nil
}

func (o *Ostrich) parseOstrichMerge(texts []string) ([]OstrichMergeInfo, error) {
	o.outputDebug("parseOstrichMerge")
	for _, text := range texts {
		o.outputDebug(fmt.Sprintf("merge texts: %s", text))
	}
	if len(texts) < 2 {
		return []OstrichMergeInfo{}, fmt.Errorf("invalid merge texts length %d", len(texts))
	}

	// getting otrich type, target line range and after text
	getStartLine := func(text string) (int, error) {
		// format: @@ -0,0 +1,9 @@
		buffs := strings.Split(texts[0], " ")
		if len(buffs) < 4 {
			return 0, fmt.Errorf("invalid terms length in merge text.%s", text)
		}
		buff := strings.Replace(buffs[1], "-", "", 1)
		buffs = strings.Split(buff, ",")
		return strconv.Atoi(buffs[0])
	}
	getOstrichType := func(texts []string) OstrichType {
		existsAdd := false
		existsDel := false

		for _, text := range texts {
			if strings.HasPrefix(text, "+") {
				existsAdd = true
				break
			}
		}
		for _, text := range texts {
			if strings.HasPrefix(text, "-") {
				existsDel = true
				break
			}
		}

		if existsAdd && existsDel {
			return OstrichTypeMod
		}

		if existsAdd {
			return OstrichTypeAdd
		}
		return OstrichTypeDel

	}
	getRemoveRange := func(baseLineNumber int, ostrichType OstrichType, texts []string) int {
		if ostrichType == OstrichTypeAdd{
			return baseLineNumber
		}

		startLine := 0
		endLine := 0
		for i, text := range texts {
			if strings.HasPrefix(text, "-") {
				startLine = i
				break
			}
		}
		for i := startLine; i < len(texts); i++ {
			if !strings.HasPrefix(texts[i], "-") {
				endLine = (i - 1)
				break

			}
		}

		return baseLineNumber - (endLine - startLine + 1)
	}

	getAddTexts := func(texts []string) []string {
		result := []string{}
		for _, text := range texts {
			if strings.HasPrefix(text, "+") {
				result = append(result, strings.Replace(text, "+", "", 1))
			}
		}
		return result
	}
	getRemoveTexts := func(texts []string) []string {
		result := []string{}
		for _, text := range texts {
			if strings.HasPrefix(text, "-") {
				result = append(result, strings.Replace(text, "-", "", 1))
			}
		}
		return result
	}
	generateMergeInfo := func(no int, lineNo int, texts []string) OstrichMergeInfo {
		o.outputDebug(fmt.Sprintf("generate merge info %d.source text line no: %d", no, lineNo))
		for _, text := range texts {
			o.outputDebug(fmt.Sprintf("\t%s", text))
		}
		ostrichType := getOstrichType(texts)
		return OstrichMergeInfo {
			no: no,
			ostrichType: ostrichType,
			targetLine: getRemoveRange(lineNo, ostrichType, texts),
			removeTexts: getRemoveTexts(texts),
			afterTexts: getAddTexts(texts),
		}
	}

	results := []OstrichMergeInfo{}
	startLineNo, err := getStartLine(texts[0])
	if err != nil {
		return []OstrichMergeInfo{}, err
	}
	o.outputDebug(fmt.Sprintf("merge start line: %d", startLineNo))
	mergeInfoNo := 0
	sourceTextLineNo := startLineNo
	buffer := []string{}
	for i, text := range texts[1:] {
		o.outputDebug(fmt.Sprintf("merge text %d: %s", i, text))
		if strings.HasPrefix(text, " ") || len(text) <= 0 {
			if len(buffer) != 0 {
				mergeInfoNo++
				mergeInfo := generateMergeInfo(mergeInfoNo, sourceTextLineNo, buffer)
				results = append(results, mergeInfo)
				buffer = buffer[:0]
			}
			sourceTextLineNo++
			continue
		}
		if strings.HasPrefix(text, "-") {
			sourceTextLineNo++
		}
		o.outputDebug("add to buffer")
		buffer = append(buffer, text)
	}
	if len(buffer) != 0{
		mergeInfoNo++
		mergeInfo := generateMergeInfo(mergeInfoNo, sourceTextLineNo, buffer)
		results = append(results, mergeInfo)
	}
	return results ,nil
}

func (o *Ostrich) applyCommit(commit Commit, git GitCommand) error {
	o.outputDebug("applyCommit")
	comment := o.generateOstrichCommentBase(commit)
	for _, ostrichFileInfo := range commit.OstrichFileInfos {
		if err := o.applyOstrichFileInfo(comment, ostrichFileInfo, git); err != nil {
			return err
		}
	}

	return nil
}

func (o *Ostrich) applyOstrichFileInfo(commentBase string, ostrichFileInfo OstrichFileInfo, git GitCommand) error {
	o.outputDebug("applyOstrichFileInfo")
	prefix, err := o.getLineCommentPrefix(ostrichFileInfo.Filename)
	if err != nil {
		return err
	}
	commentBase = prefix + " " + commentBase

	switch ostrichFileInfo.InfoType {
	case OstrichFileInfoTypeNewFile:
		return o.applyCreateOstricFile(ostrichFileInfo, git)
	case OstrichFileInfoTypeModFile:
		return o.applyEditOstricFile(commentBase, prefix, ostrichFileInfo, git)
	case OstrichFileInfoTypeDelFile:
		return o.applyRemoveOstricFile(ostrichFileInfo, git)
	}
	return nil
}

func (o *Ostrich) applyCreateOstricFile(ostrichFileInfo OstrichFileInfo, git GitCommand) error {
	o.outputDebug("applyCreateOstricFile")
	ostrichMergeInfo := ostrichFileInfo.OstrichMergeInfos[0]
	err := o.FileAccessor.WriteAll(ostrichFileInfo.Filename, ostrichMergeInfo.afterTexts)
	if err != nil {
		return err
	}
	if err := git.Add(ostrichFileInfo.Filename); err != nil {
		return err
	}
	return nil
}

func (o *Ostrich) applyEditOstricFile(commentBase string, commentPrefix string, ostrichFileInfo OstrichFileInfo, git GitCommand) error {
	o.outputDebug("applyEditOstricFile")
	contents, err := o.FileAccessor.ReadAll(ostrichFileInfo.Filename)
	if err != nil {
		return err
	}
	sort.Slice(
		ostrichFileInfo.OstrichMergeInfos,
		func(i, j int) bool {
			return ostrichFileInfo.OstrichMergeInfos[i].no > ostrichFileInfo.OstrichMergeInfos[j].no
		})
	for _, mergeInfo := range ostrichFileInfo.OstrichMergeInfos {
		contents, err = o.applyOstrichMergeInfo(commentBase, commentPrefix, contents, mergeInfo)
		if err != nil {
			return err
		}
	}
	if err := o.FileAccessor.WriteAll(ostrichFileInfo.Filename, contents); err != nil {
		return err
	}
	if err := git.Add(ostrichFileInfo.Filename); err != nil {
		return err
	}
	return nil
}
func (o *Ostrich) applyOstrichMergeInfo(commentBase string, commentPrefix string, contents []string, mergeInfo OstrichMergeInfo) ([]string, error) {
	switch mergeInfo.ostrichType {
	case OstrichTypeAdd:
		return o.applyOstrichMergeInfoAdd(commentBase, contents, mergeInfo)
	case OstrichTypeMod:
		return o.applyOstrichMergeInfoMod(commentBase, commentPrefix, contents, mergeInfo)
	case OstrichTypeDel:
		return o.applyOstrichMergeInfoDel(commentBase, commentPrefix, contents, mergeInfo)
	}

	// can not arrived here
	return []string{}, fmt.Errorf("invalid ostrich type %d", mergeInfo.ostrichType)
}

func (o *Ostrich) applyOstrichMergeInfoAdd(commentBase string, contents []string, mergeInfo OstrichMergeInfo) ([]string, error) {
	rangeComments := o.generateOstrichComment(commentBase, "ADD")
	lineIndent := o.getLineIndent(mergeInfo.afterTexts[0])

	firstHalf := contents[:mergeInfo.targetLine-1]
	latterHalf := contents[mergeInfo.targetLine + len(mergeInfo.afterTexts) -1 : len(contents)]
	resultConetnts := []string{}
	resultConetnts = append(resultConetnts, firstHalf...)
	resultConetnts = append(resultConetnts, lineIndent + rangeComments[0])
	resultConetnts = append(resultConetnts, mergeInfo.afterTexts...)
	resultConetnts = append(resultConetnts, lineIndent + rangeComments[1])
	resultConetnts = append(resultConetnts, latterHalf...)

	o.outputDebug("result contents")
	for i, text := range resultConetnts {
		o.outputDebug(fmt.Sprintf("[%d]: %s", i, text))
	}
	return resultConetnts, nil
}

func (o *Ostrich) applyOstrichMergeInfoMod(commentBase string, commentPrefix string, contents []string, mergeInfo OstrichMergeInfo) ([]string, error) {
	o.outputDebug("applyOstrichMergeInfoMod")
	rangeComments := o.generateOstrichComment(commentBase, "MOD")
	lineIndent := o.getLineIndent(mergeInfo.afterTexts[0])

	firstHalf := contents[:mergeInfo.targetLine-1]
	latterHalf := contents[mergeInfo.targetLine - 1 + len(mergeInfo.afterTexts):len(contents)]
	resultConetnts := []string{}
	resultConetnts = append(resultConetnts, firstHalf...)

	resultConetnts = append(resultConetnts, lineIndent + rangeComments[0])
	for _, row := range mergeInfo.removeTexts {
		row = strings.Replace(row, lineIndent, "", 1)
		row = lineIndent + commentPrefix + " " + row
		resultConetnts = append(resultConetnts, row)
	}
	for _, row := range mergeInfo.afterTexts {
		resultConetnts = append(resultConetnts, row)
	}
	resultConetnts = append(resultConetnts, lineIndent + rangeComments[1])
	resultConetnts = append(resultConetnts, latterHalf...)

	o.outputDebug("result contents")
	for i, text := range resultConetnts {
		o.outputDebug(fmt.Sprintf("[%d]: %s", i, text))
	}
	return resultConetnts, nil
}

func (o *Ostrich) applyOstrichMergeInfoDel(commentBase string, commentPrefix string, contents []string, mergeInfo OstrichMergeInfo) ([]string, error) {
	o.outputDebug("applyOstrichMergeInfoDel")
	rangeComments := o.generateOstrichComment(commentBase, "DEL")

	lineIndent := ""
	resultConetnts := []string{}
	latterHalf := []string{}
	if mergeInfo.targetLine >= len(contents){
		resultConetnts = append(resultConetnts, contents...) 
		lineIndent = o.getLineIndent(contents[len(contents) - 1])
	} else {
		lineIndent = o.getLineIndent(mergeInfo.removeTexts[0])
		firstHalf := contents[:mergeInfo.targetLine-1]
		o.outputDebug("front half")
		for i, text := range firstHalf {
			o.outputDebug(fmt.Sprintf("[%d]: %s", i, text))
		}
		resultConetnts = append(resultConetnts, firstHalf...)
		latterHalf = contents[mergeInfo.targetLine-1:len(contents)]
		o.outputDebug("back half")
		for i, text := range latterHalf {
			o.outputDebug(fmt.Sprintf("[%d]: %s", i, text))
		}
	}

	resultConetnts = append(resultConetnts, lineIndent + rangeComments[0])
	for _, row := range mergeInfo.removeTexts {
		row = strings.Replace(row, lineIndent, "", 1)
		row = lineIndent + commentPrefix + " " + row
		resultConetnts = append(resultConetnts, row)
	}
	resultConetnts = append(resultConetnts, lineIndent + rangeComments[1])
	resultConetnts = append(resultConetnts, latterHalf...)

	o.outputDebug("result contents")
	for i, text := range resultConetnts {
		o.outputDebug(fmt.Sprintf("[%d]: %s", i, text))
	}
	return resultConetnts, nil
}

func (o *Ostrich) applyRemoveOstricFile(ostrichFileInfo OstrichFileInfo, git GitCommand) error {
	o.outputDebug("applyRemoveOstricFile")
	if err := o.FileAccessor.RemoveFile(ostrichFileInfo.Filename); err != nil {
		return err
	}
	if err := git.Rm(ostrichFileInfo.Filename); err != nil {
		return err
	}
	return nil
}

func (o *Ostrich) getLineCommentPrefix(filename string) (string, error) {
	ext := filepath.Ext(filename)
	switch ext {
	case ".c", ".go", ".h", ".cpp":
		return "//", nil
	default:
		return "", fmt.Errorf("invalid file ext %s", ext)
	}
}

func (o *Ostrich) generateOstrichCommentBase(commit Commit) string {
	return fmt.Sprintf("%s {OSTRICH_TYPE} %s {RANGE_TAG}",
		commit.CommitDate.Format("2006/01/02"),
		commit.Author)
}

func (o *Ostrich) generateOstrichComment(commentBase string, ostrichTypeText string) []string {
	comment := strings.Replace(commentBase, "{OSTRICH_TYPE}", ostrichTypeText, 1)
	startComment := strings.Replace(comment, "{RANGE_TAG}", "START", 1)
	endComment := strings.Replace(comment, "{RANGE_TAG}", "END", 1)
	return []string{
		startComment,
		endComment,
	}

}

func (o *Ostrich) getLineIndent(line string) string {
	buffRune := []rune(line)
	spaceRune := []rune(" ")
	tabRune := []rune("\t")
	bufflength := len(buffRune)
	indent := ""
	for i := 0; i < bufflength; i++ {
		if buffRune[i] == spaceRune[0] {
			indent = indent + " "
			continue
		}
		if buffRune[i] == tabRune[0] {
			indent = indent + "\t"
			continue
		}
		break
	}
	return indent
}

func (o *Ostrich) chDir(dir string) error {
	o.outputDebug(fmt.Sprintf("chDir: %s", dir))
	return os.Chdir(dir)
}


func (o *Ostrich) checkout(branch string, git GitCommand) error {
	isCurrent, err := o.currentBranchIs(branch, git)
	if err != nil {
		return err
	}
	if isCurrent {
		return nil
	}

	if err := git.Checkout(branch); err != nil {
		return err
	}
	return nil
}

func (o *Ostrich) currentBranchIs(branch string, git GitCommand) (bool, error) {
	branches, err := git.Branch()
	if err != nil {
		return false, err
	}
	for _, branchName := range branches {
		o.outputDebug(fmt.Sprintf("check current branch: %s", branchName))
		if !strings.HasPrefix(branchName, "*") {
			continue
		}
		// no need checkout
		if strings.HasSuffix(branchName, branch) {
			return true, nil
		}
	}
	return false, nil
}

func (o *Ostrich) getRepositoryName(repository string) (string, error) {
	terms := strings.Split(repository, "/")
	if len(terms) <= 0 {
		return "", fmt.Errorf("invalid repository url %s", repository)
	}
	name := strings.Replace(terms[len(terms)-1], ".git", "", 1)
	return name, nil
}


func (o *Ostrich) outputDebug(message string) {
	log.Printf("[DEBUG]: %s", message)
}
func (o *Ostrich) showGitVersion(git GitCommand) {
	outs, _ := git.Version()
	for _, out := range outs {
		o.outputDebug(fmt.Sprintf("version output: %s", out))
	}
}
