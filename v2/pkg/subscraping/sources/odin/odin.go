// Package odin
package odin

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/projectdiscovery/subfinder/v2/pkg/subscraping"
)

// Source is the passive scraping agent
type Source struct {
	apiKeys   []string
	timeTaken time.Duration
	errors    int
	results   int
	skipped   bool
}

const baseURL string = "https://api.odin.io/v1"
const maxPerPage int = 50

// Run function returns all subdomains found with the service
func (s *Source) Run(ctx context.Context, domain string, session *subscraping.Session) <-chan subscraping.Result {
	results := make(chan subscraping.Result)
	s.errors = 0
	s.results = 0

	go func() {
		defer func(startTime time.Time) {
			s.timeTaken = time.Since(startTime)
			close(results)
		}(time.Now())

		var req request
		req.Domain = domain
		req.Limit = maxPerPage

		apiKey := subscraping.PickRandom(s.apiKeys, s.Name())
		if apiKey == "" {
			s.skipped = true
			return
		}

		headers := map[string]string{"x-api-key": apiKey}
		for {
			jsonReq, err := req.ToJSON()
			if err != nil {
				results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
				s.errors++
				return
			}
			resp, err := session.Post(ctx, fmt.Sprintf("%s/domain/subdomain/search", baseURL), "", headers, jsonReq)
			if err != nil {
				results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
				s.errors++
				session.DiscardHTTPResponse(resp)
				return
			}
			defer resp.Body.Close()

			var r response
			err = json.NewDecoder(resp.Body).Decode(&r)
			if err != nil {
				results <- subscraping.Result{Source: s.Name(), Type: subscraping.Error, Error: err}
				s.errors++
				session.DiscardHTTPResponse(resp)
				return
			}

			for _, subdomain := range r.Data {
				results <- subscraping.Result{Source: s.Name(), Type: subscraping.Subdomain, Value: subdomain}
				s.results++
			}

			req.Start = r.Pagination.Last
			if r.Pagination.Last == nil {
				break
			} else {
				session.DiscardHTTPResponse(resp)
			}
		}
	}()

	return results
}

// Name returns the name of the source
func (s *Source) Name() string {
	return "odin"
}

func (s *Source) IsDefault() bool {
	return true
}

func (s *Source) HasRecursiveSupport() bool {
	return false
}

func (s *Source) NeedsKey() bool {
	return true
}

func (s *Source) AddApiKeys(keys []string) {
	s.apiKeys = keys
}

func (s *Source) Statistics() subscraping.Statistics {
	return subscraping.Statistics{
		Errors:    s.errors,
		Results:   s.results,
		TimeTaken: s.timeTaken,
		Skipped:   s.skipped,
	}
}
