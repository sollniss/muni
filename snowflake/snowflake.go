package snowflake

import (
	"strconv"
	"sync/atomic"
	"time"
)

type gen[T ~uint64] struct {
	lastID *uint64

	epoch   uint64
	nsTicks uint64

	timeStampBits uint64
	nodeBits      uint64
	seqBits       uint64

	node uint64

	timeStampOffset uint64
	nodeOffset      uint64

	maxSeq       uint64
	maxTimestamp uint64
}

// NewInstagram returns a generator initialized with the Instagram default values.
// startTime is 2015-01-01 00:00:00 UTC
// tick is 1000000 ns (1 millisecond)
// timeStampBits is 41 bits
// nodeBits is 13 bits
// seqBits is 10 bits
//
// Note that the Instagram variation modulos the sequence by 1024 and instead of busy waiting for the next tick.
func NewInstagram[T ~uint64](nodeID uint64) *gen[T] {
	return New[T](time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), nodeID, 1*time.Millisecond, 41, 13, 10)
}

// NewDiscord returns a generator initialized with the Discord default values.
// startTime is 2015-01-01 00:00:00 UTC
// tick is 1000000 ns (1 millisecond)
// timeStampBits is 41 bits
// nodeBits is 10 bits
// seqBits is 12 bits
func NewDiscord[T ~uint64](nodeID uint64) *gen[T] {
	return New[T](time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), nodeID, 1*time.Millisecond, 41, 10, 12)
}

// NewTwitter returns a generator initialized with the Twitter default values.
// startTime is 2010-11-04 01:42:54 UTC
// tick is 1000000 ns (1 millisecond)
// timeStampBits is 41 bits
// nodeBits is 10 bits
// seqBits is 12 bits
func NewTwitter[T ~uint64](nodeID uint64) *gen[T] {
	return New[T](time.Date(2010, 11, 4, 1, 42, 54, 0, time.UTC), nodeID, 1*time.Millisecond, 41, 10, 12)
}

// NewDefault returns a generator initialized with the default bit sizes.
// timeStampBits is 41 bits
// nodeBits is 10 bits
// seqBits is 12 bits
func NewDefault[T ~uint64](startTime time.Time, nodeID uint64) *gen[T] {
	return New[T](startTime, nodeID, 1*time.Millisecond, 41, 10, 12)
}

// New returns a generator that generates IDs based on the Twitter Snowflake algorithm.
//
// timeStampBits is the bit count for the timestamp, representing the time since the chosen epoch.
// nodeBits is the bit count for the machine ID, preventing clashes.
// seqBits is the bit count for a per-machine sequence number, to allow creation of multiple snowflakes in the same epoch tick.
//
// Panics nodeID is < 0 or out of range,
// tick is < 1,
// timeStampBits + nodeBits + seqBits is > 64,
// or seqBits is 0.
func New[T ~uint64](startTime time.Time, nodeID uint64, tick time.Duration, timeStampBits uint8, nodeBits uint8, seqBits uint8) *gen[T] {
	var nodeMax uint64 = 1<<seqBits - 1
	if nodeID < 0 || nodeID > nodeMax {
		panic("snowflake: node ID must be between 0 and " + strconv.FormatUint(nodeMax, 10))
	}

	if tick < 1 {
		panic("snowflake: ticks must be > 0")
	}

	if timeStampBits+nodeBits+seqBits > 64 {
		panic("snowflake: invalid bit lengths")
	}

	if seqBits == 0 {
		panic("snowflake: seqBits must not be 0")
	}

	var i uint64 = 0
	return &gen[T]{
		timeStampBits: uint64(timeStampBits),
		epoch:         uint64(startTime.UnixNano()),
		nsTicks:       uint64(tick),
		node:          uint64(nodeID),

		nodeBits: uint64(nodeBits),

		timeStampOffset: uint64(nodeBits + seqBits),
		nodeOffset:      uint64(seqBits),

		maxSeq:       1<<seqBits - 1,
		maxTimestamp: 1<<timeStampBits - 1,

		lastID: &i,
	}
}

// ID generates a unique ID packed into an uint64.
//
// Panics if the epoch has ended.
func (g gen[T]) ID() uint64 {

	for {
		localLastID := atomic.LoadUint64(g.lastID)
		seq := localLastID & g.maxSeq
		lastIDTime := localLastID >> g.timeStampOffset

		now := g.ticksSinceEpoch()
		if now > g.maxTimestamp {
			panic("snowflake: the maximum life cycle has ended, please check startTime")
		} else if now > lastIDTime {
			seq = 0
		} else if seq >= g.maxSeq {
			continue
		} else {
			seq++
		}

		newID := (now << g.timeStampOffset) | (g.node << g.nodeOffset) | seq
		if atomic.CompareAndSwapUint64(g.lastID, localLastID, newID) {
			return newID
		}
	}
}

func (g gen[T]) ticksSinceEpoch() uint64 {
	return (uint64(time.Now().UnixNano()) - g.epoch) / g.nsTicks
}
