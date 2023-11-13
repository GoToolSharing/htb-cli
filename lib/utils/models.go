package utils

type Machine struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type Challenge struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type Username struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type Root struct {
	Machines   interface{} `json:"machines"`
	Challenges interface{} `json:"challenges"`
	Usernames  interface{} `json:"users"`
}
