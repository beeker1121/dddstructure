package dep

import "dddstructure/service/transaction/comm"

var (
	Transaction comm.Transaction
)

func RegisterTransaction(i comm.Transaction) {
	Transaction = i
}
