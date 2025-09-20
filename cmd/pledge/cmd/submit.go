package cmd

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/go-git/go-git/v6/plumbing/transport"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	errRepoDirRemoteMissing = errors.New("remote is missing from repo directory")
)

var submitCommand = &cobra.Command{
	Use:     "submit [repo-dir]",
	Aliases: []string{"sub", "s"},
	Short:   "Submit a patch to your ledger based on a repo's last commit",
	Args:    cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		ledgerRepoDirectory, err := getLedgerRepo(ctx, log)
		if err != nil {
			return err
		}

		log := log.With("ledgerRepoDirectory", ledgerRepoDirectory)

		var sourceRepoDirectory string
		if len(args) > 0 {
			sourceRepoDirectory = args[0]
		} else {
			sourceRepoDirectory, err = os.Getwd()
			if err != nil {
				return err
			}
		}

		log = log.With("sourceRepoDirectory", sourceRepoDirectory)

		log.Debug("Opening source repo")

		repo, err := git.PlainOpen(sourceRepoDirectory)
		if err != nil {
			return err
		}

		log.Debug("Getting remote URL")

		remotes, err := repo.Remotes()
		if err != nil {
			return err
		}

		var remoteURL string
		if len(remotes) > 0 && len(remotes[0].Config().URLs) > 0 {
			ep, err := transport.NewEndpoint(remotes[0].Config().URLs[0])
			if err != nil {
				return err
			}

			remoteURL = filepath.Join(filepath.Clean(ep.Host), filepath.Clean(strings.TrimPrefix(ep.Path, "/")))
		} else {
			return errRepoDirRemoteMissing
		}

		log = log.With("remoteURL", remoteURL)

		ref, err := repo.Head()
		if err != nil {
			return err
		}

		commit, err := repo.CommitObject(ref.Hash())
		if err != nil {
			return err
		}

		var patch *object.Patch
		if commit.NumParents() > 0 {
			parent, err := commit.Parent(0)
			if err != nil {
				return err
			}

			patch, err = parent.Patch(commit)
			if err != nil {
				return err
			}
		} else {
			patch, err = commit.Patch(nil)
			if err != nil {
				return err
			}
		}

		patchText := fmt.Sprintf(`From %v Mon Sep 17 00:00:00 2001
From: %v
Date: %v
Subject: [PATCH] %v---
%v
`, commit.Hash.String(), commit.Author.String(), commit.Author.When.Format(time.RFC1123Z), commit.Message, patch.String())

		log = log.With("patchText", patchText)

		log.Debug("Writing patch to ledger repo")

		fmt.Println(patchText)

		return nil
	},
}

func init() {
	viper.AutomaticEnv()

	indexCommand.AddCommand(submitCommand)
}
