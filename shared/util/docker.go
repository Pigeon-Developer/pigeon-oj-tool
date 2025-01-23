package util

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/hashicorp/go-version"
)

func getDockerVersion() string {
	version, err := exec.Command("docker", "version", "-f", "{{.Server.Version}}").Output()
	if err != nil {
		return ""
	}

	return strings.Trim(string(version), " \n")
}

func CheckDockerVersion() bool {
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

func CheckComposeVersion() bool {
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
