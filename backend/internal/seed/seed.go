package seed

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/sheet-platform/backend/internal/model"
)

// SystemConfigs 插入默认系统配置，已存在的 key 会被跳过
func SystemConfigs(db *gorm.DB) error {
	defaults := []model.SystemConfig{
		{ConfigKey: "site_name", ConfigVal: "乐谱平台", Remark: "网站名称"},
		{ConfigKey: "site_description", ConfigVal: "专业的乐谱分享平台", Remark: "网站描述"},
		{ConfigKey: "max_upload_size", ConfigVal: "50", Remark: "最大上传文件大小(MB)"},
		{ConfigKey: "allowed_file_types", ConfigVal: "pdf,mp3,wav", Remark: "允许上传的文件类型"},
		{ConfigKey: "download_points_default", ConfigVal: "10", Remark: "默认下载所需积分"},
		{ConfigKey: "register_points", ConfigVal: "100", Remark: "注册赠送积分"},
		{ConfigKey: "maintenance_mode", ConfigVal: "false", Remark: "维护模式开关"},
		{ConfigKey: "contact_email", ConfigVal: "admin@example.com", Remark: "联系邮箱"},
		{ConfigKey: "about_us", ConfigVal: "<h2>关于我们</h2><p>我们是一个专业的<strong>乐谱分享平台</strong>，致力于为音乐爱好者提供优质的乐谱资源。</p><h3>我们的使命</h3><ul><li>让音乐学习变得更加简单和有趣</li><li>为音乐爱好者提供丰富的乐谱资源</li><li>促进音乐文化的传播和交流</li></ul><h3>平台特色</h3><p>我们提供：</p><ul><li>高质量的PDF乐谱文件</li><li>音频试听功能</li><li>积分下载系统</li><li>用户交流社区</li></ul>", Remark: "关于我们"},
		{ConfigKey: "privacy_policy", ConfigVal: "<h2>隐私政策</h2><p>我们重视您的隐私保护...</p>", Remark: "隐私政策"},
		{ConfigKey: "terms_of_service", ConfigVal: "<h2>服务条款</h2><p>使用本平台即表示您同意以下条款...</p>", Remark: "服务条款"},
	}

	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "config_key"}},
		DoNothing: true,
	}).Create(&defaults).Error
}
