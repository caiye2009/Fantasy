<template>
  <div class="user-management">
    <div class="header">
      <h2>用户管理</h2>
      <el-button type="primary" @click="handleCreate">新建用户</el-button>
    </div>

    <el-table
      v-loading="loading"
      :data="users"
      style="width: 100%"
      stripe
      border
    >
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="login_id" label="登录ID" width="120" />
      <el-table-column prop="username" label="用户名" width="150" />
      <el-table-column prop="department" label="部门" width="150" />
      <el-table-column prop="role" label="角色" width="150">
        <template #default="{ row }">
          {{ getRoleLabel(row.role) }}
        </template>
      </el-table-column>
      <el-table-column prop="email" label="邮箱" width="200" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
            {{ row.status === 'active' ? '激活' : '停用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" size="small" @click="handleView(row)">
            查看
          </el-button>
          <el-button link type="primary" size="small" @click="handleEdit(row)">
            编辑
          </el-button>
          <el-button
            link
            type="danger"
            size="small"
            @click="handleDelete(row)"
            :disabled="row.role === 'admin'"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :total="total"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="handleSizeChange"
      @current-change="handlePageChange"
      style="margin-top: 20px; justify-content: center"
    />

    <!-- 新建用户对话框 -->
    <el-dialog
      v-model="createDialogVisible"
      title="新建用户"
      width="500px"
      @close="handleCreateDialogClose"
    >
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="用户名" required>
          <el-input v-model="createForm.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="部门" required>
          <el-select
            v-model="createForm.department"
            placeholder="请选择部门"
            style="width: 100%"
          >
            <el-option
              v-for="dept in departments"
              :key="dept"
              :label="dept"
              :value="dept"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="角色" required>
          <el-select
            v-model="createForm.role"
            placeholder="请选择角色"
            style="width: 100%"
          >
            <el-option
              v-for="role in roles"
              :key="role"
              :label="getRoleLabel(role)"
              :value="role"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="createForm.email" placeholder="请输入邮箱（可选）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreateSubmit" :loading="creating">
          创建
        </el-button>
      </template>
    </el-dialog>

    <!-- 创建成功对话框 -->
    <el-dialog
      v-model="successDialogVisible"
      title="用户创建成功"
      width="500px"
    >
      <el-alert
        title="请妥善保管以下信息"
        type="success"
        :closable="false"
        style="margin-bottom: 20px"
      />
      <el-descriptions :column="1" border>
        <el-descriptions-item label="登录ID">
          <span style="font-size: 18px; font-weight: bold; color: #409eff">
            {{ createdUser.login_id }}
          </span>
        </el-descriptions-item>
        <el-descriptions-item label="初始密码">
          <span style="font-size: 18px; font-weight: bold; color: #f56c6c">
            {{ createdUser.password }}
          </span>
        </el-descriptions-item>
      </el-descriptions>
      <el-alert
        title="用户首次登录时需要修改密码"
        type="warning"
        :closable="false"
        style="margin-top: 20px"
      />
      <template #footer>
        <el-button type="primary" @click="handleSuccessDialogClose">
          我知道了
        </el-button>
      </template>
    </el-dialog>

    <!-- 查看/编辑对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      :title="detailDialogMode === 'view' ? '查看用户' : '编辑用户'"
      width="600px"
      @close="handleDetailDialogClose"
    >
      <el-form :model="currentUser" label-width="100px">
        <el-form-item label="ID">
          <el-input v-model="currentUser.id" disabled />
        </el-form-item>
        <el-form-item label="登录ID">
          <el-input v-model="currentUser.login_id" disabled />
        </el-form-item>
        <el-form-item label="用户名">
          <el-input
            v-model="currentUser.username"
            :disabled="detailDialogMode === 'view'"
          />
        </el-form-item>
        <el-form-item label="部门">
          <el-select
            v-model="currentUser.department"
            :disabled="detailDialogMode === 'view'"
            style="width: 100%"
          >
            <el-option
              v-for="dept in departments"
              :key="dept"
              :label="dept"
              :value="dept"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="角色">
          <el-select
            v-model="currentUser.role"
            :disabled="detailDialogMode === 'view'"
            style="width: 100%"
          >
            <el-option
              v-for="role in roles"
              :key="role"
              :label="getRoleLabel(role)"
              :value="role"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input
            v-model="currentUser.email"
            :disabled="detailDialogMode === 'view'"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select
            v-model="currentUser.status"
            :disabled="detailDialogMode === 'view'"
            style="width: 100%"
          >
            <el-option label="激活" value="active" />
            <el-option label="未激活" value="inactive" />
            <el-option label="停用" value="suspended" />
          </el-select>
        </el-form-item>
        <el-form-item label="创建时间">
          <el-input :value="formatDate(currentUser.created_at)" disabled />
        </el-form-item>
        <el-form-item label="更新时间">
          <el-input :value="formatDate(currentUser.updated_at)" disabled />
        </el-form-item>
      </el-form>
      <template #footer v-if="detailDialogMode === 'edit'">
        <el-button @click="detailDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleUpdate" :loading="updating">
          保存
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getUserListApi,
  createUserApi,
  getDepartmentsApi,
  getRolesApi,
  updateUserApi,
  deleteUserApi,
} from '#/api/core/user'

// 数据
const users = ref<any[]>([])
const departments = ref<string[]>([])
const roles = ref<string[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 新建用户对话框
const createDialogVisible = ref(false)
const creating = ref(false)
const createForm = ref({
  username: '',
  department: '',
  role: '',
  email: '',
})

// 创建成功对话框
const successDialogVisible = ref(false)
const createdUser = ref({
  login_id: '',
  password: '',
})

// 查看/编辑对话框
const detailDialogVisible = ref(false)
const detailDialogMode = ref<'view' | 'edit'>('view')
const currentUser = ref<any>({})
const updating = ref(false)

// 角色标签映射
const roleLabels: Record<string, string> = {
  admin: '管理员',
  hr: '人力资源',
  financeDirector: '财务总监',
  finance: '财务',
  productionDirector: '生产总监',
  productionSpecialist: '生产专员',
  orderCoordinator: '订单协调员',
  salesManager: '销售经理',
  salesAssistant: '销售助理',
  user: '普通用户',
}

// 获取角色标签
const getRoleLabel = (role: string) => {
  return roleLabels[role] || role
}

// 格式化日期
const formatDate = (date: string) => {
  return date ? new Date(date).toLocaleString('zh-CN') : '-'
}

// 加载用户列表
const loadUsers = async () => {
  loading.value = true
  try {
    const offset = (currentPage.value - 1) * pageSize.value
    const res = await getUserListApi({ limit: pageSize.value, offset })
    users.value = res.users || []
    total.value = res.total || 0
  } catch (error) {
    console.error('加载用户列表失败:', error)
    ElMessage.error('加载用户列表失败')
  } finally {
    loading.value = false
  }
}

// 加载部门列表
const loadDepartments = async () => {
  try {
    const res = await getDepartmentsApi()
    departments.value = res.departments || []
  } catch (error) {
    console.error('加载部门列表失败:', error)
  }
}

// 加载角色列表
const loadRoles = async () => {
  try {
    const res = await getRolesApi()
    roles.value = res.roles || []
  } catch (error) {
    console.error('加载角色列表失败:', error)
  }
}

// 新建用户
const handleCreate = () => {
  createForm.value = {
    username: '',
    department: '',
    role: '',
    email: '',
  }
  createDialogVisible.value = true
}

// 提交新建用户
const handleCreateSubmit = async () => {
  if (
    !createForm.value.username ||
    !createForm.value.department ||
    !createForm.value.role
  ) {
    ElMessage.warning('请填写所有必填字段')
    return
  }

  creating.value = true
  try {
    const res = await createUserApi(createForm.value)
    createdUser.value = res
    createDialogVisible.value = false
    successDialogVisible.value = true
    loadUsers()
  } catch (error: any) {
    console.error('创建用户失败:', error)
    ElMessage.error(error.response?.data?.error || '创建用户失败')
  } finally {
    creating.value = false
  }
}

// 关闭新建对话框
const handleCreateDialogClose = () => {
  createForm.value = {
    username: '',
    department: '',
    role: '',
    email: '',
  }
}

// 关闭成功对话框
const handleSuccessDialogClose = () => {
  successDialogVisible.value = false
  createdUser.value = {
    login_id: '',
    password: '',
  }
}

// 查看用户
const handleView = (row: any) => {
  detailDialogMode.value = 'view'
  currentUser.value = { ...row }
  detailDialogVisible.value = true
}

// 编辑用户
const handleEdit = (row: any) => {
  detailDialogMode.value = 'edit'
  currentUser.value = { ...row }
  detailDialogVisible.value = true
}

// 保存编辑
const handleUpdate = async () => {
  updating.value = true
  try {
    const { id, login_id, created_at, updated_at, has_init_pass, ...updateData } =
      currentUser.value

    await updateUserApi(id, updateData)

    ElMessage.success('保存成功')
    detailDialogVisible.value = false
    loadUsers()
  } catch (error: any) {
    console.error('保存失败:', error)
    ElMessage.error(error.response?.data?.error || '保存失败')
  } finally {
    updating.value = false
  }
}

// 关闭详情对话框
const handleDetailDialogClose = () => {
  currentUser.value = {}
}

// 删除用户
const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除用户 "${row.username}" 吗？此操作不可恢复！`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    await deleteUserApi(row.id)
    ElMessage.success('删除成功')
    loadUsers()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error(error.response?.data?.error || '删除失败')
    }
  }
}

// 分页
const handlePageChange = () => {
  loadUsers()
}

const handleSizeChange = () => {
  currentPage.value = 1
  loadUsers()
}

// 初始化
onMounted(() => {
  loadUsers()
  loadDepartments()
  loadRoles()
})
</script>

<style scoped lang="scss">
.user-management {
  padding: 20px;
  background: #fff;
  min-height: 100vh;

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    h2 {
      margin: 0;
      font-size: 20px;
      color: #303133;
    }
  }
}
</style>
