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

	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/spf13/cobra"
)

var (
	errNoCommitFoundForToday = errors.New("no commit found for today")
)

var checkCommand = &cobra.Command{
	Use:     "check",
	Aliases: []string{"c"},
	Short:   "Check if you submitted a patch to your ledger today",
	Long:    "Fetches the last commit from your ledger repository, checks if it was committed today and prints the associated patch",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		ledgerRepo, ledgerRepoDirectory, err := verifyLedgerRepo(ctx, log)
		if err != nil {
			return err
		}

		log := log.With("ledgerRepoDirectory", ledgerRepoDirectory)

		log.DebugContext(ctx, "Getting latest commit from ledger repo")

		ref, err := ledgerRepo.Head()
		if err != nil {
			return err
		}

		commit, err := ledgerRepo.CommitObject(ref.Hash())
		if err != nil {
			return err
		}

		log.DebugContext(ctx, "Checking if latest commit was done today")

		var (
			now      = time.Now()
			today    = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			tomorrow = today.AddDate(0, 0, 1)
		)
		commitDate := commit.Author.When
		if !(commitDate.After(today) && commitDate.Before(tomorrow)) {
			return errNoCommitFoundForToday
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

		var (
			repo          string
			patchFilePath string
		)
		for _, fp := range patch.FilePatches() {
			if _, to := fp.Files(); to != nil && strings.HasSuffix(to.Path(), ".patch") {
				repo = filepath.Dir(to.Path())
				patchFilePath = to.Path()

				break
			}
		}

		log = log.With("repo", repo, "patchFilePath", patchFilePath)

		log.DebugContext(ctx, "Reading patch file for latest commit")

		logOutput, logInput := io.Pipe()
		go func() {
			defer logInput.Close()

			fmt.Fprintf(logInput, "Repo: %v\n", repo)

			tree, err := commit.Tree()
			if err != nil {
				return
			}

			file, err := tree.File(patchFilePath)
			if err != nil {
				return
			}

			reader, err := file.Reader()
			if err != nil {
				return
			}
			defer reader.Close()

			io.Copy(logInput, reader)
		}()

		pager := os.Getenv("PAGER")
		if pager == "" {
			pager = "less"
		}

		pagerCmd := exec.CommandContext(ctx, pager)
		pagerCmd.Stdin = logOutput
		pagerCmd.Stdout = os.Stdout
		pagerCmd.Stderr = os.Stderr

		log.DebugContext(ctx, "Writing patch to pager", "pager", pager)

		if err := pagerCmd.Run(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	indexCommand.AddCommand(checkCommand)
}
