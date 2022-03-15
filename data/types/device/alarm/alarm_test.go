package alarm_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/tidepool-org/platform/data"
	dataNormalizer "github.com/tidepool-org/platform/data/normalizer"
	dataTest "github.com/tidepool-org/platform/data/test"
	"github.com/tidepool-org/platform/data/types"
	"github.com/tidepool-org/platform/data/types/device"
	"github.com/tidepool-org/platform/data/types/device/alarm"
	dataTypesDeviceStatus "github.com/tidepool-org/platform/data/types/device/status"
	dataTypesDeviceStatusTest "github.com/tidepool-org/platform/data/types/device/status/test"
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
		SubType: "alarm",
	}
}

func NewAlarm() *alarm.Alarm {
	datum := alarm.New()
	datum.Device = *dataTypesDeviceTest.NewDevice()
	datum.SubType = "alarm"
	datum.AlarmType = pointer.FromString(test.RandomStringFromArray(alarm.LegacyAlarmTypes()))
	return datum
}

func NewAlarmWithStatus() *alarm.Alarm {
	var status data.Datum
	status = dataTypesDeviceStatusTest.NewStatus()
	datum := NewAlarm()
	datum.Status = &status
	return datum
}

func NewAlarmWithStatusID() *alarm.Alarm {
	datum := NewAlarm()
	datum.StatusID = pointer.FromString(dataTest.RandomID())
	return datum
}

func NewAlarmFromHandset() *alarm.Alarm {
	datum := NewAlarm()
	datum.GUID = pointer.FromString("ID123456789")
	datum.AlarmType = pointer.FromString(alarm.AlarmTypeHandset)
	datum.AlarmLevel = pointer.FromString(test.RandomStringFromArray(alarm.AlarmLevels()))
	datum.AlarmCode = pointer.FromString("code123")
	datum.AlarmLabel = pointer.FromString("label")
	datum.AckStatus = pointer.FromString(test.RandomStringFromArray(alarm.AckStatuses()))
	datum.UpdateTime = pointer.CloneString(datum.Time)
	return datum
}

func CloneAlarm(datum *alarm.Alarm) *alarm.Alarm {
	if datum == nil {
		return nil
	}
	clone := alarm.New()
	clone.Device = *dataTypesDeviceTest.CloneDevice(&datum.Device)
	clone.AlarmType = pointer.CloneString(datum.AlarmType)
	if datum.Status != nil {
		switch status := (*datum.Status).(type) {
		case *dataTypesDeviceStatus.Status:
			clone.Status = data.DatumAsPointer(dataTypesDeviceStatusTest.CloneStatus(status))
		}
	}
	clone.GUID = pointer.CloneString(datum.GUID)
	clone.StatusID = pointer.CloneString(datum.StatusID)
	clone.AlarmLevel = pointer.CloneString(datum.AlarmLevel)
	clone.AlarmCode = pointer.CloneString(datum.AlarmCode)
	clone.AlarmLabel = pointer.CloneString(datum.AlarmLabel)
	clone.AckStatus = pointer.CloneString(datum.AckStatus)
	clone.UpdateTime = pointer.CloneString(datum.UpdateTime)
	return clone
}

var _ = Describe("Change", func() {
	It("SubType is expected", func() {
		Expect(alarm.SubType).To(Equal("alarm"))
	})

	It("AlarmTypeAutoOff is expected", func() {
		Expect(alarm.AlarmTypeAutoOff).To(Equal("auto_off"))
	})

	It("AlarmTypeLowInsulin is expected", func() {
		Expect(alarm.AlarmTypeLowInsulin).To(Equal("low_insulin"))
	})

	It("AlarmTypeLowPower is expected", func() {
		Expect(alarm.AlarmTypeLowPower).To(Equal("low_power"))
	})

	It("AlarmTypeNoDelivery is expected", func() {
		Expect(alarm.AlarmTypeNoDelivery).To(Equal("no_delivery"))
	})

	It("AlarmTypeNoInsulin is expected", func() {
		Expect(alarm.AlarmTypeNoInsulin).To(Equal("no_insulin"))
	})

	It("AlarmTypeNoPower is expected", func() {
		Expect(alarm.AlarmTypeNoPower).To(Equal("no_power"))
	})

	It("AlarmTypeOcclusion is expected", func() {
		Expect(alarm.AlarmTypeOcclusion).To(Equal("occlusion"))
	})

	It("AlarmTypeOther is expected", func() {
		Expect(alarm.AlarmTypeOther).To(Equal("other"))
	})

	It("AlarmTypeOverLimit is expected", func() {
		Expect(alarm.AlarmTypeOverLimit).To(Equal("over_limit"))
	})

	It("IsAnAlarm is expected", func() {
		Expect(alarm.AlarmTypeOverLimit).To(Equal("over_limit"))
	})

	It("isAnAlarm is expected", func() {
		Expect(alarm.IsAnAlarm).To(Equal("alarm"))
	})

	It("isAnAlert is expected", func() {
		Expect(alarm.IsAnAlert).To(Equal("alert"))
	})

	It("Legacy AlarmTypes returns expected", func() {
		Expect(alarm.LegacyAlarmTypes()).To(Equal([]string{"auto_off", "low_insulin", "low_power", "no_delivery", "no_insulin", "no_power", "occlusion", "other", "over_limit"}))
	})

	It("AlarmTypes returns expected", func() {
		Expect(alarm.AlarmTypes()).To(Equal([]string{"auto_off", "low_insulin", "low_power", "no_delivery", "no_insulin", "no_power", "occlusion", "other", "over_limit", "handset"}))
	})

	Context("New", func() {
		It("returns the expected datum with all values initialized", func() {
			datum := alarm.New()
			Expect(datum).ToNot(BeNil())
			Expect(datum.Type).To(Equal("deviceEvent"))
			Expect(datum.SubType).To(Equal("alarm"))
			Expect(datum.AlarmType).To(BeNil())
			Expect(datum.Status).To(BeNil())
			Expect(datum.StatusID).To(BeNil())
			Expect(datum.AlarmLevel).To(BeNil())
			Expect(datum.AlarmCode).To(BeNil())
			Expect(datum.AlarmLabel).To(BeNil())
		})
	})

	Context("Alarm", func() {
		Context("Parse", func() {
			// TODO
		})

		Context("Validate", func() {
			DescribeTable("validates the datum",
				func(mutator func(datum *alarm.Alarm), expectedErrors ...error) {
					datum := NewAlarm()
					mutator(datum)
					dataTypesTest.ValidateWithExpectedOrigins(datum, structure.Origins(), expectedErrors...)
				},
				Entry("succeeds",
					func(datum *alarm.Alarm) {},
				),
				Entry("type missing",
					func(datum *alarm.Alarm) { datum.Type = "" },
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueEmpty(), "/type", &device.Meta{SubType: "alarm"}),
				),
				Entry("type invalid",
					func(datum *alarm.Alarm) { datum.Type = "invalidType" },
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueNotEqualTo("invalidType", "deviceEvent"), "/type", &device.Meta{Type: "invalidType", SubType: "alarm"}),
				),
				Entry("type device",
					func(datum *alarm.Alarm) { datum.Type = "deviceEvent" },
				),
				Entry("sub type missing",
					func(datum *alarm.Alarm) { datum.SubType = "" },
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueEmpty(), "/subType", &device.Meta{Type: "deviceEvent"}),
				),
				Entry("sub type invalid",
					func(datum *alarm.Alarm) { datum.SubType = "invalidSubType" },
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueNotEqualTo("invalidSubType", "alarm"), "/subType", &device.Meta{Type: "deviceEvent", SubType: "invalidSubType"}),
				),
				Entry("sub type alarm",
					func(datum *alarm.Alarm) { datum.SubType = "alarm" },
				),
				Entry("alarm type missing",
					func(datum *alarm.Alarm) { datum.AlarmType = nil },
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueNotExists(), "/alarmType", NewMeta()),
				),
				Entry("alarm type invalid",
					func(datum *alarm.Alarm) { datum.AlarmType = pointer.FromString("invalid") },
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueStringNotOneOf("invalid", []string{"auto_off", "low_insulin", "low_power", "no_delivery", "no_insulin", "no_power", "occlusion", "other", "over_limit", "handset"}), "/alarmType", NewMeta()),
				),
				Entry("alarm type auto_off",
					func(datum *alarm.Alarm) { datum.AlarmType = pointer.FromString("auto_off") },
				),
				Entry("alarm type low_insulin",
					func(datum *alarm.Alarm) { datum.AlarmType = pointer.FromString("low_insulin") },
				),
				Entry("alarm type low_power",
					func(datum *alarm.Alarm) { datum.AlarmType = pointer.FromString("low_power") },
				),
				Entry("alarm type no_delivery",
					func(datum *alarm.Alarm) { datum.AlarmType = pointer.FromString("no_delivery") },
				),
				Entry("alarm type no_insulin",
					func(datum *alarm.Alarm) { datum.AlarmType = pointer.FromString("no_insulin") },
				),
				Entry("alarm type no_power",
					func(datum *alarm.Alarm) { datum.AlarmType = pointer.FromString("no_power") },
				),
				Entry("alarm type occlusion",
					func(datum *alarm.Alarm) { datum.AlarmType = pointer.FromString("occlusion") },
				),
				Entry("alarm type other",
					func(datum *alarm.Alarm) { datum.AlarmType = pointer.FromString("other") },
				),
				Entry("alarm type over_limit",
					func(datum *alarm.Alarm) { datum.AlarmType = pointer.FromString("over_limit") },
				),
				Entry("multiple errors",
					func(datum *alarm.Alarm) {
						datum.Type = "invalidType"
						datum.SubType = "invalidSubType"
						datum.AlarmType = pointer.FromString("invalid")
					},
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueNotEqualTo("invalidType", "deviceEvent"), "/type", &device.Meta{Type: "invalidType", SubType: "invalidSubType"}),
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueNotEqualTo("invalidSubType", "alarm"), "/subType", &device.Meta{Type: "invalidType", SubType: "invalidSubType"}),
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueStringNotOneOf("invalid", []string{"auto_off", "low_insulin", "low_power", "no_delivery", "no_insulin", "no_power", "occlusion", "other", "over_limit", "handset"}), "/alarmType", &device.Meta{Type: "invalidType", SubType: "invalidSubType"}),
				),
			)

			DescribeTable("validates the datum with origin external",
				func(mutator func(datum *alarm.Alarm), expectedErrors ...error) {
					datum := NewAlarmWithStatus()
					mutator(datum)
					dataTypesTest.ValidateWithOrigin(datum, structure.OriginExternal, expectedErrors...)
				},
				Entry("succeeds",
					func(datum *alarm.Alarm) {},
				),
				Entry("status missing",
					func(datum *alarm.Alarm) { datum.Status = nil },
				),
				Entry("status valid",
					func(datum *alarm.Alarm) {
						datum.Status = data.DatumAsPointer(dataTypesDeviceStatusTest.NewStatus())
					},
				),
				Entry("status id missing",
					func(datum *alarm.Alarm) { datum.StatusID = nil },
				),
				Entry("status id exists",
					func(datum *alarm.Alarm) { datum.StatusID = pointer.FromString(dataTest.RandomID()) },
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueExists(), "/statusId", NewMeta()),
				),
				Entry("multiple errors",
					func(datum *alarm.Alarm) {
						datum.StatusID = pointer.FromString(dataTest.RandomID())
					},
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueExists(), "/statusId", NewMeta()),
				),
			)

			DescribeTable("validates the datum with origin internal/store",
				func(mutator func(datum *alarm.Alarm), expectedErrors ...error) {
					datum := NewAlarmWithStatusID()
					mutator(datum)
					dataTypesTest.ValidateWithOrigin(datum, structure.OriginInternal, expectedErrors...)
					dataTypesTest.ValidateWithOrigin(datum, structure.OriginStore, expectedErrors...)
				},
				Entry("succeeds",
					func(datum *alarm.Alarm) {},
				),
				Entry("status missing",
					func(datum *alarm.Alarm) { datum.Status = nil },
				),
				Entry("status exists",
					func(datum *alarm.Alarm) {
						datum.Status = data.DatumAsPointer(dataTypesDeviceStatusTest.NewStatus())
					},
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueExists(), "/status", NewMeta()),
				),
				Entry("status id missing",
					func(datum *alarm.Alarm) { datum.StatusID = nil },
				),
				Entry("status id invalid",
					func(datum *alarm.Alarm) { datum.StatusID = pointer.FromString("invalid") },
					errorsTest.WithPointerSourceAndMeta(data.ErrorValueStringAsIDNotValid("invalid"), "/statusId", NewMeta()),
				),
				Entry("status id valid",
					func(datum *alarm.Alarm) { datum.StatusID = pointer.FromString(dataTest.RandomID()) },
				),
				Entry("multiple errors",
					func(datum *alarm.Alarm) {
						datum.Status = data.DatumAsPointer(dataTypesDeviceStatusTest.NewStatus())
						datum.StatusID = pointer.FromString("invalid")
					},
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueExists(), "/status", NewMeta()),
					errorsTest.WithPointerSourceAndMeta(data.ErrorValueStringAsIDNotValid("invalid"), "/statusId", NewMeta()),
				),
			)

			DescribeTable("validates the datum with handset alarms",
				func(mutator func(datum *alarm.Alarm), expectedErrors ...error) {
					datum := NewAlarmFromHandset()
					mutator(datum)
					dataTypesTest.ValidateWithExpectedOrigins(datum, structure.Origins(), expectedErrors...)
				},
				Entry("succeeds",
					func(datum *alarm.Alarm) {},
				),
				Entry("GUID is missing",
					func(datum *alarm.Alarm) {
						datum.GUID = nil
					},
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueNotExists(), "/guid", NewMeta()),
				),
				Entry("GUID length out of range",
					func(datum *alarm.Alarm) {
						datum.GUID = pointer.FromString(test.RandomStringFromRange(65, 65))
					},
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorLengthNotLessThanOrEqualTo(65, 64), "/guid", NewMeta()),
				),
				Entry("GUID length in range",
					func(datum *alarm.Alarm) {
						datum.GUID = pointer.FromString(test.RandomStringFromRange(64, 64))
					},
				),
				Entry("invalid alarm level",
					func(datum *alarm.Alarm) {
						datum.AlarmLevel = pointer.FromString("invalid")
					},
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueStringNotOneOf("invalid", []string{"alarm", "alert"}), "/alarmLevel", NewMeta()),
				),
				Entry("alarm level is missing",
					func(datum *alarm.Alarm) {
						datum.AlarmLevel = nil
					},
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueNotExists(), "/alarmLevel", NewMeta()),
				),
				Entry("alarm code is missing",
					func(datum *alarm.Alarm) {
						datum.AlarmCode = nil
					},
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueNotExists(), "/alarmCode", NewMeta()),
				),
				Entry("Alarm Code length out of range",
					func(datum *alarm.Alarm) {
						datum.AlarmCode = pointer.FromString(test.RandomStringFromRange(65, 65))
					},
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorLengthNotLessThanOrEqualTo(65, 64), "/alarmCode", NewMeta()),
				),
				Entry("Alarm code length in range",
					func(datum *alarm.Alarm) {
						datum.AlarmCode = pointer.FromString(test.RandomStringFromRange(64, 64))
					},
				),
				Entry("alarm label is missing",
					func(datum *alarm.Alarm) {
						datum.AlarmLabel = nil
					},
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueNotExists(), "/alarmLabel", NewMeta()),
				),
				Entry("Alarm label out of range",
					func(datum *alarm.Alarm) {
						datum.AlarmLabel = pointer.FromString(test.RandomStringFromRange(257, 257))
					},
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorLengthNotLessThanOrEqualTo(257, 256), "/alarmLabel", NewMeta()),
				),
				Entry("Alarm label in range",
					func(datum *alarm.Alarm) {
						datum.AlarmLabel = pointer.FromString(test.RandomStringFromRange(256, 256))
					},
				),
				Entry("ackStatus is missing",
					func(datum *alarm.Alarm) {
						datum.AckStatus = nil
					},
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueNotExists(), "/ackStatus", NewMeta()),
				),
				Entry("updateTime is missing",
					func(datum *alarm.Alarm) {
						datum.UpdateTime = nil
					},
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueNotExists(), "/updateTime", NewMeta()),
				),
				Entry("Mulitple missing",
					func(datum *alarm.Alarm) {
						datum.GUID = nil
						datum.AlarmLevel = nil
						datum.AlarmCode = nil
						datum.AlarmLabel = nil
						datum.AckStatus = nil
						datum.UpdateTime = nil
					},
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueNotExists(), "/guid", NewMeta()),
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueNotExists(), "/alarmLevel", NewMeta()),
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueNotExists(), "/alarmCode", NewMeta()),
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueNotExists(), "/alarmLabel", NewMeta()),
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueNotExists(), "/ackStatus", NewMeta()),
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueNotExists(), "/updateTime", NewMeta()),
				),
				Entry("updateTime is invalid",
					func(datum *alarm.Alarm) {
						datum.UpdateTime = pointer.FromString("invalid")
					},
					errorsTest.WithPointerSourceAndMeta(structureValidator.ErrorValueStringAsTimeNotValid("invalid", types.TimeFormat), "/updateTime", NewMeta()),
				),
			)
		})

		Context("Normalize", func() {
			It("does not modify datum if status is missing", func() {
				datum := NewAlarmWithStatusID()
				expectedDatum := CloneAlarm(datum)
				normalizer := dataNormalizer.New()
				Expect(normalizer).ToNot(BeNil())
				datum.Normalize(normalizer)
				Expect(normalizer.Error()).To(BeNil())
				Expect(normalizer.Data()).To(BeEmpty())
				Expect(datum).To(Equal(expectedDatum))
			})

			It("normalizes the datum and replaces status with status id", func() {
				datumStatus := dataTypesDeviceStatusTest.NewStatus()
				datum := NewAlarmWithStatusID()
				datum.Status = data.DatumAsPointer(datumStatus)
				expectedDatum := CloneAlarm(datum)
				normalizer := dataNormalizer.New()
				Expect(normalizer).ToNot(BeNil())
				datum.Normalize(normalizer)
				Expect(normalizer.Error()).To(BeNil())
				Expect(normalizer.Data()).To(Equal([]data.Datum{datumStatus}))
				expectedDatum.Status = nil
				expectedDatum.StatusID = pointer.FromString(*datumStatus.ID)
				Expect(datum).To(Equal(expectedDatum))
			})

			It("does not modify datum if handset and status missing", func() {
				datum := NewAlarmFromHandset()
				datum.StatusID = pointer.FromString(dataTest.RandomID())
				expectedDatum := CloneAlarm(datum)
				normalizer := dataNormalizer.New()
				Expect(normalizer).ToNot(BeNil())
				datum.Normalize(normalizer)
				Expect(normalizer.Error()).To(BeNil())
				Expect(normalizer.Data()).To(BeEmpty())
				Expect(datum).To(Equal(expectedDatum))
			})

			It("normalizes the datum if handset and replaces status with status id", func() {
				datumStatus := dataTypesDeviceStatusTest.NewStatus()
				datum := NewAlarmFromHandset()
				datum.StatusID = pointer.FromString(dataTest.RandomID())
				datum.Status = data.DatumAsPointer(datumStatus)
				expectedDatum := CloneAlarm(datum)
				normalizer := dataNormalizer.New()
				Expect(normalizer).ToNot(BeNil())
				datum.Normalize(normalizer)
				Expect(normalizer.Error()).To(BeNil())
				Expect(normalizer.Data()).To(Equal([]data.Datum{datumStatus}))
				expectedDatum.Status = nil
				expectedDatum.StatusID = pointer.FromString(*datumStatus.ID)
				Expect(datum).To(Equal(expectedDatum))
			})
		})
	})
})
