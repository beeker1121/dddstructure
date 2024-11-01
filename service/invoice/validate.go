package invoice

import (
	"dddstructure/proto"
	"dddstructure/service/errors"
)

// ValidateCreateParams validates the create parameters.
func (s *Service) ValidateCreateParams(params *proto.InvoiceCreateParams) error {
	// Create a new ParamErrors.
	pes := errors.NewParamErrors()

	// Check amount.
	if params.AmountDue > 1000000 {
		pes.Add(errors.NewParamError("amount_due", errors.ErrTransactionAmountLimit))
	}

	// Return if there were parameter errors.
	if pes.Length() > 0 {
		return pes
	}

	return nil
}
