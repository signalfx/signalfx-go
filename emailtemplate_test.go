package signalfx

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/signalfx/signalfx-go/emailtemplate"
	"github.com/stretchr/testify/assert"
)

const emailTemplateRequestBody = `{
  "name": "Terraform email template",
  "triggerSubject": "{{ruleSeverity}} Alert: {{{ruleName}}}",
  "triggerBody": "Triggered at {{timestamp}}",
  "resolvedSubject": "{{ruleSeverity}} Resolved: {{{ruleName}}}",
  "resolvedBody": "Resolved at {{timestamp}}",
  "to": ["alerts@example.com"],
  "cc": ["oncall@example.com"],
  "bcc": ["audit@example.com"],
  "customHeaders": {
    "X-Test-Header": "terraform"
  }
}`

func testEmailTemplateRequest() *emailtemplate.EmailTemplate {
	return &emailtemplate.EmailTemplate{
		Name:            "Terraform email template",
		TriggerSubject:  "{{ruleSeverity}} Alert: {{{ruleName}}}",
		TriggerBody:     "Triggered at {{timestamp}}",
		ResolvedSubject: "{{ruleSeverity}} Resolved: {{{ruleName}}}",
		ResolvedBody:    "Resolved at {{timestamp}}",
		To:              []string{"alerts@example.com"},
		Cc:              []string{"oncall@example.com"},
		Bcc:             []string{"audit@example.com"},
		CustomHeaders: map[string]string{
			"X-Test-Header": "terraform",
		},
	}
}

func TestCreateEmailTemplate(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/alert/emailtemplate", verifyRequestWithJsonBody(t, "POST", true, http.StatusCreated, nil, emailTemplateRequestBody, "emailtemplate/create_success.json"))

	result, err := client.CreateEmailTemplate(context.Background(), testEmailTemplateRequest())
	assert.NoError(t, err, "Unexpected error creating email template")
	assert.Equal(t, "template-id", result.Id, "ID does not match")
	assert.Equal(t, "Terraform email template", result.Name, "Name does not match")
	assert.Equal(t, []string{"alerts@example.com"}, result.To, "To recipients do not match")
	assert.Equal(t, "terraform", result.CustomHeaders["X-Test-Header"], "Custom header does not match")
}

func TestGetEmailTemplate(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/alert/emailtemplate/template-id", verifyRequest(t, "GET", true, http.StatusOK, nil, "emailtemplate/create_success.json"))

	result, err := client.GetEmailTemplate(context.Background(), "template-id")
	assert.NoError(t, err, "Unexpected error getting email template")
	assert.Equal(t, "template-id", result.Id, "ID does not match")
	assert.Equal(t, "{{ruleSeverity}} Alert: {{{ruleName}}}", result.TriggerSubject, "Trigger subject does not match")
}

func TestUpdateEmailTemplate(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/alert/emailtemplate/template-id", verifyRequestWithJsonBody(t, "PUT", true, http.StatusOK, nil, emailTemplateRequestBody, "emailtemplate/update_success.json"))

	result, err := client.UpdateEmailTemplate(context.Background(), "template-id", testEmailTemplateRequest())
	assert.NoError(t, err, "Unexpected error updating email template")
	assert.Equal(t, "Updated Terraform email template", result.Name, "Name does not match")
	assert.Equal(t, int64(1710000100000), result.UpdatedOnMs, "Updated timestamp does not match")
}

func TestDeleteEmailTemplate(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc("/v2/alert/emailtemplate/template-id", verifyRequest(t, "DELETE", true, http.StatusNoContent, nil, ""))

	err := client.DeleteEmailTemplate(context.Background(), "template-id")
	assert.NoError(t, err, "Unexpected error deleting email template")
}

func TestSearchEmailTemplates(t *testing.T) {
	teardown := setup()
	defer teardown()

	limit := 10
	name := "Terraform"
	offset := 2
	orderBy := "-updatedOnMs"
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("offset", strconv.Itoa(offset))
	params.Add("name", name)
	params.Add("orderBy", orderBy)

	mux.HandleFunc("/v2/alert/emailtemplate", verifyRequest(t, "GET", true, http.StatusOK, params, "emailtemplate/search_success.json"))

	results, err := client.SearchEmailTemplates(context.Background(), limit, name, offset, orderBy)
	assert.NoError(t, err, "Unexpected error searching email templates")
	assert.Equal(t, 1, results.Count, "Count does not match")
	assert.Equal(t, 1, len(results.Results), "Results length does not match")
	assert.Equal(t, "template-id", results.Results[0].Id, "Result ID does not match")
}
