package services

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Protocol int

const (
	// https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml
	ProtocolICMP Protocol = 1
	ProtocolTCP  Protocol = 6
	ProtocolUDP  Protocol = 17
)

func (p *Protocol) String() string {
	switch *p {
	case ProtocolICMP:
		return "icmp"
	case ProtocolTCP:
		return "tcp"
	case ProtocolUDP:
		return "udp"
	default:
		return ""
	}
}

func (p *Protocol) FromString(s string) (Protocol, error) {
	switch strings.ToLower(s) {
	case "icmp":
		return ProtocolICMP, nil
	case "tcp":
		return ProtocolTCP, nil
	case "udp":
		return ProtocolUDP, nil
	default:
		return 0, errors.New("unknown protocol")
	}
}

type AccessType int

const (
	AccessTypeOpenSPA AccessType = iota
)

func (at *AccessType) String() string {
	switch *at {
	case AccessTypeOpenSPA:
		return "OpenSPA"
	default:
		return ""
	}
}

func (at *AccessType) FromString(s string) error {
	switch s {
	case "OpenSPA":
		*at = AccessTypeOpenSPA
	default:
		return errors.New("unknown access type")
	}
	return nil
}

type ProtoPort struct {
	Protocol Protocol
	Port     uint16
}

// Returns a string slice that contains the [proto, port]. In case the protocol
// does not use ports (eg. icmp) then just return the protocol, [proto].
func (pp *ProtoPort) StringSlice() []string {
	proto := pp.Protocol.String()
	p := strconv.Itoa(int(pp.Port))

	if proto == "icmp" {
		return []string{proto}
	}

	return []string{proto, p}
}

// Reverts a string slice into a ProtoPort (opposite of StringSlice).
func (pp *ProtoPort) FromStringSlice(s []string) error {
	if len(s) != 1 && len(s) != 2 {
		return errors.New("unknown ProtoPort string slice")
	}

	var proto Protocol
	proto, err := proto.FromString(s[0])
	if err != nil {
		return err
	}
	pp.Protocol = proto

	if len(s) == 1 {
		return nil
	}

	port, err := strconv.Atoi(s[1])
	if err != nil {
		return err
	}
	pp.Port = uint16(port)

	return nil
}

// Stringify the ProtoPort like so: 22/tcp
func (pp *ProtoPort) String() string {
	p := pp.StringSlice()

	if len(p) == 1 {
		// No port definition
		return p[0]
	} else if len(p) == 2 {
		// Port and proto split with slash
		// eg. 22/tcp
		return fmt.Sprintf("%s/%s", p[1], p[0])
	}

	// StringSlice should not return a slice that is not of length 1 or 2.
	return ""

}

type Service struct {
	Name       string
	IP         net.IP
	ProtoPort  []ProtoPort
	Tags       []string
	AccessType []AccessType
}

// Interface that implements function to gain access to a service
type ServiceAccess interface {
	Access(Service) error
}

// Returns slice of the services ports as strings
func (s *Service) ProtoPortToString() []string {
	pp := make([]string, 0, len(s.ProtoPort))
	for _, p := range s.ProtoPort {
		pp = append(pp, p.String())
	}
	return pp
}

// Returns slice of services access type as strings
func (s *Service) AccessTypeToString() []string {
	at := make([]string, 0, len(s.AccessType))
	for _, a := range s.AccessType {
		at = append(at, a.String())
	}
	return at
}
