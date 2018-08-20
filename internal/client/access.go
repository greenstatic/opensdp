package client

import (
	"errors"
	"github.com/greenstatic/opensdp/internal/openspa"
	"github.com/greenstatic/opensdp/internal/services"
	log "github.com/sirupsen/logrus"
	"time"
)

func (c *Client) Access(serv services.Service) error {

	for _, at := range serv.AccessType {
		switch at {
		case services.AccessTypeOpenSPA:
			return AccessOpenSPAService(serv, true, c.OpenSPA.Path, c.OpenSPA.OSPA)
		}
	}

	return errors.New("unsupported access type")
}

func AccessOpenSPAService(serv services.Service, continuous bool, openspaPath, ospa string) error {

	var defaultOpenSPAPort uint16 = 22211
	client := openspa.Client{
		openspaPath,
		ospa,
		serv.IP,
		defaultOpenSPAPort,
	}

	for _, port := range serv.ProtoPort {
		req := openspa.Request{
			port.Protocol.String(),
			port.Port,
			port.Port,
		}

		err := client.Send(req, continuous)
		if err != nil {
			return err
		}
	}

	return nil
}

func ConcurrentAccessServiceContinuous(c Client, srvs []services.Service) {

	failed := make(chan services.Service, 1)

	for _, srv := range srvs {
		serviceToAccess := srv
		go func() {
			err := c.Access(serviceToAccess)
			if err != nil {
				log.Error(err)
				failed <- serviceToAccess
				return
			}
		}()

		time.Sleep(200) // small delay
	}

	for {
		select {
		case srv := <-failed:
			log.WithField("serviceName", srv.Name).Error("Failed to access service")
		}
	}

}
