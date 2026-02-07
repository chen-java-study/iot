<template>
  <div class="recharge-page">
    <van-nav-bar title="充值续费" left-arrow @click-left="goBack" />

    <van-cell-group>
      <van-cell title="卡号" :value="cardInfo.card_no" />
      <van-cell title="设备号" :value="cardInfo.device_no" />
      <van-cell title="当前到期时间" :value="cardInfo.expire_date" />
      <van-cell title="充值年限" value="1年" />
      <van-cell title="充值金额" :value="'¥100.00'" />
    </van-cell-group>

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
        <van-button type="primary" plain block size="small" @click="goBack" style="margin-top: 16px;">
          返回查询
        </van-button>
      </div>
    </div>

    <!-- 支付按钮 -->
    <div v-else class="pay-section">
      <van-button type="primary" block round @click="onPay" :loading="paying" loading-text="支付中...">
        立即支付
      </van-button>
    </div>
  </div>
</template>

<script>
import { queryCard } from '@/api/card'
import { createRechargeOrder, queryPaymentStatus } from '@/api/payment'
import { closeToast, showDialog, showLoadingToast, showToast } from 'vant'
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

export default {
  name: 'Recharge',
  setup() {
    const route = useRoute()
    const router = useRouter()
    const cardInfo = ref(JSON.parse(route.query.cardInfo || '{}'))
    const paying = ref(false)
    const rechargeSuccess = ref(false)
    const rechargeResult = ref({})

    const goBack = () => {
      router.push({ name: 'Query' })
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
              reject(new Error(retries >= maxRetries ? '支付超时' : '支付失败'))
            }
          } catch (err) {
            if (retries >= maxRetries) {
              clearInterval(timer)
              reject(new Error('查询失败'))
            }
          }
        }, 2000)
      })
    }

    const onPay = async () => {
      try {
        await showDialog({
          title: '确认支付',
          message: `确认支付 ¥100.00 为卡号 ${cardInfo.value.card_no} 续费？`,
          confirmButtonText: '确认',
          cancelButtonText: '取消',
          showCancelButton: true
        })
      } catch {
        return
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
        showLoadingToast({ message: '支付处理中...', forbidClick: true, duration: 0 })

        try {
          await pollPaymentStatus(orderData.trade_no)
          closeToast()

          // 重新查询最新卡片信息获取新到期时间
          const updatedCard = await queryCard(cardInfo.value.card_no)

          rechargeResult.value = {
            amount: orderData.amount,
            trade_no: orderData.trade_no,
            new_expire_date: updatedCard.expire_date
          }
          cardInfo.value = updatedCard
          rechargeSuccess.value = true
        } catch (pollErr) {
          closeToast()
          showToast(pollErr.message || '支付未完成，请稍后查询')
        }
      } catch (error) {
        closeToast()
        showToast(error.message || '创建订单失败')
      } finally {
        paying.value = false
      }
    }

    return {
      cardInfo,
      paying,
      rechargeSuccess,
      rechargeResult,
      onPay,
      goBack
    }
  }
}
</script>

<style scoped>
.recharge-page {
  min-height: 100vh;
  background: #f7f8fa;
}

.pay-section {
  padding: 24px 16px;
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: white;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.05);
}

.success-section {
  padding: 20px 16px;
}

.success-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px 20px;
  text-align: center;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.success-icon {
  font-size: 56px;
  color: #07c160;
}

.success-title {
  font-size: 20px;
  font-weight: 600;
  color: #323233;
  margin: 12px 0 20px;
}

.success-detail {
  background: #f7f8fa;
  border-radius: 8px;
  padding: 12px 16px;
  text-align: left;
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
