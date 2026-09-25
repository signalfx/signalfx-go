package signalfx

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/signalfx/signalfx-go/template"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const templateRecordType = "#/dashify/v1/templates/Record"

const templateRequestBody = `{
  "type": "#/dashify/v1/templates/Record",
  "spec": {
    "<Dashboard>": [],
    "$import:chart0": "/v2/template/HNPr-tr_AAc"
  },
  "title": "Service overview",
  "metadata": {
    "rootElement": "Dashboard",
    "imports": ["/v2/template/HNPr-tr_AAc"]
  }
}`

const templateUpdateRequestBody = `{
  "type": "#/dashify/v1/templates/Record",
  "spec": {
    "<Dashboard>": [],
    "$import:chart0": "/v2/template/HNPr-tr_AAc"
  },
  "title": "Service overview",
  "metadata": {
    "rootElement": "Dashboard",
    "imports": ["/v2/template/HNPr-tr_AAc"]
  },
  "signalview": {
    "lastConvertedAt": null
  }
}`

const templateDatasourceRequestBody = `{
  "type": "#/dashify/v1/templates/Record",
  "spec": {
    "<Chart>": []
  },
  "title": "Request rate",
  "metadata": {
    "rootElement": "Chart",
    "datasource": {
      "type": "SPLUNK_O11Y",
      "programText": "data('http.requests').publish()"
    }
  }
}`

func testTemplateWriteRequest() *template.CreateUpdateTemplateRequest {
	rootElement := template.RootElementDashboard
	return &template.CreateUpdateTemplateRequest{
		Type:  templateRecordType,
		Spec:  json.RawMessage(`{"<Dashboard>":[],"$import:chart0":"/v2/template/HNPr-tr_AAc"}`),
		Title: "Service overview",
		Metadata: template.WriteMetadata{
			RootElement: &rootElement,
			Imports:     []string{"/v2/template/HNPr-tr_AAc"},
		},
	}
}

func TestCreateTemplateWithDatasource(t *testing.T) {
	teardown := setup()
	defer teardown()

	rootElement := template.RootElementChart
	write := &template.CreateUpdateTemplateRequest{
		Type:  templateRecordType,
		Spec:  json.RawMessage(`{"<Chart>":[]}`),
		Title: "Request rate",
		Metadata: template.WriteMetadata{
			RootElement: &rootElement,
			Datasource: &template.Datasource{
				Type:        template.DatasourceTypeSplunkObservability,
				ProgramText: "data('http.requests').publish()",
			},
		},
	}
	mux.HandleFunc(TemplateAPIURL, verifyRequestWithJsonBody(
		t,
		http.MethodPost,
		true,
		http.StatusCreated,
		nil,
		templateDatasourceRequestBody,
		"template/chart_success.json",
	))

	result, err := client.CreateTemplate(context.Background(), write)
	require.NoError(t, err)
	require.NotNil(t, result.Data)
}

func TestTemplateSLODatasourceJSON(t *testing.T) {
	datasource := template.Datasource{
		Type:        template.DatasourceTypeSplunkObservabilitySLO,
		ProgramText: "data('service.level').publish()",
		SLOID:       "example-slo-id",
	}

	payload, err := json.Marshal(datasource)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"type": "SPLUNK_O11Y_SLO",
		"programText": "data('service.level').publish()",
		"sloId": "example-slo-id"
	}`, string(payload))
}

func TestCreateTemplate(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(TemplateAPIURL, verifyRequestWithJsonBody(
		t,
		http.MethodPost,
		true,
		http.StatusCreated,
		nil,
		templateRequestBody,
		"template/dashboard_success.json",
	))

	result, err := client.CreateTemplate(context.Background(), testTemplateWriteRequest())
	require.NoError(t, err)
	require.NotNil(t, result.Data)
	assert.Equal(t, "HNPz_pNAIAE", result.Data.ID)
	assert.Equal(t, "Service overview", result.Data.Title)
	assert.Equal(t, "2026-09-05T00:44:31.408Z[UTC]", result.Data.CreatedAt)
	assert.Equal(t, []string{"/v2/directory/~templates/HNPz_pNAIAE"}, result.Data.DirectoryEntries)
	require.NotNil(t, result.Data.Metadata)
	assert.Empty(t, result.Data.Metadata.Imports)
	assert.JSONEq(t, "null", string(result.Data.SignalView))
	assert.Contains(t, string(result.Data.Spec), "futureLayoutOption")
	assert.Contains(t, string(result.Data.Spec), "/v2/template/HNPr-tr_AAc")
}

func TestGetTemplate(t *testing.T) {
	teardown := setup()
	defer teardown()

	const id = "HNPr-tr_AAc"
	mux.HandleFunc(TemplateAPIURL+"/"+id, createResponse(t, http.StatusOK, "template/chart_success.json", func(t *testing.T, r *http.Request) {
		verifyHeaders(t, r, true)
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Empty(t, r.URL.RawQuery)
	}))

	result, err := client.GetTemplate(context.Background(), id, nil)
	require.NoError(t, err)
	require.NotNil(t, result.Data)
	assert.Equal(t, id, result.Data.ID)
	require.NotNil(t, result.Data.Metadata)
	require.NotNil(t, result.Data.Metadata.RootElement)
	assert.Equal(t, template.RootElementChart, *result.Data.Metadata.RootElement)
	assert.Contains(t, string(result.Data.SignalView), "futureAssociationField")
	assert.Nil(t, result.Errors)
	assert.Nil(t, result.Includes)
	assert.Nil(t, result.Data.CreatedBy)
	require.NotNil(t, result.Data.UpdatedBy)
	assert.Equal(t, "/v2/user/example", *result.Data.UpdatedBy)
}

func TestGetTemplateWithOptions(t *testing.T) {
	teardown := setup()
	defer teardown()

	params := url.Values{"import": []string{string(template.ImportAllChildren)}}
	mux.HandleFunc(
		TemplateAPIURL+"/HNPz_pNAIAE",
		verifyRequest(t, http.MethodGet, true, http.StatusOK, params, "template/dashboard_success.json"),
	)

	result, err := client.GetTemplate(context.Background(), "HNPz_pNAIAE", &template.GetOptions{
		Imports: []template.Import{template.ImportAllChildren},
	})
	require.NoError(t, err)
	require.NotNil(t, result.Data)
	assert.Equal(t, "HNPz_pNAIAE", result.Data.ID)
}

func TestUpdateTemplate(t *testing.T) {
	teardown := setup()
	defer teardown()

	write := testTemplateWriteRequest()
	write.SignalView = json.RawMessage(`{"lastConvertedAt":null}`)
	mux.HandleFunc(TemplateAPIURL+"/HNPz_pNAIAE", verifyRequestWithJsonBody(
		t,
		http.MethodPut,
		true,
		http.StatusOK,
		nil,
		templateUpdateRequestBody,
		"template/dashboard_success.json",
	))

	result, err := client.UpdateTemplate(context.Background(), "HNPz_pNAIAE", write)
	require.NoError(t, err)
	require.NotNil(t, result.Data)
	assert.Equal(t, "HNPz_pNAIAE", result.Data.ID)
}

func TestDeleteTemplate(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(
		TemplateAPIURL+"/HNPz_pNAIAE",
		verifyRequest(t, http.MethodDelete, true, http.StatusNoContent, nil, ""),
	)

	assert.NoError(t, client.DeleteTemplate(context.Background(), "HNPz_pNAIAE"))
}

func TestSearchTemplates(t *testing.T) {
	teardown := setup()
	defer teardown()

	expectedParams := url.Values{
		"search":      []string{"service"},
		"rootElement": []string{"Dashboard", "Chart"},
		"id":          []string{"template-a", "template-b"},
		"title":       []string{"Service overview", "Request rate"},
		"orderBy":     []string{"title", "-updatedAt"},
		"offset":      []string{"3"},
		"size":        []string{"200"},
		"import":      []string{"ALL_CHILDREN"},
	}
	mux.HandleFunc(TemplateAPIURL, createResponse(t, http.StatusOK, "template/search_success.json", func(t *testing.T, r *http.Request) {
		verifyHeaders(t, r, true)
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, expectedParams, r.URL.Query())
	}))

	result, err := client.SearchTemplates(context.Background(), &template.SearchOptions{
		Search:       "service",
		RootElements: []template.RootElement{template.RootElementDashboard, template.RootElementChart},
		IDs:          []string{"template-a", "template-b"},
		Titles:       []string{"Service overview", "Request rate"},
		OrderBy:      []string{"title", "-updatedAt"},
		Offset:       3,
		Size:         200,
		Imports:      []template.Import{template.ImportAllChildren},
	})
	require.NoError(t, err)
	require.NotNil(t, result.Data)
	assert.Equal(t, int64(1), result.Data.Count)
	assert.Nil(t, result.Data.Next)
	assert.Nil(t, result.Data.Prev)
	assert.Equal(t, []string{"/v2/template/HNPz_pNAIAE"}, result.Data.Items)
	require.Len(t, result.Includes, 1)
	assert.Equal(t, "HNPz_pNAIAE", result.Includes[0].ID)
}

func TestSearchTemplatesWithNilOptions(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(TemplateAPIURL, createResponse(t, http.StatusOK, "template/search_success.json", func(t *testing.T, r *http.Request) {
		verifyHeaders(t, r, true)
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Empty(t, r.URL.RawQuery)
	}))

	result, err := client.SearchTemplates(context.Background(), nil)
	require.NoError(t, err)
	require.NotNil(t, result.Data)
}

func TestTemplateEnvelopeAllowsMissingLists(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(
		TemplateAPIURL+"/minimal-id",
		verifyRequest(t, http.MethodGet, true, http.StatusOK, nil, "template/missing_envelope_fields.json"),
	)

	result, err := client.GetTemplate(context.Background(), "minimal-id", nil)
	require.NoError(t, err)
	require.NotNil(t, result.Data)
	assert.Nil(t, result.Errors)
	assert.Nil(t, result.Includes)
	assert.Nil(t, result.Data.Metadata)
}

func TestGetTemplateNotFound(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(
		TemplateAPIURL+"/missing-id",
		verifyRequest(t, http.MethodGet, true, http.StatusNotFound, nil, "template/not_found.json"),
	)

	result, err := client.GetTemplate(context.Background(), "missing-id", nil)
	assert.Nil(t, result)
	require.Error(t, err)
	responseErr, ok := AsResponseError(err)
	require.True(t, ok)
	assert.Equal(t, http.StatusNotFound, responseErr.Code())
	assert.Equal(t, TemplateAPIURL+"/missing-id", responseErr.Route())
	assert.JSONEq(t, fixture("template/not_found.json"), responseErr.Details())
}

func TestUpdateTemplateResponseErrorDoesNotExposeRequest(t *testing.T) {
	teardown := setup()
	defer teardown()

	write := testTemplateWriteRequest()
	write.Title = "private template title"
	mux.HandleFunc(
		TemplateAPIURL+"/HNPz_pNAIAE",
		verifyRequest(t, http.MethodPut, true, http.StatusForbidden, nil, "template/not_found.json"),
	)

	result, err := client.UpdateTemplate(context.Background(), "HNPz_pNAIAE", write)
	assert.Nil(t, result)
	require.Error(t, err)
	responseErr, ok := AsResponseError(err)
	require.True(t, ok)
	assert.Equal(t, http.StatusForbidden, responseErr.Code())
	assert.Equal(t, TemplateAPIURL+"/HNPz_pNAIAE", responseErr.Route())
	assert.JSONEq(t, fixture("template/not_found.json"), responseErr.Details())
	assert.NotContains(t, err.Error(), write.Title)
}

func TestTemplateRejectsInvalidArguments(t *testing.T) {
	c := &Client{}

	result, err := c.CreateTemplate(context.Background(), nil)
	assert.Nil(t, result)
	assert.EqualError(t, err, "template write must not be nil")

	result, err = c.GetTemplate(context.Background(), "", nil)
	assert.Nil(t, result)
	assert.EqualError(t, err, "template ID must not be empty")

	result, err = c.UpdateTemplate(context.Background(), "", testTemplateWriteRequest())
	assert.Nil(t, result)
	assert.EqualError(t, err, "template ID must not be empty")

	result, err = c.UpdateTemplate(context.Background(), "valid-id", nil)
	assert.Nil(t, result)
	assert.EqualError(t, err, "template write must not be nil")

	assert.EqualError(t, c.DeleteTemplate(context.Background(), ""), "template ID must not be empty")

	for _, id := range []string{".", "..", "parent/child"} {
		result, err = c.GetTemplate(context.Background(), id, nil)
		assert.Nil(t, result)
		assert.EqualError(t, err, "template ID must be a single path segment")
	}
}

func TestTemplateOptionParamsOmitZeroValues(t *testing.T) {
	assert.Empty(t, templateGetParams(&template.GetOptions{}))
	assert.Empty(t, templateSearchParams(&template.SearchOptions{}))
	assert.NotNil(t, templateGetParams(nil))
	assert.NotNil(t, templateSearchParams(nil))
}

func TestCreateTemplateRejectsInvalidRawJSON(t *testing.T) {
	c := &Client{}
	write := testTemplateWriteRequest()
	write.Spec = json.RawMessage("{")

	result, err := c.CreateTemplate(context.Background(), write)
	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "marshal template request")
}

func TestTemplateMalformedResponse(t *testing.T) {
	teardown := setup()
	defer teardown()

	mux.HandleFunc(TemplateAPIURL+"/malformed", func(w http.ResponseWriter, r *http.Request) {
		verifyHeaders(t, r, true)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "{")
	})

	result, err := client.GetTemplate(context.Background(), "malformed", nil)
	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "decode template response")
}

func TestTemplateClosesResponseBody(t *testing.T) {
	body := &templateCloseTrackingBody{Reader: strings.NewReader(fixture("template/dashboard_success.json"))}
	httpClient := &http.Client{Transport: templateRoundTripperFunc(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       body,
			Request:    request,
		}, nil
	})}
	c, err := NewClient(TestToken, HTTPClient(httpClient))
	require.NoError(t, err)

	result, err := c.GetTemplate(context.Background(), "valid-id", nil)
	require.NoError(t, err)
	require.NotNil(t, result.Data)
	assert.True(t, body.closed)
}

type templateRoundTripperFunc func(*http.Request) (*http.Response, error)

func (f templateRoundTripperFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

type templateCloseTrackingBody struct {
	io.Reader
	closed bool
}

func (b *templateCloseTrackingBody) Close() error {
	b.closed = true
	return nil
}
