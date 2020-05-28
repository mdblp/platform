package prescriptor_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	dataNormalizer "github.com/tidepool-org/platform/data/normalizer"
	"github.com/tidepool-org/platform/data/types/bolus/prescriptor"
	dataTypesTest "github.com/tidepool-org/platform/data/types/test"
	errorsTest "github.com/tidepool-org/platform/errors/test"
	"github.com/tidepool-org/platform/pointer"
	"github.com/tidepool-org/platform/structure"
	structureValidator "github.com/tidepool-org/platform/structure/validator"
	"github.com/tidepool-org/platform/test"
)

func NewPrescriptor() *prescriptor.Prescriptor {
	datum := prescriptor.NewPrescriptor()
	datum.Prescriptor = pointer.FromString(test.RandomStringFromArray(prescriptor.Presciptors()))
	return datum
}

func ClonePrescriptor(datum *prescriptor.Prescriptor) *prescriptor.Prescriptor {
	if datum == nil {
		return nil
	}
	clone := prescriptor.NewPrescriptor()
	clone.Prescriptor = pointer.CloneString(datum.Prescriptor)
	return clone
}

var _ = Describe("Prescriptor", func() {
	Context("Addition of two digits", func() {

		It("Manual Prescriptor is expected", func() {
			Expect(prescriptor.ManualPrescriptor).To(Equal("manual"))
		})

		It("Auto Prescriptor is expected", func() {
			Expect(prescriptor.AutoPrescriptor).To(Equal("auto"))
		})

		It("Hybrid Prescriptor is expected", func() {
			Expect(prescriptor.HybridPrescriptor).To(Equal("hybrid"))
		})

		Context("NewPrescriptor", func() {
			It("is successful", func() {
				Expect(prescriptor.NewPrescriptor()).To(Equal(&prescriptor.Prescriptor{}))
			})
		})

	})
	Context("Prescriptor", func() {
		Context("Parse", func() {
			// TODO
		})

		Context("Validate", func() {
			DescribeTable("validates the datum",
				func(mutator func(datum *prescriptor.Prescriptor), expectedErrors ...error) {
					datum := NewPrescriptor()
					mutator(datum)
					dataTypesTest.ValidateWithExpectedOrigins(datum, structure.Origins(), expectedErrors...)
				},
				Entry("succeeds",
					func(datum *prescriptor.Prescriptor) {},
				),
				Entry("Valid value",
					func(datum *prescriptor.Prescriptor) {
						datum.Prescriptor = nil
					},
				),
				Entry("Invalid Prescriptor value",
					func(datum *prescriptor.Prescriptor) {
						datum.Prescriptor = pointer.FromString("invalid")
					},
					errorsTest.WithPointerSource(structureValidator.ErrorValueStringNotOneOf("invalid", []string{"auto", "manual", "hybrid"}), "/prescriptor"),
				),
			)
		})

		Context("Normalize", func() {
			DescribeTable("normalizes the datum",
				func(mutator func(datum *prescriptor.Prescriptor)) {
					for _, origin := range structure.Origins() {
						datum := NewPrescriptor()
						mutator(datum)
						expectedDatum := ClonePrescriptor(datum)
						normalizer := dataNormalizer.New()
						Expect(normalizer).ToNot(BeNil())
						datum.Normalize(normalizer.WithOrigin(origin))
						Expect(normalizer.Error()).To(BeNil())
						Expect(normalizer.Data()).To(BeEmpty())
						Expect(datum).To(Equal(expectedDatum))
					}
				},
				Entry("does not modify the datum",
					func(datum *prescriptor.Prescriptor) {},
				),
			)
		})
	})
})
