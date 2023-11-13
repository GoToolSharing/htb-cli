package vpn

type Assigned struct {
	ID           int    `json:"id"`
	FriendlyName string `json:"friendly_name"`
}

type Data struct {
	Disabled bool     `json:"disabled"`
	Assigned Assigned `json:"assigned"`
}

type Response struct {
	Status bool `json:"status"`
	Data   Data `json:"data"`
}
