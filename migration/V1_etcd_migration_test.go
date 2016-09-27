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
		etcd           db.DB
		sqlDB          *db.SqlDB
		etcdRunner     *etcdstorerunner.ETCDClusterRunner
		mysqlAllocator testrunner.DbAllocator
		sqlCfg         *config.SqlDB
		etcdConfig     *config.Etcd
	)
	Context("when database connection is successful", func() {
		BeforeEach(func() {
			mysqlAllocator = testrunner.NewMySQLAllocator()
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

			etcdConfig = &config.Etcd{
				RequireSSL: true,
				CertFile:   filepath.Join(basePath, "client.crt"),
				KeyFile:    filepath.Join(basePath, "client.key"),
				CAFile:     filepath.Join(basePath, "etcd-ca.crt"),
				NodeURLS:   etcdRunner.NodeURLS(),
			}

			etcd, err = db.NewETCD(*etcdConfig)
			Expect(err).NotTo(HaveOccurred())

			sqlCfg = &config.SqlDB{
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

			v0Migration := migration.NewV0InitMigration(sqlCfg)
			err = v0Migration.RunMigration()
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			etcdRunner.Reset()
			etcdRunner.Stop()
			etcdRunner.KillWithFire()
			etcdRunner.GoAway()
			mysqlAllocator.Delete()
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
				etcdMigration := migration.NewV1EtcdMigration(sqlCfg, etcdConfig)
				err := etcdMigration.RunMigration()
				Expect(err).ToNot(HaveOccurred())

				rg, err := sqlDB.ReadRouterGroup("1234567890")
				Expect(err).ToNot(HaveOccurred())
				Expect(rg).To(matchers.MatchRouterGroup(savedRouterGroup))
			})
		})

		Context("with http routes in etcd", func() {
			var savedRoute models.Route
			BeforeEach(func() {
				savedRoute = models.NewRoute("/route", 8333, "127.0.0.1", "log_guid", "rs", 10)
				for i := 0; i < 3; i += 1 {
					err := etcd.SaveRoute(savedRoute)
					Expect(err).NotTo(HaveOccurred())
				}
			})

			It("should successfully migrate http routes to mysql", func() {
				etcdMigration := migration.NewV1EtcdMigration(sqlCfg, etcdConfig)
				err := etcdMigration.RunMigration()
				Expect(err).ToNot(HaveOccurred())

				routes, err := sqlDB.ReadRoutes()
				Expect(err).ToNot(HaveOccurred())
				Expect(routes).To(HaveLen(1))
				Expect(routes[0]).To(matchers.MatchHttpRoute(savedRoute))
				Expect(int(routes[0].ModificationTag.Index)).To(Equal(2))
			})
		})

		Context("with tcp routes in etcd", func() {
			var tcpRoute models.TcpRouteMapping
			BeforeEach(func() {
				tcpRoute = models.NewTcpRouteMapping("router-group-guid", 3056, "127.0.0.1", 2990, 5)
				for i := 0; i < 3; i += 1 {
					err := etcd.SaveTcpRouteMapping(tcpRoute)
					Expect(err).NotTo(HaveOccurred())
				}
			})

			It("should successfully migrate tcp routes to mysql", func() {
				etcdMigration := migration.NewV1EtcdMigration(sqlCfg, etcdConfig)
				err := etcdMigration.RunMigration()
				Expect(err).ToNot(HaveOccurred())

				routes, err := sqlDB.ReadTcpRouteMappings()
				Expect(err).ToNot(HaveOccurred())
				Expect(routes).To(HaveLen(1))
				Expect(routes[0]).To(matchers.MatchTcpRoute(tcpRoute))
				Expect(int(routes[0].ModificationTag.Index)).To(Equal(2))
			})
		})
	})
})
