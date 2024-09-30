package common

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type (
	CTX_CONFIG_FILE struct{}
	CTX_CLIENT      struct{}
	CTX_API_BASE    struct{}
	CTX_AUTH_BASE   struct{}
	CTX_FORMAT      struct{}
)

const (
	FORMAT_JSON  = "json"
	FORMAT_HUMAN = "human"
)

func ConfigViper(cfgFile string) string {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		fullpath := filepath.Join(home, ".config", "quail-cli")
		if _, err := os.Stat(fullpath); os.IsNotExist(err) {
			os.MkdirAll(fullpath, 0755)
		}

		viper.AddConfigPath(fullpath)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")

		cfgFile = filepath.Join(fullpath, "config.yaml")
	}

	viper.AutomaticEnv()

	return cfgFile
}
