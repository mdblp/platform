package data

import (
	"github.com/tidepool-org/platform/net"
	"github.com/tidepool-org/platform/structure"
)

type DeduplicatorDescriptor struct {
	Name    *string `json:"name,omitempty" bson:"name,omitempty"`
	Version *string `json:"version,omitempty" bson:"version,omitempty"`
	Hash    *string `json:"hash,omitempty" bson:"hash,omitempty"`
}

func ParseDeduplicatorDescriptor(parser structure.ObjectParser) *DeduplicatorDescriptor {
	if !parser.Exists() {
		return nil
	}
	datum := NewDeduplicatorDescriptor()
	parser.Parse(datum)
	return datum
}

func ParseDeduplicatorDescriptorDEPRECATED(parser ObjectParser) *DeduplicatorDescriptor {
	if parser.Object() == nil {
		return nil
	}
	datum := NewDeduplicatorDescriptor()
	datum.ParseDEPRECATED(parser)
	parser.ProcessNotParsed()
	return datum
}

func NewDeduplicatorDescriptor() *DeduplicatorDescriptor {
	return &DeduplicatorDescriptor{}
}

func (d *DeduplicatorDescriptor) Parse(parser structure.ObjectParser) {
	d.Name = parser.String("name")
}

func (d *DeduplicatorDescriptor) ParseDEPRECATED(parser ObjectParser) {
	d.Name = parser.ParseString("name")
}

func (d *DeduplicatorDescriptor) Validate(validator structure.Validator) {
	validator.String("name", d.Name).Using(net.ReverseDomainValidator)
	validator.String("version", d.Version).Using(net.SemanticVersionValidator)
	validator.String("hash", d.Hash).NotEmpty()
}

func (d *DeduplicatorDescriptor) Normalize(normalizer structure.Normalizer) {}

func (d *DeduplicatorDescriptor) NormalizeDEPRECATED(normalizer Normalizer) {}

func (d *DeduplicatorDescriptor) HasName() bool {
	return d.Name != nil
}

func (d *DeduplicatorDescriptor) HasNameMatch(name string) bool {
	return d.Name != nil && *d.Name == name
}
