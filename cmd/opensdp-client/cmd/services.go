package cmd

import (
	"github.com/spf13/cobra"
	"github.com/greenstatic/opensdp/internal/client"
	log "github.com/sirupsen/logrus"
	"os"
	"fmt"
	"strings"
	"github.com/spf13/viper"
)

var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Returns the client's authorized services",
	Long: "Returns the client's authorized services",

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
		fmt.Printf("|%-26s|%-26s|%-18s|%-12s|%-20s|\n", "Name", "IP", "Port(s)", "Access Type", "Tag(s)")

		const dashLen = 108
		for i := 0; i < dashLen; i++ {
			fmt.Printf("-")
		}
		fmt.Printf("\n")

		for _, s := range services {
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
