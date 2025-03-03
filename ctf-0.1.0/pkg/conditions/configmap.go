package conditions

import (
	"context"
	"fmt"
	"os"

	"github.com/k8s-2025-pschoeppner/ctf/pkg/flags"
	"github.com/k8s-2025-pschoeppner/ctf/pkg/k8s"
	"github.com/k8s-2025-pschoeppner/ctf/pkg/types"
	"k8s.io/client-go/kubernetes"
)

const (
	configMapArg = "configmap"
	secretArg    = "secret"
)

func ReadFromMountedConfigMap(path string) flags.Fulfiller {
	return func(ctx context.Context, r types.Request, client kubernetes.Interface) error {
		f, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("expected file to be present at %q: %w", path, err)
		}
		r.Args[configMapArg] = string(f)
		return nil
	}
}

func ReadFromExternalConfigMap(name string) flags.Fulfiller {
	return func(ctx context.Context, r types.Request, client kubernetes.Interface) error {
		cm, err := k8s.GetConfigMap(ctx, client, name, r.PodNamespace)
		if err != nil {
			return fmt.Errorf("getting configmap %s/%s: %w", r.PodNamespace, name, err)
		}
		s := ""
		for k, v := range cm.Data {
			s += fmt.Sprintf("%s=%s\n", k, v)
		}
		r.Args[configMapArg] = s
		return nil
	}
}

func ReadFromMountedSecret(path string) flags.Fulfiller {
	return func(ctx context.Context, r types.Request, client kubernetes.Interface) error {
		f, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("expected file to be present at %q: %w", path, err)
		}
		r.Args[secretArg] = string(f)
		return nil
	}
}
