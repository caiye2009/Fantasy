/**
 * 库存管理模拟数据
 */

// 随机生成日期（最近30天内）
function randomDate(daysAgo = 30): string {
  const now = new Date()
  const past = new Date(now.getTime() - Math.random() * daysAgo * 24 * 60 * 60 * 1000)
  return past.toLocaleString('zh-CN')
}

// 随机生成数量
function randomQuantity(min = 100, max = 5000): number {
  return Math.floor(Math.random() * (max - min + 1)) + min
}

// 随机生成价格
function randomPrice(min = 10, max = 200): number {
  return Number((Math.random() * (max - min) + min).toFixed(2))
}

// 随机选择
function randomChoice<T>(arr: T[]): T {
  return arr[Math.floor(Math.random() * arr.length)]
}

// 生成订单库存数据
export function generateOrderInventory() {
  const clients = ['华润纺织', '盛虹集团', '恒力集团', '鲁泰纺织', '申洲国际']
  const products = [
    '纯棉印花布',
    '涤纶混纺布',
    '丝绸缎面',
    '牛仔布',
    '功能性面料',
    '针织布料',
    '色织提花布',
    '阻燃布料'
  ]
  const warehouses = ['主仓库', '东区仓库', '西区仓库', '成品仓库']

  const data = []
  for (let i = 1; i <= 15; i++) {
    const quantity = randomQuantity(500, 3000)
    const inboundUnitPrice = randomPrice(50, 150)
    const currentUnitPrice = inboundUnitPrice * (0.9 + Math.random() * 0.25) // ±15%波动

    const inboundValue = quantity * inboundUnitPrice
    const currentMarketValue = quantity * currentUnitPrice
    const difference = currentMarketValue - inboundValue

    data.push({
      id: i,
      orderNo: `ORD${String(i).padStart(6, '0')}`,
      clientName: randomChoice(clients),
      productName: randomChoice(products),
      quantity,
      unit: '米',
      inboundValue,
      currentMarketValue,
      difference,
      warehouseName: randomChoice(warehouses),
      location: `${randomChoice(['A', 'B', 'C', 'D'])}-${String(Math.floor(Math.random() * 20) + 1).padStart(2, '0')}-${String(Math.floor(Math.random() * 50) + 1).padStart(2, '0')}`,
      inboundDate: randomDate(60)
    })
  }

  return data
}

// 生成公共库存数据
export function generatePublicInventory() {
  const itemTypes = ['原料', '半成品', '成品', '辅料', '包装材料']
  const items = [
    { code: 'MAT001', name: '纯棉纱线', type: '原料' },
    { code: 'MAT002', name: '涤纶纤维', type: '原料' },
    { code: 'SFG001', name: '预处理胚布', type: '半成品' },
    { code: 'SFG002', name: '染色半成品', type: '半成品' },
    { code: 'FG001', name: '成品布料A', type: '成品' },
    { code: 'FG002', name: '成品布料B', type: '成品' },
    { code: 'AUX001', name: '染料包', type: '辅料' },
    { code: 'AUX002', name: '固色剂', type: '辅料' },
    { code: 'PKG001', name: '包装纸箱', type: '包装材料' },
    { code: 'PKG002', name: '塑料膜', type: '包装材料' }
  ]
  const warehouses = ['主仓库', '原料仓库', '成品仓库', '辅料仓库']

  const data = []
  for (let i = 0; i < items.length; i++) {
    const item = items[i]
    const quantity = randomQuantity(200, 2000)
    const inboundUnitPrice = randomPrice(30, 120)
    const currentUnitPrice = inboundUnitPrice * (0.85 + Math.random() * 0.35) // ±20%波动

    const inboundValue = quantity * inboundUnitPrice
    const currentMarketValue = quantity * currentUnitPrice
    const difference = currentMarketValue - inboundValue

    data.push({
      id: i + 1,
      itemCode: item.code,
      itemName: item.name,
      itemType: item.type,
      quantity,
      unit: randomChoice(['千克', '米', '件', '箱']),
      inboundValue,
      currentMarketValue,
      difference,
      warehouseName: randomChoice(warehouses),
      location: `${randomChoice(['A', 'B', 'C', 'D'])}-${String(Math.floor(Math.random() * 20) + 1).padStart(2, '0')}-${String(Math.floor(Math.random() * 50) + 1).padStart(2, '0')}`,
      inboundDate: randomDate(90)
    })
  }

  return data
}

// 生成原料库存数据
export function generateMaterialInventory() {
  const materials = [
    { code: 'RM001', name: '纯棉纱线40支', spec: '40支/2股' },
    { code: 'RM002', name: '涤纶POY', spec: '150D/48F' },
    { code: 'RM003', name: '粘胶纤维', spec: '1.5D×38mm' },
    { code: 'RM004', name: '锦纶FDY', spec: '70D/24F' },
    { code: 'RM005', name: '活性染料红', spec: '100%纯度' },
    { code: 'RM006', name: '活性染料蓝', spec: '100%纯度' },
    { code: 'RM007', name: '分散染料黄', spec: '高温型' },
    { code: 'RM008', name: '柔软剂', spec: '工业级' },
    { code: 'RM009', name: '增白剂', spec: 'VBL型' },
    { code: 'RM010', name: '防水整理剂', spec: '有机硅型' }
  ]
  const suppliers = ['江苏纺织原料', '浙江化工', '上海染料厂', '广东纤维', '山东化学']
  const warehouses = ['原料仓库A区', '原料仓库B区', '化学品仓库']

  const data = []
  for (let i = 0; i < materials.length; i++) {
    const material = materials[i]
    const quantity = randomQuantity(500, 5000)
    const inboundUnitPrice = randomPrice(20, 100)
    const currentUnitPrice = inboundUnitPrice * (0.9 + Math.random() * 0.25)

    const inboundValue = quantity * inboundUnitPrice
    const currentMarketValue = quantity * currentUnitPrice
    const difference = currentMarketValue - inboundValue

    data.push({
      id: i + 1,
      materialCode: material.code,
      materialName: material.name,
      specification: material.spec,
      quantity,
      unit: randomChoice(['千克', '吨', '升']),
      inboundValue,
      currentMarketValue,
      difference,
      supplierName: randomChoice(suppliers),
      warehouseName: randomChoice(warehouses),
      inboundDate: randomDate(45)
    })
  }

  return data
}

// 生成半成品库存数据
export function generateSemifinishedInventory() {
  const items = [
    { code: 'SFG001', name: '预处理纯棉胚布', spec: '40支/幅宽150cm' },
    { code: 'SFG002', name: '染色涤棉布', spec: '65/35混纺' },
    { code: 'SFG003', name: '印花底布', spec: '纯棉/幅宽160cm' },
    { code: 'SFG004', name: '整理后胚布', spec: '涤纶/幅宽180cm' },
    { code: 'SFG005', name: '烫金半成品', spec: '涤纶缎面' },
    { code: 'SFG006', name: '刺绣底布', spec: '纯棉帆布' },
    { code: 'SFG007', name: '复合半成品', spec: '涤纶+海绵' }
  ]
  const warehouses = ['半成品仓库A', '半成品仓库B', '加工车间仓库']

  const data = []
  for (let i = 0; i < items.length; i++) {
    const item = items[i]
    const quantity = randomQuantity(300, 2000)
    const inboundUnitPrice = randomPrice(40, 130)
    const currentUnitPrice = inboundUnitPrice * (0.92 + Math.random() * 0.20)

    const inboundValue = quantity * inboundUnitPrice
    const currentMarketValue = quantity * currentUnitPrice
    const difference = currentMarketValue - inboundValue

    data.push({
      id: i + 1,
      itemCode: item.code,
      itemName: item.name,
      specification: item.spec,
      quantity,
      unit: '米',
      inboundValue,
      currentMarketValue,
      difference,
      warehouseName: randomChoice(warehouses),
      location: `${randomChoice(['A', 'B', 'C'])}-${String(Math.floor(Math.random() * 15) + 1).padStart(2, '0')}-${String(Math.floor(Math.random() * 40) + 1).padStart(2, '0')}`,
      inboundDate: randomDate(30)
    })
  }

  return data
}

// 生成成品库存数据
export function generateFinishedInventory() {
  const products = [
    { code: 'FG001', name: '纯棉印花布-花卉系列', spec: '幅宽150cm/40支' },
    { code: 'FG002', name: '涤棉混纺工装布', spec: '幅宽160cm/CVC 60/40' },
    { code: 'FG003', name: '高档丝绸缎面', spec: '幅宽140cm/真丝100%' },
    { code: 'FG004', name: '牛仔布-深蓝', spec: '幅宽150cm/12盎司' },
    { code: 'FG005', name: '防水透气功能布', spec: '幅宽150cm/涤纶' },
    { code: 'FG006', name: '针织棉布-条纹', spec: '幅宽180cm/纯棉' },
    { code: 'FG007', name: '色织提花窗帘布', spec: '幅宽280cm/涤纶' },
    { code: 'FG008', name: '阻燃工业用布', spec: '幅宽200cm/芳纶' },
    { code: 'FG009', name: '婚纱蕾丝面料', spec: '幅宽130cm/涤纶网纱' },
    { code: 'FG010', name: '运动速干面料', spec: '幅宽150cm/锦纶' }
  ]
  const warehouses = ['成品仓库1号库', '成品仓库2号库', '成品仓库3号库', '待发货区']

  const data = []
  for (let i = 0; i < products.length; i++) {
    const product = products[i]
    const quantity = randomQuantity(200, 1500)
    const inboundUnitPrice = randomPrice(60, 180)
    const currentUnitPrice = inboundUnitPrice * (0.88 + Math.random() * 0.30)

    const inboundValue = quantity * inboundUnitPrice
    const currentMarketValue = quantity * currentUnitPrice
    const difference = currentMarketValue - inboundValue

    data.push({
      id: i + 1,
      productCode: product.code,
      productName: product.name,
      specification: product.spec,
      quantity,
      unit: '米',
      inboundValue,
      currentMarketValue,
      difference,
      warehouseName: randomChoice(warehouses),
      location: `${randomChoice(['A', 'B', 'C', 'D'])}-${String(Math.floor(Math.random() * 25) + 1).padStart(2, '0')}-${String(Math.floor(Math.random() * 60) + 1).padStart(2, '0')}`,
      inboundDate: randomDate(20)
    })
  }

  return data
}
