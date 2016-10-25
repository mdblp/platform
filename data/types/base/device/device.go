package device

/* CHECKLIST
 * [x] Uses interfaces as appropriate
 * [x] Private package variables use underscore prefix
 * [x] All parameters validated
 * [x] All errors handled
 * [x] Reviewed for concurrency safety
 * [x] Code complete
 * [x] Full test coverage
 */

import (
	"github.com/tidepool-org/platform/app"
	"github.com/tidepool-org/platform/data"
	"github.com/tidepool-org/platform/data/types/base"
)

type Device struct {
	base.Base `bson:",inline"`

	SubType string `json:"subType,omitempty" bson:"subType,omitempty"`
}

type Meta struct {
	Type    string `json:"type,omitempty"`
	SubType string `json:"subType,omitempty"`
}

func Type() string {
	return "deviceEvent"
}

func (d *Device) Init() {
	d.Base.Init()
	d.Type = Type()

	d.SubType = ""
}

func (d *Device) Meta() interface{} {
	return &Meta{
		Type:    d.Type,
		SubType: d.SubType,
	}
}

func (d *Device) Parse(parser data.ObjectParser) error {
	parser.SetMeta(d.Meta())

	return d.Base.Parse(parser)
}

func (d *Device) Validate(validator data.Validator) error {
	validator.SetMeta(d.Meta())

	if err := d.Base.Validate(validator); err != nil {
		return err
	}

	validator.ValidateString("type", &d.Type).EqualTo(Type())

	validator.ValidateString("subType", &d.SubType).NotEmpty()

	return nil
}

func (d *Device) Normalize(normalizer data.Normalizer) error {
	normalizer.SetMeta(d.Meta())

	return d.Base.Normalize(normalizer)
}

func (d *Device) IdentityFields() ([]string, error) {
	identityFields, err := d.Base.IdentityFields()
	if err != nil {
		return nil, err
	}

	if d.SubType == "" {
		return nil, app.Error("device", "sub type is empty")
	}

	return append(identityFields, d.SubType), nil
}
