package configuration

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitConfig_ReadsEnvironmentValues(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("MONGO_URI", "mongodb://127.0.0.1:27017")
	t.Setenv("MONGO_DATABASE", "user_management")
	t.Setenv("MONGO_USERNAME", "user_management_admin")
	t.Setenv("MONGO_PASSWORD", "secret")
	t.Setenv("JWT_SECRET", "jwt-secret")

	err := InitConfig()
	require.NoError(t, err)
	require.Equal(t, "9090", Env.PORT)
	require.Equal(t, "mongodb://127.0.0.1:27017", Env.MONGO_URI)
	require.Equal(t, "user_management", Env.MONGO_DATABASE)
	require.Equal(t, "user_management_admin", Env.MONGO_USERNAME)
	require.Equal(t, "secret", Env.MONGO_PASSWORD)
	require.Equal(t, "jwt-secret", Env.JWT_SECRET)
}

func TestInitConfig_MissingMongoURI_ReturnsError(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("MONGO_URI", "")
	t.Setenv("MONGO_DATABASE", "user_management")
	t.Setenv("JWT_SECRET", "jwt-secret")

	err := InitConfig()
	require.Error(t, err)
}

func TestInitConfig_DefaultPortWhenEmpty(t *testing.T) {
	previousPort := os.Getenv("PORT")
	t.Cleanup(func() {
		_ = os.Setenv("PORT", previousPort)
	})
	t.Setenv("PORT", "")
	t.Setenv("MONGO_URI", "mongodb://127.0.0.1:27017")
	t.Setenv("MONGO_DATABASE", "user_management")
	t.Setenv("JWT_SECRET", "jwt-secret")

	err := InitConfig()
	require.NoError(t, err)
	require.Equal(t, "8080", Env.PORT)
}
