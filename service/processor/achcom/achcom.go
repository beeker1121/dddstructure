package achcom

import (
	"fmt"

	"dddstructure/proto"
)

type ACHCom struct{}

func (ac *ACHCom) Process(t *proto.Transaction) error {
	fmt.Println("Processing transaction via ACHCom...")
	return nil
}

func (ac *ACHCom) Refund() error {
	return nil
}
