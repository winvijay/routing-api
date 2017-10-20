package config_test

import (
	"errors"
	"time"

	"code.cloudfoundry.org/locket"
	"code.cloudfoundry.org/routing-api/config"
	"code.cloudfoundry.org/routing-api/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Describe("NewConfigFromFile", func() {
		Context("when auth is enabled", func() {
			Context("when the file exists", func() {
				It("returns a valid Config struct", func() {
					cfg_file := "../example_config/example.yml"
					cfg, err := config.NewConfigFromFile(cfg_file, false)

					Expect(err).NotTo(HaveOccurred())
					Expect(cfg.AdminSocket).To(Equal("/some/path/to/socket.sock"))
					Expect(cfg.LogGuid).To(Equal("my_logs"))
					Expect(cfg.MetronConfig.Address).To(Equal("1.2.3.4"))
					Expect(cfg.MetronConfig.Port).To(Equal("4567"))
					Expect(cfg.StatsdClientFlushInterval).To(Equal(10 * time.Millisecond))
					Expect(cfg.OAuth.TokenEndpoint).To(Equal("127.0.0.1"))
					Expect(cfg.OAuth.Port).To(Equal(3000))
					Expect(cfg.OAuth.CACerts).To(Equal("some-ca-cert"))
					Expect(cfg.SystemDomain).To(Equal("example.com"))
					Expect(cfg.SqlDB.Username).To(Equal("username"))
					Expect(cfg.SqlDB.Password).To(Equal("password"))
					Expect(cfg.SqlDB.Port).To(Equal(1234))
					Expect(cfg.SqlDB.CACert).To(Equal("some CA cert"))
					Expect(cfg.SqlDB.SkipSSLValidation).To(Equal(cfg.OAuth.SkipSSLValidation))
					Expect(cfg.MaxTTL).To(Equal(2 * time.Minute))
					Expect(cfg.ConsulCluster.Servers).To(Equal("http://localhost:5678"))
					Expect(cfg.ConsulCluster.LockTTL).To(Equal(10 * time.Second))
					Expect(cfg.ConsulCluster.RetryInterval).To(Equal(5 * time.Second))
					Expect(cfg.SkipConsulLock).To(BeTrue())
					Expect(cfg.Locket.LocketAddress).To(Equal("http://localhost:5678"))
					Expect(cfg.Locket.LocketCACertFile).To(Equal("some-locket-ca-cert"))
					Expect(cfg.Locket.LocketClientCertFile).To(Equal("some-locket-client-cert"))
					Expect(cfg.Locket.LocketClientKeyFile).To(Equal("some-locket-client-key"))
				})

				Context("when there is no token endpoint specified", func() {
					It("returns an error", func() {
						cfg_file := "../example_config/missing_uaa_url.yml"
						_, err := config.NewConfigFromFile(cfg_file, false)
						Expect(err).To(HaveOccurred())
					})
				})
			})

			Context("when the file does not exists", func() {
				It("returns an error", func() {
					cfg_file := "notexist"
					_, err := config.NewConfigFromFile(cfg_file, false)

					Expect(err).To(HaveOccurred())
				})
			})
		})

		Context("when auth is disabled", func() {
			Context("when the file exists", func() {
				It("returns a valid config", func() {
					cfg_file := "../example_config/example.yml"
					cfg, err := config.NewConfigFromFile(cfg_file, true)

					Expect(err).NotTo(HaveOccurred())
					Expect(cfg.LogGuid).To(Equal("my_logs"))
					Expect(cfg.MetronConfig.Address).To(Equal("1.2.3.4"))
					Expect(cfg.MetronConfig.Port).To(Equal("4567"))
					Expect(cfg.StatsdClientFlushInterval).To(Equal(10 * time.Millisecond))
					Expect(cfg.OAuth.TokenEndpoint).To(Equal("127.0.0.1"))
					Expect(cfg.OAuth.Port).To(Equal(3000))
					Expect(cfg.OAuth.CACerts).To(Equal("some-ca-cert"))
				})

				Context("when there is no token endpoint url", func() {
					It("returns a valid config", func() {
						cfg_file := "../example_config/missing_uaa_url.yml"
						cfg, err := config.NewConfigFromFile(cfg_file, true)

						Expect(err).NotTo(HaveOccurred())
						Expect(cfg.LogGuid).To(Equal("my_logs"))
						Expect(cfg.MetronConfig.Address).To(Equal("1.2.3.4"))
						Expect(cfg.MetronConfig.Port).To(Equal("4567"))
						Expect(cfg.DebugAddress).To(Equal("1.2.3.4:1234"))
						Expect(cfg.MaxTTL).To(Equal(2 * time.Minute))
						Expect(cfg.StatsdClientFlushInterval).To(Equal(10 * time.Millisecond))
						Expect(cfg.OAuth.TokenEndpoint).To(BeEmpty())
						Expect(cfg.OAuth.Port).To(Equal(0))
					})
				})
			})
		})
	})

	Describe("Initialize", func() {
		var (
			cfg *config.Config
		)

		BeforeEach(func() {
			cfg = &config.Config{}
		})
		Context("when UUID property is set", func() {
			testConfig := func() string {
				return `log_guid: "my_logs"
admin_socket: "/some/path"
metrics_reporting_interval: "500ms"
uuid: "fake-uuid"
statsd_endpoint: "localhost:8125"
statsd_client_flush_interval: "10ms"
system_domain: "example.com"
router_groups:
- name: router-group-2
  reservable_ports: 1024-10000,42000
  type: udp
consul_cluster:
  url: "http://localhost:4222"
`
			}
			It("populates the value", func() {
				config := testConfig()
				err := cfg.Initialize([]byte(config), true)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.UUID).To(Equal("fake-uuid"))
			})
		})
		Context("when UUID property is not set", func() {
			testConfig := func() string {
				return `log_guid: "my_logs"
admin_socket: "/some/path"
metrics_reporting_interval: "500ms"
statsd_endpoint: "localhost:8125"
statsd_client_flush_interval: "10ms"
system_domain: "example.com"
router_groups:
- name: router-group-2
  reservable_ports: 1024-10000,42000
  type: udp
consul_cluster:
  url: "http://localhost:4222"
`
			}
			It("populates the value", func() {
				config := testConfig()
				err := cfg.Initialize([]byte(config), true)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(errors.New("No UUID is specified")))
			})
		})
		Context("when AdminSocket property is set", func() {
			testConfig := func() string {
				return `admin_socket: "/some/path"
log_guid: "my_logs"
metrics_reporting_interval: "500ms"
uuid: "fake-uuid"
statsd_endpoint: "localhost:8125"
statsd_client_flush_interval: "10ms"
system_domain: "example.com"
router_groups:
- name: router-group-2
  reservable_ports: 1024-10000,42000
  type: udp
consul_cluster:
  url: "http://localhost:4222"
`
			}
			It("populates the value", func() {
				config := testConfig()
				err := cfg.Initialize([]byte(config), true)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.AdminSocket).To(Equal("/some/path"))
			})
		})
		Context("when AdminSocket property is not set", func() {
			testConfig := func() string {
				return `log_guid: "my_logs"
metrics_reporting_interval: "500ms"
uuid: "fake-uuid"
statsd_endpoint: "localhost:8125"
statsd_client_flush_interval: "10ms"
system_domain: "example.com"
router_groups:
- name: router-group-2
  reservable_ports: 1024-10000,42000
  type: udp
consul_cluster:
  url: "http://localhost:4222"
`
			}
			It("returns an error", func() {
				config := testConfig()
				err := cfg.Initialize([]byte(config), true)
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when consul properties are not set", func() {
			testConfig := func() string {
				return `log_guid: "my_logs"
admin_socket: "/some/path"
metrics_reporting_interval: "500ms"
statsd_endpoint: "localhost:8125"
statsd_client_flush_interval: "10ms"
uuid: "fake-uuid"
system_domain: "example.com"
router_groups:
- name: router-group-2
  reservable_ports: 1024-10000,42000
  type: udp
consul_cluster:
  url: "http://localhost:4222"
`
			}
			It("populates the default value for LockTTL from locket library", func() {
				config := testConfig()
				err := cfg.Initialize([]byte(config), true)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.ConsulCluster.LockTTL).To(Equal(locket.DefaultSessionTTL))
			})
			It("populates the default value for RetryInterval from locket library", func() {
				config := testConfig()
				err := cfg.Initialize([]byte(config), true)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.ConsulCluster.RetryInterval).To(Equal(locket.RetryInterval))
			})
		})

		Context("when multiple router groups are seeded", func() {
			var expectedGroups models.RouterGroups

			testConfig := func(name string) string {
				return `log_guid: "my_logs"
admin_socket: "/some/path"
metrics_reporting_interval: "500ms"
statsd_endpoint: "localhost:8125"
statsd_client_flush_interval: "10ms"
uuid: "fake-uuid"
system_domain: "example.com"
router_groups:
- name: router-group-1
  reservable_ports: 1200
  type: tcp
- name: ` + name + `
  reservable_ports: 10000-42000
  type: tcp`
			}

			Context("with different names", func() {
				It("should not error", func() {
					config := testConfig("router-group-2")
					err := cfg.Initialize([]byte(config), true)
					Expect(err).NotTo(HaveOccurred())
					expectedGroups = models.RouterGroups{
						{
							Name:            "router-group-1",
							ReservablePorts: "1200",
							Type:            "tcp",
						},
						{
							Name:            "router-group-2",
							ReservablePorts: "10000-42000",
							Type:            "tcp",
						},
					}
					Expect(cfg.RouterGroups).To(Equal(expectedGroups))
				})
			})
		})

		Context("when router groups are seeded in the configuration file", func() {
			var expectedGroups models.RouterGroups

			testConfig := func(ports string) string {
				return `log_guid: "my_logs"
admin_socket: "/some/path"
metrics_reporting_interval: "500ms"
statsd_endpoint: "localhost:8125"
statsd_client_flush_interval: "10ms"
uuid: "fake-uuid"
system_domain: "example.com"
router_groups:
- name: router-group-1
  reservable_ports: ` + ports + `
  type: tcp
- name: router-group-2
  reservable_ports: 1024-10000,42000
  type: udp`
			}

			It("populates the router groups", func() {
				config := testConfig("12000")
				err := cfg.Initialize([]byte(config), true)
				Expect(err).NotTo(HaveOccurred())
				expectedGroups = models.RouterGroups{
					{
						Name:            "router-group-1",
						ReservablePorts: "12000",
						Type:            "tcp",
					},
					{
						Name:            "router-group-2",
						ReservablePorts: "1024-10000,42000",
						Type:            "udp",
					},
				}
				Expect(cfg.RouterGroups).To(Equal(expectedGroups))
			})

			It("returns error for invalid ports", func() {
				config := testConfig("abc")
				err := cfg.Initialize([]byte(config), true)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Port must be between 1024 and 65535"))
			})

			It("does not returns error for ports prefixed with zero", func() {
				config := testConfig("00003202-4000")
				err := cfg.Initialize([]byte(config), true)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns error for invalid port", func() {
				config := testConfig("70000")
				err := cfg.Initialize([]byte(config), true)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Port must be between 1024 and 65535"))
			})

			It("returns error for invalid ranges of ports", func() {
				config := testConfig("1024-65535,10000-20000")
				err := cfg.Initialize([]byte(config), true)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Overlapping values: [1024-65535] and [10000-20000]"))
			})

			It("returns error for invalid range of ports", func() {
				config := testConfig("1023-65530")
				err := cfg.Initialize([]byte(config), true)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Port must be between 1024 and 65535"))
			})

			It("returns error for invalid start range", func() {
				config := testConfig("1024-65535,-10000")
				err := cfg.Initialize([]byte(config), true)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("range (-10000) requires a starting port"))
			})

			It("returns error for invalid end range", func() {
				config := testConfig("10000-")
				err := cfg.Initialize([]byte(config), true)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("range (10000-) requires an ending port"))
			})

			It("returns error for invalid router group type", func() {
				missingType := `log_guid: "my_logs"
admin_socket: "/some/path"
metrics_reporting_interval: "500ms"
statsd_endpoint: "localhost:8125"
statsd_client_flush_interval: "10ms"
uuid: "fake-uuid"
system_domain: "example.com"
router_groups:
- name: router-group-1
  reservable_ports: 1024-65535`
				err := cfg.Initialize([]byte(missingType), true)
				Expect(err).To(HaveOccurred())
			})

			It("returns error for invalid router group type", func() {
				missingName := `log_guid: "my_logs"
admin_socket: "/some/path"
metrics_reporting_interval: "500ms"
statsd_endpoint: "localhost:8125"
statsd_client_flush_interval: "10ms"
uuid: "fake-uuid"
system_domain: "example.com"
router_groups:
- type: tcp
  reservable_ports: 1024-65535`
				err := cfg.Initialize([]byte(missingName), true)
				Expect(err).To(HaveOccurred())
			})

			It("returns error for missing reservable port range", func() {
				missingRouterGroup := `log_guid: "my_logs"
admin_socket: "/some/path"
metrics_reporting_interval: "500ms"
statsd_endpoint: "localhost:8125"
statsd_client_flush_interval: "10ms"
uuid: "fake-uuid"
system_domain: "example.com"
router_groups:
- type: tcp
  name: default-tcp`
				err := cfg.Initialize([]byte(missingRouterGroup), true)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Missing reservable_ports in router group:"))
			})
		})

		Context("when there are errors in the yml file", func() {
			var test_config string
			It("errors if there is no system_domain", func() {
				test_config = `log_guid: "my_logs"
admin_socket: "/some/path"
debug_address: "1.2.3.4:1234"
metron_config:
  address: "1.2.3.4"
  port: "4567"
metrics_reporting_interval: "500ms"
uuid: "fake-uuid"
statsd_endpoint: "localhost:8125"
statsd_client_flush_interval: "10ms"`
				err := cfg.Initialize([]byte(test_config), true)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("system_domain"))
			})

			Context("UAA errors", func() {
				BeforeEach(func() {
					test_config = `log_guid: "my_logs"
admin_socket: "/some/path"
debug_address: "1.2.3.4:1234"
system_domain: "example.com"
uuid: "fake-uuid"
metron_config:
  address: "1.2.3.4"
  port: "4567"
metrics_reporting_interval: "500ms"
statsd_endpoint: "localhost:8125"
statsd_client_flush_interval: "10ms"`
				})

				Context("when auth is enabled", func() {
					It("errors if no token endpoint url is found", func() {
						err := cfg.Initialize([]byte(test_config), false)
						Expect(err).To(HaveOccurred())
					})
				})

				Context("when auth is disabled", func() {
					It("it return valid config", func() {
						err := cfg.Initialize([]byte(test_config), true)
						Expect(err).NotTo(HaveOccurred())
					})
				})
			})
		})
		Context("when there are no router groups seeded in the configuration file", func() {

			testConfig := `log_guid: "my_logs"
admin_socket: "/some/path"
system_domain: "example.com"
metrics_reporting_interval: "500ms"
uuid: "fake-uuid"
statsd_endpoint: "localhost:8125"
statsd_client_flush_interval: "10ms"`

			It("does not populates the router group", func() {
				err := cfg.Initialize([]byte(testConfig), true)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.RouterGroups).To(BeNil())
			})

		})
	})
})
