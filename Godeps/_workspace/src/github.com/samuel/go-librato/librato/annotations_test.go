package librato

import (
	"testing"
	"time"
)

func TestPostAnnotation(t *testing.T) {
	cli := testClient(t)
	now := time.Now().UTC().Unix()
	ann := &Annotation{
		Title:       "Test event",
		Source:      "test",
		Description: "This event is just for testing",
		Links: []Link{{
			Rel:   "Home",
			Href:  "https://github.com/samuel/go-librato",
			Label: "Source to the package",
		}},
		StartTime: now - 60*5,
		EndTime:   now,
	}
	if id, err := cli.PostAnnotation("test_annotation", ann); err != nil {
		t.Fatal(err)
	} else if id == 0 {
		t.Fatal("Id is 0")
	} else {
		t.Logf("Annotation ID: %d", id)
		t.Logf("\t%+v", ann)
	}
}

func TestGetAnnotationStreams(t *testing.T) {
	cli := testClient(t)
	res, err := cli.GetAnnotationStreams("", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Query: %+v", res.Query)
	for _, ann := range res.Annotations {
		t.Logf("\tAnnotation: %+v", ann)
	}
}
