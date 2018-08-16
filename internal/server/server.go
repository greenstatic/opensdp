package server

import (
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"crypto/x509"
	"crypto/tls"
	"net/http"
	"net"
	"github.com/greenstatic/opensdp/internal/services"
	"github.com/greenstatic/opensdp/internal/clients"
	"encoding/json"
	"time"
)

var (
	Version = "0.1.0"
)

type Server struct {
	CAPath string
	ServerCertPath string
	ServerKeyPath string
	Bind string
	Port string
	Services []services.Service
	Clients map[string]clients.Client
}


func rootResponse(w http.ResponseWriter, req *http.Request) {
	cn := req.TLS.PeerCertificates[0].Subject.CommonName

	json.NewEncoder(w).Encode(struct {
		Success bool `json:"success"`
		Msg string `json:"msg"`
		DeviceId string `json:"deviceId"`
		Datetime string `json:"datetime"`
		Version string `json:"version"`
	}{
		true,
		"OpenSDP Server",
		cn,
		time.Now().Format(time.RFC3339),
		Version,
	})

}

func (s *Server) Start() {
	// Adapted from: https://github.com/levigross/go-mutual-tls

	certBytes, err := ioutil.ReadFile(s.CAPath)
	if err != nil {
		log.WithField("caCert", s.CAPath).Error("Failed to read CA cert")
		log.Error(err)
		return
	}

	clientCertPool := x509.NewCertPool()
	if ok := clientCertPool.AppendCertsFromPEM(certBytes); !ok {
		log.Error("Failed to add CA cert to our cert pool")
		log.Error(err)
	}

	tlsConfig := &tls.Config{
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs: clientCertPool,
		CipherSuites: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
		PreferServerCipherSuites: true,
		MinVersion: tls.VersionTLS12,
	}

	tlsConfig.BuildNameToCertificate()

	http.HandleFunc("/discover", s.discoverResponseWrapper())
	http.HandleFunc("/", rootResponse)

	httpServer := &http.Server{
		Addr: net.JoinHostPort(s.Bind, s.Port),
		TLSConfig: tlsConfig,
	}

	// Disable HTTP/2 support due to cipher suite error
	httpServer.TLSNextProto = map[string]func(*http.Server, *tls.Conn, http.Handler){}

	log.WithFields(log.Fields{
		"bind": s.Bind,
		"port": s.Port,
	}).Info("Starting server")

	log.Fatalln(httpServer.ListenAndServeTLS(s.ServerCertPath, s.ServerKeyPath))
}