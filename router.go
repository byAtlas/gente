package gente

import (
	"errors"
)

type Router interface {
	//Request and response should be JSON marshall/unmarshallabe
	Route(path string, message interface{}) (interface{}, error)
}

type DefaultRouter struct {
	paths map[string]func(interface{}) (interface{}, error)
}

func (r DefaultRouter) Route(path string, message interface{}) (interface{}, error) {
	if method, ok := r.paths[path]; ok {
		response, err := method(message)

		return response, err
	} else {
		return nil, errors.New("Couldn't find matching path")
	}
}

func (r *DefaultRouter) AddRoute(path string, method MessageHandlingFunc) error {
	if _, ok := r.paths[path]; !ok {
		return errors.New("Route already exists.") //TODO: Add additional info further up
	}

	r.paths[path] = method

	return nil
}
