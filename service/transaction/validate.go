package transaction

import (
	"dddstructure/proto"
	"dddstructure/service/errors"
)

// ValidateProcessParams validates the process parameters.
func (s *Service) ValidateProcessParams(params *proto.TransactionProcessParams) error {
	// Create a new ParamErrors.
	pes := errors.NewParamErrors()

	// Check amount.
	if params.Amount > 1000000 {
		pes.Add(errors.NewParamError("amount", errors.ErrTransactionAmountLimit))
	}

	// Return if there were parameter errors.
	if pes.Length() > 0 {
		return pes
	}

	return nil
}
