package comm

import "dddstructure/proto"

type Processor interface {
	GetProcessor(t *proto.Transaction) (proto.Processor, error)
}
