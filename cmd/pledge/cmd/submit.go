package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/go-git/go-git/v6/plumbing/transport"
	"github.com/pojntfx/the-commitment/pkg/pgp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	commitsKey = "commits"
	keyKey     = "key"
	signKeyKey = "sign"
)

var (
	errRepoDirRemoteMissing = errors.New("remote is missing from repo directory")
)

var submitCommand = &cobra.Command{
	Use:     "submit [repo-dir]",
	Aliases: []string{"sub", "s"},
	Short:   "Submit a patch to your ledger based on a repo's last -n commits",
	Args:    cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		ledgerRepo, ledgerRepoDirectory, err := verifyLedgerRepo(ctx, log)
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

		numCommits := viper.GetUint(commitsKey)
		if numCommits == 0 {
			numCommits = 1
		}

		log = log.With("sourceRepoDirectory", sourceRepoDirectory)

		log.DebugContext(ctx, "Opening source repo")

		repo, err := git.PlainOpen(sourceRepoDirectory)
		if err != nil {
			return err
		}

		log.DebugContext(ctx, "Getting remote URL")

		remotes, err := repo.Remotes()
		if err != nil {
			return err
		}

		var remoteDir string
		if len(remotes) > 0 && len(remotes[0].Config().URLs) > 0 {
			ep, err := transport.NewEndpoint(remotes[0].Config().URLs[0])
			if err != nil {
				return err
			}

			remoteDir = filepath.Join(filepath.Clean(ep.Host), filepath.Clean(strings.TrimPrefix(ep.Path, "/")))
		} else {
			return errRepoDirRemoteMissing
		}

		log = log.With("remoteDirectory", remoteDir)

		ref, err := repo.Head()
		if err != nil {
			return err
		}

		commit, err := repo.CommitObject(ref.Hash())
		if err != nil {
			return err
		}

		currentCommit := commit

		worktree, err := ledgerRepo.Worktree()
		if err != nil {
			return err
		}

		commitsToPush := 0
		for i := uint(0); i < numCommits && currentCommit != nil; i++ {
			var patch *object.Patch
			if currentCommit.NumParents() > 0 {
				parent, err := currentCommit.Parent(0)
				if err != nil {
					return err
				}

				patch, err = parent.Patch(currentCommit)
				if err != nil {
					return err
				}
			} else {
				patch, err = currentCommit.Patch(nil)
				if err != nil {
					return err
				}
			}

			patchText := fmt.Sprintf(`From %v Mon Sep 17 00:00:00 2001
From: %v
Date: %v
Subject: [PATCH] %v---
%v
`, currentCommit.Hash.String(), currentCommit.Author.String(), currentCommit.Author.When.Format(time.RFC1123Z), currentCommit.Message, patch.String())

			var (
				patchFileName = fmt.Sprintf("%v", currentCommit.Author.When.Unix()) + "-" + url.PathEscape(strings.Split(currentCommit.Message, "\n")[0]) + "-" + currentCommit.Hash.String() + ".patch"
				patchFilePath = filepath.Join(ledgerRepoDirectory, remoteDir, patchFileName)
			)

			commitLog := log.With("patchFilePath", patchFilePath)

			commitLog.DebugContext(ctx, "Writing patch to ledger repo")

			if err := os.MkdirAll(filepath.Dir(patchFilePath), os.ModePerm); err != nil {
				return err
			}

			if err := os.WriteFile(patchFilePath, []byte(patchText), os.ModePerm); err != nil {
				return err
			}

			patchFilePathRel, err := filepath.Rel(ledgerRepoDirectory, patchFilePath)
			if err != nil {
				return err
			}

			commitLog.DebugContext(ctx, "Adding patch file to git")

			_, err = worktree.Add(patchFilePathRel)
			if err != nil {
				return err
			}

			var signKey *openpgp.Entity
			if viper.GetBool(signKeyKey) {
				var keyID *string
				if v := viper.GetString(keyKey); v != "" {
					keyID = &v
				}

				commitLog = commitLog.With("keyID", keyID)

				commitLog.DebugContext(ctx, "Reading private key to sign commit with")

				secretKey, err := pgp.GetPGPSecretKey(ctx, keyID)
				if err != nil {
					return err
				}

				signKey = secretKey.GetEntity()
			}

			commitLog.DebugContext(ctx, "Committing patch file with original commit info")

			if _, err = worktree.Commit(currentCommit.Message, &git.CommitOptions{
				Author: &object.Signature{
					Name:  currentCommit.Author.Name,
					Email: currentCommit.Author.Email,
					When:  currentCommit.Author.When,
				},
				SignKey: signKey,
			}); err == nil {
				fmt.Println("Patch", patchFilePathRel, "committed successfully")

				commitsToPush++
			} else {
				if !errors.Is(err, git.ErrEmptyCommit) {
					return err
				}

				fmt.Println("Skipped committing patch", patchFilePathRel, "since it already exists in the ledger repo")
			}

			if i+1 < numCommits && currentCommit.NumParents() > 0 {
				currentCommit, err = currentCommit.Parent(0)
				if err != nil {
					return err
				}
			}
		}

		if commitsToPush <= 0 {
			fmt.Println("All patches already pushed to ledger repo remote, no need to push again")

			return nil
		}

		log.DebugContext(ctx, "Pushing commits to remote repo")

		if err := ledgerRepo.Push(&git.PushOptions{}); err != nil {
			return err
		}

		fmt.Println("All patches pushed to ledger repo remote")

		return nil
	},
}

func init() {
	submitCommand.PersistentFlags().UintP(commitsKey, "n", 1, "Number of commits to submit patches for")
	submitCommand.PersistentFlags().StringP(keyKey, "k", "", "PGP key ID for signing (default uses default key)")
	submitCommand.PersistentFlags().BoolP(signKeyKey, "s", true, "Whether to sign the commit for a patch")

	viper.AutomaticEnv()

	indexCommand.AddCommand(submitCommand)
}
