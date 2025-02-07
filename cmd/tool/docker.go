package tool

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/Pigeon-Developer/pigeon-oj-tool/shared"
	"github.com/pterm/pterm"
	"periph.io/x/host/v3/distro"
)

func DockerInstall() {
	result, _ := pterm.DefaultInteractiveConfirm.WithDefaultText("安装时会卸载之前安装的 docker.io docker-doc docker-compose podman-docker containerd runc，确认执行 docker 安装脚本吗").Show()

	// Print a blank line for better readability.
	pterm.Println()

	if !result {
		return
	}

	if distro.IsDebian() {
		vendor := "debian"
		if distro.IsUbuntu() {
			vendor = "ubuntu"
		}
		cmd := exec.Command("bash", "-l", "-c", fmt.Sprintf("bash %s/static/docker-install-%s.sh", shared.LocalPath, vendor))

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			log.Fatalf("cmd.Run() failed: %v\n", err)
		}

	} else {
		fmt.Println("暂时只支持 ubuntu 和 debian")
	}
}
