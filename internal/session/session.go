package session

type Collection struct {
	Day      string
	Date     string
	Sessions []*Session
}

type Session struct {
	Available  string
	Time       string
	Id         string
	Applicable bool
}

type Subscription struct {
	Name          string
	Remaining     string
	Date          string
	PostRequestId string
}

type SelectedSession struct {
	Day  string
	Date string
	Time string
	Id   string
}
