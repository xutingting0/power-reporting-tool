// src/api/report.js
import request from '@/utils/request'

// 报备单API

// 创建报备单
export function createReport(data) {
  return request({
    url: '/reports',
    method: 'post',
    data
  })
}

// 获取用户报备单
export function getUserReports(params) {
  return request({
    url: '/reports/user',
    method: 'get',
    params
  })
}

// 获取所有报备单（管理员）
export function getAllReports(params) {
  return request({
    url: '/admin/reports/all',
    method: 'get',
    params
  })
}

// 获取待审核报备单
export function getPendingReports() {
  return request({
    url: '/reports/pending',
    method: 'get'
  })
}

// 获取报备单详情
export function getReportById(id) {
  return request({
    url: `/reports/${id}`,
    method: 'get'
  })
}

// 更新报备单
export function updateReport(id, data) {
  return request({
    url: `/reports/${id}`,
    method: 'put',
    data
  })
}

// 提交报备单
export function submitReport(id) {
  return request({
    url: `/reports/${id}/submit`,
    method: 'post'
  })
}

// 审核通过
export function approveReport(id, reviewComments) {
  return request({
    url: `/reports/${id}/approve`,
    method: 'post',
    data: { review_comments: reviewComments }
  })
}

// 驳回
export function rejectReport(id, reviewComments) {
  return request({
    url: `/reports/${id}/reject`,
    method: 'post',
    data: { review_comments: reviewComments }
  })
}

// 删除报备单
export function deleteReport(id) {
  return request({
    url: `/reports/${id}`,
    method: 'delete'
  })
}

// 批量删除
export function batchDeleteReports(ids) {
  return request({
    url: '/reports/batch-delete',
    method: 'post',
    data: { ids }
  })
}

// 获取操作日志
export function getReportLogs(id) {
  return request({
    url: `/reports/${id}/logs`,
    method: 'get'
  })
}

// 导出报备单
export function exportReports(params) {
  return request({
    url: '/reports/export',
    method: 'get',
    params,
    responseType: 'blob' // 重要：用于下载文件
  })
}

// 获取统计信息
export function getReportStatistics(params) {
  return request({
    url: '/reports/statistics',
    method: 'get',
    params
  })
}

// 获取仪表盘统计
export function getDashboardStats() {
  return request({
    url: '/dashboard/stats',
    method: 'get'
  })
}

// 获取最近报备单
export function getRecentReports(limit = 10) {
  return request({
    url: '/reports/recent',
    method: 'get',
    params: { limit }
  })
}

// 使用模板创建报备单
export function useTemplate(templateId) {
  return request({
    url: `/templates/${templateId}/use`,
    method: 'post'
  })
}