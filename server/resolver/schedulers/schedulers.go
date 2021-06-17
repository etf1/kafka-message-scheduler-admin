package schedulers

type Resolver interface {
	List() ([]Scheduler, error)
}

type Scheduler interface {
	Name() string
}
