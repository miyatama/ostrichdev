package ostrich

type OstrichMergeInfo struct {
	no              int
	ostrichType     OstrichType
	targetLine      int      // edit start line 
	removeTexts     []string // remove or modified texts
	afterTexts      []string // add or modify texts
}

type OstrichType int

const (
	OstrichTypeAdd OstrichType = iota
	OstrichTypeMod
	OstrichTypeDel
)
