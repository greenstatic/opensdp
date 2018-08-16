package server

import (
	"net/http"
	"encoding/json"
)

type DiscoverResponseServices struct {
	Name string `json:"name"`
	IPs []string `json:"ips"`
	Ports [][]string `json:"ports"`
	Tags []string `json:"tags"`
	AccessType []string `json:"accessType"`
}

type DiscoverResponse struct {
	Success bool `json:"success"`
	DeviceId string `json:"deviceId"`
	Services []DiscoverResponseServices `json:"services"`
}


func (s *Server) discoverResponseWrapper() (func (w http.ResponseWriter, req *http.Request)) {

	return func (w http.ResponseWriter, req *http.Request) {
		cn := req.TLS.PeerCertificates[0].Subject.CommonName

		client, ok := s.Clients[cn]

		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(struct{
				Successful bool `json:"success"`
				Error string `json:"error"`
			}{
				false,
				"not authorized for any services",
			})
			return
		}

		cServices := make([]DiscoverResponseServices, 0, len(client.Services))
		for _, srv := range client.Services {
			drs := DiscoverResponseServices{}
			// Name field
			drs.Name = srv.Service.Name

			// IP field
			ips := make([]string, 0, len(srv.Service.Ips))
			for _, ip := range srv.Service.Ips {
				ips = append(ips, ip.String())
			}
			drs.IPs = ips

			// Ports field
			ports := make([][]string, 0, len(srv.Service.ProtoPort))
			for _, portCombo := range srv.Service.ProtoPort {
				ports = append(ports, portCombo.StringSlice())
			}
			drs.Ports = ports

			// Tags field
			drs.Tags = srv.Service.Tags

			// Access Type field
			accessTypes := make([]string, 0, len(srv.Service.AccessType))
			for _, at := range srv.Service.AccessType {
				accessTypes = append(accessTypes, at.String())
			}
			drs.AccessType = accessTypes

			cServices = append(cServices, drs)
		}

		json.NewEncoder(w).Encode(DiscoverResponse{true, cn, cServices})
	}

}

