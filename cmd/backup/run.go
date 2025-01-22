package backup

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Pigeon-Developer/pigeon-oj-tool/shared"
	"github.com/Pigeon-Developer/pigeon-oj-tool/shared/content"
	"github.com/Pigeon-Developer/pigeon-oj-tool/shared/dependency"
	"github.com/go-ini/ini"
	"github.com/nleeper/goment"
)

func RunHustoj() {
	// https://github.com/zhblue/hustoj/blob/master/trunk/install/bak.sh

	cfg, err := ini.LoadSources(ini.LoadOptions{
		SkipUnrecognizableLines: true,
	}, shared.HustojConfPath)

	if err != nil {
		fmt.Println("hustoj 配置文件读取失败 ", err)
		return
	}

	dbHost := cfg.Section("").Key("OJ_HOST_NAME").String()
	dbUser := cfg.Section("").Key("OJ_USER_NAME").String()
	dbPasswd := cfg.Section("").Key("OJ_PASSWORD").String()
	dbName := cfg.Section("").Key("OJ_DB_NAME").String()
	dbPort := cfg.Section("").Key("OJ_PORT_NUMBER").String()

	// 释放内置的文件到本地
	content.ExtractStatic()

	d, err := goment.New()

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

	// 之情数据清理的 sql
	exec.Command("bash", "-l", "-c", fmt.Sprintf("mysql -h %s -u%s -P %s -p%s %s < %s/static/hustoj-backup-clean.sql", dbHost, dbUser, dbPort, dbPasswd, dbName, shared.LocalPath))

	// 备份 db 文件
	exec.Command("mysqldump", "--default-character-set=utf8mb4", "-h", dbHost, "-u"+dbUser, "-P", dbPort, "-p"+dbPasswd, dbName, "--result-file", fmt.Sprintf("%s/db.sql", backupDir))

	backupList := []string{"data", "src/web", "src/core", "etc"}

	for _, path := range backupList {
		srcDir := fmt.Sprintf("/home/judge/%s", path)
		destDir := fmt.Sprintf("%s/%s", backupDir, path)

		err = os.MkdirAll(destDir, os.ModePerm)
		if err != nil {
			fmt.Println("创建备份目录失败 ", err)
			return
		}

		err = os.CopyFS(destDir, os.DirFS(srcDir))
		if err != nil {
			fmt.Println("复制文件失败 ", err)
			return
		}
	}

	dependency.InstallDepsForDebian([]string{"zstd"})

	// 打包备份目录到 tar.zstd
	exec.Command("tar", "-I", "zstd", "-cf", fmt.Sprintf("%s.tar.zst", backupDir), backupDir)

	fmt.Printf("备份完成, 文件名: %s.tar.zst\n", backupDir)
}
