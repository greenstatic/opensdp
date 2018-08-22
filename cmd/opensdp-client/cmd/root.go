package cmd

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

const defaultConfigFile = "config.yaml"

var (
	Version      = "0.1.2"
	Verbose      = false
	VerboseSplit = false
	ver          = false

	serverUrl      string
	caPath         string
	clientCertPath string
	clientKeyPath  string
	cfgFile        string

	openspaPath string
	openspaOSPA string
)

var rootCmd = &cobra.Command{
	Use:   "opensdp-client",
	Short: "OpenSDP client to gain access to authorized services",
	Long: `OpenSDP client will authenticate with your OpenSDP server and give you access
to your authorized services.`,
	Run: func(cmd *cobra.Command, args []string) {
		if ver {
			fmt.Printf("OpenSDP client version: %s\n", Version)
			return
		}
		cmd.Help()
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&serverUrl, "server", "s", "",
		"OpenSDP server")
	rootCmd.PersistentFlags().StringVar(&caPath, "ca-cert", "", "certificate of the CA")
	rootCmd.PersistentFlags().StringVarP(&clientCertPath, "certificate", "c", "client.crt",
		"client's certificate")
	rootCmd.PersistentFlags().StringVarP(&clientKeyPath, "key", "k", "client.key",
		"client's key")

	rootCmd.PersistentFlags().StringVar(&openspaPath, "openspa-path", "openspa",
		"OpenSPA path")
	rootCmd.PersistentFlags().StringVar(&openspaOSPA, "openspa-ospa", "client.ospa",
		"OpenSPA client OSPA file")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		fmt.Sprintf("config file (default: ./%s)", defaultConfigFile))
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&VerboseSplit, "verbose-split", false,
		"split output to stdout (until but not including error level) and stderr (error level)")
	rootCmd.Flags().BoolVar(&ver, "version", false, "version of the client")

	viper.BindPFlag("ca-cert", rootCmd.PersistentFlags().Lookup("ca-cert"))
	viper.BindPFlag("certificate", rootCmd.PersistentFlags().Lookup("certificate"))
	viper.BindPFlag("key", rootCmd.PersistentFlags().Lookup("key"))
	viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
	viper.BindPFlag("openspa-path", rootCmd.PersistentFlags().Lookup("openspa-path"))
	viper.BindPFlag("openspa-ospa", rootCmd.PersistentFlags().Lookup("openspa-ospa"))

	log.SetOutput(os.Stdout)
	cobra.OnInitialize(verboseSplit)
	cobra.OnInitialize(verboseLog)

	rootCmd.MarkFlagRequired("ca-cert")
	rootCmd.MarkFlagRequired("certificates")
	rootCmd.MarkFlagRequired("key")
	rootCmd.MarkFlagRequired("server")
}

// Read config values
func initConfig() {
	if cfgFile != "" {
		// Use config file path provided by the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// User default config file located inside the same dir as the executable
		exePath, err := os.Executable()
		if err != nil {
			panic(err)
		}

		viper.AddConfigPath(filepath.Dir(exePath))
		viper.SetConfigFile(defaultConfigFile)
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Error("failed to read config")
		log.Error(err)
		os.Exit(unexpectedError)
	}
}

// Used to route error level logs to stderr and the rest to stdout.
// Credits: https://github.com/sirupsen/logrus/issues/403#issuecomment-346437512
// This disables the feature of color output in case it's ran from a TTY.
type OutputSplitter struct{}

func (splitter *OutputSplitter) Write(p []byte) (n int, err error) {
	if bytes.Contains(p, []byte("level=error")) {
		return os.Stderr.Write(p)
	}
	return os.Stdout.Write(p)
}

// Enables verbose split - until error level to stdout, while error goes to stderr.
// This is to be used on cobra.OnInitialize() to enable globally for all commands
// if the verbose-split flag is present.
func verboseSplit() {
	if VerboseSplit {
		log.SetOutput(&OutputSplitter{})
	}
}

// Enables verbose logging (debug level logs). This is to be used on cobra.OnInitialize()
// to enable globally for all commands if the verbose flag is present.
func verboseLog() {
	if Verbose {
		log.SetLevel(log.DebugLevel)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(unexpectedError)
	}
}
