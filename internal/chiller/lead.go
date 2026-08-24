package chiller

type LeadRegistry struct {
	lead string
}

func NewLeadRegistry() *LeadRegistry {
	return &LeadRegistry{}
}

func (r *LeadRegistry) Claim(id string) (string, bool) {
	r.lead = id
	return id, true
}

func (r *LeadRegistry) Leader() string {
	return r.lead
}
