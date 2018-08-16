package configsyaml

import (
	"github.com/greenstatic/opensdp/internal/clients"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"errors"
	"github.com/satori/go.uuid"
	"github.com/greenstatic/opensdp/internal/services"
)

type clientFileServicePolicy struct {
	Name string
}

type clientFile struct {
	DeviceId string `yaml:"deviceId"`
	Label string
	Services []clientFileServicePolicy
}

type clientsFile struct {
	Version string
	Kind string
	Clients []clientFile
}

	func ClientsRead(path string, serv []services.Service) (map[string]clients.Client, error) {
	m := make(map[string]clients.Client)

	// Read file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse file
	cf := clientsFile{}
	err = yaml.Unmarshal(data, &cf)
	if err != nil {
		return nil, err
	}

	// Check if correct kind
	if cf.Kind != "clients" {
		return nil, errors.New("file not kind clients")
	}

	// Parse clients
	for _, c := range cf.Clients {
		clnt, err := parseClient(c, serv)
		if err != nil {
			return nil, err
		}
		m[clnt.DeviceId] = clnt
	}

	return m, nil
}

// Parses a clientFile struct into a clients.Client with resolved services
func parseClient(c clientFile, serv []services.Service) (clients.Client, error) {
	clnt := clients.Client{}

	// Parse deviceId
	if c.DeviceId == "" {
		return clients.Client{}, errors.New("clients missing deviceId")
	}

	_, err := uuid.FromString(c.DeviceId)
	if err != nil {
		return clients.Client{}, err
	}
	clnt.DeviceId = c.DeviceId

	// Parse label
	clnt.Label = clnt.Label

	// Parse client's service policy
	for _, csp := range c.Services {
		sp, err := resolveClientsFileServicePolicy(csp, serv)
		if err != nil {
			return clients.Client{}, err
		}
		clnt.Services = append(clnt.Services, sp)
	}

	return clnt, nil
}

// Resolves a clientFileServicePolicy into a clients.ServicePolicy struct
func resolveClientsFileServicePolicy(csp clientFileServicePolicy, serv []services.Service) (
	clients.ServicePolicy, error) {

		sp := clients.ServicePolicy{}

		found := false
		for _, s := range serv {
			if s.Name == csp.Name {
				sp.Service = s
				found = true
				break
			}
		}

		if !found {
			return clients.ServicePolicy{}, errors.New("non-existing service")
		}

		return sp, nil
}