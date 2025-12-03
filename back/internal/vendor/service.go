package vendor

import (
	"strconv"
	"back/pkg/es"
)

type VendorService struct {
	vendorRepo *VendorRepo
	esSync     *es.ESSync
}

func NewVendorService(vendorRepo *VendorRepo, esSync *es.ESSync) *VendorService {
	return &VendorService{
		vendorRepo: vendorRepo,
		esSync:     esSync,
	}
}

func (s *VendorService) Create(req *CreateVendorRequest) (*Vendor, error) {
	vendor := &Vendor{
		Name:    req.Name,
		Contact: req.Contact,
		Phone:   req.Phone,
		Email:   req.Email,
		Address: req.Address,
	}

	if err := s.vendorRepo.Create(vendor); err != nil {
		return nil, err
	}

	// 异步同步到 ES
	s.esSync.Index(vendor)

	return vendor, nil
}

func (s *VendorService) Get(id uint) (*Vendor, error) {
	return s.vendorRepo.GetByID(id)
}

func (s *VendorService) List() ([]Vendor, error) {
	return s.vendorRepo.List()
}

func (s *VendorService) Update(id uint, req *UpdateVendorRequest) error {
	vendor, err := s.vendorRepo.GetByID(id)
	if err != nil {
		return err
	}

	data := make(map[string]interface{})
	if req.Name != "" {
		data["name"] = req.Name
		vendor.Name = req.Name
	}
	if req.Contact != "" {
		data["contact"] = req.Contact
		vendor.Contact = req.Contact
	}
	if req.Phone != "" {
		data["phone"] = req.Phone
		vendor.Phone = req.Phone
	}
	if req.Email != "" {
		data["email"] = req.Email
		vendor.Email = req.Email
	}
	if req.Address != "" {
		data["address"] = req.Address
		vendor.Address = req.Address
	}

	if err := s.vendorRepo.Update(id, data); err != nil {
		return err
	}

	// 异步同步到 ES
	s.esSync.Update(vendor)

	return nil
}

func (s *VendorService) Delete(id uint) error {
	if err := s.vendorRepo.Delete(id); err != nil {
		return err
	}

	// 异步删除 ES 文档
	s.esSync.Delete("vendors", strconv.Itoa(int(id)))

	return nil
}