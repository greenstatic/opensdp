package cmd

import (
	"github.com/spf13/cobra"
	"github.com/greenstatic/opensdp/internal/client"
	log "github.com/sirupsen/logrus"
	"os"
	"fmt"
	"strings"
)

var (
	serverUrl string
	caPath string
	clientCertPath string
	clientKeyPath string
)

var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Returns the client's authorized services",
	Long: "Returns the client's authorized services",

	Run: func(cmd *cobra.Command, args []string) {

		ifEmptyReturnError(serverUrl, "missing server URL")
		ifEmptyReturnError(caPath, "missing ca certificate")
		ifEmptyReturnError(clientCertPath, "missing client certificate")
		ifEmptyReturnError(clientKeyPath, "missing client key")

		c := client.Client{
			serverUrl,
			caPath,
			clientCertPath,
			clientKeyPath,
			}

		services, err := c.Discover()
		if err != nil {
			log.Error("Failed to perform discover exchange")
			log.Error(err)
			os.Exit(unexpectedError)
		}

		if len(services) == 0 {
			fmt.Println("You do not have access to any services")
			return
		}

		fmt.Println("You have access to the following services:")
		fmt.Printf("|%-26s|%-26s|%-18s|%-12s|%-20s|\n", "Name", "IP(s)", "Port(s)", "Access Type", "Tag(s)")

		const dashLen = 108
		for i := 0; i < dashLen; i++ {
			fmt.Printf("-")
		}
		fmt.Printf("\n")

		for _, s := range services {
			ips := strings.Join(s.IpsToStrings(), ", ")
			ports := strings.Join(s.ProtoPortToString(), ", ")
			accessTypes := strings.Join(s.AccessTypeToString(), ", ")
			tags := strings.Join(s.Tags, ", ")
			fmt.Printf("|%-26s|%-26s|%-18s|%-12s|%-20s|\n", s.Name, ips, ports, accessTypes, tags)
		}

	},
}

func ifEmptyReturnError(variable, err string) {
	if variable == "" {
		log.Fatalf(err)
		os.Exit(badInput)
	}
}

func init() {
	servicesCmd.Flags().StringVarP(&serverUrl, "serverUrl", "s", "",
		"OpenSDP server url")
	servicesCmd.Flags().StringVar(&caPath, "ca-cert", "", "certificate of the CA")
	servicesCmd.Flags().StringVarP(&clientCertPath, "client-cert", "c", "client.crt",
		"client's certificate")
	servicesCmd.Flags().StringVarP(&clientKeyPath, "client-key", "k", "client.key",
		"client's key")


	rootCmd.AddCommand(servicesCmd)
}
