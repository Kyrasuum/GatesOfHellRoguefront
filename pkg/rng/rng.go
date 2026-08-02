package rng

import (
	"errors"
	"math"
	"math/rand"

	"roguefront/pkg/lmath"
)

type PityRng struct {
	r *rand.Rand

	target float64
	fails  float64
	offset int
}

func HashU32(x uint32) uint32 {
	x ^= x >> 16
	x *= 0x7feb352d
	x ^= x >> 15
	x *= 0x846ca68b
	x ^= x >> 16
	return x
}

func Hash2(cellX, cellZ int32, seed uint32) uint32 {
	h := uint32(cellX)*0x8da6b343 ^ uint32(cellZ)*0xd8163841 ^ seed*0xcb1ab31f
	return HashU32(h)
}

func U01(x uint32) float32 { return float32(x>>8) * (1.0 / 16777216.0) }

func NewPityRng(chance float64) (*PityRng, error) {
	if chance < 0 || math.IsNaN(chance) {
		return nil, errors.New("chance must be in [0,]")
	}
	seed := rand.Int63()
	r := rand.New(rand.NewSource(seed))

	offset := lmath.Floor(chance)

	prd := &PityRng{
		r:      r,
		target: chance - float64(offset),
		fails:  0,
		offset: offset,
	}

	return prd, nil
}

func RollOnce(chance float64) int {
	roller, err := NewPityRng(chance)
	if err != nil {
		return 0
	}
	return roller.Next()
}

func (p *PityRng) Next() int {
	if p == nil {
		return 0
	}
	roll := p.r.Float64()
	chance := p.ChanceNow()
	p.fails++
	if roll < chance {
		if chance < 0.5 {
			expected := float64(lmath.Ceil(1 / p.target))
			p.fails -= expected
		} else {
			expected := 1.0 / p.target
			p.fails -= expected
		}
		return 1 + p.offset
	}
	return p.offset
}

func (p *PityRng) ChanceNow() float64 {
	if p == nil {
		return 0.0
	}
	if p.target == 0 {
		return 0.0
	}
	if p.target == 1 {
		return 1
	}

	expected := float64(lmath.Ceil(1 / p.target))
	if p.fails < 0 {
		return p.target * (expected - lmath.Abs(p.fails)) / expected
	}

	// return p.target * float64(p.fails+1)
	return 1 / (expected - p.fails)
}

func (p *PityRng) Reset() {
	if p == nil {
		return
	}
	p.fails = 0
}

func (p *PityRng) Reseed(seed int64) {
	if p == nil {
		return
	}
	p.r.Seed(seed)
	p.Reset()
}
