package ostrich

type OstrichFileInfo struct {
	Filename          string
	InfoType          OstrichFileInfoType
	OstrichMergeInfos []OstrichMergeInfo
}

type OstrichFileInfoType int

const (
	OstrichFileInfoTypeNewFile OstrichFileInfoType = iota
	OstrichFileInfoTypeModFile
	OstrichFileInfoTypeDelFile
)
