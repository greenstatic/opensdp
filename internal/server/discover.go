package server

import (
	"encoding/json"
	"github.com/greenstatic/opensdp/internal/services"
	"net"
	"net/http"
)

type DiscoverResponseService struct {
	Name       string     `json:"name"`
	IP         string     `json:"ip"`
	Ports      [][]string `json:"ports"`
	Tags       []string   `json:"tags"`
	AccessType []string   `json:"accessType"`
}

type DiscoverResponse struct {
	Success  bool                      `json:"success"`
	DeviceId string                    `json:"deviceId"`
	Services []DiscoverResponseService `json:"services"`
}

// Fills a DiscoverResponseService struct from a services.Service struct.
func (drs *DiscoverResponseService) Create(service services.Service) error {
	// Name field
	drs.Name = service.Name

	// IP field
	drs.IP = service.IP.String()

	// Ports field
	ports := make([][]string, 0, len(service.ProtoPort))
	for _, portCombo := range service.ProtoPort {
		ports = append(ports, portCombo.StringSlice())
	}
	drs.Ports = ports

	// Tags field
	drs.Tags = service.Tags

	// Access Type field
	accessTypes := make([]string, 0, len(service.AccessType))
	for _, at := range service.AccessType {
		accessTypes = append(accessTypes, at.String())
	}
	drs.AccessType = accessTypes

	return nil
}

// Returns a services.Service struct from that data in the DiscoverResponseServices.
func (drs *DiscoverResponseService) ToService() (services.Service, error) {
	s := services.Service{}
	// Name field
	s.Name = drs.Name

	// IP field
	s.IP = net.ParseIP(drs.IP)

	// Ports field
	ports := make([]services.ProtoPort, 0, len(drs.Ports))
	for _, portCombo := range drs.Ports {
		pp := services.ProtoPort{}
		err := pp.FromStringSlice(portCombo)
		if err != nil {
			return services.Service{}, err
		}

		ports = append(ports, pp)
	}
	s.ProtoPort = ports

	// Tags field
	s.Tags = drs.Tags

	// Access Type field
	accessTypes := make([]services.AccessType, 0, len(drs.AccessType))
	for _, at := range drs.AccessType {
		var a services.AccessType
		err := a.FromString(at)

		if err != nil {
			return services.Service{}, err
		}
		accessTypes = append(accessTypes, a)
	}
	s.AccessType = accessTypes

	return s, nil
}

// Wrapper handler for the discover endpoint. The wrapper allows us to
// inject the clients slice into the handler function.
func (s *Server) discoverResponseWrapper() func(w http.ResponseWriter, req *http.Request) {

	return func(w http.ResponseWriter, req *http.Request) {
		cn := req.TLS.PeerCertificates[0].Subject.CommonName

		client, ok := s.Clients[cn]

		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(struct {
				Successful bool   `json:"success"`
				Error      string `json:"error"`
			}{
				false,
				"not authorized for any services",
			})
			return
		}

		cServices := make([]DiscoverResponseService, 0, len(client.Services))
		for _, srv := range client.Services {
			drs := DiscoverResponseService{}
			drs.Create(srv.Service)
			cServices = append(cServices, drs)
		}

		json.NewEncoder(w).Encode(DiscoverResponse{true, cn, cServices})
	}

}
