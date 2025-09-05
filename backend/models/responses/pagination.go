package responses

// PaginatedResponse is a generic structure for paginated API responses.
// It is used for documentation purposes with swaggo.
// The actual response comes from the github.com/morkid/paginate library.
type PaginatedResponse struct {
	Items      []interface{} `json:"items"`
	Total      int64         `json:"total"`
	Page       int64         `json:"page"`
	PerPage    int64         `json:"per_page"`
	TotalPages int64         `json:"total_pages"`
}
