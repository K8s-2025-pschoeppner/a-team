package conditions

import (
	"context"
	"fmt"
	"time"

	"github.com/k8s-2025-pschoeppner/ctf/pkg/flags"
	"github.com/k8s-2025-pschoeppner/ctf/pkg/k8s"
	"github.com/k8s-2025-pschoeppner/ctf/pkg/types"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/ptr"
)

type PodValidator func(pod corev1.Pod, r types.Request) error

func PodValidators(v ...PodValidator) flags.Validator {
	return func(ctx context.Context, r types.Request, client kubernetes.Interface) error {
		timeout, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		pod, err := k8s.GetPod(timeout, client, r.PodName, r.PodNamespace)
		if err != nil {
			return fmt.Errorf("getting pod: %w", err)
		}
		for _, validator := range v {
			if err := validator(*pod, r); err != nil {
				return fmt.Errorf("validating pod: %w", err)
			}
		}
		return nil
	}
}

func WithConfigmap(volumeName, configmapName string) PodValidator {
	return func(pod corev1.Pod, r types.Request) error {
		if val, found := r.Args[configMapArg]; !found || val == "" {
			return fmt.Errorf("got empty value from configmap")
		}
		for _, volume := range pod.Spec.Volumes {
			if volume.Name == volumeName {
				if volume.ConfigMap == nil {
					return fmt.Errorf("volume %q is not a configmap", volumeName)
				}
				if volume.ConfigMap.Name != configmapName {
					return fmt.Errorf("volume %q isn't mounting configmap %q", volumeName, configmapName)
				}
				return nil
			}
		}
		return fmt.Errorf("volume %q not found", volumeName)
	}
}

func WithSecret(volumeName, secretName string) PodValidator {
	return func(pod corev1.Pod, r types.Request) error {
		if val, found := r.Args[secretArg]; !found || val == "" {
			return fmt.Errorf("got empty value from secret")
		}
		for _, volume := range pod.Spec.Volumes {
			if volume.Name == volumeName {
				if volume.Secret == nil {
					return fmt.Errorf("volume %q is not a secret", volumeName)
				}
				if volume.Secret.SecretName != secretName {
					return fmt.Errorf("volume %q isn't mounting secret %q", volumeName, secretName)
				}
				return nil
			}
		}
		return fmt.Errorf("volume %q not found", volumeName)
	}
}

func WithServiceAccount() PodValidator {
	return func(pod corev1.Pod, r types.Request) error {
		if pod.Spec.ServiceAccountName == "default" {
			return fmt.Errorf("pod is using default service account")
		}
		return nil
	}
}

func WithSecurityContext() PodValidator {
	return func(pod corev1.Pod, r types.Request) error {
		if pod.Spec.SecurityContext == nil {
			return fmt.Errorf("pod has no security context")
		}
		if pod.Spec.SecurityContext.RunAsUser != ptr.To(int64(1000)) {
			return fmt.Errorf("pod is not running as user 1000")
		}
		return nil
	}
}
