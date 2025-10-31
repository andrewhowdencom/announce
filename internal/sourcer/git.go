package sourcer

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/spf13/viper"
)

// GitFetcher is an implementation of Fetcher that fetches content from a git repository.
type GitFetcher struct{}

// NewGitFetcher creates a new GitFetcher.
func NewGitFetcher() *GitFetcher {
	return &GitFetcher{}
}

// Fetch fetches the content of a URL and returns it as a byte slice.
func (f *GitFetcher) Fetch(rawURL string) ([]byte, string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse url %s: %w", rawURL, err)
	}

	// The path is in the format /<user>/<repo>/tree/<ref>/<path/to/file>
	path := strings.TrimPrefix(u.Path, "/")
	pathParts := strings.SplitN(path, "/", 5)
	if len(pathParts) < 5 {
		return nil, "", fmt.Errorf("invalid git url path: %s. Expected /<user>/<repo>/tree/<ref>/<file>", u.Path)
	}
	user := pathParts[0]
	repo := pathParts[1]
	// pathParts[2] is "tree"
	ref := pathParts[3]
	filePath := pathParts[4]

	cloneURL := fmt.Sprintf("https://%s/%s/%s.git", u.Host, user, repo)

	// Create a temporary directory
	dir, err := ioutil.TempDir("", "ruf-git-sourcer")
	if err != nil {
		return nil, "", fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(dir)

	cloneOptions := &git.CloneOptions{
		URL:          cloneURL,
		SingleBranch: true,
		Depth:        1,
	}

	username := viper.GetString(fmt.Sprintf("git.auth.%s.username", u.Host))
	token := viper.GetString(fmt.Sprintf("git.auth.%s.token", u.Host))
	if token != "" {
		cloneOptions.Auth = &http.BasicAuth{
			Username: username,
			Password: token,
		}
	}

	// Determine if the ref is a branch, tag, or commit hash
	if len(ref) == 40 {
		// This is a commit hash
		cloneOptions.ReferenceName = plumbing.HEAD
	} else {
		// Try as a branch first
		cloneOptions.ReferenceName = plumbing.NewBranchReferenceName(ref)
	}

	r, err := git.PlainClone(dir, false, cloneOptions)
	if err != nil {
		// If it failed, try as a tag (and it wasn't a commit hash)
		if len(ref) != 40 {
			cloneOptions.ReferenceName = plumbing.NewTagReferenceName(ref)
			r, err = git.PlainClone(dir, false, cloneOptions)
			if err != nil {
				return nil, "", fmt.Errorf("failed to clone repo %s with ref %s (tried as branch and tag): %w", cloneURL, ref, err)
			}
		} else {
			return nil, "", fmt.Errorf("failed to clone repo %s with commit %s: %w", cloneURL, ref, err)
		}
	}

	if len(ref) == 40 {
		w, err := r.Worktree()
		if err != nil {
			return nil, "", fmt.Errorf("failed to get worktree: %w", err)
		}
		err = w.Checkout(&git.CheckoutOptions{
			Hash: plumbing.NewHash(ref),
		})
		if err != nil {
			return nil, "", fmt.Errorf("failed to checkout commit %s: %w", ref, err)
		}
	}

	head, err := r.Head()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get head for repo %s: %w", cloneURL, err)
	}
	hash := head.Hash().String()

	w, err := r.Worktree()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get worktree for repo %s: %w", cloneURL, err)
	}

	file, err := w.Filesystem.Open(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to open file '%s' in repo %s: %w", filePath, cloneURL, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read file '%s' in repo %s: %w", filePath, cloneURL, err)
	}

	return data, hash, nil
}
