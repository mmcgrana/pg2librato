package librato

import (
	"net/url"
)

type Link struct {
	// Rel defines the relationship of the link. A link's relationship must
	// be unique within a single annotation event.
	Rel   string `json:"rel"`
	Href  string `json:"href"`
	Label string `json:"label,omitempty"`
}

// Annotation is used to record external events (e.g. a deployment) that
// typically occur at non-uniform times, yet may impact the behavior of
// monitored metrics.
type Annotation struct {
	Id          int64  `json:"id,omitempty"` // For responses. Do not include when posting.
	Title       string `json:"title"`
	Source      string `json:"source,omitempty"`
	Description string `json:"description,omitempty"`
	Links       []Link `json:"links,omitempty"`
	// StartTime is the unix timestamp indicating the the time at which the event referenced
	// by this annotation started. By default this is set to the current time if not specified.
	StartTime int64 `json:"start_time,omitempty"`
	// Endtime is the unix timestamp indicating the the time at which the event referenced
	// by this annotation ended. For events that have a duration, this is a useful way to
	// annotate the duration of the event. This parameter is optional and defaults to null if not set.
	EndTime int64 `json:"end_time,omitempty"`
}

type AnnotationStream struct {
	Name        string            `json:"name"`
	DisplayName string            `json:"display_name,omitempty"`
	Type        string            `json:"type"`
	Attributes  map[string]string `json:"attributes"` // created_by_ua
}

type AnnotationStreamsResponse struct {
	Query       *QueryResponse      `json:"query"`
	Annotations []*AnnotationStream `json:"annotations"`
}

// PostAnnotation posts an event to an annotation stream returning its ID
func (cli *Client) PostAnnotation(streamName string, ann *Annotation) (int64, error) {
	if err := cli.request("POST", annotationsURL+"/"+streamName, ann, ann); err != nil {
		return 0, err
	}
	return ann.Id, nil
}

// GetAnnotationStreams returns a list of annotation streams (not the annotations themselves)
func (cli *Client) GetAnnotationStreams(name string, page *Pagination) (*AnnotationStreamsResponse, error) {
	params := url.Values{}

	if name != "" {
		params.Set("name", name)
	}

	var ann AnnotationStreamsResponse
	return &ann, cli.request("GET", annotationsURL+"?"+page.toParams(params).Encode(), nil, &ann)
}
