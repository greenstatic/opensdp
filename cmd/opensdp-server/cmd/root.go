package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
	"bytes"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"path/filepath"
)

var (
	Version = "0.1.0"
	Verbose = false
	VerboseSplit = false
	ver = false
	cfgFile string
	servicesPath string
	clientsPath string

	caPath string
	serverCertPath string
	serverKeyPath string
	bind string
	port uint16
)

var rootCmd = &cobra.Command{
	Use:   "opensdp-server",
	Short: "OpenSDP server allows clients to access authorized services",
	Long: `OpenSDP server tells clients the necessary information to authenticate with
hidden services they are authorized to use.`,
	Run: func(cmd *cobra.Command, args []string) {
		if ver {
			fmt.Printf("OpenSDP server version: %s\n", Version)
			return
		}

		ifEmptyReturnError(viper.GetString("ca-cert"), "missing ca certificate")
		ifEmptyReturnError(viper.GetString("certificate"), "missing server certificate")
		ifEmptyReturnError(viper.GetString("key"), "missing server key")

		startServer()
	},
}

func ifEmptyReturnError(variable, err string) {
	if variable == "" {
		log.Fatalf(err)
		os.Exit(badInput)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initFlags)
	rootCmd.Flags().StringVar(&caPath, "ca-cert", "", "certificate of the CA")
	rootCmd.Flags().StringVarP(&serverCertPath, "certificate", "c", "",
		"certificate of the server")
	rootCmd.Flags().StringVarP(&serverKeyPath, "key", "k", "",
		"private key of the server")
	rootCmd.Flags().StringVarP(&bind, "bind", "b", "0.0.0.0",
		"bind server to IP")
	rootCmd.Flags().Uint16VarP(&port, "port", "p", 8443, "port to listen to")

	rootCmd.Flags().StringVar(&cfgFile, "config", "", "config file (default: ./config.yaml)")
	rootCmd.Flags().StringVar(&servicesPath, "services", "", "services file (default: ./services.yaml)")
	rootCmd.Flags().StringVar(&clientsPath, "clients", "", "clients file (default: ./clients.yaml)")

	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false,
		"verbose output")
	rootCmd.PersistentFlags().BoolVar(&VerboseSplit, "verbose-split", false,
		"split output to stdout (until but not including error level) and stderr (error level)")
	rootCmd.Flags().BoolVar(&ver, "version", false, "version of the server")

	viper.BindPFlag("ca-cert", rootCmd.Flags().Lookup("ca-cert"))
	viper.BindPFlag("certificate", rootCmd.Flags().Lookup("certificate"))
	viper.BindPFlag("key", rootCmd.Flags().Lookup("key"))
	viper.BindPFlag("bind", rootCmd.Flags().Lookup("bind"))
	viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))
	viper.BindPFlag("clients", rootCmd.Flags().Lookup("clients"))
	viper.BindPFlag("services", rootCmd.Flags().Lookup("services"))

	log.SetOutput(os.Stdout)
	cobra.OnInitialize(verboseSplit)
	cobra.OnInitialize(verboseLog)
}

// Read config values
func initConfig() {
	if cfgFile != "" {
		// Use config file path provided by the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// User default
		dir, err := os.Executable()
		if err != nil {
			panic(err)
		}

		viper.AddConfigPath(dir)
		viper.SetConfigName("config.yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Error("failed to read config")
		log.Error(err)
		os.Exit(unexpectedError)
	}
}

func initFlags() {
	dir, err := os.Executable()
	if err != nil {
		panic(err)
	}

	defaultServices := filepath.Join(dir, "services.yaml")
	if viper.GetString("services") == "" {
		viper.Set("services", defaultServices)
	}

	defaultClients := filepath.Join(dir, "clients.yaml")
	if viper.GetString("clients") == "" {
		viper.Set("clients", defaultClients)
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
