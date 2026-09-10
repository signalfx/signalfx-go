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
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/signalfx/signalfx-go/directory"
)

// DirectoryAPIURL is the base URL for interacting with Directory entries.
const DirectoryAPIURL = "/v2/directory"

const (
	directoryMaxPathSegments = 10
	directoryMaxPathLength   = 250
)

// GetDirectoryEntry returns the Directory entry for path. Path is a decoded
// logical path such as "~users/user@example.com/My Dashboards". An empty path
// reads the Directory root.
func (c *Client) GetDirectoryEntry(ctx context.Context, path string) (*directory.Result, error) {
	requestPath, err := directoryRequestPath(path, true)
	if err != nil {
		return nil, err
	}

	resp, err := c.doDirectoryRequest(ctx, http.MethodGet, requestPath, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	return decodeDirectoryResult(resp.Body)
}

// PatchDirectoryEntry creates or updates the Directory entry for path.
// Templates, when present, replaces the entry's complete membership list.
// The Directory API provides no atomic add/remove operation or stale-write
// precondition, so callers must account for concurrent updates. Setting Pinned
// to false on an entry with no Templates can make the entry unoccupied.
func (c *Client) PatchDirectoryEntry(ctx context.Context, path string, patch *directory.Patch) (*directory.Result, error) {
	if patch == nil {
		return nil, errors.New("directory patch must not be nil")
	}

	requestPath, err := directoryRequestPath(path, false)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(patch)
	if err != nil {
		return nil, fmt.Errorf("marshal directory patch: %w", err)
	}

	resp, err := c.doDirectoryRequest(ctx, http.MethodPatch, requestPath, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := newResponseError(resp, http.StatusOK, http.StatusCreated); err != nil {
		return nil, err
	}

	return decodeDirectoryResult(resp.Body)
}

// DeleteDirectoryEntry deletes the exact Directory entry at path.
//
// This low-level operation does not verify ownership. Callers should first
// read the entry, verify that the path is owned by them, and confirm that its
// Templates and Children are empty. Deleting an entry does not delete its
// Template resources, but it can remove the entry's membership metadata.
func (c *Client) DeleteDirectoryEntry(ctx context.Context, path string) error {
	requestPath, err := directoryRequestPath(path, false)
	if err != nil {
		return err
	}

	resp, err := c.doDirectoryRequest(ctx, http.MethodDelete, requestPath, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := newResponseError(resp, http.StatusNoContent); err != nil {
		return err
	}

	_, err = io.Copy(io.Discard, resp.Body)
	return err
}

func decodeDirectoryResult(body io.Reader) (*directory.Result, error) {
	payload, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("read directory response: %w", err)
	}

	result := &directory.Result{}
	if err := json.Unmarshal(payload, result); err != nil {
		return nil, fmt.Errorf("decode directory response: %w", err)
	}
	return result, nil
}

// directoryRequestPath validates a decoded logical path and returns the exact
// escaped route expected by the Directory service. It intentionally uses
// form-style escaping per segment because the service treats '+' as a space
// and '%2B' as a literal plus.
func directoryRequestPath(path string, allowRoot bool) (string, error) {
	if path == "" {
		if allowRoot {
			return DirectoryAPIURL, nil
		}
		return "", errors.New("directory root cannot be modified")
	}

	if !utf8.ValidString(path) {
		return "", errors.New("directory path must be valid UTF-8")
	}
	if strings.HasPrefix(path, "/") || strings.HasSuffix(path, "/") {
		return "", fmt.Errorf("invalid directory path %q: leading and trailing slashes are not allowed", path)
	}
	if len(utf16.Encode([]rune("/"+path))) > directoryMaxPathLength {
		return "", fmt.Errorf("invalid directory path %q: path exceeds %d characters", path, directoryMaxPathLength)
	}

	segments := strings.Split(path, "/")
	if len(segments) > directoryMaxPathSegments {
		return "", fmt.Errorf("invalid directory path %q: path exceeds %d segments", path, directoryMaxPathSegments)
	}

	escapedSegments := make([]string, len(segments))
	for i, segment := range segments {
		if err := validateDirectoryPathSegment(segment); err != nil {
			return "", fmt.Errorf("invalid directory path %q: segment %d %w", path, i+1, err)
		}

		escaped := url.QueryEscape(segment)
		escapedSegments[i] = strings.ReplaceAll(escaped, "%2A", "*")
	}

	return DirectoryAPIURL + "/" + strings.Join(escapedSegments, "/"), nil
}

func validateDirectoryPathSegment(segment string) error {
	if segment == "" {
		return errors.New("must not be empty")
	}
	if segment == "." || segment == ".." {
		return fmt.Errorf("%q is reserved", segment)
	}

	runes := []rune(segment)
	if unicode.IsSpace(runes[0]) || unicode.IsSpace(runes[len(runes)-1]) {
		return errors.New("must not start or end with whitespace")
	}
	if strings.Contains(segment[1:], "~") {
		return errors.New("must not contain '~' except as its first character")
	}

	for _, character := range runes {
		if unicode.IsControl(character) {
			return errors.New("must not contain control characters")
		}
		if strings.ContainsRune("%<>:\"\\|?", character) {
			return fmt.Errorf("must not contain %q", character)
		}
	}

	return nil
}

// doDirectoryRequest preserves the Directory service's nonstandard use of
// '+' for spaces in path segments. url.URL.Path cannot represent that spelling
// because '+' is normally a literal character in a URL path.
func (c *Client) doDirectoryRequest(ctx context.Context, method string, requestPath string, body io.Reader) (*http.Response, error) {
	baseURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, err
	}

	baseURL.RawQuery = ""
	baseURL.Fragment = ""
	baseURL.Path = strings.TrimSuffix(baseURL.Path, "/")
	if baseURL.RawPath != "" {
		baseURL.RawPath = strings.TrimSuffix(baseURL.RawPath, "/")
	}

	destination := strings.TrimSuffix(baseURL.String(), "/") + requestPath
	req, err := http.NewRequestWithContext(ctx, method, destination, body)
	if err != nil {
		return nil, err
	}
	if c.authToken != "" {
		req.Header.Set(AuthHeaderKey, c.authToken)
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", "application/json")

	return c.httpClient.Do(req)
}
