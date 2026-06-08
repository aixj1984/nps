package stream

import "encoding/base64"

const typeScreenFrame = "screen-frame"

func bytesToBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
