package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("run", func(t *testing.T) {
		environment, err := ReadDir("./testdata/env")

		expectedMap := make(Environment)
		expectedMap["BAR"] = EnvValue{Value: "bar", NeedRemove: false}
		expectedMap["EMPTY"] = EnvValue{Value: "", NeedRemove: true}
		expectedMap["FOO"] = EnvValue{Value: "   foo\nwith new line", NeedRemove: false}
		expectedMap["HELLO"] = EnvValue{Value: "\"hello\"", NeedRemove: false}
		expectedMap["UNSET"] = EnvValue{Value: "", NeedRemove: true}

		require.Nil(t, err)
		require.NotNil(t, environment)

		for key, expectedValue := range expectedMap {
			value, ok := environment[key]

			require.True(t, ok)
			require.Equal(t, expectedValue, value)
		}
	})
}
