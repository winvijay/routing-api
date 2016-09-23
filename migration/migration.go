package migration

type Migration interface {
	RunMigration() error
	Version() int
}

func RunMigrations() {
	migrations := []Migration{
		new(V1EtcdMigration),
		new(V0InitMigration),
	}

	for _, m := range migrations {
		m.RunMigration()
	}
}
