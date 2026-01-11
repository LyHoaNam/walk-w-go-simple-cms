package model

type PaginationRequest struct {
	NextPage string `query:"next_page" json:"next_page"`
	PrevPage string `query:"prev_page" json:"prev_page"`
	Limit    int    `query:"limit" json:"limit"`
	Order    string `query:"order" json:"order"`
	SortBy   string `query:"sort_by" json:"sort_by"`
}

type CursorPaginatedResponse struct {
	Data        interface{} `json:"data"`
	NextPage    string      `json:"next_page,omitempty"`
	PrevPage    string      `json:"prev_page,omitempty"`
	HasNext     bool        `json:"has_next"`
	HasPrevious bool        `json:"has_previous"`
	Limit       int         `json:"limit"`
}

// DecodedCursor represents the parsed components of a cursor string.
type DecodedCursor struct {
	TimeStamp string
	ID        int64
}
