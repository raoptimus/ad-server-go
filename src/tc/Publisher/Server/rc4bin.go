package main

import "strconv"

type Rc4Bin struct {
	key  [256]uint32
	i, j uint8
}

type KeySizeError int

func (k KeySizeError) Error() string {
	return "rc4bin: invalid key size " + strconv.Itoa(int(k))
}

var bin2Hex = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}

func NewRc4Bin(key string) (*Rc4Bin, error) {
	k := len(key)

	if k < 1 || k > 256 {
		return nil, KeySizeError(k)
	}

	var r Rc4Bin

	for i := 0; i < 256; i++ {
		r.key[i] = uint32(i)
	}

	var j uint8 = 0

	for i := 0; i < 256; i++ {
		j += uint8(r.key[i]) + key[i%k]
		r.key[i], r.key[j] = r.key[j], r.key[i]
	}

	return &r, nil
}

func (s *Rc4Bin) crypt(src string) string {
	dst := ""
	key := s.key
	i, j := s.i, s.j
	b := []byte(src)

	var ch uint32
	for _, v := range b {
		i += 1
		j += uint8(key[i])
		key[i], key[j] = key[j], key[i]
		ch = uint32(v) ^ key[uint8(key[i]+key[j])]
		dst += bin2Hex[(ch&0xf0)>>4]
		dst += bin2Hex[ch&0x0f]
	}

	return dst
}
