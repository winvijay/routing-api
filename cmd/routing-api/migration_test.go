package main_test

import (
	"fmt"

	"code.cloudfoundry.org/routing-api"
	"code.cloudfoundry.org/routing-api/config"
	"code.cloudfoundry.org/routing-api/db"
	"code.cloudfoundry.org/routing-api/models"
	"github.com/onsi/ginkgo/ginkgo/testrunner"
	"github.com/tedsuo/ifrit/ginkgomon"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Migrations", func() {
	var (
		cfg  config.Config
		etcd db.DB
	)

	BeforeEach(func() {
		cfg = config.NewConfigFromFile(routingAPIArgs.ConfigPath)

		etcd, err := db.NewETCD(cfg.Etcd)
		Expect(err).NotTo(HaveOccurred())

		err = etcd.Connect()
		Expect(err).ToNot(HaveOccurred())
	})

	JustBeforeEach(func() {
		routingAPIRunner := testrunner.New(routingAPIBinPath, routingAPIArgs)
		routingAPIProcess = ginkgomon.Invoke(routingAPIRunner)
	})

	AfterEach(func() {
		ginkgomon.Kill(routingAPIProcess)
	})

	Context("when etcd already has router groups", func() {
		var etcdRouterGroups []models.RouterGroup
		BeforeEach(func() {

			for _, rg := range cfg.RouterGroups {
				etcd.SaveRouterGroup(rg)
			}

			etcdRouterGroups = etcd.ReadRouterGroups()
		})
	})

	It("migrates all router groups with the original guids", func() {
		client := routing_api.NewClient(fmt.Sprintf("http://127.0.0.1:%d", routingAPIPort), false)
		Eventually(func() error {
			_, err := client.RouterGroups()
			return err
		}, "30s", "1s")
		routerGroups, err := client.RouterGroups()
		Expect(err).NotTo(HaveOccurred())
		Expect(len(routerGroups)).To(Equal(1))
		Expect(routerGroups[0].Guid).ToEqual(etcdRouterGroups[0])
		Expect(routerGroups[0].Name).To(Equal(DefaultRouterGroupName))
		Expect(routerGroups[0].Type).To(Equal(models.RouterGroupType("tcp")))
		Expect(routerGroups[0].ReservablePorts).To(Equal(models.ReservablePorts("1024-65535")))
	})

	Context("when routes already exist", func() {
		var (
			tcpRoute models.TcpRouteMapping
			route    models.Route
		)
		BeforeEach(func() {
			tcpRoute = models.NewTcpRouteMapping(routerGroupGuid, 52001, "1.2.3.5", 60001, 30)
			route = models.NewRoute("a.b.c", 33, "1.1.1.1", "potato", "", 55)

			etcd.SaveTcpRouteMapping(tcpRoute)
		})
		It("migrates all the tcp routes", func() {
			tcpRouteMappingsResponse, err := client.TcpRouteMappings()
			Expect(err).NotTo(HaveOccurred())
			Expect(tcpRouteMappingsResponse).NotTo(BeNil())
			mappings := TcpRouteMappings(tcpRouteMappingsResponse)
			Expect(mappings.ContainsAll(tcpRoute)).To(BeTrue())
		})

		It("migrates all the http routes", func() {
			routes, getErr = client.Routes()
			Expect(getErr).ToNot(HaveOccurred())
			Expect(routes).To(HaveLen(1))

			Expect(Routes(routes).ContainsAll(route1)).To(BeTrue())
		})

	})
})
