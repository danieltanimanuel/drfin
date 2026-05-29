package processedgmail

import "time"

type Status string

const (
	StatusSuccess Status = "success"
	StatusFailed  Status = "failed"
)

type ProcessedGmailMessage struct {
	ID             int64
	GmailMessageID string
	InternalDate   int64
	TransactionID  *int64
	Status         Status
	ErrorMessage   *string
	ProcessedAt    time.Time
}
