package models

type HomeProps struct {
	Stories     []Item
	Total       int
	PerPage     int
	CurrentPage int
	StartPage   int
	TotalPages  int
	PageNumbers []int
}
