import request from '../utils/request'

// 登录
export function login(data) {
  return request({
    url: '/api/v1/admin/login',
    method: 'post',
    data
  })
}

// 获取统计数据
export function getStatistics() {
  return request({
    url: '/api/v1/admin/statistics',
    method: 'get'
  })
}

// 卡片管理
export function getCardList(params) {
  return request({
    url: '/api/v1/admin/cards',
    method: 'get',
    params
  })
}

export function createCard(data) {
  return request({
    url: '/api/v1/admin/cards',
    method: 'post',
    data
  })
}

export function updateCard(id, data) {
  return request({
    url: `/api/v1/admin/cards/${id}`,
    method: 'put',
    data
  })
}

export function deleteCard(id) {
  return request({
    url: `/api/v1/admin/cards/${id}`,
    method: 'delete'
  })
}

// 充值记录
export function getRechargeList(params) {
  return request({
    url: '/api/v1/admin/recharges',
    method: 'get',
    params
  })
}

// 系统配置
export function getConfig() {
  return request({
    url: '/api/v1/admin/config',
    method: 'get'
  })
}

export function updateConfig(data) {
  return request({
    url: '/api/v1/admin/config',
    method: 'post',
    data
  })
}
