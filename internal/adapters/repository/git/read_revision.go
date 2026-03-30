package git

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"

	"go.redsock.ru/moti/internal/models"
)

type revisionParts struct {
	CommitHash string
	Version    string
}

func (r *gitRepo) ReadRevision(ctx context.Context, requestedVersion models.RequestedVersion) (models.Revision, error) {
	var revParts revisionParts

	var err error

	switch {
	case requestedVersion.IsGenerated() || requestedVersion.IsCommitHash():
		revParts, err = r.readRevisionByCommitHash(ctx, requestedVersion)
		if err != nil {
			return models.Revision{}, rerrors.Wrap(err, "r.readRevisionByCommitHash")
		}
	case requestedVersion.IsOmitted():
		revParts, err = r.readRevisionForLatestCommit(ctx)
		if err != nil {
			return models.Revision{}, fmt.Errorf("r.readRevisionForLatestCommit: %w", err)
		}
	default:
		// in other case use readRevisionByGitTagVersion
		revParts, err = r.readRevisionByGitTagVersion(ctx, requestedVersion)
		if err != nil {
			return models.Revision{}, fmt.Errorf("r.readRevisionByGitTagVersion: %w", err)
		}
	}

	if revParts.CommitHash == "" {
		return models.Revision{}, models.ErrVersionNotFound
	}

	revision := models.Revision{
		CommitHash: revParts.CommitHash,
		Version:    revParts.Version,
	}
	log.Debug().Interface("revision", revision).Msg("Revision")

	return revision, nil
}

func (r *gitRepo) readRevisionByGitTagVersion(
	ctx context.Context, requestedVersion models.RequestedVersion) (revisionParts, error) {
	gitTagVersion := string(requestedVersion)

	res, err := r.console.RunCmd(ctx, r.cacheDir, "git", "ls-remote", "origin", gitTagVersion)
	if err != nil {
		return revisionParts{}, models.ErrVersionNotFound
	}

	commitHash := ""

	for _, lsOut := range strings.Split(res, "\n") {
		rev := strings.Fields(lsOut)
		if len(rev) != 2 {
			continue
		}

		if strings.HasPrefix(rev[1], gitRefsTagPrefix) &&
			strings.TrimPrefix(rev[1], gitRefsTagPrefix) == gitTagVersion {
			commitHash = rev[0]

			break
		}
	}

	parts := revisionParts{
		CommitHash: commitHash,
		Version:    gitTagVersion,
	}

	return parts, nil
}

func (r *gitRepo) readRevisionForLatestCommit(ctx context.Context) (revisionParts, error) {
	headInfo, err := r.console.RunCmd(
		ctx, r.cacheDir,
		"git",
		"ls-remote",
		"origin",
		gitLatestVersionRef,
	)
	if err != nil {
		return revisionParts{}, models.ErrVersionNotFound
	}

	// got commit hash from result
	lines := strings.Split(headInfo, "\n")
	if len(lines) == 0 {
		return revisionParts{}, fmt.Errorf("invalid lines of git info: %s", headInfo)
	}

	parts := strings.Fields(lines[0])
	if len(parts) != 2 {
		return revisionParts{}, fmt.Errorf("invalid parts of git info: %s", headInfo)
	}

	commitHash := parts[0]

	// try to get git tag for this commit
	version, err := r.getTagForCommit(ctx, commitHash)
	if err != nil {
		return revisionParts{}, rerrors.Wrap(err, "getTagForCommit")
	}

	if version != "" {
		return revisionParts{
			CommitHash: commitHash,
			Version:    version,
		}, nil
	}

	generatedVersion := models.GeneratedVersionParts{
		CommitHash: commitHash,
	}

	return revisionParts{
		CommitHash: commitHash,
		Version:    generatedVersion.GetVersionString(),
	}, nil
}

func (r *gitRepo) getTagForCommit(ctx context.Context, commitHash string) (string, error) {
	tagInfo, err := r.console.RunCmd(ctx, r.cacheDir, "git", "ls-remote", "origin")
	if err != nil {
		return "", fmt.Errorf("adapters.RunCmd (ls-remote tagInfo): %w", err)
	}

	for _, lsOut := range strings.Split(tagInfo, "\n") {
		rev := strings.Fields(lsOut)
		if len(rev) != 2 {
			continue
		}

		if rev[0] != commitHash {
			continue
		}

		if strings.HasPrefix(rev[1], gitRefsTagPrefix) {
			return strings.TrimPrefix(rev[1], gitRefsTagPrefix), nil
		}
	}

	return "", nil
}

func (r *gitRepo) readRevisionByCommitHash(
	ctx context.Context, requestedVersion models.RequestedVersion) (revisionParts, error) {
	commitHash := string(requestedVersion)

	_, err := r.console.RunCmd(
		ctx, r.cacheDir,
		"git",
		"fetch", "-f",
		"origin",
		"--depth=1",
		commitHash,
	)
	if err != nil {
		return revisionParts{}, models.ErrVersionNotFound
	}

	// try to get git tag for this commit
	version, err := r.getTagForCommit(ctx, commitHash)
	if err != nil {
		return revisionParts{}, rerrors.Wrap(err, "getTagForCommit")
	}

	if version != "" {
		return revisionParts{
			CommitHash: commitHash,
			Version:    version,
		}, nil
	}

	generatedVersion := models.GeneratedVersionParts{
		CommitHash: commitHash,
	}

	return revisionParts{
		CommitHash: commitHash,
		Version:    generatedVersion.GetVersionString(),
	}, nil
}
