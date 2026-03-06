package utils

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type PayClient struct {
	AppID          string
	MchID          string
	APIV3Key       string
	SerialNo       string
	PrivateKey     *rsa.PrivateKey
	PrivateKeyPath string
	NotifyURL      string
	Client         *http.Client
}

// NewPayClient 创建微信支付客户端
func NewPayClient(appID, mchID, apiV3Key, serialNo, privateKeyPath, notifyURL string) (*PayClient, error) {
	keyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("读取私钥文件失败: %w", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("解析PEM格式私钥失败")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析私钥失败: %w", err)
	}

	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("私钥类型错误")
	}

	return &PayClient{
		AppID:          appID,
		MchID:          mchID,
		APIV3Key:       apiV3Key,
		SerialNo:       serialNo,
		PrivateKey:     rsaKey,
		PrivateKeyPath: privateKeyPath,
		NotifyURL:      notifyURL,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// PrepayRequest 统一下单请求
type PrepayRequest struct {
	AppID       string `json:"appid"`
	MchID       string `json:"mchid"`
	Description string `json:"description"`
	OutTradeNo  string `json:"out_trade_no"`
	NotifyURL   string `json:"notify_url"`
	Amount      Amount `json:"amount"`
	Payer       Payer  `json:"payer"`
}

type Amount struct {
	Total    int    `json:"total"`    // 金额，单位：分
	Currency string `json:"currency"` // 货币类型
}

type Payer struct {
	OpenID string `json:"openid"`
}

// PrepayResponse 统一下单响应
type PrepayResponse struct {
	PrepayID string `json:"prepay_id"`
	// 其他字段...
}

// CreateJSAPIPrepay 创建JSAPI支付订单
func (p *PayClient) CreateJSAPIPrepay(outTradeNo, description, openID string, amount int) (*PrepayResponse, error) {
	// 构建请求
	req := PrepayRequest{
		AppID:       p.AppID,
		MchID:       p.MchID,
		Description: description,
		OutTradeNo:  outTradeNo,
		NotifyURL:   p.NotifyURL,
		Amount: Amount{
			Total:    amount, // 单位：分
			Currency: "CNY",
		},
		Payer: Payer{
			OpenID: openID,
		},
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	// 构建签名
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonceStr := generateNonce()
	signStr := fmt.Sprintf("%s\n%s\n%s\n%s\n", "POST", "/v3/pay/transactions/jsapi", timestamp, nonceStr)
	signStr += string(reqBody) + "\n"

	signature, err := p.sign(signStr)
	if err != nil {
		return nil, fmt.Errorf("生成签名失败: %w", err)
	}

	// 发送请求
	url := "https://api.mch.weixin.qq.com/v3/pay/transactions/jsapi"
	httpReq, err := http.NewRequest("POST", url, strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("WECHATPAY2-SHA256-RSA2048 mchid=\"%s\",nonce_str=\"%s\",timestamp=\"%s\",serial_no=\"%s\",signature=\"%s\"",
		p.MchID, nonceStr, timestamp, p.SerialNo, signature))

	resp, err := p.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("微信支付下单失败: %s, body: %s", resp.Status, string(respBody))
	}

	var prepayResp PrepayResponse
	if err := json.Unmarshal(respBody, &prepayResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &prepayResp, nil
}

// GeneratePayParams 生成前端支付参数
func (p *PayClient) GeneratePayParams(prepayID string) (map[string]interface{}, error) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonceStr := generateNonce()

	// 前端签名串: appId\ntimeStamp\nnonceStr\npackage\n\n
	paySignStr := fmt.Sprintf("%s\n%s\n%s\nprepay_id=%s\n\n", p.AppID, timestamp, nonceStr, prepayID)

	signature, err := p.sign(paySignStr)
	if err != nil {
		return nil, fmt.Errorf("生成支付签名失败: %w", err)
	}

	return map[string]interface{}{
		"appId":     p.AppID,
		"timeStamp": timestamp,
		"nonceStr":  nonceStr,
		"package":   fmt.Sprintf("prepay_id=%s", prepayID),
		"signType":  "RSA",
		"paySign":   signature,
	}, nil
}

// sign 生成RSA签名
func (p *PayClient) sign(message string) (string, error) {
	h := sha256.Sum256([]byte(message))
	signature, err := rsa.SignPKCS1v15(rand.Reader, p.PrivateKey, crypto.SHA256, h[:])
	if err != nil {
		return "", fmt.Errorf("签名失败: %w", err)
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

// generateNonce 生成随机字符串
func generateNonce() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// DecryptNotification 解密微信支付回调通知
func (p *PayClient) DecryptNotification(ciphertext, nonce, associatedData string) ([]byte, error) {
	// 这里简化处理，实际需要使用 AES-256-GCM 解密
	// 使用 APIV3Key 作为密钥
	key := []byte(p.APIV3Key)
	cipherData, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("Base64解码失败: %w", err)
	}

	// AES-GCM 解密
	return decryptAESGCM(key, []byte(nonce), []byte(associatedData), cipherData)
}

// decryptAESGCM AES-256-GCM 解密
func decryptAESGCM(key, nonce, aad, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("创建AES密钥失败: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("创建GCM失败: %w", err)
	}

	// nonce 应该是 12 字节
	if len(nonce) != 12 {
		return nil, fmt.Errorf("nonce 长度错误: %d", len(nonce))
	}

	// ciphertext 包含 auth tag (16 bytes)
	if len(ciphertext) < 16 {
		return nil, fmt.Errorf("ciphertext 长度错误")
	}

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, aad)
	if err != nil {
		return nil, fmt.Errorf("解密失败: %w", err)
	}

	return plaintext, nil
}
