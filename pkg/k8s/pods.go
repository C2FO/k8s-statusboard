package k8s

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

// Pods returns the pods for a given context
func Pods(context string) ([]v1.Pod, error) {
	cs, err := getClientSet(context)
	if err != nil {
		return nil, err
	}

	ctxConfig, ok := config.Contexts[context]
	if !ok {
		return nil, fmt.Errorf("no context named %s", context)
	}

	resp, err := cs.Pods(ctxConfig.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}
