package inbox

type Inbox struct {
	Total        int    `json:"total"`
	NbElements   int    `json:"nbElements"`
	Premier      int    `json:"premier"`
	Tri          string `json:"tri"`
	PageCourante int    `json:"pageCourante"`
}
