package cloudconfig

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/cloudbase/garm-provider-common/defaults"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// helper function
func getCloudInit() *CloudInit {
	return &CloudInit{
		PackageUpgrade: true,
		SystemInfo: &SystemInfo{
			DefaultUser: DefaultUser{
				Name:   defaults.DefaultUser,
				Home:   fmt.Sprintf("/home/%s", defaults.DefaultUser),
				Shell:  defaults.DefaultUserShell,
				Groups: defaults.DefaultUserGroups,
				Sudo:   "ALL=(ALL) NOPASSWD:ALL",
			},
		},
		WriteFiles: []File{
			{
				Encoding:    "b64",
				Content:     base64.StdEncoding.EncodeToString([]byte("content")),
				Owner:       "owner",
				Path:        "path",
				Permissions: "permissions",
			},
		},
	}
}

func TestNewDefaultCloudInitConfig(t *testing.T) {
	cloudInit := getCloudInit()

	initCfg := NewDefaultCloudInitConfig()
	require.Equal(t, cloudInit.PackageUpgrade, initCfg.PackageUpgrade)
	require.Equal(t, cloudInit.SystemInfo.DefaultUser.Name, initCfg.SystemInfo.DefaultUser.Name)
	require.Equal(t, cloudInit.SystemInfo.DefaultUser.Home, initCfg.SystemInfo.DefaultUser.Home)
	require.Equal(t, cloudInit.SystemInfo.DefaultUser.Shell, initCfg.SystemInfo.DefaultUser.Shell)
	require.Equal(t, cloudInit.SystemInfo.DefaultUser.Groups, initCfg.SystemInfo.DefaultUser.Groups)
	require.Equal(t, cloudInit.SystemInfo.DefaultUser.Sudo, initCfg.SystemInfo.DefaultUser.Sudo)
}

func TestAddCACertNil(t *testing.T) {
	cloudInit := getCloudInit()

	err := cloudInit.AddCACert(nil)
	require.NoError(t, err)
	require.NoError(t, cloudInit.AddCACert(nil))
}

func TestAddCACertFailed(t *testing.T) {
	cloudInit := getCloudInit()

	err := cloudInit.AddCACert([]byte("unknown-cert"))
	require.EqualError(t, err, "failed to parse CA cert bundle")
}

func TestAddSSHKey(t *testing.T) {
	cloudInit := getCloudInit()

	cloudInit.AddSSHKey("ssh-key")
	require.Equal(t, "ssh-key", cloudInit.SSHAuthorizedKeys[0])
}

func TestAddSSHKeyNotFound(t *testing.T) {
	cloudInit := getCloudInit()
	cloudInit.SSHAuthorizedKeys = []string{""}

	cloudInit.AddSSHKey("")
	require.Equal(t, []string{""}, cloudInit.SSHAuthorizedKeys)
}

func TestAddPackage(t *testing.T) {
	cloudInit := getCloudInit()

	cloudInit.AddPackage("curl", "wget")
	require.Equal(t, "curl", cloudInit.Packages[0])
	require.Equal(t, len(cloudInit.Packages), 2)
}

func TestAddPackageNotFound(t *testing.T) {
	cloudInit := getCloudInit()
	cloudInit.Packages = []string{""}

	cloudInit.AddPackage("")
	require.Equal(t, []string{""}, cloudInit.Packages)
}

func TestAddRunCmd(t *testing.T) {
	cloudInit := getCloudInit()

	cloudInit.AddRunCmd("cmd")
	require.Equal(t, "cmd", cloudInit.RunCmd[0])
}

func TestAddFile(t *testing.T) {
	cloudInit := getCloudInit()

	cloudInit.AddFile([]byte("content"), "test-path", "test-owner", "test-permissions")
	require.Equal(t, "b64", cloudInit.WriteFiles[1].Encoding)
	require.Equal(t, "Y29udGVudA==", cloudInit.WriteFiles[1].Content)
	require.Equal(t, "test-owner", cloudInit.WriteFiles[1].Owner)
	require.Equal(t, "test-permissions", cloudInit.WriteFiles[1].Permissions)
	require.Equal(t, "test-path", cloudInit.WriteFiles[1].Path)
}

func TestSerialize(t *testing.T) {
	cloudInit := getCloudInit()

	asYaml, err := yaml.Marshal(cloudInit)
	if err != nil {
		t.Errorf("Failed to marshal cloudInit: %v", err)
	}

	serialized, err := cloudInit.Serialize()
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("%s\n%s", "#cloud-config", string(asYaml)), serialized)
}
