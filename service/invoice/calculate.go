package invoice

import (
	"dddstructure/proto"
	"dddstructure/utils"
	"strconv"
)

// Amounts defines the invoice amounts.
type Amounts struct {
	AmountDue  uint
	AmountPaid uint
}

// CalculateAmountsParams defines the parameters for the calculate amounts
// function.
type CalculateAmountsParams struct {
	LineItems []proto.InvoiceLineItem
	TaxRate   string
}

// CalculateAmounts calculates the invoice amounts.
func CalculateAmounts(params CalculateAmountsParams) (Amounts, error) {
	// Create amounts.
	amounts := Amounts{}

	// Loop through line items.
	for _, li := range params.LineItems {
		amounts.AmountDue += li.Quantity * li.Price
	}

	// Calculate tax.
	if params.TaxRate != "" {
		taxRatef64, err := strconv.ParseFloat(params.TaxRate, 64)
		if err != nil {
			return Amounts{}, err
		}

		tax := utils.PercentageFromInt(int(amounts.AmountDue), taxRatef64, 0, utils.Bankers)
		amounts.AmountDue += uint(tax)
	}

	return amounts, nil
}
