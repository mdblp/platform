package data

/* CHECKLIST
 * [x] Uses interfaces as appropriate
 * [x] Private package variables use underscore prefix
 * [x] All parameters validated
 * [x] All errors handled
 * [x] Reviewed for concurrency safety
 * [x] Code complete
 * [x] Full test coverage
 */

import "github.com/tidepool-org/platform/service"

type Normalizer interface {
	SetMeta(meta interface{})
	AppendError(reference interface{}, err *service.Error)

	AppendDatum(datum Datum)

	NormalizeBloodGlucose(units *string) BloodGlucoseNormalizer

	NewChildNormalizer(reference interface{}) Normalizer
}

type BloodGlucoseNormalizer interface {
	Units() *string
	Value(value *float64) *float64
	UnitsAndValue(value *float64) (*string, *float64)
}
