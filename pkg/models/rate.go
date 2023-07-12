package models

type Rate struct {
	from string
	to   string
	rate float64
}

func NewRate(from, to string, rate float64) *Rate {
	return &Rate{
		from: from,
		to:   to,
		rate: rate,
	}
}

func (r *Rate) From() string {
	return r.from
}

func (r *Rate) To() string {
	return r.to
}

func (r *Rate) Rate() float64 {
	return r.rate
}
