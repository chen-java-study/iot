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
      <div class="tips">
        默认账号: admin / admin123
        <el-button link type="primary" @click="showPasswordDialog" style="margin-left: 10px;">
          修改密码
        </el-button>
      </div>
    </el-card>

    <!-- 修改密码对话框 -->
    <el-dialog v-model="passwordDialogVisible" title="修改密码" width="400px">
      <el-form :model="passwordForm" label-width="80px">
        <el-form-item label="用户名">
          <el-input v-model="passwordForm.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="原密码">
          <el-input v-model="passwordForm.old_password" type="password" placeholder="请输入原密码" show-password />
        </el-form-item>
        <el-form-item label="新密码">
          <el-input v-model="passwordForm.new_password" type="password" placeholder="请输入新密码" show-password />
        </el-form-item>
        <el-form-item label="确认密码">
          <el-input v-model="passwordForm.confirm_password" type="password" placeholder="请再次输入新密码" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="passwordDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleChangePassword" :loading="passwordLoading">确认修改</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { login, changePassword } from '@/api/admin';
import { ElMessage } from 'element-plus';
import { ref } from 'vue';
import { useRouter } from 'vue-router';

export default {
  name: 'Login',
  setup() {
    const router = useRouter()
    const form = ref({ username: 'admin', password: 'admin123' })
    const loading = ref(false)

    // 修改密码相关
    const passwordDialogVisible = ref(false)
    const passwordLoading = ref(false)
    const passwordForm = ref({
      username: '',
      old_password: '',
      new_password: '',
      confirm_password: ''
    })

    const handleLogin = async () => {
      loading.value = true
      try {
        const data = await login(form.value)
        
        localStorage.setItem('admin_token', data.token)
        if (data.user_info) {
          try {
            localStorage.setItem('admin_user', JSON.stringify(data.user_info))
          } catch (e) {
            console.error('序列化user_info失败:', e)
          }
        }
        ElMessage.success('登录成功')
        router.push('/')
      } catch (error) {
        console.error('登录失败:', error)
        ElMessage.error(error?.message || '登录失败')
      } finally {
        loading.value = false
      }
    }

    const showPasswordDialog = () => {
      passwordForm.value = {
        username: '',
        old_password: '',
        new_password: '',
        confirm_password: ''
      }
      passwordDialogVisible.value = true
    }

    const handleChangePassword = async () => {
      if (!passwordForm.value.username) {
        ElMessage.warning('请输入用户名')
        return
      }
      if (!passwordForm.value.old_password) {
        ElMessage.warning('请输入原密码')
        return
      }
      if (!passwordForm.value.new_password) {
        ElMessage.warning('请输入新密码')
        return
      }
      if (passwordForm.value.new_password.length < 6) {
        ElMessage.warning('新密码长度不能少于6位')
        return
      }
      if (passwordForm.value.new_password !== passwordForm.value.confirm_password) {
        ElMessage.warning('两次输入的新密码不一致')
        return
      }

      passwordLoading.value = true
      try {
        // 先登录验证用户名和原密码
        await login({
          username: passwordForm.value.username,
          password: passwordForm.value.old_password
        })

        // 登录成功，修改密码
        await changePassword({
          old_password: passwordForm.value.old_password,
          new_password: passwordForm.value.new_password
        })

        ElMessage.success('密码修改成功，请使用新密码登录')
        passwordDialogVisible.value = false
      } catch (error) {
        console.error('修改密码失败:', error)
        ElMessage.error(error?.message || '修改密码失败，请检查用户名和原密码')
      } finally {
        passwordLoading.value = false
      }
    }

    return {
      form,
      loading,
      handleLogin,
      passwordDialogVisible,
      passwordForm,
      passwordLoading,
      showPasswordDialog,
      handleChangePassword
    }
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
