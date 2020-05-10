package web

type OstrichWebRequest struct {
	Repository    string `json:"repository"`
	FromBranch    string `json:"fromBranch"`
	CommitID      string `json:"commitId"`
	OstrichBranch string `json:"ostrichBranch"`
}
