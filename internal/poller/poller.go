package poller

import (
	"fmt"
	"time"

	"github.com/andrewhowdencom/ruf/internal/model"
	"github.com/andrewhowdencom/ruf/internal/sourcer"
)

// Poller periodically checks for updates in a list of sources.
type Poller struct {
	sourcer    sourcer.Sourcer
	interval   time.Duration
	knownState map[string]string
}

// New creates a new Poller.
func New(sourcer sourcer.Sourcer, interval time.Duration) *Poller {
	return &Poller{
		sourcer:    sourcer,
		interval:   interval,
		knownState: make(map[string]string),
	}
}

// Poll checks for updates in the sources and returns the calls from the changed URLs.
func (p *Poller) Poll(urls []string) ([]*model.Call, error) {
	var allCalls []*model.Call
	for _, url := range urls {
		calls, err := p.pollURL(url)
		if err != nil {
			// If a source can't be found, we log the error and continue.
			fmt.Printf("Error checking source %s: %v\n", url, err)
			continue
		}
		allCalls = append(allCalls, calls...)
	}
	return allCalls, nil
}

func (p *Poller) pollURL(url string) ([]*model.Call, error) {
	calls, state, err := p.sourcer.Source(url)
	if err != nil {
		return nil, err
	}

	if p.knownState[url] == state {
		return nil, nil // No change
	}

	p.knownState[url] = state
	return calls, nil
}
