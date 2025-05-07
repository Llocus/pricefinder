package types

type Product struct {
	Store       string `json:"store"`
	TotalPages  string `json:"totalPages"`
	Image       string `json:"image"`
	Name        string `json:"name"`
	Link        string `json:"link"`
	Stars       string `json:"stars"`
	StarsQty    string `json:"starsQty"`
	Price       string `json:"price"`
	Saving      string `json:"saving"`
	BeforePrice string `json:"beforePrice"`
}
