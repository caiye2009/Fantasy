package material

import (
	"strconv"
	"back/pkg/es"
)

type MaterialService struct {
	materialRepo *MaterialRepo
	esSync       *es.ESSync
}

func NewMaterialService(materialRepo *MaterialRepo, esSync *es.ESSync) *MaterialService {
	return &MaterialService{
		materialRepo: materialRepo,
		esSync:       esSync,
	}
}

func (s *MaterialService) Create(req *CreateMaterialRequest) (*Material, error) {
	material := &Material{
		Name:        req.Name,
		Spec:        req.Spec,
		Unit:        req.Unit,
		Description: req.Description,
	}

	if err := s.materialRepo.Create(material); err != nil {
		return nil, err
	}

	// 异步同步到 ES
	s.esSync.Index(material)

	return material, nil
}

func (s *MaterialService) Get(id uint) (*Material, error) {
	return s.materialRepo.GetByID(id)
}

func (s *MaterialService) List() ([]Material, error) {
	return s.materialRepo.List()
}

func (s *MaterialService) Update(id uint, req *UpdateMaterialRequest) error {
	material, err := s.materialRepo.GetByID(id)
	if err != nil {
		return err
	}

	data := make(map[string]interface{})
	if req.Name != "" {
		data["name"] = req.Name
		material.Name = req.Name
	}
	if req.Spec != "" {
		data["spec"] = req.Spec
		material.Spec = req.Spec
	}
	if req.Unit != "" {
		data["unit"] = req.Unit
		material.Unit = req.Unit
	}
	if req.Description != "" {
		data["description"] = req.Description
		material.Description = req.Description
	}

	if err := s.materialRepo.Update(id, data); err != nil {
		return err
	}

	// 异步同步到 ES
	s.esSync.Update(material)

	return nil
}

func (s *MaterialService) Delete(id uint) error {
	if err := s.materialRepo.Delete(id); err != nil {
		return err
	}

	// 异步删除 ES 文档
	s.esSync.Delete("materials", strconv.Itoa(int(id)))

	return nil
}