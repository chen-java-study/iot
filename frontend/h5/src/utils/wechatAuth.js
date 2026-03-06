const APP_ID = 'wxd7fd41032d8269a9'
const REDIRECT_URI = encodeURIComponent('https://iot4you.top/h5/')

export function getOpenId() {
  return localStorage.getItem('wechat_openid')
}

export function setOpenId(openid) {
  localStorage.setItem('wechat_openid', openid)
}

export function isWechatBrowser() {
  const ua = navigator.userAgent.toLowerCase()
  return ua.indexOf('micromessenger') > -1
}

export function redirectToAuth() {
  const authUrl = `https://open.weixin.qq.com/connect/oauth2/authorize?appid=${APP_ID}&redirect_uri=${REDIRECT_URI}&response_type=code&scope=snsapi_base#wechat_redirect`
  window.location.href = authUrl
}

export function handleAuthCode() {
  const url = new URL(window.location.href)
  const code = url.searchParams.get('code')

  if (code) {
    // 清除 URL 中的 code 参数
    url.searchParams.delete('code')
    window.history.replaceState({}, '', url.toString())
    return code
  }
  return null
}
