import request from '../utils/request'

export function queryCard(keyword) {
  return request({
    url: '/card/query',
    method: 'get',
    params: { keyword },
    silent: true  // 查不到时不弹 toast，由页面自己显示提示
  })
}
