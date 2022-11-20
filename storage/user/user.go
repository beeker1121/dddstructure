package user

// Database defines the user database interface.
type Database interface {
	Create(params *CreateParams) (*User, error)
	GetByID(id uint) (*User, error)
}

// User defines the user.
type User struct {
	ID    uint
	Name  string
	Email string
}

// CreateParams defines the create parameters.
type CreateParams struct {
	ID    uint
	Name  string
	Email string
}
