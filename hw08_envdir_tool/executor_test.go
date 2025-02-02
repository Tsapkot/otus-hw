package main

import (
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("simple test", func(t *testing.T) {
		err := os.Setenv("FAKEVAR", "lol")
		if err != nil {
			t.Errorf("error setting variable: %v", err)
		}

		var shell, flag string
		if runtime.GOOS == "windows" {
			shell = "cmd"
			flag = "/c"
		} else {
			shell = "sh"
			flag = "-c"
		}

		returnCode := RunCmd([]string{shell, flag, ""}, Environment{
			"FAKEVAR": EnvValue{"kek", false},
		})

		require.Equal(t, os.Getenv("FAKEVAR"), "kek")
		require.Equal(t, 0, returnCode)
	})
}
