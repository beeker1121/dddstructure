package interfaces

import "dddstructure/proto"

// Service defines the main business logic service interface struct that will
// be used between services to call each other.
type Service struct {
	Merchant Merchant
	Invoice  Invoice
}

type NewServiceParams struct {
	Merchant Merchant
	Invoice  Invoice
}

// NewService creates a new service.
func NewService(params NewServiceParams) *Service {
	return &Service{
		Merchant: params.Merchant,
		Invoice:  params.Invoice,
	}
}

type Merchant interface {
	Create(m *proto.Merchant) (*proto.Merchant, error)
	GetByID(id uint) (*proto.Merchant, error)
}

type Invoice interface {
	Create(i *proto.Invoice) (*proto.Invoice, error)
	GetByID(id uint) (*proto.Invoice, error)
	Update(i *proto.Invoice) error
	Pay(id uint) (*proto.Invoice, error)
}
