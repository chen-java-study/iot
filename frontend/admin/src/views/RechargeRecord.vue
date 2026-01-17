<template>
  <el-card>
    <div class="toolbar" style="margin-bottom: 20px;">
      <el-input v-model="searchForm.keyword" placeholder="搜索卡号/设备号/订单号" style="width: 300px;" @change="loadRecords" />
      <el-select v-model="searchForm.status" placeholder="支付状态" style="width: 120px; margin-left: 10px;" @change="loadRecords">
        <el-option label="全部" :value="-1" />
        <el-option label="待支付" :value="0" />
        <el-option label="已支付" :value="1" />
        <el-option label="已退款" :value="2" />
      </el-select>
      <el-date-picker
        v-model="dateRange"
        type="daterange"
        range-separator="至"
        start-placeholder="开始日期"
        end-placeholder="结束日期"
        style="margin-left: 10px;"
        value-format="YYYY-MM-DD"
        @change="loadRecords"
      />
    </div>
    
    <el-table :data="records" border v-loading="loading">
      <el-table-column prop="trade_no" label="订单号" width="180" />
      <el-table-column prop="card_no" label="卡号" width="200" />
      <el-table-column prop="device_no" label="设备号" width="180" />
      <el-table-column prop="recharge_amount" label="金额" width="100">
        <template #default="{ row }">
          ¥{{ row.recharge_amount }}
        </template>
      </el-table-column>
      <el-table-column prop="payment_status_text" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.payment_status === 1 ? 'success' : row.payment_status === 0 ? 'warning' : 'info'">
            {{ row.payment_status_text }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="paid_at" label="支付时间" width="180" />
      <el-table-column prop="created_at" label="创建时间" width="180" />
    </el-table>
    
    <el-pagination
      v-model:current-page="searchForm.page"
      v-model:page-size="searchForm.page_size"
      :total="total"
      @current-change="loadRecords"
      style="margin-top: 20px; justify-content: center"
    />
    
    <div style="margin-top: 20px; text-align: right;">
      <el-text>总充值金额: <el-text type="primary" size="large">¥{{ totalAmount }}</el-text></el-text>
    </div>
  </el-card>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { getRechargeList } from '@/api/admin'

export default {
  name: 'RechargeRecord',
  setup() {
    const loading = ref(false)
    const records = ref([])
    const total = ref(0)
    const totalAmount = ref(0)
    const dateRange = ref([])
    
    const searchForm = reactive({
      page: 1,
      page_size: 20,
      status: -1,
      keyword: '',
      start_date: '',
      end_date: ''
    })

    const loadRecords = async () => {
      loading.value = true
      try {
        if (dateRange.value && dateRange.value.length === 2) {
          searchForm.start_date = dateRange.value[0]
          searchForm.end_date = dateRange.value[1]
        } else {
          searchForm.start_date = ''
          searchForm.end_date = ''
        }
        
        const data = await getRechargeList(searchForm)
        records.value = data.list || []
        total.value = data.total || 0
        totalAmount.value = data.total_amount || 0
      } catch (error) {
        console.error('加载充值记录失败:', error)
      } finally {
        loading.value = false
      }
    }

    onMounted(() => {
      loadRecords()
    })

    return {
      loading,
      records,
      total,
      totalAmount,
      dateRange,
      searchForm,
      loadRecords
    }
  }
}
</script>
