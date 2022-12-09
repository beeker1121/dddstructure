package comm

import "dddstructure/proto"

type Invoice interface {
	Create(i *proto.Invoice) (*proto.Invoice, error)
	GetByID(id uint) (*proto.Invoice, error)
	Update(i *proto.Invoice) error
}
