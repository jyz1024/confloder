package confloder

import (
	"crypto/md5"
	"fmt"
)

func MD5(val string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(val)))
}

func BytesMD5(val []byte) string {
	return fmt.Sprintf("%x", md5.Sum(val))
}
