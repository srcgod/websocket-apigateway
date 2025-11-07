package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

type IDGenerator struct {
	prefix string
}

func NewIDGenerator(prefix string) *IDGenerator {
	return &IDGenerator{prefix: prefix}
}

func (g *IDGenerator) Generate() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%s_%d", g.prefix, time.Now().UnixNano())
	}
	return g.prefix + "_" + hex.EncodeToString(bytes)
}

var (
	ConnectionIDGenerator = NewIDGenerator("conn")
	SessionIDGenerator    = NewIDGenerator("sess")
	RequestIDGenerator    = NewIDGenerator("req")
	PingIDGenerator       = NewIDGenerator("ping")
)
