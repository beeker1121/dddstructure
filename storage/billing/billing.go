package billing

// Database defines the accounting database interface.
type Database interface {
	GetMerchantAmountsDue() ([]*MerchantAmountDue, error)
}

// MerchantAmountsDue defines defines the merchant amounts due.
type MerchantAmountDue struct {
	ID           uint
	MerchantID   uint
	MerchantName string
	AmountDue    uint
}
