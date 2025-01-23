package config

import (
	"github.com/Pigeon-Developer/pigeon-oj-tool/shared"
	"github.com/go-ini/ini"
)

type HustojConfig struct {
	OJ_HOST_NAME   string
	OJ_USER_NAME   string
	OJ_PASSWORD    string
	OJ_DB_NAME     string
	OJ_PORT_NUMBER string
}

func LoadHustojConfig() (*HustojConfig, error) {
	cfg, err := ini.LoadSources(ini.LoadOptions{
		SkipUnrecognizableLines: true,
	}, shared.HustojConfPath)

	if err != nil {
		return nil, err
	}

	p := &HustojConfig{}
	cfg.MapTo(p)

	return p, nil
}
