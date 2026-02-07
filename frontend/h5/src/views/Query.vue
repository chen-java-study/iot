<template>
  <div class="query-page">
    <van-nav-bar title="物联网卡查询" />

    <!-- 搜索区域 -->
    <div class="search-section">
      <van-search
        v-model="keyword"
        show-action
        placeholder="请输入设备号或SIM卡号"
        @search="onSearch"
      >
        <template #action>
          <van-button size="small" type="primary" @click="onSearch" :loading="searching">查询</van-button>
        </template>
      </van-search>
    </div>

    <!-- 空状态提示 -->
    <div v-if="!cardInfo && !searchError" class="empty-tips">
      <van-empty image="search" description="请输入设备号或SIM卡号进行查询" />
    </div>

    <!-- 查询失败提示 -->
    <div v-if="searchError" class="empty-tips">
      <van-empty image="error" :description="searchError" />
    </div>

    <!-- 查询结果卡片 -->
    <div v-if="cardInfo" class="result-section">
      <div class="card-result">
        <div class="card-header">
          <div class="card-title">卡片信息</div>
          <van-tag :type="getStatusType(cardInfo.status)" size="medium">{{ cardInfo.status_text }}</van-tag>
        </div>
        <div class="card-body">
          <div class="info-row">
            <span class="info-label">卡号</span>
            <span class="info-value">{{ cardInfo.card_no }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">设备号</span>
            <span class="info-value">{{ cardInfo.device_no }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">到期时间</span>
            <span class="info-value" :class="{ 'text-danger': cardInfo.days_remaining <= 30 }">
              {{ cardInfo.expire_date }}
              <span class="days-tag" v-if="cardInfo.days_remaining > 0">（剩余{{ cardInfo.days_remaining }}天）</span>
              <span class="days-tag text-danger" v-else>（已过期）</span>
            </span>
          </div>
        </div>
      </div>

      <!-- 充值成功结果 -->
      <div v-if="rechargeSuccess" class="success-section">
        <div class="success-card">
          <van-icon name="checked" class="success-icon" />
          <div class="success-title">充值成功</div>
          <div class="success-detail">
            <div class="detail-row">
              <span>充值金额</span>
              <span class="amount">¥{{ rechargeResult.amount }}</span>
            </div>
            <div class="detail-row">
              <span>新到期时间</span>
              <span class="new-expire">{{ rechargeResult.new_expire_date }}</span>
            </div>
            <div class="detail-row">
              <span>订单号</span>
              <span>{{ rechargeResult.trade_no }}</span>
            </div>
          </div>
          <van-button type="primary" plain block size="small" @click="onSearchAgain" style="margin-top: 12px;">
            重新查询
          </van-button>
        </div>
      </div>

      <!-- 充值按钮 -->
      <div v-else class="action-section">
        <van-button
          type="primary"
          block
          round
          size="large"
          @click="onRecharge"
          :loading="paying"
          loading-text="支付中..."
        >
          立即充值
        </van-button>
      </div>
    </div>
  </div>
</template>

<script>
import { queryCard } from '@/api/card';
import { createRechargeOrder, queryPaymentStatus } from '@/api/payment';
import { closeToast, showDialog, showLoadingToast, showToast } from 'vant';
import { ref } from 'vue';

export default {
  name: 'Query',
  setup() {
    const keyword = ref('')
    const cardInfo = ref(null)
    const searching = ref(false)
    const paying = ref(false)
    const searchError = ref('')
    const rechargeSuccess = ref(false)
    const rechargeResult = ref({})

    const getStatusType = (status) => {
      const types = { 1: 'success', 2: 'warning', 3: 'danger' }
      return types[status] || 'default'
    }

    // 查询卡片
    const onSearch = async () => {
      if (!keyword.value.trim()) {
        showToast('请输入设备号或SIM卡号')
        return
      }

      searching.value = true
      searchError.value = ''
      cardInfo.value = null
      rechargeSuccess.value = false

      try {
        const data = await queryCard(keyword.value.trim())
        cardInfo.value = data
      } catch (error) {
        searchError.value = '未找到对应的卡片信息，请检查输入'
      } finally {
        searching.value = false
      }
    }

    // 重新查询（充值成功后）
    const onSearchAgain = async () => {
      rechargeSuccess.value = false
      rechargeResult.value = {}
      await onSearch()
    }

    // 轮询支付状态
    const pollPaymentStatus = (tradeNo, maxRetries = 30) => {
      return new Promise((resolve, reject) => {
        let retries = 0
        const timer = setInterval(async () => {
          retries++
          try {
            const data = await queryPaymentStatus(tradeNo)
            if (data.payment_status === 1) {
              clearInterval(timer)
              resolve(data)
            } else if (data.payment_status >= 2 || retries >= maxRetries) {
              clearInterval(timer)
              reject(new Error(retries >= maxRetries ? '支付超时，请稍后查询' : '支付失败'))
            }
          } catch (err) {
            if (retries >= maxRetries) {
              clearInterval(timer)
              reject(new Error('查询支付状态失败'))
            }
          }
        }, 2000)
      })
    }

    // 充值
    const onRecharge = async () => {
      try {
        await showDialog({
          title: '确认充值',
          message: `确认为卡号 ${cardInfo.value.card_no} 进行充值续费？`,
          confirmButtonText: '确认支付',
          cancelButtonText: '取消',
          showCancelButton: true
        })
      } catch {
        return // 用户取消
      }

      paying.value = true

      try {
        showLoadingToast({ message: '正在创建订单...', forbidClick: true, duration: 0 })

        const openid = localStorage.getItem('wechat_openid') || 'test_openid'
        const orderData = await createRechargeOrder({
          card_no: cardInfo.value.card_no,
          openid
        })

        closeToast()

        // 实际项目中这里调用微信JSAPI支付
        // 简化版：模拟支付成功并轮询状态
        showLoadingToast({ message: '支付处理中...', forbidClick: true, duration: 0 })

        try {
          const statusData = await pollPaymentStatus(orderData.trade_no)
          closeToast()

          // 重新查询最新卡片信息
          const updatedCard = await queryCard(cardInfo.value.card_no)
          
          rechargeResult.value = {
            amount: orderData.amount,
            trade_no: orderData.trade_no,
            new_expire_date: updatedCard.expire_date
          }
          cardInfo.value = updatedCard
          rechargeSuccess.value = true

        } catch (pollError) {
          closeToast()
          // 支付可能成功但轮询超时，也重新查询一次
          try {
            const updatedCard = await queryCard(cardInfo.value.card_no)
            if (updatedCard.expire_date !== cardInfo.value.expire_date) {
              rechargeResult.value = {
                amount: orderData.amount,
                trade_no: orderData.trade_no,
                new_expire_date: updatedCard.expire_date
              }
              cardInfo.value = updatedCard
              rechargeSuccess.value = true
            } else {
              showToast(pollError.message || '支付未完成，请稍后重试')
            }
          } catch {
            showToast('订单已创建，订单号：' + orderData.trade_no + '，请稍后查询结果')
          }
        }
      } catch (error) {
        closeToast()
        showToast(error.message || '创建订单失败，请重试')
      } finally {
        paying.value = false
      }
    }

    return {
      keyword,
      cardInfo,
      searching,
      paying,
      searchError,
      rechargeSuccess,
      rechargeResult,
      onSearch,
      onSearchAgain,
      onRecharge,
      getStatusType
    }
  }
}
</script>

<style scoped>
.query-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding-bottom: 30px;
}

.search-section {
  padding: 20px 16px 10px;
}

.search-section :deep(.van-search) {
  border-radius: 8px;
  overflow: hidden;
}

.empty-tips {
  padding-top: 80px;
}

.empty-tips :deep(.van-empty__description) {
  color: rgba(255, 255, 255, 0.8);
}

.empty-tips :deep(.van-empty__image) {
  opacity: 0.7;
}

.result-section {
  padding: 0 16px;
  margin-top: 10px;
}

.card-result {
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e8f0 100%);
  border-bottom: 1px solid #ebedf0;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #323233;
}

.card-body {
  padding: 8px 0;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 20px;
  border-bottom: 1px solid #f7f8fa;
}

.info-row:last-child {
  border-bottom: none;
}

.info-label {
  color: #969799;
  font-size: 14px;
  flex-shrink: 0;
}

.info-value {
  color: #323233;
  font-size: 14px;
  text-align: right;
  word-break: break-all;
}

.days-tag {
  font-size: 12px;
  color: #969799;
}

.text-danger {
  color: #ee0a24 !important;
}

.action-section {
  margin-top: 24px;
  padding: 0 4px;
}

.success-section {
  margin-top: 16px;
}

.success-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px 20px;
  text-align: center;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.success-icon {
  font-size: 56px;
  color: #07c160;
}

.success-title {
  font-size: 20px;
  font-weight: 600;
  color: #323233;
  margin-top: 12px;
  margin-bottom: 20px;
}

.success-detail {
  background: #f7f8fa;
  border-radius: 8px;
  padding: 12px 16px;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  font-size: 14px;
  color: #646566;
}

.detail-row .amount {
  color: #ee0a24;
  font-weight: 600;
  font-size: 16px;
}

.detail-row .new-expire {
  color: #07c160;
  font-weight: 600;
}
</style>
