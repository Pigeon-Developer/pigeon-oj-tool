package install

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/Pigeon-Developer/pigeon-oj-tool/shared/util"
	"github.com/hashicorp/go-version"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

type Option struct {
	Name  string
	Value string
}

func showInstallModeSelect() string {
	docker := "使用 docker compose"
	rawNodejs := "直接安装在当前机器（会安装 https://volta.sh 管理 nodejs 版本，使用包管理安装 MySQL 等依赖）"
	prompt := promptui.Select{
		Label: "选择安装方式",
		Items: []string{docker, rawNodejs},
	}

	_, selected, _ := prompt.Run()

	if selected == docker {
		return "docker"
	}
	if selected == rawNodejs {
		return "raw-nodejs"
	}
	return ""
}

func showInputPrompt(label string, defaultValue string) string {
	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultValue,
	}

	result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v \n", err)
	}

	return result
}

func showNumberInputPrompt(label string, defaultValue int) int {
	result := showInputPrompt(label, fmt.Sprintf("%d", defaultValue))
	num, err := strconv.Atoi(result)
	if err != nil {
		log.Fatalf("输入的不是数字 %v \n", err)
	}
	return num
}

func handleRawInstall() {
	//
	fmt.Printf("@TODO \n")
}

func getDockerVersion() string {
	version, err := exec.Command("docker", "version", "-f", "{{.Server.Version}}").Output()
	if err != nil {
		return ""
	}

	return strings.Trim(string(version), " \n")
}

func checkDockerVersion() bool {
	ver := getDockerVersion()

	if ver == "" {
		return false
	}

	baseVersion, err := version.NewVersion("25.0.0")

	if err != nil {
		log.Fatalf("docker 最小版本解析错误 \n")
	}

	userVersion, err := version.NewVersion(ver)

	if err != nil {
		return false
	}

	return userVersion.GreaterThanOrEqual(baseVersion)
}

func getComposeVersion() string {
	out, err := exec.Command("docker", "compose", "version").Output()
	if err != nil {
		return ""
	}

	result := string(out)
	reg := regexp.MustCompile(` v(\d+\.\d+\.\d+)`)
	version := reg.FindString(result)

	return strings.Trim(string(version), " \n")
}

func checkComposeVersion() bool {
	ver := getComposeVersion()
	if ver == "" {
		return false
	}

	baseVersion, err := version.NewVersion("2.0.0")

	if err != nil {
		log.Fatalf("docker compose 最小版本解析错误 \n")
	}

	userVersion, err := version.NewVersion(ver)

	if err != nil {
		fmt.Println(err)

		return false
	}

	return userVersion.GreaterThanOrEqual(baseVersion)
}

func handleDockerInstall() {
	// 默认安装到 /etc/pigeon-oj
	// docker 目录存放对应的配置文件
	if !checkDockerVersion() {
		fmt.Println("docker 版本不满足 >= 25.0.0 ")
		return
	}

	if !checkComposeVersion() {
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
	localPort := showNumberInputPrompt("输入 pigeon-oj 服务端口号", 3000)
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

func Run(cCtx *cli.Context) {
	mode := showInstallModeSelect()

	if mode == "docker" {
		// 检查 docker&compose 版本
		handleDockerInstall()
	}
	if mode == "raw-nodejs" {
		// 系统版本，并安装必要的依赖
		handleRawInstall()
	}
}
