package async_op

import "hash/crc32"

func StrToBindId(strVal string) int {
	v := int(crc32.ChecksumIEEE([]byte(strVal)))

	if v >= 0 {
		return v
	} else {
		return -v
	}
}
