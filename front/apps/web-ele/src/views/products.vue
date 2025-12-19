<template>
  <div class="product-management">
    <DataTable
      :config="pageConfig"
      :loading="searchLoading"
      @view="handleView"
      @edit="handleEdit"
      @bulkAction="handleBulkAction"
    />

    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="900px"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="120px"
      >
        <el-form-item label="产品名称" prop="name">
          <el-input
            v-model="formData.name"
            :disabled="dialogMode === 'view'"
          />
        </el-form-item>

        <!-- 原料配置 -->
        <el-form-item label="原料配置" prop="materials">
          <el-button
            v-if="dialogMode !== 'view'"
            size="small"
            type="primary"
            @click="addMaterial"
            style="margin-bottom: 8px"
          >
            添加原料
          </el-button>

          <el-table :data="formData.materials" border>
            <el-table-column label="原料">
              <template #default="{ row }">
                <el-select
                  v-model="row.material_id"
                  :disabled="dialogMode === 'view'"
                  filterable
                  placeholder="选择原料"
                  @visible-change="onMaterialSelectVisible"
                >
                  <el-option
                    v-for="m in materialsList"
                    :key="m.id"
                    :label="m.name"
                    :value="m.id"
                  />
                </el-select>
              </template>
            </el-table-column>

            <el-table-column label="占比 %">
              <template #default="{ row }">
                <el-input-number
                  v-model="row.ratioPercent"
                  :disabled="dialogMode === 'view'"
                  :min="0"
                  :max="100"
                />
              </template>
            </el-table-column>

            <el-table-column v-if="dialogMode !== 'view'" label="操作">
              <template #default="{ $index }">
                <el-button
                  type="danger"
                  size="small"
                  @click="removeMaterial($index)"
                >
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <div
            class="ratio-summary"
            :class="{ error: totalRatio !== 100 && formData.materials.length }"
          >
            总占比：{{ totalRatio }}%
          </div>
        </el-form-item>

        <!-- 工艺配置 -->
        <el-form-item label="工艺配置" prop="processes">
          <el-select
            v-model="selectedProcessIds"
            multiple
            filterable
            placeholder="选择工艺"
            :disabled="dialogMode === 'view'"
            @visible-change="onProcessSelectVisible"
            @change="onProcessChange"
          >
            <el-option
              v-for="p in processesList"
              :key="p.id"
              :label="p.name"
              :value="p.id"
            />
          </el-select>
        </el-form-item>
      </el-form>

      <template #footer v-if="dialogMode !== 'view'">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          保存
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import DataTable from '#/components/Table/index.vue'
import { elasticsearchService } from '#/api/core/es'
import { createProductApi, updateProductApi } from '#/api/core/product'
import { useDataTable } from '#/composables/useDataTable'

/* ================= 基础 ================= */

const { searchLoading } = useDataTable('product', 20)

const pageConfig = {
  pageType: 'product',
  title: '产品管理',
  index: 'product',
  rowKey: 'id',
  actions: ['view', 'edit', 'create'],
}

/* ================= Dialog ================= */

const dialogVisible = ref(false)
const dialogMode = ref<'create' | 'edit' | 'view'>('view')
const dialogTitle = computed(() =>
  dialogMode.value === 'create'
    ? '新建产品'
    : dialogMode.value === 'edit'
    ? '编辑产品'
    : '查看产品'
)

const submitting = ref(false)
const formRef = ref<FormInstance>()

/* ================= 数据 ================= */

const formData = ref({
  name: '',
  materials: [] as { material_id: number; ratioPercent: number }[],
  processes: [] as { process_id: number }[],
})

const selectedProcessIds = ref<number[]>([])

/* ================= 校验 ================= */

const totalRatio = computed(() =>
  formData.value.materials.reduce((s, m) => s + (m.ratioPercent || 0), 0)
)

const formRules: FormRules = {
  name: [{ required: true, message: '请输入产品名称' }],
  materials: [
    {
      validator: (_, __, cb) => {
        if (!formData.value.materials.length) {
          cb(new Error('至少一个原料'))
        } else if (totalRatio.value !== 100) {
          cb(new Error('原料占比必须等于100%'))
        } else cb()
      },
    },
  ],
}

/* ================= 懒加载（关键） ================= */

const materialsList = ref<any[]>([])
const processesList = ref<any[]>([])

const materialsLoaded = ref(false)
const processesLoaded = ref(false)

const loadMaterials = async () => {
  if (materialsLoaded.value) return
  const res = await elasticsearchService.search({
    index: 'material',
    pagination: { offset: 0, size: 1000 },
  })
  materialsList.value = res.items || []
  materialsLoaded.value = true
}

const loadProcesses = async () => {
  if (processesLoaded.value) return
  const res = await elasticsearchService.search({
    index: 'process',
    pagination: { offset: 0, size: 1000 },
  })
  processesList.value = res.items || []
  processesLoaded.value = true
}

const onMaterialSelectVisible = (visible: boolean) => {
  if (visible) loadMaterials()
}

const onProcessSelectVisible = (visible: boolean) => {
  if (visible) loadProcesses()
}

/* ================= 行为 ================= */

const addMaterial = () => {
  formData.value.materials.push({ material_id: 0, ratioPercent: 0 })
}

const removeMaterial = (i: number) => {
  formData.value.materials.splice(i, 1)
}

const onProcessChange = (ids: number[]) => {
  formData.value.processes = ids.map(id => ({ process_id: id }))
}

const handleCreate = async () => {
  dialogMode.value = 'create'
  dialogVisible.value = true
  await Promise.all([loadMaterials(), loadProcesses()])
}

const handleEdit = async (row: any) => {
  dialogMode.value = 'edit'
  formData.value = {
    name: row.name,
    materials: row.materials.map((m: any) => ({
      material_id: m.material_id,
      ratioPercent: m.ratio * 100,
    })),
    processes: row.processes,
  }
  selectedProcessIds.value = row.processes.map((p: any) => p.process_id)
  dialogVisible.value = true
  await Promise.all([loadMaterials(), loadProcesses()])
}

const handleView = async (row: any) => {
  dialogMode.value = 'view'
  await handleEdit(row)
}

/* ================= 提交 ================= */

const handleSubmit = async () => {
  await formRef.value?.validate()
  submitting.value = true

  const payload = {
    name: formData.value.name,
    materials: formData.value.materials.map(m => ({
      material_id: m.material_id,
      ratio: m.ratioPercent / 100,
    })),
    processes: formData.value.processes,
  }

  try {
    dialogMode.value === 'create'
      ? await createProductApi(payload)
      : await updateProductApi(payload)

    ElMessage.success('保存成功')
    dialogVisible.value = false
    location.reload()
  } finally {
    submitting.value = false
  }
}

/* ================= 批量 ================= */

const handleBulkAction = ({ action }: any) => {
  if (action === 'create') handleCreate()
}

/* ================= reset ================= */

const handleDialogClose = () => {
  formRef.value?.clearValidate()
  formData.value = { name: '', materials: [], processes: [] }
  selectedProcessIds.value = []
}
</script>

<style scoped>
.product-management {
  height: 100%;
}
.ratio-summary.error {
  color: #f56c6c;
}
</style>
