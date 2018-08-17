package client

import (
	"encoding/json"
	"github.com/greenstatic/opensdp/internal/server"
	"github.com/greenstatic/opensdp/internal/services"
)

// Performs the discover request and returns a slice of services the
// client is authorized to access.
func (c *Client) Discover() ([]services.Service, error) {
	// Send request
	data, err := c.Request("discover")
	if err != nil {
		return nil, err
	}

	// Parse response
	dr := server.DiscoverResponse{}
	err = json.Unmarshal(data, &dr)
	if err != nil {
		return nil, err
	}

	// Convert to a services.Service slice
	srvs := make([]services.Service, 0, len(dr.Services))
	for _, drs := range dr.Services {
		srv, err := drs.ToService()
		if err != nil {
			return nil, err
		}

		srvs = append(srvs, srv)
	}

	return srvs, nil
}
