package file

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileUtils(t *testing.T) {

	t.Run(
		"GIVEN absolute path EXPECT identical absolute path", func(t *testing.T) {
			path := "/tmp/axelar/test"
			resolved := resolve(path)
			assert.Equal(
				t,
				path,
				resolved,
			)
		},
	)

	t.Run(
		"home resolver", func(t *testing.T) {
			t.Run(
				"GIVEN ~ EXPECT path to home directory", func(t *testing.T) {
					path := "~"
					expected, err := os.UserHomeDir()
					if assert.NoError(t, err) {
						resolved := resolve(path)
						assert.Equal(
							t,
							expected,
							resolved,
						)
					}
				},
			)

			t.Run(
				"GIVEN relative path starting with ~ EXPECT absolute path", func(t *testing.T) {
					path := "~/axelar/test"
					home, err := os.UserHomeDir()
					if assert.NoError(t, err) {
						expected := fmt.Sprintf("%v/axelar/test", home)
						resolved := resolve(path)
						assert.Equal(
							t,
							expected,
							resolved,
						)
					}
				},
			)

			t.Run(
				"GIVEN absolute path containing ~ EXPECT identical path", func(t *testing.T) {
					path := "/~/axelar/test"
					resolved := resolve(path)
					assert.Equal(
						t,
						path,
						resolved,
					)
				},
			)
		},
	)

	t.Run(
		"current directory resolver", func(t *testing.T) {
			t.Run(
				"GIVEN . EXPECT current directory", func(t *testing.T) {
					path := "."
					expected, err := os.Getwd()
					if assert.NoError(t, err) {
						resolved := resolve(path)
						assert.Equal(
							t,
							expected,
							resolved,
						)
					}
				},
			)

			t.Run(
				"GIVEN relative path starting with . EXPECT path to start with current directory", func(t *testing.T) {
					path := "./tmp/test"
					wd, err := os.Getwd()
					expected := fmt.Sprintf("%v/tmp/test", wd)
					if assert.NoError(t, err) {
						resolved := resolve(path)
						assert.Equal(
							t,
							expected,
							resolved,
						)
					}
				},
			)

			t.Run(
				"GIVEN absolute path containing . EXPECT identical path", func(t *testing.T) {
					path := "/./axelar/test"
					resolved := resolve(path)
					assert.Equal(
						t,
						path,
						resolved,
					)
				},
			)
		},
	)

}
