package login

import (
	"fmt"
	"log/slog"

	"github.com/quail-ink/quail-cli/cmd/common"
	"github.com/quail-ink/quail-cli/oauth"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Login to Quail using OAuth",
		Run: func(cmd *cobra.Command, args []string) {
			authBase := cmd.Context().Value(common.CTX_AUTH_BASE{}).(string)
			apiBase := cmd.Context().Value(common.CTX_API_BASE{}).(string)
			token, err := oauth.Login(authBase, apiBase)
			if err != nil {
				slog.Error("failed to login", "error", err)
				return
			}

			fullpath := common.ConfigViper("")

			viper.Set("app.access_token", token.AccessToken)
			viper.Set("app.refresh_token", token.RefreshToken)
			viper.Set("app.token_type", token.TokenType)
			viper.Set("app.expiry", token.Expiry)

			// if the config file doesn't exist, create it first
			err = viper.WriteConfigAs(fullpath)
			if err != nil {
				slog.Error("failed to save config", "error", err, "config", fullpath)
				return
			}

			fmt.Printf("Login successful. Access token saved to %s\n", fullpath)
		},
	}
}
