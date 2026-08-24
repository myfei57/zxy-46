package chiller

import "sync"

// LeadRegistry tracks which chiller currently holds the lead role. It is
// safe for concurrent use: the lead field is a single-writer slot guarded by
// a mutex so that two start commands arriving at the same time cannot both
// believe they won the lead.
type LeadRegistry struct {
	lead string
	mu   sync.Mutex
}

func NewLeadRegistry() *LeadRegistry {
	return &LeadRegistry{}
}

// Claim attempts to assign the lead role to id. It returns the id of the
// chiller that is now lead and whether id itself won the role. If the role
// is vacant (or id already holds it) id becomes the lead. If another
// chiller already holds the role that chiller stays lead and Claim reports
// isLead=false, leaving the caller in standby. The mutex serializes the
// check-and-set, so exactly one of two racing callers wins.
func (r *LeadRegistry) Claim(id string) (string, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.lead == "" || r.lead == id {
		r.lead = id
		return id, true
	}
	return r.lead, false
}

// Release surrenders the lead role if id currently holds it; otherwise it is
// a no-op. A stopped or tripped lead must release the role so a standby can
// take over instead of being blocked by a lead that can no longer run.
func (r *LeadRegistry) Release(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.lead == id {
		r.lead = ""
	}
}

func (r *LeadRegistry) Leader() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.lead
}
