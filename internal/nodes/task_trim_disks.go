// Copyright 2024 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package nodes

import (
	"fmt"
	"github.com/TeaOSLab/EdgeNode/internal/remotelogs"
	executils "github.com/TeaOSLab/EdgeNode/internal/utils/exec"
	"github.com/shirou/gopsutil/v3/load"
	"runtime"
	"time"
)

// TrimDisksTask trim ssd disks automatically
type TrimDisksTask struct {
}

// NewTrimDisksTask create new task
func NewTrimDisksTask() *TrimDisksTask {
	return &TrimDisksTask{}
}

// Start the task
func (this *TrimDisksTask) Start() {
	// execute once
	err := this.loop()
	if err != nil {
		remotelogs.Warn("TRIM_DISKS", "trim disks failed: "+err.Error())
	}

	var ticker = time.NewTicker(2 * 24 * time.Hour) // every 2 days
	for range ticker.C {
		// prevent system overload
		for i := 0; i < 24; i++ {
			stat, loadErr := load.Avg()
			if loadErr == nil && stat != nil && stat.Load1 > 15 {
				// wait load downgrade
				time.Sleep(1 * time.Hour)
			} else {
				break
			}
		}

		// run the task
		err = this.loop()
		if err != nil {
			remotelogs.Warn("TRIM_DISKS", "trim disks failed: "+err.Error())
		}
	}
}

// run the task once
func (this *TrimDisksTask) loop() error {
	if runtime.GOOS != "linux" {
		return nil
	}

	var nodeConfig = sharedNodeConfig
	if nodeConfig == nil {
		return nil
	}
	if !nodeConfig.AutoTrimDisks {
		return nil
	}

	trimExe, err := executils.LookPath("fstrim")
	if err != nil {
		return fmt.Errorf("'fstrim' command not found: %w", err)
	}

	var cmd = executils.NewCmd(trimExe, "-a").
		WithStderr()
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("'fstrim' execute failed: %s", cmd.Stderr())
	}

	return nil
}