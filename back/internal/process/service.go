package process

import (
	"strconv"
	"back/pkg/es"
)

type ProcessService struct {
	processRepo *ProcessRepo
	esSync      *es.ESSync
}

func NewProcessService(processRepo *ProcessRepo, esSync *es.ESSync) *ProcessService {
	return &ProcessService{
		processRepo: processRepo,
		esSync:      esSync,
	}
}

func (s *ProcessService) Create(req *CreateProcessRequest) (*Process, error) {
	process := &Process{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.processRepo.Create(process); err != nil {
		return nil, err
	}

	// 异步同步到 ES
	s.esSync.Index(process)

	return process, nil
}

func (s *ProcessService) Get(id uint) (*Process, error) {
	return s.processRepo.GetByID(id)
}

func (s *ProcessService) List() ([]Process, error) {
	return s.processRepo.List()
}

func (s *ProcessService) Update(id uint, req *UpdateProcessRequest) error {
	process, err := s.processRepo.GetByID(id)
	if err != nil {
		return err
	}

	data := make(map[string]interface{})
	if req.Name != "" {
		data["name"] = req.Name
		process.Name = req.Name
	}
	if req.Description != "" {
		data["description"] = req.Description
		process.Description = req.Description
	}

	if err := s.processRepo.Update(id, data); err != nil {
		return err
	}

	// 异步同步到 ES
	s.esSync.Update(process)

	return nil
}

func (s *ProcessService) Delete(id uint) error {
	if err := s.processRepo.Delete(id); err != nil {
		return err
	}

	// 异步删除 ES 文档
	s.esSync.Delete("processes", strconv.Itoa(int(id)))

	return nil
}