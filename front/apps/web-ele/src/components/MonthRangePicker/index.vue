<!-- src/components/MonthRangePicker/index.vue -->
<template>
  <el-date-picker
    v-model="dateValue"
    type="monthrange"
    :unlink-panels="unlinkPanels"
    :range-separator="rangeSeparator"
    :start-placeholder="startPlaceholder"
    :end-placeholder="endPlaceholder"
    :shortcuts="computedShortcuts"
    :clearable="clearable"
    :disabled="disabled"
    :format="format"
    :value-format="valueFormat"
    @change="handleChange"
  />
</template>

<script lang="ts" setup>
import { ref, computed, watch } from 'vue';

interface Props {
  modelValue?: string[] | null;
  unlinkPanels?: boolean;
  rangeSeparator?: string;
  startPlaceholder?: string;
  endPlaceholder?: string;
  clearable?: boolean;
  disabled?: boolean;
  format?: string;
  valueFormat?: string;
  showShortcuts?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: null,
  unlinkPanels: true,
  rangeSeparator: '至',
  startPlaceholder: '开始月份',
  endPlaceholder: '结束月份',
  clearable: true,
  disabled: false,
  format: 'YYYY-MM',
  valueFormat: 'YYYY-MM-DD',
  showShortcuts: true,
});

const emit = defineEmits<{
  'update:modelValue': [value: string[] | null];
  'change': [value: string[] | null];
}>();

const dateValue = ref<Date[] | null>(null);

// 快捷选项
const shortcuts = [
  {
    text: '本月',
    value: () => {
      const now = new Date();
      return [new Date(now.getFullYear(), now.getMonth()), now];
    },
  },
  {
    text: '本年',
    value: () => {
      const now = new Date();
      const start = new Date(now.getFullYear(), 0);
      return [start, now];
    },
  },
  {
    text: '最近3个月',
    value: () => {
      const end = new Date();
      const start = new Date();
      start.setMonth(start.getMonth() - 2);
      return [start, end];
    },
  },
  {
    text: '最近6个月',
    value: () => {
      const end = new Date();
      const start = new Date();
      start.setMonth(start.getMonth() - 5);
      return [start, end];
    },
  },
  {
    text: '最近12个月',
    value: () => {
      const end = new Date();
      const start = new Date();
      start.setMonth(start.getMonth() - 11);
      return [start, end];
    },
  },
];

const computedShortcuts = computed(() => {
  return props.showShortcuts ? shortcuts : undefined;
});

// 监听外部值变化
watch(
  () => props.modelValue,
  (newVal) => {
    if (newVal && newVal.length === 2) {
      dateValue.value = [new Date(newVal[0]), new Date(newVal[1])];
    } else {
      dateValue.value = null;
    }
  },
  { immediate: true }
);

// 处理变化
function handleChange(value: Date[] | null) {
  let result: string[] | null = null;
  
  if (value && value.length === 2) {
    // 格式化为 YYYY-MM-DD（月份第一天和最后一天）
    const startYear = value[0].getFullYear();
    const startMonth = value[0].getMonth();
    const endYear = value[1].getFullYear();
    const endMonth = value[1].getMonth();
    
    const startDate = new Date(startYear, startMonth, 1);
    const endDate = new Date(endYear, endMonth + 1, 0); // 下个月的第0天 = 当月最后一天
    
    result = [
      formatDate(startDate),
      formatDate(endDate),
    ];
  }
  
  emit('update:modelValue', result);
  emit('change', result);
}

function formatDate(date: Date): string {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
}
</script>

<style scoped>
/* 可以添加自定义样式 */
</style>