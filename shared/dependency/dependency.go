package dependency

import (
	"log"
	"os/exec"
)

// 为目标系统安装依赖

// 适配 debian 11/12 ubnuntu 24.04/22.04
func InstallDepsForDebian(packs []string) string {
	args := []string{"install", "-y"}
	for _, pack := range packs {
		args = append(args, pack)
	}
	cmd := exec.Command("apt-get", args...)
	result, err := cmd.Output()

	if err != nil {
		log.Fatalf("安装依赖失败: %v", err)
	}

	return string(result)
}
