package template

import (
	"encoding/json"

	"github.com/signalfx/signalfx-go/util"
)

// RootElement identifies the kind of document stored in a Template record.
type RootElement string

const (
	// RootElementDashboard identifies a Dashboard template document.
	RootElementDashboard RootElement = "Dashboard"
	// RootElementChart identifies a reusable Chart template document.
	RootElementChart RootElement = "Chart"
)

// Import controls which related Template records the API includes in a
// response.
type Import string

const (
	// ImportAllChildren requests all transitively imported Template records.
	ImportAllChildren Import = "ALL_CHILDREN"
)

// DatasourceType identifies the kind of datasource associated with a Template
// write request.
type DatasourceType string

const (
	// DatasourceTypeSplunkObservability identifies a SignalFlow datasource.
	DatasourceTypeSplunkObservability DatasourceType = "SPLUNK_O11Y"
	// DatasourceTypeSplunkObservabilitySLO identifies a Splunk Observability
	// service-level objective datasource.
	DatasourceTypeSplunkObservabilitySLO DatasourceType = "SPLUNK_O11Y_SLO"
)

// Datasource contains optional datasource metadata accepted when creating or
// replacing a Template. The API does not return this object in Template read
// metadata.
type Datasource struct {
	Type        DatasourceType `json:"type,omitempty"`
	ProgramText string         `json:"programText,omitempty"`
	SLOID       string         `json:"sloId,omitempty"`
}

// WriteMetadata describes a Template document and the metadata extracted from
// it for a create or replace request.
type WriteMetadata struct {
	RootElement *RootElement `json:"rootElement"`
	Imports     []string     `json:"imports,omitempty"`
	Datasource  *Datasource  `json:"datasource,omitempty"`
}

// Metadata describes the read-only metadata returned for a Template document.
type Metadata struct {
	RootElement *RootElement `json:"rootElement"`
	Imports     []string     `json:"imports"`
}

// CreateUpdateTemplateRequest is the body used to create or replace a Template record.
//
// Spec contains the polymorphic Dashify document. SignalView is omitted for a
// nil value and may contain the association update object accepted by the API.
type CreateUpdateTemplateRequest struct {
	// Type is the record type identifier expected by the Template API.
	Type       string          `json:"type"`
	Spec       json.RawMessage `json:"spec"`
	Title      string          `json:"title"`
	Metadata   WriteMetadata   `json:"metadata"`
	SignalView json.RawMessage `json:"signalview,omitempty"`
}

// Template is a Template record returned by the API.
type Template struct {
	ID               string          `json:"id"`
	Self             string          `json:"self"`
	Type             string          `json:"type"`
	Title            string          `json:"title"`
	Spec             json.RawMessage `json:"spec"`
	Metadata         *Metadata       `json:"metadata"`
	SignalView       json.RawMessage `json:"signalview"`
	DirectoryEntries []string        `json:"directoryEntries"`
	CreatedAt        string          `json:"createdAt"`
	CreatedBy        *string         `json:"createdBy"`
	UpdatedAt        string          `json:"updatedAt"`
	UpdatedBy        *string         `json:"updatedBy"`
}

type APIError = util.APIError

// Result is the response envelope for a single Template operation.
type Result struct {
	Data     *Template  `json:"data"`
	Errors   []APIError `json:"errors"`
	Includes []Template `json:"includes"`
}

// SearchPage contains the pagination links, total count, and Template URI
// references returned by a Template search.
type SearchPage struct {
	Type  string   `json:"type"`
	Self  string   `json:"self"`
	Next  *string  `json:"next"`
	Prev  *string  `json:"prev"`
	Count int64    `json:"count"`
	Items []string `json:"items"`
}

// SearchResult is the response envelope for a Template search.
type SearchResult struct {
	Data     *SearchPage `json:"data"`
	Errors   []APIError  `json:"errors"`
	Includes []Template  `json:"includes"`
}
