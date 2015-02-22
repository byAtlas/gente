package gente

import "errors"

type RouterBuilder struct {
	paths     map[string]MessageHandlingFunc
	callbacks map[string]MessageHandlingFunc
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
