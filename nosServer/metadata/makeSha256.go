package metadata

import (
	"crypto/sha256"
	"fmt"
)

func MakeSha256(buf []byte) string {
	h := sha256.New()
	h.Write(buf)
	return fmt.Sprintf("%x", h.Sum(nil))
}
