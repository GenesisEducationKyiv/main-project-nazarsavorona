package models

type Rate struct {
	From string
	To   string
	Rate float64
}

func (r *Rate) GetFrom() string {
	return r.From
}

func (r *Rate) GetTo() string {
	return r.To
}

func (r *Rate) GetRate() float64 {
	return r.Rate
}
