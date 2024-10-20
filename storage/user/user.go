package user

// Database defines the user database interface.
type Database interface {
	Create(u *User) (*User, error)
	GetByID(id uint) (*User, error)
}

// User defines a user.
type User struct {
	ID       uint
	Username string
	Email    string
}
