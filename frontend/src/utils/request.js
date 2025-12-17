// src/utils/request.js 或 src/api/index.js
import axios from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'
import store from '@/store'  

// 创建axios实例
const service = axios.create({
  baseURL: 'http://localhost:8000/api', // 后端API地址
  timeout: 10000 // 请求超时时间
})

// 请求拦截器
service.interceptors.request.use(
  config => {
    // 优先从store获取token，如果store中没有，再从localStorage获取
    const token = store.state.token || localStorage.getItem('token')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
      config.headers['X-Token'] = token
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  response => {
    return response
  },
  error => {
    if (error.response?.status === 401) {
      // token过期或无效
      store.dispatch('logout')
      router.push('/')
      ElMessage.error('登录已过期，请重新登录')
    }
    return Promise.reject(error)
  }
)

export default service