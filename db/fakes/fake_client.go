// This file was generated by counterfeiter
package fakes

import (
	"database/sql"
	"sync"

	"code.cloudfoundry.org/routing-api/db"
)

type FakeClient struct {
	CloseStub        func() error
	closeMutex       sync.RWMutex
	closeArgsForCall []struct{}
	closeReturns struct {
		result1 error
	}
	WhereStub        func(query interface{}, args ...interface{}) db.Client
	whereMutex       sync.RWMutex
	whereArgsForCall []struct {
		query interface{}
		args  []interface{}
	}
	whereReturns struct {
		result1 db.Client
	}
	CreateStub        func(value interface{}) (int64, error)
	createMutex       sync.RWMutex
	createArgsForCall []struct {
		value interface{}
	}
	createReturns struct {
		result1 int64
		result2 error
	}
	DeleteStub        func(value interface{}, where ...interface{}) (int64, error)
	deleteMutex       sync.RWMutex
	deleteArgsForCall []struct {
		value interface{}
		where []interface{}
	}
	deleteReturns struct {
		result1 int64
		result2 error
	}
	SaveStub        func(value interface{}) (int64, error)
	saveMutex       sync.RWMutex
	saveArgsForCall []struct {
		value interface{}
	}
	saveReturns struct {
		result1 int64
		result2 error
	}
	UpdateStub        func(attrs ...interface{}) (int64, error)
	updateMutex       sync.RWMutex
	updateArgsForCall []struct {
		attrs []interface{}
	}
	updateReturns struct {
		result1 int64
		result2 error
	}
	FirstStub        func(out interface{}, where ...interface{}) error
	firstMutex       sync.RWMutex
	firstArgsForCall []struct {
		out   interface{}
		where []interface{}
	}
	firstReturns struct {
		result1 error
	}
	FindStub        func(out interface{}, where ...interface{}) error
	findMutex       sync.RWMutex
	findArgsForCall []struct {
		out   interface{}
		where []interface{}
	}
	findReturns struct {
		result1 error
	}
	AutoMigrateStub        func(values ...interface{}) error
	autoMigrateMutex       sync.RWMutex
	autoMigrateArgsForCall []struct {
		values []interface{}
	}
	autoMigrateReturns struct {
		result1 error
	}
	BeginStub        func() db.Client
	beginMutex       sync.RWMutex
	beginArgsForCall []struct{}
	beginReturns struct {
		result1 db.Client
	}
	RollbackStub        func() error
	rollbackMutex       sync.RWMutex
	rollbackArgsForCall []struct{}
	rollbackReturns struct {
		result1 error
	}
	CommitStub        func() error
	commitMutex       sync.RWMutex
	commitArgsForCall []struct{}
	commitReturns struct {
		result1 error
	}
	HasTableStub        func(value interface{}) bool
	hasTableMutex       sync.RWMutex
	hasTableArgsForCall []struct {
		value interface{}
	}
	hasTableReturns struct {
		result1 bool
	}
	AddUniqueIndexStub        func(indexName string, columns ...string) (db.Client, error)
	addUniqueIndexMutex       sync.RWMutex
	addUniqueIndexArgsForCall []struct {
		indexName string
		columns   []string
	}
	addUniqueIndexReturns struct {
		result1 db.Client
		result2 error
	}
	ModelStub        func(value interface{}) db.Client
	modelMutex       sync.RWMutex
	modelArgsForCall []struct {
		value interface{}
	}
	modelReturns struct {
		result1 db.Client
	}
	ExecStub        func(query string, args ...interface{}) int64
	execMutex       sync.RWMutex
	execArgsForCall []struct {
		query string
		args  []interface{}
	}
	execReturns struct {
		result1 int64
	}
	RowsStub        func(tableName string) (*sql.Rows, error)
	rowsMutex       sync.RWMutex
	rowsArgsForCall []struct {
		tableName string
	}
	rowsReturns struct {
		result1 *sql.Rows
		result2 error
	}
	DropColumnStub        func(column string) error
	dropColumnMutex       sync.RWMutex
	dropColumnArgsForCall []struct {
		column string
	}
	dropColumnReturns struct {
		result1 error
	}
}

func (fake *FakeClient) Close() error {
	fake.closeMutex.Lock()
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct{}{})
	fake.closeMutex.Unlock()
	if fake.CloseStub != nil {
		return fake.CloseStub()
	} else {
		return fake.closeReturns.result1
	}
}

func (fake *FakeClient) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *FakeClient) CloseReturns(result1 error) {
	fake.CloseStub = nil
	fake.closeReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) Where(query interface{}, args ...interface{}) db.Client {
	fake.whereMutex.Lock()
	fake.whereArgsForCall = append(fake.whereArgsForCall, struct {
		query interface{}
		args  []interface{}
	}{query, args})
	fake.whereMutex.Unlock()
	if fake.WhereStub != nil {
		return fake.WhereStub(query, args...)
	} else {
		return fake.whereReturns.result1
	}
}

func (fake *FakeClient) WhereCallCount() int {
	fake.whereMutex.RLock()
	defer fake.whereMutex.RUnlock()
	return len(fake.whereArgsForCall)
}

func (fake *FakeClient) WhereArgsForCall(i int) (interface{}, []interface{}) {
	fake.whereMutex.RLock()
	defer fake.whereMutex.RUnlock()
	return fake.whereArgsForCall[i].query, fake.whereArgsForCall[i].args
}

func (fake *FakeClient) WhereReturns(result1 db.Client) {
	fake.WhereStub = nil
	fake.whereReturns = struct {
		result1 db.Client
	}{result1}
}

func (fake *FakeClient) Create(value interface{}) (int64, error) {
	fake.createMutex.Lock()
	fake.createArgsForCall = append(fake.createArgsForCall, struct {
		value interface{}
	}{value})
	fake.createMutex.Unlock()
	if fake.CreateStub != nil {
		return fake.CreateStub(value)
	} else {
		return fake.createReturns.result1, fake.createReturns.result2
	}
}

func (fake *FakeClient) CreateCallCount() int {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return len(fake.createArgsForCall)
}

func (fake *FakeClient) CreateArgsForCall(i int) interface{} {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return fake.createArgsForCall[i].value
}

func (fake *FakeClient) CreateReturns(result1 int64, result2 error) {
	fake.CreateStub = nil
	fake.createReturns = struct {
		result1 int64
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) Delete(value interface{}, where ...interface{}) (int64, error) {
	fake.deleteMutex.Lock()
	fake.deleteArgsForCall = append(fake.deleteArgsForCall, struct {
		value interface{}
		where []interface{}
	}{value, where})
	fake.deleteMutex.Unlock()
	if fake.DeleteStub != nil {
		return fake.DeleteStub(value, where...)
	} else {
		return fake.deleteReturns.result1, fake.deleteReturns.result2
	}
}

func (fake *FakeClient) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *FakeClient) DeleteArgsForCall(i int) (interface{}, []interface{}) {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return fake.deleteArgsForCall[i].value, fake.deleteArgsForCall[i].where
}

func (fake *FakeClient) DeleteReturns(result1 int64, result2 error) {
	fake.DeleteStub = nil
	fake.deleteReturns = struct {
		result1 int64
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) Save(value interface{}) (int64, error) {
	fake.saveMutex.Lock()
	fake.saveArgsForCall = append(fake.saveArgsForCall, struct {
		value interface{}
	}{value})
	fake.saveMutex.Unlock()
	if fake.SaveStub != nil {
		return fake.SaveStub(value)
	} else {
		return fake.saveReturns.result1, fake.saveReturns.result2
	}
}

func (fake *FakeClient) SaveCallCount() int {
	fake.saveMutex.RLock()
	defer fake.saveMutex.RUnlock()
	return len(fake.saveArgsForCall)
}

func (fake *FakeClient) SaveArgsForCall(i int) interface{} {
	fake.saveMutex.RLock()
	defer fake.saveMutex.RUnlock()
	return fake.saveArgsForCall[i].value
}

func (fake *FakeClient) SaveReturns(result1 int64, result2 error) {
	fake.SaveStub = nil
	fake.saveReturns = struct {
		result1 int64
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) Update(attrs ...interface{}) (int64, error) {
	fake.updateMutex.Lock()
	fake.updateArgsForCall = append(fake.updateArgsForCall, struct {
		attrs []interface{}
	}{attrs})
	fake.updateMutex.Unlock()
	if fake.UpdateStub != nil {
		return fake.UpdateStub(attrs...)
	} else {
		return fake.updateReturns.result1, fake.updateReturns.result2
	}
}

func (fake *FakeClient) UpdateCallCount() int {
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	return len(fake.updateArgsForCall)
}

func (fake *FakeClient) UpdateArgsForCall(i int) []interface{} {
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	return fake.updateArgsForCall[i].attrs
}

func (fake *FakeClient) UpdateReturns(result1 int64, result2 error) {
	fake.UpdateStub = nil
	fake.updateReturns = struct {
		result1 int64
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) First(out interface{}, where ...interface{}) error {
	fake.firstMutex.Lock()
	fake.firstArgsForCall = append(fake.firstArgsForCall, struct {
		out   interface{}
		where []interface{}
	}{out, where})
	fake.firstMutex.Unlock()
	if fake.FirstStub != nil {
		return fake.FirstStub(out, where...)
	} else {
		return fake.firstReturns.result1
	}
}

func (fake *FakeClient) FirstCallCount() int {
	fake.firstMutex.RLock()
	defer fake.firstMutex.RUnlock()
	return len(fake.firstArgsForCall)
}

func (fake *FakeClient) FirstArgsForCall(i int) (interface{}, []interface{}) {
	fake.firstMutex.RLock()
	defer fake.firstMutex.RUnlock()
	return fake.firstArgsForCall[i].out, fake.firstArgsForCall[i].where
}

func (fake *FakeClient) FirstReturns(result1 error) {
	fake.FirstStub = nil
	fake.firstReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) Find(out interface{}, where ...interface{}) error {
	fake.findMutex.Lock()
	fake.findArgsForCall = append(fake.findArgsForCall, struct {
		out   interface{}
		where []interface{}
	}{out, where})
	fake.findMutex.Unlock()
	if fake.FindStub != nil {
		return fake.FindStub(out, where...)
	} else {
		return fake.findReturns.result1
	}
}

func (fake *FakeClient) FindCallCount() int {
	fake.findMutex.RLock()
	defer fake.findMutex.RUnlock()
	return len(fake.findArgsForCall)
}

func (fake *FakeClient) FindArgsForCall(i int) (interface{}, []interface{}) {
	fake.findMutex.RLock()
	defer fake.findMutex.RUnlock()
	return fake.findArgsForCall[i].out, fake.findArgsForCall[i].where
}

func (fake *FakeClient) FindReturns(result1 error) {
	fake.FindStub = nil
	fake.findReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) AutoMigrate(values ...interface{}) error {
	fake.autoMigrateMutex.Lock()
	fake.autoMigrateArgsForCall = append(fake.autoMigrateArgsForCall, struct {
		values []interface{}
	}{values})
	fake.autoMigrateMutex.Unlock()
	if fake.AutoMigrateStub != nil {
		return fake.AutoMigrateStub(values...)
	} else {
		return fake.autoMigrateReturns.result1
	}
}

func (fake *FakeClient) AutoMigrateCallCount() int {
	fake.autoMigrateMutex.RLock()
	defer fake.autoMigrateMutex.RUnlock()
	return len(fake.autoMigrateArgsForCall)
}

func (fake *FakeClient) AutoMigrateArgsForCall(i int) []interface{} {
	fake.autoMigrateMutex.RLock()
	defer fake.autoMigrateMutex.RUnlock()
	return fake.autoMigrateArgsForCall[i].values
}

func (fake *FakeClient) AutoMigrateReturns(result1 error) {
	fake.AutoMigrateStub = nil
	fake.autoMigrateReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) Begin() db.Client {
	fake.beginMutex.Lock()
	fake.beginArgsForCall = append(fake.beginArgsForCall, struct{}{})
	fake.beginMutex.Unlock()
	if fake.BeginStub != nil {
		return fake.BeginStub()
	} else {
		return fake.beginReturns.result1
	}
}

func (fake *FakeClient) BeginCallCount() int {
	fake.beginMutex.RLock()
	defer fake.beginMutex.RUnlock()
	return len(fake.beginArgsForCall)
}

func (fake *FakeClient) BeginReturns(result1 db.Client) {
	fake.BeginStub = nil
	fake.beginReturns = struct {
		result1 db.Client
	}{result1}
}

func (fake *FakeClient) Rollback() error {
	fake.rollbackMutex.Lock()
	fake.rollbackArgsForCall = append(fake.rollbackArgsForCall, struct{}{})
	fake.rollbackMutex.Unlock()
	if fake.RollbackStub != nil {
		return fake.RollbackStub()
	} else {
		return fake.rollbackReturns.result1
	}
}

func (fake *FakeClient) RollbackCallCount() int {
	fake.rollbackMutex.RLock()
	defer fake.rollbackMutex.RUnlock()
	return len(fake.rollbackArgsForCall)
}

func (fake *FakeClient) RollbackReturns(result1 error) {
	fake.RollbackStub = nil
	fake.rollbackReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) Commit() error {
	fake.commitMutex.Lock()
	fake.commitArgsForCall = append(fake.commitArgsForCall, struct{}{})
	fake.commitMutex.Unlock()
	if fake.CommitStub != nil {
		return fake.CommitStub()
	} else {
		return fake.commitReturns.result1
	}
}

func (fake *FakeClient) CommitCallCount() int {
	fake.commitMutex.RLock()
	defer fake.commitMutex.RUnlock()
	return len(fake.commitArgsForCall)
}

func (fake *FakeClient) CommitReturns(result1 error) {
	fake.CommitStub = nil
	fake.commitReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) HasTable(value interface{}) bool {
	fake.hasTableMutex.Lock()
	fake.hasTableArgsForCall = append(fake.hasTableArgsForCall, struct {
		value interface{}
	}{value})
	fake.hasTableMutex.Unlock()
	if fake.HasTableStub != nil {
		return fake.HasTableStub(value)
	} else {
		return fake.hasTableReturns.result1
	}
}

func (fake *FakeClient) HasTableCallCount() int {
	fake.hasTableMutex.RLock()
	defer fake.hasTableMutex.RUnlock()
	return len(fake.hasTableArgsForCall)
}

func (fake *FakeClient) HasTableArgsForCall(i int) interface{} {
	fake.hasTableMutex.RLock()
	defer fake.hasTableMutex.RUnlock()
	return fake.hasTableArgsForCall[i].value
}

func (fake *FakeClient) HasTableReturns(result1 bool) {
	fake.HasTableStub = nil
	fake.hasTableReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakeClient) AddUniqueIndex(indexName string, columns ...string) (db.Client, error) {
	fake.addUniqueIndexMutex.Lock()
	fake.addUniqueIndexArgsForCall = append(fake.addUniqueIndexArgsForCall, struct {
		indexName string
		columns   []string
	}{indexName, columns})
	fake.addUniqueIndexMutex.Unlock()
	if fake.AddUniqueIndexStub != nil {
		return fake.AddUniqueIndexStub(indexName, columns...)
	} else {
		return fake.addUniqueIndexReturns.result1, fake.addUniqueIndexReturns.result2
	}
}

func (fake *FakeClient) AddUniqueIndexCallCount() int {
	fake.addUniqueIndexMutex.RLock()
	defer fake.addUniqueIndexMutex.RUnlock()
	return len(fake.addUniqueIndexArgsForCall)
}

func (fake *FakeClient) AddUniqueIndexArgsForCall(i int) (string, []string) {
	fake.addUniqueIndexMutex.RLock()
	defer fake.addUniqueIndexMutex.RUnlock()
	return fake.addUniqueIndexArgsForCall[i].indexName, fake.addUniqueIndexArgsForCall[i].columns
}

func (fake *FakeClient) AddUniqueIndexReturns(result1 db.Client, result2 error) {
	fake.AddUniqueIndexStub = nil
	fake.addUniqueIndexReturns = struct {
		result1 db.Client
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) Model(value interface{}) db.Client {
	fake.modelMutex.Lock()
	fake.modelArgsForCall = append(fake.modelArgsForCall, struct {
		value interface{}
	}{value})
	fake.modelMutex.Unlock()
	if fake.ModelStub != nil {
		return fake.ModelStub(value)
	} else {
		return fake.modelReturns.result1
	}
}

func (fake *FakeClient) ModelCallCount() int {
	fake.modelMutex.RLock()
	defer fake.modelMutex.RUnlock()
	return len(fake.modelArgsForCall)
}

func (fake *FakeClient) ModelArgsForCall(i int) interface{} {
	fake.modelMutex.RLock()
	defer fake.modelMutex.RUnlock()
	return fake.modelArgsForCall[i].value
}

func (fake *FakeClient) ModelReturns(result1 db.Client) {
	fake.ModelStub = nil
	fake.modelReturns = struct {
		result1 db.Client
	}{result1}
}

func (fake *FakeClient) Exec(query string, args ...interface{}) int64 {
	fake.execMutex.Lock()
	fake.execArgsForCall = append(fake.execArgsForCall, struct {
		query string
		args  []interface{}
	}{query, args})
	fake.execMutex.Unlock()
	if fake.ExecStub != nil {
		return fake.ExecStub(query, args...)
	} else {
		return fake.execReturns.result1
	}
}

func (fake *FakeClient) ExecCallCount() int {
	fake.execMutex.RLock()
	defer fake.execMutex.RUnlock()
	return len(fake.execArgsForCall)
}

func (fake *FakeClient) ExecArgsForCall(i int) (string, []interface{}) {
	fake.execMutex.RLock()
	defer fake.execMutex.RUnlock()
	return fake.execArgsForCall[i].query, fake.execArgsForCall[i].args
}

func (fake *FakeClient) ExecReturns(result1 int64) {
	fake.ExecStub = nil
	fake.execReturns = struct {
		result1 int64
	}{result1}
}

func (fake *FakeClient) Rows(tableName string) (*sql.Rows, error) {
	fake.rowsMutex.Lock()
	fake.rowsArgsForCall = append(fake.rowsArgsForCall, struct {
		tableName string
	}{tableName})
	fake.rowsMutex.Unlock()
	if fake.RowsStub != nil {
		return fake.RowsStub(tableName)
	} else {
		return fake.rowsReturns.result1, fake.rowsReturns.result2
	}
}

func (fake *FakeClient) RowsCallCount() int {
	fake.rowsMutex.RLock()
	defer fake.rowsMutex.RUnlock()
	return len(fake.rowsArgsForCall)
}

func (fake *FakeClient) RowsArgsForCall(i int) string {
	fake.rowsMutex.RLock()
	defer fake.rowsMutex.RUnlock()
	return fake.rowsArgsForCall[i].tableName
}

func (fake *FakeClient) RowsReturns(result1 *sql.Rows, result2 error) {
	fake.RowsStub = nil
	fake.rowsReturns = struct {
		result1 *sql.Rows
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) DropColumn(column string) error {
	fake.dropColumnMutex.Lock()
	fake.dropColumnArgsForCall = append(fake.dropColumnArgsForCall, struct {
		column string
	}{column})
	fake.dropColumnMutex.Unlock()
	if fake.DropColumnStub != nil {
		return fake.DropColumnStub(column)
	} else {
		return fake.dropColumnReturns.result1
	}
}

func (fake *FakeClient) DropColumnCallCount() int {
	fake.dropColumnMutex.RLock()
	defer fake.dropColumnMutex.RUnlock()
	return len(fake.dropColumnArgsForCall)
}

func (fake *FakeClient) DropColumnArgsForCall(i int) string {
	fake.dropColumnMutex.RLock()
	defer fake.dropColumnMutex.RUnlock()
	return fake.dropColumnArgsForCall[i].column
}

func (fake *FakeClient) DropColumnReturns(result1 error) {
	fake.DropColumnStub = nil
	fake.dropColumnReturns = struct {
		result1 error
	}{result1}
}

var _ db.Client = new(FakeClient)
