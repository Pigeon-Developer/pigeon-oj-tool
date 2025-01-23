package backup

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/Pigeon-Developer/pigeon-oj-tool/shared"
	"github.com/Pigeon-Developer/pigeon-oj-tool/shared/config"
	"github.com/Pigeon-Developer/pigeon-oj-tool/shared/content"
	"github.com/Pigeon-Developer/pigeon-oj-tool/shared/dependency"
	"github.com/nleeper/goment"
	"github.com/pterm/pterm"
)

func Hustoj() {
	// https://github.com/zhblue/hustoj/blob/master/trunk/install/bak.sh

	step1, _ := pterm.DefaultSpinner.Start("读取 hustoj 配置...")
	cfg, err := config.LoadHustojConfig()

	if err != nil {
		step1.Fail("hustoj 配置文件读取失败 ", err)
		return
	}
	step1.Success()

	dbHost := cfg.OJ_HOST_NAME
	dbUser := cfg.OJ_USER_NAME
	dbPasswd := cfg.OJ_PASSWORD
	dbName := cfg.OJ_DB_NAME
	dbPort := cfg.OJ_PORT_NUMBER

	// 释放内置的文件到本地
	content.ExtractStatic()

	// 展示 +8 时间
	loc := time.FixedZone("Asia/Shanghai", 8*60*60)
	d, err := goment.New(time.Now().In(loc))

	if err != nil {
		fmt.Println("时间获取失败 ", err)
		return
	}

	// 为这次备份创建一个目录名
	backName := fmt.Sprintf("%s-hustoj", d.Format("YYYYMMDD.HHmmss"))
	backupDir := fmt.Sprintf("%s/backup/%s", shared.LocalPath, backName)

	err = os.MkdirAll(backupDir, os.ModePerm)
	if err != nil {
		fmt.Println("创建备份目录失败 ", err)
		return
	}

	// 执行数据清理的 sql
	step2, _ := pterm.DefaultSpinner.Start("清理 db 老旧数据...")
	step2Result, err := exec.Command("bash", "-l", "-c", fmt.Sprintf("mysql -h %s -u%s -P %s -p%s %s < %s/static/hustoj-backup-clean.sql", dbHost, dbUser, dbPort, dbPasswd, dbName, shared.LocalPath)).CombinedOutput()
	if err != nil {
		step2.Fail("清理 db 数据失败 ", err)
		fmt.Println(string(step2Result))
		return
	}
	step2.Success()

	// 备份 db 文件
	step3, _ := pterm.DefaultSpinner.Start("备份 db 数据...")
	step3Result, err := exec.Command("mysqldump", "--default-character-set=utf8mb4", "-h", dbHost, "-u"+dbUser, "-P", dbPort, "-p"+dbPasswd, dbName, "--result-file", fmt.Sprintf("%s/db.sql", backupDir)).CombinedOutput()
	if err != nil {
		step3.Fail("备份 db 数据失败 ", err)
		fmt.Println(string(step3Result))
		return
	}
	step3.Success()

	step4, _ := pterm.DefaultSpinner.Start("生成备份文件...")
	backupList := []string{"data", "src/web", "src/core", "etc"}

	for _, path := range backupList {
		srcDir := fmt.Sprintf("/home/judge/%s", path)
		destDir := fmt.Sprintf("%s/%s", backupDir, path)

		err = os.MkdirAll(destDir, os.ModePerm)
		if err != nil {
			step4.Fail("创建备份目录失败 ", err)
			return
		}

		err = os.CopyFS(destDir, os.DirFS(srcDir))
		if err != nil {
			step4.Fail("复制文件失败 ", err)
			return
		}
	}

	dependency.InstallDepsForDebian([]string{"zstd"})

	// 打包备份目录到 tar.zstd
	exec.Command("tar", "-I", "zstd", "-cf", fmt.Sprintf("%s.tar.zst", backupDir), backupDir).Run()

	step4.Success(fmt.Sprintf("备份完成, 文件名: %s.tar.zst\n", backupDir))
}
