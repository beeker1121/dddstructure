package user

import (
	"dddstructure/proto"
	"dddstructure/service/errors"
)

// ValidateCreateParams validates the create parameters.
func (s *Service) ValidateCreateParams(params *proto.UserCreateParams) error {
	// Create a new ParamErrors.
	pes := errors.NewParamErrors()

	// Check email.
	if params.Email == "" {
		pes.Add(errors.NewParamError("email", errors.ErrUserEmailEmpty))
	}

	// Check password.
	if len(params.Password) < 8 {
		pes.Add(errors.NewParamError("password", errors.ErrUserPassword))
	}

	// Return if there were parameter errors.
	if pes.Length() > 0 {
		return pes
	}

	return nil
}
