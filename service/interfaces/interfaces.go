package interfaces

import "dddstructure/proto"

// Service defines the main business logic service interface struct that will
// be used between services to call each other.
type Service struct {
	Merchant    Merchant
	User        User
	Invoice     Invoice
	Transaction Transaction
}

// NewServiceParams defines the new service params.
type NewServiceParams struct {
	Merchant    Merchant
	User        User
	Invoice     Invoice
	Transaction Transaction
}

// NewService creates a new service.
func NewService(params NewServiceParams) *Service {
	return &Service{
		Merchant:    params.Merchant,
		User:        params.User,
		Invoice:     params.Invoice,
		Transaction: params.Transaction,
	}
}

// Merchant defines the merchant service.
type Merchant interface {
	Create(m *proto.Merchant) (*proto.Merchant, error)
	GetByID(id uint) (*proto.Merchant, error)
}

// User defines the user service.
type User interface {
	Create(u *proto.User) (*proto.User, error)
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
	Process(t *proto.Transaction) (*proto.Transaction, error)
}
