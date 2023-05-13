package models

type Movie struct {
	Tconst         string `json:"tconst"`
	TitleType      string `json:"titleType,omitempty"`
	PrimaryTitle   string `json:"primaryTitle"`
	RuntimeMinutes int    `json:"runtimeMinutes"`
	Genres         string `json:"genres"`
}

type GenreSubtotal struct {
	Genre        string `json:"genre"`
	Tconst       string `json:"tconst,omitempty"`
	PrimaryTitle string `json:"primaryTitle"`
	NumVotes     int    `json:"numVotes,omitempty"`
	Subtotal     int    `json:"subtotal"`
}

type InsertMovie struct {
	Tconst         string  `json:"tconst"`
	TitleType      string  `json:"titleType,omitempty"`
	PrimaryTitle   string  `json:"primaryTitle"`
	RuntimeMinutes int     `json:"runtimeMinutes,omitempty"`
	Genres         string  `json:"genres"`
	AverageRating  float64 `json:"averateRating"`
	NumVotes       int     `json:"numVotes,omitempty"`
}

type PaginatedMovies struct {
	Movies          []InsertMovie `json:"movies"`
	TotalRecords    int           `json:"total_records"`
	TotalPages      int           `json:"total_pages"`
	CurrentPage     int           `json:"current_page"`
	NextPage        int           `json:"next_page"`
	PreviousPage    int           `json:"previous_page"`
	Status          string        `json:"status"`
	NextPageURL     string        `json:"next_page_url,omitempty"`
	PreviousPageURL string        `json:"previous_page_url,omitempty"`
}

type PaginatedGenre struct {
	Genres          []GenreSubtotal `json:"genre"`
	TotalRecords    int             `json:"total_records"`
	TotalPages      int             `json:"total_pages"`
	CurrentPage     int             `json:"current_page"`
	NextPage        int             `json:"next_page"`
	PreviousPage    int             `json:"previous_page"`
	Status          string          `json:"status"`
	NextPageURL     string          `json:"next_page_url,omitempty"`
	PreviousPageURL string          `json:"previous_page_url,omitempty"`
}
