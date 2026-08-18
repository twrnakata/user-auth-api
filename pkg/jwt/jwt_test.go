package jwt

import (
	"testing"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"

	"github.com/twrnakata/user-auth-api/pkg/datetime"
)

func TestJWTService_CreateToken_CanBeParsedWithSameSecret(t *testing.T) {
	require.NoError(t, datetime.SetDefaultTimeZone(datetime.TimeZoneAsiaBangkok))

	jwtService, err := NewJWTService("test-secret", DefaultExpireDuration)
	require.NoError(t, err)

	token, err := jwtService.CreateToken("u-1", "Alice")
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := jwtService.ParseToken(token)
	require.NoError(t, err)
	require.Equal(t, "u-1", claims.UserID)
	require.Equal(t, "Alice", claims.Name)
	require.Equal(t, "u-1", claims.Subject)
	require.NotNil(t, claims.ExpiresAt)
	require.NotNil(t, claims.IssuedAt)

	now := datetime.GetCurrentDateTimeNow()
	require.WithinDuration(t, now.Add(DefaultExpireDuration), claims.ExpiresAt.Time, 2*time.Second)
}

func TestJWTService_ParseToken_WrongSecret_ReturnsError(t *testing.T) {
	jwtService, err := NewJWTService("test-secret", DefaultExpireDuration)
	require.NoError(t, err)

	token, err := jwtService.CreateToken("u-1", "Alice")
	require.NoError(t, err)

	otherJWTService, err := NewJWTService("other-secret", DefaultExpireDuration)
	require.NoError(t, err)

	_, err = otherJWTService.ParseToken(token)
	require.Error(t, err)
}

func TestJWTService_CreateToken_UsesHS256(t *testing.T) {
	jwtService, err := NewJWTService("test-secret", DefaultExpireDuration)
	require.NoError(t, err)

	token, err := jwtService.CreateToken("u-1", "Alice")
	require.NoError(t, err)

	parsedToken, _, err := jwtlib.NewParser().ParseUnverified(token, &Claims{})
	require.NoError(t, err)
	require.Equal(t, jwtlib.SigningMethodHS256.Alg(), parsedToken.Method.Alg())
}

func TestNewJWTService_EmptySecret_ReturnsError(t *testing.T) {
	jwtService, err := NewJWTService("  ", DefaultExpireDuration)
	require.Error(t, err)
	require.Nil(t, jwtService)
}

func TestSecretFromEnvironment_UsesFallbackWhenEmpty(t *testing.T) {
	t.Setenv(EnvSecretKey, "")
	require.Equal(t, FallbackSecret, SecretFromEnvironment())
}

func TestSecretFromEnvironment_UsesEnvValue(t *testing.T) {
	t.Setenv(EnvSecretKey, "from-env")
	require.Equal(t, "from-env", SecretFromEnvironment())
}
