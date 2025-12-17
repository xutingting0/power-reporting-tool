<template>
  <div class="create-report-container">
    <TheHeader />
    <el-main>
      <el-page-header @back="goBack">
        <template #content>
          <span class="page-title">编制报备单</span>
          <el-tag
            v-if="reportForm.status"
            :type="statusType"
            size="small"
            style="margin-left: 10px"
          >
            {{ statusText }}
          </el-tag>
        </template>
        <template #extra>
          <el-button @click="showTemplateList" type="primary" plain
            >从模板创建</el-button
          >
        </template>
      </el-page-header>

      <el-card class="form-card">
        <template #header>
          <div class="card-header">
            <span>报备信息填写</span>
            <div class="header-actions">
              <el-button type="info" plain @click="showPreview"
                >预览报备单</el-button
              >
              <el-button @click="saveAsTemplate" plain>保存为模板</el-button>
            </div>
          </div>
        </template>

        <el-form
          ref="reportFormRef"
          :model="reportForm"
          :rules="formRules"
          label-width="120px"
          label-position="top"
          status-icon
        >
          <el-row :gutter="24">
            <el-col :span="8">
              <el-form-item label="单位" prop="unit" required>
                <el-select
                  v-model="reportForm.unit"
                  placeholder="请选择单位"
                  clearable
                  filterable
                >
                  <el-option
                    v-for="unit in unitOptions"
                    :key="unit.value"
                    :label="unit.label"
                    :value="unit.value"
                  >
                    <span>{{ unit.label }}</span>
                    <el-tag
                      size="small"
                      type="info"
                      style="margin-left: 10px"
                      >{{ unit.code }}</el-tag
                    >
                  </el-option>
                </el-select>
              </el-form-item>
            </el-col>

            <el-col :span="8">
              <el-form-item label="变电站" prop="substations" required>
                <el-select
                  v-model="reportForm.substations"
                  placeholder="请选择变电站"
                  multiple
                  collapse-tags
                  collapse-tags-tooltip
                  filterable
                  clearable
                >
                  <el-option-group
                    v-for="group in substationGroups"
                    :key="group.label"
                    :label="group.label"
                  >
                    <el-option
                      v-for="substation in group.options"
                      :key="substation.value"
                      :label="substation.label"
                      :value="substation.value"
                    />
                  </el-option-group>
                </el-select>
              </el-form-item>
            </el-col>

            <el-col :span="8">
              <el-form-item label="报备单号" prop="reportNo">
                <el-input
                  v-model="reportForm.reportNo"
                  placeholder="系统自动生成"
                  readonly
                  class="readonly-input"
                >
                  <template #append>
                    <el-button
                      :icon="Refresh"
                      @click="generateReportNo"
                      title="重新生成单号"
                    />
                  </template>
                </el-input>
              </el-form-item>
            </el-col>
          </el-row>

          <el-row :gutter="24">
            <el-col :span="12">
              <el-form-item label="工作时间" prop="workTimeRange" required>
                <el-date-picker
                  v-model="reportForm.workTimeRange"
                  type="datetimerange"
                  range-separator="至"
                  start-placeholder="开始时间"
                  end-placeholder="结束时间"
                  value-format="YYYY-MM-DD HH:mm:ss"
                  :default-time="defaultTimes"
                  style="width: 100%"
                />
              </el-form-item>
            </el-col>

            <el-col :span="12">
              <el-form-item
                label="工作负责人"
                prop="responsiblePerson"
                required
              >
                <el-input
                  v-model="reportForm.responsiblePerson"
                  placeholder="请输入工作负责人"
                  clearable
                >
                  <template #prepend>
                    <el-icon><User /></el-icon>
                  </template>
                </el-input>
              </el-form-item>
            </el-col>
          </el-row>

          <el-form-item label="工作内容" prop="workContent" required>
            <div class="work-content-container">
              <el-select
                v-model="reportForm.workContent"
                placeholder="请选择工作内容"
                clearable
                filterable
                allow-create
                @change="onWorkContentChange"
                style="width: 100%"
              >
                <el-option-group
                  v-for="group in workContentGroups"
                  :key="group.label"
                  :label="group.label"
                >
                  <el-option
                    v-for="content in group.options"
                    :key="content.value"
                    :label="content.label"
                    :value="content.value"
                  >
                    <span>{{ content.label }}</span>
                    <el-tag
                      v-if="content.riskLevel"
                      size="small"
                      :type="getRiskType(content.riskLevel)"
                      style="margin-left: 10px"
                    >
                      {{ content.riskLevel }}
                    </el-tag>
                  </el-option>
                </el-option-group>
              </el-select>

              <div
                v-if="reportForm.workContent"
                class="work-content-desc"
                style="margin-top: 8px"
              >
                <el-text type="info" size="small">
                  {{ getWorkContentDesc(reportForm.workContent) }}
                </el-text>
              </div>
            </div>
          </el-form-item>

          <el-form-item label="涉及设备" prop="relatedDevices" required>
            <div class="device-container">
              <el-tooltip
                content="根据选择的工作内容自动关联设备，也可手动添加"
              >
                <el-button
                  type="primary"
                  @click="addDevice"
                  size="small"
                  :icon="Plus"
                >
                  添加设备
                </el-button>
              </el-tooltip>
              <el-button @click="clearDevices" size="small" :icon="Delete">
                清空
              </el-button>

              <el-table
                :data="reportForm.relatedDevices"
                border
                size="small"
                class="device-table"
                style="margin-top: 10px"
              >
                <el-table-column label="设备名称" prop="name" min-width="180">
                  <template #default="{ row, $index }">
                    <el-select
                      v-model="row.name"
                      placeholder="选择设备"
                      filterable
                      allow-create
                      style="width: 100%"
                      @change="onDeviceChange($index, row.name)"
                    >
                      <el-option
                        v-for="device in deviceOptions"
                        :key="device.value"
                        :label="device.label"
                        :value="device.value"
                      />
                    </el-select>
                  </template>
                </el-table-column>
                <el-table-column label="IP地址" prop="ip" min-width="150">
                  <template #default="{ row, $index }">
                    <el-input
                      v-model="row.ip"
                      placeholder="请输入IP地址"
                      @blur="validateIP(row.ip, $index)"
                    />
                  </template>
                </el-table-column>
                <el-table-column label="设备类型" prop="type" min-width="120">
                  <template #default="{ row }">
                    <el-tag size="small">{{ getDeviceType(row.name) }}</el-tag>
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="80">
                  <template #default="{ $index }">
                    <el-button
                      type="danger"
                      :icon="Delete"
                      size="small"
                      circle
                      @click="removeDevice($index)"
                    />
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-form-item>

          <el-form-item label="使用工具" prop="tools">
            <div class="tool-container">
              <el-tooltip content="添加使用的工具和软件">
                <el-button
                  type="primary"
                  @click="addTool"
                  size="small"
                  :icon="Plus"
                >
                  添加工具
                </el-button>
              </el-tooltip>

              <el-table
                :data="reportForm.tools"
                border
                size="small"
                class="tool-table"
                style="margin-top: 10px"
              >
                <el-table-column label="工具名称" prop="name" min-width="150">
                  <template #default="{ row, $index }">
                    <el-select
                      v-model="row.name"
                      placeholder="选择工具"
                      filterable
                      allow-create
                      style="width: 100%"
                      @change="onToolChange($index, row.name)"
                    >
                      <el-option
                        v-for="tool in toolOptions"
                        :key="tool.value"
                        :label="tool.label"
                        :value="tool.value"
                      />
                    </el-select>
                  </template>
                </el-table-column>
                <el-table-column label="MAC地址" prop="mac" min-width="150">
                  <template #default="{ row, $index }">
                    <el-input
                      v-model="row.mac"
                      placeholder="请输入MAC地址"
                      @blur="validateMAC(row.mac, $index)"
                    />
                  </template>
                </el-table-column>
                <el-table-column
                  label="序列号"
                  prop="serialNumber"
                  min-width="150"
                >
                  <template #default="{ row }">
                    <el-input
                      v-model="row.serialNumber"
                      placeholder="请输入序列号"
                    />
                  </template>
                </el-table-column>
                <el-table-column label="类型" prop="type" min-width="100">
                  <template #default="{ row }">
                    <el-tag size="small">{{ row.type || "其他" }}</el-tag>
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="80">
                  <template #default="{ $index }">
                    <el-button
                      type="danger"
                      :icon="Delete"
                      size="small"
                      circle
                      @click="removeTool($index)"
                    />
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-form-item>

          <el-form-item label="操作内容" prop="operation" required>
            <el-input
              v-model="reportForm.operation"
              type="textarea"
              :rows="3"
              placeholder="请输入详细的操作步骤和内容"
              maxlength="500"
              show-word-limit
            />
          </el-form-item>

          <el-form-item label="安全措施" prop="safetyMeasures" required>
            <el-input
              v-model="reportForm.safetyMeasures"
              type="textarea"
              :rows="3"
              placeholder="请输入安全措施，如：已完成病毒扫描、已断开非必要连接等"
              maxlength="500"
              show-word-limit
            />
          </el-form-item>

          <el-form-item label="备注" prop="remarks">
            <el-input
              v-model="reportForm.remarks"
              type="textarea"
              :rows="2"
              placeholder="请输入其他需要说明的事项"
              maxlength="200"
              show-word-limit
            />
          </el-form-item>

          <el-form-item>
            <div class="form-actions">
              <el-button
                type="primary"
                @click="saveReport"
                :loading="loading.save"
                :icon="DocumentAdd"
              >
                保存草稿
              </el-button>
              <el-button
                type="success"
                @click="submitReport"
                :loading="loading.submit"
                :disabled="!isFormValid"
                :icon="Promotion"
              >
                发送上报
              </el-button>
              <el-button type="warning" @click="resetForm" :icon="Refresh">
                重置表单
              </el-button>
              <el-button @click="goBack" :icon="ArrowLeft">
                返回列表
              </el-button>
            </div>
          </el-form-item>
        </el-form>
      </el-card>
    </el-main>

    <!-- 预览对话框 -->
    <el-dialog
      v-model="dialogVisible.preview"
      title="报备单预览"
      width="800px"
      :show-close="false"
    >
      <div class="preview-content">
        <div class="preview-header">
          <h3>电力监控系统作业报备单</h3>
          <div class="preview-info">
            <p><strong>报备单号：</strong>{{ previewData.reportNo }}</p>
            <p><strong>报备时间：</strong>{{ previewData.createTime }}</p>
            <p><strong>报备人：</strong>{{ previewData.reporter }}</p>
          </div>
        </div>

        <div class="preview-body">
          <el-descriptions :column="2" border size="small">
            <el-descriptions-item label="单位">{{
              previewData.unit
            }}</el-descriptions-item>
            <el-descriptions-item label="变电站">{{
              previewData.substations.join(", ")
            }}</el-descriptions-item>
            <el-descriptions-item label="工作时间" :span="2">
              {{ previewData.workTimeRange[0] }} 至
              {{ previewData.workTimeRange[1] }}
            </el-descriptions-item>
            <el-descriptions-item label="工作负责人">{{
              previewData.responsiblePerson
            }}</el-descriptions-item>
            <el-descriptions-item label="工作内容">{{
              previewData.workContent
            }}</el-descriptions-item>
          </el-descriptions>

          <h4>涉及设备：</h4>
          <el-table :data="previewData.relatedDevices" size="small" border>
            <el-table-column prop="name" label="设备名称" />
            <el-table-column prop="ip" label="IP地址" />
            <el-table-column prop="type" label="设备类型" />
          </el-table>

          <h4>使用工具：</h4>
          <el-table :data="previewData.tools" size="small" border>
            <el-table-column prop="name" label="工具名称" />
            <el-table-column prop="mac" label="MAC地址" />
            <el-table-column prop="serialNumber" label="序列号" />
            <el-table-column prop="type" label="类型" />
          </el-table>

          <h4>操作内容：</h4>
          <div class="preview-text">{{ previewData.operation }}</div>

          <h4>安全措施：</h4>
          <div class="preview-text">{{ previewData.safetyMeasures }}</div>

          <div v-if="previewData.remarks" class="preview-remarks">
            <h4>备注：</h4>
            <div class="preview-text">{{ previewData.remarks }}</div>
          </div>
        </div>
      </div>

      <template #footer>
        <el-button @click="dialogVisible.preview = false">关闭</el-button>
        <el-button type="primary" @click="handlePrint">打印</el-button>
      </template>
    </el-dialog>

    <!-- 模板选择对话框 -->
    <el-dialog
      v-model="dialogVisible.template"
      title="选择报备模板"
      width="600px"
    >
      <el-table
        :data="templateList"
        @row-click="selectTemplate"
        highlight-current-row
      >
        <el-table-column prop="name" label="模板名称" />
        <el-table-column prop="workContent" label="工作内容" />
        <el-table-column prop="devicesCount" label="设备数量" width="80" />
        <el-table-column label="创建时间" width="120">
          <template #default="{ row }">
            {{ formatDate(row.createTime) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="80">
          <template #default="{ row }">
            <el-button type="text" @click.stop="deleteTemplate(row.id)"
              >删除</el-button
            >
          </template>
        </el-table-column>
      </el-table>

      <template #footer>
        <el-button @click="dialogVisible.template = false">取消</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, reactive, computed, onMounted, nextTick, watch } from "vue";
import { useRouter, useRoute } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  Refresh,
  User,
  Plus,
  Delete,
  DocumentAdd,
  Promotion,
  ArrowLeft,
} from "@element-plus/icons-vue";
import TheHeader from "../../components/TheHeader.vue";
import * as reportAPI from "@/api/report";
import * as templateAPI from "@/api/template";
import { formatDateTime } from "@/utils/date";

export default {
  name: "CreateReport",
  components: {
    TheHeader,
  },
  setup() {
    const router = useRouter();
    const route = useRoute();
    const reportFormRef = ref();

    // 表单数据
    const reportForm = reactive({
      id: null,
      report_no: "",
      unit: "",
      substations: [],
      work_time_start: "",
      work_time_end: "",
      responsible_person: "",
      work_content: "",
      related_devices: [],
      tools: [],
      operation: "",
      safety_measures: "",
      remarks: "",
      status: "draft",
    });

    // 工作时间范围（仅用于显示）
    const workTimeRange = ref([]);

    // 加载状态
    const loading = reactive({
      save: false,
      submit: false,
      templates: false,
    });

    // 对话框状态
    const dialogVisible = reactive({
      preview: false,
      template: false,
    });

    // 预览数据
    const previewData = ref({});

    // 模板列表
    const templateList = ref([]);

    // 选项数据
    const unitOptions = [
      { value: "遵义供电局", label: "遵义供电局", code: "ZYGD" },
      { value: "贵阳供电局", label: "贵阳供电局", code: "GYGD" },
      { value: "六盘水供电局", label: "六盘水供电局", code: "LPSGD" },
      { value: "安顺供电局", label: "安顺供电局", code: "ASGD" },
      { value: "毕节供电局", label: "毕节供电局", code: "BJGD" },
    ];

    const substationGroups = [
      {
        label: "110kV变电站",
        options: [
          { value: "110kV道真变电站", label: "110kV道真变电站" },
          { value: "110kV正安变电站", label: "110kV正安变电站" },
          { value: "110kV务川变电站", label: "110kV务川变电站" },
          { value: "110kV凤冈变电站", label: "110kV凤冈变电站" },
        ],
      },
      {
        label: "220kV变电站",
        options: [
          { value: "220kV遵义变电站", label: "220kV遵义变电站" },
          { value: "220kV桐梓变电站", label: "220kV桐梓变电站" },
          { value: "220kV仁怀变电站", label: "220kV仁怀变电站" },
        ],
      },
    ];

    const workContentGroups = [
      {
        label: "安装调试",
        options: [
          {
            value: "测控装置安装调试",
            label: "测控装置安装调试",
            riskLevel: "中风险",
          },
          {
            value: "保护装置安装调试",
            label: "保护装置安装调试",
            riskLevel: "高风险",
          },
          {
            value: "自动化系统调试",
            label: "自动化系统调试",
            riskLevel: "中风险",
          },
        ],
      },
      {
        label: "维护检修",
        options: [
          { value: "设备定期巡检", label: "设备定期巡检", riskLevel: "低风险" },
          { value: "设备缺陷处理", label: "设备缺陷处理", riskLevel: "中风险" },
          {
            value: "设备预防性试验",
            label: "设备预防性试验",
            riskLevel: "高风险",
          },
        ],
      },
      {
        label: "软件操作",
        options: [
          { value: "系统参数修改", label: "系统参数修改", riskLevel: "中风险" },
          { value: "数据库维护", label: "数据库维护", riskLevel: "低风险" },
          { value: "软件升级", label: "软件升级", riskLevel: "高风险" },
        ],
      },
    ];

    const deviceOptions = [
      { value: "测控装置", label: "测控装置" },
      { value: "保护装置", label: "保护装置" },
      { value: "后台监控机", label: "后台监控机" },
      { value: "远动机", label: "远动机" },
      { value: "交换机", label: "交换机" },
      { value: "路由器", label: "路由器" },
      { value: "防火墙", label: "防火墙" },
      { value: "服务器", label: "服务器" },
    ];

    const toolOptions = [
      { value: "专用运维笔记本电脑", label: "专用运维笔记本电脑" },
      { value: "专用安全U盘", label: "专用安全U盘" },
      { value: "调试软件", label: "调试软件" },
      { value: "测试仪器", label: "测试仪器" },
      { value: "通讯设备", label: "通讯设备" },
    ];

    // 计算属性
    const statusText = computed(() => {
      const statusMap = {
        draft: "草稿",
        submitted: "已上报",
        approved: "已通过",
        rejected: "已驳回",
      };
      return statusMap[reportForm.status] || "草稿";
    });

    const statusType = computed(() => {
      const typeMap = {
        draft: "info",
        submitted: "warning",
        approved: "success",
        rejected: "danger",
      };
      return typeMap[reportForm.status] || "info";
    });

    const isFormValid = computed(() => {
      return (
        reportForm.unit &&
        reportForm.substations.length > 0 &&
        reportForm.work_time_start &&
        reportForm.work_time_end &&
        reportForm.responsible_person &&
        reportForm.work_content &&
        reportForm.related_devices.length > 0 &&
        reportForm.operation &&
        reportForm.safety_measures
      );
    });

    // 表单验证规则
    const formRules = {
      unit: [{ required: true, message: "请选择单位", trigger: "blur" }],
      substations: [
        { required: true, message: "请选择变电站", trigger: "blur" },
      ],
      responsible_person: [
        { required: true, message: "请输入工作负责人", trigger: "blur" },
      ],
      work_content: [
        { required: true, message: "请输入工作内容", trigger: "blur" },
      ],
      operation: [
        { required: true, message: "请输入操作内容", trigger: "blur" },
      ],
      safety_measures: [
        { required: true, message: "请输入安全措施", trigger: "blur" },
      ],
    };

    // 默认时间设置
    const defaultTimes = [
      new Date(2000, 1, 1, 9, 0, 0), // 开始时间默认9:00
      new Date(2000, 1, 1, 17, 0, 0), // 结束时间默认17:00
    ];

    // 方法
    const generateReportNo = () => {
      const date = new Date();
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, "0");
      const day = String(date.getDate()).padStart(2, "0");
      const random = String(Math.floor(Math.random() * 10000)).padStart(4, "0");
      reportForm.report_no = `BG${year}${month}${day}${random}`;
    };

    const getRiskType = (riskLevel) => {
      const riskMap = {
        低风险: "success",
        中风险: "warning",
        高风险: "danger",
      };
      return riskMap[riskLevel] || "info";
    };

    const getWorkContentDesc = (workContent) => {
      const descMap = {
        测控装置安装调试: "包括装置安装、参数配置、功能测试等",
        保护装置安装调试: "包括保护功能验证、定值设置等",
        设备定期巡检: "巡视检查设备运行状态，记录设备参数",
        系统参数修改: "修改系统配置参数，需谨慎操作",
      };
      return descMap[workContent] || "";
    };

    const onWorkContentChange = (value) => {
      // 根据工作内容自动关联设备
      if (value.includes("测控")) {
        reportForm.related_devices = [
          { name: "测控装置", ip: "", type: "测控设备" },
          { name: "后台监控机", ip: "", type: "计算机" },
        ];
      } else if (value.includes("保护")) {
        reportForm.related_devices = [
          { name: "保护装置", ip: "", type: "保护设备" },
          { name: "后台监控机", ip: "", type: "计算机" },
        ];
      }
    };

    const addDevice = () => {
      reportForm.related_devices.push({
        name: "",
        ip: "",
        type: "",
      });
    };

    const removeDevice = (index) => {
      reportForm.related_devices.splice(index, 1);
    };

    const clearDevices = () => {
      reportForm.related_devices = [];
    };

    const onDeviceChange = (index, name) => {
      const typeMap = {
        测控装置: "测控设备",
        保护装置: "保护设备",
        后台监控机: "计算机",
        远动机: "通讯设备",
        交换机: "网络设备",
        路由器: "网络设备",
        防火墙: "安全设备",
        服务器: "服务器",
      };
      reportForm.related_devices[index].type = typeMap[name] || "";
    };

    const validateIP = (ip, index) => {
      const ipPattern = /^(\d{1,3}\.){3}\d{1,3}$/;
      if (ip && !ipPattern.test(ip)) {
        ElMessage.warning("IP地址格式不正确");
        reportForm.related_devices[index].ip = "";
      }
    };

    const getDeviceType = (name) => {
      const typeMap = {
        测控装置: "测控设备",
        保护装置: "保护设备",
        后台监控机: "计算机",
        远动机: "通讯设备",
        交换机: "网络设备",
        路由器: "网络设备",
        防火墙: "安全设备",
        服务器: "服务器",
      };
      return typeMap[name] || "其他设备";
    };

    const addTool = () => {
      reportForm.tools.push({
        name: "",
        mac: "",
        serial_number: "",
        type: "",
      });
    };

    const removeTool = (index) => {
      reportForm.tools.splice(index, 1);
    };

    const onToolChange = (index, name) => {
      const typeMap = {
        专用运维笔记本电脑: "计算机",
        专用安全U盘: "存储设备",
        调试软件: "软件",
        测试仪器: "仪器仪表",
        通讯设备: "通讯设备",
      };
      reportForm.tools[index].type = typeMap[name] || "其他";
    };

    const validateMAC = (mac, index) => {
      const macPattern = /^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$/;
      if (mac && !macPattern.test(mac)) {
        ElMessage.warning("MAC地址格式不正确");
        reportForm.tools[index].mac = "";
      }
    };

    const saveReport = async () => {
      try {
        await reportFormRef.value.validate();

        loading.save = true;

        // 更新工作时间
        if (workTimeRange.value && workTimeRange.value.length === 2) {
          reportForm.work_time_start = workTimeRange.value[0];
          reportForm.work_time_end = workTimeRange.value[1];
        }

        // 验证安全措施是否包含必填项
        if (!validateSafetyMeasures(reportForm.safety_measures)) {
          ElMessage.error("安全措施必须包含'已确认所有工具完成病毒扫描'");
          return;
        }

        if (reportForm.id) {
          // 更新现有报备单
          await reportAPI.updateReport(reportForm.id, reportForm);
          ElMessage.success("报备单已更新");
        } else {
          // 创建新报备单
          const response = await reportAPI.createReport(reportForm);
          reportForm.id = response.report_id;
          reportForm.report_no = response.report_no;
          ElMessage.success("报备单已保存为草稿");
        }
      } catch (error) {
        console.error("保存失败:", error);
        ElMessage.error(error.response?.data?.error || "保存失败");
      } finally {
        loading.save = false;
      }
    };

    const submitReport = async () => {
      try {
        await ElMessageBox.confirm(
          "确定要上报该报备单吗？上报后需要管理员审核。",
          "确认上报",
          {
            confirmButtonText: "确定",
            cancelButtonText: "取消",
            type: "warning",
          }
        );

        await reportFormRef.value.validate();

        // 更新工作时间
        if (workTimeRange.value && workTimeRange.value.length === 2) {
          reportForm.work_time_start = workTimeRange.value[0];
          reportForm.work_time_end = workTimeRange.value[1];
        }

        // 验证安全措施是否包含必填项
        if (!validateSafetyMeasures(reportForm.safety_measures)) {
          ElMessage.error("安全措施必须包含'已确认所有工具完成病毒扫描'");
          return;
        }

        if (reportForm.status === "draft") {
          loading.submit = true;

          if (reportForm.id) {
            // 提交现有报备单
            await reportAPI.submitReport(reportForm.id);
          } else {
            // 先创建再提交
            reportForm.status = "submitted";
            await reportAPI.createReport(reportForm);
          }

          ElMessage.success("报备单已发送上报，等待管理员审核");
          router.push("/user/manage-reports");
        }
      } catch (error) {
        if (error === "cancel") {
          return;
        }
        console.error("上报失败:", error);
        ElMessage.error(error.response?.data?.error || "上报失败");
      } finally {
        loading.submit = false;
      }
    };

    // 安全措施验证函数
    const validateSafetyMeasures = (measures) => {
      const requiredPhrases = [
        "已确认所有工具完成病毒扫描",
        "病毒扫描完成",
        "完成病毒扫描",
      ];

      const measuresLower = measures.toLowerCase();
      for (const phrase of requiredPhrases) {
        if (measuresLower.includes(phrase.toLowerCase())) {
          return true;
        }
      }
      return false;
    };

    const resetForm = () => {
      ElMessageBox.confirm(
        "确定要重置表单吗？所有未保存的内容将丢失。",
        "确认重置",
        {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning",
        }
      ).then(() => {
        Object.assign(reportForm, {
          id: null,
          report_no: generateReportNo(),
          unit: "",
          substations: [],
          work_time_start: "",
          work_time_end: "",
          responsible_person: "",
          work_content: "",
          related_devices: [],
          tools: [],
          operation: "",
          safety_measures: "",
          remarks: "",
          status: "draft",
        });
        workTimeRange.value = [];
        ElMessage.success("表单已重置");
      });
    };

    const showPreview = () => {
      previewData.value = {
        report_no: reportForm.report_no || "系统自动生成",
        unit: reportForm.unit,
        substations: reportForm.substations,
        work_time_start: reportForm.work_time_start || workTimeRange.value[0],
        work_time_end: reportForm.work_time_end || workTimeRange.value[1],
        responsible_person: reportForm.responsible_person,
        work_content: reportForm.work_content,
        related_devices: reportForm.related_devices,
        tools: reportForm.tools,
        operation: reportForm.operation,
        safety_measures: reportForm.safety_measures,
        remarks: reportForm.remarks,
        created_at: new Date().toLocaleString(),
        reporter_name: localStorage.getItem("realName") || "当前用户",
      };
      dialogVisible.preview = true;
    };

    const handlePrint = () => {
      window.print();
    };

    const loadTemplates = async () => {
      try {
        loading.templates = true;
        const response = await templateAPI.getAllTemplates();
        templateList.value = response.data || response;
      } catch (error) {
        console.error("加载模板失败:", error);
        ElMessage.error("加载模板失败");
      } finally {
        loading.templates = false;
      }
    };

    const showTemplateList = () => {
      loadTemplates();
      dialogVisible.template = true;
    };

    const selectTemplate = async (template) => {
      try {
        // 使用模板创建报备单
        const response = await reportAPI.useTemplate(template.id);

        // 加载新创建的报备单
        await loadReport(response.report_id);

        ElMessage.success(`使用模板创建成功，报备单号：${response.report_no}`);
        dialogVisible.template = false;
      } catch (error) {
        console.error("使用模板失败:", error);
        ElMessage.error(error.response?.data?.error || "使用模板失败");
      }
    };

    const deleteTemplate = async (id) => {
      try {
        await ElMessageBox.confirm("确定要删除此模板吗？", "确认删除", {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning",
        });

        await templateAPI.deleteTemplate(id);
        ElMessage.success("模板删除成功");
        loadTemplates();
      } catch (error) {
        if (error !== "cancel") {
          console.error("删除模板失败:", error);
          ElMessage.error("删除模板失败");
        }
      }
    };

    const saveAsTemplate = async () => {
      if (!reportForm.work_content) {
        ElMessage.warning("请先填写工作内容");
        return;
      }

      try {
        const { value: templateName } = await ElMessageBox.prompt(
          "请输入模板名称",
          "保存为模板",
          {
            confirmButtonText: "保存",
            cancelButtonText: "取消",
            inputValue: `${reportForm.work_content}模板`,
            inputValidator: (value) => {
              if (!value) {
                return "模板名称不能为空";
              }
            },
          }
        );

        const templateData = {
          name: templateName,
          unit: reportForm.unit,
          substations: reportForm.substations,
          work_content: reportForm.work_content,
          related_devices: reportForm.related_devices,
          tools: reportForm.tools,
          operation: reportForm.operation,
          safety_measures: reportForm.safety_measures,
        };

        await templateAPI.createTemplate(templateData);
        ElMessage.success("模板保存成功");
      } catch (error) {
        if (error !== "cancel") {
          console.error("保存模板失败:", error);
          ElMessage.error("保存模板失败");
        }
      }
    };

    const formatDate = (dateString) => {
      if (!dateString) return "";
      const date = new Date(dateString);
      return date.toLocaleDateString("zh-CN");
    };

    const formatDateTime = (dateString) => {
      if (!dateString) return "";
      const date = new Date(dateString);
      return date.toLocaleString("zh-CN");
    };

    const goBack = () => {
      router.push("/user/manage-reports");
    };

    const loadReport = async (id) => {
      try {
        const response = await reportAPI.getReportById(id);
        Object.assign(reportForm, response);

        // 设置工作时间范围
        if (response.work_time_start && response.work_time_end) {
          workTimeRange.value = [
            response.work_time_start,
            response.work_time_end,
          ];
        }
      } catch (error) {
        console.error("加载报备单失败:", error);
        ElMessage.error("加载报备单失败");
      }
    };

    // 监听工作时间范围变化
    watch(workTimeRange, (newVal) => {
      if (newVal && newVal.length === 2) {
        reportForm.work_time_start = newVal[0];
        reportForm.work_time_end = newVal[1];
      }
    });

    // 生命周期
    onMounted(() => {
      generateReportNo();

      // 设置当前用户单位
      const userUnit = localStorage.getItem("userUnit");
      if (userUnit) {
        reportForm.unit = userUnit;
      }

      // 设置工作负责人为当前用户
      const realName = localStorage.getItem("realName");
      if (realName) {
        reportForm.responsible_person = realName;
      }

      // 如果有ID参数，加载对应的报备单
      const reportId = route.query.id;
      if (reportId) {
        loadReport(reportId);
      }
    });

    return {
      reportForm,
      reportFormRef,
      workTimeRange,
      loading,
      dialogVisible,
      previewData,
      templateList,
      unitOptions,
      substationGroups,
      workContentGroups,
      deviceOptions,
      toolOptions,
      statusText,
      statusType,
      isFormValid,
      formRules,
      defaultTimes,
      generateReportNo,
      getRiskType,
      getWorkContentDesc,
      onWorkContentChange,
      addDevice,
      removeDevice,
      clearDevices,
      onDeviceChange,
      validateIP,
      getDeviceType,
      addTool,
      removeTool,
      onToolChange,
      validateMAC,
      saveReport,
      submitReport,
      resetForm,
      showPreview,
      handlePrint,
      showTemplateList,
      selectTemplate,
      deleteTemplate,
      saveAsTemplate,
      formatDate,
      formatDateTime,
      goBack,
      Refresh,
      User,
      Plus,
      Delete,
      DocumentAdd,
      Promotion,
      ArrowLeft,
    };
  },
};
</script>

<style scoped>
.create-report-container {
  padding: 20px;
  background-color: #f5f7fa;
  min-height: 100vh;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
}

.form-card {
  margin-top: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.readonly-input {
  background-color: #f5f7fa;
}

.work-content-container,
.device-container,
.tool-container {
  width: 100%;
}

.device-table,
.tool-table {
  margin-top: 10px;
}

.form-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #ebeef5;
}

/* 预览样式 */
.preview-content {
  font-family: "Microsoft YaHei", sans-serif;
}

.preview-header {
  text-align: center;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 2px solid #409eff;
}

.preview-header h3 {
  color: #303133;
  margin-bottom: 16px;
}

.preview-info {
  display: flex;
  justify-content: space-between;
  color: #606266;
  font-size: 14px;
}

.preview-body {
  margin-top: 20px;
}

.preview-body h4 {
  color: #409eff;
  margin: 16px 0 8px 0;
  font-size: 15px;
}

.preview-text {
  background-color: #f8f9fa;
  padding: 12px;
  border-radius: 4px;
  border-left: 4px solid #409eff;
  margin-bottom: 16px;
  white-space: pre-wrap;
  line-height: 1.6;
}

.preview-remarks {
  margin-top: 20px;
  padding: 16px;
  background-color: #fff9e6;
  border-radius: 4px;
  border: 1px solid #ffd700;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .create-report-container {
    padding: 10px;
  }

  .el-col {
    margin-bottom: 16px;
  }

  .form-actions {
    flex-wrap: wrap;
  }

  .form-actions .el-button {
    flex: 1;
    min-width: 120px;
  }
}
</style>