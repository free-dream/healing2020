package tools

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5String(str string) string {
	m5 := md5.New()
	_, err := m5.Write([]byte(str))
	if err != nil {
		panic(err)
	}
	md5String := hex.EncodeToString(m5.Sum(nil))
	if len(md5String) > 8 {
		return md5String[:8]
	} else {
		return md5String
	}

}
