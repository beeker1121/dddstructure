package dep

import "dddstructure/service/merchant/comm"

var (
	Merchant comm.Merchant
)

func RegisterMerchant(m comm.Merchant) {
	Merchant = m
}
