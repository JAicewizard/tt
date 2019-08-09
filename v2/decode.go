package v2

//GetLocs gets the location data from an array
func GetLocs(b []byte, locs []uint64, len uint32) {
	locs[0] = 0

	for i := uint32(1); i < len; i++ {
		locs[i] = locs[i-1] + uint64(
			uint32(b[locs[i-1]]*ikeylen)+ //this is the length for the children
				uint32(Getvaluelen(b[locs[i-1]+1:locs[i-1]+1+valuelenbytes]))+ //this is the value length
				uint32(Getkeylen(b[locs[i-1]+1+valuelenbytes:locs[i-1]+1+valuelenbytes+keylenbytes]))+ //this is the key length
				2+valuelenbytes+keylenbytes) //add 4 so that we cound the length values as wel, +1 is for going to the next value
	}
}
