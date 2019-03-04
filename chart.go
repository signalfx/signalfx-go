package signalfx

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// ChartAPIURL is the base URL for interacting with charts.
const ChartAPIURL = "/v2/chart"

// Chart is a chart.
type Chart struct {
	Created          int64  `json:"created,omitempty"`
	Creator          string `json:"creator,omitempty"`
	CustomProperties string `json:"customProperties,omitempty"`
	Description      string `json:"description,omitempty"`
	LastUpdated      int64  `json:"lastUpdated,omitEmpty"`
	LastUpdatedBy    string `json:"lastUpdatedBy,omitEmpty"`
	Name             string `json:"name"`
	Options          struct {
		AreaChartOptions struct {
			ShowDataMarkers bool `json:"showDataMarkers,omitempty"`
		} `json:"areaChartOptions,omitempty"`
		Axes []struct {
			HighWaterMark      int64  `json:"highWaterMark,omitempty"`
			HighWaterMarkLabel string `json:"highWaterMarkLabel,omitempty"`
			Label              string `json:"label,omitempty"`
			LowWaterMarkLabel  string `json:"lowWaterMarkLabel,omitempty"`
			LowhWaterMark      int64  `json:"lowhWaterMark,omitempty"`
			Max                int64  `json:"max,omitempty"`
			Min                int64  `json:"min,omitempty"`
		} `json:"axes,omitempty"`
		AxisPrecision int64  `json:"axisPrecision,omitempty"`
		ColorBy       string `json:"colorBy,omitempty"`
		ColorRange    struct {
			Color string `json:"color,omitempty"`
			Max   string `json:"max,omitempty"`
			Min   string `json:"min,omitempty"`
		} `json:"colorRange,omitempty"`
		ColorScale struct {
			Inverted   bool    `json:"inverted,omitempty"`
			Thresholds []int64 `json:"thresholds,omitempty"`
		} `json:"colorScale,omitempty"`
		ColorScale2 []struct {
			Gt           int64 `json:"gt,omitempty"`
			Gte          int64 `json:"gte,omitempty"`
			Lt           int64 `json:"lt,omitempty"`
			Lte          int64 `json:"lte,omitempty"`
			PaletteIndex int64 `json:"paletteIndex,omitempty"`
		} `json:"colorScale2,omitempty"`
		DefaultPlotType       string   `json:"defaultPlotType,omitempty"`
		GroupBy               []string `json:"groupBy,omitempty"`
		HistogramChartOptions struct {
			ColorThemeIndex int64 `json:"colorThemeIndex,omitempty"`
		} `json:"histogramChartOptions,omitempty"`
		IncludeZero   bool `json:"includeZero,omitempty"`
		LegendOptions struct {
			Fields []struct {
				Enabled  bool   `json:"enabled,omitempty"`
				Property string `json:"property,omitempty"`
			} `json:"fields,omitempty"`
		} `json:"legendOptions,omitempty"`
		LineChartOptionsObject struct {
			ShowDataMarkers bool `json:"showDataMarkers,omitempty"`
		} `json:"lineChartOptionsObject,omitempty"`
		Markdown             string `json:"markdown,omitempty"`
		MaximumPrecision     int64  `json:"maximumPrecision,omitempty"`
		OnChartLegendOptions struct {
			DimensionInLegend string `json:"dimensionInLegend,omitempty"`
			ShowLegend        bool   `json:"showLegend,omitempty"`
		} `json:"onChartLegendOptions,omitempty"`
		ProgramOptions struct {
			DisableSampling   bool  `json:"disableSampling,omitempty"`
			MaxDelay          int64 `json:"maxDelay,omitempty"`
			MinimumResolution int64 `json:"minimumResolution,omitempty"`
		} `json:"programOptions,omitempty"`
		PublishLabelOptions []struct {
			DisplayName  string `json:"displayName,omitempty"`
			Label        string `json:"label,omitempty"`
			PaletteIndex int64  `json:"paletteIndex,omitempty"`
			PlotType     string `json:"plotType,omitempty"`
			ValuePrefix  string `json:"valuePrefix,omitempty"`
			ValueSuffix  string `json:"valueSuffix,omitempty"`
			ValueUnit    string `json:"valueUnit,omitempty"`
			YAxis        int64  `json:"yAxis,omitempty"`
		} `json:"publishLabelOptions,omitempty"`
		RefreshInterval        int64  `json:"refreshInterval,omitempty"`
		SecondaryVisualization string `json:"secondaryVisualization,omitempty"`
		ShowEventLines         bool   `json:"showEventLines,omitempty"`
		ShowSparkLine          bool   `json:"showSparkLine,omitempty"`
		SortBy                 string `json:"sortBy,omitempty"`
		SortDirection          string `json:"sortDirection,omitempty"`
		SortProperty           string `json:"sortProperty,omitempty"`
		Stacked                bool   `json:"stacked,omitempty"`
		Time                   struct {
			End   int64  `json:"end,omitempty"`
			Range int64  `json:"range,omitempty"`
			Start int64  `json:"start,omitempty"`
			Type  string `json:"type,omitempty"`
		} `json:"time,omitempty"`
		TimeStampHidden bool   `json:"timeStampHidden,omitempty"`
		Type            string `json:"type,omitempty"`
		UnitPrefix      string `json:"unitPrefix,omitempty"`
	} `json:"options,omitempty"`
	PackageSpecifications string   `json:"packageSpecifications,omitempty"`
	ProgramText           string   `json:"programText,omitempty"`
	Tags                  []string `json:"tags,omitempty"`
}

// ChartSearch is the result of a query for Charts
type ChartSearch struct {
	Count   int64 `json:"count,omitempty"`
	Results []Chart
}

// CreateChart creates a chart.
func (c *Client) CreateChart(chart *Chart) (*Chart, error) {
	payload, err := json.Marshal(chart)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest("POST", ChartAPIURL, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	finalChart := &Chart{}

	err = json.NewDecoder(resp.Body).Decode(finalChart)

	return finalChart, err
}

// DeleteChart deletes a chart.
func (c *Client) DeleteChart(id string) error {
	resp, err := c.doRequest("DELETE", ChartAPIURL+"/"+id, nil, nil)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("Unexpected status code: " + resp.Status)
	}

	return nil
}

// GetChart gets a chart.
func (c *Client) GetChart(id string) (*Chart, error) {
	resp, err := c.doRequest("GET", ChartAPIURL+"/"+id, nil, nil)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	finalChart := &Chart{}

	err = json.NewDecoder(resp.Body).Decode(finalChart)

	return finalChart, err
}

// UpdateChart updates a chart.
func (c *Client) UpdateChart(id string, chart *Chart) (*Chart, error) {
	payload, err := json.Marshal(chart)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest("PUT", ChartAPIURL+"/"+id, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	finalChart := &Chart{}

	err = json.NewDecoder(resp.Body).Decode(finalChart)

	return finalChart, err
}

// SearchChart searches for charts, given a query string in `name`.
func (c *Client) SearchChart(limit int, name string, offset int, tags string) (*ChartSearch, error) {
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("name", name)
	params.Add("offset", strconv.Itoa(offset))
	params.Add("tags", tags)

	resp, err := c.doRequest("GET", ChartAPIURL, params, nil)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	finalCharts := &ChartSearch{}

	err = json.NewDecoder(resp.Body).Decode(finalCharts)

	return finalCharts, err
}
