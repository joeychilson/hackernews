package types

// HomeProps is the data structure that holds the data for the home page.
type HomeProps struct {
	Stories     []Item
	Total       int
	PerPage     int
	CurrentPage int
	StartPage   int
	TotalPages  int
	PageNumbers []int
}

// ItemProps is the data structure that holds the data for the item page.
type ItemProps struct {
	Item     Item
	Comments []Item
}
