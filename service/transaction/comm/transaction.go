package comm

import "dddstructure/proto"

type Transaction interface {
	Create(i *proto.Transaction) (*proto.Transaction, error)
	GetByID(id uint) (*proto.Transaction, error)
}
