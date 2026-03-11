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

      <el-divider content-position="left">修改密码</el-divider>
      <el-form-item label="原密码">
        <el-input v-model="passwordForm.old_password" type="password" show-password style="width: 250px" placeholder="请输入原密码" />
      </el-form-item>
      <el-form-item label="新密码">
        <el-input v-model="passwordForm.new_password" type="password" show-password style="width: 250px" placeholder="请输入新密码" />
      </el-form-item>
      <el-form-item label="确认新密码">
        <el-input v-model="passwordForm.confirm_password" type="password" show-password style="width: 250px" placeholder="请再次输入新密码" />
      </el-form-item>
      <el-form-item>
        <el-button type="warning" @click="handleChangePassword" :loading="passwordLoading">修改密码</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getConfig, updateConfig, changePassword } from '@/api/admin'
import { useRouter } from 'vue-router'

export default {
  name: 'SystemConfig',
  setup() {
    const router = useRouter()
    const loading = ref(false)
    const saving = ref(false)
    const passwordLoading = ref(false)

    const config = reactive({
      alert_days: 30
    })

    const passwordForm = reactive({
      old_password: '',
      new_password: '',
      confirm_password: ''
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

    const handleChangePassword = async () => {
      if (!passwordForm.old_password) {
        ElMessage.warning('请输入原密码')
        return
      }
      if (!passwordForm.new_password) {
        ElMessage.warning('请输入新密码')
        return
      }
      if (passwordForm.new_password.length < 6) {
        ElMessage.warning('新密码长度不能少于6位')
        return
      }
      if (passwordForm.new_password !== passwordForm.confirm_password) {
        ElMessage.warning('两次输入的新密码不一致')
        return
      }

      passwordLoading.value = true
      try {
        await changePassword({
          old_password: passwordForm.old_password,
          new_password: passwordForm.new_password
        })
        ElMessage.success('密码修改成功，请使用新密码重新登录')

        // 清除 token，强制重新登录
        localStorage.removeItem('admin_token')
        localStorage.removeItem('admin_user')
        router.push('/login')
      } catch (error) {
        console.error('修改密码失败:', error)
        ElMessage.error(error?.message || '修改密码失败，原密码可能错误')
      } finally {
        passwordLoading.value = false
      }
    }

    onMounted(() => {
      loadConfig()
    })

    return {
      loading,
      saving,
      passwordLoading,
      config,
      passwordForm,
      loadConfig,
      handleSave,
      handleChangePassword
    }
  }
}
</script>
