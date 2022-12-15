package processor

import (
	"dddstructure/proto"
	"dddstructure/service/processor/achcom"
	"dddstructure/storage"
)

// idCounter handles increasing the ID.
// var idCounter uint = 1

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

func (s *Service) GetProcessor(t *proto.Transaction) proto.Processor {
	procid := "achcom"
	if procid == "achcom" {
		proc := &achcom.ACHCom{}
		return proc
	}

	return nil
}
