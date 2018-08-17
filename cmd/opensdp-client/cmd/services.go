package cmd

import (
	"fmt"
	"github.com/greenstatic/opensdp/internal/client"
	"github.com/greenstatic/opensdp/internal/services"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net"
	"os"
	"strconv"
	"strings"
)

var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Returns the client's authorized services",
	Long:  "Returns the client's authorized services",

	Run: func(cmd *cobra.Command, args []string) {

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

		// Perform the services lookup
		srvs, err := c.Discover()
		if err != nil {
			log.Error("Failed to perform discover exchange")
			log.Error(err)
			os.Exit(unexpectedError)
		}

		if len(srvs) == 0 {
			fmt.Println("You do not have access to any services")
			return
		}

		fmt.Println("You have access to the following services:")
		fmt.Printf("|%-26s|%-26s|%-18s|%-12s|%-20s|\n", "Name", "IP", "Port(s)", "Access Type", "Tag(s)")

		const dashLen = 108
		for i := 0; i < dashLen; i++ {
			fmt.Printf("-")
		}
		fmt.Printf("\n")

		for _, s := range srvs {
			ports := strings.Join(s.ProtoPortToString(), ", ")
			accessTypes := strings.Join(s.AccessTypeToString(), ", ")
			tags := strings.Join(s.Tags, ", ")
			fmt.Printf("|%-26s|%-26s|%-18s|%-12s|%-20s|\n", s.Name, s.IP.String(), ports, accessTypes, tags)
		}

	},
}

func init() {
	rootCmd.AddCommand(servicesCmd)
}

func openSdpUnlockUsingOpenSpa(c client.Client) {
	// Parse the OpenSDP server IP and port
	opensdpIp, opensdpPortStr, err := net.SplitHostPort(c.Server)
	if err != nil {
		log.Error("Failed to parse the server field into host and port")
		log.Error(err)
		os.Exit(badInput)
	}

	opensdpPortInt, err := strconv.Atoi(opensdpPortStr)
	if err != nil {
		log.Error("Failed to parse the server port field into an integer")
		log.Error(err)
		os.Exit(badInput)
	}

	// Create pseudo OpenSDP service
	opensdpService := services.Service{
		IP:        net.ParseIP(opensdpIp),
		ProtoPort: []services.ProtoPort{{Protocol: services.ProtocolTCP, Port: uint16(opensdpPortInt)}},
	}

	// Using OpenSPA request access to the OpenSDP server
	err = client.AccessOpenSPAService(opensdpService, false, c.OpenSPA.Path, c.OpenSPA.OSPA)
	if err != nil {
		log.Error("Failed to unlock the OpenSDP service port using OpenSPA")
		log.Error(err)
		os.Exit(unexpectedError)
	}
}
