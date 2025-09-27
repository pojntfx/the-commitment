# The Commitment

A contract requiring you to do one daily public OSS contribution, while allowing time for rest.

[![hydrun CI](https://github.com/pojntfx/the-commitment/actions/workflows/hydrun.yaml/badge.svg)](https://github.com/pojntfx/the-commitment/actions/workflows/hydrun.yaml)
[![Matrix](https://img.shields.io/matrix/the-commitment:matrix.org)](https://matrix.to/#/#the-commitment:matrix.org?via=matrix.org)

## Overview

The Commitment is a personal contract framework designed to encourage sustainable open-source contributions while preventing burnout through built-in safeguards.

It enables you to:

- **Consistently make daily contributions** with a requirement for at least one world-readable commit under an OSI-approved license before midnight local time
- **Track upstream contributions** through a ledger repository system that captures patches to external projects, experimental branches, or any work outside your main repository
- **Maintain work-life balance** with well-defined, registered exceptions for weekends, vacation, sick leave, and unforeseen circumstances

## Installation

In addition to the [contract versions themselves](#contract-versions), the CLI tool `pledge` is also available; it helps you submit contributions to your ledger repository. `pledge` is provided in the form of static binaries.

On Linux, you can install them like so:

```shell
$ curl -L -o /tmp/pledge "https://github.com/pojntfx/the-commitment/releases/latest/download/pledge.linux-$(uname -m)"
$ sudo install /tmp/pledge /usr/local/bin
```

On macOS, you can use the following:

```shell
$ curl -L -o /tmp/pledge "https://github.com/pojntfx/the-commitment/releases/latest/download/pledge.darwin-$(uname -m)"
$ sudo install /tmp/pledge /usr/local/bin
```

On Windows, the following should work (using PowerShell as administrator):

```PowerShell
Invoke-WebRequest https://github.com/pojntfx/the-commitment/releases/latest/download/pledge.windows-x86_64.exe -OutFile \Windows\System32\pledge.exe
```

You can find binaries for more operating systems and architectures on [GitHub releases](https://github.com/pojntfx/the-commitment/releases).

## Tutorial

### 1. Sign a Contract

There are multiple available versions of The Commitment, each with different limitations, safeguards and motivations. For most people, it is recommended to stick with the [With Time For Rest and Ledger](#with-time-for-rest-and-ledger) version, since it provides a good balance between work and life, and also includes some safeguards around force majeure. To sign the contract, simply [fork this repository](https://github.com/pojntfx/the-commitment/fork), add your name to the list of signatories, and optimally submit a PR.

### 2. Create a Ledger Repository (Optional)

The latest [With Time For Rest and Ledger](#with-time-for-rest-and-ledger) version of The Commitment introduces the concept of a ledger repository, which is a special repository that allows you to get proper attribution for contributions to external open source projects that can't be easily tracked. Please read the contract itself to learn more about this concept. The `pledge` CLI helps you use this special repository in an easy way. To set one up, simply [create a new repository](https://github.com/new) (preferably at something like `github.com/$yourusername/ledger`) and clone it to your local XDG state directory:

```shell
$ mkdir -p ~/.local/state/pledge/ledger
$ git clone git@github.com:$yourusername/ledger.git ~/.local/state/pledge/ledger
```

You can also customize the location of this repository; see the [reference](#reference) for more information.

### 3. Start Contributing

Now that you've signed up for The Commitment, all you have to do is to follow your contract version's text, which usually means creating & pushing at least one open-source commit per day. If your making a contribution to a project under your username (e.g. to `github.com/$yourusername/$yourproject`) to that project's `main` branch, there is nothing else you need to. If you're making a contribution to an upstream project or if you're working on some experiments that you might not want to merge into your `main` branch directly, instead commit and push like usual, but also submit your contribution to your ledger repository so that you get it attributed:

```shell
$ ssh-add ~/.ssh/id_rsa.pub # Make sure you add your SSH key to the SSH agent either setup PGP or disable PGP signing in git and pledge
$ git commit -s -S -m "feat: Add new feature"
$ git push
$ pledge submit # Submits a patch for your latest commit to your ledger repository
$ pledge check # Check if you submitted a patch to your ledger repository today
$ pledge log # List previously submitted patches in your ledger repository
```

For the [With Time For Rest and Ledger](#with-time-for-rest-and-ledger) version of The Commitment, the `optional-commit-days.csv` file should also be placed in this repository. You can find an example linked in the contract version's text below.

🚀 **That's it!** We hope you enjoy contributing with The Commitment.

## Contract Versions

### With Time For Rest and Ledger

This document serves as a personal contract, hereinafter referred to as "The Commitment," between the undersigned party (hereafter "Contributor") and the public.

1. **Scope of Commitment**: The Contributor agrees to make at least one world-readable, individually copyrightable commit ("Commit") to a project (hereafter "Project") for which they hold the copyright, under an Open Source Initiative (OSI)-approved license, every day before midnight local time.
2. **Ledger Repository**: The Contributor shall maintain a ledger repository containing:

   1. The `optional-commit-days.csv` file (e.g., [`optional-commit-days.csv`](https://github.com/pojntfx/ledger/blob/main/optional-commit-days.csv)) documenting optional commit days
   2. Patch files organized in directories representing the repositories to which patches were submitted (e.g., `github.com/pojntfx/senbara.git/1758264958-example-message-0f08d14c83979b1ed40afa9c6f373cdd52de1528.patch`)

   Contributions to the ledger repository satisfy the daily Commit requirement when the Contributor is making contributions to repositories they do not own or cannot write to the main branch of (e.g., upstream open-source projects, experimental branches, or kernel patches). All standard Commit requirements apply to ledger repository contributions, including being world-readable and submitted before midnight local time. Duplicate commits to both the Contribtor's ledger repository and a project as defined in Scope of Commitment are discouraged to keep commit contribution graphs clean, but not prohibited. To simplify the process of submitting contributions to the ledger repository, the [`pledge`](./cmd/pledge) tool is recommended.

3. **Frequency and Timing**: The Commit is to be made daily, with the stipulation that each Commit must be completed and submitted before the stroke of midnight, according to the local time zone of the Contributor. Commits should preferably be signed with the Contributor's PGP key to allow verification of signing time by the remote hosting the Contributor's repositories.
4. **Optional Commit Days**: The Contributor is granted the option to forgo the Commit on the following occasions:
   1. For periods officially designated by the Contributor's employer, who shall not be the Contributor themselves, as vacation or sick leave, with protocol varying as follows:
      - **Vacation**: The start and end date must be documented and made world-readable in `optional-commit-days.csv`. This documentation must occur prior to the start of the vacation period, according to the local time zone of the Contributor.
      - **Sick Leave**: Same as for vacation, except the start and end date can be updated to include sick leave on the same day the leave is requested, according to the local time zone of the Contributor.
   2. On weekends, defined as Saturday and Sunday in the Contributor's local time zone.
   3. The Contributor is granted three days per calendar year to forgo the Commit without prior registration in `optional-commit-days.csv`, to accommodate any unforeseen circumstances. Additionally, the Contributor shall not be considered in breach of this Commitment due to failure to make Commits caused by circumstances beyond reasonable control (force majeure), including but not limited to: natural disasters, widespread power or internet outages, or serious accidents. In both cases, the Contributor should retroactively document affected days in `optional-commit-days.csv` when feasible.
5. **License Agreement**: The Contributor affirms that all Commits made under this agreement will be licensed under an OSI-approved license, ensuring that the contributions align with open-source principles and standards.
6. **Term of Agreement**: This Commitment remains in effect until formally revoked or amended in writing by the Contributor.
7. **Acknowledgement of Terms**: By signing below, the Contributor acknowledges and agrees to the terms outlined in this document, committing to the frequency and conditions under which the Commits to the Project are to be made.

Signatories:

### With Time For Rest

This document serves as a personal contract, hereinafter referred to as "The Commitment," between the undersigned party (hereafter "Contributor") and the public.

1. **Scope of Commitment**: The Contributor agrees to make at least one world-readable, individually copyrightable commit ("Commit") to a project (hereafter "Project") for which they hold the copyright, under an Open Source Initiative (OSI)-approved license, every day before midnight local time.
2. **Frequency and Timing**: The Commit is to be made daily, with the stipulation that each Commit must be completed and submitted before the stroke of midnight, according to the local time zone of the Contributor.
3. **Optional Commit Days**: The Contributor is granted the option to forgo the Commit on the following occasions:
   1. For periods officially designated by the Contributor's employer, who shall not be the Contributor themselves, as vacation or sick leave, with protocol varying as follows:
      - **Vacation**: The start and end date must be documented and made world-readable in [optional-commit-days.csv](./optional-commit-days.csv). This documentation must occur prior to the start of the vacation period, according to the local time zone of the Contributor.
      - **Sick Leave**: Same as for vacation, except the start and end date can be updated to include sick leave on the same day the leave is requested, according to the local time zone of the Contributor.
   2. On weekends, defined as Saturday and Sunday in the Contributor's local time zone.
4. **License Agreement**: The Contributor affirms that all Commits made under this agreement will be licensed under an OSI-approved license, ensuring that the contributions align with open-source principles and standards.
5. **Term of Agreement**: This Commitment remains in effect until formally revoked or amended in writing by the Contributor.
6. **Acknowledgement of Terms**: By signing below, the Contributor acknowledges and agrees to the terms outlined in this document, committing to the frequency and conditions under which the Commits to the Project are to be made.

Signatories:

- Felicitas Pojtinger (@pojntfx) (signed 2024-01-01)

### Without Time For Rest

This document serves as a personal contract, hereinafter referred to as "The Commitment," between the undersigned party (hereafter "Contributor") and the public.

1. **Scope of Commitment**: The Contributor agrees to make at least one world-readable, individually copyrightable commit ("Commit") to a project (hereafter "Project") for which they hold the copyright, under an Open Source Initiative (OSI)-approved license, every day before midnight local time.
2. **Frequency and Timing**: The Commit is to be made daily, with the requirement that each Commit must be completed and submitted prior to midnight, according to the local time zone of the Contributor.
3. **No Exceptions**: The Contributor is expected to fulfill the commitment of daily Commits without exception.
4. **License Agreement**: The Contributor confirms that all Commits made under this agreement will be licensed under an OSI-approved license. This ensures that the contributions adhere to open-source principles and standards.
5. **Term of Agreement**: This Commitment is effective immediately upon signing and remains in force until formally terminated or amended in writing by the Contributor.
6. **Acknowledgement of Terms**: By signing below, the Contributor acknowledges and agrees to the terms outlined in this document, committing to the uninterrupted, daily frequency of Commits to the Project.

Signatories:

- Felicitas Pojtinger (@pojntfx) (signed 2020-09-01, revoked 2023-11-26)

## Reference

### Command Line Arguments

```shell
$ pledge --help
CLI to interact with a ledger repo for The Commitment. Allows you to submit patches to your ledger, check if your patch was submitted successfully, and view a log of your recent patch submissions.

For more information, please visit https://github.com/pojntfx/the-commitment.

Usage:
  pledge [command]

Available Commands:
  check       Check if you submitted a patch to your ledger today
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  log         List previously submitted patches in your ledger
  submit      Submit a patch to your ledger based on a [repo-dir]'s last -n commits

Flags:
  -c, --config string   Config file to use (by default pledge.yaml in the XDG config directory is read if it exists)
  -h, --help            help for pledge
  -l, --ledger string   Ledger repo directory (default "/home/pojntfx/.local/state/pledge/ledger")
  -p, --pull            Whether to pull the ledger repo before reading from or writing to it (default true)
  -v, --verbose         Whether to enable verbose logging

Use "pledge [command] --help" for more information about a command.
```

<details>
  <summary>Expand subcommand reference</summary>

#### Submit

```shell
$ pledge submit --help
Fetches the last -n commits in [repo-dir], creates patch files for them, commits the patch files to your ledger repo and pushes the ledger repo to its remote

Usage:
  pledge submit [repo-dir] [flags]

Aliases:
  submit, sub, s

Flags:
  -n, --commits uint   Number of commits in repo-dir to submit patches for (default 1)
  -h, --help           help for submit
  -k, --key string     PGP key ID for signing (default uses default key)
  -s, --sign           Whether to sign the commit for a patch (default true)

Global Flags:
  -c, --config string   Config file to use (by default pledge.yaml in the XDG config directory is read if it exists)
  -l, --ledger string   Ledger repo directory (default "/home/pojntfx/.local/state/pledge/ledger")
  -p, --pull            Whether to pull the ledger repo before reading from or writing to it (default true)
  -v, --verbose         Whether to enable verbose logging
```

#### Log

```shell
$ pledge log --help
Fetches the last commits from your ledger repository and prints the summaries for each associated patch

Usage:
  pledge log [flags]

Aliases:
  log, l

Flags:
  -h, --help   help for log

Global Flags:
  -c, --config string   Config file to use (by default pledge.yaml in the XDG config directory is read if it exists)
  -l, --ledger string   Ledger repo directory (default "/home/pojntfx/.local/state/pledge/ledger")
  -p, --pull            Whether to pull the ledger repo before reading from or writing to it (default true)
  -v, --verbose         Whether to enable verbose logging
```

#### Check

```shell
$ pledge check --help
Fetches the last commit from your ledger repository, checks if it was committed today and prints the associated patch

Usage:
  pledge check [flags]

Aliases:
  check, c

Flags:
  -h, --help   help for check

Global Flags:
  -c, --config string   Config file to use (by default pledge.yaml in the XDG config directory is read if it exists)
  -l, --ledger string   Ledger repo directory (default "/home/pojntfx/.local/state/pledge/ledger")
  -p, --pull            Whether to pull the ledger repo before reading from or writing to it (default true)
  -v, --verbose         Whether to enable verbose logging
```

</details>

### Environment Variables

All command line arguments described above can also be set using environment variables; for example, to set `--pull` to `false` with an environment variable, use `PULL=false`.

## Acknowledgements

- [GopenPGP](https://github.com/ProtonMail/gopenpgp) provides the PGP library for signing commits.
- [go-git](https://github.com/go-git/go-git) provides the Go Git implementation for making commits and interacting with the ledger repo.

## Contributing

To contribute, please use the [GitHub flow](https://guides.github.com/introduction/flow/) and follow our [Code of Conduct](./CODE_OF_CONDUCT.md).

To build and start a development version of pledge locally, run the following:

```shell
$ git clone https://github.com/pojntfx/the-commitment.git
$ cd the-commitment
$ go run ./cmd/pledge --help
```

Have any questions or need help? Chat with us [on Matrix](https://matrix.to/#/#the-commitment:matrix.org?via=matrix.org)!

## License

The Commitment (c) 2025 Felicitas Pojtinger and contributors

SPDX-License-Identifier: Apache-2.0
