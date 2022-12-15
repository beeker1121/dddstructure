package dep

import "dddstructure/service/processor/comm"

var (
	Processor comm.Processor
)

func RegisterProcessor(i comm.Processor) {
	Processor = i
}
