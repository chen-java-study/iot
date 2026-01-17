<template>
  <el-card v-loading="loading">
    <el-form label-width="180px">
      <el-divider content-position="left">充值设置</el-divider>
      <el-form-item label="充值价格(元/年)">
        <el-input-number v-model.number="config.recharge_price" :min="0" :precision="2" />
      </el-form-item>
      <el-form-item label="到期提醒天数">
        <el-input-number v-model.number="config.alert_days" :min="1" :max="365" />
      </el-form-item>

      <el-divider content-position="left">微信支付配置</el-divider>
      <el-form-item label="AppID">
        <el-input v-model="config.wechat_appid" placeholder="微信公众号或小程序AppID" />
      </el-form-item>
      <el-form-item label="商户号">
        <el-input v-model="config.wechat_mch_id" placeholder="微信支付商户号" />
      </el-form-item>
      <el-form-item label="API密钥">
        <el-input v-model="config.wechat_api_key" type="password" show-password placeholder="微信支付API密钥" />
      </el-form-item>
      <el-form-item label="APIv3密钥">
        <el-input v-model="config.wechat_api_v3_key" type="password" show-password placeholder="微信支付APIv3密钥" />
      </el-form-item>
      <el-form-item label="证书序列号">
        <el-input v-model="config.wechat_serial_no" placeholder="商户API证书序列号" />
      </el-form-item>
      <el-form-item label="支付回调地址">
        <el-input v-model="config.notify_url" placeholder="https://yourdomain.com/api/v1/payment/notify" />
      </el-form-item>

      <el-form-item>
        <el-button type="primary" @click="handleSave" :loading="saving">保存配置</el-button>
        <el-button @click="loadConfig">重置</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getConfig, updateConfig } from '@/api/admin'

export default {
  name: 'SystemConfig',
  setup() {
    const loading = ref(false)
    const saving = ref(false)
    
    const config = reactive({
      recharge_price: 100,
      alert_days: 30,
      wechat_appid: '',
      wechat_mch_id: '',
      wechat_api_key: '',
      wechat_api_v3_key: '',
      wechat_serial_no: '',
      notify_url: ''
    })

    const loadConfig = async () => {
      loading.value = true
      try {
        const data = await getConfig()
        Object.assign(config, {
          recharge_price: parseFloat(data.recharge_price) || 100,
          alert_days: parseInt(data.alert_days) || 30,
          wechat_appid: data.wechat_appid || '',
          wechat_mch_id: data.wechat_mch_id || '',
          wechat_api_key: data.wechat_api_key || '',
          wechat_api_v3_key: data.wechat_api_v3_key || '',
          wechat_serial_no: data.wechat_serial_no || '',
          notify_url: data.notify_url || ''
        })
      } catch (error) {
        console.error('加载配置失败:', error)
      } finally {
        loading.value = false
      }
    }

    const handleSave = async () => {
      saving.value = true
      try {
        await updateConfig(config)
        ElMessage.success('保存成功')
      } catch (error) {
        console.error('保存配置失败:', error)
      } finally {
        saving.value = false
      }
    }

    onMounted(() => {
      loadConfig()
    })

    return {
      loading,
      saving,
      config,
      loadConfig,
      handleSave
    }
  }
}
</script>
