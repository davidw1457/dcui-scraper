package database

// Search result structs.
type SearchResult struct {
	RecordCount int                 `json:"record_count"` //nolint:tagliatelle
	Records     SearchResultRecords `json:"records"`
	Info        SearchResultsInfo   `json:"info"`
}

type SearchResultRecords struct {
	ComicSeries []SearchResultRecordsComicseries `json:"comicseries"`
}

type SearchResultRecordsComicseries struct {
	Genres       []string `json:"genres"`
	UUID         string   `json:"uuid"`
	Title        string   `json:"title"`
	Slug         string   `json:"slug"`
	IssueCount   int      `json:"issue_count"` //nolint:tagliatelle
	BooksCount   int      `json:"books_count"` //nolint:tagliatelle
	Imprints     []string `json:"imprints"`
	VolumeCount  int      `json:"volume_count"`  //nolint:tagliatelle
	OmnibusCount int      `json:"omnibus_count"` //nolint:tagliatelle
}

type SearchResultsInfo struct {
	ComicSeries SearchResultsInfoComicseries `json:"comicseries"`
}

type SearchResultsInfoComicseries struct {
	CurrentPage      int `json:"current_page"`       //nolint:tagliatelle
	NumPages         int `json:"num_pages"`          //nolint:tagliatelle
	TotalResultCount int `json:"total_result_count"` //nolint:tagliatelle
}

// Series detail structs.
type SeriesDetail struct {
	Description string `json:"description"`
}

// Book detail structs.
type BookDetails struct {
	Page     int                 `json:"page"`
	NumPages int                 `json:"num_pages"` //nolint:tagliatelle
	Values   []BookDetailsValues `json:"values"`
	Total    int                 `json:"total"`
}

type BookDetailsValues struct {
	Tags             []BookDetailsValuesTags `json:"tags"`
	Authors          []Creator               `json:"authors"`
	CoverArtists     []Creator               `json:"cover_artists"` //nolint:tagliatelle
	Pencillers       []Creator               `json:"pencillers"`
	Inkers           []Creator               `json:"inkers"`
	Colorists        []Creator               `json:"colorist"`
	Title            string                  `json:"title"`
	Pages            int                     `json:"pages"`
	PublishDate      string                  `json:"publish_date"` //nolint:tagliatelle
	Slug             string                  `json:"slug"`
	ExclusiveToPlans []string                `json:"exclusive_to_plans"` //nolint:tagliatelle
	UUID             string                  `json:"uuid"`
	Description      string                  `json:"description"`
	PrintRelease     string                  `json:"print_release"` //nolint:tagliatelle
	Publisher        string                  `json:"publisher"`
	Imprint          string                  `json:"imprint"`
	IssueNumber      string                  `json:"issue_number"` //nolint:tagliatelle
}

type BookDetailsValuesTags struct {
	Categories []string `json:"categories"`
	Name       string   `json:"name"`
}

type Creator struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"` //nolint:tagliatelle
}

// Search POST body structs.
type SearchBody struct {
	EngineKey     string            `json:"engine_key"` //nolint:tagliatelle
	Page          int               `json:"page"`
	PerPage       int               `json:"per_page"`       //nolint:tagliatelle
	DocumentTypes []string          `json:"document_types"` //nolint:tagliatelle
	Filters       map[string]string `json:"filters"`
	SortField     map[string]string `json:"sort_field"`     //nolint:tagliatelle
	SortDirection map[string]string `json:"sort_direction"` //nolint:tagliatelle
}
