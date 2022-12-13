package comm

import "dddstructure/proto"

type Transaction interface {
	Process(t *proto.Transaction) (*proto.Transaction, error)
}
