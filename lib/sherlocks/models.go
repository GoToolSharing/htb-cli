package sherlocks

type SherlockTask struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	MaskedFlag  string `json:"masked_flag"`
	Hint        string `json:"hint"`
	Completed   bool   `json:"completed"`
}

type SherlockDataTasks struct {
	Tasks []SherlockTask `json:"data"`
}

type SherlockElement struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type SherlockData struct {
	Data []SherlockElement `json:"data"`
}

type SherlockNameID struct {
	Name string
	ID   int
}

type DownloadFile struct {
	URL       string `json:"url"`
	ExpiresIn int    `json:"expires_in"`
}
