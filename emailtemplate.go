package signalfx

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/signalfx/signalfx-go/emailtemplate"
)

// EmailTemplateAPIURL is the base URL for interacting with email templates.
const EmailTemplateAPIURL = "/v2/alert/emailtemplate"

// CreateEmailTemplate creates an email template.
func (c *Client) CreateEmailTemplate(ctx context.Context, template *emailtemplate.EmailTemplate) (*emailtemplate.EmailTemplate, error) {
	payload, err := json.Marshal(template)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, http.MethodPost, EmailTemplateAPIURL, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusCreated); err != nil {
		return nil, err
	}

	finalTemplate := &emailtemplate.EmailTemplate{}
	err = json.NewDecoder(resp.Body).Decode(finalTemplate)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalTemplate, err
}

// GetEmailTemplate gets an email template by ID.
func (c *Client) GetEmailTemplate(ctx context.Context, id string) (*emailtemplate.EmailTemplate, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, EmailTemplateAPIURL+"/"+id, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalTemplate := &emailtemplate.EmailTemplate{}
	err = json.NewDecoder(resp.Body).Decode(finalTemplate)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalTemplate, err
}

// UpdateEmailTemplate updates an email template.
func (c *Client) UpdateEmailTemplate(ctx context.Context, id string, template *emailtemplate.EmailTemplate) (*emailtemplate.EmailTemplate, error) {
	payload, err := json.Marshal(template)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, http.MethodPut, EmailTemplateAPIURL+"/"+id, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalTemplate := &emailtemplate.EmailTemplate{}
	err = json.NewDecoder(resp.Body).Decode(finalTemplate)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalTemplate, err
}

// DeleteEmailTemplate deletes an email template.
func (c *Client) DeleteEmailTemplate(ctx context.Context, id string) error {
	resp, err := c.doRequest(ctx, http.MethodDelete, EmailTemplateAPIURL+"/"+id, nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusNoContent); err != nil {
		return err
	}
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return nil
}

// SearchEmailTemplates searches email templates by name with pagination.
func (c *Client) SearchEmailTemplates(ctx context.Context, limit int, name string, offset int, orderBy string) (*emailtemplate.SearchResult, error) {
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("offset", strconv.Itoa(offset))
	if name != "" {
		params.Add("name", name)
	}
	if orderBy != "" {
		params.Add("orderBy", orderBy)
	}

	resp, err := c.doRequest(ctx, http.MethodGet, EmailTemplateAPIURL, params, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalTemplates := &emailtemplate.SearchResult{}
	err = json.NewDecoder(resp.Body).Decode(finalTemplates)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalTemplates, err
}
