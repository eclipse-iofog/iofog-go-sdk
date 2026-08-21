package apps

import "testing"

func TestResolveDeployOptions_defaultAPIVersion(t *testing.T) {
	o := resolveDeployOptions()
	if o.apiVersion != DefaultAPIVersion {
		t.Errorf("apiVersion = %q, want %q", o.apiVersion, DefaultAPIVersion)
	}
}

func TestResolveDeployOptions_withAPIVersion(t *testing.T) {
	datasance := "datasance.com" + "/v3"
	o := resolveDeployOptions(WithAPIVersion(datasance))
	if o.apiVersion != datasance {
		t.Errorf("apiVersion = %q, want %q", o.apiVersion, datasance)
	}
}

func TestApplicationExecutor_apiVersionFromOptions(t *testing.T) {
	datasance := "datasance.com" + "/v3"
	exe := newApplicationExecutor(IofogController{}, nil, "app", WithAPIVersion(datasance))
	if exe.apiVersion != datasance {
		t.Errorf("apiVersion = %q, want %q", exe.apiVersion, datasance)
	}
}

func TestMicroserviceExecutor_apiVersionDefault(t *testing.T) {
	exe := newMicroserviceExecutor(IofogController{}, nil, "app", "msvc")
	if exe.apiVersion != DefaultAPIVersion {
		t.Errorf("apiVersion = %q, want %q", exe.apiVersion, DefaultAPIVersion)
	}
}

func TestApplicationTemplateExecutor_apiVersionFromOptions(t *testing.T) {
	datasance := "datasance.com" + "/v3"
	exe := newApplicationTemplateExecutor(IofogController{}, nil, nil, "tpl", WithAPIVersion(datasance))
	if exe.apiVersion != datasance {
		t.Errorf("apiVersion = %q, want %q", exe.apiVersion, datasance)
	}
}
