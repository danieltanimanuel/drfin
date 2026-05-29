package gmail

type GmailExtractedData struct {
	Title             EmailTitle                   `json:"title"`
	DetailTransaction []GmailDetailTransaction `json:"detail_transaction"`
}

type GmailDetailTransaction struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GmailDetailTransactionv2 struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
