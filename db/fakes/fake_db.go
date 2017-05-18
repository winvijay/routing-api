// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"code.cloudfoundry.org/routing-api/db"
	"code.cloudfoundry.org/routing-api/models"
	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
)

type FakeDB struct {
	ReadRoutesStub        func() ([]models.Route, error)
	readRoutesMutex       sync.RWMutex
	readRoutesArgsForCall []struct{}
	readRoutesReturns struct {
		result1 []models.Route
		result2 error
	}
	SaveRouteStub        func(route models.Route) error
	saveRouteMutex       sync.RWMutex
	saveRouteArgsForCall []struct {
		route models.Route
	}
	saveRouteReturns struct {
		result1 error
	}
	DeleteRouteStub        func(route models.Route) error
	deleteRouteMutex       sync.RWMutex
	deleteRouteArgsForCall []struct {
		route models.Route
	}
	deleteRouteReturns struct {
		result1 error
	}
	ReadTcpRouteMappingsStub        func() ([]models.TcpRouteMapping, error)
	readTcpRouteMappingsMutex       sync.RWMutex
	readTcpRouteMappingsArgsForCall []struct{}
	readTcpRouteMappingsReturns struct {
		result1 []models.TcpRouteMapping
		result2 error
	}
	ReadFilteredTcpRouteMappingsStub        func(columnName string, values []string) ([]models.TcpRouteMapping, error)
	readFilteredTcpRouteMappingsMutex       sync.RWMutex
	readFilteredTcpRouteMappingsArgsForCall []struct {
		columnName string
		values     []string
	}
	readFilteredTcpRouteMappingsReturns struct {
		result1 []models.TcpRouteMapping
		result2 error
	}
	SaveTcpRouteMappingStub        func(tcpMapping models.TcpRouteMapping) error
	saveTcpRouteMappingMutex       sync.RWMutex
	saveTcpRouteMappingArgsForCall []struct {
		tcpMapping models.TcpRouteMapping
	}
	saveTcpRouteMappingReturns struct {
		result1 error
	}
	DeleteTcpRouteMappingStub        func(tcpMapping models.TcpRouteMapping) error
	deleteTcpRouteMappingMutex       sync.RWMutex
	deleteTcpRouteMappingArgsForCall []struct {
		tcpMapping models.TcpRouteMapping
	}
	deleteTcpRouteMappingReturns struct {
		result1 error
	}
	ReadRouterGroupsStub        func() (models.RouterGroups, error)
	readRouterGroupsMutex       sync.RWMutex
	readRouterGroupsArgsForCall []struct{}
	readRouterGroupsReturns struct {
		result1 models.RouterGroups
		result2 error
	}
	ReadRouterGroupStub        func(guid string) (models.RouterGroup, error)
	readRouterGroupMutex       sync.RWMutex
	readRouterGroupArgsForCall []struct {
		guid string
	}
	readRouterGroupReturns struct {
		result1 models.RouterGroup
		result2 error
	}
	ReadRouterGroupByNameStub        func(name string) (models.RouterGroup, error)
	readRouterGroupByNameMutex       sync.RWMutex
	readRouterGroupByNameArgsForCall []struct {
		name string
	}
	readRouterGroupByNameReturns struct {
		result1 models.RouterGroup
		result2 error
	}
	SaveRouterGroupStub        func(routerGroup models.RouterGroup) error
	saveRouterGroupMutex       sync.RWMutex
	saveRouterGroupArgsForCall []struct {
		routerGroup models.RouterGroup
	}
	saveRouterGroupReturns struct {
		result1 error
	}
	CancelWatchesStub        func()
	cancelWatchesMutex       sync.RWMutex
	cancelWatchesArgsForCall []struct{}
	WatchChangesStub        func(watchType string) (<-chan db.Event, <-chan error, context.CancelFunc)
	watchChangesMutex       sync.RWMutex
	watchChangesArgsForCall []struct {
		watchType string
	}
	watchChangesReturns struct {
		result1 <-chan db.Event
		result2 <-chan error
		result3 context.CancelFunc
	}
}

func (fake *FakeDB) ReadRoutes() ([]models.Route, error) {
	fake.readRoutesMutex.Lock()
	fake.readRoutesArgsForCall = append(fake.readRoutesArgsForCall, struct{}{})
	fake.readRoutesMutex.Unlock()
	if fake.ReadRoutesStub != nil {
		return fake.ReadRoutesStub()
	} else {
		return fake.readRoutesReturns.result1, fake.readRoutesReturns.result2
	}
}

func (fake *FakeDB) ReadRoutesCallCount() int {
	fake.readRoutesMutex.RLock()
	defer fake.readRoutesMutex.RUnlock()
	return len(fake.readRoutesArgsForCall)
}

func (fake *FakeDB) ReadRoutesReturns(result1 []models.Route, result2 error) {
	fake.ReadRoutesStub = nil
	fake.readRoutesReturns = struct {
		result1 []models.Route
		result2 error
	}{result1, result2}
}

func (fake *FakeDB) SaveRoute(route models.Route) error {
	fake.saveRouteMutex.Lock()
	fake.saveRouteArgsForCall = append(fake.saveRouteArgsForCall, struct {
		route models.Route
	}{route})
	fake.saveRouteMutex.Unlock()
	if fake.SaveRouteStub != nil {
		return fake.SaveRouteStub(route)
	} else {
		return fake.saveRouteReturns.result1
	}
}

func (fake *FakeDB) SaveRouteCallCount() int {
	fake.saveRouteMutex.RLock()
	defer fake.saveRouteMutex.RUnlock()
	return len(fake.saveRouteArgsForCall)
}

func (fake *FakeDB) SaveRouteArgsForCall(i int) models.Route {
	fake.saveRouteMutex.RLock()
	defer fake.saveRouteMutex.RUnlock()
	return fake.saveRouteArgsForCall[i].route
}

func (fake *FakeDB) SaveRouteReturns(result1 error) {
	fake.SaveRouteStub = nil
	fake.saveRouteReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeDB) DeleteRoute(route models.Route) error {
	fake.deleteRouteMutex.Lock()
	fake.deleteRouteArgsForCall = append(fake.deleteRouteArgsForCall, struct {
		route models.Route
	}{route})
	fake.deleteRouteMutex.Unlock()
	if fake.DeleteRouteStub != nil {
		return fake.DeleteRouteStub(route)
	} else {
		return fake.deleteRouteReturns.result1
	}
}

func (fake *FakeDB) DeleteRouteCallCount() int {
	fake.deleteRouteMutex.RLock()
	defer fake.deleteRouteMutex.RUnlock()
	return len(fake.deleteRouteArgsForCall)
}

func (fake *FakeDB) DeleteRouteArgsForCall(i int) models.Route {
	fake.deleteRouteMutex.RLock()
	defer fake.deleteRouteMutex.RUnlock()
	return fake.deleteRouteArgsForCall[i].route
}

func (fake *FakeDB) DeleteRouteReturns(result1 error) {
	fake.DeleteRouteStub = nil
	fake.deleteRouteReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeDB) ReadTcpRouteMappings() ([]models.TcpRouteMapping, error) {
	fake.readTcpRouteMappingsMutex.Lock()
	fake.readTcpRouteMappingsArgsForCall = append(fake.readTcpRouteMappingsArgsForCall, struct{}{})
	fake.readTcpRouteMappingsMutex.Unlock()
	if fake.ReadTcpRouteMappingsStub != nil {
		return fake.ReadTcpRouteMappingsStub()
	} else {
		return fake.readTcpRouteMappingsReturns.result1, fake.readTcpRouteMappingsReturns.result2
	}
}

func (fake *FakeDB) ReadTcpRouteMappingsCallCount() int {
	fake.readTcpRouteMappingsMutex.RLock()
	defer fake.readTcpRouteMappingsMutex.RUnlock()
	return len(fake.readTcpRouteMappingsArgsForCall)
}

func (fake *FakeDB) ReadTcpRouteMappingsReturns(result1 []models.TcpRouteMapping, result2 error) {
	fake.ReadTcpRouteMappingsStub = nil
	fake.readTcpRouteMappingsReturns = struct {
		result1 []models.TcpRouteMapping
		result2 error
	}{result1, result2}
}

func (fake *FakeDB) ReadFilteredTcpRouteMappings(columnName string, values []string) ([]models.TcpRouteMapping, error) {
	fake.readFilteredTcpRouteMappingsMutex.Lock()
	fake.readFilteredTcpRouteMappingsArgsForCall = append(fake.readFilteredTcpRouteMappingsArgsForCall, struct {
		columnName string
		values     []string
	}{columnName, values})
	fake.readFilteredTcpRouteMappingsMutex.Unlock()
	if fake.ReadFilteredTcpRouteMappingsStub != nil {
		return fake.ReadFilteredTcpRouteMappingsStub(columnName, values)
	} else {
		return fake.readFilteredTcpRouteMappingsReturns.result1, fake.readFilteredTcpRouteMappingsReturns.result2
	}
}

func (fake *FakeDB) ReadFilteredTcpRouteMappingsCallCount() int {
	fake.readFilteredTcpRouteMappingsMutex.RLock()
	defer fake.readFilteredTcpRouteMappingsMutex.RUnlock()
	return len(fake.readFilteredTcpRouteMappingsArgsForCall)
}

func (fake *FakeDB) ReadFilteredTcpRouteMappingsArgsForCall(i int) (string, []string) {
	fake.readFilteredTcpRouteMappingsMutex.RLock()
	defer fake.readFilteredTcpRouteMappingsMutex.RUnlock()
	return fake.readFilteredTcpRouteMappingsArgsForCall[i].columnName, fake.readFilteredTcpRouteMappingsArgsForCall[i].values
}

func (fake *FakeDB) ReadFilteredTcpRouteMappingsReturns(result1 []models.TcpRouteMapping, result2 error) {
	fake.ReadFilteredTcpRouteMappingsStub = nil
	fake.readFilteredTcpRouteMappingsReturns = struct {
		result1 []models.TcpRouteMapping
		result2 error
	}{result1, result2}
}

func (fake *FakeDB) SaveTcpRouteMapping(tcpMapping models.TcpRouteMapping) error {
	fake.saveTcpRouteMappingMutex.Lock()
	fake.saveTcpRouteMappingArgsForCall = append(fake.saveTcpRouteMappingArgsForCall, struct {
		tcpMapping models.TcpRouteMapping
	}{tcpMapping})
	fake.saveTcpRouteMappingMutex.Unlock()
	if fake.SaveTcpRouteMappingStub != nil {
		return fake.SaveTcpRouteMappingStub(tcpMapping)
	} else {
		return fake.saveTcpRouteMappingReturns.result1
	}
}

func (fake *FakeDB) SaveTcpRouteMappingCallCount() int {
	fake.saveTcpRouteMappingMutex.RLock()
	defer fake.saveTcpRouteMappingMutex.RUnlock()
	return len(fake.saveTcpRouteMappingArgsForCall)
}

func (fake *FakeDB) SaveTcpRouteMappingArgsForCall(i int) models.TcpRouteMapping {
	fake.saveTcpRouteMappingMutex.RLock()
	defer fake.saveTcpRouteMappingMutex.RUnlock()
	return fake.saveTcpRouteMappingArgsForCall[i].tcpMapping
}

func (fake *FakeDB) SaveTcpRouteMappingReturns(result1 error) {
	fake.SaveTcpRouteMappingStub = nil
	fake.saveTcpRouteMappingReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeDB) DeleteTcpRouteMapping(tcpMapping models.TcpRouteMapping) error {
	fake.deleteTcpRouteMappingMutex.Lock()
	fake.deleteTcpRouteMappingArgsForCall = append(fake.deleteTcpRouteMappingArgsForCall, struct {
		tcpMapping models.TcpRouteMapping
	}{tcpMapping})
	fake.deleteTcpRouteMappingMutex.Unlock()
	if fake.DeleteTcpRouteMappingStub != nil {
		return fake.DeleteTcpRouteMappingStub(tcpMapping)
	} else {
		return fake.deleteTcpRouteMappingReturns.result1
	}
}

func (fake *FakeDB) DeleteTcpRouteMappingCallCount() int {
	fake.deleteTcpRouteMappingMutex.RLock()
	defer fake.deleteTcpRouteMappingMutex.RUnlock()
	return len(fake.deleteTcpRouteMappingArgsForCall)
}

func (fake *FakeDB) DeleteTcpRouteMappingArgsForCall(i int) models.TcpRouteMapping {
	fake.deleteTcpRouteMappingMutex.RLock()
	defer fake.deleteTcpRouteMappingMutex.RUnlock()
	return fake.deleteTcpRouteMappingArgsForCall[i].tcpMapping
}

func (fake *FakeDB) DeleteTcpRouteMappingReturns(result1 error) {
	fake.DeleteTcpRouteMappingStub = nil
	fake.deleteTcpRouteMappingReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeDB) ReadRouterGroups() (models.RouterGroups, error) {
	fake.readRouterGroupsMutex.Lock()
	fake.readRouterGroupsArgsForCall = append(fake.readRouterGroupsArgsForCall, struct{}{})
	fake.readRouterGroupsMutex.Unlock()
	if fake.ReadRouterGroupsStub != nil {
		return fake.ReadRouterGroupsStub()
	} else {
		return fake.readRouterGroupsReturns.result1, fake.readRouterGroupsReturns.result2
	}
}

func (fake *FakeDB) ReadRouterGroupsCallCount() int {
	fake.readRouterGroupsMutex.RLock()
	defer fake.readRouterGroupsMutex.RUnlock()
	return len(fake.readRouterGroupsArgsForCall)
}

func (fake *FakeDB) ReadRouterGroupsReturns(result1 models.RouterGroups, result2 error) {
	fake.ReadRouterGroupsStub = nil
	fake.readRouterGroupsReturns = struct {
		result1 models.RouterGroups
		result2 error
	}{result1, result2}
}

func (fake *FakeDB) ReadRouterGroup(guid string) (models.RouterGroup, error) {
	fake.readRouterGroupMutex.Lock()
	fake.readRouterGroupArgsForCall = append(fake.readRouterGroupArgsForCall, struct {
		guid string
	}{guid})
	fake.readRouterGroupMutex.Unlock()
	if fake.ReadRouterGroupStub != nil {
		return fake.ReadRouterGroupStub(guid)
	} else {
		return fake.readRouterGroupReturns.result1, fake.readRouterGroupReturns.result2
	}
}

func (fake *FakeDB) ReadRouterGroupCallCount() int {
	fake.readRouterGroupMutex.RLock()
	defer fake.readRouterGroupMutex.RUnlock()
	return len(fake.readRouterGroupArgsForCall)
}

func (fake *FakeDB) ReadRouterGroupArgsForCall(i int) string {
	fake.readRouterGroupMutex.RLock()
	defer fake.readRouterGroupMutex.RUnlock()
	return fake.readRouterGroupArgsForCall[i].guid
}

func (fake *FakeDB) ReadRouterGroupReturns(result1 models.RouterGroup, result2 error) {
	fake.ReadRouterGroupStub = nil
	fake.readRouterGroupReturns = struct {
		result1 models.RouterGroup
		result2 error
	}{result1, result2}
}

func (fake *FakeDB) ReadRouterGroupByName(name string) (models.RouterGroup, error) {
	fake.readRouterGroupByNameMutex.Lock()
	fake.readRouterGroupByNameArgsForCall = append(fake.readRouterGroupByNameArgsForCall, struct {
		name string
	}{name})
	fake.readRouterGroupByNameMutex.Unlock()
	if fake.ReadRouterGroupByNameStub != nil {
		return fake.ReadRouterGroupByNameStub(name)
	} else {
		return fake.readRouterGroupByNameReturns.result1, fake.readRouterGroupByNameReturns.result2
	}
}

func (fake *FakeDB) ReadRouterGroupByNameCallCount() int {
	fake.readRouterGroupByNameMutex.RLock()
	defer fake.readRouterGroupByNameMutex.RUnlock()
	return len(fake.readRouterGroupByNameArgsForCall)
}

func (fake *FakeDB) ReadRouterGroupByNameArgsForCall(i int) string {
	fake.readRouterGroupByNameMutex.RLock()
	defer fake.readRouterGroupByNameMutex.RUnlock()
	return fake.readRouterGroupByNameArgsForCall[i].name
}

func (fake *FakeDB) ReadRouterGroupByNameReturns(result1 models.RouterGroup, result2 error) {
	fake.ReadRouterGroupByNameStub = nil
	fake.readRouterGroupByNameReturns = struct {
		result1 models.RouterGroup
		result2 error
	}{result1, result2}
}

func (fake *FakeDB) SaveRouterGroup(routerGroup models.RouterGroup) error {
	fake.saveRouterGroupMutex.Lock()
	fake.saveRouterGroupArgsForCall = append(fake.saveRouterGroupArgsForCall, struct {
		routerGroup models.RouterGroup
	}{routerGroup})
	fake.saveRouterGroupMutex.Unlock()
	if fake.SaveRouterGroupStub != nil {
		return fake.SaveRouterGroupStub(routerGroup)
	} else {
		return fake.saveRouterGroupReturns.result1
	}
}

func (fake *FakeDB) SaveRouterGroupCallCount() int {
	fake.saveRouterGroupMutex.RLock()
	defer fake.saveRouterGroupMutex.RUnlock()
	return len(fake.saveRouterGroupArgsForCall)
}

func (fake *FakeDB) SaveRouterGroupArgsForCall(i int) models.RouterGroup {
	fake.saveRouterGroupMutex.RLock()
	defer fake.saveRouterGroupMutex.RUnlock()
	return fake.saveRouterGroupArgsForCall[i].routerGroup
}

func (fake *FakeDB) SaveRouterGroupReturns(result1 error) {
	fake.SaveRouterGroupStub = nil
	fake.saveRouterGroupReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeDB) CancelWatches() {
	fake.cancelWatchesMutex.Lock()
	fake.cancelWatchesArgsForCall = append(fake.cancelWatchesArgsForCall, struct{}{})
	fake.cancelWatchesMutex.Unlock()
	if fake.CancelWatchesStub != nil {
		fake.CancelWatchesStub()
	}
}

func (fake *FakeDB) CancelWatchesCallCount() int {
	fake.cancelWatchesMutex.RLock()
	defer fake.cancelWatchesMutex.RUnlock()
	return len(fake.cancelWatchesArgsForCall)
}

func (fake *FakeDB) WatchChanges(watchType string) (<-chan db.Event, <-chan error, context.CancelFunc) {
	fake.watchChangesMutex.Lock()
	fake.watchChangesArgsForCall = append(fake.watchChangesArgsForCall, struct {
		watchType string
	}{watchType})
	fake.watchChangesMutex.Unlock()
	if fake.WatchChangesStub != nil {
		return fake.WatchChangesStub(watchType)
	} else {
		return fake.watchChangesReturns.result1, fake.watchChangesReturns.result2, fake.watchChangesReturns.result3
	}
}

func (fake *FakeDB) WatchChangesCallCount() int {
	fake.watchChangesMutex.RLock()
	defer fake.watchChangesMutex.RUnlock()
	return len(fake.watchChangesArgsForCall)
}

func (fake *FakeDB) WatchChangesArgsForCall(i int) string {
	fake.watchChangesMutex.RLock()
	defer fake.watchChangesMutex.RUnlock()
	return fake.watchChangesArgsForCall[i].watchType
}

func (fake *FakeDB) WatchChangesReturns(result1 <-chan db.Event, result2 <-chan error, result3 context.CancelFunc) {
	fake.WatchChangesStub = nil
	fake.watchChangesReturns = struct {
		result1 <-chan db.Event
		result2 <-chan error
		result3 context.CancelFunc
	}{result1, result2, result3}
}

var _ db.DB = new(FakeDB)
