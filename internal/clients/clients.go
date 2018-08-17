package clients

import "github.com/greenstatic/opensdp/internal/services"

type ServicePolicy struct {
	Service services.Service
}

type Client struct {
	DeviceId string
	Label    string
	Services []ServicePolicy
}
