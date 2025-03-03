package flags

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/k8s-2025-pschoeppner/ctf/pkg/types"
	"k8s.io/client-go/kubernetes"
)

type (
	Validator func(context.Context, types.Request, kubernetes.Interface) error
	Fulfiller func(context.Context, types.Request, kubernetes.Interface) error
)

type Flag struct {
	Name       string
	Value      string
	Validators []Validator
	Fulfillers []Fulfiller
	Client     kubernetes.Interface
	Logger     *slog.Logger
}

func (f Flag) Success() string {
	return fmt.Sprintf("You successfully captured flag %s: %q", f.Name, f.Value)
}

type FlagOption func(*Flag)

func NewFlag(name string, client kubernetes.Interface, logger *slog.Logger, opts ...FlagOption) *Flag {
	f := &Flag{
		Name:   name,
		Value:  "",
		Client: client,
		Logger: logger,
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func WithValidators(validators ...Validator) FlagOption {
	return func(f *Flag) {
		f.Validators = validators
	}
}

func WithFulfillers(fulfillers ...Fulfiller) FlagOption {
	return func(f *Flag) {
		f.Fulfillers = fulfillers
	}
}

func (f Flag) Handler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := types.FromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = req.Validate()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for _, validator := range f.Validators {
			err := validator(ctx, req, f.Client)
			if err != nil {
				http.Error(w, fmt.Errorf("condition failed: %w", err).Error(), http.StatusForbidden)
				return
			}
		}
		_, err = w.Write([]byte(f.Success()))
		if err != nil {
			f.Logger.ErrorContext(ctx, "write success response", slog.String("client", r.RemoteAddr), slog.String("namespace", req.PodNamespace), slog.String("pod", req.PodName), slog.String("err", err.Error()))
		}
		return
	}
}

func (f *Flag) SetValue(value string) {
	f.Value = value
}
