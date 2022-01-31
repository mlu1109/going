package going

import (
	"crypto/md5"
	"encoding/hex"
)

type Checksum func(s string) (string, error)

func DefaultChecksumFn(s string) (string, error) {
	sum := md5.Sum([]byte(s))
	h := hex.EncodeToString(sum[:])
	return h, nil
}
