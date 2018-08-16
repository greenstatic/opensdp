package client

import (
	"crypto/tls"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"crypto/x509"
	"net/http"
	"fmt"
)

type Client struct {
	ServerUrl string
	CAPath  string
	ClientCertPath string
	ClientKeyPath string
}


func (c *Client) Request() {
	url := "https://" + c.ServerUrl

	log.WithFields(log.Fields{
		"url": url,
		"ca": c.CAPath,
		"clientCert": c.ClientCertPath,
		"clientKey": c.ClientKeyPath}).Debug("Issuing services request")

	// Adapted from: https://github.com/levigross/go-mutual-tls

	cert, err := tls.LoadX509KeyPair(c.ClientCertPath, c.ClientKeyPath)
	if err != nil {
		log.Fatalln("Unable to load cert", err)
	}

	clientCACert, err := ioutil.ReadFile(c.CAPath)
	if err != nil {
		log.Fatal("Unable to open cert", err)
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
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s\n", body)
}