package directory

import "github.com/signalfx/signalfx-go/util"

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
// Pinned is a pointer so callers can distinguish an omitted field from an
// explicit false value. Templates replaces the complete membership list. A
// nil Templates slice is omitted, while a non-nil empty slice clears the list.
type Patch struct {
	Templates []string `json:"templates,omitzero"`
	Pinned    *bool    `json:"pinned,omitempty"`
}

// APIError describes an error reported inside a Directory response envelope.
type APIError = util.APIError

// Result is the response envelope returned by Directory reads and patches.
type Result struct {
	Data     *Entry     `json:"data"`
	Errors   []APIError `json:"errors"`
	Includes []Entry    `json:"includes"`
}
