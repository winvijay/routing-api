package migration_test

import (
	"code.cloudfoundry.org/routing-api/migration"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Version", func() {
	var v migration.Version
	Context("when current version is less than target version", func() {
		BeforeEach(func() {
			v.CurrentVersion = 1
			v.TargetVersion = 2
		})
		It("should indicate intent to perform migration", func() {
			Expect(v.NeedsMigration()).To(BeTrue())
		})
	})
	Context("when current version is greater than target version", func() {
		BeforeEach(func() {
			v.CurrentVersion = 3
			v.TargetVersion = 2
		})
		It("should not indicate intent to perform migration", func() {
			Expect(v.NeedsMigration()).To(BeFalse())
		})
	})
	Context("when current version is equal to target version", func() {
		BeforeEach(func() {
			v.CurrentVersion = 1
			v.TargetVersion = 1
		})
		It("should not indicate intent to perform migration", func() {
			Expect(v.NeedsMigration()).To(BeFalse())
		})
	})
})
