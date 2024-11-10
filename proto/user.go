package proto

// User defines a user.
type User struct {
	ID       uint
	Email    string
	Password string
}

// UserCreateParams defines the user create parameters.
type UserCreateParams struct {
	ID       uint
	Email    string
	Password string
}

// UserLoginParams defines the user login parameters.
type UserLoginParams struct {
	Email    string
	Password string
}
