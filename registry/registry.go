package registry

import (
	"fmt"

	"github.com/YashShekhawat/fusion/drivers"
	fusionerrors "github.com/YashShekhawat/fusion/fusionerrors"
)

type Registry struct {
	drivers map[string]drivers.Driver
}

func New() *Registry {
	return &Registry{
		drivers: make(map[string]drivers.Driver),
	}
}

func (r *Registry) Register(d drivers.Driver) error {
	if d == nil {
		return fmt.Errorf("registry: register driver: %w", fusionerrors.ErrInvalidRequest)
	}

	if _, exists := r.drivers[d.Name()]; exists {
		return fmt.Errorf("registry: duplicate driver %q: %w", d.Name(), fusionerrors.ErrDuplicateDriver)
	}
	r.drivers[d.Name()] = d
	return nil
}

func (r *Registry) Get(name string) (drivers.Driver, error) {
	if d, exists := r.drivers[name]; exists {
		return d, nil
	}
	return nil, fmt.Errorf("registry: get driver %q: %w", name, fusionerrors.ErrDriverNotFound)
}

func (r *Registry) Remove(name string) error {
	if _, exists := r.drivers[name]; !exists {
		return fmt.Errorf("registry: remove driver %q: %w", name, fusionerrors.ErrDriverNotFound)
	}
	delete(r.drivers, name)
	return nil
}

func (r *Registry) List() []drivers.Driver {
	drivers := make([]drivers.Driver, 0, len(r.drivers))
	for _, d := range r.drivers {
		drivers = append(drivers, d)
	}
	return drivers
}
