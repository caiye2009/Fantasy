package application

import (
	"context"
	"strconv"
	"time"

	"back/internal/order/domain"
	"back/internal/order/infra"
)

// ESSync ES 同步接口
type ESSync interface {
	Index(doc interface{}) error
	Update(doc interface{}) error
	Delete(indexName, docID string) error
}

// OrderService 订单应用服务
type OrderService struct {
	repo   *infra.OrderRepo
	esSync ESSync
}

// NewOrderService 创建订单服务
func NewOrderService(repo *infra.OrderRepo, esSync ESSync) *OrderService {
	return &OrderService{
		repo:   repo,
		esSync: esSync,
	}
}

// Create 创建订单
func (s *OrderService) Create(ctx context.Context, req *CreateOrderRequest) (*OrderResponse, error) {
	// 1. 检查订单编号是否重复
	exists, err := s.repo.ExistsByOrderNo(ctx, req.OrderNo)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrOrderNoDuplicate
	}

	// 2. DTO → Domain Model
	order := &domain.Order{
		OrderNo:   req.OrderNo,
		ClientID:  req.ClientID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		UnitPrice: req.UnitPrice,
		Status:    domain.OrderStatusPending,
		CreatedBy: req.CreatedBy,
	}
	order.CalculateTotalPrice()

	// 3. 领域验证
	if err := order.Validate(); err != nil {
		return nil, err
	}

	// 4. 保存到数据库
	if err := s.repo.Save(ctx, order); err != nil {
		return nil, err
	}

	// 5. 异步同步到 ES
	if s.esSync != nil {
		s.esSync.Index(order)
	}

	// 6. Domain Model → DTO
	return &OrderResponse{
		ID:         order.ID,
		OrderNo:    order.OrderNo,
		ClientID:   order.ClientID,
		ProductID:  order.ProductID,
		Quantity:   order.Quantity,
		UnitPrice:  order.UnitPrice,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
		CreatedBy:  order.CreatedBy,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}, nil
}

// Get 获取订单
func (s *OrderService) Get(ctx context.Context, id uint) (*OrderResponse, error) {
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &OrderResponse{
		ID:         order.ID,
		OrderNo:    order.OrderNo,
		ClientID:   order.ClientID,
		ProductID:  order.ProductID,
		Quantity:   order.Quantity,
		UnitPrice:  order.UnitPrice,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
		CreatedBy:  order.CreatedBy,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}, nil
}

// List 订单列表
func (s *OrderService) List(ctx context.Context, limit, offset int) (*OrderListResponse, error) {
	orders, err := s.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*OrderResponse, len(orders))
	for i, o := range orders {
		responses[i] = &OrderResponse{
			ID:         o.ID,
			OrderNo:    o.OrderNo,
			ClientID:   o.ClientID,
			ProductID:  o.ProductID,
			Quantity:   o.Quantity,
			UnitPrice:  o.UnitPrice,
			TotalPrice: o.TotalPrice,
			Status:     o.Status,
			CreatedBy:  o.CreatedBy,
			CreatedAt:  o.CreatedAt,
			UpdatedAt:  o.UpdatedAt,
		}
	}

	return &OrderListResponse{
		Total:  total,
		Orders: responses,
	}, nil
}

// Update 更新订单
func (s *OrderService) Update(ctx context.Context, id uint, req *UpdateOrderRequest) error {
	// 1. 查询订单
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	
	// 2. 更新字段（通过领域方法）
	if req.Quantity > 0 {
		if err := order.UpdateQuantity(req.Quantity); err != nil {
			return err
		}
	}
	
	if req.UnitPrice >= 0 {
		if err := order.UpdateUnitPrice(req.UnitPrice); err != nil {
			return err
		}
	}
	
	// 状态更新
	if req.Status != "" {
		switch req.Status {
		case domain.OrderStatusConfirmed:
			if err := order.Confirm(); err != nil {
				return err
			}
		case domain.OrderStatusProduction:
			if err := order.StartProduction(); err != nil {
				return err
			}
		case domain.OrderStatusCompleted:
			if err := order.Complete(); err != nil {
				return err
			}
		case domain.OrderStatusCancelled:
			if err := order.Cancel(); err != nil {
				return err
			}
		}
	}
	
	// 3. 验证
	if err := order.Validate(); err != nil {
		return err
	}
	
	// 4. 保存
	if err := s.repo.Update(ctx, order); err != nil {
		return err
	}
	
	// 5. 异步同步到 ES
	if s.esSync != nil {
		s.esSync.Update(order)
	}
	
	return nil
}

// Delete 删除订单
func (s *OrderService) Delete(ctx context.Context, id uint) error {
	// 1. 查询订单
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	
	// 2. 检查是否可以删除（领域规则）
	if !order.CanDelete() {
		return domain.ErrCannotDeleteCompleted
	}
	
	// 3. 删除
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	
	// 4. 异步删除 ES 文档
	if s.esSync != nil {
		s.esSync.Delete("orders", strconv.Itoa(int(id)))
	}
	
	return nil
}

// GetByClientID 根据客户 ID 查询订单
func (s *OrderService) GetByClientID(ctx context.Context, clientID uint) ([]*OrderResponse, error) {
	orders, err := s.repo.FindByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	responses := make([]*OrderResponse, len(orders))
	for i, o := range orders {
		responses[i] = &OrderResponse{
			ID:         o.ID,
			OrderNo:    o.OrderNo,
			ClientID:   o.ClientID,
			ProductID:  o.ProductID,
			Quantity:   o.Quantity,
			UnitPrice:  o.UnitPrice,
			TotalPrice: o.TotalPrice,
			Status:     o.Status,
			CreatedBy:  o.CreatedBy,
			CreatedAt:  o.CreatedAt,
			UpdatedAt:  o.UpdatedAt,
		}
	}

	return responses, nil
}

// GetByProductID 根据产品 ID 查询订单
func (s *OrderService) GetByProductID(ctx context.Context, productID uint) ([]*OrderResponse, error) {
	orders, err := s.repo.FindByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	responses := make([]*OrderResponse, len(orders))
	for i, o := range orders {
		responses[i] = &OrderResponse{
			ID:         o.ID,
			OrderNo:    o.OrderNo,
			ClientID:   o.ClientID,
			ProductID:  o.ProductID,
			Quantity:   o.Quantity,
			UnitPrice:  o.UnitPrice,
			TotalPrice: o.TotalPrice,
			Status:     o.Status,
			CreatedBy:  o.CreatedBy,
			CreatedAt:  o.CreatedAt,
			UpdatedAt:  o.UpdatedAt,
		}
	}

	return responses, nil
}

// GetByStatus 根据状态查询订单
func (s *OrderService) GetByStatus(ctx context.Context, status string, limit, offset int) (*OrderListResponse, error) {
	orders, err := s.repo.FindByStatus(ctx, status, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*OrderResponse, len(orders))
	for i, o := range orders {
		responses[i] = &OrderResponse{
			ID:         o.ID,
			OrderNo:    o.OrderNo,
			ClientID:   o.ClientID,
			ProductID:  o.ProductID,
			Quantity:   o.Quantity,
			UnitPrice:  o.UnitPrice,
			TotalPrice: o.TotalPrice,
			Status:     o.Status,
			CreatedBy:  o.CreatedBy,
			CreatedAt:  o.CreatedAt,
			UpdatedAt:  o.UpdatedAt,
		}
	}

	return &OrderListResponse{
		Total:  total,
		Orders: responses,
	}, nil
}

// ==================== 新版协作工作流方法 ====================

// CreateV2 创建订单（新版，支持协作工作流）
func (s *OrderService) CreateV2(ctx context.Context, req *CreateOrderRequestV2, creatorID uint, creatorName, creatorRole string) (*OrderDetailResponse, error) {
	var result *OrderDetailResponse

	err := s.repo.Transaction(ctx, func(txCtx context.Context) error {
		// 1. 检查订单编号是否重复
		exists, err := s.repo.ExistsByOrderNo(txCtx, req.OrderNo)
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrOrderNoDuplicate
		}

		// 2. 创建订单实体
		order := &domain.Order{
			OrderNo:                 req.OrderNo,
			ClientID:                req.ClientID,
			ProductID:               req.ProductID,
			RequiredQuantity:        req.RequiredQuantity,
			ProductHistoryShrinkage: req.ProductHistoryShrinkage,
			Quantity:                req.RequiredQuantity, // 临时使用，后续由生产助理设定胚布数量
			UnitPrice:               req.UnitPrice,
			Status:                  domain.OrderStatusPending,
			CreatedBy:               creatorID,
		}
		order.CalculateTotalPrice()

		// 3. 领域验证
		if err := order.Validate(); err != nil {
			return err
		}

		// 4. 保存订单
		if err := s.repo.Save(txCtx, order); err != nil {
			return err
		}

		// 5. 创建参与者记录（创建者）
		participant := &domain.OrderParticipant{
			OrderID:  order.ID,
			UserID:   creatorID,
			UserName: creatorName,
			Role:     "creator",
			IsActive: true,
		}
		if err := s.repo.CreateParticipant(txCtx, participant); err != nil {
			return err
		}

		// 6. 初始化4个进度项（全部Exists=false）
		progressTypes := []string{
			domain.ProgressTypeFabricInput,
			domain.ProgressTypeProduction,
			domain.ProgressTypeWarehouseCheck,
			domain.ProgressTypeRework,
		}
		for _, pType := range progressTypes {
			progress := &domain.OrderProgress{
				OrderID:           order.ID,
				Type:              pType,
				TargetQuantity:    0,
				CompletedQuantity: 0,
				Progress:          0,
				Exists:            false,
			}
			if err := s.repo.CreateProgress(txCtx, progress); err != nil {
				return err
			}
		}

		// 7. 记录创建订单事件
		event := createEvent(
			order.ID,
			domain.EventTypeCreateOrder,
			creatorID,
			creatorName,
			creatorRole,
			nil,
			map[string]interface{}{
				"order_no":                   order.OrderNo,
				"required_quantity":          order.RequiredQuantity,
				"product_history_shrinkage":  order.ProductHistoryShrinkage,
			},
			"创建订单",
		)
		if err := s.repo.CreateEvent(txCtx, event); err != nil {
			return err
		}

		// 8. 异步同步到 ES
		if s.esSync != nil {
			s.esSync.Index(order)
		}

		// 9. 获取完整订单详情
		detail, err := s.GetDetail(txCtx, order.ID, creatorID)
		if err != nil {
			return err
		}
		result = detail

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// AssignDepartment 分配部门（生产总监操作）
func (s *OrderService) AssignDepartment(ctx context.Context, orderID uint, req *AssignDepartmentRequest, directorID uint, directorName, directorRole string) error {
	return s.repo.Transaction(ctx, func(txCtx context.Context) error {
		// 1. 查询订单
		order, err := s.repo.FindByID(txCtx, orderID)
		if err != nil {
			return err
		}

		// 2. 检查状态
		if order.Status != domain.OrderStatusPending {
			return domain.ErrInvalidOrderStatus
		}

		// 3. 分配部门（领域方法）
		beforeData := map[string]interface{}{
			"assigned_department": order.AssignedDepartment,
		}

		if err := order.AssignDepartment(req.Department); err != nil {
			return err
		}

		afterData := map[string]interface{}{
			"assigned_department": order.AssignedDepartment,
		}

		// 4. 保存订单
		if err := s.repo.Update(txCtx, order); err != nil {
			return err
		}

		// 5. 创建参与者记录（生产总监）
		participant := &domain.OrderParticipant{
			OrderID:  orderID,
			UserID:   directorID,
			UserName: directorName,
			Role:     "production_director",
			IsActive: true,
		}
		if err := s.repo.CreateParticipant(txCtx, participant); err != nil {
			return err
		}

		// 6. 记录分配部门事件
		event := createEvent(
			orderID,
			domain.EventTypeAssignDepartment,
			directorID,
			directorName,
			directorRole,
			beforeData,
			afterData,
			"分配到"+req.Department+"部门",
		)
		if err := s.repo.CreateEvent(txCtx, event); err != nil {
			return err
		}

		// 7. 异步更新 ES
		if s.esSync != nil {
			s.esSync.Update(order)
		}

		return nil
	})
}

// AssignPersonnel 分配人员（生产助理操作）
func (s *OrderService) AssignPersonnel(ctx context.Context, orderID uint, req *AssignPersonnelRequest, assistantID uint, assistantName, assistantRole string) error {
	return s.repo.Transaction(ctx, func(txCtx context.Context) error {
		// 1. 查询订单
		order, err := s.repo.FindByID(txCtx, orderID)
		if err != nil {
			return err
		}

		// 2. 检查状态
		if order.Status != domain.OrderStatusAssigned {
			return domain.ErrInvalidOrderStatus
		}

		// 3. 更新订单状态为进行中
		if err := order.StartProgress(); err != nil {
			return err
		}

		if err := s.repo.Update(txCtx, order); err != nil {
			return err
		}

		// 4. 创建3个参与者记录（生产助理、生产专员、跟单）
		participants := []*domain.OrderParticipant{
			{
				OrderID:  orderID,
				UserID:   assistantID,
				UserName: assistantName,
				Role:     "production_assistant",
				IsActive: true,
			},
			{
				OrderID:  orderID,
				UserID:   req.ProductionSpecialistID,
				UserName: "", // 需要从用户服务查询，这里简化
				Role:     "production_specialist",
				IsActive: true,
			},
			{
				OrderID:  orderID,
				UserID:   req.OrderCoordinatorID,
				UserName: "", // 需要从用户服务查询，这里简化
				Role:     "order_coordinator",
				IsActive: true,
			},
		}

		for _, p := range participants {
			if err := s.repo.CreateParticipant(txCtx, p); err != nil {
				return err
			}
		}

		// 5. 设定胚布目标数量并激活
		fabricProgress, err := s.repo.FindProgressByType(txCtx, orderID, domain.ProgressTypeFabricInput)
		if err != nil {
			return err
		}

		fabricProgress.SetTarget(req.FabricTargetQuantity)
		fabricProgress.MarkAsExistent()

		if err := s.repo.UpdateProgress(txCtx, fabricProgress); err != nil {
			return err
		}

		// 6. 激活生产进度和验货进度
		productionProgress, err := s.repo.FindProgressByType(txCtx, orderID, domain.ProgressTypeProduction)
		if err != nil {
			return err
		}
		productionProgress.MarkAsExistent()
		productionProgress.SetTarget(order.RequiredQuantity) // 目标是成品数量
		if err := s.repo.UpdateProgress(txCtx, productionProgress); err != nil {
			return err
		}

		warehouseProgress, err := s.repo.FindProgressByType(txCtx, orderID, domain.ProgressTypeWarehouseCheck)
		if err != nil {
			return err
		}
		warehouseProgress.MarkAsExistent()
		warehouseProgress.SetTarget(order.RequiredQuantity) // 目标是成品数量
		if err := s.repo.UpdateProgress(txCtx, warehouseProgress); err != nil {
			return err
		}

		// 7. 记录分配人员事件
		event := createEvent(
			orderID,
			domain.EventTypeAssignPersonnel,
			assistantID,
			assistantName,
			assistantRole,
			nil,
			map[string]interface{}{
				"production_specialist_id": req.ProductionSpecialistID,
				"order_coordinator_id":     req.OrderCoordinatorID,
			},
			"分配生产人员",
		)
		if err := s.repo.CreateEvent(txCtx, event); err != nil {
			return err
		}

		// 8. 记录设定胚布目标数量事件
		fabricEvent := createEvent(
			orderID,
			domain.EventTypeSetFabricTarget,
			assistantID,
			assistantName,
			assistantRole,
			nil,
			map[string]interface{}{
				"fabric_target_quantity": req.FabricTargetQuantity,
			},
			"设定胚布目标数量为"+strconv.FormatFloat(req.FabricTargetQuantity, 'f', 0, 64),
		)
		if err := s.repo.CreateEvent(txCtx, fabricEvent); err != nil {
			return err
		}

		// 9. 异步更新 ES
		if s.esSync != nil {
			s.esSync.Update(order)
		}

		return nil
	})
}

// UpdateProgress 更新进度（跟单/仓管操作）
func (s *OrderService) UpdateProgress(ctx context.Context, orderID uint, progressType string, req *UpdateProgressRequest, operatorID uint, operatorName, operatorRole string) error {
	return s.repo.Transaction(ctx, func(txCtx context.Context) error {
		// 1. 查询进度
		progress, err := s.repo.FindProgressByType(txCtx, orderID, progressType)
		if err != nil {
			return err
		}

		// 2. 检查进度是否存在
		if !progress.Exists {
			return domain.ErrProgressNotExists
		}

		// 3. 记录变更前数据
		beforeData := map[string]interface{}{
			"completed_quantity": progress.CompletedQuantity,
			"progress":           progress.Progress,
		}

		// 4. 更新完成数量
		if err := progress.UpdateCompleted(req.Quantity); err != nil {
			return err
		}

		// 5. 保存进度
		if err := s.repo.UpdateProgress(txCtx, progress); err != nil {
			return err
		}

		// 6. 记录变更后数据
		afterData := map[string]interface{}{
			"completed_quantity": progress.CompletedQuantity,
			"progress":           progress.Progress,
		}

		// 7. 确定事件类型
		var eventType string
		switch progressType {
		case domain.ProgressTypeFabricInput:
			eventType = domain.EventTypeUpdateFabricInput
		case domain.ProgressTypeProduction:
			eventType = domain.EventTypeUpdateProduction
		case domain.ProgressTypeWarehouseCheck:
			eventType = domain.EventTypeUpdateWarehouseCheck
		case domain.ProgressTypeRework:
			eventType = domain.EventTypeUpdateRework
		default:
			eventType = "update_progress"
		}

		// 8. 记录事件
		description := formatQuantityChange(
			getProgressName(progressType),
			beforeData["completed_quantity"].(float64),
			afterData["completed_quantity"].(float64),
		)
		if req.Remark != "" {
			description += "（" + req.Remark + "）"
		}

		event := createEvent(
			orderID,
			eventType,
			operatorID,
			operatorName,
			operatorRole,
			beforeData,
			afterData,
			description,
		)
		if err := s.repo.CreateEvent(txCtx, event); err != nil {
			return err
		}

		return nil
	})
}

// AddDefect 录入次品（仓管操作）
func (s *OrderService) AddDefect(ctx context.Context, orderID uint, req *AddDefectRequest, warehouseID uint, warehouseName, warehouseRole string) error {
	return s.repo.Transaction(ctx, func(txCtx context.Context) error {
		// 1. 查询回修进度
		reworkProgress, err := s.repo.FindProgressByType(txCtx, orderID, domain.ProgressTypeRework)
		if err != nil {
			return err
		}

		// 2. 记录变更前数据
		beforeData := map[string]interface{}{
			"exists":          reworkProgress.Exists,
			"target_quantity": reworkProgress.TargetQuantity,
		}

		// 3. 累加次品数量到回修目标数量
		newTarget := reworkProgress.TargetQuantity + req.DefectQuantity
		reworkProgress.SetTarget(newTarget)

		// 4. 标记回修进度为存在
		if !reworkProgress.Exists {
			reworkProgress.MarkAsExistent()
		}

		// 5. 保存回修进度
		if err := s.repo.UpdateProgress(txCtx, reworkProgress); err != nil {
			return err
		}

		// 6. 记录变更后数据
		afterData := map[string]interface{}{
			"exists":          reworkProgress.Exists,
			"target_quantity": reworkProgress.TargetQuantity,
		}

		// 7. 创建仓管参与者记录（如果不存在）
		participants, err := s.repo.GetActiveParticipants(txCtx, orderID)
		if err != nil {
			return err
		}

		hasWarehouse := false
		for _, p := range participants {
			if p.Role == "warehouse" && p.UserID == warehouseID {
				hasWarehouse = true
				break
			}
		}

		if !hasWarehouse {
			participant := &domain.OrderParticipant{
				OrderID:  orderID,
				UserID:   warehouseID,
				UserName: warehouseName,
				Role:     "warehouse",
				IsActive: true,
			}
			if err := s.repo.CreateParticipant(txCtx, participant); err != nil {
				return err
			}
		}

		// 8. 记录录入次品事件
		event := createEvent(
			orderID,
			domain.EventTypeAddDefect,
			warehouseID,
			warehouseName,
			warehouseRole,
			beforeData,
			afterData,
			"录入次品数量："+strconv.FormatFloat(req.DefectQuantity, 'f', 0, 64)+"（"+req.Remark+"）",
		)
		if err := s.repo.CreateEvent(txCtx, event); err != nil {
			return err
		}

		return nil
	})
}

// GetDetail 获取订单详情（含完整参与者、进度、事件流）
func (s *OrderService) GetDetail(ctx context.Context, orderID uint, userID uint) (*OrderDetailResponse, error) {
	// 1. 查询订单基本信息（含客户和产品名称）
	orderData, err := s.repo.GetOrderWithClientProduct(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// 2. 查询订单完整信息（含关联）
	order, err := s.repo.GetOrderWithDetails(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// 3. 转换进度（添加UI元数据）
	progressResponses := make([]ProgressResponse, len(order.Progresses))
	for i, p := range order.Progresses {
		progressResponses[i] = ProgressResponse{
			ID:                p.ID,
			OrderID:           p.OrderID,
			Type:              p.Type,
			Name:              getProgressName(p.Type),
			TargetQuantity:    p.TargetQuantity,
			CompletedQuantity: p.CompletedQuantity,
			Progress:          p.Progress,
			Exists:            p.Exists,
			Icon:              getProgressIcon(p.Type),
			Color:             getProgressColor(p.Type),
			CreatedAt:         p.CreatedAt,
			UpdatedAt:         p.UpdatedAt,
		}
	}

	// 4. 转换事件
	eventResponses := make([]EventResponse, len(order.Events))
	for i, e := range order.Events {
		eventResponses[i] = EventResponse{
			ID:           e.ID,
			OrderID:      e.OrderID,
			EventType:    e.EventType,
			OperatorID:   e.OperatorID,
			OperatorName: e.OperatorName,
			OperatorRole: e.OperatorRole,
			BeforeData:   e.BeforeData,
			AfterData:    e.AfterData,
			Description:  e.Description,
			CreatedAt:    e.CreatedAt,
		}
	}

	// 5. 计算整体进度
	overallProgress := calculateOverallProgress(order.Progresses)

	// 6. 获取客户名和产品信息
	clientName, _ := orderData["client_name"].(string)
	productName, _ := orderData["product_name"].(string)
	productCode, _ := orderData["product_code"].(string)

	// 7. 构建响应
	return &OrderDetailResponse{
		ID:                      order.ID,
		OrderNo:                 order.OrderNo,
		ClientID:                order.ClientID,
		ClientName:              clientName,
		ProductID:               order.ProductID,
		ProductName:             productName,
		ProductCode:             productCode,
		ProductHistoryShrinkage: order.ProductHistoryShrinkage,
		RequiredQuantity:        order.RequiredQuantity,
		UnitPrice:               order.UnitPrice,
		TotalPrice:              order.TotalPrice,
		Status:                  order.Status,
		AssignedDepartment:      order.AssignedDepartment,
		CreatedAt:               order.CreatedAt,
		UpdatedAt:               order.UpdatedAt,
		ProgressItems:           progressResponses,
		OperationLogs:           eventResponses,
		OverallProgress:         overallProgress,
	}, nil
}

// GetEvents 获取订单事件流
func (s *OrderService) GetEvents(ctx context.Context, orderID uint) ([]EventResponse, error) {
	events, err := s.repo.GetEvents(ctx, orderID)
	if err != nil {
		return nil, err
	}

	responses := make([]EventResponse, len(events))
	for i, e := range events {
		responses[i] = EventResponse{
			ID:           e.ID,
			OrderID:      e.OrderID,
			EventType:    e.EventType,
			OperatorID:   e.OperatorID,
			OperatorName: e.OperatorName,
			OperatorRole: e.OperatorRole,
			BeforeData:   e.BeforeData,
			AfterData:    e.AfterData,
			Description:  e.Description,
			CreatedAt:    e.CreatedAt,
		}
	}

	return responses, nil
}

// ListWithDetail 获取订单列表（含完整详情）
func (s *OrderService) ListWithDetail(ctx context.Context, limit, offset int) (*OrderListDetailResponse, error) {
	// 1. 获取订单基本信息（含客户和产品名称）
	ordersData, total, err := s.repo.GetOrdersWithClientProduct(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// 2. 提取订单ID列表
	orderIDs := make([]uint, len(ordersData))
	for i, data := range ordersData {
		id, _ := data["id"].(uint)
		orderIDs[i] = id
	}

	// 3. 批量获取所有订单的进度
	allProgresses, err := s.repo.GetProgressesByOrderIDs(ctx, orderIDs)
	if err != nil {
		return nil, err
	}
	progressMap := make(map[uint][]domain.OrderProgress)
	for _, p := range allProgresses {
		progressMap[p.OrderID] = append(progressMap[p.OrderID], p)
	}

	// 4. 批量获取所有订单的事件
	allEvents, err := s.repo.GetEventsByOrderIDs(ctx, orderIDs)
	if err != nil {
		return nil, err
	}
	eventMap := make(map[uint][]domain.OrderEvent)
	for _, e := range allEvents {
		eventMap[e.OrderID] = append(eventMap[e.OrderID], e)
	}

	// 5. 组装响应
	orders := make([]*OrderDetailResponse, len(ordersData))
	for i, data := range ordersData {
		id, _ := data["id"].(uint)
		orderNo, _ := data["order_no"].(string)
		clientID, _ := data["client_id"].(uint)
		clientName, _ := data["client_name"].(string)
		productID, _ := data["product_id"].(uint)
		productName, _ := data["product_name"].(string)
		productCode, _ := data["product_code"].(string)
		requiredQuantity, _ := data["required_quantity"].(float64)
		productHistoryShrinkage, _ := data["product_history_shrinkage"].(float64)
		unitPrice, _ := data["unit_price"].(float64)
		totalPrice, _ := data["total_price"].(float64)
		status, _ := data["status"].(string)
		assignedDepartment, _ := data["assigned_department"].(string)
		createdAt, _ := data["created_at"].(time.Time)
		updatedAt, _ := data["updated_at"].(time.Time)

		// 转换进度
		progresses := progressMap[id]
		progressResponses := make([]ProgressResponse, len(progresses))
		for j, p := range progresses {
			progressResponses[j] = ProgressResponse{
				ID:                p.ID,
				OrderID:           p.OrderID,
				Type:              p.Type,
				Name:              getProgressName(p.Type),
				TargetQuantity:    p.TargetQuantity,
				CompletedQuantity: p.CompletedQuantity,
				Progress:          p.Progress,
				Exists:            p.Exists,
				Icon:              getProgressIcon(p.Type),
				Color:             getProgressColor(p.Type),
				CreatedAt:         p.CreatedAt,
				UpdatedAt:         p.UpdatedAt,
			}
		}

		// 转换事件
		events := eventMap[id]
		eventResponses := make([]EventResponse, len(events))
		for j, e := range events {
			eventResponses[j] = EventResponse{
				ID:           e.ID,
				OrderID:      e.OrderID,
				EventType:    e.EventType,
				OperatorID:   e.OperatorID,
				OperatorName: e.OperatorName,
				OperatorRole: e.OperatorRole,
				BeforeData:   e.BeforeData,
				AfterData:    e.AfterData,
				Description:  e.Description,
				CreatedAt:    e.CreatedAt,
			}
		}

		// 计算整体进度
		overallProgress := calculateOverallProgress(progresses)

		orders[i] = &OrderDetailResponse{
			ID:                      id,
			OrderNo:                 orderNo,
			ClientID:                clientID,
			ClientName:              clientName,
			ProductID:               productID,
			ProductName:             productName,
			ProductCode:             productCode,
			ProductHistoryShrinkage: productHistoryShrinkage,
			RequiredQuantity:        requiredQuantity,
			UnitPrice:               unitPrice,
			TotalPrice:              totalPrice,
			Status:                  status,
			AssignedDepartment:      assignedDepartment,
			CreatedAt:               createdAt,
			UpdatedAt:               updatedAt,
			ProgressItems:           progressResponses,
			OperationLogs:           eventResponses,
			OverallProgress:         overallProgress,
		}
	}

	return &OrderListDetailResponse{
		Total:  total,
		Orders: orders,
	}, nil
}