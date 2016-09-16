package test_helpers

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"

	"github.com/cloudfoundry/storeadapter"
	"github.com/cloudfoundry/storeadapter/storerunner/etcdstorerunner"
)

var etcdVersion = "etcdserver\":\"2.1.1"

type DbAllocator interface {
	Create(args ...int) (string, error)
	Reset(id string) error
	Delete(id string) error
}

type mysqlAllocator struct {
	sqlDB *sql.DB
}

func NewMySQLAllocator() DbAllocator {
	return &mysqlAllocator{}
}

type etcdAllocator struct {
	etcdAdapter storeadapter.StoreAdapter
	etcdRunner  *etcdstorerunner.ETCDClusterRunner
}

func NewEtcdAllocator() DbAllocator {
	return &etcdAllocator{}
}

func (a *mysqlAllocator) Create(args ...int) (string, error) {
	var err error
	sqlDBName := fmt.Sprintf("test%d", rand.Int())
	a.sqlDB, err = sql.Open("mysql", "root:password@/")
	if err != nil {
		return "", err
	}
	err = a.sqlDB.Ping()
	if err != nil {
		return "", err
	}

	_, err = a.sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", sqlDBName))
	if err != nil {
		return "", err
	}

	return sqlDBName, nil
}

func (a *mysqlAllocator) Reset(id string) error {
	_, err := a.sqlDB.Exec(fmt.Sprintf("DROP DATABASE %s", id))
	if err != nil {
		return err
	}

	_, err = a.sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", id))
	return err
}

func (a *mysqlAllocator) Delete(id string) error {
	defer a.sqlDB.Close()
	_, err := a.sqlDB.Exec(fmt.Sprintf("DROP DATABASE %s", id))
	return err
}

func (e *etcdAllocator) Create(args ...int) (string, error) {
	e.etcdRunner = etcdstorerunner.NewETCDClusterRunner(args[0], 1, nil)
	e.etcdRunner.Start()

	etcdVersionUrl := e.etcdRunner.NodeURLS()[0] + "/version"
	resp, err := http.Get(etcdVersionUrl)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// response body: {"etcdserver":"2.1.1","etcdcluster":"2.1.0"}
	if !strings.Contains(string(body), etcdVersion) {
		return "", errors.New("Incorrect etcd version")
	}

	e.etcdAdapter = e.etcdRunner.Adapter(nil)
	return "", nil
}

func (e *etcdAllocator) Reset(id string) error {
	e.etcdRunner.Reset()
	return nil
}

func (e *etcdAllocator) Delete(id string) error {
	e.etcdAdapter.Disconnect()
	e.etcdRunner.Reset()
	e.etcdRunner.Stop()
	e.etcdRunner.KillWithFire()
	return nil
}
