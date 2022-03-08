package decoder

import (
	"io"

	yaml "gopkg.in/yaml.v3"
)

// Decoder container the methods to Init a Decode a Decoder
type Decoder interface {
	Init(reader io.Reader)
	Decode(node *yaml.Node) error
}

type yamlDecoder struct {
	decoder yaml.Decoder
}

// NewYamlDecoder creates a new Decoder
func NewYamlDecoder() Decoder {
	return &yamlDecoder{}
}

func (dec *yamlDecoder) Init(reader io.Reader) {
	dec.decoder = *yaml.NewDecoder(reader)
}

func (dec *yamlDecoder) Decode(rootYamlNode *yaml.Node) error {
	return dec.decoder.Decode(rootYamlNode)
}
