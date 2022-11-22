package config

type KeyValue struct {
	Key    string
	Value  []byte
	Format string
}

type Source interface {
	Load() ([]*KeyValue, error)
}
