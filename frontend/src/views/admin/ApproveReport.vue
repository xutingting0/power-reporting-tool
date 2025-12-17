<template>
  <div class="approve-report-container">
    <TheHeader />
    <el-main>
      <el-page-header>
        <template #content>
          <span>审核报备单</span>
          <el-tag type="warning" size="small" style="margin-left: 10px">
            待审核: {{ pendingReports.length }} 个
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
            <el-form-item label="报备人">
              <el-input
                v-model="searchForm.reporter"
                placeholder="请输入报备人"
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
            <el-form-item>
              <el-button type="primary" @click="handleSearch" :icon="Search">
                查询
              </el-button>
              <el-button @click="resetSearch" :icon="Refresh">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 待审核报备单表格 -->
        <el-table
          :data="filteredReports"
          style="width: 100%"
          v-loading="loading"
          @sort-change="handleSortChange"
        >
          <el-table-column
            prop="report_no"
            label="报备单号"
            width="180"
            sortable
          >
            <template #default="{ row }">
              <el-link type="primary" @click="viewReportDetail(row)">
                {{ row.report_no }}
              </el-link>
            </template>
          </el-table-column>
          <el-table-column prop="reporter_name" label="报备人" width="100" />
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
          <el-table-column
            prop="responsible_person"
            label="负责人"
            width="100"
          />
          <el-table-column label="涉及设备" width="120">
            <template #default="{ row }">
              <el-tag type="info" size="small">
                {{ row.related_devices ? row.related_devices.length : 0 }}台
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column
            label="创建时间"
            width="150"
            sortable="custom"
            prop="created_at"
          >
            <template #default="{ row }">
              {{ formatDateTime(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <el-button
                type="primary"
                size="small"
                @click="handleApprove(row)"
              >
                通过
              </el-button>
              <el-button type="danger" size="small" @click="handleReject(row)">
                驳回
              </el-button>
              <el-button
                type="info"
                size="small"
                @click="viewReportDetail(row)"
              >
                详情
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

    <!-- 审核对话框 -->
    <el-dialog
      v-model="dialogVisible.approve"
      title="审核通过"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form :model="approveForm" label-width="80px">
        <el-form-item label="审核意见">
          <el-input
            v-model="approveForm.comments"
            type="textarea"
            :rows="3"
            placeholder="请输入审核意见（可选）"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible.approve = false">取消</el-button>
        <el-button type="primary" @click="confirmApprove">确定通过</el-button>
      </template>
    </el-dialog>

    <!-- 驳回对话框 -->
    <el-dialog
      v-model="dialogVisible.reject"
      title="驳回报备单"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form :model="rejectForm" label-width="80px">
        <el-form-item label="驳回原因" required>
          <el-input
            v-model="rejectForm.comments"
            type="textarea"
            :rows="3"
            placeholder="请输入驳回原因"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible.reject = false">取消</el-button>
        <el-button type="danger" @click="confirmReject">确定驳回</el-button>
      </template>
    </el-dialog>

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
          v-if="selectedReport"
          @click="handleApprove(selectedReport)"
        >
          通过
        </el-button>
        <el-button
          type="danger"
          v-if="selectedReport"
          @click="handleReject(selectedReport)"
        >
          驳回
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, reactive, onMounted, computed } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { Search, Refresh } from "@element-plus/icons-vue";
import TheHeader from "../../components/TheHeader.vue";
import ReportDetail from "../user/ReportDetail.vue";
import { reportAPI } from "@/utils/request.js";

export default {
  name: "ApproveReport",
  components: {
    TheHeader,
    ReportDetail,
  },
  setup() {
    const pendingReports = ref([]);
    const loading = ref(false);
    const selectedReport = ref(null);

    const searchForm = reactive({
      reportNo: "",
      reporter: "",
      substation: "",
      workTimeRange: [],
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
      approve: false,
      reject: false,
      detail: false,
    });

    const approveForm = reactive({
      comments: "",
    });

    const rejectForm = reactive({
      comments: "",
    });

    const filteredReports = computed(() => {
      let reports = pendingReports.value;

      // 按报备单号筛选
      if (searchForm.reportNo) {
        reports = reports.filter((report) =>
          report.report_no.includes(searchForm.reportNo)
        );
      }

      // 按报备人筛选
      if (searchForm.reporter) {
        reports = reports.filter((report) =>
          report.reporter_name.includes(searchForm.reporter)
        );
      }

      // 按变电站筛选
      if (searchForm.substation) {
        reports = reports.filter(
          (report) =>
            report.substations &&
            report.substations.some((sub) =>
              sub.includes(searchForm.substation)
            )
        );
      }

      // 按工作时间筛选
      if (searchForm.workTimeRange && searchForm.workTimeRange.length === 2) {
        const startTime = new Date(searchForm.workTimeRange[0]).getTime();
        const endTime = new Date(searchForm.workTimeRange[1]).getTime();
        reports = reports.filter((report) => {
          const reportStart = new Date(report.work_time_start).getTime();
          const reportEnd = new Date(report.work_time_end).getTime();
          return reportStart >= startTime && reportEnd <= endTime;
        });
      }

      // 分页
      const startIndex = (pagination.current - 1) * pagination.size;
      const endIndex = startIndex + pagination.size;
      pagination.total = reports.length;
      return reports.slice(startIndex, endIndex);
    });

    const formatDateTime = (dateTime) => {
      if (!dateTime) return "";
      return new Date(dateTime).toLocaleString("zh-CN", { hour12: false });
    };

    const loadPendingReports = async () => {
      loading.value = true;
      try {
        const response = await reportAPI.getPendingReports();
        pendingReports.value = response;
      } catch (error) {
        console.error("加载待审核报备单失败:", error);
        ElMessage.error("加载待审核报备单失败");
      } finally {
        loading.value = false;
      }
    };

    const handleSearch = () => {
      pagination.current = 1;
    };

    const resetSearch = () => {
      Object.assign(searchForm, {
        reportNo: "",
        reporter: "",
        substation: "",
        workTimeRange: [],
      });
      pagination.current = 1;
    };

    const viewReportDetail = (report) => {
      selectedReport.value = report;
      dialogVisible.detail = true;
    };

    const handleApprove = (report) => {
      selectedReport.value = report;
      approveForm.comments = "";
      dialogVisible.approve = true;
    };

    const confirmApprove = async () => {
      if (!selectedReport.value) return;

      try {
        await reportAPI.approveReport(selectedReport.value.id, {
          review_comments: approveForm.comments,
        });

        ElMessage.success("报备单审核通过");
        dialogVisible.approve = false;
        dialogVisible.detail = false;
        loadPendingReports();
      } catch (error) {
        console.error("审核通过失败:", error);
        ElMessage.error("审核通过失败");
      }
    };

    const handleReject = (report) => {
      selectedReport.value = report;
      rejectForm.comments = "";
      dialogVisible.reject = true;
    };

    const confirmReject = async () => {
      if (!rejectForm.comments) {
        ElMessage.warning("请输入驳回原因");
        return;
      }

      try {
        await reportAPI.rejectReport(selectedReport.value.id, {
          review_comments: rejectForm.comments,
        });

        ElMessage.success("报备单已驳回");
        dialogVisible.reject = false;
        dialogVisible.detail = false;
        loadPendingReports();
      } catch (error) {
        console.error("驳回失败:", error);
        ElMessage.error("驳回失败");
      }
    };

    const handleSortChange = ({ prop, order }) => {
      sortParams.field = prop;
      sortParams.order = order === "ascending" ? "asc" : "desc";
      // 这里可以调用排序API或本地排序
    };

    const handleSizeChange = (size) => {
      pagination.size = size;
      pagination.current = 1;
    };

    const handleCurrentChange = (page) => {
      pagination.current = page;
    };

    onMounted(() => {
      loadPendingReports();
    });

    return {
      pendingReports,
      filteredReports,
      loading,
      selectedReport,
      searchForm,
      pagination,
      dialogVisible,
      approveForm,
      rejectForm,
      formatDateTime,
      handleSearch,
      resetSearch,
      viewReportDetail,
      handleApprove,
      confirmApprove,
      handleReject,
      confirmReject,
      handleSortChange,
      handleSizeChange,
      handleCurrentChange,
      Search,
      Refresh,
    };
  },
};
</script>

<style scoped>
.approve-report-container {
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

.pagination-container {
  margin-top: 20px;
  text-align: right;
}
</style>