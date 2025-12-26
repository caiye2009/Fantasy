<template>
  <div class="client-management">
    <DataTable
      :config="pageConfig"
      :loading="searchLoading"
      @view="handleView"
      @bulkAction="handleBulkAction"
      @topAction="handleTopAction"
    />

    <!-- 客户详情抽屉 -->
    <el-drawer
      v-model="drawerVisible"
      title="客户详情"
      size="60%"
      @close="handleDrawerClose"
    >
      <div v-if="currentRow" class="client-detail">
        <!-- 基本信息 -->
        <el-card class="detail-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>基本信息</span>
              <el-button v-if="!isEditing" type="primary" size="small" @click="startEdit">
                修改
              </el-button>
              <div v-else>
                <el-button size="small" @click="cancelEdit">取消</el-button>
                <el-button type="primary" size="small" :loading="saving" @click="saveEdit">
                  保存
                </el-button>
              </div>
            </div>
          </template>

          <el-form v-if="isEditing" :model="editForm" label-width="100px">
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="客户代码">
                  <el-input v-model="editForm.customNo" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="内部编码">
                  <el-input v-model="editForm.customerCode" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="客户名称">
                  <el-input v-model="editForm.customName" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="英文名称">
                  <el-input v-model="editForm.customNameEn" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="业务员">
                  <el-input v-model="editForm.sales" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="联系人">
                  <el-input v-model="editForm.contactor" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="电话">
                  <el-input v-model="editForm.unitPhone" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="手机">
                  <el-input v-model="editForm.mobile" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="传真">
                  <el-input v-model="editForm.faxNum" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="邮箱">
                  <el-input v-model="editForm.email" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="国家">
                  <el-input v-model="editForm.stateChNm" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="国家代码">
                  <el-input v-model="editForm.country" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-form-item label="中文地址">
              <el-input v-model="editForm.address" type="textarea" :rows="2" />
            </el-form-item>
            <el-form-item label="英文地址">
              <el-input v-model="editForm.addressEn" type="textarea" :rows="2" />
            </el-form-item>
            <el-form-item label="所属客户">
              <el-input v-model="editForm.pyCustomName" />
            </el-form-item>
            <el-form-item label="检验要求">
              <el-input v-model="editForm.checkRequest" type="textarea" :rows="3" />
            </el-form-item>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="状态">
                  <el-input v-model="editForm.customStatus" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="输入人">
                  <el-input v-model="editForm.docMan" />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>

          <el-descriptions v-else :column="2" border>
            <el-descriptions-item label="ID">{{ currentRow.id || '-' }}</el-descriptions-item>
            <el-descriptions-item label="客户代码">{{ currentRow.customNo || '-' }}</el-descriptions-item>
            <el-descriptions-item label="内部编码">{{ currentRow.customerCode || '-' }}</el-descriptions-item>
            <el-descriptions-item label="客户名称">{{ currentRow.customName || '-' }}</el-descriptions-item>
            <el-descriptions-item label="英文名称">{{ currentRow.customNameEn || '-' }}</el-descriptions-item>
            <el-descriptions-item label="业务员">{{ currentRow.sales || '-' }}</el-descriptions-item>
            <el-descriptions-item label="联系人">{{ currentRow.contactor || '-' }}</el-descriptions-item>
            <el-descriptions-item label="电话">{{ currentRow.unitPhone || '-' }}</el-descriptions-item>
            <el-descriptions-item label="手机">{{ currentRow.mobile || '-' }}</el-descriptions-item>
            <el-descriptions-item label="传真">{{ currentRow.faxNum || '-' }}</el-descriptions-item>
            <el-descriptions-item label="邮箱">{{ currentRow.email || '-' }}</el-descriptions-item>
            <el-descriptions-item label="国家">{{ currentRow.stateChNm || '-' }}</el-descriptions-item>
            <el-descriptions-item label="国家代码">{{ currentRow.country || '-' }}</el-descriptions-item>
            <el-descriptions-item label="中文地址" :span="2">{{ currentRow.address || '-' }}</el-descriptions-item>
            <el-descriptions-item label="英文地址" :span="2">{{ currentRow.addressEn || '-' }}</el-descriptions-item>
            <el-descriptions-item label="所属客户" :span="2">{{ currentRow.pyCustomName || '-' }}</el-descriptions-item>
            <el-descriptions-item label="检验要求" :span="2">{{ currentRow.checkRequest || '-' }}</el-descriptions-item>
            <el-descriptions-item label="状态">{{ currentRow.customStatus || '-' }}</el-descriptions-item>
            <el-descriptions-item label="输入人">{{ currentRow.docMan || '-' }}</el-descriptions-item>
            <el-descriptions-item label="添加时间">{{ currentRow.inputDate ? new Date(currentRow.inputDate).toLocaleDateString('zh-CN') : '-' }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ currentRow.createdAt ? new Date(currentRow.createdAt).toLocaleString('zh-CN') : '-' }}</el-descriptions-item>
            <el-descriptions-item label="更新时间" :span="2">{{ currentRow.updatedAt ? new Date(currentRow.updatedAt).toLocaleString('zh-CN') : '-' }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </div>
    </el-drawer>

    <!-- 新增客户对话框 -->
    <el-dialog
      v-model="createDialogVisible"
      title="新增客户"
      width="800px"
      @close="handleCreateClose"
    >
      <el-form :model="createForm" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="客户代码">
              <el-input v-model="createForm.customNo" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="内部编码">
              <el-input v-model="createForm.customerCode" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="客户名称" required>
              <el-input v-model="createForm.customName" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="英文名称">
              <el-input v-model="createForm.customNameEn" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="业务员">
              <el-input v-model="createForm.sales" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="联系人">
              <el-input v-model="createForm.contactor" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="电话">
              <el-input v-model="createForm.unitPhone" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="手机">
              <el-input v-model="createForm.mobile" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="传真">
              <el-input v-model="createForm.faxNum" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="邮箱">
              <el-input v-model="createForm.email" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="国家">
              <el-input v-model="createForm.stateChNm" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="国家代码">
              <el-input v-model="createForm.country" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="中文地址">
          <el-input v-model="createForm.address" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="英文地址">
          <el-input v-model="createForm.addressEn" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="所属客户">
          <el-input v-model="createForm.pyCustomName" />
        </el-form-item>
        <el-form-item label="检验要求">
          <el-input v-model="createForm.checkRequest" type="textarea" :rows="3" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="状态">
              <el-input v-model="createForm.customStatus" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="输入人">
              <el-input v-model="createForm.docMan" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="creating" @click="handleCreateSubmit">
          确定
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
  index: 'client',
  pageSize: 20,
  columns: [
    { key: 'id', label: 'ID', width: 80, sortable: true, visible: true, order: 0, showOverflowTooltip: true },
    { key: 'customNo', label: '客户代码', width: 120, sortable: true, visible: true, order: 1, showOverflowTooltip: true },
    { key: 'customerCode', label: '内部编码', width: 120, sortable: true, visible: true, order: 2, showOverflowTooltip: true },
    { key: 'customName', label: '客户名称', width: 200, sortable: true, visible: true, order: 3, showOverflowTooltip: true },
    { key: 'customNameEn', label: '英文名称', width: 200, visible: false, order: 4, showOverflowTooltip: true },
    { key: 'sales', label: '业务员', width: 120, visible: true, order: 5, showOverflowTooltip: true },
    { key: 'contactor', label: '联系人', width: 120, visible: true, order: 6, showOverflowTooltip: true },
    { key: 'unitPhone', label: '电话', width: 130, visible: true, order: 7, showOverflowTooltip: true },
    { key: 'mobile', label: '手机', width: 130, visible: true, order: 8, showOverflowTooltip: true },
    { key: 'email', label: '邮箱', width: 180, visible: false, order: 9, showOverflowTooltip: true },
    { key: 'stateChNm', label: '国家', width: 100, visible: true, order: 10, showOverflowTooltip: true },
    { key: 'country', label: '国家代码', width: 100, visible: false, order: 11, showOverflowTooltip: true },
    { key: 'address', label: '中文地址', width: 250, visible: true, order: 12, showOverflowTooltip: true },
    { key: 'addressEn', label: '英文地址', width: 250, visible: false, order: 13, showOverflowTooltip: true },
    { key: 'faxNum', label: '传真', width: 130, visible: false, order: 14, showOverflowTooltip: true },
    { key: 'pyCustomName', label: '所属客户', width: 150, visible: false, order: 15, showOverflowTooltip: true },
    { key: 'checkRequest', label: '检验要求', width: 200, visible: false, order: 16, showOverflowTooltip: true },
    { key: 'customStatus', label: '状态', width: 100, visible: true, order: 17, showOverflowTooltip: true },
    { key: 'docMan', label: '输入人', width: 100, visible: false, order: 18, showOverflowTooltip: true },
    {
      key: 'inputDate',
      label: '添加时间',
      width: 120,
      sortable: true,
      visible: true,
      order: 19,
      showOverflowTooltip: true,
      formatter: (v: string) => (v ? new Date(v).toLocaleDateString('zh-CN') : '-'),
    },
    {
      key: 'createdAt',
      label: '创建时间',
      width: 180,
      sortable: true,
      visible: false,
      order: 20,
      showOverflowTooltip: true,
      formatter: (v: string) => (v ? new Date(v).toLocaleString('zh-CN') : '-'),
    },
    {
      key: 'updatedAt',
      label: '更新时间',
      width: 180,
      sortable: true,
      visible: false,
      order: 21,
      showOverflowTooltip: true,
      formatter: (v: string) => (v ? new Date(v).toLocaleString('zh-CN') : '-'),
    },
  ] as ColumnConfig[],
  filters: [
    {
      key: 'sales',
      label: '业务员',
      type: 'select',
      placeholder: '请选择业务员',
      fetchOptions: async () => {
        try {
          const response = await elasticsearchService.search({
            index: 'client',
            pagination: { offset: 0, size: 0 },
            aggRequests: {
              sales: {
                type: 'terms',
                field: 'sales',
                size: 50,
              },
            },
          })
          const buckets = response.aggregations?.sales?.buckets || []
          return buckets.map((bucket: any) => ({
            label: String(bucket.key),
            value: bucket.key,
          }))
        } catch (error) {
          console.error('加载业务员选项失败:', error)
          return []
        }
      },
    },
    {
      key: 'country',
      label: '国家',
      type: 'select',
      placeholder: '请选择国家',
      fetchOptions: async () => {
        try {
          const response = await elasticsearchService.search({
            index: 'client',
            pagination: { offset: 0, size: 0 },
            aggRequests: {
              country: {
                type: 'terms',
                field: 'country',
                size: 100,
              },
            },
          })
          const buckets = response.aggregations?.country?.buckets || []
          return buckets.map((bucket: any) => ({
            label: String(bucket.key),
            value: bucket.key,
          }))
        } catch (error) {
          console.error('加载国家选项失败:', error)
          return []
        }
      },
    },
    {
      key: 'customStatus',
      label: '客户状态',
      type: 'select',
      placeholder: '请选择客户状态',
      fetchOptions: async () => {
        try {
          const response = await elasticsearchService.search({
            index: 'client',
            pagination: { offset: 0, size: 0 },
            aggRequests: {
              customStatus: {
                type: 'terms',
                field: 'customStatus',
                size: 10,
              },
            },
          })
          const buckets = response.aggregations?.customStatus?.buckets || []
          return buckets.map((bucket: any) => ({
            label: String(bucket.key),
            value: bucket.key,
          }))
        } catch (error) {
          console.error('加载客户状态选项失败:', error)
          return []
        }
      },
    },
  ] as FilterConfig[],
  topActions: [
    { key: 'create', label: '新增客户', type: 'primary' },
  ],
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

// ----- 抽屉状态 -----
const drawerVisible = ref(false);
const currentRow = ref<any>(null);
const saving = ref(false);
const isEditing = ref(false);
const editForm = ref<any>({});

// 查看
const handleView = (row: any) => {
  currentRow.value = { ...row };
  isEditing.value = false;
  drawerVisible.value = true;
};

// 开始编辑
const startEdit = () => {
  if (currentRow.value) {
    editForm.value = { ...currentRow.value };
    isEditing.value = true;
  }
};

// 取消编辑
const cancelEdit = () => {
  isEditing.value = false;
  editForm.value = {};
};

// 保存编辑
const saveEdit = async () => {
  if (!currentRow.value._id) {
    ElMessage.error('缺少必要的ID信息');
    return;
  }

  saving.value = true;
  try {
    const { _id, id, createdAt, updatedAt, ...updateData } = editForm.value;
    await elasticsearchService.update(_id, 'client', updateData);

    ElMessage.success('保存成功');
    isEditing.value = false;
    currentRow.value = { ...currentRow.value, ...updateData };
    window.location.reload();
  } catch (error) {
    console.error(error);
    ElMessage.error('保存失败');
  } finally {
    saving.value = false;
  }
};

// 顶部操作
const handleTopAction = ({ action }: { action: string }) => {
  if (action === 'create') {
    handleCreate();
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

// 新增客户
const createDialogVisible = ref(false);
const createForm = ref<any>({
  customNo: '',
  customerCode: '',
  customName: '',
  customNameEn: '',
  sales: '',
  contactor: '',
  unitPhone: '',
  mobile: '',
  faxNum: '',
  email: '',
  stateChNm: '',
  country: '',
  address: '',
  addressEn: '',
  pyCustomName: '',
  checkRequest: '',
  customStatus: '',
  docMan: '',
});
const creating = ref(false);

const handleCreate = () => {
  createDialogVisible.value = true;
};

const handleCreateSubmit = async () => {
  if (!createForm.value.customName) {
    ElMessage.error('请输入客户名称');
    return;
  }

  creating.value = true;
  try {
    await elasticsearchService.create('client', createForm.value);
    ElMessage.success('新增成功');
    createDialogVisible.value = false;
    setTimeout(() => window.location.reload(), 500);
  } catch (error) {
    console.error(error);
    ElMessage.error('新增失败');
  } finally {
    creating.value = false;
  }
};

const handleCreateClose = () => {
  createForm.value = {
    customNo: '',
    customerCode: '',
    customName: '',
    customNameEn: '',
    sales: '',
    contactor: '',
    unitPhone: '',
    mobile: '',
    faxNum: '',
    email: '',
    stateChNm: '',
    country: '',
    address: '',
    addressEn: '',
    pyCustomName: '',
    checkRequest: '',
    customStatus: '',
    docMan: '',
  };
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

// 抽屉关闭
const handleDrawerClose = () => {
  currentRow.value = null;
  isEditing.value = false;
  editForm.value = {};
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
