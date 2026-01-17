<template>
  <div class="recharge-page">
    <van-nav-bar title="充值续费" left-arrow @click-left="$router.back()" />

    <van-cell-group>
      <van-cell title="卡号" :value="cardInfo.card_no" />
      <van-cell title="当前到期时间" :value="cardInfo.expire_date" />
      <van-cell title="充值年限" value="1年" />
      <van-cell title="充值金额" :value="'¥100.00'" />
    </van-cell-group>

    <div class="pay-section">
      <van-button type="primary" block @click="onPay">立即支付</van-button>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { createRechargeOrder } from '@/api/payment'
import { showLoadingToast, closeToast, showSuccessToast, showFailToast } from 'vant'

export default {
  name: 'Recharge',
  setup() {
    const route = useRoute()
    const cardInfo = ref(JSON.parse(route.query.cardInfo || '{}'))

    const onPay = async () => {
      try {
        showLoadingToast({ message: '正在创建订单...', forbidClick: true })

        // 实际项目中需要获取真实的openid
        const openid = localStorage.getItem('wechat_openid') || 'test_openid'

        const data = await createRechargeOrder({
          card_no: cardInfo.value.card_no,
          openid
        })

        closeToast()
        // 这里应该调用微信JSAPI支付
        // 简化处理，实际需要调用 WeixinJSBridge
        showSuccessToast('订单创建成功，订单号：' + data.trade_no)

      } catch (error) {
        closeToast()
        showFailToast('支付失败')
      }
    }

    return {
      cardInfo,
      onPay
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
  padding: 16px;
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: white;
}
</style>
