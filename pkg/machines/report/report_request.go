package report

// Data transfer object for missing update value type
type MissingUpdate struct {
	UpdateId string
	Severity int
	Duration string
}

// Data transfer object for report request
type ReportRequest struct {
	MachineName    string
	MissingUpdates []MissingUpdate
}
