// Copyright (c) 2022 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package fx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestDecorateSuccess(t *testing.T) {
	type Logger struct {
		Name string
	}

	t.Run("decorate something from Module", func(t *testing.T) {
		redis := fx.Module("redis",
			fx.Provide(func() *Logger {
				return &Logger{Name: "redis"}
			}),
		)

		testRedis := fx.Module("testRedis",
			redis,
			fx.Decorate(func() *Logger {
				return &Logger{Name: "testRedis"}
			}),
			fx.Invoke(func(l *Logger) {
				assert.Equal(t, "testRedis", l.Name)
			}),
		)

		app := fxtest.New(t,
			testRedis,
			fx.Invoke(func(l *Logger) {
				assert.Equal(t, "redis", l.Name)
			}),
		)
		defer app.RequireStart().RequireStop()
	})

	t.Run("decorate a dependency from root", func(t *testing.T) {
		redis := fx.Module("redis",
			fx.Decorate(func() *Logger {
				return &Logger{Name: "redis"}
			}),
			fx.Invoke(func(l *Logger) {
				assert.Equal(t, "redis", l.Name)
			}),
		)
		app := fxtest.New(t,
			redis,
			fx.Provide(func() *Logger {
				assert.Fail(t, "should not run this")
				return &Logger{Name: "root"}
			}),
		)
		defer app.RequireStart().RequireStop()
	})

	t.Run("use a decorator in root", func(t *testing.T) {
		redis := fx.Module("redis",
			fx.Invoke(func(l *Logger) {
				assert.Equal(t, "decorated logger", l.Name)
			}),
		)
		logger := fx.Module("logger",
			fx.Provide(func() *Logger {
				return &Logger{Name: "logger"}
			}),
		)
		app := fxtest.New(t,
			redis,
			logger,
			fx.Decorate(func(l *Logger) *Logger {
				return &Logger{Name: "decorated" + l.Name}
			}),
		)
		defer app.RequireStart().RequireStop()
	})
}
