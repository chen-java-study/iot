<template>
  <div>
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card><el-statistic title="卡片总数" :value="stats.total_cards || 0" /></el-card>
      </el-col>
      <el-col :span="6">
        <el-card><el-statistic title="充值总金额" :value="stats.total_recharge_amount || 0" prefix="¥" /></el-card>
      </el-col>
      <el-col :span="6">
        <el-card><el-statistic title="今日收入" :value="stats.today_amount || 0" prefix="¥" /></el-card>
      </el-col>
      <el-col :span="6">
        <el-card><el-statistic title="本月收入" :value="stats.month_amount || 0" prefix="¥" /></el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { getStatistics } from '@/api/admin'

export default {
  name: 'Dashboard',
  setup() {
    const stats = ref({})

    const loadStats = async () => {
      try {
        const data = await getStatistics()
        stats.value = data
      } catch (error) {
        console.error('加载统计失败:', error)
      }
    }

    onMounted(() => {
      loadStats()
    })

    return { stats }
  }
}
</script>
