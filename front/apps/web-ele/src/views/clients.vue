<template>
  <div class="client-management">
    <DataTable
      :config="pageConfig"
      :loading="searchLoading"
      @view="handleView"
      @edit="handleEdit"
      @bulkAction="handleBulkAction"
    />

    <!-- 详情/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="800px"
      @close="handleDialogClose"
    >
      <el-form :model="currentRow" label-width="100px" style="max-height: 70vh; overflow-y: auto;">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="ID">
              <el-input v-model="currentRow.id" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="客户代码">
              <el-input v-model="currentRow.customNo" :disabled="dialogMode === 'view'" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="内部编码">
              <el-input v-model="currentRow.customerCode" :disabled="dialogMode === 'view'" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="业务员">
              <el-input v-model="currentRow.sales" :disabled="dialogMode === 'view'" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="客户名称">
          <el-input v-model="currentRow.customName" :disabled="dialogMode === 'view'" />
        </el-form-item>

        <el-form-item label="英文名称">
          <el-input v-model="currentRow.customNameEn" :disabled="dialogMode === 'view'" />
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="联系人">
              <el-input v-model="currentRow.contactor" :disabled="dialogMode === 'view'" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="电话">
              <el-input v-model="currentRow.unitPhone" :disabled="dialogMode === 'view'" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="手机">
              <el-input v-model="currentRow.mobile" :disabled="dialogMode === 'view'" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="传真">
              <el-input v-model="currentRow.faxNum" :disabled="dialogMode === 'view'" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="邮箱">
          <el-input v-model="currentRow.email" :disabled="dialogMode === 'view'" />
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="国家">
              <el-input v-model="currentRow.stateChNm" :disabled="dialogMode === 'view'" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="国家代码">
              <el-input v-model="currentRow.country" :disabled="dialogMode === 'view'" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="中文地址">
          <el-input
            v-model="currentRow.address"
            type="textarea"
            :rows="2"
            :disabled="dialogMode === 'view'"
          />
        </el-form-item>

        <el-form-item label="英文地址">
          <el-input
            v-model="currentRow.addressEn"
            type="textarea"
            :rows="2"
            :disabled="dialogMode === 'view'"
          />
        </el-form-item>

        <el-form-item label="所属客户">
          <el-input v-model="currentRow.pyCustomName" :disabled="dialogMode === 'view'" />
        </el-form-item>

        <el-form-item label="检验要求">
          <el-input
            v-model="currentRow.checkRequest"
            type="textarea"
            :rows="3"
            :disabled="dialogMode === 'view'"
          />
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="状态">
              <el-input v-model="currentRow.customStatus" :disabled="dialogMode === 'view'" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="输入人">
              <el-input v-model="currentRow.docMan" :disabled="dialogMode === 'view'" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="添加时间" v-if="currentRow.inputDate">
              <el-input
                :value="new Date(currentRow.inputDate).toLocaleDateString('zh-CN')"
                disabled
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="创建时间" v-if="currentRow.createdAt">
              <el-input
                :value="new Date(currentRow.createdAt).toLocaleString('zh-CN')"
                disabled
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="更新时间" v-if="currentRow.updatedAt">
          <el-input
            :value="new Date(currentRow.updatedAt).toLocaleString('zh-CN')"
            disabled
          />
        </el-form-item>
      </el-form>

      <template #footer v-if="dialogMode === 'edit'">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave" :loading="saving">
          保存
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { ElMessage } from 'element-plus';
import DataTable from '#/components/Table/index.vue';
import { elasticsearchService } from '#/api/core/es';
import type {
  PageConfig,
  ColumnConfig,
  FilterConfig,
  BulkAction,
} from '#/components/Table/types';

import { useDataTable } from '#/composables/useDataTable';
const { searchLoading } = useDataTable('client', 20);

// 页面配置
const pageConfig: PageConfig = {
  pageType: 'client',
  title: '客户管理',
  index: 'clients',
  pageSize: 20,
  columns: [
    { key: 'id', label: 'ID', width: 80, sortable: true, visible: true, order: 0 },
    { key: 'customNo', label: '客户代码', width: 120, sortable: true, visible: true, order: 1 },
    { key: 'customerCode', label: '内部编码', width: 120, sortable: true, visible: true, order: 2 },
    { key: 'customName', label: '客户名称', width: 200, sortable: true, visible: true, order: 3 },
    { key: 'customNameEn', label: '英文名称', width: 200, visible: false, order: 4 },
    { key: 'sales', label: '业务员', width: 120, visible: true, order: 5 },
    { key: 'contactor', label: '联系人', width: 120, visible: true, order: 6 },
    { key: 'unitPhone', label: '电话', width: 130, visible: true, order: 7 },
    { key: 'mobile', label: '手机', width: 130, visible: true, order: 8 },
    { key: 'email', label: '邮箱', width: 180, visible: false, order: 9 },
    { key: 'stateChNm', label: '国家', width: 100, visible: true, order: 10 },
    { key: 'country', label: '国家代码', width: 100, visible: false, order: 11 },
    { key: 'address', label: '中文地址', width: 250, visible: true, order: 12 },
    { key: 'addressEn', label: '英文地址', width: 250, visible: false, order: 13 },
    { key: 'faxNum', label: '传真', width: 130, visible: false, order: 14 },
    { key: 'pyCustomName', label: '所属客户', width: 150, visible: false, order: 15 },
    { key: 'checkRequest', label: '检验要求', width: 200, visible: false, order: 16 },
    { key: 'customStatus', label: '状态', width: 100, visible: true, order: 17 },
    { key: 'docMan', label: '输入人', width: 100, visible: false, order: 18 },
    {
      key: 'inputDate',
      label: '添加时间',
      width: 120,
      sortable: true,
      visible: true,
      order: 19,
      formatter: (v: string) => (v ? new Date(v).toLocaleDateString('zh-CN') : '-'),
    },
    {
      key: 'createdAt',
      label: '创建时间',
      width: 180,
      sortable: true,
      visible: false,
      order: 20,
      formatter: (v: string) => (v ? new Date(v).toLocaleString('zh-CN') : '-'),
    },
    {
      key: 'updatedAt',
      label: '更新时间',
      width: 180,
      sortable: true,
      visible: false,
      order: 21,
      formatter: (v: string) => (v ? new Date(v).toLocaleString('zh-CN') : '-'),
    },
  ] as ColumnConfig[],
  filters: [
    {
      key: 'level',
      label: '客户等级',
      type: 'select',
      placeholder: '请选择客户等级',
      options: [], // 后端会提供接口返回选项
    },
    {
      key: 'type',
      label: '客户类型',
      type: 'select',
      placeholder: '请选择客户类型',
      options: [], // 后端会提供接口返回选项
    },
    {
      key: 'status',
      label: '客户状态',
      type: 'select',
      placeholder: '请选择客户状态',
      options: [], // 后端会提供接口返回选项
    },
  ] as FilterConfig[],
  bulkActions: [
    {
      key: 'delete',
      label: '批量删除',
      type: 'danger',
      confirm: true,
      confirmMessage: '确定要删除选中的客户吗？此操作不可恢复！',
    },
    { key: 'export', label: '导出数据', type: 'primary' },
  ] as BulkAction[],
};

// ----- 对话框状态 -----
const dialogVisible = ref(false);
const dialogMode = ref<'view' | 'edit'>('view');
const currentRow = ref<any>({});
const saving = ref(false);

const dialogTitle = computed(() =>
  dialogMode.value === 'view' ? '查看客户' : '编辑客户'
);

// 查看
const handleView = (row: any) => {
  dialogMode.value = 'view';
  currentRow.value = { ...row };
  dialogVisible.value = true;
};

// 编辑
const handleEdit = (row: any) => {
  dialogMode.value = 'edit';
  currentRow.value = { ...row };
  dialogVisible.value = true;
};

// 保存
const handleSave = async () => {
  if (!currentRow.value._id) {
    ElMessage.error('缺少必要的ID信息');
    return;
  }
  saving.value = true;

  try {
    const { _id, id, createdAt, updatedAt, ...updateData } = currentRow.value;
    await elasticsearchService.update(_id, 'clients', updateData);

    ElMessage.success('保存成功');
    dialogVisible.value = false;
    window.location.reload();
  } catch (error) {
    console.error(error);
  } finally {
    saving.value = false;
  }
};

// 批量操作
const handleBulkAction = async ({ action, rows }: { action: string; rows: any[] }) => {
  if (!rows.length) return ElMessage.warning('请选择数据');

  const ids = rows.map((r) => r._id).filter(Boolean);
  if (!ids.length) return ElMessage.error('缺少 ID');

  switch (action) {
    case 'delete':
      await handleBulkDelete(ids);
      break;
    case 'export':
      await handleExport(rows);
      break;
  }
};

// 批量删除
const handleBulkDelete = async (ids: string[]) => {
  const result = await elasticsearchService.bulkDelete(ids, 'clients');

  if (result.success) {
    ElMessage.success(`已删除 ${result.successCount} 条`);
    if (result.failedCount > 0) {
      ElMessage.warning(`${result.failedCount} 条删除失败`);
      console.error(result.errors);
    }
    setTimeout(() => window.location.reload(), 500);
  }
};

// 导出
const handleExport = async (rows: any[]) => {
  const data = rows.map(({ _id, ...rest }) => rest);
  const headers = Object.keys(data[0] || {});
  const csvContent = [
    headers.join(','),
    ...data.map((row) =>
      headers.map((h) => JSON.stringify(row[h] || '')).join(',')
    ),
  ].join('\n');

  const blob = new Blob(['\ufeff' + csvContent], {
    type: 'text/csv;charset=utf-8;',
  });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');

  link.href = url;
  link.download = `clients_${Date.now()}.csv`;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(url);

  ElMessage.success('导出成功');
};

// 对话框关闭
const handleDialogClose = () => {
  currentRow.value = {};
  saving.value = false;
};
</script>

<style scoped lang="scss">
.client-management {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f7fa;
  overflow: hidden;

  :deep(.data-table-container) {
    flex: 1;
    margin: 20px;
    overflow: hidden;
  }
}
</style>
