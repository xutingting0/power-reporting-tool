// src/api/template.js
import request from '@/utils/request'

// 模板API

// 获取所有模板
export function getAllTemplates() {
  return request({
    url: '/templates',
    method: 'get'
  })
}

// 获取模板详情
export function getTemplateById(id) {
  return request({
    url: `/templates/${id}`,
    method: 'get'
  })
}

// 创建模板（管理员）
export function createTemplate(data) {
  return request({
    url: '/admin/templates',
    method: 'post',
    data
  })
}

// 更新模板（管理员）
export function updateTemplate(id, data) {
  return request({
    url: `/admin/templates/${id}`,
    method: 'put',
    data
  })
}

// 删除模板（管理员）
export function deleteTemplate(id) {
  return request({
    url: `/admin/templates/${id}`,
    method: 'delete'
  })
}

// 使用模板创建报备单
export function useTemplate(id) {
  return request({
    url: `/templates/${id}/use`,
    method: 'post'
  })
}