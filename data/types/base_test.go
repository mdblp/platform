package types_test

import (
	"sort"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	associationTest "github.com/tidepool-org/platform/association/test"
	"github.com/tidepool-org/platform/data"
	dataNormalizer "github.com/tidepool-org/platform/data/normalizer"
	dataTest "github.com/tidepool-org/platform/data/test"
	"github.com/tidepool-org/platform/data/types"
	dataTypesTest "github.com/tidepool-org/platform/data/types/test"
	errorsTest "github.com/tidepool-org/platform/errors/test"
	"github.com/tidepool-org/platform/metadata"
	metadataTest "github.com/tidepool-org/platform/metadata/test"
	"github.com/tidepool-org/platform/net"
	originTest "github.com/tidepool-org/platform/origin/test"
	"github.com/tidepool-org/platform/pointer"
	"github.com/tidepool-org/platform/structure"
	structureParser "github.com/tidepool-org/platform/structure/parser"
	structureValidator "github.com/tidepool-org/platform/structure/validator"
	"github.com/tidepool-org/platform/test"
	timeZone "github.com/tidepool-org/platform/time/zone"
	timeZoneTest "github.com/tidepool-org/platform/time/zone/test"
	"github.com/tidepool-org/platform/user"
	userTest "github.com/tidepool-org/platform/user/test"
)

var _ = Describe("Base", func() {
	Context("New", func() {
		It("creates a new datum with all values initialized", func() {
			typ := dataTypesTest.NewType()
			datum := types.New(typ)
			Expect(datum.Active).To(BeFalse())
			Expect(datum.Annotations).To(BeNil())
			Expect(datum.Associations).To(BeNil())
			Expect(datum.ArchivedDataSetID).To(BeNil())
			Expect(datum.ArchivedTime).To(BeNil())
			Expect(datum.ClockDriftOffset).To(BeNil())
			Expect(datum.ConversionOffset).To(BeNil())
			Expect(datum.CreatedTime).To(BeNil())
			Expect(datum.CreatedUserID).To(BeNil())
			Expect(datum.Deduplicator).To(BeNil())
			Expect(datum.DeletedTime).To(BeNil())
			Expect(datum.DeletedUserID).To(BeNil())
			Expect(datum.DeviceID).To(BeNil())
			Expect(datum.DeviceTime).To(BeNil())
			Expect(datum.GUID).To(BeNil())
			Expect(datum.ID).To(BeNil())
			Expect(datum.ModifiedTime).To(BeNil())
			Expect(datum.ModifiedUserID).To(BeNil())
			Expect(datum.Notes).To(BeNil())
			Expect(datum.Origin).To(BeNil())
			Expect(datum.Payload).To(BeNil())
			Expect(datum.Source).To(BeNil())
			Expect(datum.Tags).To(BeNil())
			Expect(datum.Time).To(BeNil())
			Expect(datum.TimeZoneName).To(BeNil())
			Expect(datum.TimeZoneOffset).To(BeNil())
			Expect(datum.Type).To(Equal(typ))
			Expect(datum.UploadID).To(BeNil())
			Expect(datum.UserID).To(BeNil())
			Expect(datum.VersionInternal).To(Equal(0))
		})
	})

	Context("with new datum", func() {
		var typ string
		var datum types.Base

		BeforeEach(func() {
			typ = dataTypesTest.NewType()
			datum = types.New(typ)
		})

		Context("Meta", func() {
			It("returns the meta with type", func() {
				Expect(datum.Meta()).To(Equal(&types.Meta{Type: typ}))
			})
		})
	})

	Context("Base", func() {
		Context("Parse", func() {
			var parsedBase types.Base
			It("parses eventId when guid is missing", func() {
				parsedBase = types.New("")
				object := map[string]interface{}{"eventId": "1234"}
				parser := structureParser.NewObject(&object)
				parsedBase.Parse(parser)
				Expect(*parsedBase.GUID).To(Equal("1234"))
				Expect(parser.Error()).ToNot(HaveOccurred())
			})
			It("parses eventId when guid is empty", func() {
				parsedBase = types.New("")
				object := map[string]interface{}{"eventId": "1234", "guid": ""}
				parser := structureParser.NewObject(&object)
				parsedBase.Parse(parser)
				Expect(*parsedBase.GUID).To(Equal("1234"))
				Expect(parser.Error()).ToNot(HaveOccurred())
			})
			It("doesn't parses eventId when guid is not empty", func() {
				parsedBase = types.New("")
				object := map[string]interface{}{"eventId": "1234", "guid": "4567"}
				parser := structureParser.NewObject(&object)
				parsedBase.Parse(parser)
				Expect(*parsedBase.GUID).To(Equal("4567"))
				Expect(parser.Error()).ToNot(HaveOccurred())
			})
			It("parses guid", func() {
				parsedBase = types.New("")
				object := map[string]interface{}{"guid": "4567"}
				parser := structureParser.NewObject(&object)
				parsedBase.Parse(parser)
				Expect(*parsedBase.GUID).To(Equal("4567"))
				Expect(parser.Error()).ToNot(HaveOccurred())
			})
			It("parses without error with empty guid and eventId", func() {
				parsedBase = types.New("")
				object := map[string]interface{}{}
				parser := structureParser.NewObject(&object)
				parsedBase.Parse(parser)
				Expect(parsedBase.GUID).To(BeNil())
				Expect(parser.Error()).ToNot(HaveOccurred())
			})
		})

		Context("Validate", func() {
			DescribeTable("validates the datum",
				func(mutator func(datum *types.Base), expectedOrigins []structure.Origin, expectedErrors ...error) {
					datum := dataTypesTest.NewBase()
					mutator(datum)
					dataTypesTest.ValidateWithExpectedOrigins(datum, expectedOrigins, expectedErrors...)
				},
				Entry("succeeds",
					func(datum *types.Base) {},
					structure.Origins(),
				),
				Entry("active true",
					func(datum *types.Base) { datum.Active = true },
					structure.Origins(),
				),
				Entry("active false",
					func(datum *types.Base) { datum.Active = false },
					structure.Origins(),
				),
				Entry("annotations missing",
					func(datum *types.Base) { datum.Annotations = nil },
					structure.Origins(),
				),
				Entry("annotations exist",
					func(datum *types.Base) { datum.Annotations = metadataTest.RandomMetadataArray() },
					structure.Origins(),
				),
				Entry("associations missing",
					func(datum *types.Base) { datum.Associations = nil },
					structure.Origins(),
				),
				Entry("associations invalid",
					func(datum *types.Base) { (*datum.Associations)[0].Type = nil },
					structure.Origins(),
					errorsTest.WithPointerSource(structureValidator.ErrorValueNotExists(), "/associations/0/type"),
				),
				Entry("associations valid",
					func(datum *types.Base) { datum.Associations = associationTest.RandomAssociationArray() },
					structure.Origins(),
				),
				Entry("archived data set id missing",
					func(datum *types.Base) { datum.ArchivedDataSetID = nil },
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueNotExists(), "/archivedDatasetId"),
				),
				Entry("archived data set id empty",
					func(datum *types.Base) { datum.ArchivedDataSetID = pointer.FromString("") },
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueEmpty(), "/archivedDatasetId"),
				),
				Entry("archived data set id invalid",
					func(datum *types.Base) { datum.ArchivedDataSetID = pointer.FromString("invalid") },
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(data.ErrorValueStringAsSetIDNotValid("invalid"), "/archivedDatasetId"),
				),
				Entry("archived data set id valid",
					func(datum *types.Base) { datum.ArchivedDataSetID = pointer.FromString(dataTest.RandomSetID()) },
					structure.Origins(),
				),
				Entry("archived time missing; archived data set id missing",
					func(datum *types.Base) {
						datum.ArchivedDataSetID = nil
						datum.ArchivedTime = nil
					},
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
				),
				Entry("archived time missing; archived data set id exists",
					func(datum *types.Base) {
						datum.ArchivedDataSetID = pointer.FromString(dataTest.RandomSetID())
						datum.ArchivedTime = nil
					},
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueExists(), "/archivedDatasetId"),
				),
				Entry("archived time invalid",
					func(datum *types.Base) { datum.ArchivedTime = pointer.FromString("invalid") },
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueStringAsTimeNotValid("invalid", time.RFC3339Nano), "/archivedTime"),
				),
				Entry("archived time not after created time",
					func(datum *types.Base) {
						datum.ArchivedTime = pointer.FromString(test.PastFarTime().Format(time.RFC3339Nano))
						datum.CreatedTime = pointer.FromString(test.FutureFarTime().Format(time.RFC3339Nano))
						datum.DeletedTime = pointer.FromString(test.FutureFarTime().Format(time.RFC3339Nano))
						datum.ModifiedTime = pointer.FromString(test.FutureFarTime().Format(time.RFC3339Nano))
					},
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotAfter(test.PastFarTime(), test.FutureFarTime()), "/archivedTime"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotBeforeNow(test.FutureFarTime()), "/createdTime"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotBeforeNow(test.FutureFarTime()), "/deletedTime"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotBeforeNow(test.FutureFarTime()), "/modifiedTime"),
				),
				Entry("archived time not before now",
					func(datum *types.Base) {
						datum.ArchivedTime = pointer.FromString(test.FutureFarTime().Format(time.RFC3339Nano))
						datum.DeletedTime = pointer.FromString(test.FutureFarTime().Format(time.RFC3339Nano))
						datum.ModifiedTime = pointer.FromString(test.FutureFarTime().Format(time.RFC3339Nano))
					},
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotBeforeNow(test.FutureFarTime()), "/archivedTime"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotBeforeNow(test.FutureFarTime()), "/deletedTime"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotBeforeNow(test.FutureFarTime()), "/modifiedTime"),
				),
				Entry("clock drift offset missing",
					func(datum *types.Base) { datum.ClockDriftOffset = nil },
					structure.Origins(),
				),
				Entry("clock drift offset; out of range (lower)",
					func(datum *types.Base) { datum.ClockDriftOffset = pointer.FromInt(-86400001) },
					structure.Origins(),
					errorsTest.WithPointerSource(structureValidator.ErrorValueNotInRange(-86400001, -86400000, 86400000), "/clockDriftOffset"),
				),
				Entry("clock drift offset; in range (lower)",
					func(datum *types.Base) { datum.ClockDriftOffset = pointer.FromInt(-86400000) },
					structure.Origins(),
				),
				Entry("clock drift offset; in range (upper)",
					func(datum *types.Base) { datum.ClockDriftOffset = pointer.FromInt(86400000) },
					structure.Origins(),
				),
				Entry("clock drift offset; out of range (upper)",
					func(datum *types.Base) { datum.ClockDriftOffset = pointer.FromInt(86400001) },
					structure.Origins(),
					errorsTest.WithPointerSource(structureValidator.ErrorValueNotInRange(86400001, -86400000, 86400000), "/clockDriftOffset"),
				),
				Entry("conversion offset missing",
					func(datum *types.Base) { datum.ConversionOffset = nil },
					structure.Origins(),
				),
				Entry("conversion offset exists",
					func(datum *types.Base) { datum.ConversionOffset = pointer.FromInt(dataTypesTest.NewConversionOffset()) },
					structure.Origins(),
				),
				Entry("created user id missing",
					func(datum *types.Base) { datum.CreatedUserID = nil },
					structure.Origins(),
				),
				Entry("created user id empty",
					func(datum *types.Base) { datum.CreatedUserID = pointer.FromString("") },
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueEmpty(), "/createdUserId"),
				),
				Entry("created user id invalid",
					func(datum *types.Base) { datum.CreatedUserID = pointer.FromString("invalid") },
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(user.ErrorValueStringAsIDNotValid("invalid"), "/createdUserId"),
				),
				Entry("created user id valid",
					func(datum *types.Base) { datum.CreatedUserID = pointer.FromString(userTest.RandomID()) },
					structure.Origins(),
				),
				Entry("created time missing; created user id missing",
					func(datum *types.Base) {
						datum.CreatedTime = nil
						datum.CreatedUserID = nil
					},
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueNotExists(), "/createdTime"),
				),
				Entry("created time missing; created user id exists",
					func(datum *types.Base) {
						datum.CreatedTime = nil
						datum.CreatedUserID = pointer.FromString(userTest.RandomID())
					},
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueNotExists(), "/createdTime"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueExists(), "/createdUserId"),
				),
				Entry("created time invalid",
					func(datum *types.Base) { datum.CreatedTime = pointer.FromString("invalid") },
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueStringAsTimeNotValid("invalid", time.RFC3339Nano), "/createdTime"),
				),
				Entry("created time not before now",
					func(datum *types.Base) {
						datum.ArchivedTime = pointer.FromString(test.FutureFarTime().Format(time.RFC3339Nano))
						datum.CreatedTime = pointer.FromString(test.FutureFarTime().Format(time.RFC3339Nano))
						datum.DeletedTime = pointer.FromString(test.FutureFarTime().Format(time.RFC3339Nano))
						datum.ModifiedTime = pointer.FromString(test.FutureFarTime().Format(time.RFC3339Nano))
					},
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotBeforeNow(test.FutureFarTime()), "/archivedTime"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotBeforeNow(test.FutureFarTime()), "/createdTime"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotBeforeNow(test.FutureFarTime()), "/deletedTime"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotBeforeNow(test.FutureFarTime()), "/modifiedTime"),
				),
				Entry("deleted user id missing",
					func(datum *types.Base) { datum.DeletedUserID = nil },
					structure.Origins(),
				),
				Entry("deleted user id empty",
					func(datum *types.Base) { datum.DeletedUserID = pointer.FromString("") },
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueEmpty(), "/deletedUserId"),
				),
				Entry("deleted user id invalid",
					func(datum *types.Base) { datum.DeletedUserID = pointer.FromString("invalid") },
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(user.ErrorValueStringAsIDNotValid("invalid"), "/deletedUserId"),
				),
				Entry("deleted user id valid",
					func(datum *types.Base) { datum.DeletedUserID = pointer.FromString(userTest.RandomID()) },
					structure.Origins(),
				),
				Entry("deleted time missing; deleted user id missing",
					func(datum *types.Base) {
						datum.DeletedTime = nil
						datum.DeletedUserID = nil
					},
					structure.Origins(),
				),
				Entry("deleted time missing; deleted user id exists",
					func(datum *types.Base) {
						datum.DeletedTime = nil
						datum.DeletedUserID = pointer.FromString(userTest.RandomID())
					},
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueExists(), "/deletedUserId"),
				),
				Entry("deleted time invalid",
					func(datum *types.Base) { datum.DeletedTime = pointer.FromString("invalid") },
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueStringAsTimeNotValid("invalid", time.RFC3339Nano), "/deletedTime"),
				),
				Entry("deleted time not after archived time",
					func(datum *types.Base) {
						datum.ArchivedTime = pointer.FromString(test.FutureFarTime().Format(time.RFC3339Nano))
						datum.DeletedTime = pointer.FromString(test.PastFarTime().Format(time.RFC3339Nano))
						datum.ModifiedTime = pointer.FromString(test.FutureFarTime().Format(time.RFC3339Nano))
					},
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotBeforeNow(test.FutureFarTime()), "/archivedTime"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotAfter(test.PastFarTime(), test.FutureFarTime()), "/deletedTime"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotBeforeNow(test.FutureFarTime()), "/modifiedTime"),
				),
				Entry("deleted time not after created time",
					func(datum *types.Base) {
						datum.ArchivedTime = pointer.FromString(test.FutureFarTime().Format(time.RFC3339Nano))
						datum.CreatedTime = pointer.FromString(test.FutureFarTime().Format(time.RFC3339Nano))
						datum.DeletedTime = pointer.FromString(test.PastFarTime().Format(time.RFC3339Nano))
						datum.ModifiedTime = pointer.FromString(test.FutureFarTime().Format(time.RFC3339Nano))
					},
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotBeforeNow(test.FutureFarTime()), "/archivedTime"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotBeforeNow(test.FutureFarTime()), "/createdTime"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotAfter(test.PastFarTime(), test.FutureFarTime()), "/deletedTime"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotBeforeNow(test.FutureFarTime()), "/modifiedTime"),
				),
				Entry("deleted time not after modified time",
					func(datum *types.Base) {
						datum.DeletedTime = pointer.FromString(test.PastFarTime().Format(time.RFC3339Nano))
						datum.ModifiedTime = pointer.FromString(test.FutureFarTime().Format(time.RFC3339Nano))
					},
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotAfter(test.PastFarTime(), test.FutureFarTime()), "/deletedTime"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotBeforeNow(test.FutureFarTime()), "/modifiedTime"),
				),
				Entry("deleted time not before now",
					func(datum *types.Base) {
						datum.DeletedTime = pointer.FromString(test.FutureFarTime().Format(time.RFC3339Nano))
					},
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotBeforeNow(test.FutureFarTime()), "/deletedTime"),
				),
				Entry("deduplicator missing",
					func(datum *types.Base) { datum.Deduplicator = nil },
					structure.Origins(),
				),
				Entry("deduplicator invalid",
					func(datum *types.Base) { datum.Deduplicator.Name = pointer.FromString("invalid") },
					structure.Origins(),
					errorsTest.WithPointerSource(net.ErrorValueStringAsReverseDomainNotValid("invalid"), "/deduplicator/name"),
				),
				Entry("deduplicator valid",
					func(datum *types.Base) { datum.Deduplicator = dataTest.RandomDeduplicatorDescriptor() },
					structure.Origins(),
				),
				Entry("device id missing",
					func(datum *types.Base) { datum.DeviceID = nil },
					structure.Origins(),
				),
				Entry("device id empty",
					func(datum *types.Base) { datum.DeviceID = pointer.FromString("") },
					structure.Origins(),
					errorsTest.WithPointerSource(structureValidator.ErrorValueEmpty(), "/deviceId"),
				),
				Entry("device id valid",
					func(datum *types.Base) { datum.DeviceID = pointer.FromString(dataTest.NewDeviceID()) },
					structure.Origins(),
				),
				Entry("device time missing",
					func(datum *types.Base) { datum.DeviceTime = nil },
					structure.Origins(),
				),
				Entry("device time invalid",
					func(datum *types.Base) { datum.DeviceTime = pointer.FromString("invalid") },
					structure.Origins(),
					errorsTest.WithPointerSource(structureValidator.ErrorValueStringAsTimeNotValid("invalid", "2006-01-02T15:04:05"), "/deviceTime"),
				),
				Entry("device time valid",
					func(datum *types.Base) {
						datum.DeviceTime = pointer.FromString(test.RandomTime().Format("2006-01-02T15:04:05"))
					},
					structure.Origins(),
				),
				Entry("id missing",
					func(datum *types.Base) { datum.ID = nil },
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueNotExists(), "/id"),
				),
				Entry("id empty",
					func(datum *types.Base) { datum.ID = pointer.FromString("") },
					structure.Origins(),
					errorsTest.WithPointerSource(structureValidator.ErrorValueEmpty(), "/id"),
				),
				Entry("id invalid",
					func(datum *types.Base) { datum.ID = pointer.FromString("invalid") },
					structure.Origins(),
					errorsTest.WithPointerSource(data.ErrorValueStringAsIDNotValid("invalid"), "/id"),
				),
				Entry("id valid",
					func(datum *types.Base) { datum.ID = pointer.FromString(dataTest.RandomID()) },
					structure.Origins(),
				),
				Entry("modified user id missing",
					func(datum *types.Base) { datum.ModifiedUserID = nil },
					structure.Origins(),
				),
				Entry("modified user id empty",
					func(datum *types.Base) { datum.ModifiedUserID = pointer.FromString("") },
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueEmpty(), "/modifiedUserId"),
				),
				Entry("modified user id invalid",
					func(datum *types.Base) { datum.ModifiedUserID = pointer.FromString("invalid") },
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(user.ErrorValueStringAsIDNotValid("invalid"), "/modifiedUserId"),
				),
				Entry("modified user id valid",
					func(datum *types.Base) { datum.ModifiedUserID = pointer.FromString(userTest.RandomID()) },
					structure.Origins(),
				),
				Entry("modified time missing; modified user id missing",
					func(datum *types.Base) {
						datum.ArchivedTime = nil
						datum.ArchivedDataSetID = nil
						datum.ModifiedTime = nil
						datum.ModifiedUserID = nil
					},
					structure.Origins(),
				),
				Entry("modified time missing; modified user id exists",
					func(datum *types.Base) {
						datum.ArchivedTime = nil
						datum.ArchivedDataSetID = nil
						datum.ModifiedTime = nil
						datum.ModifiedUserID = pointer.FromString(userTest.RandomID())
					},
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueExists(), "/modifiedUserId"),
				),
				Entry("modified time missing; modified user id missing; archived time exists",
					func(datum *types.Base) {
						datum.ArchivedTime = datum.ModifiedTime
						datum.ModifiedTime = nil
						datum.ModifiedUserID = nil
					},
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueNotExists(), "/modifiedTime"),
				),
				Entry("modified time missing; modified user id exists; archived time exists",
					func(datum *types.Base) {
						datum.ArchivedTime = datum.ModifiedTime
						datum.ModifiedTime = nil
						datum.ModifiedUserID = pointer.FromString(userTest.RandomID())
					},
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueNotExists(), "/modifiedTime"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueExists(), "/modifiedUserId"),
				),
				Entry("modified time invalid",
					func(datum *types.Base) {
						datum.ModifiedTime = pointer.FromString("invalid")
					},
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueStringAsTimeNotValid("invalid", time.RFC3339Nano), "/modifiedTime"),
				),
				Entry("modified time not after archived time",
					func(datum *types.Base) {
						datum.ArchivedTime = pointer.FromString(test.FutureFarTime().Format(time.RFC3339Nano))
						datum.DeletedTime = pointer.FromString(test.FutureFarTime().Format(time.RFC3339Nano))
						datum.ModifiedTime = pointer.FromString(test.PastFarTime().Format(time.RFC3339Nano))
					},
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotBeforeNow(test.FutureFarTime()), "/archivedTime"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotBeforeNow(test.FutureFarTime()), "/deletedTime"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotAfter(test.PastFarTime(), test.FutureFarTime()), "/modifiedTime"),
				),
				Entry("modified time not after created time",
					func(datum *types.Base) {
						datum.ArchivedTime = pointer.FromString(test.FutureFarTime().Format(time.RFC3339Nano))
						datum.DeletedTime = pointer.FromString(test.FutureFarTime().Format(time.RFC3339Nano))
						datum.CreatedTime = pointer.FromString(test.FutureFarTime().Format(time.RFC3339Nano))
						datum.ModifiedTime = pointer.FromString(test.PastFarTime().Format(time.RFC3339Nano))
					},
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotBeforeNow(test.FutureFarTime()), "/archivedTime"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotBeforeNow(test.FutureFarTime()), "/createdTime"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotBeforeNow(test.FutureFarTime()), "/deletedTime"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotAfter(test.PastFarTime(), test.FutureFarTime()), "/modifiedTime"),
				),
				Entry("modified time not before now",
					func(datum *types.Base) {
						datum.DeletedTime = pointer.FromString(test.FutureFarTime().Format(time.RFC3339Nano))
						datum.ModifiedTime = pointer.FromString(test.FutureFarTime().Format(time.RFC3339Nano))
					},
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotBeforeNow(test.FutureFarTime()), "/deletedTime"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueTimeNotBeforeNow(test.FutureFarTime()), "/modifiedTime"),
				),
				Entry("notes missing",
					func(datum *types.Base) { datum.Notes = nil },
					structure.Origins(),
				),
				Entry("notes empty",
					func(datum *types.Base) { datum.Notes = pointer.FromStringArray([]string{}) },
					structure.Origins(),
					errorsTest.WithPointerSource(structureValidator.ErrorValueEmpty(), "/notes"),
				),
				Entry("notes length; in range (upper)",
					func(datum *types.Base) { datum.Notes = pointer.FromStringArray(dataTypesTest.NewNotes(100, 100)) },
					structure.Origins(),
				),
				Entry("notes length; out of range (upper)",
					func(datum *types.Base) { datum.Notes = pointer.FromStringArray(dataTypesTest.NewNotes(101, 101)) },
					structure.Origins(),
					errorsTest.WithPointerSource(structureValidator.ErrorLengthNotLessThanOrEqualTo(101, 100), "/notes"),
				),
				Entry("notes note empty",
					func(datum *types.Base) {
						datum.Notes = pointer.FromStringArray(append([]string{dataTypesTest.NewNote(1, 1000), "", dataTypesTest.NewNote(1, 1000), ""}, dataTypesTest.NewNotes(0, 96)...))
					},
					structure.Origins(),
					errorsTest.WithPointerSource(structureValidator.ErrorValueEmpty(), "/notes/1"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueEmpty(), "/notes/3"),
				),
				Entry("notes note length; in range (upper)",
					func(datum *types.Base) {
						datum.Notes = pointer.FromStringArray(append([]string{dataTypesTest.NewNote(1000, 1000), dataTypesTest.NewNote(1, 1000), dataTypesTest.NewNote(1000, 1000)}, dataTypesTest.NewNotes(0, 97)...))
					},
					structure.Origins(),
				),
				Entry("notes note length; out of range (upper)",
					func(datum *types.Base) {
						datum.Notes = pointer.FromStringArray(append([]string{dataTypesTest.NewNote(1001, 1001), dataTypesTest.NewNote(1, 1000), dataTypesTest.NewNote(1001, 1001)}, dataTypesTest.NewNotes(0, 97)...))
					},
					structure.Origins(),
					errorsTest.WithPointerSource(structureValidator.ErrorLengthNotLessThanOrEqualTo(1001, 1000), "/notes/0"),
					errorsTest.WithPointerSource(structureValidator.ErrorLengthNotLessThanOrEqualTo(1001, 1000), "/notes/2"),
				),
				Entry("origin missing",
					func(datum *types.Base) { datum.Origin = nil },
					structure.Origins(),
				),
				Entry("origin invalid",
					func(datum *types.Base) { datum.Origin.Name = pointer.FromString("") },
					structure.Origins(),
					errorsTest.WithPointerSource(structureValidator.ErrorValueEmpty(), "/origin/name"),
				),
				Entry("origin valid",
					func(datum *types.Base) { datum.Origin = originTest.RandomOrigin() },
					structure.Origins(),
				),
				Entry("payload missing",
					func(datum *types.Base) { datum.Payload = nil },
					structure.Origins(),
				),
				Entry("payload invalid",
					func(datum *types.Base) { datum.Payload = metadata.NewMetadata() },
					structure.Origins(),
					errorsTest.WithPointerSource(structureValidator.ErrorValueEmpty(), "/payload"),
				),
				Entry("payload valid",
					func(datum *types.Base) { datum.Payload = metadataTest.RandomMetadata() },
					structure.Origins(),
				),
				Entry("source missing",
					func(datum *types.Base) { datum.Source = nil },
					structure.Origins(),
				),
				Entry("source invalid",
					func(datum *types.Base) { datum.Source = pointer.FromString("invalid") },
					structure.Origins(),
					errorsTest.WithPointerSource(structureValidator.ErrorValueNotEqualTo("invalid", "carelink"), "/source"),
				),
				Entry("source valid",
					func(datum *types.Base) { datum.Source = pointer.FromString("carelink") },
					structure.Origins(),
				),
				Entry("tags missing",
					func(datum *types.Base) { datum.Tags = nil },
					structure.Origins(),
				),
				Entry("tags empty",
					func(datum *types.Base) { datum.Tags = pointer.FromStringArray([]string{}) },
					structure.Origins(),
					errorsTest.WithPointerSource(structureValidator.ErrorValueEmpty(), "/tags"),
				),
				Entry("tags length; in range (upper)",
					func(datum *types.Base) { datum.Tags = pointer.FromStringArray(dataTypesTest.NewTags(100, 100)) },
					structure.Origins(),
				),
				Entry("tags length; out of range (upper)",
					func(datum *types.Base) { datum.Tags = pointer.FromStringArray(dataTypesTest.NewTags(101, 101)) },
					structure.Origins(),
					errorsTest.WithPointerSource(structureValidator.ErrorLengthNotLessThanOrEqualTo(101, 100), "/tags"),
				),
				Entry("tags tag empty",
					func(datum *types.Base) {
						datum.Tags = pointer.FromStringArray(append([]string{dataTypesTest.NewTag(100, 100), ""}, dataTypesTest.NewTags(0, 98)...))
					},
					structure.Origins(),
					errorsTest.WithPointerSource(structureValidator.ErrorValueEmpty(), "/tags/1"),
				),
				Entry("tags tag length; in range (upper)",
					func(datum *types.Base) {
						datum.Tags = pointer.FromStringArray(append([]string{dataTypesTest.NewTag(100, 100)}, dataTypesTest.NewTags(0, 99)...))
					},
					structure.Origins(),
				),
				Entry("tags tag length; out of range (upper)",
					func(datum *types.Base) {
						datum.Tags = pointer.FromStringArray(append([]string{dataTypesTest.NewTag(101, 101)}, dataTypesTest.NewTags(0, 99)...))
					},
					structure.Origins(),
					errorsTest.WithPointerSource(structureValidator.ErrorLengthNotLessThanOrEqualTo(101, 100), "/tags/0"),
				),
				Entry("tags tag duplicate",
					func(datum *types.Base) {
						tags := dataTypesTest.NewTags(5, 99)
						datum.Tags = pointer.FromStringArray(append([]string{tags[4]}, tags...))
					},
					structure.Origins(),
					errorsTest.WithPointerSource(structureValidator.ErrorValueDuplicate(), "/tags/5"),
				),
				Entry("time missing",
					func(datum *types.Base) { datum.Time = nil },
					structure.Origins(),
					errorsTest.WithPointerSource(structureValidator.ErrorValueNotExists(), "/time"),
				),
				Entry("time invalid",
					func(datum *types.Base) { datum.Time = pointer.FromString("invalid") },
					structure.Origins(),
					errorsTest.WithPointerSource(structureValidator.ErrorValueStringAsTimeNotValid("invalid", time.RFC3339Nano), "/time"),
				),
				Entry("time valid",
					func(datum *types.Base) { datum.Time = pointer.FromString(test.RandomTime().Format(time.RFC3339Nano)) },
					structure.Origins(),
				),
				Entry("time zone name missing",
					func(datum *types.Base) { datum.TimeZoneName = nil },
					structure.Origins(),
				),
				Entry("time zone name empty",
					func(datum *types.Base) { datum.TimeZoneName = pointer.FromString("") },
					structure.Origins(),
					errorsTest.WithPointerSource(structureValidator.ErrorValueEmpty(), "/timezone"),
				),
				Entry("time zone name invalid",
					func(datum *types.Base) { datum.TimeZoneName = pointer.FromString("invalid") },
					structure.Origins(),
					errorsTest.WithPointerSource(timeZone.ErrorValueStringAsNameNotValid("invalid"), "/timezone"),
				),
				Entry("time zone name valid",
					func(datum *types.Base) { datum.TimeZoneName = pointer.FromString(timeZoneTest.RandomName()) },
					structure.Origins(),
				),
				Entry("time zone offset missing",
					func(datum *types.Base) { datum.TimeZoneOffset = nil },
					structure.Origins(),
				),
				Entry("time zone offset; out of range (lower)",
					func(datum *types.Base) { datum.TimeZoneOffset = pointer.FromInt(-10081) },
					structure.Origins(),
					errorsTest.WithPointerSource(structureValidator.ErrorValueNotInRange(-10081, -10080, 10080), "/timezoneOffset"),
				),
				Entry("time zone offset; in range (lower)",
					func(datum *types.Base) { datum.TimeZoneOffset = pointer.FromInt(-10080) },
					structure.Origins(),
				),
				Entry("time zone offset; in range (upper)",
					func(datum *types.Base) { datum.TimeZoneOffset = pointer.FromInt(10080) },
					structure.Origins(),
				),
				Entry("time zone offset; out of range (upper)",
					func(datum *types.Base) { datum.TimeZoneOffset = pointer.FromInt(10081) },
					structure.Origins(),
					errorsTest.WithPointerSource(structureValidator.ErrorValueNotInRange(10081, -10080, 10080), "/timezoneOffset"),
				),
				Entry("type empty",
					func(datum *types.Base) { datum.Type = "" },
					structure.Origins(),
					errorsTest.WithPointerSource(structureValidator.ErrorValueEmpty(), "/type"),
				),
				Entry("type valid",
					func(datum *types.Base) {
						datum.Type = test.RandomStringFromRangeAndCharset(1, 16, test.CharsetAlphaNumeric)
					},
					structure.Origins(),
				),
				Entry("upload id missing",
					func(datum *types.Base) { datum.UploadID = nil },
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueNotExists(), "/uploadId"),
				),
				Entry("upload id empty",
					func(datum *types.Base) { datum.UploadID = pointer.FromString("") },
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueEmpty(), "/uploadId"),
				),
				Entry("upload id invalid",
					func(datum *types.Base) { datum.UploadID = pointer.FromString("invalid") },
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(data.ErrorValueStringAsSetIDNotValid("invalid"), "/uploadId"),
				),
				Entry("upload id valid",
					func(datum *types.Base) { datum.UploadID = pointer.FromString(dataTest.RandomSetID()) },
					structure.Origins(),
				),
				Entry("user id missing",
					func(datum *types.Base) { datum.UserID = nil },
					[]structure.Origin{structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueNotExists(), "/_userId"),
				),
				Entry("user id empty",
					func(datum *types.Base) { datum.UserID = pointer.FromString("") },
					[]structure.Origin{structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueEmpty(), "/_userId"),
				),
				Entry("user id invalid",
					func(datum *types.Base) { datum.UserID = pointer.FromString("invalid") },
					[]structure.Origin{structure.OriginStore},
					errorsTest.WithPointerSource(user.ErrorValueStringAsIDNotValid("invalid"), "/_userId"),
				),
				Entry("user id valid",
					func(datum *types.Base) { datum.UserID = pointer.FromString(userTest.RandomID()) },
					structure.Origins(),
				),
				Entry("version; out of range (lower)",
					func(datum *types.Base) { datum.VersionInternal = -1 },
					[]structure.Origin{structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueNotGreaterThanOrEqualTo(-1, 0), "/_version"),
				),
				Entry("version; in range (lower)",
					func(datum *types.Base) { datum.VersionInternal = 0 },
					structure.Origins(),
				),
				Entry("multiple errors with store origin",
					func(datum *types.Base) {
						datum.UserID = nil
						datum.VersionInternal = -1
					},
					[]structure.Origin{structure.OriginStore},
					errorsTest.WithPointerSource(structureValidator.ErrorValueNotExists(), "/_userId"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueNotGreaterThanOrEqualTo(-1, 0), "/_version"),
				),
				Entry("multiple errors with internal origin",
					func(datum *types.Base) {
						datum.ArchivedDataSetID = pointer.FromString("invalid")
						datum.ArchivedTime = pointer.FromString("invalid")
						datum.CreatedTime = pointer.FromString("invalid")
						datum.CreatedUserID = pointer.FromString("invalid")
						datum.DeletedTime = pointer.FromString("invalid")
						datum.DeletedUserID = pointer.FromString("invalid")
						datum.ID = nil
						datum.ModifiedTime = pointer.FromString("invalid")
						datum.ModifiedUserID = pointer.FromString("invalid")
						datum.UploadID = nil
					},
					[]structure.Origin{structure.OriginInternal, structure.OriginStore},
					errorsTest.WithPointerSource(data.ErrorValueStringAsSetIDNotValid("invalid"), "/archivedDatasetId"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueStringAsTimeNotValid("invalid", time.RFC3339Nano), "/archivedTime"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueStringAsTimeNotValid("invalid", time.RFC3339Nano), "/createdTime"),
					errorsTest.WithPointerSource(user.ErrorValueStringAsIDNotValid("invalid"), "/createdUserId"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueStringAsTimeNotValid("invalid", time.RFC3339Nano), "/deletedTime"),
					errorsTest.WithPointerSource(user.ErrorValueStringAsIDNotValid("invalid"), "/deletedUserId"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueNotExists(), "/id"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueStringAsTimeNotValid("invalid", time.RFC3339Nano), "/modifiedTime"),
					errorsTest.WithPointerSource(user.ErrorValueStringAsIDNotValid("invalid"), "/modifiedUserId"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueNotExists(), "/uploadId"),
				),
				Entry("multiple errors with external origin",
					func(datum *types.Base) {
						datum.ClockDriftOffset = pointer.FromInt(-86400001)
						datum.Deduplicator.Name = pointer.FromString("invalid")
						datum.DeviceID = pointer.FromString("")
						datum.DeviceTime = pointer.FromString("invalid")
						datum.ID = pointer.FromString("")
						datum.Notes = pointer.FromStringArray([]string{})
						datum.Origin.Name = pointer.FromString("")
						datum.Source = pointer.FromString("invalid")
						datum.Tags = pointer.FromStringArray([]string{})
						datum.Time = nil
						datum.TimeZoneName = pointer.FromString("")
						datum.TimeZoneOffset = pointer.FromInt(-10081)
						datum.Type = ""
					},
					structure.Origins(),
					errorsTest.WithPointerSource(structureValidator.ErrorValueNotInRange(-86400001, -86400000, 86400000), "/clockDriftOffset"),
					errorsTest.WithPointerSource(net.ErrorValueStringAsReverseDomainNotValid("invalid"), "/deduplicator/name"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueEmpty(), "/deviceId"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueStringAsTimeNotValid("invalid", "2006-01-02T15:04:05"), "/deviceTime"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueEmpty(), "/id"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueEmpty(), "/notes"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueEmpty(), "/origin/name"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueNotEqualTo("invalid", "carelink"), "/source"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueEmpty(), "/tags"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueNotExists(), "/time"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueEmpty(), "/timezone"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueNotInRange(-10081, -10080, 10080), "/timezoneOffset"),
					errorsTest.WithPointerSource(structureValidator.ErrorValueEmpty(), "/type"),
				),
			)
		})

		Context("Normalize", func() {
			DescribeTable("normalizes the datum",
				func(mutator func(datum *types.Base), expectator func(datum *types.Base, expectedDatum *types.Base)) {
					for _, origin := range structure.Origins() {
						datum := dataTypesTest.NewBase()
						mutator(datum)
						expectedDatum := dataTypesTest.CloneBase(datum)
						normalizer := dataNormalizer.New()
						Expect(normalizer).ToNot(BeNil())
						datum.Normalize(normalizer.WithOrigin(origin))
						Expect(normalizer.Error()).To(BeNil())
						Expect(normalizer.Data()).To(BeEmpty())
						if expectator != nil {
							expectator(datum, expectedDatum)
						}
						Expect(datum).To(Equal(expectedDatum))
					}
				},
				Entry("tags sorted",
					func(datum *types.Base) {},
					func(datum *types.Base, expectedDatum *types.Base) {
						sort.Strings(*expectedDatum.Tags)
					},
				),
			)

			DescribeTable("normalizes the datum with origin external",
				func(mutator func(datum *types.Base), expectator func(datum *types.Base, expectedDatum *types.Base)) {
					datum := dataTypesTest.NewBase()
					mutator(datum)
					expectedDatum := dataTypesTest.CloneBase(datum)
					normalizer := dataNormalizer.New()
					Expect(normalizer).ToNot(BeNil())
					datum.Normalize(normalizer.WithOrigin(structure.OriginExternal))
					Expect(normalizer.Error()).To(BeNil())
					Expect(normalizer.Data()).To(BeEmpty())
					if expectator != nil {
						expectator(datum, expectedDatum)
					}
					Expect(datum).To(Equal(expectedDatum))
				},
				Entry("id missing",
					func(datum *types.Base) { datum.ID = nil },
					func(datum *types.Base, expectedDatum *types.Base) {
						Expect(datum.ID).ToNot(BeNil())
						Expect(datum.ID).ToNot(Equal(expectedDatum.ID))
						expectedDatum.ID = datum.ID
						sort.Strings(*expectedDatum.Tags)
					},
				),
				Entry("all missing",
					func(datum *types.Base) {
						*datum = types.New("")
					},
					func(datum *types.Base, expectedDatum *types.Base) {
						Expect(datum.ID).ToNot(BeNil())
						Expect(datum.ID).ToNot(Equal(expectedDatum.ID))
						expectedDatum.ID = datum.ID
					},
				),
			)

			DescribeTable("normalizes the datum with origin internal/store",
				func(mutator func(datum *types.Base), expectator func(datum *types.Base, expectedDatum *types.Base)) {
					for _, origin := range []structure.Origin{structure.OriginInternal, structure.OriginStore} {
						datum := dataTypesTest.NewBase()
						mutator(datum)
						expectedDatum := dataTypesTest.CloneBase(datum)
						normalizer := dataNormalizer.New()
						Expect(normalizer).ToNot(BeNil())
						datum.Normalize(normalizer.WithOrigin(origin))
						Expect(normalizer.Error()).To(BeNil())
						Expect(normalizer.Data()).To(BeEmpty())
						if expectator != nil {
							expectator(datum, expectedDatum)
						}
						Expect(datum).To(Equal(expectedDatum))
					}
				},
				Entry("id missing",
					func(datum *types.Base) { datum.ID = nil },
					func(datum *types.Base, expectedDatum *types.Base) {
						sort.Strings(*expectedDatum.Tags)
					},
				),
				Entry("all missing",
					func(datum *types.Base) {
						*datum = types.New("")
					},
					nil,
				),
				Entry("Time field set with offset, wrong TimeZoneName, wrong TimeZoneOffset",
					func(datum *types.Base) {
						datum.Time = pointer.FromString("2020-01-01T15:00:00+0100")
						datum.TimeZoneName = pointer.FromString("America/Chicago")
						datum.TimeZoneOffset = pointer.FromInt(0)
						datum.DeviceTime = pointer.FromString("2020-01-01T15:00:00")
					},
					func(datum *types.Base, expectedDatum *types.Base) {
						expectedDatum.Time = pointer.FromString("2020-01-01T14:00:00Z")
						expectedDatum.TimeZoneName = pointer.FromString("Etc/GMT-1")
						expectedDatum.TimeZoneOffset = pointer.FromInt(60)
					},
				),
				Entry("Time field set with offset, wrong TimeZoneName, correct TimeZoneOffset",
					func(datum *types.Base) {
						datum.Time = pointer.FromString("2020-01-01T15:00:00-0300")
						datum.TimeZoneName = pointer.FromString("America/Chicago")
						datum.TimeZoneOffset = pointer.FromInt(-180)
						datum.DeviceTime = pointer.FromString("2020-01-01T15:00:00")
					},
					func(datum *types.Base, expectedDatum *types.Base) {
						expectedDatum.Time = pointer.FromString("2020-01-01T18:00:00Z")
						expectedDatum.TimeZoneName = pointer.FromString("Etc/GMT+3")
						expectedDatum.TimeZoneOffset = pointer.FromInt(-180)
					},
				),
				Entry("Time field set with offset, no TimeZoneName, wrong TimeZoneOffset",
					func(datum *types.Base) {
						datum.Time = pointer.FromString("2020-01-01T15:00:00-0100")
						datum.TimeZoneName = nil
						datum.TimeZoneOffset = pointer.FromInt(180)
						datum.DeviceTime = pointer.FromString("2020-01-01T15:00:00")
					},
					func(datum *types.Base, expectedDatum *types.Base) {
						expectedDatum.Time = pointer.FromString("2020-01-01T16:00:00Z")
						expectedDatum.TimeZoneName = pointer.FromString("Etc/GMT+1")
						expectedDatum.TimeZoneOffset = pointer.FromInt(-60)
					},
				),
				Entry("Time field set as iso date, no TimeZoneName, correct TimeZoneOffset",
					func(datum *types.Base) {
						datum.Time = pointer.FromString("2020-01-01T15:00:00Z")
						datum.TimeZoneName = nil
						datum.TimeZoneOffset = pointer.FromInt(180)
						datum.DeviceTime = pointer.FromString("2020-01-01T12:00:00")
					},
					func(datum *types.Base, expectedDatum *types.Base) {
						expectedDatum.Time = pointer.FromString("2020-01-01T15:00:00Z")
						expectedDatum.TimeZoneName = pointer.FromString("Etc/GMT-3")
						expectedDatum.TimeZoneOffset = pointer.FromInt(180)
						expectedDatum.DeviceTime = pointer.FromString("2020-01-01T18:00:00")
					},
				),
				Entry("Time field set as iso date, no TimeZoneName, no TimeZoneOffset",
					func(datum *types.Base) {
						datum.Time = pointer.FromString("2020-01-01T15:00:00Z")
						datum.TimeZoneName = nil
						datum.TimeZoneOffset = nil
						datum.DeviceTime = pointer.FromString("2020-01-01T15:00:00")
					},
					func(datum *types.Base, expectedDatum *types.Base) {
						expectedDatum.Time = pointer.FromString("2020-01-01T15:00:00Z")
						expectedDatum.TimeZoneName = pointer.FromString("Etc/GMT")
						expectedDatum.TimeZoneOffset = pointer.FromInt(0)
						expectedDatum.DeviceTime = pointer.FromString("2020-01-01T15:00:00")
					},
				),
				Entry("Time field set with offset , correct TimeZoneName and wrong TimeZoneOffset",
					func(datum *types.Base) {
						datum.Time = pointer.FromString("2020-01-01T15:00:00-0600")
						datum.TimeZoneName = pointer.FromString("America/Chicago")
						datum.TimeZoneOffset = pointer.FromInt(0)
						datum.DeviceTime = pointer.FromString("2020-01-01T15:00:00")
					},
					func(datum *types.Base, expectedDatum *types.Base) {
						expectedDatum.Time = pointer.FromString("2020-01-01T21:00:00Z")
						expectedDatum.TimeZoneName = pointer.FromString("America/Chicago")
						expectedDatum.TimeZoneOffset = pointer.FromInt(-360)
					},
				),
				Entry("Time field set with malformed UTC, correct TimeZoneName, correct TimeZoneOffset",
					func(datum *types.Base) {
						datum.Time = pointer.FromString("2020-01-01T15:00:00+0000")
						datum.TimeZoneName = pointer.FromString("America/Chicago")
						datum.TimeZoneOffset = pointer.FromInt(-360)
						datum.DeviceTime = pointer.FromString("2020-01-01T09:00:00")
					},
					func(datum *types.Base, expectedDatum *types.Base) {
						expectedDatum.Time = pointer.FromString("2020-01-01T15:00:00Z")
						expectedDatum.TimeZoneName = pointer.FromString("America/Chicago")
						expectedDatum.TimeZoneOffset = pointer.FromInt(-360)
					},
				),
				Entry("DeviceTime unset",
					func(datum *types.Base) {
						datum.Time = pointer.FromString("2020-01-01T15:00:00Z")
						datum.TimeZoneName = pointer.FromString("America/Chicago")
						datum.TimeZoneOffset = pointer.FromInt(-360)
						datum.DeviceTime = nil
					},
					func(datum *types.Base, expectedDatum *types.Base) {
						expectedDatum.Time = pointer.FromString("2020-01-01T15:00:00Z")
						expectedDatum.TimeZoneName = pointer.FromString("America/Chicago")
						expectedDatum.TimeZoneOffset = pointer.FromInt(-360)
						expectedDatum.DeviceTime = pointer.FromString("2020-01-01T09:00:00")
					},
				),
				Entry("DeviceTime set with incorrect value",
					func(datum *types.Base) {
						datum.Time = pointer.FromString("2020-01-01T15:00:00Z")
						datum.TimeZoneName = pointer.FromString("America/Chicago")
						datum.TimeZoneOffset = pointer.FromInt(-360)
						datum.DeviceTime = pointer.FromString("2014-10-15T22:00:00")
					},
					func(datum *types.Base, expectedDatum *types.Base) {
						expectedDatum.Time = pointer.FromString("2020-01-01T15:00:00Z")
						expectedDatum.TimeZoneName = pointer.FromString("America/Chicago")
						expectedDatum.TimeZoneOffset = pointer.FromInt(-360)
						expectedDatum.DeviceTime = pointer.FromString("2020-01-01T09:00:00")
					},
				),
			)
		})
	})

	Context("with new, initialized datum", func() {
		var datum *types.Base

		BeforeEach(func() {
			datum = dataTypesTest.NewBase()
		})

		Context("IdentityFields", func() {
			It("returns error if user id is missing", func() {
				datum.UserID = nil
				identityFields, err := datum.IdentityFields()
				Expect(err).To(MatchError("user id is missing"))
				Expect(identityFields).To(BeEmpty())
			})

			It("returns error if user id is empty", func() {
				datum.UserID = pointer.FromString("")
				identityFields, err := datum.IdentityFields()
				Expect(err).To(MatchError("user id is empty"))
				Expect(identityFields).To(BeEmpty())
			})

			It("returns error if device id is missing", func() {
				datum.DeviceID = nil
				identityFields, err := datum.IdentityFields()
				Expect(err).To(MatchError("device id is missing"))
				Expect(identityFields).To(BeEmpty())
			})

			It("returns error if device id is empty", func() {
				datum.DeviceID = pointer.FromString("")
				identityFields, err := datum.IdentityFields()
				Expect(err).To(MatchError("device id is empty"))
				Expect(identityFields).To(BeEmpty())
			})

			It("returns error if time is missing", func() {
				datum.Time = nil
				identityFields, err := datum.IdentityFields()
				Expect(err).To(MatchError("time is missing"))
				Expect(identityFields).To(BeEmpty())
			})

			It("returns error if time is empty", func() {
				datum.Time = pointer.FromString("")
				identityFields, err := datum.IdentityFields()
				Expect(err).To(MatchError("time is empty"))
				Expect(identityFields).To(BeEmpty())
			})

			It("returns error if type is empty", func() {
				datum.Type = ""
				identityFields, err := datum.IdentityFields()
				Expect(err).To(MatchError("type is empty"))
				Expect(identityFields).To(BeEmpty())
			})

			It("returns the expected identity fields", func() {
				identityFields, err := datum.IdentityFields()
				Expect(err).ToNot(HaveOccurred())
				Expect(identityFields).To(Equal([]string{*datum.UserID, *datum.DeviceID, *datum.Time, datum.Type}))
			})
		})

		Context("GetPayload", func() {
			It("gets the payload", func() {
				Expect(datum.GetPayload()).To(Equal(datum.Payload))
			})
		})

		Context("SetUserID", func() {
			It("sets the user id", func() {
				userID := pointer.FromString(userTest.RandomID())
				datum.SetUserID(userID)
				Expect(datum.UserID).To(Equal(userID))
			})
		})

		Context("SetDataSetID", func() {
			It("sets the data set id", func() {
				dataSetID := pointer.FromString(dataTest.RandomSetID())
				datum.SetDataSetID(dataSetID)
				Expect(datum.UploadID).To(Equal(dataSetID))
			})
		})

		Context("SetActive", func() {
			It("sets active to true", func() {
				datum.SetActive(true)
				Expect(datum.Active).To(BeTrue())
			})

			It("sets active to false", func() {
				datum.SetActive(false)
				Expect(datum.Active).To(BeFalse())
			})
		})

		Context("SetDeviceID", func() {
			It("sets the device id", func() {
				deviceID := pointer.FromString(dataTest.NewDeviceID())
				datum.SetDeviceID(deviceID)
				Expect(datum.DeviceID).To(Equal(deviceID))
			})
		})

		Context("SetCreatedTime", func() {
			It("sets the created time", func() {
				createdTime := pointer.FromString(test.RandomTime().Format(time.RFC3339Nano))
				datum.SetCreatedTime(createdTime)
				Expect(datum.CreatedTime).To(Equal(createdTime))
			})
		})

		Context("SetCreatedUserID", func() {
			It("sets the created user id", func() {
				createdUserID := pointer.FromString(userTest.RandomID())
				datum.SetCreatedUserID(createdUserID)
				Expect(datum.CreatedUserID).To(Equal(createdUserID))
			})
		})

		Context("SetModifiedTime", func() {
			It("sets the modified time", func() {
				modifiedTime := pointer.FromString(test.RandomTime().Format(time.RFC3339Nano))
				datum.SetModifiedTime(modifiedTime)
				Expect(datum.ModifiedTime).To(Equal(modifiedTime))
			})
		})

		Context("SetModifiedUserID", func() {
			It("sets the modified user id", func() {
				modifiedUserID := pointer.FromString(userTest.RandomID())
				datum.SetModifiedUserID(modifiedUserID)
				Expect(datum.ModifiedUserID).To(Equal(modifiedUserID))
			})
		})

		Context("SetDeletedTime", func() {
			It("sets the deleted time", func() {
				deletedTime := pointer.FromString(test.RandomTime().Format(time.RFC3339Nano))
				datum.SetDeletedTime(deletedTime)
				Expect(datum.DeletedTime).To(Equal(deletedTime))
			})
		})

		Context("SetDeletedUserID", func() {
			It("sets the deleted user id", func() {
				deletedUserID := pointer.FromString(userTest.RandomID())
				datum.SetDeletedUserID(deletedUserID)
				Expect(datum.DeletedUserID).To(Equal(deletedUserID))
			})
		})

		Context("DeduplicatorDescriptor", func() {
			It("gets the deduplicator descriptor", func() {
				Expect(datum.DeduplicatorDescriptor()).To(Equal(datum.Deduplicator))
			})
		})

		Context("SetDeduplicatorDescriptor", func() {
			It("sets the deduplicator descriptor", func() {
				deduplicatorDescriptor := dataTest.RandomDeduplicatorDescriptor()
				datum.SetDeduplicatorDescriptor(deduplicatorDescriptor)
				Expect(datum.Deduplicator).To(Equal(deduplicatorDescriptor))
			})
		})
	})
})
