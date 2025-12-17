<template>
  <div class="report-detail">
    <!-- 基本信息 -->
    <el-descriptions :column="2" border size="small" class="mb-20">
      <el-descriptions-item label="报备单号">
        <el-tag type="primary">{{ report.report_no }}</el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="状态">
        <el-tag :type="getStatusType(report.status)" size="small">
          {{ getStatusText(report.status) }}
        </el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="报备人">
        {{ report.reporter_name }}
      </el-descriptions-item>
      <el-descriptions-item label="报备时间">
        {{ formatDateTime(report.created_at) }}
      </el-descriptions-item>
      <el-descriptions-item label="单位">{{
        report.unit
      }}</el-descriptions-item>
      <el-descriptions-item label="变电站">
        {{ report.substations ? report.substations.join(", ") : "" }}
      </el-descriptions-item>
      <el-descriptions-item label="工作时间" :span="2">
        {{ formatDateTime(report.work_time_start) }} 至
        {{ formatDateTime(report.work_time_end) }}
      </el-descriptions-item>
      <el-descriptions-item label="工作负责人">
        {{ report.responsible_person }}
      </el-descriptions-item>
      <el-descriptions-item label="工作内容">
        {{ report.work_content }}
      </el-descriptions-item>
    </el-descriptions>

    <!-- 涉及设备 -->
    <div class="section-title">涉及设备</div>
    <el-table
      :data="report.related_devices"
      border
      size="small"
      class="mb-20"
      v-if="report.related_devices && report.related_devices.length > 0"
    >
      <el-table-column prop="name" label="设备名称" width="150" />
      <el-table-column prop="ip" label="IP地址" width="120" />
      <el-table-column prop="type" label="设备类型" width="100" />
    </el-table>
    <el-empty v-else description="无涉及设备" :image-size="60" />

    <!-- 使用工具 -->
    <div class="section-title">使用工具</div>
    <el-table
      :data="report.tools"
      border
      size="small"
      class="mb-20"
      v-if="report.tools && report.tools.length > 0"
    >
      <el-table-column prop="name" label="工具名称" width="150" />
      <el-table-column prop="mac" label="MAC地址" width="150" />
      <el-table-column prop="serial_number" label="序列号" width="150" />
      <el-table-column prop="type" label="类型" width="100" />
    </el-table>
    <el-empty v-else description="无使用工具" :image-size="60" />

    <!-- 操作内容 -->
    <div class="section-title">操作内容</div>
    <div class="content-box mb-20">{{ report.operation }}</div>

    <!-- 安全措施 -->
    <div class="section-title">安全措施</div>
    <div class="content-box mb-20">{{ report.safety_measures }}</div>

    <!-- 备注 -->
    <div v-if="report.remarks" class="mb-20">
      <div class="section-title">备注</div>
      <div class="content-box">{{ report.remarks }}</div>
    </div>

    <!-- 审核信息 -->
    <div v-if="report.status === 'approved' || report.status === 'rejected'">
      <div class="section-title">审核信息</div>
      <el-descriptions :column="2" border size="small">
        <el-descriptions-item label="审核人">
          {{ report.reviewer_name }}
        </el-descriptions-item>
        <el-descriptions-item label="审核时间">
          {{ formatDateTime(report.reviewed_at) }}
        </el-descriptions-item>
        <el-descriptions-item label="审核意见" :span="2">
          <div class="review-comments">{{ report.review_comments }}</div>
        </el-descriptions-item>
      </el-descriptions>
    </div>
  </div>
</template>

<script>
import { defineProps } from "vue";

export default {
  name: "ReportDetail",
  props: {
    report: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
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

    return {
      getStatusText,
      getStatusType,
      formatDateTime,
    };
  },
};
</script>

<style scoped>
.report-detail {
  padding: 10px;
}

.mb-20 {
  margin-bottom: 20px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #409eff;
  margin: 20px 0 10px 0;
  padding-bottom: 5px;
  border-bottom: 1px solid #eee;
}

.content-box {
  background-color: #f8f9fa;
  padding: 15px;
  border-radius: 4px;
  border-left: 4px solid #409eff;
  line-height: 1.6;
  white-space: pre-wrap;
}

.review-comments {
  color: #e6a23c;
  font-weight: 500;
}
</style>