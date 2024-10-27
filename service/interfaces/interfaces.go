package interfaces

import "dddstructure/proto"

// Service defines the main business logic service interface struct that will
// be used between services to call each other.
type Service struct {
	User        User
	Invoice     Invoice
	Transaction Transaction
}

// NewServiceParams defines the new service params.
type NewServiceParams struct {
	User        User
	Invoice     Invoice
	Transaction Transaction
}

// NewService creates a new service.
func NewService(params NewServiceParams) *Service {
	return &Service{
		User:        params.User,
		Invoice:     params.Invoice,
		Transaction: params.Transaction,
	}
}

// User defines the user service.
type User interface {
	Create(u *proto.UserCreateParams) (*proto.User, error)
	GetByID(id uint) (*proto.User, error)
}

// Invoice defines the invoice service.
type Invoice interface {
	Create(i *proto.Invoice) (*proto.Invoice, error)
	GetByID(id uint) (*proto.Invoice, error)
	Update(i *proto.Invoice) error
	Pay(id uint) (*proto.Invoice, error)
}

// Transaction defines the transaction service.
type Transaction interface {
	Process(t *proto.TransactionProcessParams) (*proto.Transaction, error)
}
