package localrepo

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"gitlab.com/gitlab-org/gitaly/v14/internal/git"
	"gitlab.com/gitlab-org/gitaly/v14/internal/gitalyssh"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
)

// Remote provides functionality of the 'remote' git sub-command.
type Remote struct {
	repo *Repo
}

// Add adds a new remote to the repository.
func (remote Remote) Add(ctx context.Context, name, url string, opts git.RemoteAddOpts) error {
	if err := validateNotBlank(name, "name"); err != nil {
		return err
	}

	if err := validateNotBlank(url, "url"); err != nil {
		return err
	}

	var stderr bytes.Buffer
	if err := remote.repo.ExecAndWait(ctx,
		git.SubSubCmd{
			Name:   "remote",
			Action: "add",
			Flags:  buildRemoteAddOptsFlags(opts),
			Args:   []string{name, url},
		},
		git.WithStderr(&stderr),
		git.WithRefTxHook(ctx, remote.repo, remote.repo.cfg),
	); err != nil {
		switch {
		case isExitWithCode(err, 3):
			// In Git v2.30.0 and newer (https://gitlab.com/git-vcs/git/commit/9144ba4cf52)
			return git.ErrAlreadyExists
		case isExitWithCode(err, 128) && bytes.HasPrefix(stderr.Bytes(), []byte("fatal: remote "+name+" already exists")):
			// ..in older versions we parse stderr
			return git.ErrAlreadyExists
		}

		return err
	}

	return nil
}

func buildRemoteAddOptsFlags(opts git.RemoteAddOpts) []git.Option {
	var flags []git.Option
	for _, b := range opts.RemoteTrackingBranches {
		flags = append(flags, git.ValueFlag{Name: "-t", Value: b})
	}

	if opts.DefaultBranch != "" {
		flags = append(flags, git.ValueFlag{Name: "-m", Value: opts.DefaultBranch})
	}

	if opts.Fetch {
		flags = append(flags, git.Flag{Name: "-f"})
	}

	if opts.Tags != git.RemoteAddOptsTagsDefault {
		flags = append(flags, git.Flag{Name: opts.Tags.String()})
	}

	if opts.Mirror != git.RemoteAddOptsMirrorDefault {
		flags = append(flags, git.ValueFlag{Name: "--mirror", Value: opts.Mirror.String()})
	}

	return flags
}

// Remove removes a named remote from the repository configuration.
func (remote Remote) Remove(ctx context.Context, name string) error {
	if err := validateNotBlank(name, "name"); err != nil {
		return err
	}

	var stderr bytes.Buffer
	if err := remote.repo.ExecAndWait(ctx,
		git.SubSubCmd{
			Name:   "remote",
			Action: "remove",
			Args:   []string{name},
		},
		git.WithStderr(&stderr),
		git.WithRefTxHook(ctx, remote.repo, remote.repo.cfg),
	); err != nil {
		switch {
		case isExitWithCode(err, 2):
			// In Git v2.30.0 and newer (https://gitlab.com/git-vcs/git/commit/9144ba4cf52)
			return git.ErrNotFound
		case isExitWithCode(err, 128) && strings.HasPrefix(stderr.String(), "fatal: No such remote"):
			// ..in older versions we parse stderr
			return git.ErrNotFound
		}

		return err
	}

	return nil
}

// SetURL sets the URL for a given remote.
func (remote Remote) SetURL(ctx context.Context, name, url string, opts git.SetURLOpts) error {
	if err := validateNotBlank(name, "name"); err != nil {
		return err
	}

	if err := validateNotBlank(url, "url"); err != nil {
		return err
	}

	var stderr bytes.Buffer
	if err := remote.repo.ExecAndWait(ctx,
		git.SubSubCmd{
			Name:   "remote",
			Action: "set-url",
			Flags:  buildSetURLOptsFlags(opts),
			Args:   []string{name, url},
		},
		git.WithStderr(&stderr),
		git.WithRefTxHook(ctx, remote.repo, remote.repo.cfg),
	); err != nil {
		switch {
		case isExitWithCode(err, 2):
			// In Git v2.30.0 and newer (https://gitlab.com/git-vcs/git/commit/9144ba4cf52)
			return git.ErrNotFound
		case isExitWithCode(err, 128) && strings.HasPrefix(stderr.String(), "fatal: No such remote"):
			// ..in older versions we parse stderr
			return git.ErrNotFound
		}

		return err
	}

	return nil
}

// Exists determines whether a given named remote exists.
func (remote Remote) Exists(ctx context.Context, name string) (bool, error) {
	cmd, err := remote.repo.Exec(ctx,
		git.SubCmd{Name: "remote"},
		git.WithRefTxHook(ctx, remote.repo, remote.repo.cfg),
	)
	if err != nil {
		return false, err
	}

	found := false
	scanner := bufio.NewScanner(cmd)
	for scanner.Scan() {
		if scanner.Text() == name {
			found = true
			break
		}
	}

	return found, cmd.Wait()
}

func buildSetURLOptsFlags(opts git.SetURLOpts) []git.Option {
	if opts.Push {
		return []git.Option{git.Flag{Name: "--push"}}
	}

	return nil
}

// FetchOptsTags controls what tags needs to be imported on fetch.
type FetchOptsTags string

func (t FetchOptsTags) String() string {
	return string(t)
}

var (
	// FetchOptsTagsDefault enables importing of tags only on fetched branches.
	FetchOptsTagsDefault = FetchOptsTags("")
	// FetchOptsTagsAll enables importing of every tag from the remote repository.
	FetchOptsTagsAll = FetchOptsTags("--tags")
	// FetchOptsTagsNone disables importing of tags from the remote repository.
	FetchOptsTagsNone = FetchOptsTags("--no-tags")
)

// FetchOpts is used to configure invocation of the 'FetchRemote' command.
type FetchOpts struct {
	// Env is a list of env vars to pass to the cmd.
	Env []string
	// CommandOptions is a list of options to use with 'git' command.
	CommandOptions []git.CmdOpt
	// Prune if set fetch removes any remote-tracking references that no longer exist on the remote.
	// https://git-scm.com/docs/git-fetch#Documentation/git-fetch.txt---prune
	Prune bool
	// Force if set fetch overrides local references with values from remote that's
	// doesn't have the previous commit as an ancestor.
	// https://git-scm.com/docs/git-fetch#Documentation/git-fetch.txt---force
	Force bool
	// Verbose controls how much information is written to stderr. The list of
	// refs updated by the fetch will only be listed if verbose is true.
	// https://git-scm.com/docs/git-fetch#Documentation/git-fetch.txt---quiet
	// https://git-scm.com/docs/git-fetch#Documentation/git-fetch.txt---verbose
	Verbose bool
	// Tags controls whether tags will be fetched as part of the remote or not.
	// https://git-scm.com/docs/git-fetch#Documentation/git-fetch.txt---tags
	// https://git-scm.com/docs/git-fetch#Documentation/git-fetch.txt---no-tags
	Tags FetchOptsTags
	// Stderr if set it would be used to redirect stderr stream into it.
	Stderr io.Writer
}

// ErrFetchFailed indicates that the fetch has failed.
type ErrFetchFailed struct {
	err error
}

// Error returns the error message.
func (e ErrFetchFailed) Error() string {
	return e.err.Error()
}

// FetchRemote fetches changes from the specified remote. Returns an ErrFetchFailed error in case
// the fetch itself failed.
func (repo *Repo) FetchRemote(ctx context.Context, remoteName string, opts FetchOpts) error {
	if err := validateNotBlank(remoteName, "remoteName"); err != nil {
		return err
	}

	var stderr bytes.Buffer
	if opts.Stderr == nil {
		opts.Stderr = &stderr
	}

	commandOptions := []git.CmdOpt{
		git.WithEnv(opts.Env...),
		git.WithStderr(opts.Stderr),
		git.WithDisabledHooks(),
	}
	commandOptions = append(commandOptions, opts.CommandOptions...)

	cmd, err := repo.gitCmdFactory.New(ctx, repo,
		git.SubCmd{
			Name:  "fetch",
			Flags: opts.buildFlags(),
			Args:  []string{remoteName},
		},
		commandOptions...,
	)
	if err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return ErrFetchFailed{errorWithStderr(err, stderr.Bytes())}
	}

	return nil
}

// FetchInternal performs a fetch from an internal Gitaly-hosted repository. Returns an
// ErrFetchFailed error in case git-fetch(1) failed.
func (repo *Repo) FetchInternal(
	ctx context.Context,
	remoteRepo *gitalypb.Repository,
	refspecs []string,
	opts FetchOpts,
) error {
	if len(refspecs) == 0 {
		return fmt.Errorf("fetch internal called without refspecs")
	}

	env, err := gitalyssh.UploadPackEnv(ctx, repo.cfg, &gitalypb.SSHUploadPackRequest{
		Repository:       remoteRepo,
		GitConfigOptions: []string{"uploadpack.allowAnySHA1InWant=true"},
	})
	if err != nil {
		return fmt.Errorf("fetch internal: %w", err)
	}

	var stderr bytes.Buffer
	if opts.Stderr == nil {
		opts.Stderr = &stderr
	}

	commandOptions := []git.CmdOpt{
		git.WithEnv(append(env, opts.Env...)...),
		git.WithStderr(opts.Stderr),
		git.WithRefTxHook(ctx, repo, repo.cfg),
		// We've observed performance issues when fetching into big repositories part of an
		// object pool. The root cause of this seems to be the connectivity check, which by
		// default will also include references of any alternates. Given that object pools
		// often have hundreds of thousands of references, this is quite expensive to
		// compute. Below config entry will disable listing of alternate refs: they
		// shouldn't even be included in the negotiation phase, so they aren't going to
		// matter in the connectivity check either.
		git.WithConfig(git.ConfigPair{Key: "core.alternateRefsCommand", Value: "exit 0 #"}),
	}
	commandOptions = append(commandOptions, opts.CommandOptions...)

	if err := repo.ExecAndWait(ctx,
		git.SubCmd{
			Name:  "fetch",
			Flags: append(opts.buildFlags(), git.Flag{Name: "--atomic"}),
			Args:  append([]string{gitalyssh.GitalyInternalURL}, refspecs...),
		},
		commandOptions...,
	); err != nil {
		return ErrFetchFailed{errorWithStderr(err, stderr.Bytes())}
	}

	return nil
}

func (opts FetchOpts) buildFlags() []git.Option {
	flags := []git.Option{}

	if !opts.Verbose {
		flags = append(flags, git.Flag{Name: "--quiet"})
	}

	if opts.Prune {
		flags = append(flags, git.Flag{Name: "--prune"})
	}

	if opts.Force {
		flags = append(flags, git.Flag{Name: "--force"})
	}

	if opts.Tags != FetchOptsTagsDefault {
		flags = append(flags, git.Flag{Name: opts.Tags.String()})
	}

	return flags
}

func validateNotBlank(val, name string) error {
	if strings.TrimSpace(val) == "" {
		return fmt.Errorf("%w: %q is blank or empty", git.ErrInvalidArg, name)
	}
	return nil
}

func envGitSSHCommand(cmd string) string {
	return "GIT_SSH_COMMAND=" + cmd
}

// PushOptions are options that can be configured for a push.
type PushOptions struct {
	// SSHCommand is the command line to use for git's SSH invocation. The command line is used
	// as is and must be verified by the caller to be safe.
	SSHCommand string
	// Force decides whether to force push all of the refspecs.
	Force bool
	// Config is the Git configuration which gets passed to the git-push(1) invocation.
	// Configuration is set up via `WithConfigEnv()`, so potential credentials won't be leaked
	// via the command line.
	Config []git.ConfigPair
}

// Push force pushes the refspecs to the remote.
func (repo *Repo) Push(ctx context.Context, remote string, refspecs []string, options PushOptions) error {
	if len(refspecs) == 0 {
		return errors.New("refspecs to push must be explicitly specified")
	}

	var env []string
	if options.SSHCommand != "" {
		env = append(env, envGitSSHCommand(options.SSHCommand))
	}

	var flags []git.Option
	if options.Force {
		flags = append(flags, git.Flag{Name: "--force"})
	}

	stderr := &bytes.Buffer{}
	if err := repo.ExecAndWait(ctx,
		git.SubCmd{
			Name:  "push",
			Flags: flags,
			Args:  append([]string{remote}, refspecs...),
		},
		git.WithStderr(stderr),
		git.WithEnv(env...),
		git.WithConfigEnv(options.Config...),
	); err != nil {
		return fmt.Errorf("git push: %w, stderr: %q", err, stderr)
	}

	return nil
}
