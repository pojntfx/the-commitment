package cmd

import (
	"context"
	"errors"
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

var (
	errNoCommitFoundForToday = errors.New("no commit found for today")
)

var checkCommand = &cobra.Command{
	Use:     "check",
	Aliases: []string{"c"},
	Short:   "Check if a commit has been made for today",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		ledgerRepo, ledgerRepoDirectory, err := verifyLedgerRepo(ctx, log)
		if err != nil {
			return err
		}

		log := log.With("ledgerRepoDirectory", ledgerRepoDirectory)

		log.Debug("Checking for today's commit in ledger repo")

		ref, err := ledgerRepo.Head()
		if err != nil {
			return err
		}

		commitLog, err := ledgerRepo.Log(&git.LogOptions{From: ref.Hash()})
		if err != nil {
			return err
		}
		defer commitLog.Close()

		var (
			now      = time.Now()
			today    = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			tomorrow = today.AddDate(0, 0, 1)

			todayCommit *object.Commit
		)
		if err = commitLog.ForEach(func(c *object.Commit) error {
			commitDate := c.Author.When
			if commitDate.After(today) && commitDate.Before(tomorrow) {
				todayCommit = c

				return storer.ErrStop
			}

			if commitDate.Before(today) {
				return storer.ErrStop
			}

			return nil
		}); err != nil {
			return err
		}

		if todayCommit == nil {
			return errNoCommitFoundForToday
		}

		logOutput, logInput := io.Pipe()
		go func() {
			defer logInput.Close()

			if tree, err := todayCommit.Tree(); err == nil {
				tree.Files().ForEach(func(f *object.File) error {
					if strings.HasSuffix(f.Name, ".patch") {
						repo := filepath.Dir(f.Name)
						if content, err := f.Contents(); err == nil {
							fmt.Fprintf(logInput, "Repo: %v\n%s", repo, content)
						}

						return storer.ErrStop
					}
					return nil
				})
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
	indexCommand.AddCommand(checkCommand)
}
