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

// Move this

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Structure pour représenter le JSON entier
type JsonResponse struct {
	Status bool            `json:"status"`
	Data   map[string]Item `json:"data"`
}
