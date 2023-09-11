package cloudconfig

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInstallRunnerScript(t *testing.T) {
	script, err := InstallRunnerScript(InstallRunnerParams{}, "linux", "test-template")
	require.NoError(t, err)
	require.Equal(t, "test-template", string(script))
}

func TestInstallRunnerScriptParsingTemplateFailed(t *testing.T) {
	_, err := InstallRunnerScript(InstallRunnerParams{}, "linux", "{{ .test")
	require.Error(t, err)
	require.Contains(t, err.Error(), "parsing template: template: :1: unclosed action")
}

func TestInstallRunnerScriptRenderingTemplateFailed(t *testing.T) {
	_, err := InstallRunnerScript(InstallRunnerParams{}, "linux", "{{ .test }}")
	require.Error(t, err)
	require.Contains(t, err.Error(), "rendering template: template: :1:3: executing")
}
