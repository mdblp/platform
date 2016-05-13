package ketone_test

import (
	"github.com/tidepool-org/platform/pvn/data/types/base/testing"
	"github.com/tidepool-org/platform/service"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
)

var _ = Describe("Blood Ketone", func() {

	var rawObject = testing.RawBaseObject()

	rawObject["type"] = "bloodKetone"
	rawObject["units"] = "mmol/L"
	rawObject["value"] = 5

	Context("units", func() {

		DescribeTable("units when", testing.ExpectFieldNotValid,
			Entry("empty", rawObject, "units", "", []*service.Error{&service.Error{}}),
			Entry("not one of the predefined values", rawObject, "units", "wrong", []*service.Error{&service.Error{}}),
		)

		DescribeTable("valid when", testing.ExpectFieldIsValid,
			Entry("mmol/l", rawObject, "units", "mmol/l"),
			Entry("mmol/L", rawObject, "units", "mmol/L"),
			Entry("mg/dl", rawObject, "units", "mg/dl"),
			Entry("mg/dL", rawObject, "units", "mg/dL"),
		)

	})

	Context("value", func() {

		DescribeTable("value when", testing.ExpectFieldNotValid,
			Entry("less than 0", rawObject, "value", -0.1, []*service.Error{&service.Error{}}),
			Entry("greater than 1000", rawObject, "value", 1000.1, []*service.Error{&service.Error{}}),
		)

		DescribeTable("valid when", testing.ExpectFieldIsValid,
			Entry("above 0", rawObject, "value", 0.1),
			Entry("below 1000", rawObject, "value", 999.99),
		)

	})

})