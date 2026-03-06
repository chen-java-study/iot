import axios from 'axios'
import { showToast } from 'vant'

const request = axios.create({
  baseURL: '/api/v1',
  timeout: 30000
})

request.interceptors.response.use(
  response => {
    const res = response.data
    if (res.code !== 200) {
      // silent: true 时不弹 toast，由调用方自行处理
      if (!response.config.silent) {
        showToast(res.message || '请求失败')
      }
      return Promise.reject(new Error(res.message || '请求失败'))
    }
    return res.data
  },
  error => {
    if (!error.config?.silent) {
      showToast(error.message || '网络错误')
    }
    return Promise.reject(error)
  }
)

export default request
