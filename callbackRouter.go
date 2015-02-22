package gente

type CallbackRouter interface {
	CallbackForPath(string) (MessageHandlingFunc, error)
}
