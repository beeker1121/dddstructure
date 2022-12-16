package processor_test

import (
	"database/sql"
	"testing"

	"dddstructure/dep"
	"dddstructure/proto"
	"dddstructure/service"
	"dddstructure/service/processor/achcom"
	"dddstructure/storage/mysql"
)

func TestGetProcessor(t *testing.T) {
	// Create a new MySQL storage implementation.
	store := mysql.New(&sql.DB{})

	// Create a new service.
	serv := service.New(store)

	// Register dependencies.
	dep.RegisterMerchant(serv.Merchant)
	dep.RegisterUser(serv.User)
	dep.RegisterInvoice(serv.Invoice)
	dep.RegisterProcessor(serv.Processor)
	dep.RegisterTransaction(serv.Transaction)

	// Get processor.
	proc, err := serv.Processor.GetProcessor(&proto.Transaction{
		ProcessorType: "achcom",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check type of processor.
	if _, ok := proc.(*achcom.ACHCom); !ok {
		t.Error("Could not type assert prcoessor interface to type 'achcom.ACHCom'")
	}
}
