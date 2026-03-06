import request from '../utils/request'

export function createRechargeOrder(data) {
  return request({
    url: '/v1/payment/create',
    method: 'post',
    data
  })
}

export function queryPaymentStatus(tradeNo) {
  return request({
    url: '/v1/payment/status',
    method: 'get',
    params: { trade_no: tradeNo }
  })
}

