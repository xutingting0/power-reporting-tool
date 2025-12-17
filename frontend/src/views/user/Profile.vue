<template>
  <div class="profile-container">
    <TheHeader />
    <el-main>
      <el-page-header @back="goBack">
        <template #content>
          <span>个人中心</span>
        </template>
      </el-page-header>

      <el-row :gutter="20">
        <!-- 个人信息 -->
        <el-col :span="8">
          <el-card class="profile-card">
            <template #header>
              <div class="card-header">
                <span>个人信息</span>
                <el-button
                  type="primary"
                  size="small"
                  @click="handleEditProfile"
                >
                  编辑信息
                </el-button>
              </div>
            </template>

            <div class="user-info">
              <div class="avatar-container">
                <el-avatar :size="100" :src="userInfo.avatar || defaultAvatar">
                  {{ userInfo.real_name?.charAt(0) || "U" }}
                </el-avatar>
              </div>

              <div class="info-details">
                <div class="info-item">
                  <span class="label">账号：</span>
                  <span class="value">{{ userInfo.username }}</span>
                </div>
                <div class="info-item">
                  <span class="label">姓名：</span>
                  <span class="value">{{ userInfo.real_name }}</span>
                </div>
                <div class="info-item">
                  <span class="label">角色：</span>
                  <el-tag
                    :type="userInfo.role === 'admin' ? 'danger' : 'success'"
                    size="small"
                  >
                    {{ userInfo.role === "admin" ? "管理员" : "普通用户" }}
                  </el-tag>
                </div>
                <div class="info-item">
                  <span class="label">单位：</span>
                  <span class="value">{{ userInfo.unit }}</span>
                </div>
                <div class="info-item">
                  <span class="label">岗位：</span>
                  <span class="value">{{ userInfo.position }}</span>
                </div>
                <div class="info-item">
                  <span class="label">电话：</span>
                  <span class="value">{{ userInfo.phone }}</span>
                </div>
                <div class="info-item">
                  <span class="label">注册时间：</span>
                  <span class="value">{{
                    formatDateTime(userInfo.created_at)
                  }}</span>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>

        <!-- 我的统计 -->
        <el-col :span="16">
          <el-card class="stats-card">
            <template #header>
              <span>我的统计</span>
            </template>

            <el-row :gutter="20">
              <el-col :span="6">
                <el-statistic title="总报备单数" :value="stats.total">
                  <template #suffix>
                    <el-icon><Document /></el-icon>
                  </template>
                </el-statistic>
              </el-col>
              <el-col :span="6">
                <el-statistic title="草稿" :value="stats.draft">
                  <template #suffix>
                    <el-icon><Edit /></el-icon>
                  </template>
                </el-statistic>
              </el-col>
              <el-col :span="6">
                <el-statistic title="待审核" :value="stats.submitted">
                  <template #suffix>
                    <el-icon><Clock /></el-icon>
                  </template>
                </el-statistic>
              </el-col>
              <el-col :span="6">
                <el-statistic title="已通过" :value="stats.approved">
                  <template #suffix>
                    <el-icon><CircleCheck /></el-icon>
                  </template>
                </el-statistic>
              </el-col>
            </el-row>

            <div class="recent-reports">
              <h4>最近报备单</h4>
              <el-table :data="recentReports" size="small" style="width: 100%">
                <el-table-column
                  prop="report_no"
                  label="报备单号"
                  width="150"
                />
                <el-table-column prop="work_content" label="工作内容" />
                <el-table-column prop="status" label="状态" width="80">
                  <template #default="{ row }">
                    <el-tag :type="getStatusType(row.status)" size="small">
                      {{ getStatusText(row.status) }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="80">
                  <template #default="{ row }">
                    <el-button
                      type="text"
                      size="small"
                      @click="viewReport(row.id)"
                    >
                      查看
                    </el-button>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-card>

          <!-- 修改密码 -->
          <el-card class="password-card" style="margin-top: 20px">
            <template #header>
              <span>修改密码</span>
            </template>

            <el-form
              ref="passwordFormRef"
              :model="passwordForm"
              :rules="passwordRules"
              label-width="100px"
            >
              <el-form-item label="原密码" prop="oldPassword" required>
                <el-input
                  v-model="passwordForm.oldPassword"
                  type="password"
                  placeholder="请输入原密码"
                  show-password
                />
              </el-form-item>
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
              <el-form-item>
                <el-button type="primary" @click="changePassword">
                  修改密码
                </el-button>
              </el-form-item>
            </el-form>
          </el-card>
        </el-col>
      </el-row>
    </el-main>

    <!-- 编辑个人信息对话框 -->
    <el-dialog
      v-model="dialogVisible.editProfile"
      title="编辑个人信息"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="editFormRef"
        :model="editForm"
        :rules="editRules"
        label-width="100px"
      >
        <el-form-item label="姓名" prop="real_name" required>
          <el-input v-model="editForm.real_name" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="电话" prop="phone">
          <el-input v-model="editForm.phone" placeholder="请输入电话" />
        </el-form-item>
        <el-form-item label="岗位" prop="position">
          <el-input v-model="editForm.position" placeholder="请输入岗位" />
        </el-form-item>
        <el-form-item label="单位" prop="unit">
          <el-select
            v-model="editForm.unit"
            placeholder="请选择单位"
            filterable
            style="width: 100%"
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
        <el-button @click="dialogVisible.editProfile = false">取消</el-button>
        <el-button type="primary" @click="updateProfile">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from "vue";
import { useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import { Document, Edit, Clock, CircleCheck } from "@element-plus/icons-vue";
import TheHeader from "../../components/TheHeader.vue";
import { userAPI, reportAPI } from "@/utils/request.js";

export default {
  name: "Profile",
  components: {
    TheHeader,
  },
  setup() {
    const router = useRouter();
    const passwordFormRef = ref();
    const editFormRef = ref();

    const userInfo = ref({});
    const stats = reactive({
      total: 0,
      draft: 0,
      submitted: 0,
      approved: 0,
      rejected: 0,
    });
    const recentReports = ref([]);

    const passwordForm = reactive({
      oldPassword: "",
      newPassword: "",
      confirmPassword: "",
    });

    const editForm = reactive({
      real_name: "",
      phone: "",
      position: "",
      unit: "",
    });

    const dialogVisible = reactive({
      editProfile: false,
    });

    const defaultAvatar =
      "https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png";

    const unitOptions = [
      { value: "遵义供电局", label: "遵义供电局" },
      { value: "贵阳供电局", label: "贵阳供电局" },
      { value: "六盘水供电局", label: "六盘水供电局" },
      { value: "安顺供电局", label: "安顺供电局" },
      { value: "毕节供电局", label: "毕节供电局" },
    ];

    const passwordRules = {
      oldPassword: [
        { required: true, message: "请输入原密码", trigger: "blur" },
      ],
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

    const editRules = {
      real_name: [{ required: true, message: "请输入姓名", trigger: "blur" }],
      phone: [
        {
          pattern: /^1[3-9]\d{9}$/,
          message: "请输入正确的手机号码",
          trigger: "blur",
        },
      ],
    };

    const getStatusText = (status) => {
      const statusMap = {
        draft: "草稿",
        submitted: "待审核",
        approved: "已通过",
        rejected: "已驳回",
      };
      return statusMap[status] || status;
    };

    const getStatusType = (status) => {
      const typeMap = {
        draft: "info",
        submitted: "warning",
        approved: "success",
        rejected: "danger",
      };
      return typeMap[status] || "info";
    };

    const formatDateTime = (dateTime) => {
      if (!dateTime) return "";
      return new Date(dateTime).toLocaleString("zh-CN", { hour12: false });
    };

    const loadUserInfo = async () => {
      try {
        const response = await userAPI.getCurrentUser();
        userInfo.value = response;

        // 初始化编辑表单
        Object.assign(editForm, {
          real_name: response.real_name,
          phone: response.phone || "",
          position: response.position || "",
          unit: response.unit || "",
        });
      } catch (error) {
        console.error("加载用户信息失败:", error);
        ElMessage.error("加载用户信息失败");
      }
    };

    const loadUserStats = async () => {
      try {
        const response = await reportAPI.getUserReports({
          page: 1,
          page_size: 100,
        });

        const reports = response.data || response;

        // 统计
        stats.total = reports.length;
        stats.draft = reports.filter((r) => r.status === "draft").length;
        stats.submitted = reports.filter(
          (r) => r.status === "submitted"
        ).length;
        stats.approved = reports.filter((r) => r.status === "approved").length;
        stats.rejected = reports.filter((r) => r.status === "rejected").length;

        // 最近5条报备单
        recentReports.value = reports.slice(0, 5);
      } catch (error) {
        console.error("加载用户统计失败:", error);
      }
    };

    const handleEditProfile = () => {
      dialogVisible.editProfile = true;
    };

    const updateProfile = async () => {
      try {
        await editFormRef.value.validate();

        await userAPI.updateUserProfile(editForm);
        ElMessage.success("个人信息更新成功");
        dialogVisible.editProfile = false;
        loadUserInfo();
      } catch (error) {
        console.error("更新个人信息失败:", error);
        ElMessage.error("更新个人信息失败");
      }
    };

    const changePassword = async () => {
      try {
        await passwordFormRef.value.validate();

        // 验证原密码
        if (passwordForm.oldPassword !== localStorage.getItem("password")) {
          ElMessage.warning("原密码错误");
          return;
        }

        await userAPI.updateUserProfile({
          password: passwordForm.newPassword,
        });

        ElMessage.success("密码修改成功");

        // 清空表单
        Object.assign(passwordForm, {
          oldPassword: "",
          newPassword: "",
          confirmPassword: "",
        });
        passwordFormRef.value.resetFields();
      } catch (error) {
        console.error("修改密码失败:", error);
        ElMessage.error("修改密码失败");
      }
    };

    const viewReport = (reportId) => {
      router.push({ path: "/user/manage-reports" });
    };

    const goBack = () => {
      router.back();
    };

    onMounted(() => {
      loadUserInfo();
      loadUserStats();
    });

    return {
      passwordFormRef,
      editFormRef,
      userInfo,
      stats,
      recentReports,
      passwordForm,
      editForm,
      dialogVisible,
      defaultAvatar,
      unitOptions,
      passwordRules,
      editRules,
      getStatusText,
      getStatusType,
      formatDateTime,
      handleEditProfile,
      updateProfile,
      changePassword,
      viewReport,
      goBack,
      Document,
      Edit,
      Clock,
      CircleCheck,
    };
  },
};
</script>

<style scoped>
.profile-container {
  padding: 20px;
  background-color: #f5f7fa;
  min-height: 100vh;
}

.profile-card,
.stats-card,
.password-card {
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.user-info {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.avatar-container {
  margin-bottom: 20px;
}

.info-details {
  width: 100%;
}

.info-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #eee;
}

.info-item:last-child {
  border-bottom: none;
}

.info-item .label {
  color: #666;
  font-weight: 500;
}

.info-item .value {
  color: #333;
}

.recent-reports {
  margin-top: 20px;
}

.recent-reports h4 {
  margin-bottom: 10px;
  color: #409eff;
}
</style>