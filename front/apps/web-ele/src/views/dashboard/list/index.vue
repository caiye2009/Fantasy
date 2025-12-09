<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { requestClient } from '#/api/request';
import MonthRangePicker from '@/components/MonthRangePicker/index.vue';

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
  dateRange: null as string[] | null, // 使用月份范围选择器返回的格式
});

// 获取客户列表
async function fetchCustomerList() {
  try {
    const response = await requestClient.get('/return-analysis/customers');
    customerList.value = response;
  } catch (error) {
    console.error('获取客户列表失败:', error);
    alert('获取客户列表失败');
  }
}

// 获取退货分析数据
async function fetchAnalysisData() {
  loading.value = true;
  try {
    const response = await requestClient.post('/return-analysis/analysis', {
      customerNo: searchForm.value.customerNo,
      dateRange: {
        start: searchForm.value.dateRange?.[0] || '',
        end: searchForm.value.dateRange?.[1] || '',
      },
    });
    analysisData.value = response;
  } catch (error: any) {
    console.error('获取分析数据失败:', error);
    alert(error?.message || '获取分析数据失败');
  } finally {
    loading.value = false;
  }
}

// 重置搜索
function handleReset() {
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
  fetchCustomerList();
});
</script>

<template>
  <div class="p-5">
    <!-- 搜索区域 -->
    <div class="mb-4 rounded-lg bg-white p-6 shadow">
      <div class="flex flex-wrap items-end gap-4">
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
        
        <div class="flex flex-col">
          <label class="mb-2 text-sm font-medium text-gray-700">月份范围</label>
          <MonthRangePicker 
            v-model="searchForm.dateRange"
            start-placeholder="开始月份"
            end-placeholder="结束月份"
          />
        </div>
        
        <button
          class="rounded-md bg-blue-600 px-6 py-2 text-white transition hover:bg-blue-700 disabled:bg-gray-400"
          :disabled="loading"
          @click="fetchAnalysisData"
        >
          {{ loading ? '查询中...' : '查询' }}
        </button>
        <button
          class="rounded-md border border-gray-300 bg-white px-6 py-2 transition hover:bg-gray-50"
          @click="handleReset"
        >
          重置
        </button>
      </div>
      <div class="mt-2 text-xs text-gray-500">
        * 月份范围可选，不选表示查询全部时间；支持快捷选择：本月、本年、最近3/6/12个月
      </div>
    </div>

    <!-- 结果展示 -->
    <div v-if="analysisData" class="space-y-4">
      <!-- 查询条件回显 -->
      <div class="rounded-lg bg-blue-50 p-4">
        <div class="text-sm text-gray-700">
          <span class="font-semibold">查询条件：</span>
          <span v-if="analysisData.queryConditions.customerName">
            客户：{{ analysisData.queryConditions.customerName }}
          </span>
          <span v-else>客户：全部</span>
          <span class="ml-4">
            时间：
            <template v-if="analysisData.queryConditions.dateRange.start">
              {{ analysisData.queryConditions.dateRange.start }} 至 {{ analysisData.queryConditions.dateRange.end }}
            </template>
            <template v-else>全部时间</template>
          </span>
        </div>
      </div>

      <!-- 概览卡片 -->
      <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
        <div class="rounded-lg bg-white p-6 shadow">
          <div class="text-sm text-gray-600">总订单数</div>
          <div class="mt-2 text-3xl font-bold text-gray-900">
            {{ formatNumber(analysisData.totalOrders, 0) }} <span class="text-lg">单</span>
          </div>
        </div>
        <div class="rounded-lg bg-white p-6 shadow">
          <div class="text-sm text-gray-600">退货订单数</div>
          <div class="mt-2 text-3xl font-bold text-red-600">
            {{ formatNumber(analysisData.amountStats.returnedOrderCount, 0) }} <span class="text-lg">单</span>
          </div>
        </div>
      </div>

      <!-- 米数维度 -->
      <div class="rounded-lg bg-white p-6 shadow">
        <h3 class="mb-4 text-lg font-semibold text-gray-900">米数统计</h3>
        <div class="grid grid-cols-1 gap-4 md:grid-cols-4">
          <div class="rounded-lg bg-blue-50 p-4">
            <div class="text-sm text-gray-600">订单总米数</div>
            <div class="mt-2 text-2xl font-bold text-blue-600">
              {{ formatNumber(analysisData.meterStats.totalMeters) }} <span class="text-sm">米</span>
            </div>
          </div>
          <div class="rounded-lg bg-orange-50 p-4">
            <div class="text-sm text-gray-600">退货总米数</div>
            <div class="mt-2 text-2xl font-bold text-orange-600">
              {{ formatNumber(analysisData.meterStats.returnedMeters) }} <span class="text-sm">米</span>
            </div>
          </div>
          <div class="rounded-lg bg-red-50 p-4">
            <div class="text-sm text-gray-600">退货率</div>
            <div class="mt-2 text-2xl font-bold" :class="analysisData.meterStats.returnRate === 'N/A' ? 'text-gray-400' : 'text-red-600'">
              {{ analysisData.meterStats.returnRate }}
            </div>
          </div>
          <div class="rounded-lg bg-gray-50 p-4">
            <div class="text-sm text-gray-600">米数订单数</div>
            <div class="mt-2 text-2xl font-bold text-gray-700">
              {{ formatNumber(analysisData.meterStats.orderCount, 0) }} <span class="text-sm">单</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 重量维度 -->
      <div class="rounded-lg bg-white p-6 shadow">
        <h3 class="mb-4 text-lg font-semibold text-gray-900">重量统计</h3>
        <div class="grid grid-cols-1 gap-4 md:grid-cols-4">
          <div class="rounded-lg bg-green-50 p-4">
            <div class="text-sm text-gray-600">订单总重量</div>
            <div class="mt-2 text-2xl font-bold text-green-600">
              {{ formatNumber(analysisData.weightStats.totalWeight) }} <span class="text-sm">kg</span>
            </div>
          </div>
          <div class="rounded-lg bg-yellow-50 p-4">
            <div class="text-sm text-gray-600">退货总重量</div>
            <div class="mt-2 text-2xl font-bold text-yellow-600">
              {{ formatNumber(analysisData.weightStats.returnedWeight) }} <span class="text-sm">kg</span>
            </div>
          </div>
          <div class="rounded-lg bg-red-50 p-4">
            <div class="text-sm text-gray-600">退货率</div>
            <div class="mt-2 text-2xl font-bold" :class="analysisData.weightStats.returnRate === 'N/A' ? 'text-gray-400' : 'text-red-600'">
              {{ analysisData.weightStats.returnRate }}
            </div>
          </div>
          <div class="rounded-lg bg-gray-50 p-4">
            <div class="text-sm text-gray-600">重量订单数</div>
            <div class="mt-2 text-2xl font-bold text-gray-700">
              {{ formatNumber(analysisData.weightStats.orderCount, 0) }} <span class="text-sm">单</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 金额维度 -->
      <div class="rounded-lg bg-white p-6 shadow">
        <h3 class="mb-4 text-lg font-semibold text-gray-900">退款金额统计</h3>
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div class="rounded-lg bg-purple-50 p-4">
            <div class="text-sm text-gray-600">退款总额（RMB）</div>
            <div class="mt-2 text-2xl font-bold text-purple-600">
              {{ formatCurrency(analysisData.amountStats.totalAmountRMB) }}
            </div>
            <div class="mt-1 text-xs text-gray-500">USD 已按汇率 8 换算</div>
          </div>
          <div class="rounded-lg bg-gray-50 p-4">
            <div class="text-sm text-gray-600">退货订单数</div>
            <div class="mt-2 text-2xl font-bold text-gray-700">
              {{ formatNumber(analysisData.amountStats.returnedOrderCount, 0) }} <span class="text-sm">单</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="rounded-lg bg-white p-12 text-center shadow">
      <div class="text-gray-400">
        <svg class="mx-auto h-16 w-16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
        </svg>
        <div class="mt-4 text-lg">请选择查询条件并点击查询按钮</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 自定义样式 */
</style>