package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/pgvillage-tools/chainsmith/internal/config"
	"github.com/pgvillage-tools/chainsmith/internal/version"
	"github.com/pgvillage-tools/chainsmith/pkg/tls"
)

// Create the root command object with version information
var rootCmd = &cobra.Command{
	Use:   "chainsmith",
	Short: "Chainsmith - A simple certificate chain manager",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use --help to see available commands.")
	},
	Version: version.GetAppVersion(),
}

func init() {
	rootCmd.PersistentFlags().String("config", os.Getenv("CMG_CONFIGFILE"), "Path to the config file")
	err := viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	if err != nil {
		panic(fmt.Errorf("init failed: %w", err).Error())
	}
	rootCmd.AddCommand(issueCmd, listCmd, revokeCmd)
}

// issueCmd generates CA and certificates based on a configuration file read with Viper.
var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Generate CA and certificates based on the configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := loadConfig(viper.GetString("config"))
		if err != nil {
			return err
		}
		return run(*config)
	},
}

// listCmd reads the configuration file using Viper and lists all issued certificates
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all issued certificates",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig(viper.GetString("config"))
		if err != nil {
			return err
		}
		fmt.Println("Issued Certificates:")
		for name, certCfg := range cfg.Certificates {
			fmt.Printf("- %s (%s)\n", name, certCfg.CommonName)
		}
		return nil
	},
}

// revokeCmd takes the name of a certificate to revoke and deletes the corresponding files.
// It also removes the certificate from the configuration file.
var revokeCmd = &cobra.Command{
	Use:   "revoke <certificate_name>",
	Short: "Revoke a certificate",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		certName := args[0]
		cfg, err := loadConfig(viper.GetString("config"))
		if err != nil {
			return err
		}

		certCfg, exists := cfg.Certificates[certName]
		if !exists {
			return fmt.Errorf("certificate '%s' not found", certName)
		}

		if err := os.Remove(certCfg.CertPath); err != nil {
			return fmt.Errorf("failed to delete certificate file: %v", err)
		}
		if err := os.Remove(certCfg.KeyPath); err != nil {
			return fmt.Errorf("failed to delete key file: %v", err)
		}

		fmt.Printf("Certificate '%s' revoked successfully.\n", certName)
		return nil
	},
}

// loadConfig unmarshals the configuration file into a Config struct using Viper and returns it.
func loadConfig(configPath string) (*config.Config, error) {
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// The run function does the heavy lifting of generating the certificate chain.
// It generates a root certificate authority (CA), an intermediate CA,
// and a set of certificates based on the provided configuration.
// It uses the `tls` package to create the certificates and keys.
// At this point it will still regenerate the whole chain.
// Todo: Add a flag to only regenerate the leaf certificates (Issues #2).
//
// Parameters:
//   - cfg: A config.Config object containing paths and settings for
//     generating the certificates and keys.
//
// Returns:
//   - error: An error if any step in the certificate generation process fails,
//     or nil if the operation completes successfully.
func run(cfg config.Config) error {
	rootCert, rootKey, err := tls.GenerateCA(cfg.RootCAPath, cfg.RootCAPath+".key", nil, nil, true)
	if err != nil {
		return err
	}

	intermediateCert, intermediateKey, err := tls.GenerateCA(cfg.IntermediateCAPath, cfg.IntermediateCAPath+".key", rootCert, rootKey, false)
	if err != nil {
		return err
	}

	for name, certCfg := range cfg.Certificates {
		log.Printf("Generating certificate for %s...", name)
		if err := tls.GenerateCert(certCfg.CertPath, certCfg.KeyPath, intermediateCert, intermediateKey, certCfg.CommonName); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
