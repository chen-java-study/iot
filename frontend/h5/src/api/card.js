import request from '../utils/request'

export function queryCard(keyword) {
  return request({
    url: '/card/query',
    method: 'get',
    params: { keyword }
  })
}
