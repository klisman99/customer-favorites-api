package model

type Product struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Image  string  `json:"image"`
	Price  float64 `json:"price"`
	Rating Rating  `json:"rating"`
}

type Rating struct {
	Rate  float64 `json:"rate"`
	Count int     `json:"count"`
}
