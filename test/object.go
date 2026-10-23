package test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"time"

	"github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ObjectFormat int

const (
	ObjectFormatBSON ObjectFormat = iota
	ObjectFormatJSON
)

func ExpectSerializedObjectBSON(object interface{}, expected interface{}) {
	gomega.Expect(object).ToNot(gomega.BeNil())
	gomega.Expect(expected).ToNot(gomega.BeNil())
	bites, err := bson.Marshal(object)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(bites).ToNot(gomega.BeNil())
	output := map[string]interface{}{}
	reg := bson.NewRegistry()
	reg.RegisterTypeMapEntry(bson.TypeDateTime, reflect.TypeOf(time.Time{}))
	reg.RegisterTypeMapEntry(bson.TypeArray, reflect.TypeOf([]interface{}{}))
	decoder := bson.NewDecoder(bson.NewDocumentReader(bytes.NewReader(bites)))
	decoder.SetRegistry(reg)
	gomega.Expect(decoder.Decode(&output)).To(gomega.Succeed())
	gomega.Expect(output).To(gomega.Equal(expected), "Unexpected serialized BSON")
}

func ExpectSerializedObjectJSON(object interface{}, expected interface{}) {
	gomega.Expect(object).ToNot(gomega.BeNil())
	gomega.Expect(expected).ToNot(gomega.BeNil())
	bites, err := json.Marshal(object)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
	gomega.Expect(bites).ToNot(gomega.BeNil())
	output := map[string]interface{}{}
	gomega.Expect(json.Unmarshal(bites, &output)).To(gomega.Succeed())
	gomega.Expect(output).To(gomega.Equal(expected), "Unexpected serialized JSON")
}
