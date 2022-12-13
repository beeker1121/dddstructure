package dep

import "dddstructure/service/transaction/comm"

var (
	Transaction comm.Transaction
)

func RegisterTransaction(t comm.Transaction) {
	Transaction = t
}
