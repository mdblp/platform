package test

import (
	"bytes"
	"io/ioutil"
	"time"

	"github.com/onsi/gomega"
	gomegaGstruct "github.com/onsi/gomega/gstruct"
	gomegaTypes "github.com/onsi/gomega/types"

	"github.com/tidepool-org/platform/blob"
	"github.com/tidepool-org/platform/crypto"
	cryptoTest "github.com/tidepool-org/platform/crypto/test"
	netTest "github.com/tidepool-org/platform/net/test"
	"github.com/tidepool-org/platform/pointer"
	requestTest "github.com/tidepool-org/platform/request/test"
	"github.com/tidepool-org/platform/test"
	userTest "github.com/tidepool-org/platform/user/test"
)

func RandomID() string {
	return blob.NewID()
}

func RandomStatuses() []string {
	return test.RandomStringArrayFromRangeAndArrayWithoutDuplicates(1, len(blob.Statuses()), blob.Statuses())
}

func RandomFilter() *blob.Filter {
	datum := &blob.Filter{}
	datum.MediaType = pointer.FromStringArray(netTest.RandomMediaTypes(1, 3))
	datum.Status = pointer.FromStringArray(RandomStatuses())
	return datum
}

func NewObjectFromFilter(datum *blob.Filter, objectFormat test.ObjectFormat) map[string]interface{} {
	if datum == nil {
		return nil
	}
	object := map[string]interface{}{}
	if datum.MediaType != nil {
		object["mediaType"] = test.NewObjectFromStringArray(*datum.MediaType, objectFormat)
	}
	if datum.Status != nil {
		object["status"] = test.NewObjectFromStringArray(*datum.Status, objectFormat)
	}
	return object
}

func RandomContent() *blob.Content {
	content := test.RandomBytes()
	datum := &blob.Content{}
	datum.Body = ioutil.NopCloser(bytes.NewReader(content))
	datum.DigestMD5 = pointer.FromString(crypto.Base64EncodedMD5Hash(content))
	datum.MediaType = pointer.FromString(netTest.RandomMediaType())
	return datum
}

func RandomBlob() *blob.Blob {
	datum := &blob.Blob{}
	datum.ID = pointer.FromString(RandomID())
	datum.UserID = pointer.FromString(userTest.RandomID())
	datum.DigestMD5 = pointer.FromString(cryptoTest.RandomBase64EncodedMD5Hash())
	datum.MediaType = pointer.FromString(netTest.RandomMediaType())
	datum.Size = pointer.FromInt(test.RandomIntFromRange(1, 100*1024*1024))
	datum.Status = pointer.FromString(test.RandomStringFromArray(blob.Statuses()))
	datum.CreatedTime = pointer.FromTime(test.RandomTimeFromRange(test.RandomTimeMinimum(), time.Now()))
	if *datum.Status == blob.StatusAvailable {
		datum.ModifiedTime = pointer.FromTime(test.RandomTimeFromRange(*datum.CreatedTime, time.Now()))
	}
	datum.Revision = pointer.FromInt(requestTest.RandomRevision())
	return datum
}

func CloneBlob(datum *blob.Blob) *blob.Blob {
	if datum == nil {
		return nil
	}
	clone := &blob.Blob{}
	clone.ID = pointer.CloneString(datum.ID)
	clone.UserID = pointer.CloneString(datum.UserID)
	clone.DigestMD5 = pointer.CloneString(datum.DigestMD5)
	clone.MediaType = pointer.CloneString(datum.MediaType)
	clone.Size = pointer.CloneInt(datum.Size)
	clone.Status = pointer.CloneString(datum.Status)
	clone.CreatedTime = pointer.CloneTime(datum.CreatedTime)
	clone.ModifiedTime = pointer.CloneTime(datum.ModifiedTime)
	clone.DeletedTime = pointer.CloneTime(datum.DeletedTime)
	clone.Revision = pointer.CloneInt(datum.Revision)
	return clone
}

func NewObjectFromBlob(datum *blob.Blob, objectFormat test.ObjectFormat) map[string]interface{} {
	if datum == nil {
		return nil
	}
	object := map[string]interface{}{}
	if datum.ID != nil {
		object["id"] = test.NewObjectFromString(*datum.ID, objectFormat)
	}
	if datum.UserID != nil {
		object["userId"] = test.NewObjectFromString(*datum.UserID, objectFormat)
	}
	if datum.DigestMD5 != nil {
		object["digestMD5"] = test.NewObjectFromString(*datum.DigestMD5, objectFormat)
	}
	if datum.MediaType != nil {
		object["mediaType"] = test.NewObjectFromString(*datum.MediaType, objectFormat)
	}
	if datum.Size != nil {
		object["size"] = test.NewObjectFromInt(*datum.Size, objectFormat)
	}
	if datum.Status != nil {
		object["status"] = test.NewObjectFromString(*datum.Status, objectFormat)
	}
	if datum.CreatedTime != nil {
		object["createdTime"] = test.NewObjectFromTime(*datum.CreatedTime, objectFormat)
	}
	if datum.ModifiedTime != nil {
		object["modifiedTime"] = test.NewObjectFromTime(*datum.ModifiedTime, objectFormat)
	}
	if datum.DeletedTime != nil {
		object["deletedTime"] = test.NewObjectFromTime(*datum.DeletedTime, objectFormat)
	}
	if datum.Revision != nil {
		object["revision"] = test.NewObjectFromInt(*datum.Revision, objectFormat)
	}
	return object
}

func MatchBlob(datum *blob.Blob) gomegaTypes.GomegaMatcher {
	if datum == nil {
		return gomega.BeNil()
	}
	return gomegaGstruct.PointTo(gomegaGstruct.MatchAllFields(gomegaGstruct.Fields{
		"ID":           gomega.Equal(datum.ID),
		"UserID":       gomega.Equal(datum.UserID),
		"DigestMD5":    gomega.Equal(datum.DigestMD5),
		"MediaType":    gomega.Equal(datum.MediaType),
		"Size":         gomega.Equal(datum.Size),
		"Status":       gomega.Equal(datum.Status),
		"CreatedTime":  test.MatchTime(datum.CreatedTime),
		"ModifiedTime": test.MatchTime(datum.ModifiedTime),
		"DeletedTime":  test.MatchTime(datum.DeletedTime),
		"Revision":     gomega.Equal(datum.Revision),
	}))
}

func RandomBlobArray(minimumLength int, maximumLength int) blob.BlobArray {
	datum := make(blob.BlobArray, test.RandomIntFromRange(minimumLength, maximumLength))
	for index := range datum {
		datum[index] = RandomBlob()
	}
	return datum
}

func MatchBlobArray(datum blob.BlobArray) gomegaTypes.GomegaMatcher {
	matchers := []gomegaTypes.GomegaMatcher{}
	for _, d := range datum {
		matchers = append(matchers, MatchBlob(d))
	}
	return test.MatchArray(matchers)
}
