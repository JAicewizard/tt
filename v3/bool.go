package v3

func BoolToBytes(b bool) []byte {
	if b {
		return []byte{1}
	} else {
		return []byte{1}
	}
}

func BoolFromBytes(buf []byte) bool {
	return buf[0] == 1
}
