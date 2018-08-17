package client

import (
	"crypto/tls"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"crypto/x509"
	"net/http"
	netUrl "net/url"
)

type Client struct {
	Server string
	CAPath  string
	ClientCertPath string
	ClientKeyPath string
	OpenSPA OpenSPADetails
}

type OpenSPADetails struct {
	Path string
	OSPA string
}

// Perform a GET request on the urlpath of the client server. Return the
// response as a byte slice.
func (c *Client) Request(urlpath string) ([]byte, error) {

	// Build url
	urlRawStr := "https://" + c.Server
	urlParsed, err := netUrl.Parse(urlRawStr)
	if err != nil {
		log.WithField("url", urlRawStr).Error("Failed to build server url")
		return nil, err
	}

	url := urlParsed.String()
	if url[len(url) - 1:] != "/" {
		url += "/"
	}
	url += urlpath

	log.WithFields(log.Fields{
		"url": url,
		"ca": c.CAPath,
		"clientCert": c.ClientCertPath,
		"clientKey": c.ClientKeyPath}).Debug("Issuing services request")

	// Adapted from: https://github.com/levigross/go-mutual-tls

	cert, err := tls.LoadX509KeyPair(c.ClientCertPath, c.ClientKeyPath)
	if err != nil {
		log.Error("Unable to load client keypair")
		return nil, err
	}

	clientCACert, err := ioutil.ReadFile(c.CAPath)
	if err != nil {
		log.Error("Unable to open ca certificate")
		return nil, err
	}

	// Trust only the CA certificate
	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(clientCACert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      clientCertPool,
		ServerName: "OpenSDP-server",
	}

	tlsConfig.BuildNameToCertificate()

	client := http.Client{
		Transport: &http.Transport{TLSClientConfig: tlsConfig},
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Error("Failed to connect to the server")
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read response from server")
		return nil, err
	}

	log.WithFields(log.Fields{
		"url": url,
		"responseLength": len(body),
	}).Debug("Successfully connected to the server")

	return body, nil
}