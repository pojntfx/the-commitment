package cmd

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/go-git/go-git/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	verboseKey = "verbose"
	configKey  = "config"
	ledgerKey  = "ledger"
	pullKey    = "pull"
)

var (
	log          *slog.Logger
	indexCommand = &cobra.Command{
		Use:   "pledge",
		Short: "CLI to interact with a ledger repo for The Commitment",
		Long: `CLI to interact with a ledger repo for The Commitment. Allows you to submit patches to your ledger, check if your patch was submitted successfully, and view a log of your recent patch submissions.

For more information, please visit https://github.com/pojntfx/the-commitment.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			opts := &slog.HandlerOptions{}
			if viper.GetBool(verboseKey) {
				opts.Level = slog.LevelDebug
			}
			log = slog.New(slog.NewJSONHandler(os.Stderr, opts))

			if viper.IsSet(configKey) {
				viper.SetConfigFile(viper.GetString(configKey))

				log.DebugContext(cmd.Context(), "Config key set, reading from file", "path", viper.GetViper().ConfigFileUsed())

				if err := viper.ReadInConfig(); err != nil {
					return err
				}
			} else {
				configBase := xdg.ConfigHome
				configName := cmd.Root().Use

				viper.SetConfigName(configName)
				viper.AddConfigPath(configBase)

				log.DebugContext(cmd.Context(), "Config key not set, reading from default location", "path", filepath.Join(configBase, configName))

				if err := viper.ReadInConfig(); err != nil && !errors.As(err, &viper.ConfigFileNotFoundError{}) {
					return err
				}
			}

			return nil
		},
	}
)

func verifyLedgerRepo(ctx context.Context, log *slog.Logger) (*git.Repository, string, error) {
	ledgerRepoDirectory := viper.GetString(ledgerKey)

	log = log.With("ledgerRepoDirectory", ledgerRepoDirectory)

	log.DebugContext(ctx, "Checking whether ledger repo directory exists")

	if _, err := os.Stat(ledgerRepoDirectory); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.ErrorContext(ctx, "Ledger repo directory does not exist, please clone it first", "err", err)

			return nil, "", err
		}

		return nil, "", err
	}

	log.DebugContext(ctx, "Opening ledger repo")

	repo, err := git.PlainOpen(ledgerRepoDirectory)
	if err != nil {
		return nil, "", err
	}

	if viper.GetBool(pullKey) {
		log.DebugContext(ctx, "Pulling ledger repo")

		worktree, err := repo.Worktree()
		if err != nil {
			return nil, "", err
		}

		if err := worktree.Pull(&git.PullOptions{}); err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
			return nil, "", err
		}
	}

	return repo, ledgerRepoDirectory, nil
}

func Execute() error {
	indexCommand.PersistentFlags().BoolP(verboseKey, "v", false, "Whether to enable verbose logging")
	indexCommand.PersistentFlags().StringP(configKey, "c", "", "Config file to use (by default "+indexCommand.Use+".yaml in the XDG config directory is read if it exists)")
	indexCommand.PersistentFlags().StringP(ledgerKey, "l", filepath.Join(xdg.StateHome, indexCommand.Use, "ledger"), "Ledger repo directory")
	indexCommand.PersistentFlags().BoolP(pullKey, "p", true, "Whether to pull the ledger repo before reading from or writing to it")

	if err := viper.BindPFlags(indexCommand.PersistentFlags()); err != nil {
		return err
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	return indexCommand.Execute()
}
