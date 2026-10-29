package mongo

import (
	"testing"

	"github.com/tidepool-org/platform/data/types/basal/automated"
	"github.com/tidepool-org/platform/data/types/blood/glucose/continuous"
	"github.com/tidepool-org/platform/data/types/blood/glucose/selfmonitored"
	"github.com/tidepool-org/platform/data/types/device/prime"
	"github.com/tidepool-org/platform/data/types/settings/basalsecurity"
	"github.com/tidepool-org/platform/data/types/settings/pump"
	"github.com/tidepool-org/platform/data/types/upload"
)

func TestIsDatumToBucket(t *testing.T) {
	if !isDatumToBucket(continuous.New()) {
		t.Error("cbg must be written to a bucket")
	}
	if !isDatumToBucket(automated.New()) {
		t.Error("basal must be written to a bucket")
	}
	if !isDatumToBucket(prime.New()) {
		t.Error("deviceEvent must be written to a bucket")
	}
	if !isDatumToBucket(basalsecurity.New()) {
		t.Error("basalSecurity must be written to a bucket")
	}
	if isDatumToBucket(selfmonitored.New()) {
		t.Error("smbg must stay in deviceData")
	}
	if isDatumToBucket(pump.New()) {
		t.Error("pumpSettings must stay in deviceData")
	}
	if isDatumToBucket(upload.New()) {
		t.Error("upload must stay in deviceData")
	}
}
