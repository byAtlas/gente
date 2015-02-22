package gente

import (
	"errors"
)

type Router interface {
	//Request and response should be JSON marshall/unmarshallabe
	Route(path string, message interface{}) (interface{}, error)
	CallbackForRoute(path string) (MessageHandlingFunc, error)
}

type RouterBuilder struct {
	paths     map[string]MessageHandlingFunc
	callbacks map[string]MessageHandlingFunc
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

func (r *RouterBuilder) AddRoute(path string, method MessageHandlingFunc) error {
	if r.paths == nil {
		r.paths = make(map[string]MessageHandlingFunc)
	}

	if _, ok := r.paths[path]; ok {
		return errors.New("Route already exists.") //TODO: Add additional info further up
	}

	r.paths[path] = method

	return nil
}

func (r *RouterBuilder) AddRouteWithCallback(path string, method MessageHandlingFunc, callback MessageHandlingFunc) error {
	err := r.AddRoute(path, method)

	if err != nil {
		return err
	}

	err = r.AddCallbackForRoute(path, callback)

	if err != nil {
		delete(r.paths, path)
		return err
	}

	return nil
}

func (r *RouterBuilder) AddCallbackForRoute(path string, callback MessageHandlingFunc) error {
	if r.callbacks == nil {
		r.paths = make(map[string]MessageHandlingFunc)
	}

	if _, ok := r.callbacks[path]; ok {
		return errors.New("Route already exists.") //TODO: Add additional info further up
	}

	return nil
}

func (r *RouterBuilder) Finalize() Router {
	return defaultRouter{paths: r.paths, callbacks: r.callbacks}
}
