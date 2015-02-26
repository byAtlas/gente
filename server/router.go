package gente

import "errors"

type Router interface {
	//Request and response should be JSON marshall/unmarshallabe
	Route(path string, message interface{}) (interface{}, error)
	CallbackForRoute(path string) (MessageHandlingFunc, error)
}

type defaultRouter struct {
	paths     map[string]MessageHandlingFunc
	callbacks map[string]MessageHandlingFunc
}

func (r defaultRouter) Route(path string, message interface{}) (interface{}, error) {
	if method, ok := r.paths[path]; ok {
		response, err := method(message)

		return response, err
	} else {
		return nil, errors.New("Couldn't find matching path")
	}
}

func (r defaultRouter) CallbackForRoute(path string) (MessageHandlingFunc, error) {
	if method, ok := r.callbacks[path]; ok {
		return method, nil
	} else {
		return nil, errors.New("Couldn't find matching path")
	}
}
