package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dbyond-nl/chainsmithgo/internal/config"
	"github.com/dbyond-nl/chainsmithgo/pkg/tls"
)

var rootCmd = &cobra.Command{
	Use:   "chainsmith",
	Short: "Chainsmith - A simple certificate chain manager",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use --help to see available commands.")
	},
}

func init() {
	rootCmd.PersistentFlags().String("config", "configs/config.yml", "Path to the config file")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	rootCmd.AddCommand(issueCmd, listCmd, revokeCmd)
}

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Generate CA and certificates based on the configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all issued certificates",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
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

var revokeCmd = &cobra.Command{
	Use:   "revoke <certificate_name>",
	Short: "Revoke a certificate",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		certName := args[0]
		cfg, err := loadConfig()
		if err != nil {
			return err
		}

		certCfg, exists := cfg.Certificates[certName]
		if !exists {
			return fmt.Errorf("Certificate '%s' not found", certName)
		}

		if err := os.Remove(certCfg.CertPath); err != nil {
			return fmt.Errorf("Failed to delete certificate file: %v", err)
		}
		if err := os.Remove(certCfg.KeyPath); err != nil {
			return fmt.Errorf("Failed to delete key file: %v", err)
		}

		fmt.Printf("Certificate '%s' revoked successfully.\n", certName)
		return nil
	},
}

func loadConfig() (*config.Config, error) {
	viper.SetConfigFile(viper.GetString("config"))
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func run() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

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
