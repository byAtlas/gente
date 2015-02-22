package gente

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"testing"
)

func TestReplyUsesCallback(t *testing.T) {
	inboundPipe := make(chan []byte)

	msgId := uuid.NewUUID()
	callbacks := make(map[string]MessageHandlingFunc)
	callCount := 0

	callbacks[msgId.String()] = func(msg interface{}) (interface{}, error) {
		callCount++
		return nil, nil
	}

	bp := boundJsonCallbackPipeline{
		JsonCallbackPipeline: JsonCallbackPipeline{},
		inbound:              inboundPipe,
		outbound:             make(chan []byte),
		callbacks:            callbacks,
	}

	go bp.inboundLoop()

	msg := SockMessage{
		Id:      uuid.NewUUID(),
		ReplyTo: msgId,
	}

	msgBytes, _ := json.Marshal(msg)

	inboundPipe <- msgBytes

	if callCount != 1 {
		t.Error("Callback not invoked exactly once")
	}
}
