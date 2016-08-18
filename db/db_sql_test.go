package db_test

import (
	"code.cloudfoundry.org/routing-api/config"
	"code.cloudfoundry.org/routing-api/db"
	"code.cloudfoundry.org/routing-api/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SqlDB", func() {
	var (
		cfg   *config.SqlDB
		sqlDB *db.SqlDB
		err   error
	)
	BeforeEach(func() {
		cfg = &config.SqlDB{
			Username: "root",
			Password: "password",
			Schema:   "routing_api",
			Host:     "localhost",
			Port:     3306,
			Type:     "mysql",
		}
	})

	Describe("Connection", func() {
		JustBeforeEach(func() {
			sqlDB, err = db.NewSqlDB(cfg)
		})

		It("returns a sql db client", func() {
			Expect(err).ToNot(HaveOccurred())
			Expect(sqlDB).ToNot(BeNil())
		})

		Context("when config is nil", func() {
			BeforeEach(func() {
				cfg = nil
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(sqlDB).To(BeNil())
			})
		})

		Context("when authentication fails", func() {
			BeforeEach(func() {
				cfg.Password = "wrong_password"
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(sqlDB).To(BeNil())
			})
		})

		Context("when connecting to SQL DB fails", func() {
			BeforeEach(func() {
				cfg.Port = 1234
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(sqlDB).To(BeNil())

			})
		})
	})

	FDescribe("ReadRouterGroups", func() {
		var (
			routerGroups models.RouterGroups
			err          error
			rg           models.RouterGroupDB
		)
		BeforeEach(func() {
			sqlDB, err = db.NewSqlDB(cfg)
		})

		JustBeforeEach(func() {
			routerGroups, err = sqlDB.ReadRouterGroups()
		})

		Context("when there are router groups", func() {
			BeforeEach(func() {
				rg = models.RouterGroupDB{
					Guid:            "some-guid",
					Name:            "rg-1",
					Type:            "tcp",
					ReservablePorts: "120",
				}
				sqlDB.Client.Create(&rg)
			})

			AfterEach(func() {
				sqlDB.Client.Delete(&rg)
			})

			It("returns list of router groups", func() {
				Expect(routerGroups).ToNot(BeNil())
				Expect(len(routerGroups)).To(BeNumerically(">", 0))
				Expect(routerGroups[0].Guid).To(Equal(rg.Guid))
				Expect(routerGroups[0].Name).To(Equal(rg.Name))
				Expect(string(routerGroups[0].ReservablePorts)).To(Equal(rg.ReservablePorts))
				Expect(string(routerGroups[0].Type)).To(Equal(rg.Type))
			})

			Context("when there are no router groups", func() {
				BeforeEach(func() {
					sqlDB.Client.Where("1 = 1").Delete(&models.RouterGroupDB{})
				})

				It("returns an empty slice", func() {
					Expect(routerGroups).ToNot(BeNil())
					Expect(routerGroups).To(HaveLen(0))
				})
			})

			Context("when there is a connection error", func() {

				BeforeEach(func() {
					sqlDB.Client.Where("1 = 1").Delete(&models.RouterGroupDB{})
				})

				It("returns an empty slice", func() {
					Expect(routerGroups).ToNot(BeNil())
					Expect(routerGroups).To(HaveLen(0))
				})
			})
		})
	})
})
