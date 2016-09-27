package migration_test

import (
	"fmt"

	"code.cloudfoundry.org/routing-api/cmd/routing-api/testrunner"
	"code.cloudfoundry.org/routing-api/config"
	"code.cloudfoundry.org/routing-api/db"
	"code.cloudfoundry.org/routing-api/migration"
	"code.cloudfoundry.org/routing-api/migration/fakes"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Migration", func() {
	var (
		sqlDB                 *db.SqlDB
		sqlCfg                *config.SqlDB
		mysqlAllocator        testrunner.DbAllocator
		fakeMigration         *fakes.FakeMigration
		fakeLastMigration     *fakes.FakeMigration
		migrations            []migration.Migration
		lastMigrationVersion  int = 10
		firstMigrationVersion int = 1
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

		fakeMigration = new(fakes.FakeMigration)
		fakeLastMigration = new(fakes.FakeMigration)

		fakeMigration.VersionReturns(firstMigrationVersion)
		fakeLastMigration.VersionReturns(lastMigrationVersion)
		migrations = []migration.Migration{}
		migrations = append(migrations, fakeLastMigration)
	})

	AfterEach(func() {
		mysqlAllocator.Delete()
		fmt.Println("RUNNING THE AFTER EACH")
	})

	Context("when no migration table exists", func() {
		It("should create the migration table and set the target version to last migration version", func() {
			err := migration.RunMigrations(sqlCfg, migrations)
			Expect(err).ToNot(HaveOccurred())
			gormClient := sqlDB.Client.(*gorm.DB)
			Expect(gormClient.HasTable(&migration.MigrationData{})).To(BeTrue())

			var migrationVersions []migration.MigrationData
			gormClient.Find(&migrationVersions)

			Expect(migrationVersions).To(HaveLen(1))

			migrationVersion := migrationVersions[0]
			Expect(migrationVersion.MigrationKey).To(Equal(migration.MigrationKey))
			Expect(migrationVersion.CurrentVersion).To(Equal(lastMigrationVersion))
			Expect(migrationVersion.TargetVersion).To(Equal(lastMigrationVersion))
		})
		It("should run all the migrations up to the current version", func() {
			err := migration.RunMigrations(sqlCfg, migrations)
			Expect(err).ToNot(HaveOccurred())
			Expect(fakeMigration.RunMigrationCallCount()).To(Equal(len(migrations)))
		})
	})

	Context("when a migration table exists", func() {
		BeforeEach(func() {
			gormClient := sqlDB.Client.(*gorm.DB)
			gormClient.AutoMigrate(&migration.MigrationData{})
		})

		Context("when a migration is necessary", func() {
			Context("when another routing-api has already started the migration", func() {
				BeforeEach(func() {
					migrationData := migration.MigrationData{
						MigrationKey:   migration.MigrationKey,
						CurrentVersion: -1,
						TargetVersion:  lastMigrationVersion,
					}

					err := sqlDB.Client.Create(migrationData).Error
					Expect(err).ToNot(HaveOccurred())
				})

				It("should not update the migration data", func() {
					err := migration.RunMigrations(sqlCfg, migrations)
					Expect(err).ToNot(HaveOccurred())

					var migrationVersions []migration.MigrationData
					sqlDB.Client.Find(&migrationVersions)

					Expect(migrationVersions).To(HaveLen(1))

					migrationVersion := migrationVersions[0]
					Expect(migrationVersion.MigrationKey).To(Equal(migration.MigrationKey))
					Expect(migrationVersion.CurrentVersion).To(Equal(-1))
					Expect(migrationVersion.TargetVersion).To(Equal(lastMigrationVersion))
				})

				It("should not run any migrations", func() {
					err := migration.RunMigrations(sqlCfg, migrations)
					Expect(err).ToNot(HaveOccurred())

					Expect(fakeMigration.RunMigrationCallCount()).To(BeZero())
				})
			})

			Context("when the migration has not been started", func() {
				BeforeEach(func() {
					migrationData := migration.MigrationData{
						MigrationKey:   migration.MigrationKey,
						CurrentVersion: 1,
						TargetVersion:  1,
					}

					err := sqlDB.Client.Create(migrationData).Error
					Expect(err).ToNot(HaveOccurred())
				})

				It("should update the migration data with the target version", func() {
					err := migration.RunMigrations(sqlCfg, migrations)
					Expect(err).ToNot(HaveOccurred())

					var migrationVersions []migration.MigrationData
					sqlDB.Client.Find(&migrationVersions)

					Expect(migrationVersions).To(HaveLen(1))

					migrationVersion := migrationVersions[0]
					Expect(migrationVersion.MigrationKey).To(Equal(migration.MigrationKey))
					Expect(migrationVersion.CurrentVersion).To(Equal(lastMigrationVersion))
					Expect(migrationVersion.TargetVersion).To(Equal(lastMigrationVersion))
				})

				It("should run all the migrations up to the current version", func() {
					err := migration.RunMigrations(sqlCfg, migrations)
					Expect(err).ToNot(HaveOccurred())
					Expect(fakeMigration.RunMigrationCallCount()).To(Equal(0))
					Expect(fakeLastMigration.RunMigrationCallCount()).To(Equal(1))
				})
			})
		})

		Context("when a migration is unnecessary", func() {
			BeforeEach(func() {
				migrationData := migration.MigrationData{
					MigrationKey:   migration.MigrationKey,
					CurrentVersion: lastMigrationVersion,
					TargetVersion:  lastMigrationVersion,
				}

				err := sqlDB.Client.Create(migrationData).Error
				Expect(err).ToNot(HaveOccurred())
			})

			It("should not update the migration data", func() {
				err := migration.RunMigrations(sqlCfg, migrations)
				Expect(err).ToNot(HaveOccurred())

				var migrationVersions []migration.MigrationData
				sqlDB.Client.Find(&migrationVersions)

				Expect(migrationVersions).To(HaveLen(1))

				migrationVersion := migrationVersions[0]
				Expect(migrationVersion.MigrationKey).To(Equal(migration.MigrationKey))
				Expect(migrationVersion.CurrentVersion).To(Equal(lastMigrationVersion))
				Expect(migrationVersion.TargetVersion).To(Equal(lastMigrationVersion))
			})

			It("should not run any migrations", func() {
				err := migration.RunMigrations(sqlCfg, migrations)
				Expect(err).ToNot(HaveOccurred())

				Expect(fakeMigration.RunMigrationCallCount()).To(BeZero())
				Expect(fakeLastMigration.RunMigrationCallCount()).To(BeZero())
			})
		})
	})
})
