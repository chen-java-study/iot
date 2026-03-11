import axios from 'axios'
import { ElMessage } from 'element-plus'

const request = axios.create({
  baseURL: '/',
  timeout: 30000
})

request.interceptors.request.use(
  config => {
    const token = localStorage.getItem('admin_token')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  error => Promise.reject(error)
)

request.interceptors.response.use(
  response => {
    const res = response.data
    if (res.code !== 200) {
      ElMessage.error(res.message || '请求失败')
      if (res.code === 401) {
        // 登录接口本身返回 401 时，不强制跳转，避免吞掉错误提示
        if (!response.config?.url?.includes('/api/v1/admin/login')) {
          localStorage.removeItem('admin_token')
          window.location.href = '/login'
        }
      }
      // 保留完整的响应信息方便错误处理
      const error = new Error(res.message || '请求失败')
      error.response = { data: res }
      return Promise.reject(error)
    }
    return res.data
  },
  error => {
    const backendMsg = error.response?.data?.message || error.response?.data?.Message
    ElMessage.error(backendMsg || error.message || '网络错误')
    if (error.response && error.response.status === 401) {
      // 登录接口 401 时不跳转；其余接口 401 仍跳转登录
      if (!error.config?.url?.includes('/api/v1/admin/login')) {
        localStorage.removeItem('admin_token')
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  }
)

export default request
