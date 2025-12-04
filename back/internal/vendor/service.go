package vendor

import (
	"context"
	"strconv"
	
	"back/pkg/es"
	"back/pkg/repo"
	"gorm.io/gorm"
)

type VendorService struct {
	db     *gorm.DB
	esSync *es.ESSync
}

func NewVendorService(db *gorm.DB, esSync *es.ESSync) *VendorService {
	return &VendorService{
		db:     db,
		esSync: esSync,
	}
}

func (s *VendorService) Create(ctx context.Context, req *CreateVendorRequest) (*Vendor, error) {
	vendor := &Vendor{
		Name:    req.Name,
		Contact: req.Contact,
		Phone:   req.Phone,
		Email:   req.Email,
		Address: req.Address,
	}

	vendorRepo := repo.NewRepo[Vendor](s.db)
	if err := vendorRepo.Create(ctx, vendor); err != nil {
		return nil, err
	}

	s.esSync.Index(vendor)
	return vendor, nil
}

func (s *VendorService) Get(ctx context.Context, id uint) (*Vendor, error) {
	vendorRepo := repo.NewRepo[Vendor](s.db)
	return vendorRepo.GetByID(ctx, id)
}

func (s *VendorService) List(ctx context.Context, limit, offset int) ([]Vendor, error) {
	vendorRepo := repo.NewRepo[Vendor](s.db)
	return vendorRepo.List(ctx, limit, offset)
}

func (s *VendorService) Update(ctx context.Context, id uint, req *UpdateVendorRequest) error {
	vendorRepo := repo.NewRepo[Vendor](s.db)
	
	vendor, err := vendorRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	fields := make(map[string]interface{})
	if req.Name != "" {
		fields["name"] = req.Name
		vendor.Name = req.Name
	}
	if req.Contact != "" {
		fields["contact"] = req.Contact
		vendor.Contact = req.Contact
	}
	if req.Phone != "" {
		fields["phone"] = req.Phone
		vendor.Phone = req.Phone
	}
	if req.Email != "" {
		fields["email"] = req.Email
		vendor.Email = req.Email
	}
	if req.Address != "" {
		fields["address"] = req.Address
		vendor.Address = req.Address
	}

	if len(fields) == 0 {
		return nil
	}

	if err := vendorRepo.UpdateFields(ctx, id, fields); err != nil {
		return err
	}

	s.esSync.Update(vendor)
	return nil
}

func (s *VendorService) Delete(ctx context.Context, id uint) error {
	vendorRepo := repo.NewRepo[Vendor](s.db)
	
	if err := vendorRepo.Delete(ctx, id); err != nil {
		return err
	}

	s.esSync.Delete("vendors", strconv.Itoa(int(id)))
	return nil
}