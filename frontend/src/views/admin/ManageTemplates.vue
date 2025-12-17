<template>
  <div class="manage-templates-container">
    <TheHeader />
    <el-main>
      <el-page-header>
        <template #content>
          <span>模板管理</span>
          <el-tag type="info" size="small" style="margin-left: 10px">
            共 {{ pagination.total }} 个模板
          </el-tag>
        </template>
      </el-page-header>
      <el-card>
        <!-- 查询表单 -->
        <div class="search-form">
          <el-form :model="searchForm" label-width="80px" inline>
            <el-form-item label="模板名称">
              <el-input
                v-model="searchForm.name"
                placeholder="请输入模板名称"
                clearable
              />
            </el-form-item>
            <el-form-item label="工作内容">
              <el-input
                v-model="searchForm.workContent"
                placeholder="请输入工作内容"
                clearable
              />
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
            <el-form-item label="创建人">
              <el-input
                v-model="searchForm.creator"
                placeholder="请输入创建人"
                clearable
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

        <!-- 批量操作 -->
        <div class="batch-operation" v-if="selectedTemplates.length > 0">
          <el-space>
            <span>已选择 {{ selectedTemplates.length }} 个模板</span>
            <el-button type="danger" @click="batchDelete">批量删除</el-button>
          </el-space>
        </div>

        <!-- 模板表格 -->
        <el-table
          :data="templateData"
          style="width: 100%"
          v-loading="loading"
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column prop="name" label="模板名称" width="150" sortable />
          <el-table-column prop="unit" label="单位" width="120" />
          <el-table-column label="变电站" width="150">
            <template #default="{ row }">
              {{ row.substations ? row.substations.join(", ") : "" }}
            </template>
          </el-table-column>
          <el-table-column
            prop="work_content"
            label="工作内容"
            min-width="200"
            show-overflow-tooltip
          />
          <el-table-column label="涉及设备" width="120">
            <template #default="{ row }">
              <el-popover
                placement="top"
                trigger="hover"
                width="300"
                v-if="row.related_devices && row.related_devices.length > 0"
              >
                <template #reference>
                  <el-tag type="info" size="small">
                    {{ row.related_devices.length }}台设备
                  </el-tag>
                </template>
                <div>
                  <p
                    v-for="(device, index) in row.related_devices"
                    :key="index"
                  >
                    {{ device.name }} ({{ device.ip }})
                  </p>
                </div>
              </el-popover>
              <span v-else>无</span>
            </template>
          </el-table-column>
          <el-table-column label="使用工具" width="100">
            <template #default="{ row }">
              <el-tag type="info" size="small">
                {{ row.tools ? row.tools.length : 0 }}个
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="创建人" width="100">
            <template #default="{ row }">
              {{ row.creator_name || "系统" }}
            </template>
          </el-table-column>
          <el-table-column label="创建时间" width="150">
            <template #default="{ row }">
              {{ formatDateTime(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="更新时间" width="150">
            <template #default="{ row }">
              {{ formatDateTime(row.updated_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" size="small" @click="handleEdit(row)">
                编辑
              </el-button>
              <el-button type="success" size="small" @click="useTemplate(row)">
                使用
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

        <!-- 新增模板表单 -->
        <div class="add-template-form">
          <h3>{{ isEditing ? "编辑模板" : "新增模板" }}</h3>
          <el-form
            ref="templateFormRef"
            :model="templateForm"
            :rules="templateRules"
            label-width="120px"
          >
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="模板名称" prop="name" required>
                  <el-input
                    v-model="templateForm.name"
                    placeholder="请输入模板名称"
                    clearable
                  />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="单位" prop="unit">
                  <el-select
                    v-model="templateForm.unit"
                    placeholder="请选择单位"
                    filterable
                    clearable
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
              </el-col>
            </el-row>

            <el-form-item label="变电站" prop="substations">
              <el-select
                v-model="templateForm.substations"
                placeholder="请选择变电站"
                multiple
                collapse-tags
                collapse-tags-tooltip
                filterable
                clearable
                style="width: 100%"
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

            <el-form-item label="工作内容" prop="work_content" required>
              <el-input
                v-model="templateForm.work_content"
                type="textarea"
                :rows="3"
                placeholder="请输入工作内容"
                maxlength="500"
                show-word-limit
              />
            </el-form-item>

            <el-form-item label="涉及设备">
              <div class="device-container">
                <el-button
                  type="primary"
                  @click="addDevice"
                  size="small"
                  :icon="Plus"
                >
                  添加设备
                </el-button>
                <el-button @click="clearDevices" size="small" :icon="Delete">
                  清空
                </el-button>

                <el-table
                  :data="templateForm.related_devices"
                  border
                  size="small"
                  class="device-table"
                  style="margin-top: 10px"
                >
                  <el-table-column label="设备名称" prop="name" min-width="150">
                    <template #default="{ row, $index }">
                      <el-select
                        v-model="row.name"
                        placeholder="选择设备"
                        filterable
                        allow-create
                        style="width: 100%"
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
                  <el-table-column label="IP地址" prop="ip" min-width="120">
                    <template #default="{ row, $index }">
                      <el-input v-model="row.ip" placeholder="请输入IP地址" />
                    </template>
                  </el-table-column>
                  <el-table-column label="设备类型" prop="type" min-width="100">
                    <template #default="{ row }">
                      <el-select
                        v-model="row.type"
                        placeholder="选择类型"
                        style="width: 100%"
                      >
                        <el-option label="测控设备" value="测控设备" />
                        <el-option label="保护设备" value="保护设备" />
                        <el-option label="计算机" value="计算机" />
                        <el-option label="通讯设备" value="通讯设备" />
                        <el-option label="网络设备" value="网络设备" />
                        <el-option label="安全设备" value="安全设备" />
                        <el-option label="服务器" value="服务器" />
                      </el-select>
                    </template>
                  </el-table-column>
                  <el-table-column label="操作" width="60">
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

            <el-form-item label="使用工具">
              <div class="tool-container">
                <el-button
                  type="primary"
                  @click="addTool"
                  size="small"
                  :icon="Plus"
                >
                  添加工具
                </el-button>

                <el-table
                  :data="templateForm.tools"
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
                  <el-table-column label="MAC地址" prop="mac" min-width="120">
                    <template #default="{ row, $index }">
                      <el-input v-model="row.mac" placeholder="请输入MAC地址" />
                    </template>
                  </el-table-column>
                  <el-table-column
                    label="序列号"
                    prop="serial_number"
                    min-width="120"
                  >
                    <template #default="{ row }">
                      <el-input
                        v-model="row.serial_number"
                        placeholder="请输入序列号"
                      />
                    </template>
                  </el-table-column>
                  <el-table-column label="类型" prop="type" min-width="100">
                    <template #default="{ row }">
                      <el-select
                        v-model="row.type"
                        placeholder="选择类型"
                        style="width: 100%"
                      >
                        <el-option label="计算机" value="计算机" />
                        <el-option label="存储设备" value="存储设备" />
                        <el-option label="软件" value="软件" />
                        <el-option label="仪器仪表" value="仪器仪表" />
                        <el-option label="通讯设备" value="通讯设备" />
                      </el-select>
                    </template>
                  </el-table-column>
                  <el-table-column label="操作" width="60">
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

            <el-form-item label="操作内容" prop="operation">
              <el-input
                v-model="templateForm.operation"
                type="textarea"
                :rows="3"
                placeholder="请输入操作内容"
                maxlength="500"
                show-word-limit
              />
            </el-form-item>

            <el-form-item label="安全措施" prop="safety_measures">
              <el-input
                v-model="templateForm.safety_measures"
                type="textarea"
                :rows="3"
                placeholder="请输入安全措施"
                maxlength="500"
                show-word-limit
              />
            </el-form-item>

            <el-form-item>
              <el-button type="primary" @click="saveTemplate" :icon="Check">
                {{ isEditing ? "保存修改" : "添加模板" }}
              </el-button>
              <el-button @click="resetTemplateForm" :icon="Refresh">
                重置
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-card>
    </el-main>
  </div>
</template>

<script>
import { ref, reactive, onMounted, computed } from "vue";
import { useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import { Search, Refresh, Plus, Delete, Check } from "@element-plus/icons-vue";
import TheHeader from "../../components/TheHeader.vue";
import { templateAPI } from "@/utils/request.js";

export default {
  name: "ManageTemplates",
  components: {
    TheHeader,
  },
  setup() {
    const router = useRouter();
    const templateFormRef = ref();

    const templateData = ref([]);
    const loading = ref(false);
    const selectedTemplates = ref([]);
    const isEditing = ref(false);

    const searchForm = reactive({
      name: "",
      workContent: "",
      unit: "",
      creator: "",
    });

    const pagination = reactive({
      current: 1,
      size: 10,
      total: 0,
    });

    const templateForm = reactive({
      id: null,
      name: "",
      unit: "",
      substations: [],
      work_content: "",
      related_devices: [],
      tools: [],
      operation: "",
      safety_measures: "",
    });

    const unitOptions = [
      { value: "遵义供电局", label: "遵义供电局" },
      { value: "贵阳供电局", label: "贵阳供电局" },
      { value: "六盘水供电局", label: "六盘水供电局" },
      { value: "安顺供电局", label: "安顺供电局" },
      { value: "毕节供电局", label: "毕节供电局" },
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

    const templateRules = {
      name: [
        { required: true, message: "请输入模板名称", trigger: "blur" },
        { max: 50, message: "模板名称不能超过50个字符", trigger: "blur" },
      ],
      work_content: [
        { required: true, message: "请输入工作内容", trigger: "blur" },
      ],
    };

    const formatDateTime = (dateTime) => {
      if (!dateTime) return "";
      return new Date(dateTime).toLocaleString("zh-CN", { hour12: false });
    };

    const loadTemplates = async () => {
      loading.value = true;
      try {
        const params = {
          page: pagination.current,
          page_size: pagination.size,
          name: searchForm.name,
          work_content: searchForm.workContent,
          unit: searchForm.unit,
          creator: searchForm.creator,
        };

        const response = await templateAPI.getAllTemplates(params);
        templateData.value = response.data || response;
        pagination.total = response.total || response.length;
      } catch (error) {
        console.error("加载模板失败:", error);
        ElMessage.error("加载模板失败");
      } finally {
        loading.value = false;
      }
    };

    const handleSearch = () => {
      pagination.current = 1;
      loadTemplates();
    };

    const resetSearch = () => {
      Object.assign(searchForm, {
        name: "",
        workContent: "",
        unit: "",
        creator: "",
      });
      pagination.current = 1;
      loadTemplates();
    };

    const handleSelectionChange = (selection) => {
      selectedTemplates.value = selection;
    };

    const addDevice = () => {
      templateForm.related_devices.push({
        name: "",
        ip: "",
        type: "测控设备",
      });
    };

    const removeDevice = (index) => {
      templateForm.related_devices.splice(index, 1);
    };

    const clearDevices = () => {
      templateForm.related_devices = [];
    };

    const addTool = () => {
      templateForm.tools.push({
        name: "",
        mac: "",
        serial_number: "",
        type: "计算机",
      });
    };

    const removeTool = (index) => {
      templateForm.tools.splice(index, 1);
    };

    const saveTemplate = async () => {
      try {
        await templateFormRef.value.validate();

        const templateDataToSend = {
          name: templateForm.name,
          unit: templateForm.unit,
          substations: templateForm.substations,
          work_content: templateForm.work_content,
          related_devices: templateForm.related_devices,
          tools: templateForm.tools,
          operation: templateForm.operation,
          safety_measures: templateForm.safety_measures,
        };

        if (isEditing.value) {
          await templateAPI.updateTemplate(templateForm.id, templateDataToSend);
          ElMessage.success("模板更新成功");
        } else {
          await templateAPI.createTemplate(templateDataToSend);
          ElMessage.success("模板添加成功");
        }

        resetTemplateForm();
        loadTemplates();
      } catch (error) {
        console.error("保存模板失败:", error);
        ElMessage.error("保存模板失败");
      }
    };

    const resetTemplateForm = () => {
      Object.assign(templateForm, {
        id: null,
        name: "",
        unit: "",
        substations: [],
        work_content: "",
        related_devices: [],
        tools: [],
        operation: "",
        safety_measures: "",
      });
      isEditing.value = false;
      templateFormRef.value?.clearValidate();
    };

    const handleEdit = (row) => {
      isEditing.value = true;
      Object.assign(templateForm, {
        id: row.id,
        name: row.name,
        unit: row.unit || "",
        substations: row.substations || [],
        work_content: row.work_content,
        related_devices: row.related_devices || [],
        tools: row.tools || [],
        operation: row.operation || "",
        safety_measures: row.safety_measures || "",
      });

      // 滚动到表单位置
      document
        .querySelector(".add-template-form")
        .scrollIntoView({ behavior: "smooth" });
    };

    const useTemplate = (row) => {
      ElMessageBox.confirm("确定要使用此模板创建报备单吗？", "确认使用模板", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "info",
      }).then(() => {
        // 将模板数据传递给创建报备单页面
        const templateData = {
          unit: row.unit,
          substations: row.substations,
          work_content: row.work_content,
          related_devices: row.related_devices,
          tools: row.tools,
          operation: row.operation,
          safety_measures: row.safety_measures,
        };

        // 保存到localStorage，在CreateReport页面读取
        localStorage.setItem("templateData", JSON.stringify(templateData));
        router.push("/user/create-report");
      });
    };

    const handleDelete = async (row) => {
      try {
        await ElMessageBox.confirm(
          `确定要删除模板 "${row.name}" 吗？`,
          "确认删除",
          {
            confirmButtonText: "确定",
            cancelButtonText: "取消",
            type: "warning",
          }
        );

        await templateAPI.deleteTemplate(row.id);
        ElMessage.success("模板删除成功");
        loadTemplates();
      } catch (error) {
        if (error === "cancel") return;
        console.error("删除模板失败:", error);
        ElMessage.error("删除模板失败");
      }
    };

    const batchDelete = async () => {
      if (selectedTemplates.value.length === 0) {
        ElMessage.warning("请选择要删除的模板");
        return;
      }

      try {
        await ElMessageBox.confirm(
          `确定要删除选中的 ${selectedTemplates.value.length} 个模板吗？`,
          "确认批量删除",
          {
            confirmButtonText: "确定",
            cancelButtonText: "取消",
            type: "warning",
          }
        );

        const deletePromises = selectedTemplates.value.map((template) =>
          templateAPI.deleteTemplate(template.id)
        );
        await Promise.all(deletePromises);
        ElMessage.success("批量删除成功");
        selectedTemplates.value = [];
        loadTemplates();
      } catch (error) {
        if (error === "cancel") return;
        console.error("批量删除失败:", error);
        ElMessage.error("批量删除失败");
      }
    };

    const handleSizeChange = (size) => {
      pagination.size = size;
      pagination.current = 1;
      loadTemplates();
    };

    const handleCurrentChange = (page) => {
      pagination.current = page;
      loadTemplates();
    };

    onMounted(() => {
      loadTemplates();
    });

    return {
      templateFormRef,
      templateData,
      loading,
      selectedTemplates,
      isEditing,
      searchForm,
      pagination,
      templateForm,
      unitOptions,
      substationGroups,
      deviceOptions,
      toolOptions,
      templateRules,
      formatDateTime,
      handleSearch,
      resetSearch,
      handleSelectionChange,
      addDevice,
      removeDevice,
      clearDevices,
      addTool,
      removeTool,
      saveTemplate,
      resetTemplateForm,
      handleEdit,
      useTemplate,
      handleDelete,
      batchDelete,
      handleSizeChange,
      handleCurrentChange,
      Search,
      Refresh,
      Plus,
      Delete,
      Check,
    };
  },
};
</script>

<style scoped>
.manage-templates-container {
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

.add-template-form {
  margin-top: 30px;
  padding: 20px;
  background-color: #fff;
  border-radius: 4px;
  border-top: 1px solid #ebeef5;
}

.add-template-form h3 {
  margin-top: 0;
  margin-bottom: 20px;
  color: #409eff;
}

.device-table,
.tool-table {
  margin-top: 10px;
}

.pagination-container {
  margin-top: 20px;
  text-align: right;
}
</style>