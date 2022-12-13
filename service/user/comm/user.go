package comm

import "dddstructure/proto"

type User interface {
	Create(u *proto.User) (*proto.User, error)
	GetByID(id uint) (*proto.User, error)
}
