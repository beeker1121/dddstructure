package proto

type Processor interface {
	Process(p *Transaction) error
	Refund() error
}
