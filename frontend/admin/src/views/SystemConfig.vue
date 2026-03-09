<template>
  <el-card v-loading="loading">
    <el-form label-width="180px">
      <el-divider content-position="left">充值设置</el-divider>
      <el-form-item label="到期提醒天数">
        <el-input-number v-model.number="config.alert_days" :min="1" :max="365" />
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
      alert_days: 30
    })

    const loadConfig = async () => {
      loading.value = true
      try {
        const data = await getConfig()
        Object.assign(config, {
          alert_days: parseInt(data.alert_days) || 30
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
        await updateConfig({
          alert_days: String(config.alert_days)
        })
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
