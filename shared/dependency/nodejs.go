package dependency

import (
	"log"
	"os/exec"

	"github.com/Pigeon-Developer/pigeon-oj-tool/shared/util"
)

// 使用 volta 管理 nodejs 版本
func InstallNodejsManager() {
	exec.Command("bach", "-l", "-c", "curl https://get.volta.sh | bash").Run()

	// 设置 pnpm 启用变量
	if util.FileExists("$HOME/.bashrc") {
		exec.Command("bach", "-l", "-c", "echo 'export VOLTA_FEATURE_PNPM=1' >> $HOME/.bashrc").Run()
	} else {
		log.Fatalln("无法找到 $HOME/.bashrc 文件")
	}
}
