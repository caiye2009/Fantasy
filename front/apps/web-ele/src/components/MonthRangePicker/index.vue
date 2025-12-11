<script lang="ts" setup>
import { ref, watch } from 'vue';
import { ElDatePicker } from 'element-plus';
import type { Dayjs } from 'dayjs';

interface Props {
  modelValue?: string[] | null;
  startPlaceholder?: string;
  endPlaceholder?: string;
  size?: 'large' | 'default' | 'small';
  rangeSeparator?: string;
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: null,
  startPlaceholder: '开始月份',
  endPlaceholder: '结束月份',
  size: 'default',
  rangeSeparator: '至',
});

const emit = defineEmits<{
  'update:modelValue': [value: string[] | null];
}>();

// 内部值直接使用数组
const internalValue = ref<string[] | null>(null);

// 快捷选项
const shortcuts = [
  {
    text: '本月',
    value: () => {
      const now = new Date();
      return [now, now];
    },
  },
  {
    text: '本年',
    value: () => {
      const end = new Date();
      const start = new Date(new Date().getFullYear(), 0);
      return [start, end];
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

// 初始化
if (props.modelValue && props.modelValue.length === 2) {
  internalValue.value = [...props.modelValue];
  console.log('MonthRangePicker 初始化:', internalValue.value);
}

// 监听 props 变化
watch(
  () => props.modelValue,
  (newValue) => {
    console.log('MonthRangePicker props 变化:', newValue);
    if (newValue && newValue.length === 2) {
      internalValue.value = [...newValue];
    } else {
      internalValue.value = null;
    }
  },
);

// 监听内部值变化并向外发送
watch(internalValue, (newValue) => {
  console.log('MonthRangePicker 内部值变化:', newValue);
  if (newValue && newValue.length === 2) {
    emit('update:modelValue', [...newValue]);
  } else {
    emit('update:modelValue', null);
  }
});
</script>

<template>
  <ElDatePicker
    v-model="internalValue"
    type="monthrange"
    unlink-panels
    :range-separator="rangeSeparator"
    :start-placeholder="startPlaceholder"
    :end-placeholder="endPlaceholder"
    :shortcuts="shortcuts"
    :size="size"
    value-format="YYYY-MM"
    format="YYYY-MM"
    class="month-range-picker"
  />
</template>

<style scoped>
.month-range-picker {
  width: 280px;
}
</style>