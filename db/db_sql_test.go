package db_test

import (
	"errors"

	"code.cloudfoundry.org/routing-api/config"
	"code.cloudfoundry.org/routing-api/db"
	"code.cloudfoundry.org/routing-api/db/fakes"
	"code.cloudfoundry.org/routing-api/models"
	"github.com/jinzhu/gorm"
	"github.com/nu7hatch/gouuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SqlDB", func() {
	var (
		sqlDB *db.SqlDB
		err   error
	)
	BeforeEach(func() {
		sqlCfg = &config.SqlDB{
			Username: "root",
			Password: "password",
			Schema:   sqlDBName,
			Host:     "localhost",
			Port:     3306,
			Type:     "mysql",
		}
		dbSQL, err := db.NewSqlDB(sqlCfg)
		Expect(err).ToNot(HaveOccurred())
		sqlDB = dbSQL.(*db.SqlDB)
	})

	Describe("Connection", func() {
		var sqlDB db.DB
		JustBeforeEach(func() {
			sqlDB, err = db.NewSqlDB(sqlCfg)
		})

		It("returns a sql db client", func() {
			Expect(err).ToNot(HaveOccurred())
			Expect(sqlDB).ToNot(BeNil())
		})

		Context("when config is nil", func() {
			BeforeEach(func() {
				sqlCfg = nil
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(sqlDB).To(BeNil())
			})
		})

		Context("when authentication fails", func() {
			BeforeEach(func() {
				sqlCfg.Password = "wrong_password"
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(sqlDB).To(BeNil())
			})
		})

		Context("when connecting to SQL DB fails", func() {
			BeforeEach(func() {
				sqlCfg.Port = 1234
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

		JustBeforeEach(func() {
			routerGroups, err = sqlDB.ReadRouterGroups()
		})

		Context("when there are router groups", func() {
			BeforeEach(func() {
				rg = models.RouterGroupDB{
					Guid:            newUuid(),
					Name:            "rg-1",
					Type:            "tcp",
					ReservablePorts: "120",
				}
				Expect(sqlDB.Client.Create(&rg).Error).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				Expect(sqlDB.Client.Delete(&rg).Error).ToNot(HaveOccurred())
			})

			It("returns list of router groups", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(routerGroups).ToNot(BeNil())
				Expect(len(routerGroups)).To(BeNumerically(">", 0))
				Expect(routerGroups).Should(ContainElement(rg.ToRouterGroup()))
			})
		})

		Context("when there are no router groups", func() {
			BeforeEach(func() {
				sqlDB.Client.Where("1 = 1").Delete(&models.RouterGroupDB{})
			})

			It("returns an empty slice", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(routerGroups).ToNot(BeNil())
				Expect(routerGroups).To(HaveLen(0))
			})
		})

		Context("when there is a connection error", func() {
			BeforeEach(func() {
				fakeClient := &fakes.FakeClient{}
				fakeClient.FindReturns(&gorm.DB{Error: errors.New("connection refused")})
				sqlDB.Client = fakeClient
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
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

		JustBeforeEach(func() {
			routerGroup, err = sqlDB.ReadRouterGroup(routerGroupId)
		})

		Context("when router group exists", func() {
			BeforeEach(func() {
				routerGroupId = newUuid()
				rg = models.RouterGroupDB{
					Guid:            routerGroupId,
					Name:            "rg-1",
					Type:            "tcp",
					ReservablePorts: "120",
				}
				Expect(sqlDB.Client.Create(&rg).Error).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				Expect(sqlDB.Client.Delete(&rg).Error).ToNot(HaveOccurred())
			})

			It("returns the router group", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(routerGroup.Guid).To(Equal(rg.Guid))
				Expect(routerGroup.Name).To(Equal(rg.Name))
				Expect(string(routerGroup.ReservablePorts)).To(Equal(rg.ReservablePorts))
				Expect(string(routerGroup.Type)).To(Equal(rg.Type))
			})
		})

		Context("when router group doesn't exist", func() {
			BeforeEach(func() {
				routerGroupId = newUuid()
			})

			It("returns an empty struct", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(routerGroup).To(Equal(models.RouterGroup{}))
			})
		})
	})

	Describe("SaveRouterGroup", func() {
		var (
			routerGroup   models.RouterGroup
			err           error
			routerGroupId string
		)
		BeforeEach(func() {
			routerGroupId = newUuid()
			routerGroup = models.RouterGroup{
				Guid:            routerGroupId,
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
					Guid:            routerGroupId,
					Name:            "rg-1",
					Type:            "tcp",
					ReservablePorts: "120",
				})
			})

			AfterEach(func() {
				sqlDB.Client.Delete(&models.RouterGroupDB{
					Guid: routerGroupId,
				})
			})

			It("updates the existing router group", func() {
				Expect(err).ToNot(HaveOccurred())
				rg, err := sqlDB.ReadRouterGroup(routerGroup.Guid)
				Expect(err).ToNot(HaveOccurred())

				Expect(rg.Guid).To(Equal(routerGroup.Guid))
				Expect(rg.Name).To(Equal(routerGroup.Name))
				Expect(rg.ReservablePorts).To(Equal(routerGroup.ReservablePorts))
				Expect(rg.Type).To(Equal(routerGroup.Type))
			})
		})

		Context("when router group doesn't exist", func() {
			It("creates the router group", func() {
				Expect(err).ToNot(HaveOccurred())
				rg, err := sqlDB.ReadRouterGroup(routerGroup.Guid)
				Expect(err).ToNot(HaveOccurred())
				Expect(rg.Guid).To(Equal(routerGroup.Guid))
				Expect(rg.Name).To(Equal(routerGroup.Name))
				Expect(rg.ReservablePorts).To(Equal(routerGroup.ReservablePorts))
				Expect(rg.Type).To(Equal(routerGroup.Type))
			})
		})
	})

	Describe("SaveTcpRouteMapping", func() {
		var (
			err           error
			routerGroupId string
			tcpRoute      models.TcpRouteMapping
		)
		BeforeEach(func() {
			routerGroupId = newUuid()
			tcpRoute = models.NewTcpRouteMapping(routerGroupId, 3056, "127.0.0.1", 2990, 5)
			tcpRoute.ModificationTag = models.ModificationTag{Guid: "some-tag", Index: 10}
			// tcpRoute = models.TcpRouteMapping{
			// 	RouteModel:      models.RouteModel{ExpiresAt: time.Now().Add(2 * time.Second)},
			// 	TcpRoute:        models.TcpRoute{RouterGroupGuid: routerGroupId, ExternalPort: 2990},
			// 	HostPort:        3056,
			// 	HostIP:          "127.0.0.1",
			// 	ModificationTag:
			// }
		})

		JustBeforeEach(func() {
			err = sqlDB.SaveTcpRouteMapping(tcpRoute)
		})

		AfterEach(func() {
			sqlDB.Client.Delete(&tcpRoute)
		})
		Context("when tcp route exists", func() {
			BeforeEach(func() {
				sqlDB.Client.Create(&tcpRoute)
				tcpRoute.ModificationTag.Index = 15
			})

			It("updates the existing router group", func() {
				Expect(err).ToNot(HaveOccurred())
				var dbTcpRoute models.TcpRouteMapping
				sqlDB.Client.Where("host_ip = ?", "127.0.0.1").First(&dbTcpRoute)
				Expect(dbTcpRoute).ToNot(BeNil())
				Expect(dbTcpRoute.ModificationTag.Index).To(BeNumerically("==", 15))
			})
		})

		Context("when tcp route doesn't exist", func() {
			It("creates a tcp route", func() {
				Expect(err).ToNot(HaveOccurred())
				var dbTcpRoute models.TcpRouteMapping
				err = sqlDB.Client.Where("host_ip = ?", "127.0.0.1").First(&dbTcpRoute).Error
				Expect(err).ToNot(HaveOccurred())
				Expect(dbTcpRoute.TcpRouteMappingEntity).To(Equal(tcpRoute.TcpRouteMappingEntity))
			})
		})

	})

	// Describe("ReadTcpRouteMappings", func() {
	// 	var (
	// 		err           error
	// 		routerGroupId string
	// 		tcpRoute      models.TcpRouteMapping
	// 		tcpRoutes     []models.TcpRouteMapping
	// 	)
	// 	BeforeEach(func() {
	// 		routerGroupId = newUuid()
	// 		tcpRoute = models.TcpRouteMapping{
	// 			TcpRouteMappingEntity: TcpRouteMappingEntity{
	// 				TcpRoute:        models.TcpRoute{RouterGroupGuid: routerGroupId, ExternalPort: 2990},
	// 				HostPort:        3056,
	// 				HostIP:          "127.0.0.1",
	// 				ModificationTag: models.ModificationTag{Guid: "some-tag", Index: 10},
	// 			},
	// 		}
	// 	})

	// 	JustBeforeEach(func() {
	// 		tcpRoutes, err = sqlDB.ReadTcpRouteMappings()
	// 	})

	// 	Context("when at least one tcp route exists", func() {
	// 		BeforeEach(func() {
	// 			Expect(sqlDB.Client.Create(&tcpRoute).Error).ToNot(HaveOccurred())
	// 		})

	// 		AfterEach(func() {
	// 			Expect(sqlDB.Client.Delete(&tcpRoute).Error).ToNot(HaveOccurred())
	// 		})

	// 		It("returns the tcp routes", func() {
	// 			Expect(err).ToNot(HaveOccurred())
	// 			Expect(tcpRoutes).To(ContainElement(tcpRoute))
	// 		})
	// 	})

	// 	Context("when tcp route doesn't exist", func() {

	// 		It("returns an empty array", func() {
	// 			Expect(err).ToNot(HaveOccurred())
	// 			Expect(tcpRoutes).To(Equal([]models.TcpRouteMapping{}))
	// 		})
	// 	})
	// })

	// Describe("DeleteTcpRouteMapping", func() {
	// 	var (
	// 		err           error
	// 		routerGroupId string
	// 		tcpRoute      models.TcpRouteMapping
	// 		tcpRoutes     []models.TcpRouteMapping
	// 	)
	// 	BeforeEach(func() {
	// 		routerGroupId = newUuid()
	// 		tcpRoute = models.TcpRouteMapping{
	// 			TcpRoute:        models.TcpRoute{RouterGroupGuid: routerGroupId, ExternalPort: 2990},
	// 			HostPort:        3056,
	// 			HostIP:          "127.0.0.1",
	// 			ModificationTag: models.ModificationTag{Guid: "some-tag", Index: 10},
	// 		}
	// 	})

	// 	JustBeforeEach(func() {
	// 		err = sqlDB.DeleteTcpRouteMapping(tcpRoute)
	// 	})

	// 	Context("when at least one tcp route exists", func() {
	// 		BeforeEach(func() {
	// 			Expect(sqlDB.Client.Create(&tcpRoute).Error).ToNot(HaveOccurred())
	// 		})

	// 		AfterEach(func() {
	// 			Expect(sqlDB.Client.Delete(&tcpRoute).Error).ToNot(HaveOccurred())
	// 		})

	// 		It("returns the tcp routes", func() {
	// 			Expect(err).ToNot(HaveOccurred())

	// 			tcpRoutes, err = sqlDB.ReadTcpRouteMappings()
	// 			Expect(err).ToNot(HaveOccurred())
	// 			Expect(tcpRoutes).ToNot(ContainElement(tcpRoute))
	// 		})
	// 	})

	// 	Context("when the tcp route doesn't exist", func() {

	// 		It("returns an error", func() {
	// 			Expect(err).To(HaveOccurred())
	// 			Expect(err).Should(MatchError(db.DeleteError))
	// 		})
	// 	})
	// })
	Describe("Methods not implemented", func() {
		It("returns an error", func() {
			err := sqlDB.SaveRoute(models.Route{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("function not implemented:"))
			Expect(err.Error()).To(ContainSubstring("SaveRoute"))
		})
	})
})

func newUuid() string {
	u, err := uuid.NewV4()
	Expect(err).ToNot(HaveOccurred())
	return u.String()
}
