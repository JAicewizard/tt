package v3

//BoolToBytes converts a bool into bytes
func BoolToBytes(b bool) []byte {
	if b {
		return []byte{1}
	}
	return []byte{0}
}

//BoolFromBytes converts bytes into a bool
func BoolFromBytes(buf []byte) bool {
	return buf[0] == 1
}
