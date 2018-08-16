package configsyaml

import (
	"github.com/greenstatic/opensdp/internal/services"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"errors"
	"net"
	"fmt"
	"strings"
	"strconv"
	log "github.com/sirupsen/logrus"
)

type ports []string

type serviceFile struct {
	Name string
	Ips []string
	Ports []ports
	Tags []string
	AccessType []string `yaml:"accessType"`
}

type servicesFile struct {
	Version string
	Kind string
	Services []serviceFile
}

func ServicesRead(path string) ([]services.Service, error) {

	// Read file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse file
	sf := servicesFile{}
	err = yaml.Unmarshal(data, &sf)
	if err != nil {
		return nil, err
	}

	// Check if correct kind
	if sf.Kind != "services" {
		return nil, errors.New("file not kind services")
	}

	// Parse service
	allServices := make([]services.Service, 0, len(sf.Services))
	for _, s := range sf.Services {
		serv, err := parseService(s)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed parsing service: %s", err))
		}

		log.WithField("name", serv.Name).Debug("Loaded service configuration")

		allServices = append(allServices, serv)
	}

	return allServices, nil
}

func parseService(s serviceFile) (services.Service, error) {
	serv := services.Service{}

	if s.Name == "" {
		return services.Service{}, errors.New("service missing field name")
	}
	serv.Name = s.Name

	var err error

	// Parse IPs
	if len(s.Ips) == 0 {
		return services.Service{}, errors.New("missing field ips")
	}

	serv.Ips, err = parseIps(s.Ips)
	if err != nil {
		return services.Service{}, errors.New(fmt.Sprintf("failed to parse ip: %s", err))
	}

	// Parse protocols & ports
	if len(s.Ports) == 0 {
		return services.Service{}, errors.New("missing field ports")
	}
	serv.ProtoPort, err = parseProtocolAndPort(s.Ports)
	if err != nil {
		return services.Service{}, errors.New(fmt.Sprintf("bad field ports: %s", err))
	}

	// Parse tags
	serv.Tags = s.Tags

	// Parse access types
	if len(s.AccessType) == 0 {
		return services.Service{}, errors.New("missing field access types")
	}
	serv.AccessType, err = parseAccessTypes(s.AccessType)
	if err != nil {
		return services.Service{}, errors.New(fmt.Sprintf("failed to parse access type: %s", err))
	}

	return serv, nil
}

// Parses a list of strings into a slice of net.IP's
func parseIps(ipsStr []string) ([]net.IP, error) {
	ips := make([]net.IP, 0, len(ipsStr))
	for _, ipStr := range ipsStr {
		ip := net.ParseIP(ipStr)
		if ip == nil {
			return nil, errors.New("bad ip")
		}
		ips = append(ips, ip)
	}

	return ips, nil
}

// Parses a ProtocolPort filed
func parseProtocolAndPort(protoPortStrngs []ports) ([]services.ProtoPort, error) {

	protoPorts := make([]services.ProtoPort, 0, len(protoPortStrngs))

	for _, ppStr := range protoPortStrngs {
		pp := services.ProtoPort{}

		if len(ppStr) != 1 && len(ppStr) != 2 {
			return nil, errors.New("bad protocol port combo")
		}

		var protocol services.Protocol
		switch strings.ToUpper(ppStr[0]) {
		case "TCP":
			protocol = services.ProtocolTCP
		case "UDP":
			protocol = services.ProtocolUDP
		case "ICMP":
			protocol = services.ProtocolICMP
		default:
			return nil, errors.New("unknown protocol")
		}

		pp.Protocol = protocol

		if protocol == services.ProtocolICMP && len(ppStr) != 1 {
			return nil, errors.New("icmp has no ports")
		}

		if protocol != services.ProtocolICMP && len(ppStr) != 2 {
			return nil, errors.New("missing port field in ports entry")
		}

		// Parse port, but not for ICMP
		if protocol != services.ProtocolICMP {
			port, err := strconv.Atoi(ppStr[1])
			if err != nil {
				return nil, errors.New("bad port value")
			}
			pp.Port = uint16(port)
		}

		protoPorts = append(protoPorts, pp)
	}

	return protoPorts, nil
}

// Parses a list of access types into a services.AccessType slice
func parseAccessTypes(atStr []string) ([]services.AccessType, error) {
	aTypes := make([]services.AccessType, 0, len(atStr))

	for _, aType := range atStr {
		switch aType {
		case "OpenSPA":
			aTypes = append(aTypes, services.AccessTypeOpenSPA)
		default:
			return nil, errors.New("unknown access type")
		}
	}

	return aTypes, nil
}