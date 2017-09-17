package cmd

// Row contains the data in each csv row
type Row struct {
	ActivationDate   string
	DeactivationDate string
}

// Result stores the real activation date
type Result struct {
	RealActivationDate string
}
