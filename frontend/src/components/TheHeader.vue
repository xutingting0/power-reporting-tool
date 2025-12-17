<template>
  <div class="header-container">
    <el-container>
      <!-- 顶部导航栏 -->
      <el-header class="main-header">
        <!-- Logo区域 -->
        <div class="logo-area">
          <div class="logo">
            <el-icon size="24" color="#409eff"><Lightning /></el-icon>
            <span class="logo-text">电力监控系统工作报备工具</span>
          </div>
          <div class="system-time">
            <el-icon><Clock /></el-icon>
            <span>{{ currentTime }}</span>
          </div>
        </div>

        <!-- 用户信息和操作 -->
        <div class="user-area">
          <el-dropdown @command="handleCommand">
            <div class="user-info">
              <el-avatar :size="36" :src="user.avatar" class="user-avatar">
                {{
                  user.real_name?.charAt(0) || user.username?.charAt(0) || "U"
                }}
              </el-avatar>
              <div class="user-details">
                <div class="user-name">
                  {{ user.real_name || user.username }}
                </div>
                <div class="user-role">
                  <el-tag
                    :type="user.role === 'admin' ? 'danger' : 'success'"
                    size="small"
                  >
                    {{ user.role === "admin" ? "管理员" : "普通用户" }}
                  </el-tag>
                </div>
              </div>
              <el-icon class="user-arrow"><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>
                  <span>个人中心</span>
                </el-dropdown-item>
                <el-dropdown-item command="changePassword" divided>
                  <el-icon><Key /></el-icon>
                  <span>修改密码</span>
                </el-dropdown-item>
                <el-dropdown-item command="logout">
                  <el-icon><SwitchButton /></el-icon>
                  <span>退出登录</span>
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <!-- 水平导航菜单 -->
      <div class="nav-container">
        <el-menu
          :default-active="activeMenu"
          class="horizontal-menu"
          mode="horizontal"
          @select="handleMenuSelect"
        >
          <!-- 管理员菜单 -->
          <template v-if="user.role === 'admin'">
            <el-menu-item index="/admin/approve-report">
              <el-icon><DocumentChecked /></el-icon>
              <span>审核报备单</span>
            </el-menu-item>
            <el-menu-item index="/admin/view-reports">
              <el-icon><Files /></el-icon>
              <span>查看报备单</span>
            </el-menu-item>
            <el-menu-item index="/admin/manage-users">
              <el-icon><UserFilled /></el-icon>
              <span>账号管理</span>
            </el-menu-item>
            <el-menu-item index="/admin/manage-templates">
              <el-icon><Grid /></el-icon>
              <span>模板管理</span>
            </el-menu-item>
            <el-menu-item index="/admin/statistics">
              <el-icon><DataLine /></el-icon>
              <span>统计分析</span>
            </el-menu-item>
            <el-menu-item index="/admin/system">
              <el-icon><Setting /></el-icon>
              <span>系统管理</span>
            </el-menu-item>
          </template>

          <!-- 普通用户菜单 -->
          <template v-else>
            <el-menu-item index="/user/create-report">
              <el-icon><Edit /></el-icon>
              <span>编制报备单</span>
            </el-menu-item>
            <el-menu-item index="/user/manage-reports">
              <el-icon><DocumentCopy /></el-icon>
              <span>管理报备单</span>
            </el-menu-item>
            <el-menu-item index="/user/report-status">
              <el-icon><Search /></el-icon>
              <span>报备单状态</span>
            </el-menu-item>
            <el-menu-item index="/user/profile">
              <el-icon><User /></el-icon>
              <span>个人中心</span>
            </el-menu-item>
          </template>
        </el-menu>

        <!-- 快速操作区域 -->
        <div class="quick-actions" v-if="user.role === 'user'">
          <el-button
            type="primary"
            size="small"
            @click="navigateTo('/user/create-report')"
            :icon="Plus"
          >
            新建报备单
          </el-button>
        </div>
      </div>

      <!-- 面包屑导航 -->
      <div class="breadcrumb-container" v-if="showBreadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item v-for="item in breadcrumbItems" :key="item.path">
            <span
              @click="item.path ? navigateTo(item.path) : null"
              :class="{ 'breadcrumb-link': item.path }"
            >
              {{ item.title }}
            </span>
          </el-breadcrumb-item>
        </el-breadcrumb>
      </div>
    </el-container>
  </div>
</template>

<script>
import { computed, onMounted, onUnmounted, ref } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useStore } from "vuex";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  Lightning,
  Clock,
  User,
  ArrowDown,
  SwitchButton,
  DocumentChecked,
  Files,
  UserFilled,
  Grid,
  DataLine,
  Setting,
  Edit,
  DocumentCopy,
  Search,
  Key,
  Plus,
} from "@element-plus/icons-vue";

export default {
  name: "TheHeader",
  setup() {
    const router = useRouter();
    const route = useRoute();
    const store = useStore();
    const currentTime = ref("");

    // 获取用户信息
    const user = computed(() => store.state.userInfo);

    // 计算当前激活的菜单
    const activeMenu = computed(() => route.path);

    // 面包屑导航
    const breadcrumbItems = computed(() => {
      const items = [];
      const routeName = route.name;
      const routePath = route.path;

      // 根据当前路由生成面包屑
      if (routePath.includes("/admin/")) {
        items.push({ title: "管理员", path: "/admin/approve-report" });

        if (routePath.includes("approve-report")) {
          items.push({ title: "审核报备单", path: "" });
        } else if (routePath.includes("view-reports")) {
          items.push({ title: "查看报备单", path: "" });
        } else if (routePath.includes("manage-users")) {
          items.push({ title: "账号管理", path: "" });
        } else if (routePath.includes("manage-templates")) {
          items.push({ title: "模板管理", path: "" });
        } else if (routePath.includes("statistics")) {
          items.push({ title: "统计分析", path: "" });
        } else if (routePath.includes("system")) {
          items.push({ title: "系统管理", path: "" });
        }
      } else if (routePath.includes("/user/")) {
        items.push({ title: "用户中心", path: "/user/create-report" });

        if (routePath.includes("create-report")) {
          items.push({ title: "编制报备单", path: "" });
        } else if (routePath.includes("manage-reports")) {
          items.push({ title: "管理报备单", path: "" });
        } else if (routePath.includes("report-status")) {
          items.push({ title: "报备单状态", path: "" });
        } else if (routePath.includes("profile")) {
          items.push({ title: "个人中心", path: "" });
        }
      }

      return items;
    });

    const showBreadcrumb = computed(() => {
      return !["/", "/login"].includes(route.path);
    });

    // 更新时间
    const updateTime = () => {
      const now = new Date();
      const year = now.getFullYear();
      const month = String(now.getMonth() + 1).padStart(2, "0");
      const day = String(now.getDate()).padStart(2, "0");
      const hour = String(now.getHours()).padStart(2, "0");
      const minute = String(now.getMinutes()).padStart(2, "0");
      const second = String(now.getSeconds()).padStart(2, "0");
      currentTime.value = `${year}-${month}-${day} ${hour}:${minute}:${second}`;
    };

    // 菜单选择处理
    const handleMenuSelect = (index) => {
      navigateTo(index);
    };

    // 导航到指定路径
    const navigateTo = (path) => {
      if (path && path !== route.path) {
        router.push(path);
      }
    };

    // 下拉菜单命令处理
    const handleCommand = (command) => {
      switch (command) {
        case "profile":
          navigateTo("/user/profile");
          break;
        case "changePassword":
          handleChangePassword();
          break;
        case "logout":
          handleLogout();
          break;
      }
    };

    // 修改密码
    const handleChangePassword = () => {
      ElMessageBox.prompt("请输入新密码", "修改密码", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        inputType: "password",
        inputPlaceholder: "请输入6-20位新密码",
        inputValidator: (value) => {
          if (!value) {
            return "密码不能为空";
          }
          if (value.length < 6 || value.length > 20) {
            return "密码长度应为6-20位";
          }
          return true;
        },
      })
        .then(({ value }) => {
          // 这里应该调用API修改密码
          ElMessage.success("密码修改成功，请重新登录");
          setTimeout(() => {
            handleLogout();
          }, 1500);
        })
        .catch(() => {
          // 用户取消
        });
    };

    // 退出登录
    const handleLogout = () => {
      ElMessageBox.confirm("确定要退出登录吗？", "确认退出", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      })
        .then(() => {
          store.dispatch("logout");
          router.push("/");
          ElMessage.success("已退出登录");
        })
        .catch(() => {
          // 用户取消
        });
    };

    // 初始化时间
    let timer = null;
    onMounted(() => {
      updateTime();
      timer = setInterval(updateTime, 1000);
    });

    onUnmounted(() => {
      if (timer) {
        clearInterval(timer);
      }
    });

    return {
      user,
      currentTime,
      activeMenu,
      breadcrumbItems,
      showBreadcrumb,
      handleMenuSelect,
      handleCommand,
      navigateTo,
      handleChangePassword,
      handleLogout,
      Lightning,
      Clock,
      User,
      ArrowDown,
      SwitchButton,
      DocumentChecked,
      Files,
      UserFilled,
      Grid,
      DataLine,
      Setting,
      Edit,
      DocumentCopy,
      Search,
      Key,
      Plus,
    };
  },
};
</script>

<style scoped>
.header-container {
  background: linear-gradient(135deg, #001529 0%, #002140 100%);
  color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  position: relative;
  z-index: 1000;
}

/* 主头部样式 */
.main-header {
  height: 64px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  background: transparent;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

/* Logo区域 */
.logo-area {
  display: flex;
  align-items: center;
  gap: 40px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 20px;
  font-weight: 600;
  color: #fff;
  transition: opacity 0.3s;
}

.logo:hover {
  opacity: 0.9;
  cursor: pointer;
}

.logo-text {
  background: linear-gradient(90deg, #409eff, #66b1ff);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.system-time {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.8);
}

/* 用户信息区域 */
.user-area {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.user-info:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

.user-avatar {
  background: linear-gradient(135deg, #409eff, #66b1ff);
  font-weight: 600;
}

.user-details {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.user-name {
  font-size: 14px;
  font-weight: 500;
  color: #fff;
}

.user-role {
  margin-top: 2px;
}

.user-arrow {
  color: rgba(255, 255, 255, 0.6);
  font-size: 14px;
}

/* 导航菜单容器 */
.nav-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  background: rgba(0, 21, 41, 0.95);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

/* 水平菜单样式 */
.horizontal-menu {
  flex: 1;
  background: transparent;
  border-bottom: none;
  min-height: 48px;
}

.horizontal-menu :deep(.el-menu-item) {
  height: 48px;
  line-height: 48px;
  color: rgba(255, 255, 255, 0.8);
  border-bottom: 2px solid transparent;
  transition: all 0.3s;
  margin: 0 4px;
  border-radius: 4px 4px 0 0;
}

.horizontal-menu :deep(.el-menu-item) i {
  margin-right: 8px;
  font-size: 16px;
}

.horizontal-menu :deep(.el-menu-item:hover) {
  color: #fff;
  background-color: rgba(64, 158, 255, 0.1);
  border-bottom-color: rgba(64, 158, 255, 0.3);
}

.horizontal-menu :deep(.el-menu-item.is-active) {
  color: #409eff;
  background-color: rgba(64, 158, 255, 0.15);
  border-bottom-color: #409eff;
  font-weight: 500;
}

/* 快速操作区域 */
.quick-actions {
  padding: 8px 0;
}

/* 面包屑导航 */
.breadcrumb-container {
  padding: 12px 20px;
  background: linear-gradient(
    to bottom,
    rgba(0, 21, 41, 0.7),
    rgba(0, 21, 41, 0.4)
  );
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

:deep(.el-breadcrumb) {
  font-size: 14px;
}

:deep(.el-breadcrumb__inner) {
  color: rgba(255, 255, 255, 0.6);
  transition: color 0.3s;
}

:deep(.el-breadcrumb__inner:hover) {
  color: #409eff;
}

:deep(.el-breadcrumb__separator) {
  color: rgba(255, 255, 255, 0.4);
}

.breadcrumb-link {
  cursor: pointer;
  transition: color 0.3s;
}

.breadcrumb-link:hover {
  color: #409eff;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .logo-text {
    font-size: 18px;
  }

  .horizontal-menu :deep(.el-menu-item) {
    margin: 0 2px;
    padding: 0 12px;
  }
}

@media (max-width: 992px) {
  .logo-text {
    display: none;
  }

  .system-time {
    display: none;
  }

  .horizontal-menu :deep(.el-menu-item) {
    padding: 0 8px;
  }

  .horizontal-menu :deep(.el-menu-item span) {
    font-size: 13px;
  }
}

@media (max-width: 768px) {
  .main-header,
  .nav-container {
    padding: 0 12px;
  }

  .user-name {
    display: none;
  }

  .user-role {
    display: none;
  }

  .horizontal-menu :deep(.el-menu-item span) {
    display: none;
  }

  .horizontal-menu :deep(.el-menu-item i) {
    margin-right: 0;
    font-size: 18px;
  }

  .quick-actions {
    display: none;
  }
}
</style>

<style>
/* 全局覆盖el-dropdown的样式 */
.el-dropdown-menu {
  background: linear-gradient(135deg, #001529 0%, #002140 100%);
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.el-dropdown-menu__item {
  color: rgba(255, 255, 255, 0.8);
}

.el-dropdown-menu__item i {
  margin-right: 8px;
  color: #409eff;
}

.el-dropdown-menu__item:hover {
  background-color: rgba(64, 158, 255, 0.1);
  color: #409eff;
}

.el-dropdown-menu__item--divided {
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}
</style>