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

// Fortreses

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Structure pour représenter le JSON entier
type JsonResponse struct {
	Status bool            `json:"status"`
	Data   map[string]Item `json:"data"`
}

// Endgames
type Endgame struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type EndgameJsonResponse struct {
	Status bool      `json:"status"`
	Data   []Endgame `json:"data"`
}

// Prolabs

type Lab struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Structure pour représenter la section 'data' du JSON
type Data struct {
	Labs []Lab `json:"labs"`
}

// Structure pour représenter le JSON entier
type ProlabJsonResponse struct {
	Status bool `json:"status"`
	Data   Data `json:"data"`
}
