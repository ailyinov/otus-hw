package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRunCmd(t *testing.T) {
	cases := []struct {
		name string
		cmd  []string
		env  Environment
	}{
		{
			name: "echo",
			cmd: []string{
				"/bin/sh",
				"-c",
				"echo $FOO",
			},
			env: Environment{
				"FOO": EnvValue{
					Value:      "OUT >>>--- foooo-o-o-o-o",
					NeedRemove: false,
				},
			},
		}, {
			name: "echo $HOME",
			cmd: []string{
				"/bin/sh",
				"-c",
				"echo $HOME",
			},
			env: Environment{
				"HOME": EnvValue{
					NeedRemove: true,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			code := RunCmd(c.cmd, c.env)
			require.Equal(t, 0, code)
		})
	}

	t.Run("testdata", func(t *testing.T) {
		env, err := ReadDir("testdata/env")
		require.NoError(t, err)
		str := []string{"/bin/sh", "-c", "echo $PATH"}
		RunCmd(str, env)
	})
}
