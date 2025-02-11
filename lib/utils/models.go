package utils

type Machine struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

type Challenge struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

type Username struct {
	ID    int    `json:"id"`
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

// Challenges

type ChallengeFinder struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ChallengeResponseFinder struct {
	ChallengesFinder []ChallengeFinder `json:"challenges"`
}

// Activity

type Activity struct {
	CreatedAt  string `json:"created_at"`
	Date       string `json:"date"`
	DateDiff   string `json:"date_diff"`
	UserID     int    `json:"user_id"`
	UserName   string `json:"user_name"`
	UserAvatar string `json:"user_avatar"`
	Type       string `json:"type"`
}

type InfoActivity struct {
	Activities []Activity `json:"activity"`
}

type DataActivity struct {
	Info InfoActivity `json:"info"`
}
