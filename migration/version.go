package migration

type Version struct {
	CurrentVersion int
	TargetVersion  int
}

func (v *Version) NeedsMigration() bool {
	return v.CurrentVersion < v.TargetVersion
}
