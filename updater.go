package updater

import (
	"context"
	"sync"

	"github.com/pkg/errors"
)

type Updater struct {
	mu       sync.Mutex
	m        map[MeansType]Means
	UpdateIf UpdateCondition
}

func New() *Updater {
	return &Updater{
		m: map[MeansType]Means{},
	}
}

func (u *Updater) RegisterMeans(m Means) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if _, ok := u.m[m.Type()]; ok {
		return errors.Errorf("duplicated means: %s", m.Type())
	}
	u.m[m.Type()] = m
	return nil
}

func (u *Updater) UpdateBy(typ MeansType) error {
	m, ok := u.m[typ]
	if !ok {
		return errors.Errorf("no such means: %s", typ)
	}
	_, err := m.Update(context.TODO())
	return err
}
