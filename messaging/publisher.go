package messaging

import "bytes"

type Publisher interface {
	Publish(buffer bytes.Buffer) error
}