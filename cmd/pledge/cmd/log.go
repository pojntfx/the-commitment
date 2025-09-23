package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/go-git/go-git/v6/plumbing/storer"
	"github.com/spf13/cobra"
)

var logCommand = &cobra.Command{
	Use:     "log",
	Aliases: []string{"l"},
	Short:   "Display all previous commits from the ledger repo",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		ledgerRepo, ledgerRepoDirectory, err := verifyLedgerRepo(ctx, log)
		if err != nil {
			return err
		}

		log := log.With("ledgerRepoDirectory", ledgerRepoDirectory)

		log.Debug("Getting commit log from ledger repo")

		ref, err := ledgerRepo.Head()
		if err != nil {
			return err
		}

		commitLog, err := ledgerRepo.Log(&git.LogOptions{From: ref.Hash()})
		if err != nil {
			return err
		}
		defer commitLog.Close()

		logOutput, logInput := io.Pipe()
		go func() {
			defer logInput.Close()

			if err = commitLog.ForEach(func(c *object.Commit) error {
				var repo string
				if tree, err := c.Tree(); err == nil {
					tree.Files().ForEach(func(f *object.File) error {
						if strings.HasSuffix(f.Name, ".patch") {
							repo = filepath.Dir(f.Name)

							return storer.ErrStop
						}
						return nil
					})
				}

				fmt.Fprintf(logInput, "Repo: %v\nFrom %v Mon Sep 17 00:00:00 2001\nFrom: %v\nDate: %v\nSubject: [PATCH] %v---\n\n", repo, c.Hash.String(), c.Author.String(), c.Author.When.Format(time.RFC1123Z), c.Message)

				return nil
			}); err != nil {
				panic(err)
			}
		}()

		pager := os.Getenv("PAGER")
		if pager == "" {
			pager = "less"
		}

		pagerCmd := exec.CommandContext(ctx, pager)
		pagerCmd.Stdin = logOutput
		pagerCmd.Stdout = os.Stdout
		pagerCmd.Stderr = os.Stderr

		log.Debug("Writing log output to pager", "pager", pager)

		if err := pagerCmd.Run(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	indexCommand.AddCommand(logCommand)
}
