package billing

// Database defines the accounting database interface.
type Database interface {
	GetMerchantAmountsDue() ([]*MerchantAmountsDue, error)
}

// MerchantAmountsDue defines defines the merchant amounts due.
type MerchantAmountsDue struct {
	ID           uint
	MerchantID   uint
	MerchantName string
	AmountDue    uint
}
