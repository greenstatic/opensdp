package client

import (
	"crypto/tls"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"crypto/x509"
	"net/http"
	"fmt"
	netUrl "net/url"
)

type Client struct {
	ServerUrl string
	CAPath  string
	ClientCertPath string
	ClientKeyPath string
}


func (c *Client) Request() error {

	// Build url
	urlRawStr := "https://" + c.ServerUrl
	urlParsed, err := netUrl.Parse(urlRawStr)
	if err != nil {
		log.WithField("url", urlRawStr).Error("Failed to build server url")
		return err
	}

	url := urlParsed.String()
	if url[len(url) - 1:] != "/" {
		url += "/"
	}
	url += "discover"

	log.WithFields(log.Fields{
		"url": url,
		"ca": c.CAPath,
		"clientCert": c.ClientCertPath,
		"clientKey": c.ClientKeyPath}).Debug("Issuing services request")

	// Adapted from: https://github.com/levigross/go-mutual-tls

	cert, err := tls.LoadX509KeyPair(c.ClientCertPath, c.ClientKeyPath)
	if err != nil {
		log.Fatalln("Unable to load cert", err)
		return err
	}

	clientCACert, err := ioutil.ReadFile(c.CAPath)
	if err != nil {
		log.Fatal("Unable to open cert", err)
		return err
	}

	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(clientCACert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      clientCertPool,
	}

	tlsConfig.BuildNameToCertificate()

	client := http.Client{
		Transport: &http.Transport{TLSClientConfig: tlsConfig},
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Println("Unable to connect to server", err)
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s\n", body)

	return nil
}