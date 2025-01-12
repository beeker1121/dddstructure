package invoice

import (
	"dddstructure/proto"
	"dddstructure/service/errors"
)

// ValidateCreateParams validates the create parameters.
func (s *Service) ValidateCreateParams(params *proto.InvoiceCreateParams) error {
	// Create a new ParamErrors.
	pes := errors.NewParamErrors()

	// Check payment methods.
	if len(params.PaymentMethods) == 0 {
		pes.Add(errors.NewParamError("payment_methods", errors.ErrInvoicePaymentMethodRequired))
	}

	for _, v := range params.PaymentMethods {
		if v != "card" && v != "ach" {
			pes.Add(errors.NewParamError("payment_methods", errors.ErrInvoicePaymentMethodInvalid))
			break
		}
	}

	// Return if there were parameter errors.
	if pes.Length() > 0 {
		return pes
	}

	return nil
}

// ValidateGetParams validates the get parameters.
func (s *Service) ValidateGetParams(params *proto.InvoiceGetParams) error {
	// Create a new ParamErrors.
	pes := errors.NewParamErrors()

	// Check validation here.

	// Return if there were parameter errors.
	if pes.Length() > 0 {
		return pes
	}

	return nil
}

// ValidateUpdateParams validates the update parameters.
func (s *Service) ValidateUpdateParams(params *proto.InvoiceUpdateParams) error {
	// Create a new ParamErrors.
	pes := errors.NewParamErrors()

	// Check payment methods.
	if params.PaymentMethods != nil {
		if len(*params.PaymentMethods) == 0 {
			pes.Add(errors.NewParamError("payment_methods", errors.ErrInvoicePaymentMethodRequired))
		}

		for _, v := range *params.PaymentMethods {
			if v != "card" && v != "ach" {
				pes.Add(errors.NewParamError("payment_methods", errors.ErrInvoicePaymentMethodInvalid))
				break
			}
		}
	}

	// Return if there were parameter errors.
	if pes.Length() > 0 {
		return pes
	}

	return nil
}

// ValidateUpdateForTransactionParams validates the update parameters.
func (s *Service) ValidateUpdateForTransactionParams(params *proto.InvoiceUpdateForTransactionParams) error {
	// Create a new ParamErrors.
	pes := errors.NewParamErrors()

	// Check amount.
	if params.AmountDue != nil {
		if *params.AmountDue > 1000000 {
			pes.Add(errors.NewParamError("amount_due", errors.ErrInvoiceAmountDueLimit))
		}
	}

	// Return if there were parameter errors.
	if pes.Length() > 0 {
		return pes
	}

	return nil
}

// ValidatePayParams validates the pay parameters.
func (s *Service) ValidatePayParams(params *proto.InvoicePayParams) error {
	// Create a new ParamErrors.
	pes := errors.NewParamErrors()

	// Check amount.
	if params.Amount > 1000000 {
		pes.Add(errors.NewParamError("amount_due", errors.ErrInvoiceAmountDueLimit))
	}

	// Return if there were parameter errors.
	if pes.Length() > 0 {
		return pes
	}

	return nil
}
