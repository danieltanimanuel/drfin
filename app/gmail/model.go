package gmail

type GmailExtractedData struct {
	Title             string                   `json:"title"`
	DetailTransaction []GmailDetailTransaction `json:"detail_transaction"`
}

type GmailDetailTransaction struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
