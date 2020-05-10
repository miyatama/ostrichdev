package ostrich

import (
	"time"
)

type Commit struct {
	Message          string
	Author           string
	CommitDate       time.Time
	OstrichFileInfos []OstrichFileInfo
}
