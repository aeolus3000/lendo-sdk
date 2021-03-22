package messaging

import "bytes"

type AcknowledgeFunc func() error
type NotAcknowledgeFunc func(requeue bool) error
type RejectFunc func(requeue bool) error

type Message struct {
	Body           *bytes.Buffer
	Acknowledge    AcknowledgeFunc
	NotAcknowledge NotAcknowledgeFunc
	Reject         RejectFunc
}
