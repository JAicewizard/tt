package tt

import (
	"encoding/binary"

	v2 "github.com/JAicewizard/tt/v2"
	v1 "github.com/JAicewizard/tt/v1"
)

func Decodev2(b []byte, d *Data) (err error) {
	vlen := binary.BigEndian.Uint32(b[len(b)-4:len(b)-0])
	locs := make([]uint64, vlen)
	v2.GetLocs(b, locs, vlen)

	//decoding the actual values
	var v v1.Value
	v.FromBytes(b[locs[vlen-1]:])

	if *d == nil {
		*d = make(Data, len(v.Children)*1)
	}

	data := d
	childs := v.Children

	for ck := range childs {
		var err error
		v.FromBytes(b[locs[childs[ck]]:])

		err = valueToMapv1(&v, *data, locs, b)

		if err != nil {
			return err
		}
	}
	return nil
}