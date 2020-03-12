// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	blobstore "magma/orc8r/cloud/go/blobstore"

	mock "github.com/stretchr/testify/mock"

	storage "magma/orc8r/cloud/go/storage"
)

// TransactionalBlobStorage is an autogenerated mock type for the TransactionalBlobStorage type
type TransactionalBlobStorage struct {
	mock.Mock
}

// Commit provides a mock function with given fields:
func (_m *TransactionalBlobStorage) Commit() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateOrUpdate provides a mock function with given fields: networkID, blobs
func (_m *TransactionalBlobStorage) CreateOrUpdate(networkID string, blobs []blobstore.Blob) error {
	ret := _m.Called(networkID, blobs)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []blobstore.Blob) error); ok {
		r0 = rf(networkID, blobs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: networkID, ids
func (_m *TransactionalBlobStorage) Delete(networkID string, ids []storage.TypeAndKey) error {
	ret := _m.Called(networkID, ids)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []storage.TypeAndKey) error); ok {
		r0 = rf(networkID, ids)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: networkID, id
func (_m *TransactionalBlobStorage) Get(networkID string, id storage.TypeAndKey) (blobstore.Blob, error) {
	ret := _m.Called(networkID, id)

	var r0 blobstore.Blob
	if rf, ok := ret.Get(0).(func(string, storage.TypeAndKey) blobstore.Blob); ok {
		r0 = rf(networkID, id)
	} else {
		r0 = ret.Get(0).(blobstore.Blob)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, storage.TypeAndKey) error); ok {
		r1 = rf(networkID, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetExistingKeys provides a mock function with given fields: keys, filter
func (_m *TransactionalBlobStorage) GetExistingKeys(keys []string, filter blobstore.SearchFilter) ([]string, error) {
	ret := _m.Called(keys, filter)

	var r0 []string
	if rf, ok := ret.Get(0).(func([]string, blobstore.SearchFilter) []string); ok {
		r0 = rf(keys, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]string, blobstore.SearchFilter) error); ok {
		r1 = rf(keys, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMany provides a mock function with given fields: networkID, ids
func (_m *TransactionalBlobStorage) GetMany(networkID string, ids []storage.TypeAndKey) ([]blobstore.Blob, error) {
	ret := _m.Called(networkID, ids)

	var r0 []blobstore.Blob
	if rf, ok := ret.Get(0).(func(string, []storage.TypeAndKey) []blobstore.Blob); ok {
		r0 = rf(networkID, ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]blobstore.Blob)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, []storage.TypeAndKey) error); ok {
		r1 = rf(networkID, ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IncrementVersion provides a mock function with given fields: networkID, id
func (_m *TransactionalBlobStorage) IncrementVersion(networkID string, id storage.TypeAndKey) error {
	ret := _m.Called(networkID, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, storage.TypeAndKey) error); ok {
		r0 = rf(networkID, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListKeys provides a mock function with given fields: networkID, typeVal
func (_m *TransactionalBlobStorage) ListKeys(networkID string, typeVal string) ([]string, error) {
	ret := _m.Called(networkID, typeVal)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string, string) []string); ok {
		r0 = rf(networkID, typeVal)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(networkID, typeVal)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Rollback provides a mock function with given fields:
func (_m *TransactionalBlobStorage) Rollback() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Search provides a mock function with given fields: filter
func (_m *TransactionalBlobStorage) Search(filter blobstore.SearchFilter) (map[string][]blobstore.Blob, error) {
	ret := _m.Called(filter)

	var r0 map[string][]blobstore.Blob
	if rf, ok := ret.Get(0).(func(blobstore.SearchFilter) map[string][]blobstore.Blob); ok {
		r0 = rf(filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string][]blobstore.Blob)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(blobstore.SearchFilter) error); ok {
		r1 = rf(filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
