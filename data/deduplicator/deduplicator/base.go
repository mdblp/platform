package deduplicator

import (
	"context"

	"github.com/tidepool-org/platform/data"
	dataStoreDEPRECATED "github.com/tidepool-org/platform/data/storeDEPRECATED"
	dataTypesUpload "github.com/tidepool-org/platform/data/types/upload"
	"github.com/tidepool-org/platform/errors"
	"github.com/tidepool-org/platform/net"
	"github.com/tidepool-org/platform/pointer"
)

type Base struct {
	name    string
	version string
}

func NewBase(name string, version string) (*Base, error) {
	if name == "" {
		return nil, errors.New("name is missing")
	} else if !net.IsValidReverseDomain(name) {
		return nil, errors.New("name is invalid")
	}
	if version == "" {
		return nil, errors.New("version is missing")
	} else if !net.IsValidSemanticVersion(version) {
		return nil, errors.New("version is invalid")
	}

	return &Base{
		name:    name,
		version: version,
	}, nil
}

func (b *Base) New(dataSet *dataTypesUpload.Upload) (bool, error) {
	return b.Get(dataSet)
}

func (b *Base) Get(dataSet *dataTypesUpload.Upload) (bool, error) {
	if dataSet == nil {
		return false, errors.New("data set is missing")
	}

	return dataSet.HasDeduplicatorNameMatch(b.name), nil
}

func (b *Base) Open(ctx context.Context, session dataStoreDEPRECATED.DataSession, dataSet *dataTypesUpload.Upload) (*dataTypesUpload.Upload, error) {
	if ctx == nil {
		return nil, errors.New("context is missing")
	}
	if session == nil {
		return nil, errors.New("session is missing")
	}
	if dataSet == nil {
		return nil, errors.New("data set is missing")
	}

	update := data.NewDataSetUpdate()
	update.Active = pointer.FromBool(dataSet.Active)
	update.Deduplicator = data.NewDeduplicatorDescriptor()
	update.Deduplicator.Name = pointer.FromString(b.name)
	update.Deduplicator.Version = pointer.FromString(b.version)
	return session.UpdateDataSet(ctx, *dataSet.UploadID, update)
}

func (b *Base) AddData(ctx context.Context, session dataStoreDEPRECATED.DataSession, dataSet *dataTypesUpload.Upload, dataSetData data.Data) error {
	if ctx == nil {
		return errors.New("context is missing")
	}
	if session == nil {
		return errors.New("session is missing")
	}
	if dataSet == nil {
		return errors.New("data set is missing")
	}
	if dataSetData == nil {
		return errors.New("data set data is missing")
	}

	return session.CreateDataSetData(ctx, dataSet, dataSetData)
}

func (b *Base) DeleteData(ctx context.Context, session dataStoreDEPRECATED.DataSession, dataSet *dataTypesUpload.Upload, selectors *data.Selectors) error {
	if ctx == nil {
		return errors.New("context is missing")
	}
	if session == nil {
		return errors.New("session is missing")
	}
	if dataSet == nil {
		return errors.New("data set is missing")
	}
	if selectors == nil {
		return errors.New("selectors is missing")
	}

	return session.DestroyDataSetData(ctx, dataSet, selectors)
}

func (b *Base) Close(ctx context.Context, session dataStoreDEPRECATED.DataSession, dataSet *dataTypesUpload.Upload) error {
	if ctx == nil {
		return errors.New("context is missing")
	}
	if session == nil {
		return errors.New("session is missing")
	}
	if dataSet == nil {
		return errors.New("data set is missing")
	}

	update := data.NewDataSetUpdate()
	update.Active = pointer.FromBool(true)
	if _, err := session.UpdateDataSet(ctx, *dataSet.UploadID, update); err != nil {
		return err
	}

	return session.ActivateDataSetData(ctx, dataSet, nil)
}

func (b *Base) Delete(ctx context.Context, session dataStoreDEPRECATED.DataSession, dataSet *dataTypesUpload.Upload, doPurge bool) error {
	if ctx == nil {
		return errors.New("context is missing")
	}
	if session == nil {
		return errors.New("session is missing")
	}
	if dataSet == nil {
		return errors.New("data set is missing")
	}

	return session.DeleteDataSet(ctx, dataSet, doPurge)
}
