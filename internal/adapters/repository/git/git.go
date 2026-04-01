package git

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"
	"golang.org/x/net/html"

	"go.redsock.ru/moti/internal/adapters/repository"
)

var _ repository.Repo = (*GitRepo)(nil)

// GitRepo implements repository.Repo interface
type GitRepo struct {
	// RemoteURL full repository RemoteURL address with schema
	RemoteURL string
	// CacheDir local cache directory for store repository
	CacheDir string
	// Console for call external commands
	Console Console
}

const (
	// for omitted package version. HEAD is git key word.
	gitLatestVersionRef = "HEAD"
	// tag prefix on output of ls-remote command
	gitRefsTagPrefix = "refs/tags/"
)

// Some links from go mod:
// cmd/go/internal/modfetch/codehost/git.go:65 - create work dir
// cmd/go/internal/modfetch/codehost/git.go:137 - git's struct

// Console temporary interface for Console commands, must be replaced from core.Console.
type Console interface {
	RunCmd(ctx context.Context, dir string, command string, commandParams ...string) (string, error)
}

// New returns GitRepo instance
// remote: full RemoteURL address without schema
func New(ctx context.Context, remote string, cacheDir string, console Console) (repository.Repo, error) {
	repo := &GitRepo{
		CacheDir: cacheDir,
		Console:  console,
	}

	var err error

	repo.RemoteURL, err = GetRemote(ctx, remote)
	if err != nil {
		return nil, rerrors.Wrap(err)
	}

	_, err = os.Stat(filepath.Join(repo.CacheDir, "objects"))
	if err == nil {
		return repo, nil
	}

	_, err = repo.Console.RunCmd(ctx, repo.CacheDir, "git", "init", "--bare")
	if err != nil {
		return nil, rerrors.Wrap(err, "adapters.RunCmd (init): %w", err)
	}

	_, err = repo.Console.RunCmd(ctx, repo.CacheDir, "git", "remote", "add", "origin", repo.RemoteURL)
	if err != nil {
		return nil, rerrors.Wrap(err, "adapters.RunCmd (add origin)")
	}

	return repo, nil
}

func GetRemote(ctx context.Context, remoteURL string) (string, error) {
	if !strings.HasPrefix(remoteURL, "http") {
		remoteURL = "https://" + remoteURL
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, remoteURL, nil)
	if err != nil {
		return "", rerrors.Wrap(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", rerrors.Wrap(err, "reading page")
	}

	defer func() {
		cErr := resp.Body.Close()
		if cErr != nil {
			log.Error().
				Err(cErr).
				Msg("failed to close response body")
		}
	}()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Panic().
			Err(err).
			Msg("failed to parse html")
	}

	var nodeWalker func(node *html.Node)

	nodeWalker = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "meta" {
			remoteURL = updateRemoteFromMeta(node, remoteURL)
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			nodeWalker(c)
		}
	}

	nodeWalker(doc)

	return remoteURL, nil
}

func updateRemoteFromMeta(node *html.Node, remoteURL string) string {
	metaName, content := getMetaNameAndContent(node)
	if metaName == "go-import" && content != "" {
		parts := strings.Fields(content)
		if len(parts) == 3 {
			return parts[2]
		}
	}

	return remoteURL
}

func getMetaNameAndContent(n *html.Node) (metaName, content string) {
	for _, attr := range n.Attr {
		if attr.Key == "name" && attr.Val == "go-import" {
			metaName = attr.Val
		}

		if attr.Key == "content" {
			content = attr.Val
		}
	}

	return metaName, content
}
