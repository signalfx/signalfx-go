package signalfx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/signalfx/signalfx-go/template"
)

// TemplateAPIURL is the base URL for interacting with Template API records.
const TemplateAPIURL = "/v2/template"

var (
	errEmptyTemplateID   = errors.New("template ID must not be empty")
	errInvalidTemplateID = errors.New("template ID must be a single path segment")
	errNilTemplateWrite  = errors.New("template write must not be nil")
)

// CreateTemplate creates a Template record.
func (c *Client) CreateTemplate(ctx context.Context, write *template.CreateUpdateTemplateRequest) (*template.Result, error) {
	if write == nil {
		return nil, errNilTemplateWrite
	}

	result := &template.Result{}
	if err := c.executeTemplateRequest(ctx, http.MethodPost, TemplateAPIURL, http.StatusCreated, write, nil, result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetTemplate gets a Template record by ID.
func (c *Client) GetTemplate(ctx context.Context, id string, options *template.GetOptions) (*template.Result, error) {
	path, err := templateAPIPath(id)
	if err != nil {
		return nil, err
	}

	result := &template.Result{}
	if err := c.executeTemplateRequest(ctx, http.MethodGet, path, http.StatusOK, nil, templateGetParams(options), result); err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateTemplate fully replaces a Template record.
func (c *Client) UpdateTemplate(ctx context.Context, id string, write *template.CreateUpdateTemplateRequest) (*template.Result, error) {
	path, err := templateAPIPath(id)
	if err != nil {
		return nil, err
	}
	if write == nil {
		return nil, errNilTemplateWrite
	}

	result := &template.Result{}
	if err := c.executeTemplateRequest(ctx, http.MethodPut, path, http.StatusOK, write, nil, result); err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteTemplate deletes a Template record.
func (c *Client) DeleteTemplate(ctx context.Context, id string) error {
	path, err := templateAPIPath(id)
	if err != nil {
		return err
	}

	return c.executeTemplateRequest(ctx, http.MethodDelete, path, http.StatusNoContent, nil, nil, nil)
}

// SearchTemplates searches Template records.
func (c *Client) SearchTemplates(ctx context.Context, options *template.SearchOptions) (*template.SearchResult, error) {
	result := &template.SearchResult{}
	if err := c.executeTemplateRequest(ctx, http.MethodGet, TemplateAPIURL, http.StatusOK, nil, templateSearchParams(options), result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) executeTemplateRequest(ctx context.Context, method string, path string, expectedStatus int, write *template.CreateUpdateTemplateRequest, params url.Values, result any) error {
	var body io.Reader
	if write != nil {
		payload, err := json.Marshal(write)
		if err != nil {
			return fmt.Errorf("marshal template request: %w", err)
		}
		body = bytes.NewReader(payload)
	}

	resp, err := c.doRequest(ctx, method, path, params, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := newResponseError(resp, expectedStatus); err != nil {
		return err
	}

	if result == nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		return fmt.Errorf("decode template response: %w", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)

	return nil
}

func templateAPIPath(id string) (string, error) {
	if id == "" {
		return "", errEmptyTemplateID
	}
	if id == "." || id == ".." || strings.Contains(id, "/") {
		return "", errInvalidTemplateID
	}
	return TemplateAPIURL + "/" + id, nil
}

func templateGetParams(options *template.GetOptions) url.Values {
	if options == nil {
		return url.Values{}
	}

	params := url.Values{}
	addTemplateImports(params, options.Imports)
	return params
}

func templateSearchParams(options *template.SearchOptions) url.Values {
	if options == nil {
		return url.Values{}
	}

	params := url.Values{}
	if options.Search != "" {
		params.Set("search", options.Search)
	}
	for _, rootElement := range options.RootElements {
		params.Add("rootElement", string(rootElement))
	}
	for _, id := range options.IDs {
		params.Add("id", id)
	}
	for _, title := range options.Titles {
		params.Add("title", title)
	}
	for _, orderBy := range options.OrderBy {
		params.Add("orderBy", orderBy)
	}
	if options.Offset != 0 {
		params.Set("offset", strconv.Itoa(options.Offset))
	}
	if options.Size != 0 {
		params.Set("size", strconv.Itoa(options.Size))
	}
	addTemplateImports(params, options.Imports)

	return params
}

func addTemplateImports(params url.Values, imports []template.Import) {
	for _, include := range imports {
		params.Add("import", string(include))
	}
}
