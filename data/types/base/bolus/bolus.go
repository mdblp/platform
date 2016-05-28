package bolus

/* CHECKLIST
 * [ ] Uses interfaces as appropriate
 * [ ] Private package variables use underscore prefix
 * [ ] All parameters validated
 * [ ] All errors handled
 * [ ] Reviewed for concurrency safety
 * [ ] Code complete
 * [ ] Full test coverage
 */

import (
	"github.com/tidepool-org/platform/app"
	"github.com/tidepool-org/platform/data"
	"github.com/tidepool-org/platform/data/types/base"
)

type Bolus struct {
	base.Base `bson:",inline"`

	SubType string `json:"subType,omitempty" bson:"subType,omitempty"`
}

type Meta struct {
	Type    string `json:"type,omitempty"`
	SubType string `json:"subType,omitempty"`
}

func Type() string {
	return "bolus"
}

func New(subType string) (*Bolus, error) {
	if subType == "" {
		return nil, app.Error("basal", "sub type is missing")
	}

	bolusBase, err := base.New(Type())
	if err != nil {
		return nil, err
	}

	return &Bolus{
		Base:    *bolusBase,
		SubType: subType,
	}, nil
}

func (b *Bolus) Meta() interface{} {
	return &Meta{
		Type:    b.Type,
		SubType: b.SubType,
	}
}

func (b *Bolus) Parse(parser data.ObjectParser) error {
	parser.SetMeta(b.Meta())

	return b.Base.Parse(parser)
}

func (b *Bolus) Validate(validator data.Validator) error {
	validator.SetMeta(b.Meta())

	return b.Base.Validate(validator)
}

func (b *Bolus) Normalize(normalizer data.Normalizer) error {
	normalizer.SetMeta(b.Meta())

	return b.Base.Normalize(normalizer)
}
