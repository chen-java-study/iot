<template>
  <div class="query-page">
    <van-nav-bar title="物联网卡查询" />

    <div class="search-section">
      <van-search
        v-model="keyword"
        show-action
        placeholder="请输入卡号或设备号"
        @search="onSearch"
      >
        <template #action>
          <van-button size="small" type="primary" @click="onSearch">查询</van-button>
        </template>
      </van-search>
    </div>

    <div v-if="!cardInfo" class="empty-tips">
      <van-empty description="请输入卡号或设备号查询" />
    </div>

    <div v-else class="card-info">
      <van-cell-group>
        <van-cell title="卡号" :value="cardInfo.card_no" />
        <van-cell title="设备号" :value="cardInfo.device_no" />
        <van-cell title="运营商" :value="cardInfo.operator" />
        <van-cell title="到期时间" :value="cardInfo.expire_date" />
        <van-cell title="剩余天数" :value="cardInfo.days_remaining + '天'" />
        <van-cell title="状态">
          <template #value>
            <van-tag :type="getStatusType(cardInfo.status)">{{ cardInfo.status_text }}</van-tag>
          </template>
        </van-cell>
      </van-cell-group>

      <div class="action-section">
        <van-button type="primary" block @click="goToRecharge">充值续费</van-button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { queryCard } from '@/api/card'
import { showToast, showLoadingToast, closeToast } from 'vant'

export default {
  name: 'Query',
  setup() {
    const router = useRouter()
    const keyword = ref('')
    const cardInfo = ref(null)

    const getStatusType = (status) => {
      const types = { 1: 'success', 2: 'warning', 3: 'danger' }
      return types[status] || 'default'
    }

    const onSearch = async () => {
      if (!keyword.value.trim()) {
        showToast('请输入卡号或设备号')
        return
      }

      showLoadingToast({ message: '查询中...', forbidClick: true })

      try {
        const data = await queryCard(keyword.value.trim())
        cardInfo.value = data
        closeToast()
      } catch (error) {
        closeToast()
      }
    }

    const goToRecharge = () => {
      router.push({
        name: 'Recharge',
        query: { cardInfo: JSON.stringify(cardInfo.value) }
      })
    }

    return {
      keyword,
      cardInfo,
      onSearch,
      getStatusType,
      goToRecharge
    }
  }
}
</script>

<style scoped>
.query-page {
  min-height: 100vh;
  background: #f7f8fa;
}

.search-section {
  padding: 16px;
  background: white;
}

.empty-tips {
  padding-top: 100px;
}

.card-info {
  margin-top: 16px;
}

.action-section {
  padding: 16px;
}
</style>
