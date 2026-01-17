<template>
  <el-container>
    <el-aside width="200px">
      <div class="logo">IoT卡管理</div>
      <el-menu :default-active="$route.path" router>
        <el-menu-item index="/">
          <span>首页仪表盘</span>
        </el-menu-item>
        <el-menu-item index="/cards">
          <span>卡片管理</span>
        </el-menu-item>
        <el-menu-item index="/recharges">
          <span>充值记录</span>
        </el-menu-item>
        <el-menu-item index="/config">
          <span>系统配置</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header>
        <div class="header-right">
          <span>欢迎，管理员</span>
          <el-button type="text" @click="handleLogout">退出</el-button>
        </div>
      </el-header>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
import { useRouter } from 'vue-router'

export default {
  name: 'Layout',
  setup() {
    const router = useRouter()

    const handleLogout = () => {
      localStorage.removeItem('admin_token')
      localStorage.removeItem('admin_user')
      router.push('/login')
    }

    return { handleLogout }
  }
}
</script>

<style scoped>
.el-aside {
  background: #304156;
  color: #fff;
}
.logo {
  height: 60px;
  line-height: 60px;
  text-align: center;
  font-size: 20px;
  font-weight: bold;
  color: #fff;
  background: #2a3f54;
}
.el-header {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0,21,41,.08);
}
.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}
.el-main {
  background: #f0f2f5;
}
</style>
