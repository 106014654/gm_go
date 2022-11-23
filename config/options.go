package config

import (
	"fmt"
	"gm_go/encoding"
	"strings"

	_ "gm_go/encoding/json"
	_ "gm_go/encoding/proto"
	_ "gm_go/encoding/xml"
	_ "gm_go/encoding/yaml"
)

type Decoder func(*KeyValue, map[string]interface{}) error

type Option func(*options)

type options struct {
	sources []Source
	decoder Decoder
}

func WithSource(s ...Source) Option {
	return func(o *options) {
		o.sources = s
	}
}

func defaultDecoder(src *KeyValue, target map[string]interface{}) error {
	if src.Format == "" {
		// expand key "aaa.bbb" into map[aaa]map[bbb]interface{}
		keys := strings.Split(src.Key, ".")
		for i, k := range keys {
			if i == len(keys)-1 {
				target[k] = src.Value
			} else {
				sub := make(map[string]interface{})
				target[k] = sub
				target = sub
			}
		}
		return nil
	}

	if codec := encoding.GetCodec(src.Format); codec != nil {
		return codec.Unmarshal(src.Value, &target)
	}
	return fmt.Errorf("unsupported key: %s format: %s", src.Key, src.Format)
}
