<template>
  <div class="statistics-container">
    <TheHeader />
    <el-main>
      <el-page-header>
        <template #content>
          <span>统计分析</span>
        </template>
        <template #extra>
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            @change="loadStatistics"
          />
        </template>
      </el-page-header>

      <el-row :gutter="20" class="statistics-row">
        <!-- 数据概览 -->
        <el-col :span="24">
          <el-card class="statistics-card">
            <template #header>
              <div class="card-header">
                <span>数据概览</span>
                <el-tag type="info"
                  >统计周期: {{ dateRange[0] }} 至 {{ dateRange[1] }}</el-tag
                >
              </div>
            </template>
            <el-row :gutter="20">
              <el-col :span="6">
                <el-statistic title="总报备单数" :value="statistics.total">
                  <template #suffix>
                    <el-icon><Document /></el-icon>
                  </template>
                </el-statistic>
              </el-col>
              <el-col :span="6">
                <el-statistic title="待审核" :value="statistics.pending">
                  <template #suffix>
                    <el-icon><Clock /></el-icon>
                  </template>
                </el-statistic>
              </el-col>
              <el-col :span="6">
                <el-statistic title="已通过" :value="statistics.approved">
                  <template #suffix>
                    <el-icon><CircleCheck /></el-icon>
                  </template>
                </el-statistic>
              </el-col>
              <el-col :span="6">
                <el-statistic title="已驳回" :value="statistics.rejected">
                  <template #suffix>
                    <el-icon><CircleClose /></el-icon>
                  </template>
                </el-statistic>
              </el-col>
            </el-row>
          </el-card>
        </el-col>

        <!-- 按状态统计 -->
        <el-col :span="12">
          <el-card class="statistics-card">
            <template #header>
              <span>报备单状态分布</span>
            </template>
            <div ref="statusChart" style="height: 300px"></div>
          </el-card>
        </el-col>

        <!-- 按单位统计 -->
        <el-col :span="12">
          <el-card class="statistics-card">
            <template #header>
              <span>各单位报备情况</span>
            </template>
            <div ref="unitChart" style="height: 300px"></div>
          </el-card>
        </el-col>

        <!-- 按工作内容统计 -->
        <el-col :span="24">
          <el-card class="statistics-card">
            <template #header>
              <span>热门工作内容排行</span>
            </template>
            <el-table
              :data="workContentStats"
              style="width: 100%"
              v-loading="loading"
            >
              <el-table-column prop="rank" label="排名" width="80">
                <template #default="{ $index }">
                  <el-tag :type="$index < 3 ? 'danger' : 'info'">
                    {{ $index + 1 }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="work_content" label="工作内容" />
              <el-table-column prop="count" label="报备次数" width="120">
                <template #default="{ row }">
                  <div class="progress-container">
                    <div
                      class="progress-bar"
                      :style="{ width: row.percentage + '%' }"
                    ></div>
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="count" label="数量" width="100">
                <template #default="{ row }"> {{ row.count }} 次 </template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-col>

        <!-- 按时间趋势 -->
        <el-col :span="24">
          <el-card class="statistics-card">
            <template #header>
              <span>报备趋势分析</span>
            </template>
            <div ref="trendChart" style="height: 350px"></div>
          </el-card>
        </el-col>

        <!-- 详细统计表格 -->
        <el-col :span="24">
          <el-card class="statistics-card">
            <template #header>
              <span>详细统计</span>
              <el-button type="primary" size="small" @click="exportStatistics">
                <el-icon><Download /></el-icon>
                导出数据
              </el-button>
            </template>
            <el-table
              :data="detailedStats"
              style="width: 100%"
              v-loading="loading"
            >
              <el-table-column prop="date" label="日期" width="120" />
              <el-table-column label="报备数量" width="100">
                <template #default="{ row }">
                  {{ row.total }}
                </template>
              </el-table-column>
              <el-table-column label="单位分布">
                <template #default="{ row }">
                  <div class="unit-tags">
                    <el-tag
                      v-for="unit in row.units"
                      :key="unit.name"
                      size="small"
                      type="info"
                      class="unit-tag"
                    >
                      {{ unit.name }}: {{ unit.count }}
                    </el-tag>
                  </div>
                </template>
              </el-table-column>
              <el-table-column label="状态分布" width="200">
                <template #default="{ row }">
                  <div class="status-tags">
                    <el-tag
                      v-for="status in row.statuses"
                      :key="status.name"
                      :type="getStatusType(status.name)"
                      size="small"
                      class="status-tag"
                    >
                      {{ getStatusText(status.name) }}: {{ status.count }}
                    </el-tag>
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-col>
      </el-row>
    </el-main>
  </div>
</template>

<script>
import { ref, reactive, onMounted, nextTick, onUnmounted } from "vue";
import { ElMessage } from "element-plus";
import * as echarts from "echarts"; // 修改这里：直接导入echarts
import {
  Document,
  Clock,
  CircleCheck,
  CircleClose,
  Download,
} from "@element-plus/icons-vue";
import TheHeader from "../../components/TheHeader.vue";
import { reportAPI } from "@/utils/request.js";

export default {
  name: "Statistics",
  components: {
    TheHeader,
  },
  setup() {
    const loading = ref(false);
    const statusChart = ref(null);
    const unitChart = ref(null);
    const trendChart = ref(null);

    let statusChartInstance = null;
    let unitChartInstance = null;
    let trendChartInstance = null;

    // 设置默认日期范围（最近30天）
    const defaultEndDate = new Date();
    const defaultStartDate = new Date();
    defaultStartDate.setDate(defaultStartDate.getDate() - 30);

    const dateRange = ref([
      defaultStartDate.toISOString().split("T")[0],
      defaultEndDate.toISOString().split("T")[0],
    ]);

    const statistics = reactive({
      total: 0,
      pending: 0,
      approved: 0,
      rejected: 0,
    });

    const workContentStats = ref([]);
    const detailedStats = ref([]);

    const customColors = [
      { color: "#f56c6c", percentage: 20 },
      { color: "#e6a23c", percentage: 40 },
      { color: "#5cb87a", percentage: 60 },
      { color: "#1989fa", percentage: 80 },
      { color: "#6f7ad3", percentage: 100 },
    ];

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

    const loadStatistics = async () => {
      loading.value = true;
      try {
        const params = {
          start_date: dateRange.value[0],
          end_date: dateRange.value[1],
        };

        // 模拟数据
        const mockData = {
          overview: {
            total: 125,
            pending: 15,
            approved: 85,
            rejected: 25,
          },
          statusStats: {
            draft: 25,
            submitted: 15,
            approved: 85,
            rejected: 25,
          },
          unitStats: {
            遵义供电局: 45,
            贵阳供电局: 32,
            六盘水供电局: 28,
            安顺供电局: 12,
            毕节供电局: 8,
          },
          workContentStats: [
            { work_content: "设备巡检", count: 35, percentage: 70 },
            { work_content: "设备维护", count: 28, percentage: 56 },
            { work_content: "故障处理", count: 20, percentage: 40 },
            { work_content: "系统调试", count: 18, percentage: 36 },
            { work_content: "参数修改", count: 12, percentage: 24 },
          ],
          trendStats: [
            { date: "2024-01-01", total: 5, approved: 3, pending: 2 },
            { date: "2024-01-02", total: 8, approved: 6, pending: 2 },
            { date: "2024-01-03", total: 12, approved: 10, pending: 2 },
            { date: "2024-01-04", total: 7, approved: 5, pending: 2 },
            { date: "2024-01-05", total: 15, approved: 12, pending: 3 },
          ],
          detailedStats: [
            {
              date: "2024-01-01",
              total: 5,
              units: [
                { name: "遵义供电局", count: 3 },
                { name: "贵阳供电局", count: 2 },
              ],
              statuses: [
                { name: "approved", count: 3 },
                { name: "submitted", count: 2 },
              ],
            },
            // 更多模拟数据...
          ],
        };

        // 更新统计数据
        Object.assign(statistics, mockData.overview);
        workContentStats.value = mockData.workContentStats.map(
          (item, index) => ({
            ...item,
            rank: index + 1,
          })
        );
        detailedStats.value = mockData.detailedStats;

        // 渲染图表
        nextTick(() => {
          renderStatusChart(mockData.statusStats);
          renderUnitChart(mockData.unitStats);
          renderTrendChart(mockData.trendStats);
        });
      } catch (error) {
        console.error("加载统计数据失败:", error);
        ElMessage.error("加载统计数据失败");
      } finally {
        loading.value = false;
      }
    };

    const renderStatusChart = (statusStats) => {
      if (!statusChart.value) return;

      // 销毁旧的图表实例
      if (statusChartInstance) {
        statusChartInstance.dispose();
      }

      statusChartInstance = echarts.init(statusChart.value);
      const options = {
        tooltip: {
          trigger: "item",
          formatter: "{a} <br/>{b}: {c} ({d}%)",
        },
        legend: {
          orient: "vertical",
          left: "left",
          data: Object.keys(statusStats).map((key) => getStatusText(key)),
        },
        series: [
          {
            name: "报备单状态",
            type: "pie",
            radius: ["40%", "70%"],
            avoidLabelOverlap: false,
            itemStyle: {
              borderRadius: 10,
              borderColor: "#fff",
              borderWidth: 2,
            },
            label: {
              show: false,
              position: "center",
            },
            emphasis: {
              label: {
                show: true,
                fontSize: "20",
                fontWeight: "bold",
              },
            },
            labelLine: {
              show: false,
            },
            data: Object.entries(statusStats).map(([key, value]) => ({
              name: getStatusText(key),
              value: value,
              itemStyle: {
                color: getPieColor(key),
              },
            })),
          },
        ],
      };

      statusChartInstance.setOption(options);

      // 响应式
      window.addEventListener("resize", () => {
        if (statusChartInstance) {
          statusChartInstance.resize();
        }
      });
    };

    const renderUnitChart = (unitStats) => {
      if (!unitChart.value) return;

      // 销毁旧的图表实例
      if (unitChartInstance) {
        unitChartInstance.dispose();
      }

      unitChartInstance = echarts.init(unitChart.value);
      const options = {
        tooltip: {
          trigger: "axis",
          axisPointer: {
            type: "shadow",
          },
        },
        grid: {
          left: "3%",
          right: "4%",
          bottom: "3%",
          containLabel: true,
        },
        xAxis: {
          type: "category",
          data: Object.keys(unitStats),
          axisLabel: {
            rotate: 45,
          },
        },
        yAxis: {
          type: "value",
        },
        series: [
          {
            name: "报备数量",
            type: "bar",
            data: Object.values(unitStats),
            itemStyle: {
              color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                { offset: 0, color: "#83bff6" },
                { offset: 0.5, color: "#188df0" },
                { offset: 1, color: "#188df0" },
              ]),
            },
            emphasis: {
              itemStyle: {
                color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                  { offset: 0, color: "#2378f7" },
                  { offset: 0.7, color: "#2378f7" },
                  { offset: 1, color: "#83bff6" },
                ]),
              },
            },
          },
        ],
      };

      unitChartInstance.setOption(options);

      // 响应式
      window.addEventListener("resize", () => {
        if (unitChartInstance) {
          unitChartInstance.resize();
        }
      });
    };

    const renderTrendChart = (trendStats) => {
      if (!trendChart.value) return;

      // 销毁旧的图表实例
      if (trendChartInstance) {
        trendChartInstance.dispose();
      }

      trendChartInstance = echarts.init(trendChart.value);
      const dates = trendStats.map((item) => item.date);
      const totalData = trendStats.map((item) => item.total);
      const approvedData = trendStats.map((item) => item.approved || 0);
      const pendingData = trendStats.map((item) => item.pending || 0);

      const options = {
        tooltip: {
          trigger: "axis",
          axisPointer: {
            type: "cross",
            label: {
              backgroundColor: "#6a7985",
            },
          },
        },
        legend: {
          data: ["总报备数", "已通过", "待审核"],
        },
        grid: {
          left: "3%",
          right: "4%",
          bottom: "3%",
          containLabel: true,
        },
        xAxis: {
          type: "category",
          boundaryGap: false,
          data: dates,
        },
        yAxis: {
          type: "value",
        },
        series: [
          {
            name: "总报备数",
            type: "line",
            smooth: true,
            lineStyle: {
              width: 4,
            },
            areaStyle: {
              opacity: 0.5,
            },
            emphasis: {
              focus: "series",
            },
            data: totalData,
          },
          {
            name: "已通过",
            type: "line",
            smooth: true,
            lineStyle: {
              width: 3,
            },
            emphasis: {
              focus: "series",
            },
            data: approvedData,
          },
          {
            name: "待审核",
            type: "line",
            smooth: true,
            lineStyle: {
              width: 3,
            },
            emphasis: {
              focus: "series",
            },
            data: pendingData,
          },
        ],
      };

      trendChartInstance.setOption(options);

      // 响应式
      window.addEventListener("resize", () => {
        if (trendChartInstance) {
          trendChartInstance.resize();
        }
      });
    };

    const getPieColor = (status) => {
      const colorMap = {
        draft: "#909399",
        submitted: "#e6a23c",
        approved: "#67c23a",
        rejected: "#f56c6c",
      };
      return colorMap[status] || "#409EFF";
    };

    const exportStatistics = async () => {
      try {
        ElMessage.success("导出功能已触发，请稍候...");
        // 实际项目中这里会调用API导出数据
      } catch (error) {
        console.error("导出失败:", error);
        ElMessage.error("导出失败");
      }
    };

    const cleanupCharts = () => {
      if (statusChartInstance) {
        statusChartInstance.dispose();
        statusChartInstance = null;
      }
      if (unitChartInstance) {
        unitChartInstance.dispose();
        unitChartInstance = null;
      }
      if (trendChartInstance) {
        trendChartInstance.dispose();
        trendChartInstance = null;
      }

      // 移除事件监听
      const resizeHandler = () => {};
      window.removeEventListener("resize", resizeHandler);
    };

    onMounted(() => {
      loadStatistics();
    });

    onUnmounted(() => {
      cleanupCharts();
    });

    return {
      loading,
      statusChart,
      unitChart,
      trendChart,
      dateRange,
      statistics,
      workContentStats,
      detailedStats,
      customColors,
      getStatusText,
      getStatusType,
      loadStatistics,
      exportStatistics,
      Document,
      Clock,
      CircleCheck,
      CircleClose,
      Download,
    };
  },
};
</script>

<style scoped>
.statistics-container {
  padding: 20px;
  background-color: #f5f7fa;
  min-height: 100vh;
}

.statistics-row {
  margin-top: 20px;
}

.statistics-card {
  margin-bottom: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.progress-container {
  width: 100%;
  height: 20px;
  background-color: #f0f0f0;
  border-radius: 10px;
  overflow: hidden;
}

.progress-bar {
  height: 100%;
  background: linear-gradient(90deg, #67c23a, #5cb87a);
  transition: width 0.3s ease;
}

.unit-tags,
.status-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
}

.unit-tag,
.status-tag {
  margin: 2px;
}
</style>