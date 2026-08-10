package pkg

import (
	"fmt"
	"strings"

	"github.com/smartwalle/alipay/v3"

	"github.com/sheet-platform/backend/internal/config"
)

func NewAlipayClient(cfg *config.AlipayConfig) (*alipay.Client, error) {
	privateKey := ensurePEM(cfg.AppPrivateKey, "PRIVATE KEY")
	publicKey := ensurePEM(cfg.AlipayPublicKey, "PUBLIC KEY")

	client, err := alipay.New(cfg.AppID, privateKey, !cfg.Sandbox)
	if err != nil {
		return nil, fmt.Errorf("初始化支付宝客户端失败: %w", err)
	}

	if err := client.LoadAliPayPublicKey(publicKey); err != nil {
		return nil, fmt.Errorf("加载支付宝公钥失败: %w", err)
	}

	return client, nil
}

func ensurePEM(key, label string) string {
	if strings.Contains(key, "-----BEGIN") {
		return key
	}
	return fmt.Sprintf("-----BEGIN %s-----\n%s\n-----END %s-----", label, key, label)
}
