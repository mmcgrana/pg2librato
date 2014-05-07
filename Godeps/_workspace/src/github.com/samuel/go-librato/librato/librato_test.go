package librato

import (
	"os"
	"testing"
)

func testClient(t *testing.T) *Client {
	username := os.Getenv("LIBRATO_TEST_USERNAME")
	token := os.Getenv("LIBRATO_TEST_TOKEN")
	if username == "" || token == "" {
		t.Skip("LIBRATO_TEST_USERNAME or LIBRATO_TEST_TOKEN unset")
	}
	return &Client{
		Username: username,
		Token:    token,
	}
}
