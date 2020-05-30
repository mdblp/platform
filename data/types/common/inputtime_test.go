package commontypes_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	commontypes "github.com/tidepool-org/platform/data/types/common"
	commontypesTest "github.com/tidepool-org/platform/data/types/common/test"
	"github.com/tidepool-org/platform/test"

	dataNormalizer "github.com/tidepool-org/platform/data/normalizer"
	dataTypesTest "github.com/tidepool-org/platform/data/types/test"
	errorsTest "github.com/tidepool-org/platform/errors/test"
	"github.com/tidepool-org/platform/pointer"
	"github.com/tidepool-org/platform/structure"
	structureValidator "github.com/tidepool-org/platform/structure/validator"
)

func NewInputTime() *commontypes.InputTime {
	datum := commontypesTest.NewInputTime()
	timeReference := test.RandomTime()
	datum.InputTime = pointer.FromString(timeReference.Format(time.RFC3339Nano))
	return datum
}

func CloneInputTime(datum *commontypes.InputTime) *commontypes.InputTime {
	if datum == nil {
		return nil
	}
	clone := commontypesTest.NewInputTime()
	clone.InputTime = pointer.CloneString(datum.InputTime)
	return clone
}

var _ = Describe("InputTime", func() {

	Context("NewInputTime", func() {
		It("is successful", func() {
			Expect(commontypes.NewInputTime()).To(Equal(&commontypes.InputTime{}))
		})
	})

	Context("InputTime", func() {
		Context("Parse", func() {
			// TODO
		})

		Context("Validate", func() {
			DescribeTable("validates the datum",
				func(mutator func(datum *commontypes.InputTime), expectedErrors ...error) {
					datum := NewInputTime()
					mutator(datum)
					dataTypesTest.ValidateWithExpectedOrigins(datum, structure.Origins(), expectedErrors...)
				},
				Entry("succeeds",
					func(datum *commontypes.InputTime) {},
				),
				Entry("Valid inputTime",
					func(datum *commontypes.InputTime) {
						datum.InputTime = pointer.FromString(test.RandomTime().Format(time.RFC3339Nano))
					},
				),
				Entry("invalid inputTime",
					func(datum *commontypes.InputTime) {
						datum.InputTime = pointer.FromString("invalid")
					},
					errorsTest.WithPointerSource(structureValidator.ErrorValueStringAsTimeNotValid("invalid", time.RFC3339Nano), "/inputTime"),
				),
			)
		})

		Context("Normalize", func() {
			DescribeTable("normalizes the datum",
				func(mutator func(datum *commontypes.InputTime)) {
					for _, origin := range structure.Origins() {
						datum := NewInputTime()
						mutator(datum)
						expectedDatum := CloneInputTime(datum)
						normalizer := dataNormalizer.New()
						Expect(normalizer).ToNot(BeNil())
						datum.Normalize(normalizer.WithOrigin(origin))
						Expect(normalizer.Error()).To(BeNil())
						Expect(normalizer.Data()).To(BeEmpty())
						Expect(datum).To(Equal(expectedDatum))
					}
				},
				Entry("does not modify the datum",
					func(datum *commontypes.InputTime) {},
				),
			)
		})
	})
})
