package services

import (
	"net"
	"strconv"
	"fmt"
)

type Protocol int
const (
	// https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml
	ProtocolICMP Protocol = 1
	ProtocolTCP Protocol = 6
	ProtocolUDP Protocol = 17
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


type ProtoPort struct {
	Protocol Protocol
	Port uint16
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
	Name string
	Ips []net.IP
	ProtoPort []ProtoPort
	Tags []string
	AccessType []AccessType
}
