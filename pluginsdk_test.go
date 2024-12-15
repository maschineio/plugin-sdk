package plugin_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"maschine.io/core/context"
	pluginsdk "maschine.io/plugin-sdk"
)

func TestGetResourceManager(t *testing.T) {
	rm1 := pluginsdk.GetResourceManager()
	rm2 := pluginsdk.GetResourceManager()
	assert.Equal(t, rm1, rm2, "GetResourceManager should return the same instance")
}

func TestRegisterLambdaFn(t *testing.T) {
	rm := pluginsdk.GetResourceManager()
	fn := func(ctx *context.Context) (any, error) {
		return "test", nil
	}

	err := rm.RegisterLambdaFn("testFn", fn)
	assert.NoError(t, err, "RegisterLambdaFn should not return an error")

	err = rm.RegisterLambdaFn("testFn", fn)
	assert.Error(t, err, "RegisterLambdaFn should return an error for duplicate registration")
}

// func TestGetFn(t *testing.T) {
// 	rm := pluginsdk.GetResourceManager()
// 	fn := func(ctx *context.Context) (any, error) {
// 		return "test", nil
// 	}
// 	rm.RegisterLambdaFn("testFn", fn)
//
// 	retrievedFn := rm.GetFn("testFn")
// 	assert.NotNil(t, retrievedFn, "GetFn should return the registered function")
// 	assert.Equal(t, fn, retrievedFn, "GetFn should return the correct function")
//
// 	nilFn := rm.GetFn("nonExistentFn")
// 	assert.Nil(t, nilFn, "GetFn should return nil for non-existent function")
// }

func TestResourceNames(t *testing.T) {
	rm := pluginsdk.GetResourceManager()
	rm.RegisterLambdaFn("testFn1", func(ctx *context.Context) (any, error) { return "test1", nil })
	rm.RegisterLambdaFn("testFn2", func(ctx *context.Context) (any, error) { return "test2", nil })

	names := rm.ResourceNames()
	assert.ElementsMatch(t, names, []string{"testFn1", "testFn2"}, "ResourceNames should return the correct list of resource names")
}

// func TestLoadPlugins(t *testing.T) {
// 	// Mock plugin loading by creating temporary .so files
// 	tempDir := os.TempDir()
// 	pluginFile := filepath.Join(tempDir, "mockplugin.so")
// 	os.Create(pluginFile)
// 	defer os.Remove(pluginFile)
//
// 	rm := pluginsdk.GetResourceManager()
// 	err := rm.LoadPlugins(tempDir)
// 	assert.NoError(t, err, "LoadPlugins should not return an error")
// }
