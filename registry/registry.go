package registry

import (
	"fmt"

	"github.com/YashShekhawat/fusion/drivers"
)

type Registry struct {
	drivers map[string]drivers.Drivers
}

func New() *Registry {
	return &Registry{
		drivers: make(map[string]drivers.Drivers),
	}
}

func (r *Registry) Register(d drivers.Drivers) error {
	if d == nil {
		return fmt.Errorf("driver cannot be nil")
	}

	if _, exists := r.drivers[d.Name()]; exists {
		return fmt.Errorf("driver with name %s already exists", d.Name())
	}
	r.drivers[d.Name()] = d
	return nil
}

func (r *Registry) Get(name string) (drivers.Drivers, error) {
	fmt.Println("Registry: Searching for driver:", name)
	if d, exists := r.drivers[name]; exists {
		return d, nil
	}
	return nil, fmt.Errorf("driver with name %s not found", name)
}

func (r *Registry) Remove(name string) error {
	if _, exists := r.drivers[name]; !exists {
		return fmt.Errorf("Driver with name %s not found", name)
	}
	delete(r.drivers, name)
	return nil
}

func (r *Registry) List() []drivers.Drivers {
	drivers := make([]drivers.Drivers, 0, len(r.drivers))
	for _, d := range r.drivers {
		drivers = append(drivers, d)
	}
	return drivers
}
