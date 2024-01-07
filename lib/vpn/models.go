package vpn

type Assigned struct {
	ID               int    `json:"id"`
	FriendlyName     string `json:"friendly_name"`
	CurrentClients   int    `json:"current_clients"`
	Location         string `json:"location"`
	LocationFriendly string `json:"location_type_friendly"`
}

type Data struct {
	Disabled bool     `json:"disabled"`
	Assigned Assigned `json:"assigned"`
}

type Response struct {
	Status bool `json:"status"`
	Data   Data `json:"data"`
}
