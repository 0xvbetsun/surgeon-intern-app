//go:generate mockery --name IAdder
package organisations

import "github.com/vbetsun/surgeon-intern-app/pkg/repo/organisations"

type (
	IAdder interface {
		Add(organisation *APIOrganisation) (*APIOrganisation, error)
	}
)

type Adder struct {
	repo organisations.IAdder
}

func NewAdder(repo organisations.IAdder) IAdder {
	return &Adder{repo: repo}
}

func (a *Adder) Add(organisation *APIOrganisation) (*APIOrganisation, error) {
	panic("implement me")
}
