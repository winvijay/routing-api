package migration_test

import (
	"fmt"

	"code.cloudfoundry.org/routing-api/config"
	"code.cloudfoundry.org/routing-api/migration"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("V0InitMigration", func() {
	var (
		sqlDBName string
		sqlCfg    *config.SqlDB
	)
	BeforeEach(func() {
		sqlDBName = fmt.Sprintf("test%d", GinkgoParallelNode())
		sqlCfg = &config.SqlDB{
			Username: "root",
			Password: "password",
			Schema:   sqlDBName,
			Host:     "localhost",
			Port:     3306,
			Type:     "mysql",
		}
	})
	Context("when valid sql config is passed", func() {
		var v0Migration *migration.V0InitMigration
		BeforeEach(func() {
			v0Migration = migration.NewV0InitMigration(*sqlCfg)
		})
		AfterEach(func() {
			//cleanup schema
		})
		It("should successfully create correct schema", func() {
			v0Migration.RunMigration()
			//  use sqldb client to verify schema
		})
	})
})
