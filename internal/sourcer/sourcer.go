package sourcer

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/andrewhowdencom/ruf/internal/model"
	"gopkg.in/yaml.v3"
)

// Fetcher defines the interface for fetching content from a URL.
type Fetcher interface {
	Fetch(url string) ([]byte, string, error)
}

// CompositeFetcher is a fetcher that can handle multiple schemes.
type CompositeFetcher struct {
	fetchers map[string]Fetcher
}

// NewCompositeFetcher creates a new CompositeFetcher.
func NewCompositeFetcher() *CompositeFetcher {
	return &CompositeFetcher{
		fetchers: make(map[string]Fetcher),
	}
}

// AddFetcher adds a new fetcher for a given scheme.
func (f *CompositeFetcher) AddFetcher(scheme string, fetcher Fetcher) {
	f.fetchers[scheme] = fetcher
}

// Fetch fetches the content of a URL and returns it as a byte slice.
func (f *CompositeFetcher) Fetch(rawURL string) ([]byte, string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse url %s: %w", rawURL, err)
	}

	fetcher, ok := f.fetchers[u.Scheme]
	if !ok {
		return nil, "", fmt.Errorf("unsupported scheme: %s", u.Scheme)
	}

	return fetcher.Fetch(rawURL)
}

// HTTPFetcher is an implementation of Fetcher that fetches content over HTTP.
type HTTPFetcher struct {
	client *http.Client
}

// NewHTTPFetcher creates a new HTTPFetcher.
func NewHTTPFetcher() *HTTPFetcher {
	return &HTTPFetcher{
		client: &http.Client{},
	}
}

// Fetch fetches the content of a URL and returns it as a byte slice.
func (f *HTTPFetcher) Fetch(url string) ([]byte, string, error) {
	resp, err := f.client.Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch url %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to fetch url %s: status code %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	// Prefer ETag, but fall back to Last-Modified.
	var state string
	if etag := resp.Header.Get("ETag"); etag != "" {
		state = etag
	} else if lastModified := resp.Header.Get("Last-Modified"); lastModified != "" {
		state = lastModified
	} else {
		state = fmt.Sprintf("%x", sha256.Sum256(body))
	}

	return body, state, nil
}

// FileFetcher is an implementation of Fetcher that fetches content from a local file.
type FileFetcher struct{}

// NewFileFetcher creates a new FileFetcher.
func NewFileFetcher() *FileFetcher {
	return &FileFetcher{}
}

// Fetch fetches the content of a URL and returns it as a byte slice.
func (f *FileFetcher) Fetch(rawURL string) ([]byte, string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse url %s: %w", rawURL, err)
	}

	data, err := os.ReadFile(u.Path)
	if err != nil {
		return nil, "", err
	}

	return data, fmt.Sprintf("%x", sha256.Sum256(data)), nil
}

// Parser defines the interface for parsing content into a list of calls.
type Parser interface {
	Parse(data []byte) ([]*model.Call, error)
}

// YAMLParser is an implementation of Parser that parses YAML content.
type YAMLParser struct{}

// NewYAMLParser creates a new YAMLParser.
func NewYAMLParser() *YAMLParser {
	return &YAMLParser{}
}

// Parse parses a YAML byte slice and returns a list of calls.
func (p *YAMLParser) Parse(data []byte) ([]*model.Call, error) {
	var calls []*model.Call
	if err := yaml.Unmarshal(data, &calls); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}
	return calls, nil
}

// Sourcer is an interface that defines the methods for sourcing calls.
type Sourcer interface {
	Source(url string) ([]*model.Call, string, error)
}

// sourcer is the concrete implementation of the Sourcer interface.
type sourcer struct {
	fetcher Fetcher
	parser  Parser
}

// NewSourcer creates a new Sourcer.
func NewSourcer(fetcher Fetcher, parser Parser) Sourcer {
	return &sourcer{
		fetcher: fetcher,
		parser:  parser,
	}
}

// Source fetches and parses calls from a URL.
func (s *sourcer) Source(url string) ([]*model.Call, string, error) {
	data, state, err := s.fetcher.Fetch(url)
	if err != nil {
		return nil, "", err
	}

	calls, err := s.parser.Parse(data)
	if err != nil {
		return nil, "", err
	}

	return calls, state, nil
}
