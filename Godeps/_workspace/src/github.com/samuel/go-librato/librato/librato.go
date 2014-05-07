/*
Package librato is a library for Librato's metrics service API.

Example:

	metrics = &librato.Client{"login@email.com", "token"}
	metrics := &librato.Metrics{
		Counters: []librato.Metric{librato.Metric{"name", 123, "source"}},
		Gauges: []librato.Gauge{},
	}
	metrics.SendMetrics(metrics)
*/
package librato

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const (
	userAgent = "go-librato/1.0"
)

const (
	annotationsURL = "https://metrics-api.librato.com/v1/annotations"
	metricsURL     = "https://metrics-api.librato.com/v1/metrics"
	servicesURL    = "https://metrics-api.librato.com/v1/services"
	usersURL       = "https://metrics-api.librato.com/v1/users"
)

type Sort string

const (
	Ascending  Sort = "asc"
	Descending Sort = "desc"
)

type QueryResponse struct {
	// Length is the maximum number of resources to return in the response.
	Length int `json:"length"`
	// Offset is the index into the entire result set at which the current
	// response begins. E.g. if a total of 20 resources match the query,
	// and the offset is 5, the response begins with the sixth resource.
	Offset int `json:"offset"`
	// Total is the total number of resources owned by the user.
	Total int `json:"total"`
	// Found is the number of resources owned by the user that satisfy
	// the specified query parameters. found will be less than or equal
	// to total. Additionally if length is less than found, the response
	// is a subset of the resources matching the specified query parameters.
	Found int `json:"found"`
}

// Pagination sets the pagination options on get requests
type Pagination struct {
	// Offset specifies how many results to skip for the first
	// returned result. Defaults to 0.
	Offset int
	// Length specifies how many resources should be returned. The
	// maximum permissible (and the default) length is 100.
	Length int
	// Orderby the specified attribute. Permissible set of orderby
	// attributes and the default value varies with resource type.
	Orderby string
	// Sort is the order in which the results should be ordered. Permissible
	// values are asc (ascending) and desc (descending). Defaults to asc.
	Sort Sort
}

type TimeInterval struct {
	// StarTime is the unix timestamp indicating the start time of the desired interval.
	StartTime int64 `json:"start_time,omitempty"`
	// EndTime is the unix timestamp indicating the end time of the desired
	// interval. If left unspecified it defaults to the current time.
	EndTime int64 `json:"end_time,omitempty"`
	// Count is the number of measurements desired. When specified as
	// N in conjunction with StartTime, the response contains the first N
	// measurements after StartTime. When specified as N in conjunction with
	// EndTime, the response contains the last N measurements before EndTime.
	Count int `json:"count,omitempty"`
	// A resolution for the response as measured in seconds. If the original
	// measurements were reported at a higher resolution than specified in
	// the request, the response contains averaged measurements.
	Resolution int `json:"resolution,omitempty"`
}

type ErrTypes struct {
	Params  map[string]interface{} `json:"params"`
	Request []string               `json:"request"`
	System  []string               `json:"system"`
}

type ErrResponse struct {
	StatusCode int
	Errors     ErrTypes `json:"errors"`
}

func (e *ErrResponse) Error() string {
	return fmt.Sprintf("librato: error %d: %+v", e.StatusCode, e.Errors)
}

type Client struct {
	Username string
	Token    string
}

func (cli *Client) request(method string, url string, req, res interface{}) error {
	if method == "GET" && req != nil {
		errors.New("librato: req must be nil for GET requests")
	}

	var body io.Reader
	if req != nil {
		buf := &bytes.Buffer{}
		body = buf
		if err := json.NewEncoder(buf).Encode(req); err != nil {
			return err
		}
	}

	httpReq, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	if httpReq != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	httpReq.Header.Set("User-Agent", userAgent)
	httpReq.SetBasicAuth(cli.Username, cli.Token)
	httpRes, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode >= 400 {
		errRes := &ErrResponse{StatusCode: httpRes.StatusCode}
		if err := json.NewDecoder(httpRes.Body).Decode(errRes); err != nil {
			return err
		}
		return errRes
	}

	if res != nil {
		if err := json.NewDecoder(httpRes.Body).Decode(res); err != nil {
			return err
		}
	}

	return nil
}

func (page *Pagination) toParams(params url.Values) url.Values {
	if params == nil {
		params = url.Values{}
	}
	if page == nil {
		return params
	}
	if page.Offset > 0 {
		params.Set("offset", strconv.Itoa(page.Offset))
	}
	if page.Length > 0 {
		params.Set("length", strconv.Itoa(page.Length))
	}
	if page.Orderby != "" {
		params.Set("orderby", page.Orderby)
	}
	if page.Sort != "" {
		params.Set("sort", string(page.Sort))
	}
	return params
}
