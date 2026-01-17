import request from '../utils/request'

export function createRechargeOrder(data) {
  return request({
    url: '/payment/create',
    method: 'post',
    data
  })
}

export function queryPaymentStatus(tradeNo) {
  return request({
    url: '/payment/status',
    method: 'get',
    params: { trade_no: tradeNo }
  })
}
