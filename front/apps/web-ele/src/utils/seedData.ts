/**
 * 种子数据脚本 - 批量插入原料和工艺数据
 * 使用方法：在浏览器控制台调用 window.seedData()
 */

import { requestClient } from '#/api/request'

// 原料数据（10条）
const materials = [
  { code: 'MAT-001', name: '40支纯棉胚布', spec: '40支/2股', unit: '米', category: '胚布', currentPrice: 12.5, description: '优质纯棉原料，适合印花' },
  { code: 'MAT-002', name: '60支涤棉混纺胚布', spec: '60支/涤棉65/35', unit: '米', category: '胚布', currentPrice: 18.8, description: '混纺面料，耐磨性强' },
  { code: 'MAT-003', name: '活性染料-红', spec: '100%纯度', unit: '千克', category: '染料', currentPrice: 85.0, description: '高品质活性染料，着色力强' },
  { code: 'MAT-004', name: '活性染料-蓝', spec: '100%纯度', unit: '千克', category: '染料', currentPrice: 88.0, description: '环保染料，色牢度高' },
  { code: 'MAT-005', name: '柔软剂', spec: '工业级', unit: '千克', category: '助剂', currentPrice: 32.0, description: '改善面料手感' },
  { code: 'MAT-006', name: '固色剂', spec: '高效型', unit: '千克', category: '助剂', currentPrice: 28.5, description: '提高色牢度' },
  { code: 'MAT-007', name: '增白剂', spec: 'VBL型', unit: '千克', category: '助剂', currentPrice: 45.0, description: '增强面料白度' },
  { code: 'MAT-008', name: '防水整理剂', spec: '有机硅型', unit: '千克', category: '助剂', currentPrice: 65.0, description: '防水透气功能整理' },
  { code: 'MAT-009', name: '涤纶POY', spec: '150D/48F', unit: '千克', category: '纤维', currentPrice: 22.0, description: '涤纶预取向丝' },
  { code: 'MAT-010', name: '粘胶纤维', spec: '1.5D×38mm', unit: '千克', category: '纤维', currentPrice: 28.0, description: '粘胶短纤维，吸湿性好' }
]

// 工艺数据（10条）
const processes = [
  { name: '染色', description: '面料染色工艺，根据客户要求调配颜色', currentPrice: 15.0 },
  { name: '印花', description: '滚筒印花或数码印花工艺', currentPrice: 20.0 },
  { name: '定型', description: '高温定型，固定面料尺寸和形态', currentPrice: 8.0 },
  { name: '预缩', description: '预缩处理，减少成品缩水率', currentPrice: 6.5 },
  { name: '丝光', description: '棉布丝光处理，增加光泽度', currentPrice: 12.0 },
  { name: '柔软整理', description: '使用柔软剂改善手感', currentPrice: 5.0 },
  { name: '防水整理', description: '防水透气功能整理', currentPrice: 10.0 },
  { name: '涂层', description: 'PU或PVC涂层加工', currentPrice: 18.0 },
  { name: '压光', description: '面料压光处理，增加光泽', currentPrice: 7.0 },
  { name: '烫金', description: '热转印烫金工艺', currentPrice: 25.0 }
]

/**
 * 创建原料
 */
async function createMaterial(data: { code: string; name: string; spec: string; unit: string; category: string; currentPrice: number; description: string }) {
  try {
    const response = await requestClient.post('/material', data)
    console.log(`✅ 原料创建成功: ${data.name}`, response.data)
    return response.data
  } catch (error: any) {
    console.error(`❌ 原料创建失败: ${data.name}`, error.response?.data || error.message)
    throw error
  }
}

/**
 * 创建工艺
 */
async function createProcess(data: { name: string; description: string; currentPrice: number }) {
  try {
    const response = await requestClient.post('/process', data)
    console.log(`✅ 工艺创建成功: ${data.name}`, response.data)
    return response.data
  } catch (error: any) {
    console.error(`❌ 工艺创建失败: ${data.name}`, error.response?.data || error.message)
    throw error
  }
}

/**
 * 批量插入所有数据
 */
export async function seedData() {
  console.log('🚀 开始插入种子数据...')

  let materialCount = 0
  let processCount = 0

  // 1. 插入原料数据
  console.log('\n📦 正在插入原料数据...')
  for (const material of materials) {
    try {
      await createMaterial(material)
      materialCount++
      // 避免请求过快，延迟200ms
      await new Promise(resolve => setTimeout(resolve, 200))
    } catch (error) {
      // 继续插入下一条
    }
  }

  // 2. 插入工艺数据
  console.log('\n🔧 正在插入工艺数据...')
  for (const process of processes) {
    try {
      await createProcess(process)
      processCount++
      // 避免请求过快，延迟200ms
      await new Promise(resolve => setTimeout(resolve, 200))
    } catch (error) {
      // 继续插入下一条
    }
  }

  console.log(`\n✨ 数据插入完成！`)
  console.log(`   - 原料: ${materialCount}/${materials.length} 条`)
  console.log(`   - 工艺: ${processCount}/${processes.length} 条`)

  return {
    materials: materialCount,
    processes: processCount
  }
}

// 导出到全局供控制台调用
if (typeof window !== 'undefined') {
  (window as any).seedData = seedData
}
