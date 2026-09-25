package template

// GetOptions controls related records included with a Template read.
type GetOptions struct {
	Imports []Import
}

// SearchOptions contains optional Template search filters and pagination.
// Slice fields are encoded as repeated query parameters in their given order.
type SearchOptions struct {
	Search       string
	RootElements []RootElement
	IDs          []string
	Titles       []string
	OrderBy      []string
	Offset       int
	Size         int
	Imports      []Import
}
