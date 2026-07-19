package scheduler

import (
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/sheet-platform/backend/internal/pkg"
)

type task struct {
	name     string
	interval time.Duration
	run      func(*gorm.DB)
	runFirst bool // 启动时立即执行一次
}

var tasks []task

func register(name string, interval time.Duration, run func(*gorm.DB), runFirst bool) {
	tasks = append(tasks, task{name: name, interval: interval, run: run, runFirst: runFirst})
}

func init() {
	register("archive_pv", 1*time.Hour, pkg.ArchivePV, true)
	// 后续定时任务在这里加一行即可：
	// register("cleanup_expired_codes", 5*time.Minute, pkg.CleanExpiredCodes, true)
}

// Start 启动所有注册的定时任务，返回一个 stop 函数用于优雅停止
func Start(db *gorm.DB, logger *zap.Logger) (stop func()) {
	done := make([]chan struct{}, len(tasks))

	for i, t := range tasks {
		done[i] = make(chan struct{})
		t := t
		ch := done[i]

		go func() {
			if t.runFirst {
				t.run(db)
			}
			ticker := time.NewTicker(t.interval)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					t.run(db)
				case <-ch:
					return
				}
			}
		}()

		logger.Info("scheduler: task started", zap.String("name", t.name), zap.Duration("interval", t.interval))
	}

	return func() {
		for _, ch := range done {
			close(ch)
		}
		logger.Info("scheduler: all tasks stopped")
	}
}
