package librato

type Service struct {
	ID       int               `json:"id"`
	Type     string            `json:"type"`
	Settings map[string]string `json:"settings"`
	Title    string            `json:"title"`
}

type ServicesResponse struct {
	Query    *QueryResponse `json:"query"`
	Services []*Service     `json:"service"`
}

// GetServices returns all services created by the user.
// http://dev.librato.com/v1/get/services
func (cli *Client) GetServices(page *Pagination) (*ServicesResponse, error) {
	var svc ServicesResponse
	return &svc, cli.request("GET", servicesURL+"?"+page.toParams(nil).Encode(), nil, &svc)
}
