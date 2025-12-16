<script lang="ts" setup>
import { ref, onMounted, watch } from 'vue';
import { requestClient } from '#/api/request';
import MonthRangePicker from '#/components/MonthRangePicker/index.vue';

// 数据类型定义
interface CustomerOption {
  customerNo: string;
  customerName: string;
}

interface MeterStats {
  totalMeters: number;
  returnedMeters: number;
  returnRate: string;
  orderCount: number;
}

interface WeightStats {
  totalWeight: number;
  returnedWeight: number;
  returnRate: string;
  orderCount: number;
}

interface AmountStats {
  totalAmountRMB: number;
  returnedOrderCount: number;
}

interface AnalysisResponse {
  queryConditions: {
    dateRange: {
      start: string;
      end: string;
    };
    customerNo: string;
    customerName: string;
  };
  meterStats: MeterStats;
  weightStats: WeightStats;
  amountStats: AmountStats;
  totalOrders: number;
}

// 响应式数据
const customerList = ref<CustomerOption[]>([]);
const analysisData = ref<AnalysisResponse | null>(null);
const loading = ref(false);

// 搜索表单
const searchForm = ref({
  customerNo: '',
  dateRange: null as string[] | null,
});

// 监听 dateRange 变化（调试用）
watch(
  () => searchForm.value.dateRange,
  (newValue) => {
    console.log('退货分析页面 dateRange 变化:', newValue);
  },
);

// 将月份格式转换为完整日期格式（只有日期，没有时间）
function convertMonthToDate(month: string, isEnd: boolean = false): string {
  // month 格式: "2025-01"
  const [year, monthNum] = month.split('-');
  
  if (isEnd) {
    // 结束日期：取该月最后一天
    const lastDay = new Date(Number(year), Number(monthNum), 0).getDate();
    return `${year}-${monthNum}-${String(lastDay).padStart(2, '0')}`;
  } else {
    // 开始日期：取该月第一天
    return `${year}-${monthNum}-01`;
  }
}

// 获取客户列表
async function fetchCustomerList() {
  try {
    const response = await requestClient.get('/return-analysis/customers');
    customerList.value = response.data;
    console.log('客户列表获取成功:', response.data);
  } catch (error) {
    console.error('获取客户列表失败:', error);
    alert('获取客户列表失败');
  }
}

// 获取退货分析数据
async function fetchAnalysisData() {
  console.log('开始查询，当前参数:', {
    customerNo: searchForm.value.customerNo,
    dateRange: searchForm.value.dateRange,
  });

  loading.value = true;
  try {
    // 构建请求数据
    const requestData: any = {
      customerNo: searchForm.value.customerNo || '',
    };

    // 只有在选择了日期范围时才添加 dateRange 字段
    if (searchForm.value.dateRange && searchForm.value.dateRange.length === 2) {
      const startDate = convertMonthToDate(searchForm.value.dateRange[0], false);
      const endDate = convertMonthToDate(searchForm.value.dateRange[1], true);
      
      requestData.dateRange = {
        start: startDate,
        end: endDate,
      };
    } else {
      // 没有选择日期范围，发送空字符串
      requestData.dateRange = {
        start: '',
        end: '',
      };
    }

    console.log('发送请求数据:', JSON.stringify(requestData, null, 2));

    const response = await requestClient.post('/return-analysis/analysis', requestData);
    analysisData.value = response.data;
    console.log('分析数据获取成功:', response.data);
  } catch (error: any) {
    console.error('获取分析数据失败:', error);
    alert(error?.msg || error?.message || '获取分析数据失败');
  } finally {
    loading.value = false;
  }
}

// 重置搜索
function handleReset() {
  console.log('重置搜索');
  searchForm.value = {
    customerNo: '',
    dateRange: null,
  };
  analysisData.value = null;
}

// 格式化数字（带千分位）
function formatNumber(num: number, decimals: number = 1): string {
  return num.toLocaleString('zh-CN', {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  });
}

// 格式化金额
function formatCurrency(num: number): string {
  return num.toLocaleString('zh-CN', {
    style: 'currency',
    currency: 'CNY',
  });
}

// 页面加载时获取客户列表
onMounted(() => {
  console.log('退货分析页面已挂载');
  fetchCustomerList();
});
</script>

<template>
  <div class="p-5">
    <!-- 搜索区域 -->
    <div class="mb-4 rounded-lg bg-white p-6 shadow">
      <div class="flex flex-wrap items-end gap-4">
        <!-- 客户选择 -->
        <div class="flex flex-col">
          <label class="mb-2 text-sm font-medium text-gray-700">客户选择</label>
          <select
            v-model="searchForm.customerNo"
            class="w-48 rounded-md border border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
          >
            <option value="">全部客户</option>
            <option 
              v-for="customer in customerList" 
              :key="customer.customerNo" 
              :value="customer.customerNo"
            >
              {{ customer.customerName }}
            </option>
          </select>
        </div>
        
        <!-- 月份范围选择器 -->
        <div class="flex flex-col">
          <label class="mb-2 text-sm font-medium text-gray-700">月份范围</label>
          <MonthRangePicker 
            v-model="searchForm.dateRange"
            start-placeholder="开始月份"
            end-placeholder="结束月份"
            range-separator="至"
          />
        </div>
        
        <!-- 操作按钮 -->
        <button
          class="rounded-md bg-blue-600 px-6 py-2 text-white transition hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed"
          :disabled="loading"
          @click="fetchAnalysisData"
        >
          {{ loading ? '查询中...' : '查询' }}
        </button>
        <button
          class="rounded-md border border-gray-300 bg-white px-6 py-2 text-gray-700 transition hover:bg-gray-50"
          @click="handleReset"
        >
          重置
        </button>
      </div>
      
      <!-- 提示信息 -->
      <div class="mt-3 text-xs text-gray-500">
        <svg class="mr-1 inline-block h-3.5 w-3.5" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
        </svg>
        月份范围可选，不选表示查询全部时间；支持快捷选择：本月、本年、最近3/6/12个月
      </div>
    </div>

    <!-- 结果展示区域 -->
    <div v-if="analysisData" class="space-y-4">
      <!-- 查询条件回显 -->
      <div class="rounded-lg bg-gradient-to-r from-blue-50 to-blue-100 p-4 shadow-sm">
        <div class="flex items-center text-sm text-gray-700">
          <svg class="mr-2 h-5 w-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span class="font-semibold">查询条件：</span>
          <span class="ml-2">
            客户：<span class="font-medium">{{ analysisData.queryConditions.customerName || '全部' }}</span>
          </span>
          <span class="ml-4">
            时间：
            <span class="font-medium">
              <template v-if="analysisData.queryConditions.dateRange.start">
                {{ analysisData.queryConditions.dateRange.start }} 至 {{ analysisData.queryConditions.dateRange.end }}
              </template>
              <template v-else>全部时间</template>
            </span>
          </span>
        </div>
      </div>

      <!-- 概览卡片 -->
      <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
        <div class="group rounded-lg bg-white p-6 shadow transition-shadow hover:shadow-lg">
          <div class="flex items-center justify-between">
            <div>
              <div class="text-sm text-gray-600">总订单数</div>
              <div class="mt-2 text-3xl font-bold text-gray-900">
                {{ formatNumber(analysisData.totalOrders, 0) }} <span class="text-lg text-gray-600">单</span>
              </div>
            </div>
            <div class="rounded-full bg-gray-100 p-3 transition-colors group-hover:bg-gray-200">
              <svg class="h-8 w-8 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
          </div>
        </div>
        <div class="group rounded-lg bg-white p-6 shadow transition-shadow hover:shadow-lg">
          <div class="flex items-center justify-between">
            <div>
              <div class="text-sm text-gray-600">退货订单数</div>
              <div class="mt-2 text-3xl font-bold text-red-600">
                {{ formatNumber(analysisData.amountStats.returnedOrderCount, 0) }} <span class="text-lg text-red-400">单</span>
              </div>
            </div>
            <div class="rounded-full bg-red-100 p-3 transition-colors group-hover:bg-red-200">
              <svg class="h-8 w-8 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
          </div>
        </div>
      </div>

      <!-- 米数维度 -->
      <div class="rounded-lg bg-white p-6 shadow">
        <div class="mb-4 flex items-center">
          <div class="mr-2 h-1 w-8 rounded bg-blue-600"></div>
          <h3 class="text-lg font-semibold text-gray-900">米数统计</h3>
        </div>
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <div class="rounded-lg bg-gradient-to-br from-blue-50 to-blue-100 p-4 transition-transform hover:scale-105">
            <div class="text-sm text-gray-600">订单总米数</div>
            <div class="mt-2 text-2xl font-bold text-blue-600">
              {{ formatNumber(analysisData.meterStats.totalMeters) }} <span class="text-sm text-blue-400">米</span>
            </div>
          </div>
          <div class="rounded-lg bg-gradient-to-br from-orange-50 to-orange-100 p-4 transition-transform hover:scale-105">
            <div class="text-sm text-gray-600">退货总米数</div>
            <div class="mt-2 text-2xl font-bold text-orange-600">
              {{ formatNumber(analysisData.meterStats.returnedMeters) }} <span class="text-sm text-orange-400">米</span>
            </div>
          </div>
          <div class="rounded-lg bg-gradient-to-br from-red-50 to-red-100 p-4 transition-transform hover:scale-105">
            <div class="text-sm text-gray-600">退货率</div>
            <div class="mt-2 text-2xl font-bold" :class="analysisData.meterStats.returnRate === 'N/A' ? 'text-gray-400' : 'text-red-600'">
              {{ analysisData.meterStats.returnRate }}
            </div>
          </div>
          <div class="rounded-lg bg-gradient-to-br from-gray-50 to-gray-100 p-4 transition-transform hover:scale-105">
            <div class="text-sm text-gray-600">米数订单数</div>
            <div class="mt-2 text-2xl font-bold text-gray-700">
              {{ formatNumber(analysisData.meterStats.orderCount, 0) }} <span class="text-sm text-gray-500">单</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 重量维度 -->
      <div class="rounded-lg bg-white p-6 shadow">
        <div class="mb-4 flex items-center">
          <div class="mr-2 h-1 w-8 rounded bg-green-600"></div>
          <h3 class="text-lg font-semibold text-gray-900">重量统计</h3>
        </div>
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <div class="rounded-lg bg-gradient-to-br from-green-50 to-green-100 p-4 transition-transform hover:scale-105">
            <div class="text-sm text-gray-600">订单总重量</div>
            <div class="mt-2 text-2xl font-bold text-green-600">
              {{ formatNumber(analysisData.weightStats.totalWeight) }} <span class="text-sm text-green-400">kg</span>
            </div>
          </div>
          <div class="rounded-lg bg-gradient-to-br from-yellow-50 to-yellow-100 p-4 transition-transform hover:scale-105">
            <div class="text-sm text-gray-600">退货总重量</div>
            <div class="mt-2 text-2xl font-bold text-yellow-600">
              {{ formatNumber(analysisData.weightStats.returnedWeight) }} <span class="text-sm text-yellow-400">kg</span>
            </div>
          </div>
          <div class="rounded-lg bg-gradient-to-br from-red-50 to-red-100 p-4 transition-transform hover:scale-105">
            <div class="text-sm text-gray-600">退货率</div>
            <div class="mt-2 text-2xl font-bold" :class="analysisData.weightStats.returnRate === 'N/A' ? 'text-gray-400' : 'text-red-600'">
              {{ analysisData.weightStats.returnRate }}
            </div>
          </div>
          <div class="rounded-lg bg-gradient-to-br from-gray-50 to-gray-100 p-4 transition-transform hover:scale-105">
            <div class="text-sm text-gray-600">重量订单数</div>
            <div class="mt-2 text-2xl font-bold text-gray-700">
              {{ formatNumber(analysisData.weightStats.orderCount, 0) }} <span class="text-sm text-gray-500">单</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 金额维度 -->
      <div class="rounded-lg bg-white p-6 shadow">
        <div class="mb-4 flex items-center">
          <div class="mr-2 h-1 w-8 rounded bg-purple-600"></div>
          <h3 class="text-lg font-semibold text-gray-900">退款金额统计</h3>
        </div>
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div class="rounded-lg bg-gradient-to-br from-purple-50 to-purple-100 p-6 transition-transform hover:scale-105">
            <div class="flex items-start justify-between">
              <div>
                <div class="text-sm text-gray-600">退款总额（RMB）</div>
                <div class="mt-2 text-3xl font-bold text-purple-600">
                  {{ formatCurrency(analysisData.amountStats.totalAmountRMB) }}
                </div>
                <div class="mt-2 text-xs text-gray-500">* USD 已按汇率 8 换算</div>
              </div>
              <div class="rounded-full bg-purple-200 p-3">
                <svg class="h-6 w-6 text-purple-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
            </div>
          </div>
          <div class="rounded-lg bg-gradient-to-br from-gray-50 to-gray-100 p-6 transition-transform hover:scale-105">
            <div class="flex items-start justify-between">
              <div>
                <div class="text-sm text-gray-600">退货订单数</div>
                <div class="mt-2 text-3xl font-bold text-gray-700">
                  {{ formatNumber(analysisData.amountStats.returnedOrderCount, 0) }} <span class="text-lg text-gray-500">单</span>
                </div>
              </div>
              <div class="rounded-full bg-gray-200 p-3">
                <svg class="h-6 w-6 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
                </svg>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="rounded-lg bg-white p-12 text-center shadow">
      <div class="text-gray-400">
        <svg class="mx-auto h-16 w-16 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
        </svg>
        <div class="mt-4 text-lg font-medium text-gray-500">暂无数据</div>
        <div class="mt-2 text-sm text-gray-400">请选择查询条件并点击查询按钮</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.space-y-4 > * {
  animation: fadeIn 0.3s ease-out;
}
</style>