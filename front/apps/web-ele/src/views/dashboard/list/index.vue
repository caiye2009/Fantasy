<script lang="ts" setup>
import { ref } from 'vue';

// 造一些示例数据
const tableData = ref([
  {
    id: 1,
    username: 'admin',
    realName: '管理员',
    department: '技术部',
    role: '超级管理员',
    email: 'admin@example.com',
    phone: '13800138000',
    status: 1,
    createTime: '2024-01-15 10:30:00',
  },
  {
    id: 2,
    username: 'zhangsan',
    realName: '张三',
    department: '开发部',
    role: '开发工程师',
    email: 'zhangsan@example.com',
    phone: '13800138001',
    status: 1,
    createTime: '2024-02-20 14:20:00',
  },
  {
    id: 3,
    username: 'lisi',
    realName: '李四',
    department: '产品部',
    role: '产品经理',
    email: 'lisi@example.com',
    phone: '13800138002',
    status: 1,
    createTime: '2024-03-10 09:15:00',
  },
  {
    id: 4,
    username: 'wangwu',
    realName: '王五',
    department: '设计部',
    role: 'UI设计师',
    email: 'wangwu@example.com',
    phone: '13800138003',
    status: 0,
    createTime: '2024-04-05 16:45:00',
  },
  {
    id: 5,
    username: 'zhaoliu',
    realName: '赵六',
    department: '测试部',
    role: '测试工程师',
    email: 'zhaoliu@example.com',
    phone: '13800138004',
    status: 1,
    createTime: '2024-05-12 11:00:00',
  },
  {
    id: 6,
    username: 'sunqi',
    realName: '孙七',
    department: '运维部',
    role: '运维工程师',
    email: 'sunqi@example.com',
    phone: '13800138005',
    status: 1,
    createTime: '2024-06-18 13:30:00',
  },
  {
    id: 7,
    username: 'zhouba',
    realName: '周八',
    department: '市场部',
    role: '市场专员',
    email: 'zhouba@example.com',
    phone: '13800138006',
    status: 1,
    createTime: '2024-07-22 15:20:00',
  },
  {
    id: 8,
    username: 'wujiu',
    realName: '吴九',
    department: '人事部',
    role: '人事专员',
    email: 'wujiu@example.com',
    phone: '13800138007',
    status: 0,
    createTime: '2024-08-08 10:10:00',
  },
]);

// 搜索表单数据
const searchForm = ref({
  username: '',
  realName: '',
  department: '',
});

// 分页
const currentPage = ref(1);
const pageSize = ref(10);

// 计算当前页数据
const paginatedData = ref(tableData.value);

// 操作方法
function handleEdit(row: any) {
  alert('编辑用户: ' + row.realName);
}

function handleDelete(row: any) {
  if (confirm(`确定要删除用户 ${row.realName} 吗？`)) {
    const index = tableData.value.findIndex(item => item.id === row.id);
    if (index > -1) {
      tableData.value.splice(index, 1);
      paginatedData.value = tableData.value;
    }
  }
}

function handleView(row: any) {
  alert('查看用户: ' + row.realName);
}

function handleAdd() {
  alert('添加新用户');
}

function handleSearch() {
  console.log('搜索:', searchForm.value);
  // 可以在这里添加搜索逻辑
}

function handleReset() {
  searchForm.value = {
    username: '',
    realName: '',
    department: '',
  };
}
</script>

<template>
  <div class="p-5">
    <!-- 搜索区域 -->
    <div class="mb-4 rounded-lg bg-white p-6 shadow">
      <div class="flex flex-wrap items-end gap-4">
        <div class="flex flex-col">
          <label class="mb-2 text-sm font-medium text-gray-700">用户名</label>
          <input
            v-model="searchForm.username"
            type="text"
            placeholder="请输入用户名"
            class="w-48 rounded-md border border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
          />
        </div>
        <div class="flex flex-col">
          <label class="mb-2 text-sm font-medium text-gray-700">姓名</label>
          <input
            v-model="searchForm.realName"
            type="text"
            placeholder="请输入姓名"
            class="w-48 rounded-md border border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
          />
        </div>
        <div class="flex flex-col">
          <label class="mb-2 text-sm font-medium text-gray-700">部门</label>
          <select
            v-model="searchForm.department"
            class="w-48 rounded-md border border-gray-300 px-3 py-2 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
          >
            <option value="">全部</option>
            <option value="技术部">技术部</option>
            <option value="开发部">开发部</option>
            <option value="产品部">产品部</option>
            <option value="设计部">设计部</option>
            <option value="测试部">测试部</option>
            <option value="运维部">运维部</option>
            <option value="市场部">市场部</option>
            <option value="人事部">人事部</option>
          </select>
        </div>
        <button
          class="rounded-md bg-blue-600 px-6 py-2 text-white transition hover:bg-blue-700"
          @click="handleSearch"
        >
          搜索
        </button>
        <button
          class="rounded-md border border-gray-300 bg-white px-6 py-2 transition hover:bg-gray-50"
          @click="handleReset"
        >
          重置
        </button>
      </div>
    </div>

    <!-- 工具栏和表格 -->
    <div class="rounded-lg bg-white p-6 shadow">
      <!-- 工具栏 -->
      <div class="mb-4 flex items-center justify-between">
        <button
          class="rounded-md bg-blue-600 px-6 py-2 text-white transition hover:bg-blue-700"
          @click="handleAdd"
        >
          + 新增用户
        </button>
        <div class="text-sm text-gray-600">
          共 {{ tableData.length }} 条数据
        </div>
      </div>

      <!-- 表格 -->
      <div class="overflow-x-auto">
        <table class="w-full border-collapse">
          <thead>
            <tr class="bg-gray-50">
              <th class="border border-gray-200 px-4 py-3 text-left text-sm font-semibold text-gray-700">序号</th>
              <th class="border border-gray-200 px-4 py-3 text-left text-sm font-semibold text-gray-700">用户名</th>
              <th class="border border-gray-200 px-4 py-3 text-left text-sm font-semibold text-gray-700">姓名</th>
              <th class="border border-gray-200 px-4 py-3 text-left text-sm font-semibold text-gray-700">部门</th>
              <th class="border border-gray-200 px-4 py-3 text-left text-sm font-semibold text-gray-700">角色</th>
              <th class="border border-gray-200 px-4 py-3 text-left text-sm font-semibold text-gray-700">邮箱</th>
              <th class="border border-gray-200 px-4 py-3 text-left text-sm font-semibold text-gray-700">手机号</th>
              <th class="border border-gray-200 px-4 py-3 text-left text-sm font-semibold text-gray-700">状态</th>
              <th class="border border-gray-200 px-4 py-3 text-left text-sm font-semibold text-gray-700">创建时间</th>
              <th class="border border-gray-200 px-4 py-3 text-left text-sm font-semibold text-gray-700">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr 
              v-for="(item, index) in paginatedData" 
              :key="item.id" 
              class="transition hover:bg-gray-50"
            >
              <td class="border border-gray-200 px-4 py-3 text-sm">{{ index + 1 }}</td>
              <td class="border border-gray-200 px-4 py-3 text-sm">{{ item.username }}</td>
              <td class="border border-gray-200 px-4 py-3 text-sm">{{ item.realName }}</td>
              <td class="border border-gray-200 px-4 py-3 text-sm">{{ item.department }}</td>
              <td class="border border-gray-200 px-4 py-3 text-sm">{{ item.role }}</td>
              <td class="border border-gray-200 px-4 py-3 text-sm">{{ item.email }}</td>
              <td class="border border-gray-200 px-4 py-3 text-sm">{{ item.phone }}</td>
              <td class="border border-gray-200 px-4 py-3 text-sm">
                <span 
                  :class="item.status === 1 ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
                  class="rounded-full px-3 py-1 text-xs font-medium"
                >
                  {{ item.status === 1 ? '正常' : '禁用' }}
                </span>
              </td>
              <td class="border border-gray-200 px-4 py-3 text-sm">{{ item.createTime }}</td>
              <td class="border border-gray-200 px-4 py-3 text-sm">
                <div class="flex gap-3">
                  <button
                    class="text-blue-600 hover:text-blue-800 hover:underline"
                    @click="handleView(item)"
                  >
                    查看
                  </button>
                  <button
                    class="text-blue-600 hover:text-blue-800 hover:underline"
                    @click="handleEdit(item)"
                  >
                    编辑
                  </button>
                  <button
                    class="text-red-600 hover:text-red-800 hover:underline"
                    @click="handleDelete(item)"
                  >
                    删除
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- 分页 -->
      <div class="mt-4 flex items-center justify-between">
        <div class="text-sm text-gray-600">
          显示第 1 - {{ paginatedData.length }} 条，共 {{ tableData.length }} 条
        </div>
        <div class="flex gap-2">
          <button class="rounded border border-gray-300 px-3 py-1 text-sm hover:bg-gray-50">
            上一页
          </button>
          <button class="rounded bg-blue-600 px-3 py-1 text-sm text-white">
            1
          </button>
          <button class="rounded border border-gray-300 px-3 py-1 text-sm hover:bg-gray-50">
            下一页
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 自定义样式 */
</style>