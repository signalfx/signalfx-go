package directory

// EntryType is the schema type returned for Directory entries.
const EntryType = "https://schema.splunkdev.com/dashify/v1/directory/Entry"

// Entry describes one logical Directory path and its Template memberships.
type Entry struct {
	Self      string   `json:"self"`
	Type      string   `json:"type"`
	Path      string   `json:"path"`
	Label     string   `json:"label"`
	Templates []string `json:"templates"`
	Ancestors []string `json:"ancestors"`
	Children  []string `json:"children"`
	Pinned    bool     `json:"pinned"`
	Identity  bool     `json:"identity"`
	Canonical bool     `json:"canonical"`
}

// Patch contains the writable fields for a Directory entry.
//
// Both fields are pointers so callers can distinguish an omitted field from
// an explicit false value or empty list. Templates replaces the complete
// membership list; it is not an atomic add or remove operation. A pointer to
// a nil slice encodes as JSON null, which clears the list just like an empty
// slice does.
type Patch struct {
	Templates *[]string `json:"templates,omitempty"`
	Pinned    *bool     `json:"pinned,omitempty"`
}

// APIError describes an error reported inside a Directory response envelope:
// a status code as a string and a human-readable message.
type APIError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// Result is the response envelope returned by Directory reads and patches.
type Result struct {
	Data     *Entry     `json:"data"`
	Errors   []APIError `json:"errors"`
	Includes []Entry    `json:"includes"`
}
