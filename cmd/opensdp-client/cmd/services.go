package cmd

import (
	"github.com/spf13/cobra"
	"github.com/greenstatic/opensdp/internal/client"
	log "github.com/sirupsen/logrus"
	"os"
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

		c.Request()
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
