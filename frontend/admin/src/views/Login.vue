<template>
  <div class="login-container">
    <el-card class="login-box">
      <h2>物联网卡管理系统</h2>
      <el-form :model="form" @submit.prevent="handleLogin">
        <el-form-item>
          <el-input v-model="form.username" placeholder="用户名" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="form.password" type="password" placeholder="密码" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleLogin" :loading="loading" style="width: 100%">
            登录
          </el-button>
        </el-form-item>
      </el-form>
      <div class="tips">默认账号: admin / admin123</div>
    </el-card>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { login } from '@/api/admin'

export default {
  name: 'Login',
  setup() {
    const router = useRouter()
    const form = ref({ username: 'admin', password: 'admin123' })
    const loading = ref(false)

    const handleLogin = async () => {
      loading.value = true
      try {
        const data = await login(form.value)
        localStorage.setItem('admin_token', data.token)
        localStorage.setItem('admin_user', JSON.stringify(data.user_info))
        ElMessage.success('登录成功')
        router.push('/')
      } catch (error) {
        console.error('登录失败:', error)
      } finally {
        loading.value = false
      }
    }

    return { form, loading, handleLogin }
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}
.login-box {
  width: 400px;
  padding: 40px;
}
h2 {
  text-align: center;
  margin-bottom: 30px;
}
.tips {
  text-align: center;
  color: #999;
  font-size: 12px;
}
</style>
