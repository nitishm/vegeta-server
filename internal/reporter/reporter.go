package reporter

type IReporter interface {
}

type reporter struct {
}

func NewReporter() *reporter {
	return &reporter{}
}
