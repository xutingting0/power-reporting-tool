<template>
  <div class="system-container">
    <TheHeader />
    <el-main>
      <el-page-header>
        <template #content>
          <span>系统管理</span>
        </template>
      </el-page-header>

      <el-row :gutter="20" class="system-row">
        <!-- 系统信息 -->
        <el-col :span="8">
          <el-card class="system-card">
            <template #header>
              <span>系统信息</span>
            </template>
            <div class="system-info">
              <div class="info-item">
                <span class="label">系统名称：</span>
                <span class="value">电力监控系统工作报备工具</span>
              </div>
              <div class="info-item">
                <span class="label">版本号：</span>
                <span class="value">v1.0.0</span>
              </div>
              <div class="info-item">
                <span class="label">后端服务：</span>
                <el-tag type="success" size="small">运行正常</el-tag>
              </div>
              <div class="info-item">
                <span class="label">数据库：</span>
                <el-tag type="success" size="small">SQLite</el-tag>
              </div>
              <div class="info-item">
                <span class="label">启动时间：</span>
                <span class="value">{{ systemInfo.startTime }}</span>
              </div>
              <div class="info-item">
                <span class="label">运行时长：</span>
                <span class="value">{{ systemInfo.uptime }}</span>
              </div>
            </div>
          </el-card>
        </el-col>

        <!-- 数据统计 -->
        <el-col :span="8">
          <el-card class="system-card">
            <template #header>
              <span>数据统计</span>
            </template>
            <div class="data-stats">
              <div class="stat-item">
                <el-statistic title="用户总数" :value="systemInfo.userCount">
                  <template #suffix>
                    <el-icon><User /></el-icon>
                  </template>
                </el-statistic>
              </div>
              <div class="stat-item">
                <el-statistic
                  title="报备单总数"
                  :value="systemInfo.reportCount"
                >
                  <template #suffix>
                    <el-icon><Document /></el-icon>
                  </template>
                </el-statistic>
              </div>
              <div class="stat-item">
                <el-statistic
                  title="模板总数"
                  :value="systemInfo.templateCount"
                >
                  <template #suffix>
                    <el-icon><Files /></el-icon>
                  </template>
                </el-statistic>
              </div>
            </div>
          </el-card>
        </el-col>

        <!-- 系统状态 -->
        <el-col :span="8">
          <el-card class="system-card">
            <template #header>
              <span>系统状态</span>
            </template>
            <div class="system-status">
              <div class="status-item">
                <span class="label">CPU使用率：</span>
                <el-progress
                  :percentage="systemStatus.cpu"
                  :color="getStatusColor(systemStatus.cpu)"
                />
              </div>
              <div class="status-item">
                <span class="label">内存使用率：</span>
                <el-progress
                  :percentage="systemStatus.memory"
                  :color="getStatusColor(systemStatus.memory)"
                />
              </div>
              <div class="status-item">
                <span class="label">磁盘使用率：</span>
                <el-progress
                  :percentage="systemStatus.disk"
                  :color="getStatusColor(systemStatus.disk)"
                />
              </div>
            </div>
          </el-card>
        </el-col>

        <!-- 系统设置 -->
        <el-col :span="12">
          <el-card class="system-card">
            <template #header>
              <span>系统设置</span>
            </template>
            <div class="system-settings">
              <el-form :model="settingsForm" label-width="120px">
                <el-form-item label="系统名称">
                  <el-input
                    v-model="settingsForm.systemName"
                    placeholder="请输入系统名称"
                  />
                </el-form-item>
                <el-form-item label="默认密码">
                  <el-input
                    v-model="settingsForm.defaultPassword"
                    type="password"
                    placeholder="新用户的默认密码"
                    show-password
                  />
                </el-form-item>
                <el-form-item label="会话超时">
                  <el-input-number
                    v-model="settingsForm.sessionTimeout"
                    :min="30"
                    :max="480"
                    :step="30"
                  />
                  <span style="margin-left: 10px; color: #666">分钟</span>
                </el-form-item>
                <el-form-item label="数据保留">
                  <el-select v-model="settingsForm.dataRetention">
                    <el-option label="1个月" value="1" />
                    <el-option label="3个月" value="3" />
                    <el-option label="6个月" value="6" />
                    <el-option label="1年" value="12" />
                    <el-option label="永久" value="0" />
                  </el-select>
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" @click="saveSettings"
                    >保存设置</el-button
                  >
                  <el-button @click="resetSettings">重置</el-button>
                </el-form-item>
              </el-form>
            </div>
          </el-card>
        </el-col>

        <!-- 系统操作 -->
        <el-col :span="12">
          <el-card class="system-card">
            <template #header>
              <span>系统操作</span>
            </template>
            <div class="system-actions">
              <div class="action-item">
                <el-button
                  type="primary"
                  @click="handleBackup"
                  :icon="Download"
                >
                  备份数据库
                </el-button>
                <span class="action-desc">备份当前所有数据到本地文件</span>
              </div>
              <div class="action-item">
                <el-button
                  type="warning"
                  @click="handleClearCache"
                  :icon="Delete"
                >
                  清除缓存
                </el-button>
                <span class="action-desc">清除系统缓存，释放内存</span>
              </div>
              <div class="action-item">
                <el-button
                  type="danger"
                  @click="handleSystemRestart"
                  :icon="Refresh"
                >
                  重启系统
                </el-button>
                <span class="action-desc">重启系统服务，需要管理员权限</span>
              </div>
              <div class="action-item">
                <el-button
                  type="info"
                  @click="handleExportLogs"
                  :icon="Document"
                >
                  导出日志
                </el-button>
                <span class="action-desc">导出系统运行日志</span>
              </div>
            </div>
          </el-card>
        </el-col>

        <!-- 系统日志 -->
        <el-col :span="24">
          <el-card class="system-card">
            <template #header>
              <span>最近系统日志</span>
              <el-button type="text" @click="refreshLogs" :icon="Refresh">
                刷新
              </el-button>
            </template>
            <el-table :data="systemLogs" style="width: 100%" height="300">
              <el-table-column prop="time" label="时间" width="180" />
              <el-table-column prop="level" label="级别" width="100">
                <template #default="{ row }">
                  <el-tag :type="getLogLevelType(row.level)" size="small">
                    {{ row.level }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="module" label="模块" width="120" />
              <el-table-column prop="message" label="内容" />
              <el-table-column prop="user" label="用户" width="120" />
            </el-table>
          </el-card>
        </el-col>
      </el-row>
    </el-main>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  User,
  Document,
  Files,
  Download,
  Delete,
  Refresh,
} from "@element-plus/icons-vue";
import TheHeader from "../../components/TheHeader.vue";

export default {
  name: "System",
  components: {
    TheHeader,
  },
  setup() {
    const systemInfo = reactive({
      startTime: "",
      uptime: "",
      userCount: 0,
      reportCount: 0,
      templateCount: 0,
    });

    const systemStatus = reactive({
      cpu: 0,
      memory: 0,
      disk: 0,
    });

    const settingsForm = reactive({
      systemName: "电力监控系统工作报备工具",
      defaultPassword: "123456",
      sessionTimeout: 120,
      dataRetention: "6",
    });

    const systemLogs = ref([]);

    const formatDateTime = (dateTime) => {
      if (!dateTime) return "";
      return new Date(dateTime).toLocaleString("zh-CN", { hour12: false });
    };

    const getStatusColor = (percentage) => {
      if (percentage < 60) return "#67c23a";
      if (percentage < 80) return "#e6a23c";
      return "#f56c6c";
    };

    const getLogLevelType = (level) => {
      const typeMap = {
        INFO: "success",
        WARN: "warning",
        ERROR: "danger",
        DEBUG: "info",
      };
      return typeMap[level] || "info";
    };

    const loadSystemInfo = () => {
      // 模拟系统信息
      const startTime = new Date();
      startTime.setHours(startTime.getHours() - 2); // 2小时前启动

      systemInfo.startTime = formatDateTime(startTime);
      systemInfo.uptime = "2小时15分钟";
      systemInfo.userCount = 24;
      systemInfo.reportCount = 156;
      systemInfo.templateCount = 8;

      // 模拟系统状态
      systemStatus.cpu = 35;
      systemStatus.memory = 62;
      systemStatus.disk = 42;

      // 模拟系统日志
      systemLogs.value = [
        {
          time: formatDateTime(new Date()),
          level: "INFO",
          module: "用户管理",
          message: "管理员创建了新用户: testuser",
          user: "admin",
        },
        {
          time: formatDateTime(new Date(Date.now() - 300000)),
          level: "WARN",
          module: "报备管理",
          message: "检测到异常登录尝试",
          user: "system",
        },
        {
          time: formatDateTime(new Date(Date.now() - 600000)),
          level: "INFO",
          module: "系统",
          message: "系统备份完成",
          user: "system",
        },
        {
          time: formatDateTime(new Date(Date.now() - 900000)),
          level: "ERROR",
          module: "数据库",
          message: "数据库连接超时，已自动重连",
          user: "system",
        },
        {
          time: formatDateTime(new Date(Date.now() - 1200000)),
          level: "INFO",
          module: "模板管理",
          message: "用户user1使用了模板创建报备单",
          user: "user1",
        },
      ];
    };

    const saveSettings = () => {
      ElMessage.success("系统设置已保存");
    };

    const resetSettings = () => {
      Object.assign(settingsForm, {
        systemName: "电力监控系统工作报备工具",
        defaultPassword: "123456",
        sessionTimeout: 120,
        dataRetention: "6",
      });
      ElMessage.info("设置已重置为默认值");
    };

    const handleBackup = async () => {
      try {
        await ElMessageBox.confirm(
          "确定要备份数据库吗？备份文件将保存在服务器上。",
          "确认备份",
          {
            confirmButtonText: "确定",
            cancelButtonText: "取消",
            type: "warning",
          }
        );

        // 模拟备份过程
        ElMessage.info("开始备份数据库...");
        setTimeout(() => {
          ElMessage.success("数据库备份完成");
        }, 2000);
      } catch (error) {
        if (error === "cancel") return;
      }
    };

    const handleClearCache = async () => {
      try {
        await ElMessageBox.confirm(
          "确定要清除系统缓存吗？这不会影响用户数据。",
          "确认清除缓存",
          {
            confirmButtonText: "确定",
            cancelButtonText: "取消",
            type: "warning",
          }
        );

        ElMessage.success("缓存清除完成");
      } catch (error) {
        if (error === "cancel") return;
      }
    };

    const handleSystemRestart = async () => {
      try {
        await ElMessageBox.confirm(
          "确定要重启系统吗？重启期间用户将无法访问系统。",
          "确认重启系统",
          {
            confirmButtonText: "确定",
            cancelButtonText: "取消",
            type: "danger",
          }
        );

        ElMessage.info("系统正在重启...");
        setTimeout(() => {
          ElMessage.success("系统重启完成");
          // 在实际项目中，这里会触发系统重启
        }, 3000);
      } catch (error) {
        if (error === "cancel") return;
      }
    };

    const handleExportLogs = async () => {
      try {
        await ElMessageBox.confirm("确定要导出系统日志吗？", "确认导出日志", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "info",
        });

        // 模拟导出过程
        ElMessage.info("正在导出日志文件...");
        setTimeout(() => {
          ElMessage.success("日志导出完成");
        }, 1500);
      } catch (error) {
        if (error === "cancel") return;
      }
    };

    const refreshLogs = () => {
      loadSystemInfo();
      ElMessage.success("日志已刷新");
    };

    onMounted(() => {
      loadSystemInfo();
    });

    return {
      systemInfo,
      systemStatus,
      settingsForm,
      systemLogs,
      getStatusColor,
      getLogLevelType,
      saveSettings,
      resetSettings,
      handleBackup,
      handleClearCache,
      handleSystemRestart,
      handleExportLogs,
      refreshLogs,
      User,
      Document,
      Files,
      Download,
      Delete,
      Refresh,
    };
  },
};
</script>

<style scoped>
.system-container {
  padding: 20px;
  background-color: #f5f7fa;
  min-height: 100vh;
}

.system-row {
  margin-top: 20px;
}

.system-card {
  margin-bottom: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.system-info .info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  padding-bottom: 10px;
  border-bottom: 1px solid #eee;
}

.system-info .info-item:last-child {
  border-bottom: none;
}

.system-info .label {
  color: #666;
  font-weight: 500;
}

.system-info .value {
  color: #333;
}

.data-stats .stat-item {
  margin-bottom: 20px;
}

.data-stats .stat-item:last-child {
  margin-bottom: 0;
}

.system-status .status-item {
  margin-bottom: 20px;
}

.system-status .status-item:last-child {
  margin-bottom: 0;
}

.system-status .label {
  display: block;
  margin-bottom: 8px;
  color: #666;
}

.system-actions .action-item {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #eee;
}

.system-actions .action-item:last-child {
  margin-bottom: 0;
  padding-bottom: 0;
  border-bottom: none;
}

.action-desc {
  margin-left: 15px;
  color: #666;
  font-size: 14px;
}
</style>