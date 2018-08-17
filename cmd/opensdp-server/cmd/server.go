package cmd

import (
	"github.com/greenstatic/opensdp/internal/configsyaml"
	"github.com/greenstatic/opensdp/internal/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strconv"
)

func startServer() {
	portStr := strconv.Itoa(int(viper.GetInt("port")))

	servicesPath := viper.GetString("services")
	srvs, err := configsyaml.ServicesRead(servicesPath)
	if err != nil {
		log.WithField("services", servicesPath).Error("Failed to read services")
		log.Error(err)
		os.Exit(unexpectedError)
	}

	clientsPath := viper.GetString("clients")
	clnts, err := configsyaml.ClientsRead(clientsPath, srvs)
	if err != nil {
		log.WithField("clients", clientsPath).Error("Failed to read clients")
		log.Error(err)
		os.Exit(unexpectedError)
	}

	s := server.Server{
		viper.GetString("ca-cert"),
		viper.GetString("certificate"),
		viper.GetString("key"),
		viper.GetString("bind"),
		portStr,
		srvs,
		clnts,
	}

	s.Start()
}
