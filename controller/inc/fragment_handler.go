package inc

func (i Inc) DeleteFragmentsByPath(path string) error {

	panic("not implemented")

}

func (i Inc) CronJobDeleteFragments() error {

	// Checks if the job is currently being done. The job should be only handled by the first server (?)

	// See implementation code in db_cron_test.go

	panic("not implemented")
}
