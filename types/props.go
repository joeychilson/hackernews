package types

// FeedProps is the data structure that holds the data for the feed pages.
type FeedProps struct {
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
