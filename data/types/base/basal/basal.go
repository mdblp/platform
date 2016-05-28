package basal

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

type Basal struct {
	base.Base `bson:",inline"`

	DeliveryType string `json:"deliveryType,omitempty" bson:"deliveryType,omitempty"`
}

type Meta struct {
	Type         string `json:"type,omitempty"`
	DeliveryType string `json:"deliveryType,omitempty"`
}

func Type() string {
	return "basal"
}

func New(deliveryType string) (*Basal, error) {
	if deliveryType == "" {
		return nil, app.Error("basal", "delivery type is missing")
	}

	basalBase, err := base.New(Type())
	if err != nil {
		return nil, err
	}

	return &Basal{
		Base:         *basalBase,
		DeliveryType: deliveryType,
	}, nil
}

func (b *Basal) Meta() interface{} {
	return &Meta{
		Type:         b.Type,
		DeliveryType: b.DeliveryType,
	}
}

func (b *Basal) Parse(parser data.ObjectParser) error {
	parser.SetMeta(b.Meta())

	return b.Base.Parse(parser)
}

func (b *Basal) Validate(validator data.Validator) error {
	validator.SetMeta(b.Meta())

	return b.Base.Validate(validator)
}

func (b *Basal) Normalize(normalizer data.Normalizer) error {
	normalizer.SetMeta(b.Meta())

	return b.Base.Normalize(normalizer)
}
