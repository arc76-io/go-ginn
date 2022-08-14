package xor

import (
	"crypto/sha256"
	"encoding/binary"
	"math"
)

type Shift64 struct {
	seed int64
}

func NewShift64(seed int64) *Shift64 {
	seed += 1
	seed ^= seed >> 12
	seed ^= seed << 25
	seed ^= seed >> 27
	seed *= 0x2345F4914F6CFD1E

	return &Shift64 {
		seed: seed,
	}
}

func (ref *Shift64) shift() int64 {
	ref.seed ^= ref.seed >> 12
	ref.seed ^= ref.seed << 25
	ref.seed ^= ref.seed >> 27
	ref.seed *= 0x2345F4914F6CFD1E
	return ref.seed
}

func (ref *Shift64) Int32(n int32) int32 {
	x := float64(ref.shift()) / 9223372036854775808
	return int32(math.Abs(x * float64(n)) - 1.0)
}

func GetSeed(str string) int64 {
	sum := sha256.Sum256([]byte(str))
	da := binary.BigEndian.Uint64(sum[0:8])
	db := binary.BigEndian.Uint64(sum[8:16])
	dc := binary.BigEndian.Uint64(sum[16:24])
	dd := binary.BigEndian.Uint64(sum[24:32])
	return int64(da * db * dc * dd)
}