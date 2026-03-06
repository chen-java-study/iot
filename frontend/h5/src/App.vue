<template>
  <router-view />
</template>

<script>
import { isWechatBrowser, handleAuthCode, getOpenId, setOpenId, redirectToAuth } from './utils/wechatAuth'
import { getOpenIdByCode } from './api/payment'
import { showLoadingToast, closeToast } from 'vant'

export default {
  name: 'App',
  mounted() {
    this.handleWechatAuth()
  },
  methods: {
    async handleWechatAuth() {
      // 只在微信浏览器中处理
      if (!isWechatBrowser()) {
        return
      }

      // 检查是否已有 openid
      if (getOpenId()) {
        return
      }

      // 检查 URL 中是否有 code
      const code = handleAuthCode()
      if (code) {
        // 用 code 换 openid
        try {
          showLoadingToast({ message: '正在授权...', forbidClick: true, duration: 0 })
          const data = await getOpenIdByCode(code)
          closeToast()
          if (data && data.openid) {
            setOpenId(data.openid)
          }
        } catch (e) {
          closeToast()
          console.error('获取openid失败:', e)
        }
      } else {
        // 没有 code，也没有 openid，跳转授权
        redirectToAuth()
      }
    }
  }
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
  background-color: #f7f8fa;
}

#app {
  min-height: 100vh;
}
</style>
