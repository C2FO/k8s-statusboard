package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	// Load the gcp plugin for authing against gcp clusters
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var config *api.Config

// LoadConfig should be called before calling any of the package methods. It
// loads the config into memory from disk.
func LoadConfig() error {
	var err error
	pathOptions := clientcmd.NewDefaultPathOptions()
	config, err = pathOptions.GetStartingConfig()
	return err
}

func getClientSet(context string) (*kubernetes.Clientset, error) {
	// Set the context so we target the right cluster
	clientConfig := clientcmd.NewDefaultClientConfig(
		*config,
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		},
	)
	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
