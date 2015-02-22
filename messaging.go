package gente

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"github.com/Sirupsen/logrus"
)

type SockMessage struct {
	Id      uuid.UUID
	ReplyTo uuid.UUID
	Path    string
	Body    interface{} //Will be by handler method.
}

type MessageHandlingFunc func(interface{}) (interface{}, error)

type MessagePipeline interface {
	Register(inbound chan []byte, outbound chan []byte)
}

type JsonCallbackPipeline struct {
	log    *logrus.Logger
	router Router
}

type boundJsonCallbackPipeline struct {
	JsonCallbackPipeline

	inbound  chan []byte
	outbound chan []byte

	callbacks map[string]MessageHandlingFunc
}

func (p *JsonCallbackPipeline) Register(inbound chan []byte, outbound chan []byte) {
	pipeline := &boundJsonCallbackPipeline{
		inbound:              inbound,
		outbound:             outbound,
		JsonCallbackPipeline: *p}
	go pipeline.inboundLoop()
}

func (p *boundJsonCallbackPipeline) inboundLoop() {
	for msgBytes, ok := <-p.inbound; ok; msgBytes, ok = <-p.inbound {
		msg := &SockMessage{}
		err := json.Unmarshal(msgBytes, msg)

		if err != nil {
			p.log.WithField("Error", err).Error("Error destructuring inbound message")
		}

		if msg.ReplyTo != nil {
			if fn, ok := p.callbacks[msg.ReplyTo.String()]; ok {
				fn(msg.Body)
			}

			continue
		}

		response, err := p.router.Route(msg.Path, msg.Body)

		if err != nil {
			p.log.WithFields(logrus.Fields{
				"MsgId":       msg.Id,
				"Error":       err,
				"Path":        msg.Path,
				"MessageBody": msg.Body,
			}).Error("Router returned error.")

			p.replyWithError(msg.Id)
		}

		if response != nil {
			if callbackRouter, ok := p.router.(CallbackRouter); ok {
				fn, err := callbackRouter.CallbackForPath(msg.Path)
				if err != nil {
					p.log.WithFields(logrus.Fields{
						"MsgId":       msg.Id,
						"Error":       err,
						"Path":        msg.Path,
						"MessageBody": msg.Body,
					}).Error("Router returned error.")

					p.replyWithError(msg.Id)
				}

				p.replyWithCallback(msg.Id, response, fn)
			} else {
				p.reply(msg.Id, response)
			}
		}
	}
}

func stubReply(toMsgId uuid.UUID, message interface{}) SockMessage {
	return SockMessage{
		Id:      uuid.NewUUID(),
		ReplyTo: toMsgId,
		Body:    message,
	}
}

func (p *boundJsonCallbackPipeline) replyWithError(toMsgId uuid.UUID) {
	p.reply(toMsgId, "Error handling request.")
}

func (p *boundJsonCallbackPipeline) reply(toMsgId uuid.UUID, message interface{}) {
	msg := stubReply(toMsgId, message)

	msgBytes, err := json.Marshal(msg)

	if err != nil {
		p.log.WithFields(logrus.Fields{
			"replyTo":   toMsgId,
			"replyBody": message,
		}).Error("Couldn't marshal reply.")
		return
	}

	p.outbound <- msgBytes
}

func (p *boundJsonCallbackPipeline) replyWithCallback(toMsgId uuid.UUID, message interface{}, callback MessageHandlingFunc) {

}
