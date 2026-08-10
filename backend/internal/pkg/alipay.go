package pkg

import (
	"errors"

	"github.com/smartwalle/alipay/v3"

	"github.com/sheet-platform/backend/internal/config"
)

func NewAlipayClient(cfg *config.AlipayConfig) (*alipay.Client, error) {
	if cfg.AppPrivateKey == "" {
		return nil, errors.New("支付宝应用私钥未配置")
	}
	if cfg.AlipayPublicKey == "" {
		return nil, errors.New("支付宝公钥未配置")
	}

	client, err := alipay.New(cfg.AppID, cfg.AppPrivateKey, !cfg.Sandbox)
	if err != nil {
		return nil, err
	}

	if err := client.LoadAliPayPublicKey(cfg.AlipayPublicKey); err != nil {
		return nil, err
	}

	return client, nil
}
