package mode_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	dataNormalizer "github.com/tidepool-org/platform/data/normalizer"
	dataTypesCommonTest "github.com/tidepool-org/platform/data/types/common/test"
	"github.com/tidepool-org/platform/data/types/device"
	"github.com/tidepool-org/platform/data/types/device/mode"
	dataTypesDeviceTest "github.com/tidepool-org/platform/data/types/device/test"
	dataTypesTest "github.com/tidepool-org/platform/data/types/test"
	errorsTest "github.com/tidepool-org/platform/errors/test"
	"github.com/tidepool-org/platform/pointer"
	"github.com/tidepool-org/platform/structure"
	structureValidator "github.com/tidepool-org/platform/structure/validator"
	"github.com/tidepool-org/platform/test"
)

func NewMeta() interface{} {
	return &device.Meta{
		Type:    "deviceEvent",
		SubType: "zen",
	}
}

func NewMode() *mode.Mode {
	datum := mode.New(mode.ZenMode)
	datum.Device = *dataTypesDeviceTest.NewDevice()
	datum.SubType = mode.ZenMode
	datum.EventID = pointer.FromString("123456789")
	datum.Duration = dataTypesCommonTest.NewDuration()
	datum.InputTime = dataTypesCommonTest.NewInputTime()
	return datum
}

func CloneMode(datum *mode.Mode) *mode.Mode {
	if datum == nil {
		return nil
	}
	clone := mode.New(datum.SubType)
	clone.Device = *dataTypesDeviceTest.CloneDevice(&datum.Device)
	clone.EventID = pointer.FromString("123456789")
	clone.Duration = dataTypesCommonTest.CloneDuration(datum.Duration)
	clone.InputTime = dataTypesCommonTest.CloneInputTime(datum.InputTime)
	return clone
}

var _ = Describe("Change", func() {
	It("SubType is expected", func() {
		Expect(mode.ConfidentialMode).To(Equal("confidential"))
		Expect(mode.ZenMode).To(Equal("zen"))
		Expect(mode.Warmup).To(Equal("warmup"))
	})

	Context("New", func() {
		It("returns the expected datum with all Zen values initialized", func() {
			datum := mode.New(mode.ZenMode)
			Expect(datum).ToNot(BeNil())
			Expect(datum.Type).To(Equal("deviceEvent"))
			Expect(datum.SubType).To(Equal("zen"))
			Expect(datum.InputTime.InputTime).To(BeNil())
		})
		It("returns the expected datum with all confidential values initialized", func() {
			datum := mode.New(mode.ConfidentialMode)
			Expect(datum).ToNot(BeNil())
			Expect(datum.Type).To(Equal("deviceEvent"))
			Expect(datum.SubType).To(Equal("confidential"))
			Expect(datum.InputTime.InputTime).To(BeNil())
		})
		It("returns the expected datum with all warmup values initialized", func() {
			datum := mode.New(mode.Warmup)
			Expect(datum).ToNot(BeNil())
			Expect(datum.Type).To(Equal("deviceEvent"))
			Expect(datum.SubType).To(Equal("warmup"))
			Expect(datum.InputTime.InputTime).To(BeNil())
		})
	})

	Context("Mode", func() {
		Context("Parse", func() {
			// TODO
		})

		Context("Validate", func() {
			DescribeTable("validates the datum",
				func(mutator func(datum *mode.Mode), expectedErrors ...error) {
					datum := NewMode()
					mutator(datum)
					dataTypesTest.ValidateWithExpectedOrigins(datum, structure.Origins(), expectedErrors...)
				},
				Entry("succeeds",
					func(datum *mode.Mode) {},
				),
				Entry("type missing",
					func(datum *mode.Mode) { datum.Type = "" },
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueEmpty(), "/type", &device.Meta{SubType: "zen"}),
				),
				Entry("type invalid",
					func(datum *mode.Mode) { datum.Type = "invalidType" },
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueNotEqualTo("invalidType", "deviceEvent"), "/type", &device.Meta{Type: "invalidType", SubType: "zen"}),
				),
				Entry("type device",
					func(datum *mode.Mode) { datum.Type = "deviceEvent" },
				),
				Entry("sub type missing",
					func(datum *mode.Mode) { datum.SubType = "" },
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueEmpty(), "/subType", &device.Meta{Type: "deviceEvent"}),
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueStringNotOneOf("", []string{"confidential", "zen", "warmup", "loopMode"}), "/subType", &device.Meta{Type: "deviceEvent", SubType: ""}),
				),
				Entry("sub type invalid",
					func(datum *mode.Mode) { datum.SubType = "invalidSubType" },
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueStringNotOneOf("invalidSubType", []string{"confidential", "zen", "warmup", "loopMode"}), "/subType", &device.Meta{Type: "deviceEvent", SubType: "invalidSubType"}),
				),
				Entry("sub type zen",
					func(datum *mode.Mode) { datum.SubType = "zen" },
				),
				Entry("multiple errors",
					func(datum *mode.Mode) {
						datum.Type = "invalidType"
						datum.SubType = "invalidSubType"
					},
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueNotEqualTo("invalidType", "deviceEvent"), "/type", &device.Meta{Type: "invalidType", SubType: "invalidSubType"}),
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueStringNotOneOf("invalidSubType", []string{"confidential", "zen", "warmup", "loopMode"}), "/subType", &device.Meta{Type: "invalidType", SubType: "invalidSubType"}),
				),
				Entry("EventId is missing",
					func(datum *mode.Mode) {
						datum.EventID = nil
					},
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueNotExists(), "/eventId", &device.Meta{Type: "deviceEvent", SubType: "zen"}),
				),
				Entry("EventId is missing",
					func(datum *mode.Mode) {
						datum.SubType = "confidential"
						datum.EventID = nil
					},
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueNotExists(), "/eventId", &device.Meta{Type: "deviceEvent", SubType: "confidential"}),
				),
				Entry("EventId is missing",
					func(datum *mode.Mode) {
						datum.SubType = "warmup"
						datum.EventID = nil
					},
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueNotExists(), "/eventId", &device.Meta{Type: "deviceEvent", SubType: "warmup"}),
				),
				Entry("inputTime is missing",
					func(datum *mode.Mode) {
						datum.SubType = "confidential"
						datum.InputTime.InputTime = nil
					},
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueNotExists(), "/inputTime", &device.Meta{Type: "deviceEvent", SubType: "confidential"}),
				),
				Entry("Valid inputTime",
					func(datum *mode.Mode) {
						datum.InputTime.InputTime = pointer.FromString(test.RandomTime().Format(time.RFC3339Nano))
					},
				),
				Entry("InputTime invalid",
					func(datum *mode.Mode) {
						datum.InputTime.InputTime = pointer.FromString("invalid")
					},
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueStringAsTimeNotValid("invalid", time.RFC3339Nano), "/inputTime", NewMeta()),
				),
				Entry("duration missing",
					func(datum *mode.Mode) {
						datum.Duration = nil
					},
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueNotExists(), "/duration", &device.Meta{Type: "deviceEvent", SubType: "zen"}),
				),
				Entry("duration missing",
					func(datum *mode.Mode) {
						datum.SubType = "warmup"
						datum.Duration = nil
					},
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueNotExists(), "/duration", &device.Meta{Type: "deviceEvent", SubType: "warmup"}),
				),
				Entry("duration missing",
					func(datum *mode.Mode) {
						datum.SubType = "confidential"
						datum.Duration = nil
					},
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueNotExists(), "/duration", &device.Meta{Type: "deviceEvent", SubType: "confidential"}),
				),
				Entry("succeeds when duration is missing and sub type is loopMode",
					func(datum *mode.Mode) {
						datum.SubType = "loopMode"
						datum.Duration = nil
					},
				),
				Entry("sub type invalid",
					func(datum *mode.Mode) { datum.SubType = "invalidSubType" },
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueStringNotOneOf("invalidSubType", []string{"confidential", "zen", "warmup", "loopMode"}), "/subType", &device.Meta{Type: "deviceEvent", SubType: "invalidSubType"}),
				),
			)

			DescribeTable("validates the datum with origin internal/store",
				func(mutator func(datum *mode.Mode), expectedErrors ...error) {
					datum := NewMode()
					mutator(datum)
					dataTypesTest.ValidateWithOrigin(datum, structure.OriginInternal, expectedErrors...)
					dataTypesTest.ValidateWithOrigin(datum, structure.OriginStore, expectedErrors...)
				},
				Entry("succeeds",
					func(datum *mode.Mode) {},
				),
			)
		})

		Context("Normalize", func() {
			DescribeTable("normalizes the datum",
				func(mutator func(datum *mode.Mode)) {
					for _, origin := range structure.Origins() {
						datum := NewMode()
						mutator(datum)
						expectedDatum := CloneMode(datum)
						normalizer := dataNormalizer.New()
						Expect(normalizer).ToNot(BeNil())
						datum.Normalize(normalizer.WithOrigin(origin))
						Expect(normalizer.Error()).To(BeNil())
						Expect(normalizer.Data()).To(BeEmpty())
						Expect(datum).To(Equal(expectedDatum))
					}
				},
				Entry("does not modify the datum",
					func(datum *mode.Mode) {},
				),
			)
		})
	})
})
