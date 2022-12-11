package dep

import "dddstructure/service/user/comm"

var (
	User comm.User
)

func RegisterUser(u comm.User) {
	User = u
}
