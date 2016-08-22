package db_test

import (
	"fmt"

	"code.cloudfoundry.org/routing-api/config"
	"code.cloudfoundry.org/routing-api/db"
	"code.cloudfoundry.org/routing-api/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("SqlDB", func() {
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

	Describe("ReadRouterGroups", func() {
		var (
			routerGroups models.RouterGroups
			err          error
			rg           models.RouterGroupDB
		)
		BeforeEach(func() {
			sqlDB, err = db.NewSqlDB(cfg)
			Expect(err).ToNot(HaveOccurred())
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
			})

			It("returns an error", func() {
			})
		})
	})

	Describe("ReadRouterGroup", func() {
		var (
			routerGroup   models.RouterGroup
			err           error
			rg            models.RouterGroupDB
			routerGroupId string
		)
		BeforeEach(func() {
			sqlDB, err = db.NewSqlDB(cfg)
			Expect(err).ToNot(HaveOccurred())
		})

		JustBeforeEach(func() {
			routerGroup, err = sqlDB.ReadRouterGroup(routerGroupId)
		})

		Context("when router group exists", func() {
			BeforeEach(func() {
				routerGroupId = "some-guid"
				rg = models.RouterGroupDB{
					Guid:            routerGroupId,
					Name:            "rg-1",
					Type:            "tcp",
					ReservablePorts: "120",
				}
				sqlDB.Client.Create(&rg)
			})

			AfterEach(func() {
				sqlDB.Client.Delete(&rg)
			})

			It("returns the router group", func() {
				Expect(routerGroup.Guid).To(Equal(rg.Guid))
				Expect(routerGroup.Name).To(Equal(rg.Name))
				Expect(string(routerGroup.ReservablePorts)).To(Equal(rg.ReservablePorts))
				Expect(string(routerGroup.Type)).To(Equal(rg.Type))
			})
		})

		Context("when router group doesn't exist", func() {
			BeforeEach(func() {
				routerGroupId = "some-other-guid"
			})

			It("returns an empty struct", func() {
				Expect(routerGroup).To(Equal(models.RouterGroup{}))
			})
		})

		Context("when there is a connection error", func() {
			BeforeEach(func() {
			})

			It("returns an error", func() {
			})
		})
	})

	FDescribe("SaveRouterGroup", func() {
		var (
			routerGroup models.RouterGroup
			err         error
		)
		BeforeEach(func() {
			sqlDB, err = db.NewSqlDB(cfg)
			Expect(err).ToNot(HaveOccurred())
			routerGroup = models.RouterGroup{
				Guid:            "some-guid",
				Name:            "router-group-1",
				Type:            "tcp",
				ReservablePorts: "65000-65002",
			}
		})

		JustBeforeEach(func() {
			err = sqlDB.SaveRouterGroup(routerGroup)
		})

		Context("when router group exists", func() {
			BeforeEach(func() {
				sqlDB.Client.Create(&models.RouterGroupDB{
					Guid:            "some-guid",
					Name:            "rg-1",
					Type:            "tcp",
					ReservablePorts: "120",
				})
			})

			AfterEach(func() {
				sqlDB.Client.Delete(&models.RouterGroupDB{
					Guid: "some-guid",
				})

			})

			It("updates the existing router group", func() {
				rg, err := sqlDB.ReadRouterGroup(routerGroup.Guid)
				Expect(err).ToNot(HaveOccurred())

				Expect(rg.Guid).To(Equal(routerGroup.Guid))
				Expect(rg.Name).To(Equal(routerGroup.Name))
				Expect(rg.ReservablePorts).To(Equal(routerGroup.ReservablePorts))
				Expect(rg.Type).To(Equal(routerGroup.Type))
			})
		})

		FContext("when router group doesn't exist", func() {
			It("creates the router group", func() {
				rg, err := sqlDB.ReadRouterGroup(routerGroup.Guid)
				fmt.Printf("created ...............%#v", rg)
				Expect(err).ToNot(HaveOccurred())
				Expect(rg.Guid).To(Equal(routerGroup.Guid))
				Expect(rg.Name).To(Equal(routerGroup.Name))
				Expect(rg.ReservablePorts).To(Equal(routerGroup.ReservablePorts))
				Expect(rg.Type).To(Equal(routerGroup.Type))
			})
		})

		Context("when there is a connection error", func() {
			BeforeEach(func() {
			})

			It("returns an error", func() {
			})
		})
	})
})
