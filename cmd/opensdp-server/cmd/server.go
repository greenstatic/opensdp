package cmd

import (
	"github.com/greenstatic/opensdp/internal/server"
	"strconv"
	"github.com/spf13/viper"
	"github.com/greenstatic/opensdp/internal/configsyaml"
	log "github.com/sirupsen/logrus"
	"os"
)

func startServer() {
	portStr := strconv.Itoa(int(viper.GetInt("port")))

	servicesPath := viper.GetString("services")
	srvs, err := configsyaml.Read(servicesPath)
	if err != nil {
		log.WithField("services", servicesPath).Error("Failed to read services")
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
	}

	s.Start()
}