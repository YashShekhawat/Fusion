package registry

import (
	"errors"
	"testing"

	"github.com/YashShekhawat/fusion/drivers/mock"
	fusionerrors "github.com/YashShekhawat/fusion/fusionerrors"
)

func TestRegisterAndGetUseFusionErrors(t *testing.T) {
	r := New()
	m := mock.New()

	if err := r.Register(m); err != nil {
		t.Fatalf("expected registration to succeed, got %v", err)
	}

	if err := r.Register(m); !errors.Is(err, fusionerrors.ErrDuplicateDriver) {
		t.Fatalf("expected duplicate driver error, got %v", err)
	}

	if _, err := r.Get(m.Name()); err != nil {
		t.Fatalf("expected driver lookup to succeed, got %v", err)
	}

	if _, err := r.Get("missing"); !errors.Is(err, fusionerrors.ErrDriverNotFound) {
		t.Fatalf("expected driver not found error, got %v", err)
	}
}
