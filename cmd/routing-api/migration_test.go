package main_test

import (
	"fmt"

	"code.cloudfoundry.org/routing-api"
	"code.cloudfoundry.org/routing-api/cmd/routing-api/testrunner"
	"code.cloudfoundry.org/routing-api/config"
	"code.cloudfoundry.org/routing-api/db"
	"code.cloudfoundry.org/routing-api/models"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/ginkgomon"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = XDescribe("Migrations", func() {
	var (
		cfg               config.Config
		etcd              db.DB
		routingAPIProcess ifrit.Process
		etcdRouterGroups  []models.RouterGroup
	)

	BeforeEach(func() {
		var err error
		cfg, err = config.NewConfigFromFile(routingAPIArgs.ConfigPath, false)
		Expect(err).ToNot(HaveOccurred())
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
		BeforeEach(func() {
			var err error

			for _, rg := range cfg.RouterGroups {
				etcd.SaveRouterGroup(rg)
			}

			etcdRouterGroups, err = etcd.ReadRouterGroups()
			Expect(err).ToNot(HaveOccurred())
		})
	})

	It("migrates all router groups with the original guids", func() {
		client := routing_api.NewClient(fmt.Sprintf("http://127.0.0.1:%d", routingAPIPort), false)
		Eventually(func() error {
			_, err := client.RouterGroups()
			if err != nil {
				fmt.Printf("err", err, err.Error())
			}
			return err
		}, "60s", "1s").Should(BeNil())
		routerGroups, err := client.RouterGroups()
		Expect(err).NotTo(HaveOccurred())
		Expect(len(routerGroups)).To(Equal(1))
		Expect(routerGroups[0].Guid).To(Equal(etcdRouterGroups[0]))
		Expect(routerGroups[0].Name).To(Equal(DefaultRouterGroupName))
		Expect(routerGroups[0].Type).To(Equal(models.RouterGroupType("tcp")))
		Expect(routerGroups[0].ReservablePorts).To(Equal(models.ReservablePorts("1024-65535")))
	})

	XContext("when routes already exist", func() {
		var (
			tcpRoute models.TcpRouteMapping
			route    models.Route
		)
		BeforeEach(func() {
			Expect(len(etcdRouterGroups)).To(Equal(1))
			routerGroupGuid := etcdRouterGroups[0].Guid
			tcpRoute = models.NewTcpRouteMapping(routerGroupGuid, 52001, "1.2.3.5", 60001, 30)
			route = models.NewRoute("a.b.c", 33, "1.1.1.1", "potato", "", 55)

			etcd.SaveTcpRouteMapping(tcpRoute)
			etcd.SaveRoute(route)
		})
		It("migrates all routes", func() {
			tcpRouteMappings, err := client.TcpRouteMappings()
			Expect(err).NotTo(HaveOccurred())
			Expect(tcpRouteMappings).NotTo(BeNil())
			Expect(tcpRouteMappings).To(ContainElement(tcpRoute))

			routes, err := client.Routes()
			Expect(err).ToNot(HaveOccurred())
			Expect(routes).To(HaveLen(1))

			Expect(routes).To(ContainElement(route))
		})
	})
})
