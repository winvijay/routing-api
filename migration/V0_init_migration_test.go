package migration_test

import (
	"code.cloudfoundry.org/routing-api/cmd/routing-api/testrunner"
	"code.cloudfoundry.org/routing-api/config"
	"code.cloudfoundry.org/routing-api/db"
	"code.cloudfoundry.org/routing-api/migration"
	"code.cloudfoundry.org/routing-api/models"

	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("V0InitMigration", func() {
	var (
		sqlDB          *db.SqlDB
		sqlCfg         *config.SqlDB
		mysqlAllocator testrunner.DbAllocator
	)
	BeforeEach(func() {
		mysqlAllocator = testrunner.NewMySQLAllocator()
		mysqlSchema, err := mysqlAllocator.Create()
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
	})

	AfterEach(func() {
		mysqlAllocator.Delete()
	})

	Context("when valid sql config is passed", func() {
		var v0Migration *migration.V0InitMigration
		BeforeEach(func() {
			v0Migration = migration.NewV0InitMigration(sqlCfg)
		})

		It("should successfully create correct schema", func() {
			err := v0Migration.RunMigration()
			Expect(err).ToNot(HaveOccurred())
			gormClient := sqlDB.Client.(*gorm.DB)
			Expect(gormClient.HasTable(&models.RouterGroupDB{})).To(BeTrue())
			Expect(gormClient.HasTable(&models.TcpRouteMapping{})).To(BeTrue())
			Expect(gormClient.HasTable(&models.Route{})).To(BeTrue())
		})
	})
})
