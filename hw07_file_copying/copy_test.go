package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const output = "./testdata/output.txt"

func cleanup() {
	err := os.Remove(output)
	if err != nil {
		return
	}
}

func TestCopy0_0(t *testing.T) {
	defer cleanup()
	err := Copy("testdata/input.txt", output, 0, 0)
	fact, _ := os.ReadFile(output)
	expected, _ := os.ReadFile("testdata/out_offset0_limit0.txt")
	require.Equal(t, expected, fact)
	require.Nil(t, err)
}

func TestCopy0_10(t *testing.T) {
	defer cleanup()
	err := Copy("testdata/input.txt", output, 0, 10)
	fact, _ := os.ReadFile(output)
	expected, _ := os.ReadFile("testdata/out_offset0_limit10.txt")
	require.Equal(t, expected, fact)
	require.Nil(t, err)
}

func TestCopy0_1000(t *testing.T) {
	defer cleanup()
	err := Copy("testdata/input.txt", output, 0, 1000)
	fact, _ := os.ReadFile(output)
	expected, _ := os.ReadFile("testdata/out_offset0_limit1000.txt")
	require.Equal(t, expected, fact)
	require.Nil(t, err)
}

func TestCopy0_10000(t *testing.T) {
	defer cleanup()
	err := Copy("testdata/input.txt", output, 0, 10000)
	fact, _ := os.ReadFile(output)
	expected, _ := os.ReadFile("testdata/out_offset0_limit10000.txt")
	require.Equal(t, expected, fact)
	require.Nil(t, err)
}

func TestCopy100_1000(t *testing.T) {
	defer cleanup()
	err := Copy("testdata/input.txt", output, 100, 1000)
	fact, _ := os.ReadFile(output)
	expected, _ := os.ReadFile("testdata/out_offset100_limit1000.txt")
	require.Equal(t, expected, fact)
	require.Nil(t, err)
}

func TestCopy6000_1000(t *testing.T) {
	defer cleanup()
	err := Copy("testdata/input.txt", output, 6000, 1000)
	expected, _ := os.ReadFile("testdata/out_offset6000_limit1000.txt")
	fact, _ := os.ReadFile(output)
	require.Equal(t, expected, fact)
	require.Nil(t, err)
}

func TestCopyOverlap(t *testing.T) {
	defer cleanup()
	err := Copy("testdata/input.txt", "testdata/input.txt", 0, 0)
	require.ErrorIs(t, err, ErrFileOverlap)
}
