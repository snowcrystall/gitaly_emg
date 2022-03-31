package git

import (
	"fmt"
	"log"
	"strings"
)

const (
	// scNoRefUpdates denotes a command which will never update refs
	scNoRefUpdates = 1 << iota
	// scNoEndOfOptions denotes a command which doesn't know --end-of-options
	scNoEndOfOptions
	// scGeneratesPackfiles denotes a command which may generate packfiles
	scGeneratesPackfiles
)

type commandDescription struct {
	flags                  uint
	opts                   []GlobalOption
	validatePositionalArgs func([]string) error
}

// commandDescriptions is a curated list of Git command descriptions for special
// git.ExecCommandFactory validation logic
var commandDescriptions = map[string]commandDescription{
	"am": {},
	"apply": {
		flags: scNoRefUpdates,
	},
	"archive": {
		// git-archive(1) does not support disambiguating options from paths from revisions.
		flags: scNoRefUpdates | scNoEndOfOptions,
	},
	"blame": {
		// git-blame(1) does not support disambiguating options from paths from revisions.
		flags: scNoRefUpdates | scNoEndOfOptions,
	},
	"bundle": {
		flags: scNoRefUpdates | scGeneratesPackfiles,
	},
	"cat-file": {
		flags: scNoRefUpdates,
	},
	"check-ref-format": {
		// git-check-ref-format(1) uses a hand-rolled option parser which doesn't support
		// `--end-of-options`.
		flags: scNoRefUpdates | scNoEndOfOptions,
	},
	"checkout": {
		// git-checkout(1) does not support disambiguating options from paths from
		// revisions.
		flags: scNoEndOfOptions,
	},
	"clone": {
		flags: scGeneratesPackfiles,
	},
	"commit": {
		flags: 0,
	},
	"commit-graph": {
		flags: scNoRefUpdates,
	},
	"commit-tree": {
		flags: scNoRefUpdates,
	},
	"config": {
		flags: scNoRefUpdates,
	},
	"count-objects": {
		flags: scNoRefUpdates,
	},
	"diff": {
		flags: scNoRefUpdates,
	},
	"diff-tree": {
		flags: scNoRefUpdates,
	},
	"fetch": {
		flags: 0,

		opts: []GlobalOption{
			// When fetching objects from an untrusted source, we want to always assert
			// that all objects are valid. Please refer to the receive-pack
			// description with regards to why we ignore some checks.
			ConfigPair{Key: "fetch.fsckObjects", Value: "true"},
			ConfigPair{Key: "fetch.fsck.badTimezone", Value: "ignore"},
			ConfigPair{Key: "fetch.fsck.missingSpaceBeforeDate", Value: "ignore"},
			// While git-fetch(1) by default won't write commit graphs, both CNG and
			// Omnibus set this value to true. This has caused performance issues when
			// doing internal fetches, and furthermore it's not encouraged to run such
			// maintenance tasks on "normal" Git operations. Instead, writing commit
			// graphs should be done in our housekeeping RPCs, which already know to do
			// so. So let's disable writing commit graphs on fetches -- if it really is
			// required, we can enable it on a case-by-case basis.
			ConfigPair{Key: "fetch.writeCommitGraph", Value: "false"},
		},
	},
	"for-each-ref": {
		flags: scNoRefUpdates,
	},
	"format-patch": {
		flags: scNoRefUpdates,
	},
	"fsck": {
		flags: scNoRefUpdates,
	},
	"gc": {
		flags: scNoRefUpdates | scGeneratesPackfiles,
	},
	"grep": {
		// git-grep(1) does not support disambiguating options from paths from
		// revisions.
		flags: scNoRefUpdates | scNoEndOfOptions,
	},
	"hash-object": {
		flags: scNoRefUpdates,
	},
	"init": {
		flags: scNoRefUpdates,
		opts: []GlobalOption{
			// We're not prepared for a world where the user has configured the default
			// branch to be something different from "master" in Gitaly's git
			// configuration. There explicitly override it on git-init.
			ConfigPair{Key: "init.defaultBranch", Value: DefaultBranch},
		},
	},
	"linguist": {
		// linguist is not a native Git command, so we cannot use --end-of-options.
		flags: scNoEndOfOptions,
	},
	"log": {
		flags: scNoRefUpdates,
	},
	"ls-remote": {
		flags: scNoRefUpdates,
	},
	"ls-tree": {
		flags: scNoRefUpdates,
	},
	"merge-base": {
		flags: scNoRefUpdates,
	},
	"merge-file": {
		flags: scNoRefUpdates,
	},
	"mktag": {
		flags: scNoRefUpdates,
	},
	"multi-pack-index": {
		flags: scNoRefUpdates,
	},
	"pack-refs": {
		flags: scNoRefUpdates,
	},
	"pack-objects": {
		flags: scNoRefUpdates | scGeneratesPackfiles,
	},
	"push": {
		flags: scNoRefUpdates,
	},
	"receive-pack": {
		flags: 0,
		opts: append([]GlobalOption{
			// In case the repository belongs to an object pool, we want to prevent
			// Git from including the pool's refs in the ref advertisement. We do
			// this by rigging core.alternateRefsCommand to produce no output.
			// Because Git itself will append the pool repository directory, the
			// command ends with a "#". The end result is that Git runs `/bin/sh -c 'exit 0 # /path/to/pool.git`.
			ConfigPair{Key: "core.alternateRefsCommand", Value: "exit 0 #"},

			// When receiving objects from an untrusted source, we want to always assert
			// that all objects are valid.
			ConfigPair{Key: "receive.fsckObjects", Value: "true"},

			// In the past, there was a bug in git that caused users to
			// create commits with invalid timezones. As a result, some
			// histories contain commits that do not match the spec. As we
			// fsck received packfiles by default, any push containing such
			// a commit will be rejected. As this is a mostly harmless
			// issue, we add the following flag to ignore this check.
			ConfigPair{Key: "receive.fsck.badTimezone", Value: "ignore"},

			// git-fsck(1) complains in case a signature does not have a space
			// between mail and date. The most common case where this can be hit
			// is in case the date is missing completely. This error is harmless
			// enough and we cope just fine parsing such signatures, so we can
			// ignore this error.
			ConfigPair{Key: "receive.fsck.missingSpaceBeforeDate", Value: "ignore"},

			// Make git-receive-pack(1) advertise the push options
			// capability to clients.
			ConfigPair{Key: "receive.advertisePushOptions", Value: "true"},
		}, hiddenReceivePackRefPrefixes()...),
	},
	"remote": {
		// While git-remote(1)'s `add` subcommand does support `--end-of-options`,
		// `remove` doesn't.
		flags: scNoEndOfOptions,
	},
	"repack": {
		flags: scNoRefUpdates | scGeneratesPackfiles,
		opts: []GlobalOption{
			// Write bitmap indices when packing objects, which
			// speeds up packfile creation for fetches.
			ConfigPair{Key: "repack.writeBitmaps", Value: "true"},
		},
	},
	"rev-list": {
		// We cannot use --end-of-options here because pseudo revisions like `--all`
		// and `--not` count as options.
		flags: scNoRefUpdates | scNoEndOfOptions,
		validatePositionalArgs: func(args []string) error {
			for _, arg := range args {
				// git-rev-list(1) supports pseudo-revision arguments which can be
				// intermingled with normal positional arguments. Given that these
				// pseudo-revisions have leading dashes, normal validation would
				// refuse them as positional arguments. We thus override validation
				// for two of these which we are using in our codebase. There are
				// more, but we can add them at a later point if they're ever
				// required.
				if arg == "--all" || arg == "--not" {
					continue
				}
				if err := validatePositionalArg(arg); err != nil {
					return fmt.Errorf("rev-list: %w", err)
				}
			}
			return nil
		},
	},
	"rev-parse": {
		// --end-of-options is echoed by git-rev-parse(1) if used without
		// `--verify`.
		flags: scNoRefUpdates | scNoEndOfOptions,
	},
	"show": {
		flags: scNoRefUpdates,
	},
	"show-ref": {
		flags: scNoRefUpdates,
	},
	"symbolic-ref": {
		flags: 0,
	},
	"tag": {
		flags: 0,
	},
	"update-ref": {
		flags: 0,
	},
	"upload-archive": {
		// git-upload-archive(1) has a handrolled parser which always interprets the
		// first argument as directory, so we cannot use `--end-of-options`.
		flags: scNoRefUpdates | scNoEndOfOptions,
	},
	"upload-pack": {
		flags: scNoRefUpdates | scGeneratesPackfiles,
		opts: []GlobalOption{
			ConfigPair{Key: "uploadpack.allowFilter", Value: "true"},
			// Enables the capability to request individual SHA1's from the
			// remote repo.
			ConfigPair{Key: "uploadpack.allowAnySHA1InWant", Value: "true"},
		},
	},
	"version": {
		flags: scNoRefUpdates,
	},
	"worktree": {
		flags: 0,
	},
}

func init() {
	// This is the poor-mans static assert that all internal ref prefixes are properly hidden
	// from git-receive-pack(1) such that they cannot be written to when the user pushes.
	receivePackDesc, ok := commandDescriptions["receive-pack"]
	if !ok {
		log.Fatal("could not find command description of git-receive-pack(1)")
	}

	hiddenRefs := map[string]bool{}
	for _, opt := range receivePackDesc.opts {
		configPair, ok := opt.(ConfigPair)
		if !ok {
			continue
		}
		if configPair.Key != "receive.hideRefs" {
			continue
		}

		hiddenRefs[configPair.Value] = true
	}

	for _, internalRef := range InternalRefPrefixes {
		if !hiddenRefs[internalRef] {
			log.Fatalf("command description of receive-pack is missing hidden ref %q", internalRef)
		}
	}
}

// mayUpdateRef indicates if a command is known to update references.
// This is useful to determine if a command requires reference hook
// configuration. A non-exhaustive list of commands is consulted to determine if
// refs are updated. When unknown, true is returned to err on the side of
// caution.
func (c commandDescription) mayUpdateRef() bool {
	return c.flags&scNoRefUpdates == 0
}

// mayGeneratePackfiles indicates if a command is known to generate
// packfiles. This is used in order to inject packfile configuration.
func (c commandDescription) mayGeneratePackfiles() bool {
	return c.flags&scGeneratesPackfiles != 0
}

// supportsEndOfOptions indicates whether a command can handle the
// `--end-of-options` option.
func (c commandDescription) supportsEndOfOptions() bool {
	return c.flags&scNoEndOfOptions == 0
}

// args validates the given flags and arguments and, if valid, returns the complete command line.
func (c commandDescription) args(flags []Option, args []string, postSepArgs []string) ([]string, error) {
	var commandArgs []string

	for _, o := range flags {
		args, err := o.OptionArgs()
		if err != nil {
			return nil, err
		}
		commandArgs = append(commandArgs, args...)
	}

	if c.supportsEndOfOptions() {
		commandArgs = append(commandArgs, "--end-of-options")
	}

	if c.validatePositionalArgs != nil {
		if err := c.validatePositionalArgs(args); err != nil {
			return nil, err
		}
	} else {
		for _, a := range args {
			if err := validatePositionalArg(a); err != nil {
				return nil, err
			}
		}
	}
	commandArgs = append(commandArgs, args...)

	if len(postSepArgs) > 0 {
		commandArgs = append(commandArgs, "--")
	}

	// post separator args do not need any validation
	commandArgs = append(commandArgs, postSepArgs...)

	return commandArgs, nil
}

func validatePositionalArg(arg string) error {
	if strings.HasPrefix(arg, "-") {
		return fmt.Errorf("positional arg %q cannot start with dash '-': %w", arg, ErrInvalidArg)
	}
	return nil
}

func hiddenReceivePackRefPrefixes() []GlobalOption {
	var cps []GlobalOption

	for _, ns := range InternalRefPrefixes {
		cps = append(cps, ConfigPair{Key: "receive.hideRefs", Value: ns})
	}

	return cps
}
