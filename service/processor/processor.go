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
	switch t.ProcessorType {
	case "achcom":
		return &achcom.ACHCom{}, nil
	}

	return nil, errors.New("Could not find processor type")
}
