<template>
  <div class="manage-users-container">
    <TheHeader />
    <el-main>
      <el-page-header>
        <template #content>
          <span>账号管理</span>
          <el-tag type="info" size="small" style="margin-left: 10px">
            共 {{ pagination.total }} 个用户
          </el-tag>
        </template>
      </el-page-header>
      <el-card>
        <!-- 查询表单 -->
        <div class="search-form">
          <el-form :model="searchForm" label-width="80px" inline>
            <el-form-item label="账号">
              <el-input
                v-model="searchForm.username"
                placeholder="请输入账号"
                clearable
              />
            </el-form-item>
            <el-form-item label="姓名">
              <el-input
                v-model="searchForm.realName"
                placeholder="请输入姓名"
                clearable
              />
            </el-form-item>
            <el-form-item label="角色">
              <el-select
                v-model="searchForm.role"
                placeholder="请选择角色"
                clearable
              >
                <el-option label="全部" value="" />
                <el-option label="管理员" value="admin" />
                <el-option label="普通用户" value="user" />
              </el-select>
            </el-form-item>
            <el-form-item label="单位">
              <el-select
                v-model="searchForm.unit"
                placeholder="请选择单位"
                clearable
                filterable
              >
                <el-option
                  v-for="unit in unitOptions"
                  :key="unit.value"
                  :label="unit.label"
                  :value="unit.value"
                />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSearch" :icon="Search">
                查询
              </el-button>
              <el-button @click="resetSearch" :icon="Refresh">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 批量操作 -->
        <div class="batch-operation" v-if="selectedUsers.length > 0">
          <el-space>
            <span>已选择 {{ selectedUsers.length }} 个用户</span>
            <el-button type="danger" @click="batchDelete">批量删除</el-button>
            <el-button type="warning" @click="batchResetPassword">
              批量重置密码
            </el-button>
          </el-space>
        </div>

        <!-- 用户表格 -->
        <el-table
          :data="userData"
          style="width: 100%"
          v-loading="loading"
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column prop="username" label="账号" width="120" sortable />
          <el-table-column prop="real_name" label="姓名" width="100" />
          <el-table-column prop="phone" label="电话" width="130" />
          <el-table-column prop="position" label="岗位" width="120" />
          <el-table-column prop="role" label="角色" width="100">
            <template #default="{ row }">
              <el-tag
                :type="row.role === 'admin' ? 'danger' : 'success'"
                size="small"
              >
                {{ row.role === "admin" ? "管理员" : "普通用户" }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="unit" label="单位" width="150" />
          <el-table-column label="创建时间" width="150">
            <template #default="{ row }">
              {{ formatDateTime(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="状态" width="80">
            <template #default="{ row }">
              <el-tag type="success" size="small">正常</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" size="small" @click="handleEdit(row)">
                编辑
              </el-button>
              <el-button
                type="warning"
                size="small"
                @click="handleResetPassword(row)"
              >
                重置密码
              </el-button>
              <el-button type="danger" size="small" @click="handleDelete(row)">
                删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <div class="pagination-container">
          <el-pagination
            v-model:current-page="pagination.current"
            v-model:page-size="pagination.size"
            :total="pagination.total"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>

        <!-- 新增用户表单 -->
        <div class="add-user-form">
          <h3>新增用户</h3>
          <el-form
            ref="userFormRef"
            :model="userForm"
            :rules="userRules"
            label-width="100px"
            inline
          >
            <el-form-item label="账号" prop="username" required>
              <el-input
                v-model="userForm.username"
                placeholder="请输入账号"
                clearable
              />
            </el-form-item>
            <el-form-item label="密码" prop="password" required>
              <el-input
                v-model="userForm.password"
                type="password"
                placeholder="请输入密码"
                show-password
              />
            </el-form-item>
            <el-form-item label="确认密码" prop="confirmPassword" required>
              <el-input
                v-model="userForm.confirmPassword"
                type="password"
                placeholder="请确认密码"
                show-password
              />
            </el-form-item>
            <el-form-item label="姓名" prop="real_name" required>
              <el-input
                v-model="userForm.real_name"
                placeholder="请输入姓名"
                clearable
              />
            </el-form-item>
            <el-form-item label="电话" prop="phone">
              <el-input
                v-model="userForm.phone"
                placeholder="请输入电话"
                clearable
              />
            </el-form-item>
            <el-form-item label="岗位" prop="position">
              <el-input
                v-model="userForm.position"
                placeholder="请输入岗位"
                clearable
              />
            </el-form-item>
            <el-form-item label="角色" prop="role" required>
              <el-select v-model="userForm.role" placeholder="请选择角色">
                <el-option label="普通用户" value="user" />
                <el-option label="管理员" value="admin" />
              </el-select>
            </el-form-item>
            <el-form-item label="单位" prop="unit">
              <el-select
                v-model="userForm.unit"
                placeholder="请选择单位"
                filterable
                clearable
              >
                <el-option
                  v-for="unit in unitOptions"
                  :key="unit.value"
                  :label="unit.label"
                  :value="unit.value"
                />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="addUser" :icon="Plus">
                添加用户
              </el-button>
              <el-button @click="resetUserForm" :icon="Refresh">
                重置
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-card>
    </el-main>

    <!-- 编辑用户对话框 -->
    <el-dialog
      v-model="dialogVisible.edit"
      :title="`编辑用户 - ${editForm.real_name || ''}`"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="editFormRef"
        :model="editForm"
        :rules="userRules"
        label-width="100px"
      >
        <el-form-item label="账号" prop="username">
          <el-input v-model="editForm.username" readonly />
        </el-form-item>
        <el-form-item label="姓名" prop="real_name" required>
          <el-input v-model="editForm.real_name" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="电话" prop="phone">
          <el-input v-model="editForm.phone" placeholder="请输入电话" />
        </el-form-item>
        <el-form-item label="岗位" prop="position">
          <el-input v-model="editForm.position" placeholder="请输入岗位" />
        </el-form-item>
        <el-form-item label="角色" prop="role" required>
          <el-select v-model="editForm.role" placeholder="请选择角色">
            <el-option label="普通用户" value="user" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-form-item>
        <el-form-item label="单位" prop="unit">
          <el-select
            v-model="editForm.unit"
            placeholder="请选择单位"
            filterable
          >
            <el-option
              v-for="unit in unitOptions"
              :key="unit.value"
              :label="unit.label"
              :value="unit.value"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible.edit = false">取消</el-button>
        <el-button type="primary" @click="updateUser">保存</el-button>
      </template>
    </el-dialog>

    <!-- 重置密码对话框 -->
    <el-dialog
      v-model="dialogVisible.resetPassword"
      title="重置密码"
      width="400px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="passwordFormRef"
        :model="passwordForm"
        :rules="passwordRules"
        label-width="80px"
      >
        <el-form-item label="新密码" prop="newPassword" required>
          <el-input
            v-model="passwordForm.newPassword"
            type="password"
            placeholder="请输入新密码"
            show-password
          />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword" required>
          <el-input
            v-model="passwordForm.confirmPassword"
            type="password"
            placeholder="请确认新密码"
            show-password
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible.resetPassword = false">取消</el-button>
        <el-button type="primary" @click="confirmResetPassword">
          确定重置
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { Search, Refresh, Plus } from "@element-plus/icons-vue";
import TheHeader from "../../components/TheHeader.vue";
import { userAPI } from "@/utils/request.js";

export default {
  name: "ManageUsers",
  components: {
    TheHeader,
  },
  setup() {
    const userFormRef = ref();
    const editFormRef = ref();
    const passwordFormRef = ref();

    const userData = ref([]);
    const loading = ref(false);
    const selectedUsers = ref([]);

    const searchForm = reactive({
      username: "",
      realName: "",
      role: "",
      unit: "",
    });

    const pagination = reactive({
      current: 1,
      size: 10,
      total: 0,
    });

    const userForm = reactive({
      username: "",
      password: "",
      confirmPassword: "",
      real_name: "",
      phone: "",
      position: "",
      role: "user",
      unit: "",
    });

    const editForm = reactive({
      id: null,
      username: "",
      real_name: "",
      phone: "",
      position: "",
      role: "",
      unit: "",
    });

    const passwordForm = reactive({
      userId: null,
      username: "",
      newPassword: "",
      confirmPassword: "",
    });

    const dialogVisible = reactive({
      edit: false,
      resetPassword: false,
    });

    const unitOptions = [
      { value: "遵义供电局", label: "遵义供电局" },
      { value: "贵阳供电局", label: "贵阳供电局" },
      { value: "六盘水供电局", label: "六盘水供电局" },
      { value: "安顺供电局", label: "安顺供电局" },
      { value: "毕节供电局", label: "毕节供电局" },
      { value: "铜仁供电局", label: "铜仁供电局" },
      { value: "黔西南供电局", label: "黔西南供电局" },
      { value: "黔东南供电局", label: "黔东南供电局" },
      { value: "黔南供电局", label: "黔南供电局" },
    ];

    const userRules = {
      username: [
        { required: true, message: "请输入账号", trigger: "blur" },
        { min: 3, max: 20, message: "账号长度在3-20个字符", trigger: "blur" },
      ],
      password: [
        { required: true, message: "请输入密码", trigger: "blur" },
        { min: 6, max: 20, message: "密码长度在6-20个字符", trigger: "blur" },
      ],
      confirmPassword: [
        { required: true, message: "请确认密码", trigger: "blur" },
        {
          validator: (rule, value, callback) => {
            if (value !== userForm.password) {
              callback(new Error("两次输入密码不一致"));
            } else {
              callback();
            }
          },
          trigger: "blur",
        },
      ],
      real_name: [{ required: true, message: "请输入姓名", trigger: "blur" }],
      phone: [
        {
          pattern: /^1[3-9]\d{9}$/,
          message: "请输入正确的手机号码",
          trigger: "blur",
        },
      ],
      role: [{ required: true, message: "请选择角色", trigger: "change" }],
    };

    const passwordRules = {
      newPassword: [
        { required: true, message: "请输入新密码", trigger: "blur" },
        { min: 6, max: 20, message: "密码长度在6-20个字符", trigger: "blur" },
      ],
      confirmPassword: [
        { required: true, message: "请确认新密码", trigger: "blur" },
        {
          validator: (rule, value, callback) => {
            if (value !== passwordForm.newPassword) {
              callback(new Error("两次输入密码不一致"));
            } else {
              callback();
            }
          },
          trigger: "blur",
        },
      ],
    };

    const formatDateTime = (dateTime) => {
      if (!dateTime) return "";
      return new Date(dateTime).toLocaleString("zh-CN", { hour12: false });
    };

    const loadUsers = async () => {
      loading.value = true;
      try {
        const params = {
          page: pagination.current,
          page_size: pagination.size,
          username: searchForm.username,
          real_name: searchForm.realName,
          role: searchForm.role,
          unit: searchForm.unit,
        };

        const response = await userAPI.getAllUsers(params);
        userData.value = response.data || response;
        pagination.total = response.total || response.length;
      } catch (error) {
        console.error("加载用户失败:", error);
        ElMessage.error("加载用户失败");
      } finally {
        loading.value = false;
      }
    };

    const handleSearch = () => {
      pagination.current = 1;
      loadUsers();
    };

    const resetSearch = () => {
      Object.assign(searchForm, {
        username: "",
        realName: "",
        role: "",
        unit: "",
      });
      pagination.current = 1;
      loadUsers();
    };

    const handleSelectionChange = (selection) => {
      selectedUsers.value = selection;
    };

    const addUser = async () => {
      try {
        await userFormRef.value.validate();

        // 检查账号是否已存在
        const existingUser = userData.value.find(
          (user) => user.username === userForm.username
        );
        if (existingUser) {
          ElMessage.warning("该账号已存在");
          return;
        }

        const userDataToSend = {
          username: userForm.username,
          password: userForm.password,
          real_name: userForm.real_name,
          phone: userForm.phone,
          position: userForm.position,
          role: userForm.role,
          unit: userForm.unit,
        };

        await userAPI.createUser(userDataToSend);
        ElMessage.success("用户添加成功");
        resetUserForm();
        loadUsers();
      } catch (error) {
        console.error("添加用户失败:", error);
        ElMessage.error("添加用户失败");
      }
    };

    const resetUserForm = () => {
      Object.assign(userForm, {
        username: "",
        password: "",
        confirmPassword: "",
        real_name: "",
        phone: "",
        position: "",
        role: "user",
        unit: "",
      });
      userFormRef.value?.clearValidate();
    };

    const handleEdit = (row) => {
      Object.assign(editForm, {
        id: row.id,
        username: row.username,
        real_name: row.real_name,
        phone: row.phone,
        position: row.position,
        role: row.role,
        unit: row.unit,
      });
      dialogVisible.edit = true;
    };

    const updateUser = async () => {
      try {
        await editFormRef.value.validate();

        const userDataToUpdate = {
          real_name: editForm.real_name,
          phone: editForm.phone,
          position: editForm.position,
          role: editForm.role,
          unit: editForm.unit,
        };

        await userAPI.updateUser(editForm.username, userDataToUpdate);
        ElMessage.success("用户信息更新成功");
        dialogVisible.edit = false;
        loadUsers();
      } catch (error) {
        console.error("更新用户失败:", error);
        ElMessage.error("更新用户失败");
      }
    };

    const handleDelete = async (row) => {
      try {
        // 不能删除自己
        if (row.username === localStorage.getItem("user")) {
          ElMessage.warning("不能删除当前登录用户");
          return;
        }

        await ElMessageBox.confirm(
          `确定要删除用户 "${row.real_name}" (${row.username}) 吗？`,
          "确认删除",
          {
            confirmButtonText: "确定",
            cancelButtonText: "取消",
            type: "warning",
          }
        );

        await userAPI.deleteUser(row.username);
        ElMessage.success("用户删除成功");
        loadUsers();
      } catch (error) {
        if (error === "cancel") return;
        console.error("删除用户失败:", error);
        ElMessage.error("删除用户失败");
      }
    };

    const batchDelete = async () => {
      if (selectedUsers.value.length === 0) {
        ElMessage.warning("请选择要删除的用户");
        return;
      }

      try {
        // 检查是否包含当前登录用户
        const currentUser = localStorage.getItem("user");
        const hasCurrentUser = selectedUsers.value.some(
          (user) => user.username === currentUser
        );

        if (hasCurrentUser) {
          ElMessage.warning("不能删除当前登录用户");
          return;
        }

        await ElMessageBox.confirm(
          `确定要删除选中的 ${selectedUsers.value.length} 个用户吗？`,
          "确认批量删除",
          {
            confirmButtonText: "确定",
            cancelButtonText: "取消",
            type: "warning",
          }
        );

        const deletePromises = selectedUsers.value.map((user) =>
          userAPI.deleteUser(user.username)
        );
        await Promise.all(deletePromises);
        ElMessage.success("批量删除成功");
        selectedUsers.value = [];
        loadUsers();
      } catch (error) {
        if (error === "cancel") return;
        console.error("批量删除失败:", error);
        ElMessage.error("批量删除失败");
      }
    };

    const handleResetPassword = (row) => {
      Object.assign(passwordForm, {
        userId: row.id,
        username: row.username,
        newPassword: "",
        confirmPassword: "",
      });
      dialogVisible.resetPassword = true;
    };

    const confirmResetPassword = async () => {
      try {
        await passwordFormRef.value.validate();

        await userAPI.updateUser(passwordForm.username, {
          password: passwordForm.newPassword,
        });

        ElMessage.success("密码重置成功");
        dialogVisible.resetPassword = false;
      } catch (error) {
        console.error("重置密码失败:", error);
        ElMessage.error("重置密码失败");
      }
    };

    const batchResetPassword = async () => {
      if (selectedUsers.value.length === 0) {
        ElMessage.warning("请选择要重置密码的用户");
        return;
      }

      try {
        await ElMessageBox.confirm(
          `确定要为选中的 ${selectedUsers.value.length} 个用户重置密码吗？新密码将设置为"123456"`,
          "确认批量重置密码",
          {
            confirmButtonText: "确定",
            cancelButtonText: "取消",
            type: "warning",
          }
        );

        const resetPromises = selectedUsers.value.map((user) =>
          userAPI.updateUser(user.username, { password: "123456" })
        );
        await Promise.all(resetPromises);
        ElMessage.success("批量重置密码成功");
        selectedUsers.value = [];
      } catch (error) {
        if (error === "cancel") return;
        console.error("批量重置密码失败:", error);
        ElMessage.error("批量重置密码失败");
      }
    };

    const handleSizeChange = (size) => {
      pagination.size = size;
      pagination.current = 1;
      loadUsers();
    };

    const handleCurrentChange = (page) => {
      pagination.current = page;
      loadUsers();
    };

    onMounted(() => {
      loadUsers();
    });

    return {
      userFormRef,
      editFormRef,
      passwordFormRef,
      userData,
      loading,
      selectedUsers,
      searchForm,
      pagination,
      userForm,
      editForm,
      passwordForm,
      dialogVisible,
      unitOptions,
      userRules,
      passwordRules,
      formatDateTime,
      handleSearch,
      resetSearch,
      handleSelectionChange,
      addUser,
      resetUserForm,
      handleEdit,
      updateUser,
      handleDelete,
      batchDelete,
      handleResetPassword,
      confirmResetPassword,
      batchResetPassword,
      handleSizeChange,
      handleCurrentChange,
      Search,
      Refresh,
      Plus,
    };
  },
};
</script>

<style scoped>
.manage-users-container {
  padding: 20px;
  background-color: #f5f7fa;
  min-height: 100vh;
}

.search-form {
  margin-bottom: 20px;
  padding: 20px;
  background-color: #fff;
  border-radius: 4px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.batch-operation {
  margin-bottom: 20px;
  padding: 10px 20px;
  background-color: #fff;
  border-radius: 4px;
  border: 1px solid #ebeef5;
}

.add-user-form {
  margin-top: 30px;
  padding: 20px;
  background-color: #fff;
  border-radius: 4px;
  border-top: 1px solid #ebeef5;
}

.add-user-form h3 {
  margin-top: 0;
  margin-bottom: 20px;
  color: #409eff;
}

.pagination-container {
  margin-top: 20px;
  text-align: right;
}
</style>