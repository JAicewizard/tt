package v2

func BoolToBytes(b *bool, buf *[]byte) {
	if *b {
		*buf = []byte{1}
	} else {
		*buf = []byte{0}
	}
}

func BoolFromBytes(buf []byte) bool {
	return buf[0] == 1
}
