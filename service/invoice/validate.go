package invoice

import (
	"errors"

	"dddstructure/proto"
	serverrors "dddstructure/service/errors"
)

// ValidateCreateParams validates the create parameters.
func (s *Service) ValidateCreateParams(params *proto.InvoiceCreateParams) error {
	// Create a new ParamErrors.
	pes := serverrors.NewParamErrors()

	// Check bill to.
	if len(params.BillTo.FirstName) > 255 {
		pes.Add(serverrors.NewParamError("bill_to.first_name", errors.New("first name must be less than 255 characters")))
	}

	// Check payment methods.
	if len(params.PaymentMethods) == 0 {
		pes.Add(serverrors.NewParamError("payment_methods", serverrors.ErrInvoicePaymentMethodRequired))
	}

	for _, v := range params.PaymentMethods {
		if v != "card" && v != "ach" {
			pes.Add(serverrors.NewParamError("payment_methods", serverrors.ErrInvoicePaymentMethodInvalid))
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
	pes := serverrors.NewParamErrors()

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
	pes := serverrors.NewParamErrors()

	// Check payment methods.
	if params.PaymentMethods != nil {
		if len(*params.PaymentMethods) == 0 {
			pes.Add(serverrors.NewParamError("payment_methods", serverrors.ErrInvoicePaymentMethodRequired))
		}

		for _, v := range *params.PaymentMethods {
			if v != "card" && v != "ach" {
				pes.Add(serverrors.NewParamError("payment_methods", serverrors.ErrInvoicePaymentMethodInvalid))
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
	pes := serverrors.NewParamErrors()

	// Check amount.
	if params.AmountDue != nil {
		if *params.AmountDue > 1000000 {
			pes.Add(serverrors.NewParamError("amount_due", serverrors.ErrInvoiceAmountDueLimit))
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
	pes := serverrors.NewParamErrors()

	// Check amount.
	if params.Amount > 1000000 {
		pes.Add(serverrors.NewParamError("amount_due", serverrors.ErrInvoiceAmountDueLimit))
	}

	// Return if there were parameter errors.
	if pes.Length() > 0 {
		return pes
	}

	return nil
}
