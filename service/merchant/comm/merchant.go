package comm

import "dddstructure/proto"

type Merchant interface {
	Create(m *proto.Merchant) (*proto.Merchant, error)
	GetByID(id uint) (*proto.Merchant, error)
}
