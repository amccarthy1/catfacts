package catfacts

// Breed is a description of a specific cat breed
type Breed struct {
	Breed   string `json:"breed"`
	Country string `json:"country"`
	Origin  string `json:"origin"`
	Coat    string `json:"coat"`
	Pattern string `json:"pattern"`
}

// Fact is an interesting fact about cats
type Fact struct {
	Fact   string `json:"fact"`
	Length int    `json:"length"`
}

type pagination struct {
	Total int `json:"total"`
	// PerPage is inexplicably a string?
	PerPage     string `json:"per_page"`
	CurrentPage int    `json:"current_page"`
	LastPage    int    `json:"last_page"`
	From        int    `json:"from"`
	To          int    `json:"to"`
	// Don't use these URLs, they ignore limit params
	NextPageURL *string `json:"next_page_url"`
	PrevPageURL *string `json:"prev_page_url"`
}

type breeds struct {
	pagination
	Data []Breed `json:"data"`
}

type facts struct {
	pagination
	Data []Fact `json:"data"`
}
