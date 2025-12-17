<template>
  <div class="manage-reports-container">
    <TheHeader />
    <el-main>
      <el-page-header @back="goBack">
        <template #content>
          <span>管理报备单</span>
          <el-tag type="info" size="small" style="margin-left: 10px">
            共 {{ pagination.total }} 条记录
          </el-tag>
        </template>
      </el-page-header>
      <el-card>
        <!-- 查询表单 -->
        <div class="search-form">
          <el-form :model="searchForm" label-width="80px" inline>
            <el-form-item label="报备单号">
              <el-input
                v-model="searchForm.reportNo"
                placeholder="请输入报备单号"
                clearable
              />
            </el-form-item>
            <el-form-item label="变电站">
              <el-input
                v-model="searchForm.substation"
                placeholder="请输入变电站"
                clearable
              />
            </el-form-item>
            <el-form-item label="工作时间">
              <el-date-picker
                v-model="searchForm.workTimeRange"
                type="datetimerange"
                range-separator="至"
                start-placeholder="开始时间"
                end-placeholder="结束时间"
                value-format="YYYY-MM-DD HH:mm:ss"
                style="width: 240px"
              />
            </el-form-item>
            <el-form-item label="状态">
              <el-select
                v-model="searchForm.status"
                placeholder="请选择状态"
                clearable
              >
                <el-option label="全部" value="" />
                <el-option label="草稿" value="draft" />
                <el-option label="已提交" value="submitted" />
                <el-option label="已通过" value="approved" />
                <el-option label="已驳回" value="rejected" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSearch" :icon="Search">
                查询
              </el-button>
              <el-button @click="resetSearch" :icon="Refresh">重置</el-button>
              <el-button type="success" @click="exportReports" :icon="Download">
                导出
              </el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 批量操作 -->
        <div class="batch-operation" v-if="selectedReports.length > 0">
          <el-space>
            <span>已选择 {{ selectedReports.length }} 项</span>
            <el-button type="danger" @click="batchDelete"> 批量删除 </el-button>
            <el-button
              type="success"
              @click="batchSubmit"
              v-if="canBatchSubmit"
            >
              批量提交
            </el-button>
          </el-space>
        </div>

        <!-- 报备单表格 -->
        <el-table
          :data="reportData"
          style="width: 100%"
          v-loading="loading"
          @selection-change="handleSelectionChange"
          @sort-change="handleSortChange"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column
            prop="report_no"
            label="报备单号"
            width="180"
            sortable
          >
            <template #default="{ row }">
              <el-link type="primary" @click="viewReportDetail(row.id)">
                {{ row.report_no }}
              </el-link>
            </template>
          </el-table-column>
          <el-table-column prop="unit" label="单位" width="120" />
          <el-table-column label="变电站" width="150">
            <template #default="{ row }">
              {{ row.substations ? row.substations.join(", ") : "" }}
            </template>
          </el-table-column>
          <el-table-column label="工作时间" width="200">
            <template #default="{ row }">
              {{ formatDateTime(row.work_time_start) }}<br />
              {{ formatDateTime(row.work_time_end) }}
            </template>
          </el-table-column>
          <el-table-column
            prop="work_content"
            label="工作内容"
            min-width="150"
            show-overflow-tooltip
          />
          <el-table-column label="涉及设备" width="120">
            <template #default="{ row }">
              <el-tag type="info" size="small">
                {{ row.related_devices ? row.related_devices.length : 0 }}台
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100" sortable>
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)" size="small">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column
            label="审核意见"
            min-width="120"
            show-overflow-tooltip
          >
            <template #default="{ row }">
              {{ row.review_comments || "暂无" }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="250" fixed="right">
            <template #default="{ row }">
              <el-button
                type="primary"
                size="small"
                @click="handleEdit(row)"
                v-if="row.status === 'draft' || row.status === 'rejected'"
              >
                编辑
              </el-button>
              <el-button
                type="danger"
                size="small"
                @click="handleDelete(row)"
                v-if="row.status === 'draft'"
              >
                删除
              </el-button>
              <el-button
                type="success"
                size="small"
                @click="handleSubmit(row)"
                v-if="row.status === 'draft'"
              >
                提交上报
              </el-button>
              <el-button
                type="info"
                size="small"
                @click="viewReportDetail(row.id)"
              >
                详情
              </el-button>
              <el-button
                type="warning"
                size="small"
                @click="viewReportLogs(row.id)"
              >
                日志
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
      </el-card>
    </el-main>

    <!-- 报备单详情对话框 -->
    <el-dialog
      v-model="dialogVisible.detail"
      :title="`报备单详情 - ${selectedReport?.report_no || ''}`"
      width="800px"
      :close-on-click-modal="false"
    >
      <ReportDetail v-if="selectedReport" :report="selectedReport" />
      <template #footer>
        <el-button @click="dialogVisible.detail = false">关闭</el-button>
        <el-button
          type="primary"
          v-if="
            selectedReport?.status === 'draft' ||
            selectedReport?.status === 'rejected'
          "
          @click="handleEdit(selectedReport)"
        >
          编辑
        </el-button>
        <el-button
          type="success"
          v-if="selectedReport?.status === 'draft'"
          @click="handleSubmit(selectedReport)"
        >
          提交上报
        </el-button>
      </template>
    </el-dialog>

    <!-- 操作日志对话框 -->
    <el-dialog
      v-model="dialogVisible.logs"
      title="操作日志"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-timeline v-if="reportLogs.length > 0">
        <el-timeline-item
          v-for="log in reportLogs"
          :key="log.id"
          :timestamp="formatDateTime(log.created_at)"
          placement="top"
        >
          <el-card shadow="hover">
            <p>{{ log.operator_name }} {{ getActionText(log.action) }}</p>
            <p v-if="log.comment" style="color: #666; margin-top: 5px">
              备注: {{ log.comment }}
            </p>
          </el-card>
        </el-timeline-item>
      </el-timeline>
      <el-empty v-else description="暂无操作日志" />
      <template #footer>
        <el-button @click="dialogVisible.logs = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, reactive, onMounted, computed } from "vue";
import { useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import { Search, Refresh, Download } from "@element-plus/icons-vue";
import TheHeader from "../../components/TheHeader.vue";
import ReportDetail from "./ReportDetail.vue";
import { reportAPI } from "@/utils/request.js";

export default {
  name: "ManageReports",
  components: {
    TheHeader,
    ReportDetail,
  },
  setup() {
    const router = useRouter();
    const reportData = ref([]);
    const loading = ref(false);
    const selectedReport = ref(null);
    const selectedReports = ref([]);
    const reportLogs = ref([]);

    const searchForm = reactive({
      reportNo: "",
      substation: "",
      workTimeRange: [],
      status: "",
    });

    const pagination = reactive({
      current: 1,
      size: 10,
      total: 0,
    });

    const sortParams = reactive({
      field: "",
      order: "",
    });

    const dialogVisible = reactive({
      detail: false,
      logs: false,
    });

    // 计算属性：是否可以批量提交
    const canBatchSubmit = computed(() => {
      return selectedReports.value.every((report) => report.status === "draft");
    });

    const getStatusText = (status) => {
      const statusMap = {
        draft: "草稿",
        submitted: "已提交",
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

    const getActionText = (action) => {
      const actionMap = {
        create: "创建了报备单",
        update: "修改了报备单",
        submit: "提交了报备单",
        approve: "通过了报备单",
        reject: "驳回了报备单",
        delete: "删除了报备单",
      };
      return actionMap[action] || action;
    };

    const formatDateTime = (dateTime) => {
      if (!dateTime) return "";
      return new Date(dateTime).toLocaleString("zh-CN", { hour12: false });
    };

    const loadReports = async () => {
      loading.value = true;
      try {
        const params = {
          page: pagination.current,
          page_size: pagination.size,
          report_no: searchForm.reportNo,
          substation: searchForm.substation,
          status: searchForm.status,
          ...sortParams,
        };

        if (searchForm.workTimeRange && searchForm.workTimeRange.length === 2) {
          params.start_date = searchForm.workTimeRange[0];
          params.end_date = searchForm.workTimeRange[1];
        }

        const response = await reportAPI.getUserReports(params);
        reportData.value = response.data || response;
        pagination.total = response.total || response.length;
      } catch (error) {
        console.error("加载报备单失败:", error);
        ElMessage.error("加载报备单失败");
      } finally {
        loading.value = false;
      }
    };

    const handleSearch = () => {
      pagination.current = 1;
      loadReports();
    };

    const resetSearch = () => {
      Object.assign(searchForm, {
        reportNo: "",
        substation: "",
        workTimeRange: [],
        status: "",
      });
      pagination.current = 1;
      loadReports();
    };

    const exportReports = async () => {
      try {
        const params = {
          report_no: searchForm.reportNo,
          substation: searchForm.substation,
          status: searchForm.status,
        };

        if (searchForm.workTimeRange && searchForm.workTimeRange.length === 2) {
          params.start_date = searchForm.workTimeRange[0];
          params.end_date = searchForm.workTimeRange[1];
        }

        const response = await reportAPI.exportReports(params);

        // 创建下载链接
        const blob = new Blob([response], {
          type: "application/vnd.ms-excel",
        });
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement("a");
        link.href = url;
        link.download = `我的报备单_${new Date().getTime()}.xlsx`;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        window.URL.revokeObjectURL(url);

        ElMessage.success("导出成功");
      } catch (error) {
        console.error("导出失败:", error);
        ElMessage.error("导出失败");
      }
    };

    const viewReportDetail = async (reportId) => {
      try {
        const response = await reportAPI.getReportById(reportId);
        selectedReport.value = response;
        dialogVisible.detail = true;
      } catch (error) {
        console.error("获取报备单详情失败:", error);
        ElMessage.error("获取报备单详情失败");
      }
    };

    const viewReportLogs = async (reportId) => {
      try {
        const response = await reportAPI.getReportLogs(reportId);
        reportLogs.value = response;
        dialogVisible.logs = true;
      } catch (error) {
        console.error("获取操作日志失败:", error);
        ElMessage.error("获取操作日志失败");
      }
    };

    const handleEdit = (row) => {
      router.push({ path: "/user/create-report", query: { id: row.id } });
    };

    const handleDelete = async (row) => {
      try {
        await ElMessageBox.confirm(
          "确定要删除此报备单吗？删除后无法恢复。",
          "确认删除",
          {
            confirmButtonText: "确定",
            cancelButtonText: "取消",
            type: "warning",
          }
        );

        await reportAPI.deleteReport(row.id);
        ElMessage.success("删除成功");
        loadReports();
      } catch (error) {
        if (error === "cancel") return;
        console.error("删除失败:", error);
        ElMessage.error("删除失败");
      }
    };

    const handleSubmit = async (row) => {
      try {
        await ElMessageBox.confirm(
          "确定要提交此报备单吗？提交后需要管理员审核。",
          "确认提交",
          {
            confirmButtonText: "确定",
            cancelButtonText: "取消",
            type: "warning",
          }
        );

        await reportAPI.submitReport(row.id);
        ElMessage.success("提交成功，等待管理员审核");
        loadReports();
      } catch (error) {
        if (error === "cancel") return;
        console.error("提交失败:", error);
        ElMessage.error("提交失败");
      }
    };

    const batchDelete = async () => {
      if (selectedReports.value.length === 0) {
        ElMessage.warning("请选择要删除的报备单");
        return;
      }

      try {
        await ElMessageBox.confirm(
          `确定要删除选中的 ${selectedReports.value.length} 个报备单吗？删除后无法恢复。`,
          "确认批量删除",
          {
            confirmButtonText: "确定",
            cancelButtonText: "取消",
            type: "warning",
          }
        );

        const reportIds = selectedReports.value.map((report) => report.id);
        await reportAPI.batchDeleteReports(reportIds);
        ElMessage.success("批量删除成功");
        selectedReports.value = [];
        loadReports();
      } catch (error) {
        if (error === "cancel") return;
        console.error("批量删除失败:", error);
        ElMessage.error("批量删除失败");
      }
    };

    const batchSubmit = async () => {
      if (selectedReports.value.length === 0) {
        ElMessage.warning("请选择要提交的报备单");
        return;
      }

      try {
        await ElMessageBox.confirm(
          `确定要批量提交选中的 ${selectedReports.value.length} 个报备单吗？`,
          "确认批量提交",
          {
            confirmButtonText: "确定",
            cancelButtonText: "取消",
            type: "warning",
          }
        );

        // 批量提交
        const submitPromises = selectedReports.value.map((report) =>
          reportAPI.submitReport(report.id)
        );
        await Promise.all(submitPromises);
        ElMessage.success("批量提交成功");
        selectedReports.value = [];
        loadReports();
      } catch (error) {
        console.error("批量提交失败:", error);
        ElMessage.error("批量提交失败");
      }
    };

    const handleSelectionChange = (selection) => {
      selectedReports.value = selection;
    };

    const handleSortChange = ({ prop, order }) => {
      sortParams.field = prop;
      sortParams.order = order === "ascending" ? "asc" : "desc";
      loadReports();
    };

    const handleSizeChange = (size) => {
      pagination.size = size;
      pagination.current = 1;
      loadReports();
    };

    const handleCurrentChange = (page) => {
      pagination.current = page;
      loadReports();
    };

    const goBack = () => {
      router.back();
    };

    onMounted(() => {
      loadReports();
    });

    return {
      reportData,
      loading,
      selectedReport,
      selectedReports,
      reportLogs,
      searchForm,
      pagination,
      dialogVisible,
      canBatchSubmit,
      getStatusText,
      getStatusType,
      getActionText,
      formatDateTime,
      handleSearch,
      resetSearch,
      exportReports,
      viewReportDetail,
      viewReportLogs,
      handleEdit,
      handleDelete,
      handleSubmit,
      batchDelete,
      batchSubmit,
      handleSelectionChange,
      handleSortChange,
      handleSizeChange,
      handleCurrentChange,
      goBack,
      Search,
      Refresh,
      Download,
    };
  },
};
</script>

<style scoped>
.manage-reports-container {
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

.pagination-container {
  margin-top: 20px;
  text-align: right;
}
</style>