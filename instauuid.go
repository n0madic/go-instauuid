package instauuid

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sync"
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
	mu            sync.Mutex
}

// NewGenerator initializes a new Generator with a given shard ID and epoch
func NewGenerator(shardID uint32, epoch int64) *Generator {
	if shardID > maxShard {
		panic(fmt.Sprintf("Shard ID exceeds the maximum allowed value: %d", maxShard))
	}

	if epoch == 0 {
		epoch = instagramEpoch
	}

	g := &Generator{
		epoch:         epoch,
		lastTimestamp: 0,
		sequence:      0,
		shardID:       shardID,
	}
	return g
}

// GenerateID generates the next unique ID based on timestamp, shard ID, and sequence
func (g *Generator) GenerateID() uint64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	timestamp := time.Now().UnixMilli() - g.epoch

	if timestamp == g.lastTimestamp {
		g.sequence = (g.sequence + 1) & maxSeq
		if g.sequence == 0 {
			timestamp = g.waitNextMillis(timestamp)
		}
	} else if timestamp > g.lastTimestamp {
		g.sequence = 0
	} else {
		timestamp = g.waitNextMillis(g.lastTimestamp)
	}

	g.lastTimestamp = timestamp

	id := (uint64(timestamp) << timeShift) |
		(uint64(g.shardID) << shardShift) |
		uint64(g.sequence)

	return id
}

func (g *Generator) waitNextMillis(lastTimestamp int64) int64 {
	timestamp := time.Now().UnixMilli() - g.epoch
	for timestamp <= lastTimestamp {
		time.Sleep(100 * time.Microsecond)
		timestamp = time.Now().UnixMilli() - g.epoch
	}
	return timestamp
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
