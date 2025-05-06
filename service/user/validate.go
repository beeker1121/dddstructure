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

// ValidateLoginParams validates the login parameters.
func (s *Service) ValidateLoginParams(params *proto.UserLoginParams) error {
	// Create a new ParamErrors.
	pes := errors.NewParamErrors()

	// Check email.
	if params.Email == "" {
		pes.Add(errors.NewParamError("email", errors.ErrUserEmailEmpty))
	}

	// Return if there were parameter errors.
	if pes.Length() > 0 {
		return pes
	}

	return nil
}

// ValidateUpdateParams validates the update parameters.
func (s *Service) ValidateUpdateParams(params *proto.UserUpdateParams) error {
	// Update a new ParamErrors.
	pes := errors.NewParamErrors()

	// Check email.
	if params.Email != nil {
		if *params.Email == "" {
			pes.Add(errors.NewParamError("email", errors.ErrUserEmailEmpty))
		}
	}

	// Check password.
	if params.Password != nil {
		if len(*params.Password) < 8 {
			pes.Add(errors.NewParamError("password", errors.ErrUserPassword))
		}
	}

	// Return if there were parameter errors.
	if pes.Length() > 0 {
		return pes
	}

	return nil
}
