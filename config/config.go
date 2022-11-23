package config

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"sync"
)

var (
	ErrNotFound = errors.New("key not found")
)

type Config interface {
	Load() error
	Scan(v interface{}) error
	Value(key string) Value
}

type config struct {
	opts   options
	reader Reader
	cached sync.Map
}

func New(opts ...Option) Config {
	o := options{
		decoder: defaultDecoder,
	}

	for _, opt := range opts {
		opt(&o)
	}

	return &config{
		opts:   o,
		reader: newReader(o),
	}
}

func (c *config) Load() error {
	for _, src := range c.opts.sources {
		kvs, err := src.Load()
		if err != nil {
			return err
		}
		for _, v := range kvs {
			log.Debugf("config loaded: %s format: %s", v.Key, v.Format)
		}
		if err = c.reader.Merge(kvs...); err != nil {
			log.Errorf("failed to merge config source: %v", err)
			return err
		}

	}
	return nil
}

func (c *config) Scan(v interface{}) error {
	data, err := c.reader.Source()
	if err != nil {
		return err
	}
	return unmarshalJSON(data, v)
}

func (c *config) Value(key string) Value {
	if v, ok := c.cached.Load(key); ok {
		return v.(Value)
	}
	if v, ok := c.reader.Value(key); ok {
		c.cached.Store(key, v)
		return v
	}
	return &errValue{err: ErrNotFound}
}
