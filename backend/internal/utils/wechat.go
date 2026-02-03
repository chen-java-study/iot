package utils

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// WechatPayValidator 微信支付签名验证器
type WechatPayValidator struct {
	mchID      string
	apiV3Key   string
	privateKey *rsa.PrivateKey
}

// NewWechatPayValidator 创建微信支付验证器
func NewWechatPayValidator(mchID, apiV3Key, privateKeyPath string) (*WechatPayValidator, error) {
	keyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("读取私钥文件失败: %w", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, errors.New("解析PEM格式私钥失败")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析私钥失败: %w", err)
	}

	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("私钥类型错误")
	}

	return &WechatPayValidator{
		mchID:      mchID,
		apiV3Key:   apiV3Key,
		privateKey: rsaKey,
	}, nil
}

// VerifySignature 验证微信支付回调签名
func (v *WechatPayValidator) VerifySignature(timestamp, nonce, body, signature string, wechatPublicKey *rsa.PublicKey) error {
	// 构造验签字符串
	message := fmt.Sprintf("%s\n%s\n%s\n", timestamp, nonce, body)

	// 解码签名
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("解码签名失败: %w", err)
	}

	// 计算消息哈希
	hashed := sha256.Sum256([]byte(message))

	// 验证签名
	err = rsa.VerifyPKCS1v15(wechatPublicKey, crypto.SHA256, hashed[:], signatureBytes)
	if err != nil {
		return fmt.Errorf("签名验证失败: %w", err)
	}

	return nil
}

// DecryptAESGCM 解密微信支付回调数据 (AES-256-GCM)
func (v *WechatPayValidator) DecryptAESGCM(ciphertext, nonce, associatedData string) ([]byte, error) {
	// 这里需要实现 AES-256-GCM 解密
	// 使用 apiV3Key 作为密钥
	// 实际实现需要引入 crypto/aes 和 crypto/cipher 包
	return nil, errors.New("请实现 AES-256-GCM 解密")
}

// ParseNotifyHeaders 从请求头解析验签所需信息
func ParseNotifyHeaders(headers map[string]string) (timestamp, nonce, signature, serialNo string, err error) {
	timestamp = headers["Wechatpay-Timestamp"]
	nonce = headers["Wechatpay-Nonce"]
	signature = headers["Wechatpay-Signature"]
	serialNo = headers["Wechatpay-Serial"]

	if timestamp == "" || nonce == "" || signature == "" || serialNo == "" {
		return "", "", "", "", errors.New("缺少必要的微信支付回调头信息")
	}

	return timestamp, nonce, signature, serialNo, nil
}

// ReadBody 读取请求体
func ReadBody(body io.Reader) (string, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ValidateNotifyRequest 验证微信支付回调请求的完整性
func ValidateNotifyRequest(timestamp, nonce, body string) error {
	// 验证时间戳，防止重放攻击（5分钟内有效）
	// 这里简化处理，实际应该解析时间戳并比较
	if timestamp == "" {
		return errors.New("时间戳为空")
	}

	if nonce == "" {
		return errors.New("随机串为空")
	}

	if strings.TrimSpace(body) == "" {
		return errors.New("请求体为空")
	}

	return nil
}
