package updater

import (
	"sync"

	"github.com/pkg/errors"
)

type Updater struct {
	mu sync.Mutex
	m  map[MeansType]*Means
}

func New() *Updater {
	return &Updater{}
}

func (u *Updater) RegisterMeans(m *Means) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if _, ok := u.m[m.typ]; !ok {
		return errors.Errorf("duplicated means: %s", m.typ)
	}
	u.m[m.typ] = m
	return nil
}

func (u *Updater) UpdateBy(typ MeansType) error {
	m, ok := u.m[typ]
	if !ok {
		return errors.Errorf("no such means: %s", typ)
	}
	return m.Update()
}
