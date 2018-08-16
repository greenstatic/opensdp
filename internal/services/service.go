package services

import "net"

type Protocol int
const (
	// https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml
	ProtocolICMP Protocol = 1
	ProtocolTCP Protocol = 6
	ProtocolUDP Protocol = 17
)

type AccessType int
const (
	AccessTypeOpenSPA AccessType = iota
)

type ProtoPort struct {
	Protocol Protocol
	Port uint16
}

type Service struct {
	Name string
	Ips []net.IP
	ProtoPort []ProtoPort
	Tags []string
	AccessType []AccessType
}
