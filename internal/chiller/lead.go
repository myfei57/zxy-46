package chiller

import "sync"

type LeadRegistry struct {
	mu   sync.Mutex
	lead string
}

func NewLeadRegistry() *LeadRegistry {
	return &LeadRegistry{}
}

func (r *LeadRegistry) Claim(id string) (string, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.lead == "" {
		r.lead = id
	}
	return r.lead, r.lead == id
}

func (r *LeadRegistry) Leader() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.lead
}
