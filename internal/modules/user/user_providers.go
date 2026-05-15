package user

import (
	"daymark/config"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

func SetupProviders(cfg *config.Config) {
	if cfg == nil {
		return
	}

	goth.UseProviders(

		google.New(
			cfg.GOOGLE_CLIENT_ID,
			cfg.GOOGLE_CLIENT_SECRET,
			cfg.GOOGLE_CALLBACK_URL,
		),

		github.New(
			cfg.GITHUB_CLIENT_ID,
			cfg.GITHUB_CLIENT_SECRET,
			cfg.GITHUB_CALLBACK_URL,
		),
	)
}
