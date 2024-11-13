package database

// Search result structs.
type SearchResult struct {
	Record_count int //nolint:revive,stylecheck
	Records      SearchResultRecords
	Info         SearchResultsInfo
}

type SearchResultRecords struct {
	Comicseries []SearchResultRecordsComicseries
}

type SearchResultRecordsComicseries struct {
	Genres        []string
	UUID          string
	Title         string
	Slug          string
	Issue_count   int //nolint:revive,stylecheck
	Books_count   int //nolint:revive,stylecheck
	Imprints      []string
	Volume_count  int //nolint:revive,stylecheck
	Omnibus_count int //nolint:revive,stylecheck
}

type SearchResultsInfo struct {
	Comicseries SearchResultsInfoComicseries
}

type SearchResultsInfoComicseries struct {
	Current_page       int //nolint:revive,stylecheck
	Num_pages          int //nolint:revive,stylecheck
	Total_result_count int //nolint:revive,stylecheck
}

// Series detail structs.
type SeriesDetail struct {
	Description string
}

// Book detail structs.
type BookDetails struct {
	Page      int
	Num_pages int //nolint:revive,stylecheck
	Values    []BookDetailsValues
	Total     int
}

type BookDetailsValues struct {
	Tags               []BookDetailsValuesTags
	Authors            []Creator
	Cover_artists      []Creator //nolint:revive,stylecheck
	Pencillers         []Creator
	Inkers             []Creator
	Colorists          []Creator
	Title              string
	Pages              int
	Publish_date       string //nolint:revive,stylecheck
	Slug               string
	Exclusive_to_plans []string //nolint:revive,stylecheck
	UUID               string
	Description        string
	Print_release      string //nolint:revive,stylecheck
	Publisher          string
	Imprint            string
	Issue_number       string //nolint:revive,stylecheck
}

type BookDetailsValuesTags struct {
	Categories []string
	Name       string
}

type Creator struct {
	Name         string
	Display_name string //nolint:revive,stylecheck
}
