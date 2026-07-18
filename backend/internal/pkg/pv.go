package pkg

import (
	"sync/atomic"
	"time"

	"gorm.io/gorm"

	"github.com/sheet-platform/backend/internal/model"
)

var todayPV atomic.Int64

func IncrementPV() {
	todayPV.Add(1)
}

func TodayPVCount() int64 {
	return todayPV.Load()
}

// ArchivePV 每小时将内存中的 PV 归档到 site_pv_history 表
func ArchivePV(db *gorm.DB) {
	count := todayPV.Swap(0)
	if count == 0 {
		return
	}

	today := time.Now().Format("2006-01-02")
	db.Where("date = ?", today).FirstOrCreate(&model.SitePvHistory{Date: today})
	db.Model(&model.SitePvHistory{}).
		Where("date = ?", today).
		UpdateColumn("page_views", gorm.Expr("page_views + ?", count))
}

// GetPvStats 返回总 PV 和今日 PV
func GetPvStats(db *gorm.DB) (totalPV, todayPVCount int64) {
	var total int64
	db.Model(&model.SitePvHistory{}).Select("COALESCE(SUM(page_views), 0)").Scan(&total)

	today := time.Now().Format("2006-01-02")
	var todayRecord model.SitePvHistory
	db.Where("date = ?", today).First(&todayRecord)

	todayTotal := todayRecord.PageViews + TodayPVCount()
	return total + TodayPVCount(), todayTotal
}
