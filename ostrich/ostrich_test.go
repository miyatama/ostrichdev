package ostrich

import (
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

type DummyFileAcccessor struct{}

func (d *DummyFileAcccessor) ReadAll(filepath string) ([]string, error) {
	return []string{}, nil
}

func (d *DummyFileAcccessor) WriteAll(filepath string, contents []string) error {
	return nil
}

func (d *DummyFileAcccessor) RemoveFile(filepath string) error {
	return nil
}
func TestParseCommit(t *testing.T) {
	ostrich := Ostrich{
		Repository:    "",
		FromBranch:    "",
		OstrichBranch: "",
		CommitId:      "",
		FileAccessor:  &DummyFileAcccessor{},
	}
	t.Run("empty string", func(t *testing.T) {
		commitTexts := []string{}
		_, err := ostrich.parseCommit(commitTexts)
		if err == nil {
			t.Fatal("not return error")
		}
	})
	t.Run("add file commit", func(t *testing.T) {
		b, err := ioutil.ReadFile("../testdata/add_file_commit_text.txt")
		if err != nil {
			t.Fatalf("can not read test data.%#v", err)
		}
		text := string(b)
		commitTexts := strings.Split(text, "\n")
		commit, err := ostrich.parseCommit(commitTexts)
		if err != nil {
			t.Fatalf("reterned error %#v", err)

		}
		if commit.Author != "unknown" {
			t.Fatalf("invalid author %s", commit.Author)
		}
		expectDate, _ := time.Parse("2006-01-02", "2020-03-31")
		if commit.CommitDate.Equal(expectDate) {
			t.Fatalf("invalid commit date %s", commit.CommitDate.Format(time.RFC3339))
		}
		if commit.Message != "add main.go" {
			t.Fatalf("invalid commit message %s", commit.Message)
		}
		if len(commit.OstrichFileInfos) != 1 {
			t.Fatalf(
				"invalid ostrich file infos.ostrich file info length is %d",
				len(commit.OstrichFileInfos))

		}
		if commit.OstrichFileInfos[0].Filename != "./main.go" {
			t.Fatalf("invalid filename %s", commit.OstrichFileInfos[0].Filename)
		}

		if commit.OstrichFileInfos[0].InfoType != OstrichFileInfoTypeNewFile {
			t.Fatalf("invalid info type.expect %d, result %d",
				OstrichFileInfoTypeNewFile,
				commit.OstrichFileInfos[0].InfoType,
			)
		}

		expectMergeInfosLength := 1
		if len(commit.OstrichFileInfos[0].OstrichMergeInfos) != expectMergeInfosLength {
			t.Fatalf(
				"invalid ostrich merge info length.expect: %d, result: %d",
				expectMergeInfosLength,
				len(commit.OstrichFileInfos[0].OstrichMergeInfos))
		}

		// when new file mode then valid info is after texts only
		ostrichMergeInfo := commit.OstrichFileInfos[0].OstrichMergeInfos[0]
		expectTexts := []string{
			"package main",
			"",
			"import (",
			"       \"fmt\"",
			")",
			"",
			"func main() {",
			"       fmt.Println(\"hello\")",
			"}",
		}
		if len(ostrichMergeInfo.afterTexts) != len(expectTexts) {
			t.Fatalf(
				"invalid after texts length.expect %d, result %d",
				len(expectTexts),
				len(ostrichMergeInfo.afterTexts),
			)
		}
	})
	t.Run("remove file commit", func(t *testing.T) {
		b, err := ioutil.ReadFile("../testdata/remove_file_commit_text.txt")
		if err != nil {
			t.Fatal("can not read test data")
		}
		text := string(b)
		commitTexts := strings.Split(text, "\n")
		commit, err := ostrich.parseCommit(commitTexts)
		if err != nil {
			t.Fatalf("reterned error %#v", err)

		}
		if commit.Message != "remove miyata.txt" {
			t.Fatalf("invalid commit message %s", commit.Message)
		}
		if len(commit.OstrichFileInfos) != 1 {
			t.Fatalf(
				"invalid ostrich file infos.ostrich file info length is %d",
				len(commit.OstrichFileInfos))

		}

		if commit.OstrichFileInfos[0].Filename != "./miyata.txt" {
			t.Fatalf("invalid filename %s", commit.OstrichFileInfos[0].Filename)
		}

		if commit.OstrichFileInfos[0].InfoType != OstrichFileInfoTypeDelFile {
			t.Fatalf("invalid info type.expect %d, result %d",
				OstrichFileInfoTypeDelFile,
				commit.OstrichFileInfos[0].InfoType,
			)
		}

	})
	t.Run("modify file commit", func(t *testing.T) {
		b, err := ioutil.ReadFile("../testdata/mod_file_commit_text.txt")
		if err != nil {
			t.Fatal("can not read test data")
		}
		text := string(b)
		commitTexts := strings.Split(text, "\n")
		commit, err := ostrich.parseCommit(commitTexts)
		if err != nil {
			t.Fatalf("reterned error %#v", err)

		}
		if commit.Message != "mod print message" {
			t.Fatalf("invalid commit message %s", commit.Message)
		}
		if len(commit.OstrichFileInfos) != 1 {
			t.Fatalf("invalid ostrich file infos.ostrich file info length is %d", len(commit.OstrichFileInfos))

		}

		if commit.OstrichFileInfos[0].Filename != "./main.go" {
			t.Fatalf("invalid filename.expect %s, result %s", "./main.go", commit.OstrichFileInfos[0].Filename)
		}

		if commit.OstrichFileInfos[0].InfoType != OstrichFileInfoTypeModFile {
			t.Fatalf("invalid info type.expect %d, result %d",
				OstrichFileInfoTypeNewFile,
				commit.OstrichFileInfos[0].InfoType,
			)
		}
		expectMergetInfoLength := 1
		if len(commit.OstrichFileInfos[0].OstrichMergeInfos) != expectMergetInfoLength {
			t.Fatalf(
				"invalid ostrich merge info length.expect: %d, result: %d",
				expectMergetInfoLength,
				len(commit.OstrichFileInfos[0].OstrichMergeInfos))
		}

		ostrichMergeInfo := commit.OstrichFileInfos[0].OstrichMergeInfos[0]
		expectTexts := []string{
			"       fmt.Println(\"hello world\")",
		}
		if len(ostrichMergeInfo.afterTexts) != len(expectTexts) {
			t.Fatalf(
				"invalid after texts length.expect %d, result %d",
				len(expectTexts),
				len(ostrichMergeInfo.afterTexts),
			)
		}
		for i, text := range expectTexts {
			if text != ostrichMergeInfo.afterTexts[i] {
				t.Fatalf(
					"invalid after text %d, expect %s, result %s",
					i,
					text,
					ostrichMergeInfo.afterTexts[i],
				)
			}
		}
		expectRemoveTexts := []string{
			"       fmt.Println(\"hello\")",
		}
		if len(ostrichMergeInfo.removeTexts) != len(expectRemoveTexts) {
			t.Fatalf(
				"invalid remove texts length.expect %d, result %d",
				len(expectRemoveTexts),
				len(ostrichMergeInfo.removeTexts),
			)
		}
		for i, text := range expectRemoveTexts {
			if text != ostrichMergeInfo.removeTexts[i] {
				t.Fatalf(
					"invalid remove text %d, expect %s, result %s",
					i,
					text,
					ostrichMergeInfo.removeTexts[i],
				)
			}
		}

		if ostrichMergeInfo.ostrichType != OstrichTypeMod {
			t.Fatalf(
				"invalid ostrich type.expect %d, result %d",
				OstrichTypeMod,
				ostrichMergeInfo.ostrichType)
		}
		expectTargetLineRange := 8
		if expectTargetLineRange != ostrichMergeInfo.targetLine {
			t.Fatalf(
				"invalid target line , expect: %d, result: %d",
				expectTargetLineRange,
				ostrichMergeInfo.targetLine)
		}
	})
	t.Run("edit two point add", func(t *testing.T) {
		b, err := ioutil.ReadFile("../testdata/mod_file_commit_text_add_twe_parts.txt")
		if err != nil {
			t.Fatal("can not read test data")
		}
		text := string(b)
		commitTexts := strings.Split(text, "\n")
		commit, err := ostrich.parseCommit(commitTexts)
		if err != nil {
			t.Fatalf("reterned error %#v", err)

		}
		if commit.Message != "mod print message" {
			t.Fatalf("invalid commit message %s", commit.Message)
		}
		if len(commit.OstrichFileInfos) != 1 {
			t.Fatalf("invalid ostrich file infos.ostrich file info length is %d", len(commit.OstrichFileInfos))

		}

		if commit.OstrichFileInfos[0].Filename != "./main.go" {
			t.Fatalf("invalid filename.expect %s, result %s", "./main.go", commit.OstrichFileInfos[0].Filename)
		}

		if commit.OstrichFileInfos[0].InfoType != OstrichFileInfoTypeModFile {
			t.Fatalf("invalid info type.expect %d, result %d",
				OstrichFileInfoTypeNewFile,
				commit.OstrichFileInfos[0].InfoType,
			)
		}
		expectMergeInfosLength := 2
		if len(commit.OstrichFileInfos[0].OstrichMergeInfos) != expectMergeInfosLength {
			t.Fatalf(
				"invalid ostrich merge info length.expect: %d, result: %d",
				expectMergeInfosLength,
				len(commit.OstrichFileInfos[0].OstrichMergeInfos))
		}
		expectTexts := [][]string{
			[]string{
				"     test2()",
			},
			[]string{
				"",
				"func test2() {",
				"     fmt.Println(\"this is test\")",
				"}",
			},
		}
		expectTargetLines := []int{
			10,
			15,
		}
		expectRemoveTexts := [][]string{
			[]string{},
			[]string{},
		}

		for i, ostrichMergeInfo := range commit.OstrichFileInfos[0].OstrichMergeInfos{
			expectText := expectTexts[i]
			if len(ostrichMergeInfo.afterTexts) != len(expectText) {
				t.Fatalf(
					"invalid after texts length.expect: %d, result: %d",
					len(expectText),
					len(ostrichMergeInfo.afterTexts),
				)
			}
			for i, text := range expectText {
				if text != ostrichMergeInfo.afterTexts[i] {
					t.Fatalf(
						"invalid after text %d, expect: %s, result: %s",
						i,
						text,
						ostrichMergeInfo.afterTexts[i],
					)
				}
			}
			removeText := expectRemoveTexts[i]
			if len(ostrichMergeInfo.removeTexts) != len(removeText) {
				t.Fatalf(
					"invalid remove texts length.expect: %d, result: %d",
					len(removeText),
					len(ostrichMergeInfo.removeTexts),
				)
			}

			if ostrichMergeInfo.ostrichType != OstrichTypeAdd {
				t.Fatalf(
					"invalid ostrich type.expect: %d, result: %d",
					OstrichTypeAdd,
					ostrichMergeInfo.ostrichType)
			}
			expectTargetLine := expectTargetLines[i]
			if expectTargetLine != ostrichMergeInfo.targetLine {
				t.Fatalf(
					"invalid target line, expect: %d, result: %d",
					expectTargetLine,
					ostrichMergeInfo.targetLine)
			}
		}

	})
	t.Run("modify file commit remove two line", func(t *testing.T) {
		b, err := ioutil.ReadFile("../testdata/mod_file_commit_text_remove_multi_line.txt")

		if err != nil {
			t.Fatal("can not read test data")
		}
		text := string(b)
		commitTexts := strings.Split(text, "\n")
		commit, err := ostrich.parseCommit(commitTexts)
		if err != nil {
			t.Fatalf("reterned error %#v", err)

		}
		if commit.Message != "mod print message" {
			t.Fatalf("invalid commit message %s", commit.Message)
		}
		if len(commit.OstrichFileInfos) != 1 {
			t.Fatalf("invalid ostrich file infos.ostrich file info length is %d", len(commit.OstrichFileInfos))

		}

		if commit.OstrichFileInfos[0].Filename != "./main.go" {
			t.Fatalf("invalid filename.expect %s, result %s", "./main.go", commit.OstrichFileInfos[0].Filename)
		}

		if commit.OstrichFileInfos[0].InfoType != OstrichFileInfoTypeModFile {
			t.Fatalf("invalid info type.expect %d, result %d",
				OstrichFileInfoTypeNewFile,
				commit.OstrichFileInfos[0].InfoType,
			)
		}
		expectMergetInfoLength := 1
		if len(commit.OstrichFileInfos[0].OstrichMergeInfos) != expectMergetInfoLength {
			t.Fatalf(
				"invalid ostrich merge info length.expect: %d, result: %d",
				expectMergetInfoLength,
				len(commit.OstrichFileInfos[0].OstrichMergeInfos))
		}

		ostrichMergeInfo := commit.OstrichFileInfos[0].OstrichMergeInfos[0]
		expectTexts := []string{
			"       fmt.Println(\"hello world\")",
		}
		if len(ostrichMergeInfo.afterTexts) != len(expectTexts) {
			t.Fatalf(
				"invalid after texts length.expect %d, result %d",
				len(expectTexts),
				len(ostrichMergeInfo.afterTexts),
			)
		}
		for i, text := range expectTexts {
			if text != ostrichMergeInfo.afterTexts[i] {
				t.Fatalf(
					"invalid after text %d, expect %s, result %s",
					i,
					text,
					ostrichMergeInfo.afterTexts[i],
				)
			}
		}
		expectRemoveTexts := []string{
			"       fmt.Println(\"hello\")",
			"       fmt.Println(\"hello2\")",
			"       fmt.Println(\"hello3\")",
		}
		if len(ostrichMergeInfo.removeTexts) != len(expectRemoveTexts) {
			t.Fatalf(
				"invalid remove texts length.expect %d, result %d",
				len(expectRemoveTexts),
				len(ostrichMergeInfo.removeTexts),
			)
		}
		for i, text := range expectRemoveTexts {
			if text != ostrichMergeInfo.removeTexts[i] {
				t.Fatalf(
					"invalid remove text %d, expect %s, result %s",
					i,
					text,
					ostrichMergeInfo.removeTexts[i],
				)
			}
		}
		if ostrichMergeInfo.ostrichType != OstrichTypeMod {
			t.Fatalf(
				"invalid ostrich type.expect %d, result %d",
				OstrichTypeMod,
				ostrichMergeInfo.ostrichType)
		}
		expectTargetLines := 8
		if expectTargetLines != ostrichMergeInfo.targetLine {
			t.Fatalf(
			"invalid target line, expect: %d, result: %d",
				expectTargetLines,
				ostrichMergeInfo.targetLine)
		}
	})
}

func TestGetLineCommentPrefix(t *testing.T) {
	ostrich := Ostrich{
		Repository:    "",
		FromBranch:    "",
		OstrichBranch: "",
		CommitId:      "",
		FileAccessor:  &DummyFileAcccessor{},
	}
	t.Run("valid file type", func(t *testing.T) {
		filenames := []string{
			"./main.c",
			"./main.h",
			"./main.cpp",
			"./main.go",
		}
		expectResults := []string{
			"//",
			"//",
			"//",
			"//",
		}

		for i, filename := range filenames {
			result, err := ostrich.getLineCommentPrefix(filename)
			if err != nil {
				t.Fatalf("invalid return error %s", err.Error())
			}
			if expectResults[i] != result {
				t.Fatalf("invalid return.expect: %s, result: %s", expectResults[i], result)
			}
		}
	})
	t.Run("invalid file type", func(t *testing.T) {
		filename := "miyata"
		_, err := ostrich.getLineCommentPrefix(filename)
		if err == nil {
			t.Fatalf("invalid return error.error is nil")
		}
	})
}

func TestApplyOstrichMergeInfoAdd(t *testing.T) {
	ostrich := Ostrich{
		Repository:    "",
		FromBranch:    "",
		OstrichBranch: "",
		CommitId:      "",
		FileAccessor:  &DummyFileAcccessor{},
	}
	t.Run("Add First line", func(t *testing.T) {
		comment := "// 2020/04/18 {OSTRICH_TYPE} miyatama {RANGE_TAG}"
		contents := []string{
			"add text 01",
			"add text 02",
			"add text 03",
			"row 001",
			"row 002",
			"row 003",
		}
		mergeInfo := OstrichMergeInfo{
			no:              1,
			ostrichType:     OstrichTypeAdd,
			targetLine: 1,
			afterTexts: []string{
				"add text 01",
				"add text 02",
				"add text 03",
			},
		}
		resultContents, err := ostrich.applyOstrichMergeInfoAdd(comment, contents, mergeInfo)
		if err != nil {
			t.Fatalf("return error %#v", err)
		}
		expectContents := []string{
			"// 2020/04/18 ADD miyatama START",
			"add text 01",
			"add text 02",
			"add text 03",
			"// 2020/04/18 ADD miyatama END",
			"row 001",
			"row 002",
			"row 003",
		}
		if len(expectContents) != len(resultContents) {
			t.Fatalf("invalid result contents row length.expect %d, result %d.", len(expectContents), len(resultContents))
		}
		for i, expectRow := range expectContents {
			if expectRow != resultContents[i] {
				t.Fatalf("invalid result contents %d row.expect %s, result %s.",
					i,
					expectRow,
					resultContents[i])
			}
		}

	})
	t.Run("Add second line", func(t *testing.T) {
		comment := "// 2020/04/18 {OSTRICH_TYPE} miyatama {RANGE_TAG}"
		contents := []string{
			"    row 001",
			"    row new 01",
			"    row new 02",
			"    row 002",
			"    row 003",
		}
		mergeInfo := OstrichMergeInfo{
			no:              1,
			ostrichType:     OstrichTypeAdd,
			targetLine: 2,
			afterTexts: []string{
				"    row new 01",
				"    row new 02",
			},
		}
		resultContents, err := ostrich.applyOstrichMergeInfoAdd(comment, contents, mergeInfo)
		if err != nil {
			t.Fatalf("return error %#v", err)
		}
		expectContents := []string{
			"    row 001",
			"    // 2020/04/18 ADD miyatama START",
			"    row new 01",
			"    row new 02",
			"    // 2020/04/18 ADD miyatama END",
			"    row 002",
			"    row 003",
		}
		if len(expectContents) != len(resultContents) {
			t.Fatalf("invalid result contents row length.expect %d, result %d.", len(expectContents), len(resultContents))
		}
		for i, expectRow := range expectContents {
			if expectRow != resultContents[i] {
				t.Fatalf("invalid result contents %d row.expect %s, result %s.",
					i,
					expectRow,
					resultContents[i])
			}
		}

	})
}

func TestApplyOstrichMergeInfoMod(t *testing.T) {
	ostrich := Ostrich{
		Repository:    "",
		FromBranch:    "",
		OstrichBranch: "",
		CommitId:      "",
		FileAccessor:  &DummyFileAcccessor{},
	}
	t.Run("remove first row", func(t *testing.T) {

		comment := "// 2020/04/18 {OSTRICH_TYPE} miyatama {RANGE_TAG}"
		commentPrefix := "//"
		contents := []string{
			"add text 01",
			"add text 02",
			"add text 03",
			"row 001",
			"row 002",
			"row 003",
		}
		mergeInfo := OstrichMergeInfo{
			no:              1,
			ostrichType:     OstrichTypeAdd,
			targetLine: 1,
			afterTexts: []string{
				"add text 01",
				"add text 02",
				"add text 03",
			},
			removeTexts: []string{
				"remove row 01",
				"remove row 02",
			},
		}
		resultContents, err := ostrich.applyOstrichMergeInfoMod(comment, commentPrefix, contents, mergeInfo)
		if err != nil {
			t.Fatalf("invalid return error %#v", err)
		}
		expectContents := []string{
			"// 2020/04/18 MOD miyatama START",
			"// remove row 01",
			"// remove row 02",
			"add text 01",
			"add text 02",
			"add text 03",
			"// 2020/04/18 MOD miyatama END",
			"row 001",
			"row 002",
			"row 003",
		}
		if len(expectContents) != len(resultContents) {
			t.Fatalf("invalid result contents row length.expect %d, result %d.", len(expectContents), len(resultContents))
		}
		for i, expectRow := range expectContents {
			if expectRow != resultContents[i] {
				t.Fatalf("invalid result contents %d row.expect %s, result %s.",
					i,
					expectRow,
					resultContents[i])
			}
		}
	})
	t.Run("remove last row", func(t *testing.T) {

		comment := "// 2020/04/18 {OSTRICH_TYPE} miyatama {RANGE_TAG}"
		commentPrefix := "//"
		contents := []string{
			"row 001",
			"row 002",
			"row 003",
			"add text 01",
			"add text 02",
			"add text 03",
		}
		mergeInfo := OstrichMergeInfo{
			no:              1,
			ostrichType:     OstrichTypeAdd,
			targetLine: 4,
			afterTexts: []string{
				"add text 01",
				"add text 02",
				"add text 03",
			},
			removeTexts: []string{
				"remove row 01",
				"remove row 02",
			},
		}
		resultContents, err := ostrich.applyOstrichMergeInfoMod(comment, commentPrefix, contents, mergeInfo)
		if err != nil {
			t.Fatalf("invalid return error %#v", err)
		}
		expectContents := []string{
			"row 001",
			"row 002",
			"row 003",
			"// 2020/04/18 MOD miyatama START",
			"// remove row 01",
			"// remove row 02",
			"add text 01",
			"add text 02",
			"add text 03",
			"// 2020/04/18 MOD miyatama END",
		}
		if len(expectContents) != len(resultContents) {
			t.Fatalf("invalid result contents row length.expect %d, result %d.", len(expectContents), len(resultContents))
		}
		for i, expectRow := range expectContents {
			if expectRow != resultContents[i] {
				t.Fatalf("invalid result contents %d row.expect %s, result %s.",
					i,
					expectRow,
					resultContents[i])
			}
		}
	})
}

func TestApplyOstrichMergeInfoDel(t *testing.T) {
	ostrich := Ostrich{
		Repository:    "",
		FromBranch:    "",
		OstrichBranch: "",
		CommitId:      "",
		FileAccessor:  &DummyFileAcccessor{},
	}
	t.Run("remove first row", func(t *testing.T) {

		comment := "// 2020/04/18 {OSTRICH_TYPE} miyatama {RANGE_TAG}"
		commentPrefix := "//"
		contents := []string{
			"row 001",
			"row 002",
			"row 003",
		}
		mergeInfo := OstrichMergeInfo{
			no:              1,
			ostrichType:     OstrichTypeAdd,
			targetLine: 1,
			afterTexts:      []string{},
			removeTexts:     []string{
				"remove row1",
				"remove row2",
			},
		}
		resultContents, err := ostrich.applyOstrichMergeInfoDel(comment, commentPrefix, contents, mergeInfo)
		if err != nil {
			t.Fatalf("invalid return error %#v", err)
		}
		expectContents := []string{
			"// 2020/04/18 DEL miyatama START",
			"// remove row1",
			"// remove row2",
			"// 2020/04/18 DEL miyatama END",
			"row 001",
			"row 002",
			"row 003",
		}
		if len(expectContents) != len(resultContents) {
			t.Fatalf("invalid result contents row length.expect %d, result %d.", len(expectContents), len(resultContents))
		}
		for i, expectRow := range expectContents {
			if expectRow != resultContents[i] {
				t.Fatalf("invalid result contents %d row.expect %s, result %s.",
					i,
					expectRow,
					resultContents[i])
			}
		}
	})
	t.Run("remove last row", func(t *testing.T) {

		comment := "// 2020/04/18 {OSTRICH_TYPE} miyatama {RANGE_TAG}"
		commentPrefix := "//"
		contents := []string{
			"row 001",
			"row 002",
			"row 003",
		}
		mergeInfo := OstrichMergeInfo{
			no:              1,
			ostrichType:     OstrichTypeAdd,
			targetLine: 4,
			afterTexts:      []string{},
			removeTexts:     []string{
				"remove row1",
				"remove row2",
			},
		}
		resultContents, err := ostrich.applyOstrichMergeInfoDel(comment, commentPrefix, contents, mergeInfo)
		if err != nil {
			t.Fatalf("invalid return error %#v", err)
		}
		expectContents := []string{
			"row 001",
			"row 002",
			"row 003",
			"// 2020/04/18 DEL miyatama START",
			"// remove row1",
			"// remove row2",
			"// 2020/04/18 DEL miyatama END",
		}
		if len(expectContents) != len(resultContents) {
			t.Fatalf("invalid result contents row length.expect %d, result %d.", len(expectContents), len(resultContents))
		}
		for i, expectRow := range expectContents {
			if expectRow != resultContents[i] {
				t.Fatalf("invalid result contents %d row.expect %s, result %s.",
					i,
					expectRow,
					resultContents[i])
			}
		}
	})
}
