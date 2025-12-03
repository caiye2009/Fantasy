package client

import (
	"strconv"
	"back/pkg/es"
)

type ClientService struct {
	clientRepo *ClientRepo
	esSync     *es.ESSync
}

func NewClientService(clientRepo *ClientRepo, esSync *es.ESSync) *ClientService {
	return &ClientService{
		clientRepo: clientRepo,
		esSync:     esSync,
	}
}

func (s *ClientService) Create(req *CreateClientRequest) (*Client, error) {
	client := &Client{
		Name:    req.Name,
		Contact: req.Contact,
		Phone:   req.Phone,
		Email:   req.Email,
		Address: req.Address,
	}

	if err := s.clientRepo.Create(client); err != nil {
		return nil, err
	}

	// 异步同步到 ES
	s.esSync.Index(client)

	return client, nil
}

func (s *ClientService) Get(id uint) (*Client, error) {
	return s.clientRepo.GetByID(id)
}

func (s *ClientService) List() ([]Client, error) {
	return s.clientRepo.List()
}

func (s *ClientService) Update(id uint, req *UpdateClientRequest) error {
	client, err := s.clientRepo.GetByID(id)
	if err != nil {
		return err
	}

	data := make(map[string]interface{})
	if req.Name != "" {
		data["name"] = req.Name
		client.Name = req.Name
	}
	if req.Contact != "" {
		data["contact"] = req.Contact
		client.Contact = req.Contact
	}
	if req.Phone != "" {
		data["phone"] = req.Phone
		client.Phone = req.Phone
	}
	if req.Email != "" {
		data["email"] = req.Email
		client.Email = req.Email
	}
	if req.Address != "" {
		data["address"] = req.Address
		client.Address = req.Address
	}

	if err := s.clientRepo.Update(id, data); err != nil {
		return err
	}

	// 异步同步到 ES
	s.esSync.Update(client)

	return nil
}

func (s *ClientService) Delete(id uint) error {
	if err := s.clientRepo.Delete(id); err != nil {
		return err
	}

	// 异步删除 ES 文档
	s.esSync.Delete("clients", strconv.Itoa(int(id)))

	return nil
}