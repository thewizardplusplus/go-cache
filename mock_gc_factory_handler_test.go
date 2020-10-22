// Code generated by mockery v1.0.0. DO NOT EDIT.

package cache

import gc "github.com/thewizardplusplus/go-cache/gc"
import hashmap "github.com/thewizardplusplus/go-hashmap"
import mock "github.com/stretchr/testify/mock"
import models "github.com/thewizardplusplus/go-cache/models"

// MockGCFactoryHandler is an autogenerated mock type for the GCFactoryHandler type
type MockGCFactoryHandler struct {
	mock.Mock
}

// NewGC provides a mock function with given fields: storage, clock
func (_m *MockGCFactoryHandler) NewGC(storage hashmap.Storage, clock models.Clock) gc.GC {
	ret := _m.Called(storage, clock)

	var r0 gc.GC
	if rf, ok := ret.Get(0).(func(hashmap.Storage, models.Clock) gc.GC); ok {
		r0 = rf(storage, clock)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gc.GC)
		}
	}

	return r0
}
