package main

import (
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestReadDir(t *testing.T) {
	cases := []struct {
		name string
		dir  string
		res  Environment
	}{
		{
			name: "envdir",
			dir:  "testdata/env",
			res: Environment{
				"BAR": EnvValue{
					Value:      "bar",
					NeedRemove: false,
				},
				"EMPTY": EnvValue{
					NeedRemove: false,
				},
				"FOO": EnvValue{
					Value:      "   foo\u0000with new line",
					NeedRemove: false,
				},
				"HELLO": EnvValue{
					Value:      "\"hello\"",
					NeedRemove: false,
				},
				"UNSET": EnvValue{
					NeedRemove: true,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			env, err := ReadDir(c.dir)
			require.NoError(t, err)
			eq := reflect.DeepEqual(c.res, env)
			require.True(t, eq)
		})
	}
}
