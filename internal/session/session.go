package session

type Collection struct {
	Day      string
	Date     string
	Sessions []*Session
}

type Session struct {
	Available  string
	Limit      string
	Time       string
	Applicable bool
}

type Subscription struct {
	Name          string
	Remaining     string
	Date          string
	PostRequestId string
}
