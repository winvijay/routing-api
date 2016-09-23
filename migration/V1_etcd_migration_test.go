package migration_test

import (
	"path"
	"path/filepath"

	"code.cloudfoundry.org/routing-api/cmd/routing-api/testrunner"
	"code.cloudfoundry.org/routing-api/config"
	"code.cloudfoundry.org/routing-api/db"
	"code.cloudfoundry.org/routing-api/matchers"
	"code.cloudfoundry.org/routing-api/migration"
	"code.cloudfoundry.org/routing-api/models"

	"github.com/cloudfoundry/storeadapter/storerunner/etcdstorerunner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("V1EtcdMigration", func() {
	var (
		etcd       db.DB
		sqlDB      *db.SqlDB
		etcdRunner *etcdstorerunner.ETCDClusterRunner
	)
	Context("when database connection is successful", func() {
		BeforeEach(func() {
			// wrrite router groups
			// read from mysql and verify the data
			mysqlAllocator := testrunner.NewMySQLAllocator()
			mysqlSchema, err := mysqlAllocator.Create()
			Expect(err).NotTo(HaveOccurred())

			basePath, err := filepath.Abs(path.Join("..", "fixtures", "etcd-certs"))
			Expect(err).NotTo(HaveOccurred())

			serverSSLConfig := &etcdstorerunner.SSLConfig{
				CertFile: filepath.Join(basePath, "server.crt"),
				KeyFile:  filepath.Join(basePath, "server.key"),
				CAFile:   filepath.Join(basePath, "etcd-ca.crt"),
			}

			etcdPort := 4001 + GinkgoParallelNode()
			etcdRunner = etcdstorerunner.NewETCDClusterRunner(etcdPort, 1, serverSSLConfig)
			etcdRunner.Start()

			etcdConfig := config.Etcd{
				RequireSSL: true,
				CertFile:   filepath.Join(basePath, "client.crt"),
				KeyFile:    filepath.Join(basePath, "client.key"),
				CAFile:     filepath.Join(basePath, "etcd-ca.crt"),
				NodeURLS:   etcdRunner.NodeURLS(),
			}

			etcd, err = db.NewETCD(etcdConfig)
			Expect(err).NotTo(HaveOccurred())

			sqlCfg := &config.SqlDB{
				Username: "root",
				Password: "password",
				Schema:   mysqlSchema,
				Host:     "localhost",
				Port:     3306,
				Type:     "mysql",
			}

			dbSQL, err := db.NewSqlDB(sqlCfg)
			Expect(err).ToNot(HaveOccurred())
			sqlDB = dbSQL.(*db.SqlDB)

		})
		AfterEach(func() {
			etcdRunner.Reset()
			etcdRunner.Stop()
			etcdRunner.KillWithFire()
			etcdRunner.GoAway()
			sqlDB.Client.Delete(models.RouterGroupsDB{})
		})

		Context("with router groups in etcd", func() {
			var savedRouterGroup models.RouterGroup
			BeforeEach(func() {
				savedRouterGroup = models.RouterGroup{
					Name:            "router-group-1",
					Type:            "tcp",
					Guid:            "1234567890",
					ReservablePorts: "10-20,25",
				}
				err := etcd.SaveRouterGroup(savedRouterGroup)
				Expect(err).NotTo(HaveOccurred())
			})
			It("should successfully migrate router groups to mysql", func() {
				etcdMigration := migration.NewV1EtcdMigration(sqlDB, etcd)
				err := etcdMigration.RunMigration()
				Expect(err).ToNot(HaveOccurred())

				rg, err := sqlDB.ReadRouterGroup("1234567890")
				Expect(err).ToNot(HaveOccurred())
				Expect(rg).To(matchers.MatchRouterGroup(savedRouterGroup))
			})
		})
	})
})
