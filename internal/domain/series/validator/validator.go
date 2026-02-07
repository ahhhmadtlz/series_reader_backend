package validator

type Repository interface {
}

type Validator struct {
	repo Repository
}

func New(repo Repository) Validator {
	return Validator{
		repo: repo,
	}
}