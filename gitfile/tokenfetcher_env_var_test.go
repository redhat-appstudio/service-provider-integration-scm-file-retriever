package gitfile

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCreateHeaderStructFromEnv(t *testing.T) {
	//t.Setenv("TOKEN", "abcd_foo") // only since 1.17
	os.Setenv("TOKEN", "abcd_foo")
	defer os.Unsetenv("TOKEN")
	headerStruct := new(EnvVarTokenFetcher).BuildHeader("https://github.com/any/test.git")
	assert.Equal(t, "Bearer abcd_foo", headerStruct.Authorization, "Authorization header value mismatch")
}
