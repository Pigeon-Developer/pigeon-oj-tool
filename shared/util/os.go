package util

import "runtime"

func IsLinux() bool {
	return runtime.GOOS == "linux"
}
