package web


type WebRequest struct {
	Action WebAction
	Info OstrichWebRequest
}

type WebAction int

const (
	WebRequestActionOstrich WebAction = iota
	WebRequestActionDone
)
