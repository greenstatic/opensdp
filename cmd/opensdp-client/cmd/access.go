package cmd

import (
	"github.com/greenstatic/opensdp/internal/client"
	"github.com/greenstatic/opensdp/internal/services"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var (
	all bool
)

var accessCmd = &cobra.Command{
	Use:   "access",
	Short: "Performs access handshake for authorized service",
	Long:  "Performs access handshake for authorized service",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 && !all {
			log.Error("Missing service name (or use -a for all services)")
			os.Exit(badInput)
			return
		}

		openspaD := client.OpenSPADetails{
			viper.GetString("openspa-path"),
			viper.GetString("openspa-ospa"),
		}

		c := client.Client{
			viper.GetString("server"),
			viper.GetString("ca-cert"),
			viper.GetString("certificate"),
			viper.GetString("key"),
			openspaD,
		}

		// Unlock the OpenSDP service using OpenSPA
		openSdpUnlockUsingOpenSpa(c)
		// OpenSDP server should be available now

		srvs, err := c.Discover()
		if err != nil {
			log.Error("Failed to perform discover exchange")
			log.Error(err)
			os.Exit(unexpectedError)
		}

		if all {
			log.WithField("count", len(srvs)).Info("Gaining access to all authorized services")
			client.ConcurrentAccessServiceContinuous(c, srvs)
		} else {
			srv := findService(srvs, args[0])
			client.ConcurrentAccessServiceContinuous(c, []services.Service{srv})
		}
	},
}

func init() {
	accessCmd.Flags().BoolVarP(&all, "all", "a", false, "Access all services you have access to")

	rootCmd.AddCommand(accessCmd)
}

// Finds the service by name from the slice of all services
func findService(srvs []services.Service, name string) services.Service {
	var serv services.Service
	found := false
	for _, s := range srvs {
		if s.Name == name {
			serv = s
			found = true
			break
		}
	}

	if !found {
		log.WithField("service", name).Warning("Unknown service")

		srvsName := make([]string, 0, len(srvs))
		for _, s := range srvs {
			srvsName = append(srvsName, s.Name)
		}

		log.WithField("services", strings.Join(srvsName, ", ")).Info("You have access to these services")
		os.Exit(unknownService)
	}

	return serv
}
