package processor

import (
	"dddstructure/proto"
	"dddstructure/service/processor/achcom"
	"dddstructure/storage"
	"errors"
)

// Service defines the user service.
type Service struct {
	s *storage.Storage
}

// New creates a new service.
func New(s *storage.Storage) *Service {
	return &Service{
		s: s,
	}
}

func (s *Service) GetProcessor(t *proto.Transaction) (proto.Processor, error) {
	if t.ProcessorType == "achcom" {
		proc := &achcom.ACHCom{}
		return proc, nil
	}

	return nil, errors.New("Could not find processor type")
}
