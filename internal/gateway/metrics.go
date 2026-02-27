package gateway

import (
	"sync"
)

type Stats struct {
	TotalRequests int64 `json:"total_requests"`
	Accepted      int64 `json:"accepted"`
	RateLimit     int64 `json:"rate_limited"`
	Unauthorized  int64 `json:"unauthorized"`
	sync.Mutex
}

var GlobalStats = &Stats{}

func (s *Stats) RecordRequest(status int) {
	s.Lock()
	defer s.Unlock()
	s.TotalRequests++
	switch status {
	case 200:
		s.Accepted++
	case 401:
		s.Unauthorized++
	case 429:
		s.RateLimit++
	}
}
