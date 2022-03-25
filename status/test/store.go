package test

import "context"

func NewStoreStatusReporter() *StoreStatusReporter {
	return &StoreStatusReporter{}
}

type StoreStatus struct {
	Ping string
}

func OkStoreStatus() interface{} {
	return &StoreStatus{Ping: "OK"}
}

type StoreStatusReporter struct {
	sts interface{}
}

func (r *StoreStatusReporter) SetStatus(sts interface{}) {
	r.sts = sts
}

func (r *StoreStatusReporter) Status(ctx context.Context) interface{} {
	return r.sts
}
