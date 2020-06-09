package common_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/tidepool-org/platform/data/types/common"
	dataTypesCommonTest "github.com/tidepool-org/platform/data/types/common/test"

	dataNormalizer "github.com/tidepool-org/platform/data/normalizer"
	dataTypesTest "github.com/tidepool-org/platform/data/types/test"
	errorsTest "github.com/tidepool-org/platform/errors/test"
	"github.com/tidepool-org/platform/pointer"
	"github.com/tidepool-org/platform/structure"
	structureValidator "github.com/tidepool-org/platform/structure/validator"
)

var _ = Describe("EventType", func() {

	Context("NewEventType", func() {
		It("is successful", func() {
			Expect(common.NewEventType()).To(Equal(&common.EventType{}))
		})

		It("Events returns expected", func() {
			Expect(common.Events()).To(Equal([]string{"start", "stop"}))
		})
	})

	Context("EventType", func() {
		Context("Parse", func() {
			// TODO
		})

		Context("Validate", func() {
			DescribeTable("validates the datum",
				func(mutator func(datum *common.EventType), expectedErrors ...error) {
					datum := dataTypesCommonTest.NewEventType()
					mutator(datum)
					dataTypesTest.ValidateWithExpectedOrigins(datum, structure.Origins(), expectedErrors...)
				},
				Entry("succeeds",
					func(datum *common.EventType) {},
				),
				Entry("Valid eventType, start",
					func(datum *common.EventType) {
						datum.EventType = pointer.FromString(common.StartEvent)
					},
				),
				Entry("Valid eventType, stop",
					func(datum *common.EventType) {
						datum.EventType = pointer.FromString(common.StopEvent)
					},
				),
				Entry("invalid eventType",
					func(datum *common.EventType) {
						datum.EventType = pointer.FromString("invalid")
					},
					errorsTest.WithPointerSource(structureValidator.ErrorValueStringNotOneOf("invalid", common.Events()), "/eventType"),
				),
			)
		})

		Context("Normalize", func() {
			DescribeTable("normalizes the datum",
				func(mutator func(datum *common.EventType)) {
					for _, origin := range structure.Origins() {
						datum := dataTypesCommonTest.NewEventType()
						mutator(datum)
						expectedDatum := dataTypesCommonTest.CloneEventType(datum)
						normalizer := dataNormalizer.New()
						Expect(normalizer).ToNot(BeNil())
						datum.Normalize(normalizer.WithOrigin(origin))
						Expect(normalizer.Error()).To(BeNil())
						Expect(normalizer.Data()).To(BeEmpty())
						Expect(datum).To(Equal(expectedDatum))
					}
				},
				Entry("does not modify the datum",
					func(datum *common.EventType) {},
				),
			)
		})
	})
})
