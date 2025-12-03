package order

import (
	"strconv"
	"back/pkg/es"
)

type OrderService struct {
	orderRepo *OrderRepo
	esSync    *es.ESSync
}

func NewOrderService(orderRepo *OrderRepo, esSync *es.ESSync) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		esSync:    esSync,
	}
}

func (s *OrderService) Create(req *CreateOrderRequest) (*Order, error) {
	order := &Order{
		OrderNo:    req.OrderNo,
		ClientID:   req.ClientID,
		ProductID:  req.ProductID,
		Quantity:   req.Quantity,
		UnitPrice:  req.UnitPrice,
		TotalPrice: req.Quantity * req.UnitPrice,
		Status:     "pending",
		CreatedBy:  req.CreatedBy,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	// 异步同步到 ES
	s.esSync.Index(order)

	return order, nil
}

func (s *OrderService) Get(id uint) (*Order, error) {
	return s.orderRepo.GetByID(id)
}

func (s *OrderService) List() ([]Order, error) {
	return s.orderRepo.List()
}

func (s *OrderService) Update(id uint, req *UpdateOrderRequest) error {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return err
	}

	data := make(map[string]interface{})
	if req.Status != "" {
		data["status"] = req.Status
		order.Status = req.Status
	}
	if req.Quantity > 0 {
		data["quantity"] = req.Quantity
		order.Quantity = req.Quantity
		order.TotalPrice = req.Quantity * order.UnitPrice
		data["total_price"] = order.TotalPrice
	}

	if err := s.orderRepo.Update(id, data); err != nil {
		return err
	}

	// 异步同步到 ES
	s.esSync.Update(order)

	return nil
}

func (s *OrderService) Delete(id uint) error {
	if err := s.orderRepo.Delete(id); err != nil {
		return err
	}

	// 异步删除 ES 文档
	s.esSync.Delete("orders", strconv.Itoa(int(id)))

	return nil
}