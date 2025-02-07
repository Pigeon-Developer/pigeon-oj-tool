package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Pigeon-Developer/pigeon-oj-tool/cmd/backup"
	"github.com/Pigeon-Developer/pigeon-oj-tool/cmd/install"
	"github.com/Pigeon-Developer/pigeon-oj-tool/cmd/tool"
	"github.com/Pigeon-Developer/pigeon-oj-tool/shared/content"
	"github.com/urfave/cli/v2"
)

// 工具默认存放到 /etc/pigeon-oj-tool
func main() {
	app := &cli.App{
		Name:  "pojt",
		Usage: "pigeon-oj 安装升级维护工具",
		Commands: []*cli.Command{
			{
				Name:    "install",
				Aliases: []string{"i"},
				Usage:   "安装一个新的 pigeon-oj",
				Action: func(cCtx *cli.Context) error {
					install.Pigeonoj(cCtx)
					return nil
				},
			},
			{
				Name:    "upgrade",
				Aliases: []string{"u"},
				Usage:   "更新当前 pigeon-oj 版本",
				Action: func(cCtx *cli.Context) error {
					fmt.Println("completed task: ", cCtx.Args().First())
					return nil
				},
			},
			{
				Name:    "backup",
				Aliases: []string{"b"},
				Usage:   "备份当前 pigeon-oj 的数据",
				Action: func(cCtx *cli.Context) error {
					fmt.Println("completed task: ", cCtx.Args().First())
					return nil
				},
				Subcommands: []*cli.Command{
					{
						Name:  "hustoj",
						Usage: "备份 hustoj 的数据",
						Action: func(cCtx *cli.Context) error {
							backup.Hustoj()
							return nil
						},
					},
				},
			},
			{
				Name:    "tool",
				Aliases: []string{"t"},
				Usage:   "安装外部依赖",
				Action: func(cCtx *cli.Context) error {
					return nil
				},
				Subcommands: []*cli.Command{
					{
						Name:  "docker",
						Usage: "安装 docker",
						Action: func(cCtx *cli.Context) error {
							tool.DockerInstall()
							return nil
						},
					},
				},
			},
			{
				Name:    "setup",
				Aliases: []string{"s"},
				Usage:   "初始化 pigeon-oj-tool 配置",
				Action: func(cCtx *cli.Context) error {
					// 释放静态资源到本地
					content.ExtractStatic()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
