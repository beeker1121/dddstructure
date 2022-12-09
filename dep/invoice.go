package dep

import "dddstructure/service/invoice/comm"

var (
	Invoice comm.Invoice
)

func RegisterInvoice(i comm.Invoice) {
	Invoice = i
}
