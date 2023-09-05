package util

import (
	"os"
	"path"
	"testing"

	runnerErrors "github.com/cloudbase/garm-provider-common/errors"
	"github.com/cloudbase/garm-provider-common/params"
	"github.com/stretchr/testify/require"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func TestResolveToGithubArch(t *testing.T) {
	ghArch, err := ResolveToGithubArch("amd64")
	require.NoError(t, err)
	require.Equal(t, "x64", ghArch)
}

func TestResolveToGithubArchUnknown(t *testing.T) {
	arch := "some-unknown-arch"

	_, err := ResolveToGithubArch(arch)
	require.Error(t, err)
	require.EqualError(t, runnerErrors.NewNotFoundError("arch %s is unknown", arch), err.Error())
}

func TestResolveToGithubOSType(t *testing.T) {
	ghOSType, err := ResolveToGithubOSType("linux")
	require.NoError(t, err)
	require.Equal(t, "linux", ghOSType)
}

func TestResolveToGithubOSTypeUnknown(t *testing.T) {
	osType := "some-unknown-os"

	_, err := ResolveToGithubOSType(osType)
	require.Error(t, err)
	require.EqualError(t, runnerErrors.NewNotFoundError("os %s is unknown", osType), err.Error())
}

func TestResolveToGithubTag(t *testing.T) {
	ghOSTag, err := ResolveToGithubTag("linux")
	require.NoError(t, err)
	require.Equal(t, "Linux", ghOSTag)
}

func TestResolveToGithubTagUnknown(t *testing.T) {
	osTag := params.OSType("some-unknown-os")

	_, err := ResolveToGithubTag(osTag)
	require.Error(t, err)
	require.EqualError(t, runnerErrors.NewNotFoundError("os %s is unknown", osTag), err.Error())
}

func TestIsValidEmail(t *testing.T) {
	validEmail := "test@example.com"

	isValid := IsValidEmail(validEmail)
	require.True(t, isValid)
	require.Equal(t, true, isValid)
}

func TestIsInvalidEmail(t *testing.T) {
	validEmail := "invalid-email"

	isValid := IsValidEmail(validEmail)
	require.False(t, isValid)
	require.Equal(t, false, isValid)
}

func TestIsAlphanumeric(t *testing.T) {
	validString := "test123"

	isValid := IsAlphanumeric(validString)
	require.True(t, isValid)
	require.Equal(t, true, isValid)
}

func TestIsAlphanumericInvalid(t *testing.T) {
	validString := "test@123"

	isValid := IsAlphanumeric(validString)
	require.False(t, isValid)
	require.Equal(t, false, isValid)
}

func TestGetLoggingWriterEmptyLogFile(t *testing.T) {
	writer, err := GetLoggingWriter("")
	require.NoError(t, err)
	require.Equal(t, os.Stdout, writer)
}

func TestGetLoggingWriterValidLogFile(t *testing.T) {
	// Create a temporary directory for log files.
	dir, err := os.MkdirTemp("", "garm-log")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %s", err)
	}
	// Remove the temporary directory when the test finishes.
	t.Cleanup(func() { os.RemoveAll(dir) })

	// Create a log file path within the temporary directory.
	logFile := path.Join(dir, "test.log")

	writer, err := GetLoggingWriter(logFile)
	require.NoError(t, err)

	// Assert that the writer is of type *lumberjack.Logger and it has the
	// expected settings.
	logger, ok := writer.(*lumberjack.Logger)
	require.True(t, ok)
	require.Equal(t, logFile, logger.Filename)
	require.Equal(t, 500, logger.MaxSize)
	require.Equal(t, 3, logger.MaxBackups)
	require.Equal(t, 28, logger.MaxAge)
	require.Equal(t, true, logger.Compress)
}

func TestGetLoggingWriterFailedToCreateLogFolder(t *testing.T) {
	// Add a log file path that includes a directory that does not exist.
	logFile := "/non-existent-system-dir/test.log"

	_, err := GetLoggingWriter(logFile)
	require.Error(t, err)
	require.EqualError(t, err, "failed to create log folder")
}

func TestGetLoggingWriterPermisionDenied(t *testing.T) {
	// Create a temporary directory for testing
	dir, err := os.MkdirTemp("", "test-dir")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %s", err)
	}
	// Remove the temporary directory when the test finishes.
	t.Cleanup(func() { os.RemoveAll(dir) })

	// Remove execute permission from the temporary directory
	err = os.Chmod(dir, 0644)
	if err != nil {
		t.Fatalf("failed to remove execute permission from temporary directory: %s", err)
	}

	_, err = GetLoggingWriter(path.Join(dir, "non-existing-folder", "test.log"))
	require.Error(t, err)
	require.EqualError(t, err, "failed to create log folder")
}
