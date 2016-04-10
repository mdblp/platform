package basal

import (
	"reflect"

	"github.com/tidepool-org/platform/Godeps/_workspace/src/gopkg.in/bluesuncorp/validator.v8"
	"github.com/tidepool-org/platform/data/types"

	"github.com/tidepool-org/platform/validate"
)

func init() {
	types.GetPlatformValidator().RegisterValidation(durationField.Tag, DurationValidator)
	types.GetPlatformValidator().RegisterValidation(deliveryTypeField.Tag, DeliveryTypeValidator)
}

type Base struct {
	DeliveryType *string `json:"deliveryType" bson:"deliveryType" valid:"basaldeliverytype"`
	Duration     *int    `json:"duration,omitempty" bson:"duration,omitempty" valid:"omitempty,basalduration"`
	types.Base   `bson:",inline"`
}

type SuppressedBasal struct {
	Type         *string  `json:"type" bson:"type" valid:"required"`
	DeliveryType *string  `json:"deliveryType" bson:"deliveryType" valid:"basaldeliverytype"`
	ScheduleName *string  `json:"scheduleName" bson:"scheduleName" valid:"omitempty,required"`
	Rate         *float64 `json:"rate" bson:"rate" valid:"omitempty,basalrate"`
}

const Name = "basal"

var (
	deliveryTypeField = types.DatumFieldInformation{
		DatumField: &types.DatumField{Name: "deliveryType"},
		Tag:        "basaldeliverytype",
		Message:    "Must be one of scheduled, suspend, temp",
		Allowed:    types.Allowed{"scheduled": true, "suspend": true, "temp": true},
	}

	durationField = types.IntDatumField{
		DatumField:      &types.DatumField{Name: "duration"},
		Tag:             "basalduration",
		Message:         "Must be greater than 0",
		AllowedIntRange: &types.AllowedIntRange{LowerLimit: 0},
	}

	failureReasons = validate.ErrorReasons{
		deliveryTypeField.Tag: deliveryTypeField.Message,
		rateField.Tag:         rateField.Message,
		durationField.Tag:     durationField.Message,
		percentField.Tag:      percentField.Message,
	}
)

func makeSuppressed(datum types.Datum, errs validate.ErrorProcessing) *SuppressedBasal {
	return &SuppressedBasal{
		Type:         datum.ToString("type", errs),
		DeliveryType: datum.ToString(deliveryTypeField.Name, errs),
		ScheduleName: datum.ToString(scheduleNameField.Name, errs),
		Rate:         datum.ToFloat64(rateField.Name, errs),
	}
}

func Build(datum types.Datum, errs validate.ErrorProcessing) interface{} {

	base := &Base{
		DeliveryType: datum.ToString(deliveryTypeField.Name, errs),
		Duration:     datum.ToInt(durationField.Name, errs),
		Base:         types.BuildBase(datum, errs),
	}

	switch *base.DeliveryType {
	case "scheduled":
		return base.makeScheduled(datum, errs)
	case "suspend":
		return base.makeSuspend(datum, errs)
	case "temp":
		return base.makeTemporary(datum, errs)
	default:
		types.GetPlatformValidator().SetErrorReasons(failureReasons).Struct(base, errs)
		return base
	}
}

func DurationValidator(v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {
	duration, ok := field.Interface().(int)
	if !ok {
		return false
	}
	return duration > durationField.LowerLimit
}

func DeliveryTypeValidator(v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {
	deliveryType, ok := field.Interface().(string)
	if !ok {
		return false
	}
	_, ok = deliveryTypeField.Allowed[deliveryType]
	return ok
}