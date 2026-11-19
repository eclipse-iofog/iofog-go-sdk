package apps

const DefaultAPIVersion = "iofog.org/v3"

type deployOptions struct {
	apiVersion string
}

// DeployOption configures deploy YAML emission.
type DeployOption func(*deployOptions)

// WithAPIVersion sets the apiVersion field emitted in deploy YAML.
func WithAPIVersion(v string) DeployOption {
	return func(o *deployOptions) {
		o.apiVersion = v
	}
}

func resolveDeployOptions(opts ...DeployOption) deployOptions {
	o := deployOptions{apiVersion: DefaultAPIVersion}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}
