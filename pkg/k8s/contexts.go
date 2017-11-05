package k8s

// Contexts returns a list of contexts available based on the kubernetes config
func Contexts() []string {
	contexts := make([]string, 0)
	for k := range config.Contexts {
		contexts = append(contexts, k)
	}
	return contexts
}
