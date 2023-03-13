package helper

import (
	"crypto/md5"
	"fmt"
)

// GetMd5
// 生成MD5
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
