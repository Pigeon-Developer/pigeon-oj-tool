package install

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Pigeon-Developer/pigeon-oj-tool/shared/util"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
)

func showInstallModeSelect() string {
	docker := "使用 docker compose"
	rawNodejs := "直接安装在当前机器（会安装 https://volta.sh 管理 nodejs 版本，使用包管理安装 MySQL 等依赖）"

	selected, _ := pterm.DefaultInteractiveSelect.WithFilter(false).WithOptions([]string{docker, rawNodejs}).Show("选择安装方式")

	if selected == docker {
		return "docker"
	}
	if selected == rawNodejs {
		return "raw-nodejs"
	}
	return ""
}

func handlePigeonojRawInstall() {
	//
	fmt.Printf("@TODO \n")
}

func handlePigeonojDockerInstall() {
	// 默认安装到 /etc/pigeon-oj
	// docker 目录存放对应的配置文件
	if !util.CheckDockerVersion() {
		fmt.Println("docker 版本不满足 >= 25.0.0 ")
		return
	}

	if !util.CheckComposeVersion() {
		fmt.Println("docker compose 版本不满足 >= 2.0.0 ")
		return
	}

	// 创建 /etc/pigeon-oj/docker 目录
	wd := "/etc/pigeon-oj/docker"
	err := os.MkdirAll(wd, 0755)

	if err != nil {
		fmt.Printf("创建 %s 目录失败 %v \n", wd, err)
		return
	}

	// 获取一个版本让用户选，或者使用用户指定的版本
	// 这里先固定使用 20250120.0214
	// 未来应该是从一个网站上，用户选一个版本，生成最终的安装指令

	// 创建 compose.yaml
	composeConfigPath := "/etc/pigeon-oj/docker/compose.yaml"
	if util.FileExists(composeConfigPath) {
		fmt.Println("compose.yaml 已存在，可能是之前已经安装过")
		return
	}

	file, err := os.OpenFile(composeConfigPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("创建 %s 失败 %v \n", composeConfigPath, err)
		return
	}

	imageVersion := "20250120.0214"
	localPort := util.ShowNumberInputPrompt("输入 pigeon-oj 服务端口号", 3000)
	composeContent := fmt.Sprintf(`
services:
  server:
    image: pigeonojdev/quick-start:%s
    ports:
      - %d:3000

`, imageVersion, localPort)
	file.WriteString(composeContent)

	cmd := exec.Command("docker", "compose", "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = wd
	cmd.Run()
}

func Pigeonoj(cCtx *cli.Context) {
	if !util.IsLinux() {
		fmt.Println("只支持 Linux 系统安装")
		return
	}
	mode := showInstallModeSelect()

	if mode == "docker" {
		// 检查 docker&compose 版本
		handlePigeonojDockerInstall()
	}
	if mode == "raw-nodejs" {
		// 系统版本，并安装必要的依赖
		handlePigeonojRawInstall()
	}
}
