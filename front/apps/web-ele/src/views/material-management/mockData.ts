import type { Material, SupplierQuote, PriceHistory, StockRecord } from './types'

// 原料列表
export const mockMaterials: Material[] = [
  {
    id: 'mat-001',
    code: 'FB-40S-001',
    name: '40支纯棉胚布',
    category: '胚布',
    unit: 'm',
    currentPrice: 12.5,
    status: 'active',
    updatedBy: '张三',
    createdAt: '2025-01-10 09:00:00',
    updatedAt: '2025-12-18 14:30:00'
  },
  {
    id: 'mat-002',
    code: 'FB-60S-002',
    name: '60支涤棉混纺胚布',
    category: '胚布',
    unit: 'm',
    currentPrice: 18.8,
    status: 'active',
    updatedBy: '李四',
    createdAt: '2025-01-15 10:20:00',
    updatedAt: '2025-12-19 09:15:00'
  },
  {
    id: 'mat-003',
    code: 'RAN-001',
    name: '活性染料-红',
    category: '染料',
    unit: 'kg',
    currentPrice: 85.0,
    status: 'active',
    updatedBy: '王五',
    createdAt: '2025-02-01 11:00:00',
    updatedAt: '2025-12-17 16:00:00'
  },
  {
    id: 'mat-004',
    code: 'RAN-002',
    name: '活性染料-蓝',
    category: '染料',
    unit: 'kg',
    currentPrice: 88.0,
    status: 'active',
    updatedBy: '王五',
    createdAt: '2025-02-01 11:00:00',
    updatedAt: '2025-12-19 10:30:00'
  },
  {
    id: 'mat-005',
    code: 'ZJ-001',
    name: '柔软剂',
    category: '助剂',
    unit: 'kg',
    currentPrice: 32.0,
    status: 'active',
    updatedBy: '赵六',
    createdAt: '2025-02-10 14:00:00',
    updatedAt: '2025-12-18 11:00:00'
  },
  {
    id: 'mat-006',
    code: 'ZJ-002',
    name: '固色剂',
    category: '助剂',
    unit: 'kg',
    currentPrice: 28.5,
    status: 'active',
    updatedBy: '赵六',
    createdAt: '2025-02-10 14:00:00',
    updatedAt: '2025-12-16 15:20:00'
  }
]

// 供应商报价记录（按原料分组，从最新到最旧）
export const mockSupplierQuotes: Record<string, SupplierQuote[]> = {
  'mat-001': [
    {
      id: 'sq-001-1',
      materialId: 'mat-001',
      supplierName: '华润纺织',
      price: 12.5,
      quotedBy: '张三',
      quotedAt: '2025-12-18 14:30:00',
      remark: '年末补货价格'
    },
    {
      id: 'sq-001-2',
      materialId: 'mat-001',
      supplierName: '天虹纺织',
      price: 12.8,
      quotedBy: '李四',
      quotedAt: '2025-12-10 10:00:00'
    },
    {
      id: 'sq-001-3',
      materialId: 'mat-001',
      supplierName: '华润纺织',
      price: 13.0,
      quotedBy: '张三',
      quotedAt: '2025-11-25 09:20:00',
      remark: '11月报价'
    },
    {
      id: 'sq-001-4',
      materialId: 'mat-001',
      supplierName: '锦绣纺织',
      price: 12.2,
      quotedBy: '王五',
      quotedAt: '2025-11-15 16:00:00'
    }
  ],
  'mat-002': [
    {
      id: 'sq-002-1',
      materialId: 'mat-002',
      supplierName: '天虹纺织',
      price: 18.8,
      quotedBy: '李四',
      quotedAt: '2025-12-19 09:15:00'
    },
    {
      id: 'sq-002-2',
      materialId: 'mat-002',
      supplierName: '华润纺织',
      price: 19.2,
      quotedBy: '张三',
      quotedAt: '2025-12-05 14:00:00'
    },
    {
      id: 'sq-002-3',
      materialId: 'mat-002',
      supplierName: '天虹纺织',
      price: 19.0,
      quotedBy: '李四',
      quotedAt: '2025-11-20 11:30:00'
    }
  ],
  'mat-003': [
    {
      id: 'sq-003-1',
      materialId: 'mat-003',
      supplierName: '德美化工',
      price: 85.0,
      quotedBy: '王五',
      quotedAt: '2025-12-17 16:00:00'
    },
    {
      id: 'sq-003-2',
      materialId: 'mat-003',
      supplierName: '华美化工',
      price: 86.5,
      quotedBy: '赵六',
      quotedAt: '2025-12-10 10:00:00'
    },
    {
      id: 'sq-003-3',
      materialId: 'mat-003',
      supplierName: '德美化工',
      price: 82.0,
      quotedBy: '王五',
      quotedAt: '2025-11-28 13:20:00'
    }
  ],
  'mat-004': [
    {
      id: 'sq-004-1',
      materialId: 'mat-004',
      supplierName: '德美化工',
      price: 88.0,
      quotedBy: '王五',
      quotedAt: '2025-12-19 10:30:00'
    },
    {
      id: 'sq-004-2',
      materialId: 'mat-004',
      supplierName: '华美化工',
      price: 89.0,
      quotedBy: '赵六',
      quotedAt: '2025-12-08 15:00:00'
    }
  ],
  'mat-005': [
    {
      id: 'sq-005-1',
      materialId: 'mat-005',
      supplierName: '长虹化工',
      price: 32.0,
      quotedBy: '赵六',
      quotedAt: '2025-12-18 11:00:00'
    },
    {
      id: 'sq-005-2',
      materialId: 'mat-005',
      supplierName: '东方化工',
      price: 33.5,
      quotedBy: '李四',
      quotedAt: '2025-12-01 09:30:00'
    }
  ],
  'mat-006': [
    {
      id: 'sq-006-1',
      materialId: 'mat-006',
      supplierName: '长虹化工',
      price: 28.5,
      quotedBy: '赵六',
      quotedAt: '2025-12-16 15:20:00'
    },
    {
      id: 'sq-006-2',
      materialId: 'mat-006',
      supplierName: '东方化工',
      price: 29.0,
      quotedBy: '李四',
      quotedAt: '2025-11-30 10:00:00'
    }
  ]
}

// 历史价格数据（按原料分组）
export const mockPriceHistory: Record<string, PriceHistory[]> = {
  'mat-001': [
    {
      id: 'ph-001-1',
      materialId: 'mat-001',
      price: 12.5,
      effectiveDate: '2025-12-01',
      supplier: '华润纺织',
      purchaseQuantity: 5000,
      remark: '年末补货'
    },
    {
      id: 'ph-001-2',
      materialId: 'mat-001',
      price: 12.8,
      effectiveDate: '2025-11-01',
      supplier: '华润纺织',
      purchaseQuantity: 3000
    },
    {
      id: 'ph-001-3',
      materialId: 'mat-001',
      price: 12.2,
      effectiveDate: '2025-10-01',
      supplier: '华润纺织',
      purchaseQuantity: 4000
    },
    {
      id: 'ph-001-4',
      materialId: 'mat-001',
      price: 13.0,
      effectiveDate: '2025-09-01',
      supplier: '华润纺织',
      purchaseQuantity: 6000,
      remark: '旺季涨价'
    }
  ],
  'mat-002': [
    {
      id: 'ph-002-1',
      materialId: 'mat-002',
      price: 18.8,
      effectiveDate: '2025-12-01',
      supplier: '天虹纺织',
      purchaseQuantity: 2000
    },
    {
      id: 'ph-002-2',
      materialId: 'mat-002',
      price: 19.2,
      effectiveDate: '2025-11-01',
      supplier: '天虹纺织',
      purchaseQuantity: 2500
    },
    {
      id: 'ph-002-3',
      materialId: 'mat-002',
      price: 18.5,
      effectiveDate: '2025-10-01',
      supplier: '天虹纺织',
      purchaseQuantity: 3000
    }
  ],
  'mat-003': [
    {
      id: 'ph-003-1',
      materialId: 'mat-003',
      price: 85.0,
      effectiveDate: '2025-12-01',
      supplier: '德美化工',
      purchaseQuantity: 100
    },
    {
      id: 'ph-003-2',
      materialId: 'mat-003',
      price: 82.0,
      effectiveDate: '2025-11-01',
      supplier: '德美化工',
      purchaseQuantity: 150
    },
    {
      id: 'ph-003-3',
      materialId: 'mat-003',
      price: 88.0,
      effectiveDate: '2025-10-01',
      supplier: '德美化工',
      purchaseQuantity: 120
    }
  ]
}

// 库存记录数据（按原料分组）
export const mockStockRecords: Record<string, StockRecord[]> = {
  'mat-001': [
    {
      id: 'sr-001-1',
      materialId: 'mat-001',
      type: 'in',
      quantity: 5000,
      unitPrice: 12.5,
      totalValue: 62500,
      operator: '张三',
      remark: '年末补货',
      createdAt: '2025-12-18 14:30:00'
    },
    {
      id: 'sr-001-2',
      materialId: 'mat-001',
      type: 'out',
      quantity: 2000,
      unitPrice: 12.5,
      totalValue: 25000,
      relatedOrder: 'ORD-2025-12-001',
      operator: '李四',
      createdAt: '2025-12-17 10:00:00'
    },
    {
      id: 'sr-001-3',
      materialId: 'mat-001',
      type: 'in',
      quantity: 3000,
      unitPrice: 12.8,
      totalValue: 38400,
      operator: '张三',
      createdAt: '2025-11-28 09:20:00'
    },
    {
      id: 'sr-001-4',
      materialId: 'mat-001',
      type: 'out',
      quantity: 1500,
      unitPrice: 12.8,
      totalValue: 19200,
      relatedOrder: 'ORD-2025-11-025',
      operator: '李四',
      createdAt: '2025-11-25 15:30:00'
    }
  ],
  'mat-002': [
    {
      id: 'sr-002-1',
      materialId: 'mat-002',
      type: 'in',
      quantity: 2000,
      unitPrice: 18.8,
      totalValue: 37600,
      operator: '张三',
      createdAt: '2025-12-19 09:15:00'
    },
    {
      id: 'sr-002-2',
      materialId: 'mat-002',
      type: 'out',
      quantity: 1000,
      unitPrice: 18.8,
      totalValue: 18800,
      relatedOrder: 'ORD-2025-12-002',
      operator: '李四',
      createdAt: '2025-12-15 11:00:00'
    }
  ],
  'mat-003': [
    {
      id: 'sr-003-1',
      materialId: 'mat-003',
      type: 'in',
      quantity: 100,
      unitPrice: 85.0,
      totalValue: 8500,
      operator: '王五',
      createdAt: '2025-12-17 16:00:00'
    },
    {
      id: 'sr-003-2',
      materialId: 'mat-003',
      type: 'out',
      quantity: 50,
      unitPrice: 85.0,
      totalValue: 4250,
      relatedOrder: 'ORD-2025-12-001',
      operator: '赵六',
      createdAt: '2025-12-16 13:30:00'
    }
  ]
}
