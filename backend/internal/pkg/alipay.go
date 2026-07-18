package pkg

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/sheet-platform/backend/internal/config"
)

// alipaySign 对参数做 RSA2（SHA256WithRSA）签名
func alipaySign(params map[string]string, privateKeyB64 string) (string, error) {
	// 1. 去掉 sign 字段，按 key 字母序排序
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "sign" || k == "sign_type" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 2. 拼接 key=value&key=value（不 URL 编码）
	var parts []string
	for _, k := range keys {
		v := params[k]
		if v == "" {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	signStr := strings.Join(parts, "&")

	// 3. 解析 PKCS#8 私钥（base64 编码）
	der, err := base64.StdEncoding.DecodeString(privateKeyB64)
	if err != nil {
		return "", fmt.Errorf("解码私钥失败: %w", err)
	}
	key, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		return "", fmt.Errorf("解析私钥失败: %w", err)
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("私钥不是 RSA 类型")
	}

	// 4. SHA256 哈希 → RSA 签名
	hash := sha256.Sum256([]byte(signStr))
	sig, err := rsa.SignPKCS1v15(rand.Reader, rsaKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("签名失败: %w", err)
	}

	return base64.StdEncoding.EncodeToString(sig), nil
}

// BuildPayForm 生成支付宝电脑网站支付表单 HTML
func BuildPayForm(cfg *config.AlipayConfig, outTradeNo, totalAmount, subject, body, sheetID string) (string, error) {
	bizContent := fmt.Sprintf(
		`{"out_trade_no":"%s","total_amount":"%s","subject":"%s","body":"%s","product_code":"FAST_INSTANT_TRADE_PAY"}`,
		outTradeNo, totalAmount, subject, body,
	)

	returnUrl := strings.ReplaceAll(cfg.ReturnUrl, "{id}", sheetID)

	params := map[string]string{
		"app_id":      cfg.AppID,
		"method":      "alipay.trade.page.pay",
		"charset":     cfg.Charset,
		"sign_type":   cfg.SignType,
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"version":     "1.0",
		"notify_url":  cfg.NotifyUrl,
		"return_url":  returnUrl,
		"biz_content": bizContent,
	}

	sign, err := alipaySign(params, cfg.AppPrivateKey)
	if err != nil {
		return "", err
	}
	params["sign"] = sign

	// 构建自动提交的 HTML 表单
	var formFields strings.Builder
	for k, v := range params {
		formFields.WriteString(fmt.Sprintf(
			`<input type="hidden" name="%s" value="%s" />`,
			k, escapeHTML(v),
		))
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="utf-8" /><title>跳转支付...</title></head>
<body onload="document.forms[0].submit()">
<form action="%s" method="post">%s<input type="submit" value="正在跳转..." style="display:none" /></form>
</body>
</html>`, cfg.GatewayUrl, formFields.String())

	return html, nil
}

func escapeHTML(s string) string {
	escaped := url.QueryEscape(s)
	escaped = strings.ReplaceAll(escaped, "+", "%20")
	return escaped
}
