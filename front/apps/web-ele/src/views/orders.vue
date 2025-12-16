<template>
  <div>
    <!-- 触发按钮 -->
    <el-button type="primary" @click="openModal">创建订单</el-button>

    <!-- Modal 对话框 - 大尺寸 -->
    <el-dialog
      v-model="dialogVisible"
      title="创建订单"
      width="90%"
      :close-on-click-modal="false"
      @close="handleClose"
      style="max-width: 1400px"
    >
      <div class="order-creation-wrapper">
        <!-- 步骤选择区域 -->
        <div class="steps-container" :class="{ 'slide-left': showPreview }">
          <div class="order-creation-flow">
            <!-- 1. 选择客户 -->
            <div class="step-section">
              <div 
                class="step-header" 
                :class="{ 'active': currentStep === 1, 'completed': selectedCustomer }"
                @click="goToStep(1)"
              >
                <div v-if="currentStep === 1 || !selectedCustomer" class="step-title">
                  <span class="step-number">1</span>
                  <span class="step-name">选择客户</span>
                </div>
                <div v-else class="step-summary-inline">
                  <span class="step-number-small">1</span>
                  <span class="summary-value">{{ selectedCustomer.name }} - {{ selectedCustomer.phone }}</span>
                </div>
              </div>
              
              <div v-if="currentStep === 1" class="step-content">
                <el-form :model="form" label-width="100px">
                  <el-form-item label="客户名称">
                    <el-select 
                      v-model="form.customerId" 
                      placeholder="请选择客户"
                      @change="handleCustomerChange"
                      style="width: 100%"
                      size="large"
                    >
                      <el-option
                        v-for="customer in customers"
                        :key="customer.id"
                        :label="customer.name"
                        :value="customer.id"
                      >
                        <div style="display: flex; justify-content: space-between">
                          <span>{{ customer.name }}</span>
                          <span style="color: #8492a6; font-size: 13px">{{ customer.phone }}</span>
                        </div>
                      </el-option>
                    </el-select>
                  </el-form-item>
                </el-form>
              </div>
            </div>

            <!-- 2. 选择产品 -->
            <div class="step-section">
              <div 
                class="step-header" 
                :class="{ 
                  'active': currentStep === 2, 
                  'completed': selectedProduct,
                  'disabled': !selectedCustomer 
                }"
                @click="goToStep(2)"
              >
                <div v-if="currentStep === 2 || !selectedProduct" class="step-title">
                  <span class="step-number">2</span>
                  <span class="step-name">选择产品</span>
                </div>
                <div v-else class="step-summary-inline">
                  <span class="step-number-small">2</span>
                  <span class="summary-value">{{ selectedProduct.name }} × {{ form.quantity }} (¥{{ selectedProduct.price }}/件)</span>
                </div>
              </div>
              
              <div v-if="currentStep === 2" class="step-content">
                <el-form :model="form" label-width="100px">
                  <el-row :gutter="20">
                    <el-col :span="12">
                      <el-form-item label="产品名称">
                        <el-select 
                          v-model="form.productId" 
                          placeholder="请选择产品"
                          @change="handleProductChange"
                          style="width: 100%"
                          size="large"
                        >
                          <el-option
                            v-for="product in products"
                            :key="product.id"
                            :label="product.name"
                            :value="product.id"
                          >
                            <div style="display: flex; justify-content: space-between">
                              <span>{{ product.name }}</span>
                              <span style="color: #8492a6; font-size: 13px">¥{{ product.price }}</span>
                            </div>
                          </el-option>
                        </el-select>
                      </el-form-item>
                    </el-col>
                    
                    <el-col :span="12">
                      <el-form-item label="购买数量">
                        <el-input-number 
                          v-model="form.quantity" 
                          :min="1" 
                          :max="100"
                          style="width: 100%"
                          size="large"
                          @change="handleQuantityChange"
                        />
                      </el-form-item>
                    </el-col>
                  </el-row>
                </el-form>
              </div>
            </div>

            <!-- 3. 选择库存 -->
            <div class="step-section">
              <div 
                class="step-header" 
                :class="{ 
                  'active': currentStep === 3, 
                  'completed': selectedWarehouse,
                  'disabled': !selectedProduct 
                }"
                @click="goToStep(3)"
              >
                <div v-if="currentStep === 3 || !selectedWarehouse" class="step-title">
                  <span class="step-number">3</span>
                  <span class="step-name">选择库存</span>
                </div>
                <div v-else class="step-summary-inline">
                  <span class="step-number-small">3</span>
                  <span class="summary-value">{{ selectedWarehouse.name }} (库存: {{ selectedWarehouse.stock }})</span>
                </div>
              </div>
              
              <div v-if="currentStep === 3" class="step-content">
                <el-form :model="form" label-width="100px">
                  <el-form-item label="发货仓库">
                    <el-select 
                      v-model="form.warehouseId" 
                      placeholder="请选择仓库"
                      @change="handleWarehouseChange"
                      style="width: 100%"
                      size="large"
                    >
                      <el-option
                        v-for="warehouse in warehouses"
                        :key="warehouse.id"
                        :label="warehouse.name"
                        :value="warehouse.id"
                      >
                        <div style="display: flex; justify-content: space-between">
                          <span>{{ warehouse.name }}</span>
                          <span style="color: #8492a6; font-size: 13px">库存: {{ warehouse.stock }}</span>
                        </div>
                      </el-option>
                    </el-select>
                  </el-form-item>
                </el-form>
              </div>
            </div>
          </div>
        </div>

        <!-- 预览区域 -->
        <div class="preview-container" :class="{ 'show-preview': showPreview }">
          <div class="preview-content-full">
            <h3 style="margin-bottom: 24px; color: #303133; font-size: 20px;">订单信息预览</h3>
            
            <el-row :gutter="40">
              <el-col :span="8">
                <div class="info-card">
                  <div class="info-card-title">
                    <el-icon style="margin-right: 8px"><User /></el-icon>
                    客户信息
                  </div>
                  <div class="info-card-content">
                    <div class="info-item">
                      <span class="info-label">客户姓名：</span>
                      <span class="info-value">{{ selectedCustomer?.name }}</span>
                    </div>
                    <div class="info-item">
                      <span class="info-label">联系电话：</span>
                      <span class="info-value">{{ selectedCustomer?.phone }}</span>
                    </div>
                    <div class="info-item">
                      <span class="info-label">收货地址：</span>
                      <span class="info-value">{{ selectedCustomer?.address }}</span>
                    </div>
                  </div>
                </div>
              </el-col>
              
              <el-col :span="8">
                <div class="info-card">
                  <div class="info-card-title">
                    <el-icon style="margin-right: 8px"><ShoppingCart /></el-icon>
                    产品信息
                  </div>
                  <div class="info-card-content">
                    <div class="info-item">
                      <span class="info-label">产品名称：</span>
                      <span class="info-value">{{ selectedProduct?.name }}</span>
                    </div>
                    <div class="info-item">
                      <span class="info-label">产品编号：</span>
                      <span class="info-value">{{ selectedProduct?.code }}</span>
                    </div>
                    <div class="info-item">
                      <span class="info-label">产品单价：</span>
                      <span class="info-value">¥{{ selectedProduct?.price }}</span>
                    </div>
                    <div class="info-item">
                      <span class="info-label">购买数量：</span>
                      <span class="info-value">{{ form.quantity }}</span>
                    </div>
                  </div>
                </div>
              </el-col>
              
              <el-col :span="8">
                <div class="info-card">
                  <div class="info-card-title">
                    <el-icon style="margin-right: 8px"><Box /></el-icon>
                    库存信息
                  </div>
                  <div class="info-card-content">
                    <div class="info-item">
                      <span class="info-label">发货仓库：</span>
                      <span class="info-value">{{ selectedWarehouse?.name }}</span>
                    </div>
                    <div class="info-item">
                      <span class="info-label">仓库地址：</span>
                      <span class="info-value">{{ selectedWarehouse?.address }}</span>
                    </div>
                    <div class="info-item">
                      <span class="info-label">当前库存：</span>
                      <span class="info-value">{{ selectedWarehouse?.stock }}</span>
                    </div>
                  </div>
                </div>
              </el-col>
            </el-row>
            
            <div class="total-amount-card">
              <div class="total-label">订单总金额</div>
              <div class="total-value">¥{{ totalAmount }}</div>
            </div>
            
            <div class="preview-actions">
              <el-button size="large" @click="backToSteps">返回</el-button>
              <el-button type="primary" size="large" @click="submitOrder" :loading="submitting">
                确认下单
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { ElMessage } from 'element-plus';
import { User, ShoppingCart, Box } from '@element-plus/icons-vue';

// 对话框显示状态
const dialogVisible = ref(false);
const currentStep = ref(1);
const submitting = ref(false);
const showPreview = ref(false);

// 表单数据
const form = ref({
  customerId: '',
  productId: '',
  quantity: 1,
  warehouseId: ''
});

// 模拟客户数据
const customers = ref([
  { id: 1, name: '张三', phone: '13800138001', address: '北京市朝阳区xxx街道xxx号' },
  { id: 2, name: '李四', phone: '13800138002', address: '上海市浦东新区xxx路xxx号' },
  { id: 3, name: '王五', phone: '13800138003', address: '广州市天河区xxx大道xxx号' },
  { id: 4, name: '赵六', phone: '13800138004', address: '深圳市南山区xxx街xxx号' }
]);

// 模拟产品数据
const products = ref([
  { id: 1, name: 'iPhone 15 Pro', code: 'IP15P-001', price: 7999, category: '手机' },
  { id: 2, name: 'MacBook Pro 14', code: 'MBP14-001', price: 14999, category: '电脑' },
  { id: 3, name: 'AirPods Pro 2', code: 'APP2-001', price: 1899, category: '耳机' },
  { id: 4, name: 'iPad Air', code: 'IPA-001', price: 4799, category: '平板' }
]);

// 模拟仓库数据
const warehouses = ref([
  { id: 1, name: '北京仓库', address: '北京市顺义区物流园区A区', stock: 150 },
  { id: 2, name: '上海仓库', address: '上海市青浦区工业园区B区', stock: 200 },
  { id: 3, name: '广州仓库', address: '广州市白云区物流中心C区', stock: 180 },
  { id: 4, name: '深圳仓库', address: '深圳市龙岗区科技园D区', stock: 120 }
]);

// 已选择的数据
const selectedCustomer = computed(() => 
  customers.value.find(c => c.id === form.value.customerId)
);

const selectedProduct = computed(() => 
  products.value.find(p => p.id === form.value.productId)
);

const selectedWarehouse = computed(() => 
  warehouses.value.find(w => w.id === form.value.warehouseId)
);

// 计算总金额
const totalAmount = computed(() => {
  if (selectedProduct.value && form.value.quantity) {
    return (selectedProduct.value.price * form.value.quantity).toFixed(2);
  }
  return 0;
});

// 打开 Modal
const openModal = () => {
  dialogVisible.value = true;
  currentStep.value = 1;
};

// 关闭 Modal
const handleClose = () => {
  // 重置表单
  form.value = {
    customerId: '',
    productId: '',
    quantity: 1,
    warehouseId: ''
  };
  currentStep.value = 1;
};

// 处理客户选择 - 自动进入下一步
const handleCustomerChange = () => {
  if (form.value.customerId) {
    currentStep.value = 2;
  }
};

// 处理产品选择
const handleProductChange = () => {
  // 当产品和数量都有值时自动进入下一步
  if (form.value.productId && form.value.quantity) {
    currentStep.value = 3;
  }
};

// 处理数量变化
const handleQuantityChange = () => {
  // 当产品和数量都有值时自动进入下一步
  if (form.value.productId && form.value.quantity) {
    currentStep.value = 3;
  }
};

// 处理仓库选择 - 显示预览
const handleWarehouseChange = () => {
  if (form.value.warehouseId) {
    showPreview.value = true;
  }
};

// 跳转到指定步骤
const goToStep = (step) => {
  // 只能跳转到已完成的步骤
  if (step === 1) {
    currentStep.value = step;
  } else if (step === 2 && selectedCustomer.value) {
    currentStep.value = step;
  } else if (step === 3 && selectedProduct.value) {
    currentStep.value = step;
  }
};

// 返回步骤选择
const backToSteps = () => {
  showPreview.value = false;
};

// 关闭 Modal
// const handleClose = () => {
//   // 重置表单
//   form.value = {
//     customerId: '',
//     productId: '',
//     quantity: 1,
//     warehouseId: ''
//   };
//   currentStep.value = 1;
//   showPreview.value = false;
// };

// 提交订单
const submitOrder = async () => {
  submitting.value = true;
  
  // 模拟提交延迟
  setTimeout(() => {
    const orderData = {
      customer: selectedCustomer.value,
      product: selectedProduct.value,
      quantity: form.value.quantity,
      warehouse: selectedWarehouse.value,
      totalAmount: totalAmount.value,
      orderTime: new Date().toLocaleString()
    };
    
    console.log('订单数据:', orderData);
    
    ElMessage.success('订单创建成功！');
    submitting.value = false;
    dialogVisible.value = false;
    handleClose();
  }, 1000);
};
</script>

<style scoped>
.order-creation-wrapper {
  position: relative;
  overflow: hidden;
  min-height: 60vh;
}

.steps-container {
  transition: transform 0.4s ease-in-out;
  padding: 10px 0;
  min-height: 60vh;
}

.steps-container.slide-left {
  transform: translateX(-100%);
}

.preview-container {
  position: absolute;
  top: 0;
  left: 100%;
  width: 100%;
  height: 100%;
  transition: transform 0.4s ease-in-out;
  padding: 20px 0;
}

.preview-container.show-preview {
  transform: translateX(-100%);
}

.preview-content-full {
  padding: 20px;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.preview-actions {
  margin-top: 32px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 20px;
  border-top: 1px solid #e4e7ed;
}

.order-creation-flow {
  padding: 10px 0;
  min-height: 60vh;
}

.step-section {
  margin-bottom: 16px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  overflow: hidden;
  transition: all 0.3s;
}

.step-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 18px 24px;
  background-color: #f5f7fa;
  cursor: pointer;
  transition: all 0.3s;
}

.step-header:hover:not(.disabled) {
  background-color: #ecf5ff;
}

.step-header.active {
  background-color: #409eff;
  color: white;
  padding: 20px 24px;
}

.step-header.completed {
  background-color: #f0f9ff;
  border-left: 4px solid #67c23a;
}

.step-header.disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.step-title {
  display: flex;
  align-items: center;
  gap: 12px;
}

.step-number {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background-color: #409eff;
  color: white;
  font-weight: bold;
  font-size: 16px;
}

.step-header.active .step-number {
  background-color: white;
  color: #409eff;
}

.step-header.completed .step-number {
  background-color: #67c23a;
}

.step-name {
  font-size: 18px;
  font-weight: 500;
}

.step-summary-inline {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
}

.step-number-small {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background-color: #67c23a;
  color: white;
  font-weight: bold;
  font-size: 14px;
  flex-shrink: 0;
}

.summary-value {
  color: #303133;
  font-weight: 500;
  font-size: 16px;
}

.step-content {
  padding: 32px 24px;
  background-color: white;
}

.preview-content {
  padding: 40px 32px;
}

.preview-section {
  max-width: 100%;
}

.info-card {
  background: linear-gradient(135deg, #f5f7fa 0%, #ffffff 100%);
  border: 1px solid #e4e7ed;
  border-radius: 12px;
  padding: 24px;
  height: 100%;
  transition: all 0.3s;
}

.info-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  transform: translateY(-2px);
}

.info-card-title {
  display: flex;
  align-items: center;
  font-size: 16px;
  font-weight: 600;
  color: #409eff;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 2px solid #409eff;
}

.info-card-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.info-item {
  display: flex;
  line-height: 1.8;
}

.info-label {
  color: #909399;
  min-width: 80px;
  font-size: 14px;
}

.info-value {
  color: #303133;
  font-weight: 500;
  flex: 1;
  font-size: 14px;
}

.total-amount-card {
  margin-top: 40px;
  padding: 32px;
  background: linear-gradient(135deg, #fff3e0 0%, #ffe0b2 100%);
  border: 2px solid #ff9800;
  border-radius: 12px;
  text-align: center;
}

.total-label {
  font-size: 16px;
  color: #666;
  margin-bottom: 12px;
}

.total-value {
  font-size: 36px;
  font-weight: bold;
  color: #f56c6c;
}

:deep(.el-dialog__body) {
  padding: 20px 30px;
}

:deep(.el-form-item__label) {
  font-weight: 500;
}

:deep(.el-select) {
  width: 100%;
}

:deep(.el-input-number) {
  width: 100%;
}

:deep(.el-input-number .el-input__inner) {
  text-align: left;
}
</style>