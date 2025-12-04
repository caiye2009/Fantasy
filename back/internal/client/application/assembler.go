package application

import "back/internal/client/domain"

// ToClient DTO → Domain Model
func ToClient(req *CreateClientRequest) *domain.Client {
	return &domain.Client{
		Name:    req.Name,
		Contact: req.Contact,
		Phone:   req.Phone,
		Email:   req.Email,
		Address: req.Address,
	}
}

// ToClientResponse Domain Model → DTO
func ToClientResponse(c *domain.Client) *ClientResponse {
	return &ClientResponse{
		ID:        c.ID,
		Name:      c.Name,
		Contact:   c.Contact,
		Phone:     c.Phone,
		Email:     c.Email,
		Address:   c.Address,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// ToClientListResponse Domain Models → List DTO
func ToClientListResponse(clients []*domain.Client, total int64) *ClientListResponse {
	responses := make([]*ClientResponse, len(clients))
	for i, c := range clients {
		responses[i] = ToClientResponse(c)
	}
	
	return &ClientListResponse{
		Total:   total,
		Clients: responses,
	}
}