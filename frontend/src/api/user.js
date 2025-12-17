// src/api/user.js
import request from '@/utils/request'

// 用户API

// 登录
export function login(data) {
  return request({
    url: '/login',
    method: 'post',
    data
  })
}

// 登出
export function logout() {
  return request({
    url: '/logout',
    method: 'post'
  })
}

// 获取当前用户信息
export function getCurrentUser() {
  return request({
    url: '/users/current',
    method: 'get'
  })
}

// 更新用户资料
export function updateUserProfile(data) {
  return request({
    url: '/users/profile',
    method: 'put',
    data
  })
}

// 获取所有用户（管理员）
export function getAllUsers() {
  return request({
    url: '/users',
    method: 'get'
  })
}

// 创建用户（管理员）
export function createUser(data) {
  return request({
    url: '/admin/users',
    method: 'post',
    data
  })
}

// 更新用户（管理员）
export function updateUser(username, data) {
  return request({
    url: `/admin/users/${username}`,
    method: 'put',
    data
  })
}

// 删除用户（管理员）
export function deleteUser(username) {
  return request({
    url: `/admin/users/${username}`,
    method: 'delete'
  })
}