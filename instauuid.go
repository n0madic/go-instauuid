package instauuid

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"
)

const (
	instagramEpoch = int64(1314220021721) // 2011-08-24T21:07:01Z
	seqBits        = 10
	shardBits      = 13
	maxSeq         = (1 << seqBits) - 1   // 1023
	maxShard       = (1 << shardBits) - 1 // 8191
	timeShift      = seqBits + shardBits  // 23
	shardShift     = seqBits
)

// Generator structure to generate IDs
type Generator struct {
	epoch         int64
	lastTimestamp int64
	sequence      uint32
	shardID       uint32
}

// NewGenerator initializes a new Generator with a given shard ID and epoch
func NewGenerator(shardID uint32, epoch int64) *Generator {
	if shardID > maxShard {
		panic(fmt.Sprintf("Shard ID exceeds the maximum allowed value: %d", maxShard))
	}
	if epoch == 0 {
		epoch = instagramEpoch
	}
	return &Generator{
		epoch:   epoch,
		shardID: shardID,
	}
}

// GenerateID returns a new unique ID
func (g *Generator) GenerateID() uint64 {
	for {
		timestamp := time.Now().UnixMilli() - g.epoch
		lastTimestamp := atomic.LoadInt64(&g.lastTimestamp)
		if timestamp == lastTimestamp {
			seq := atomic.AddUint32(&g.sequence, 1) & maxSeq
			if seq == 0 {
				continue // Wait for the next millisecond
			}
			return g.assembleID(uint64(timestamp), seq)
		}
		if timestamp > lastTimestamp {
			atomic.StoreInt64(&g.lastTimestamp, timestamp)
			atomic.StoreUint32(&g.sequence, 0)
			return g.assembleID(uint64(timestamp), 0)
		}
		// If we're here, timestamp < lastTimestamp. This should be rare.
		time.Sleep(time.Millisecond)
	}
}
func (g *Generator) assembleID(timestamp uint64, seq uint32) uint64 {
	return (timestamp << timeShift) |
		(uint64(g.shardID) << shardShift) |
		uint64(seq)
}

// GenerateBase64 generates a Base64 encoded ID
func (g *Generator) GenerateBase64() string {
	id := g.GenerateID()
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, id)
	return base64.RawURLEncoding.EncodeToString(buf)
}

// GenerateHex generates a hexadecimal string ID
func (g *Generator) GenerateHex() string {
	id := g.GenerateID()
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, id)
	return hex.EncodeToString(buf)
}

// GenerateBuffer generates a byte buffer ID (Little Endian)
func (g *Generator) GenerateBuffer() []byte {
	id := g.GenerateID()
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, id)
	return buf
}

// GenerateBufferBE generates a byte buffer ID (Big Endian)
func (g *Generator) GenerateBufferBE() []byte {
	id := g.GenerateID()
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, id)
	return buf
}
