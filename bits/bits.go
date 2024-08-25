package bits

import "errors"

// FromUint16 converts a uint16 value to a slice of bools where every set bit is
// represented by a true value
func FromUint16(data uint16) []bool {
	b := make([]bool, 0, 16)
	for i := 0; i < 16; i++ {
		b = append(b, (data&0x8000) > 0)
		data <<= 1
	}

	return b
}

// ToUint16 converts a slice of bools to a uint16
func ToUint16(bools []bool) (uint16, error) {
	if len(bools) != 16 {
		return 0, errors.New("len of bools must be exactly 16")
	}

	var i uint16
	for _, b := range bools {
		i <<= 1
		if b {
			i |= 0x0001
		}
	}

	return i, nil
}

// FromBytess converts a slice of bytes to a slice of bools where every
// set bit is represented by a true value
func FromBytes(data []byte) []bool {
	b := make([]bool, 0, len(data)*8)
	for _, d := range data {
		for i := 0; i < 8; i++ {
			b = append(b, (d&0x80) > 0)
			d <<= 1
		}
	}

	return b
}

// ToBytes converts a slice of bools to a byte slice
func ToBytes(bools []bool) ([]byte, error) {
	if len(bools)%8 != 0 {
		return nil, errors.New("len of bools must be divsible by 8")
	}

	var ret []byte
	var item byte
	for i, b := range bools {
		item <<= 1
		if b {
			item |= 0x0001
		}
		if (i+1)%8 == 0 {
			ret = append(ret, item)
			item = 0
		}
	}

	return ret, nil
}
