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
	Create(params *proto.UserCreateParams) (*proto.User, error)
	Login(params *proto.UserLoginParams) (*proto.User, error)
	GetByID(id uint) (*proto.User, error)
	Update(params *proto.UserUpdateParams) (*proto.User, error)
}

// Invoice defines the invoice service.
type Invoice interface {
	Create(params *proto.InvoiceCreateParams) (*proto.Invoice, error)
	Get(params *proto.InvoiceGetParams) ([]*proto.Invoice, error)
	GetCount(params *proto.InvoiceGetParams) (uint, error)
	GetByID(id uint) (*proto.Invoice, error)
	Update(params *proto.InvoiceUpdateParams) (*proto.Invoice, error)
	Pay(id uint, params *proto.InvoicePayParams) (*proto.Invoice, error)
}

// Transaction defines the transaction service.
type Transaction interface {
	Process(params *proto.TransactionProcessParams) (*proto.Transaction, error)
}
