import type { Order, OperationLog } from './types'

// 生成操作日志
const generateLogs = (orderId: string): OperationLog[] => {
  return [
    {
      id: `${orderId}-log-1`,
      orderId,
      operator: 'user001',
      operatorName: '张三',
      role: 'sales',
      operationType: 'create_order',
      operationName: '创建订单',
      beforeData: null,
      afterData: { requiredQuantity: 1000 },
      description: '创建订单，需求数量 1000 件',
      createdAt: '2025-12-15 09:30:00'
    },
    {
      id: `${orderId}-log-2`,
      orderId,
      operator: 'user002',
      operatorName: '李四',
      role: 'follower',
      operationType: 'update_fabric_input',
      operationName: '更新胚布投入',
      beforeData: { completedQuantity: 0 },
      afterData: { completedQuantity: 500 },
      description: '胚布投入数量：0 → 500',
      createdAt: '2025-12-16 10:15:00'
    },
    {
      id: `${orderId}-log-3`,
      orderId,
      operator: 'user002',
      operatorName: '李四',
      role: 'follower',
      operationType: 'update_production',
      operationName: '更新生产进度',
      beforeData: { completedQuantity: 0 },
      afterData: { completedQuantity: 300 },
      description: '生产完成数量：0 → 300',
      createdAt: '2025-12-17 14:20:00'
    },
    {
      id: `${orderId}-log-4`,
      orderId,
      operator: 'user003',
      operatorName: '王五',
      role: 'warehouse',
      operationType: 'update_warehouse_check',
      operationName: '更新验货进度',
      beforeData: { completedQuantity: 0 },
      afterData: { completedQuantity: 250 },
      description: '验货完成数量：0 → 250',
      createdAt: '2025-12-18 11:30:00'
    }
  ]
}

// 生成有次品的日志（包含回修）
const generateLogsWithDefect = (orderId: string): OperationLog[] => {
  return [
    {
      id: `${orderId}-log-1`,
      orderId,
      operator: 'user001',
      operatorName: '张三',
      role: 'sales',
      operationType: 'create_order',
      operationName: '创建订单',
      beforeData: null,
      afterData: { requiredQuantity: 2000 },
      description: '创建订单，需求数量 2000 件',
      createdAt: '2025-12-10 09:00:00'
    },
    {
      id: `${orderId}-log-2`,
      orderId,
      operator: 'user002',
      operatorName: '李四',
      role: 'follower',
      operationType: 'update_fabric_input',
      operationName: '更新胚布投入',
      beforeData: { completedQuantity: 0 },
      afterData: { completedQuantity: 2200 },
      description: '胚布投入数量：0 → 2200',
      createdAt: '2025-12-11 10:00:00'
    },
    {
      id: `${orderId}-log-3`,
      orderId,
      operator: 'user002',
      operatorName: '李四',
      role: 'follower',
      operationType: 'update_production',
      operationName: '更新生产进度',
      beforeData: { completedQuantity: 0 },
      afterData: { completedQuantity: 1800 },
      description: '生产完成数量：0 → 1800',
      createdAt: '2025-12-14 16:00:00'
    },
    {
      id: `${orderId}-log-4`,
      orderId,
      operator: 'user003',
      operatorName: '王五',
      role: 'warehouse',
      operationType: 'update_warehouse_check',
      operationName: '更新验货进度',
      beforeData: { completedQuantity: 0 },
      afterData: { completedQuantity: 1600 },
      description: '验货完成数量：0 → 1600',
      createdAt: '2025-12-16 14:00:00'
    },
    {
      id: `${orderId}-log-5`,
      orderId,
      operator: 'user003',
      operatorName: '王五',
      role: 'warehouse',
      operationType: 'add_defect',
      operationName: '录入次品',
      beforeData: { defectQuantity: 0 },
      afterData: { defectQuantity: 200 },
      description: '发现次品 200 件',
      createdAt: '2025-12-16 14:30:00'
    },
    {
      id: `${orderId}-log-6`,
      orderId,
      operator: 'system',
      operatorName: '系统',
      role: 'system',
      operationType: 'generate_rework',
      operationName: '生成回修进度',
      beforeData: null,
      afterData: { targetQuantity: 200, completedQuantity: 0 },
      description: '自动生成回修进度，目标数量 200 件',
      createdAt: '2025-12-16 14:30:01'
    },
    {
      id: `${orderId}-log-7`,
      orderId,
      operator: 'user002',
      operatorName: '李四',
      role: 'follower',
      operationType: 'update_rework',
      operationName: '更新回修进度',
      beforeData: { completedQuantity: 0 },
      afterData: { completedQuantity: 150 },
      description: '回修完成数量：0 → 150',
      createdAt: '2025-12-18 10:00:00'
    }
  ]
}

// 计算订单整体进度
const calculateOverallProgress = (order: Order): number => {
  const existingItems = order.progressItems.filter(item => item.exists)
  if (existingItems.length === 0) return 0
  const totalProgress = existingItems.reduce((sum, item) => sum + item.progress, 0)
  return Math.round(totalProgress / existingItems.length)
}

// 模拟订单数据
export const mockOrders: Order[] = [
  {
    id: 'order-001',
    orderNo: 'ORD-2025-12-001',
    clientName: '华润纺织有限公司',
    productName: '纯棉府绸染色布',
    productCode: 'PC-40S-001',
    productHistoryShrinkage: 8.5, // 历史缩率 8.5%
    requiredQuantity: 1000,
    createdAt: '2025-12-15 09:30:00',
    updatedAt: '2025-12-18 11:30:00',
    status: 'in_progress',
    progressItems: [
      {
        type: 'fabric_input',
        name: '胚布投入进度',
        targetQuantity: 1085, // 考虑缩率：1000 * (1 + 8.5%) = 1085
        completedQuantity: 500,
        progress: 46, // 500 / 1085 ≈ 46%
        exists: true,
        icon: 'lucide:package-open',
        color: '#409EFF'
      },
      {
        type: 'production',
        name: '生产进度',
        targetQuantity: 1000,
        completedQuantity: 300,
        progress: 30, // 300 / 1000 = 30%
        exists: true,
        icon: 'lucide:factory',
        color: '#67C23A'
      },
      {
        type: 'warehouse_check',
        name: '仓库验货进度',
        targetQuantity: 1000,
        completedQuantity: 250,
        progress: 25, // 250 / 1000 = 25%
        exists: true,
        icon: 'lucide:clipboard-check',
        color: '#E6A23C'
      },
      {
        type: 'rework',
        name: '回修进度',
        targetQuantity: 0,
        completedQuantity: 0,
        progress: 0,
        exists: false, // 无次品，不存在
        icon: 'lucide:wrench',
        color: '#F56C6C'
      }
    ],
    overallProgress: 0, // 后面计算
    operationLogs: generateLogs('order-001')
  },
  {
    id: 'order-002',
    orderNo: 'ORD-2025-12-002',
    clientName: '天虹纺织集团',
    productName: '涤纶混纺印花布',
    productCode: 'PC-65T35C-002',
    productHistoryShrinkage: 6.2,
    requiredQuantity: 2000,
    createdAt: '2025-12-10 09:00:00',
    updatedAt: '2025-12-18 10:00:00',
    status: 'in_progress',
    progressItems: [
      {
        type: 'fabric_input',
        name: '胚布投入进度',
        targetQuantity: 2124, // 2000 * (1 + 6.2%)
        completedQuantity: 2200,
        progress: 100, // 已投入足够
        exists: true,
        icon: 'lucide:package-open',
        color: '#409EFF'
      },
      {
        type: 'production',
        name: '生产进度',
        targetQuantity: 2000,
        completedQuantity: 1800,
        progress: 90, // 1800 / 2000 = 90%
        exists: true,
        icon: 'lucide:factory',
        color: '#67C23A'
      },
      {
        type: 'warehouse_check',
        name: '仓库验货进度',
        targetQuantity: 2000,
        completedQuantity: 1600,
        progress: 80, // 1600 / 2000 = 80%
        exists: true,
        icon: 'lucide:clipboard-check',
        color: '#E6A23C'
      },
      {
        type: 'rework',
        name: '回修进度',
        targetQuantity: 200, // 次品数量
        completedQuantity: 150,
        progress: 75, // 150 / 200 = 75%
        exists: true, // 有次品，存在
        icon: 'lucide:wrench',
        color: '#F56C6C'
      }
    ],
    overallProgress: 0,
    operationLogs: generateLogsWithDefect('order-002')
  },
  {
    id: 'order-003',
    orderNo: 'ORD-2025-12-003',
    clientName: '锦兴服饰有限公司',
    productName: '精梳棉双面布',
    productCode: 'PC-40S-JMC-003',
    productHistoryShrinkage: 7.8,
    requiredQuantity: 1500,
    createdAt: '2025-12-12 10:00:00',
    updatedAt: '2025-12-19 09:00:00',
    status: 'completed',
    progressItems: [
      {
        type: 'fabric_input',
        name: '胚布投入进度',
        targetQuantity: 1617, // 1500 * (1 + 7.8%)
        completedQuantity: 1617,
        progress: 100,
        exists: true,
        icon: 'lucide:package-open',
        color: '#409EFF'
      },
      {
        type: 'production',
        name: '生产进度',
        targetQuantity: 1500,
        completedQuantity: 1500,
        progress: 100,
        exists: true,
        icon: 'lucide:factory',
        color: '#67C23A'
      },
      {
        type: 'warehouse_check',
        name: '仓库验货进度',
        targetQuantity: 1500,
        completedQuantity: 1500,
        progress: 100,
        exists: true,
        icon: 'lucide:clipboard-check',
        color: '#E6A23C'
      },
      {
        type: 'rework',
        name: '回修进度',
        targetQuantity: 0,
        completedQuantity: 0,
        progress: 0,
        exists: false, // 无次品
        icon: 'lucide:wrench',
        color: '#F56C6C'
      }
    ],
    overallProgress: 0,
    operationLogs: [
      {
        id: 'order-003-log-1',
        orderId: 'order-003',
        operator: 'user001',
        operatorName: '张三',
        role: 'sales',
        operationType: 'create_order',
        operationName: '创建订单',
        beforeData: null,
        afterData: { requiredQuantity: 1500 },
        description: '创建订单，需求数量 1500 件',
        createdAt: '2025-12-12 10:00:00'
      }
    ]
  },
  {
    id: 'order-004',
    orderNo: 'ORD-2025-12-004',
    clientName: '美佳服饰',
    productName: '雪纺印花面料',
    productCode: 'PC-XF-004',
    productHistoryShrinkage: 5.5,
    requiredQuantity: 800,
    createdAt: '2025-12-18 14:00:00',
    updatedAt: '2025-12-19 08:00:00',
    status: 'in_progress',
    progressItems: [
      {
        type: 'fabric_input',
        name: '胚布投入进度',
        targetQuantity: 844, // 800 * (1 + 5.5%)
        completedQuantity: 200,
        progress: 24,
        exists: true,
        icon: 'lucide:package-open',
        color: '#409EFF'
      },
      {
        type: 'production',
        name: '生产进度',
        targetQuantity: 800,
        completedQuantity: 0,
        progress: 0,
        exists: true,
        icon: 'lucide:factory',
        color: '#67C23A'
      },
      {
        type: 'warehouse_check',
        name: '仓库验货进度',
        targetQuantity: 800,
        completedQuantity: 0,
        progress: 0,
        exists: true,
        icon: 'lucide:clipboard-check',
        color: '#E6A23C'
      },
      {
        type: 'rework',
        name: '回修进度',
        targetQuantity: 0,
        completedQuantity: 0,
        progress: 0,
        exists: false,
        icon: 'lucide:wrench',
        color: '#F56C6C'
      }
    ],
    overallProgress: 0,
    operationLogs: [
      {
        id: 'order-004-log-1',
        orderId: 'order-004',
        operator: 'user001',
        operatorName: '张三',
        role: 'sales',
        operationType: 'create_order',
        operationName: '创建订单',
        beforeData: null,
        afterData: { requiredQuantity: 800 },
        description: '创建订单，需求数量 800 件',
        createdAt: '2025-12-18 14:00:00'
      },
      {
        id: 'order-004-log-2',
        orderId: 'order-004',
        operator: 'user002',
        operatorName: '李四',
        role: 'follower',
        operationType: 'update_fabric_input',
        operationName: '更新胚布投入',
        beforeData: { completedQuantity: 0 },
        afterData: { completedQuantity: 200 },
        description: '胚布投入数量：0 → 200',
        createdAt: '2025-12-19 08:00:00'
      }
    ]
  }
]

// 计算每个订单的整体进度
mockOrders.forEach(order => {
  order.overallProgress = calculateOverallProgress(order)
})
