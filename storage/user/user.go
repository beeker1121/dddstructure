package user

// Database defines the user database interface.
type Database interface {
	Create(u *User) (*User, error)
	GetByID(id uint) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(u *User) (*User, error)
}

// User defines a user.
type User struct {
	ID       uint
	Email    string
	Password string
}
