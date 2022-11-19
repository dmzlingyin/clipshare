package utils

import (
	"os"
	"strings"
)

func IsFile(s string) bool {
	if strings.HasPrefix(strings.TrimSpace(s), "file://") {
		// 从linux文件目录复制的文件
		s = s[7:]
	} else if strings.HasPrefix(strings.TrimSpace(s), "desktop:/") {
		// 从KDE桌面复制的文件
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return false
		}
		s = homeDir + s[8:]
	}

	fi, err := os.Stat(s)
	if err != nil {
		return false
	}
	return fi.Mode().IsRegular()
}
