<template>
  <div class="detail-page">
    <van-nav-bar title="卡片详情" left-arrow @click-left="$router.back()" />
    <van-cell-group>
      <van-cell title="卡号" :value="cardInfo.card_no" />
      <van-cell title="设备号" :value="cardInfo.device_no" />
    </van-cell-group>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRoute } from 'vue-router'

export default {
  name: 'Detail',
  setup() {
    const route = useRoute()
    const cardInfo = ref(JSON.parse(route.query.cardInfo || '{}'))
    return { cardInfo }
  }
}
</script>
