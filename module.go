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

package fx

import (
	"fmt"

	"go.uber.org/dig"
	"go.uber.org/fx/fxevent"
	"go.uber.org/fx/internal/fxreflect"
)

// A container represents a set of constructors to provide
// dependencies, and a set of functions to invoke once all the
// dependencies have been initialized.
//
// This definition corresponds to the dig.Container and dig.Scope.
type container interface {
	Invoke(interface{}, ...dig.InvokeOption) error
	Provide(interface{}, ...dig.ProvideOption) error
	Decorate(interface{}, ...dig.DecorateOption) error
}

// Module is a named group of zero or more fx.Options. A Module is a
// dependency graph with limited scope, and can be used for scoping
// graph modifications (not implemented yet).
func Module(name string, opts ...Option) Option {
	mo := moduleOption{
		name:    name,
		options: opts,
	}
	return mo
}

type moduleOption struct {
	name    string
	options []Option
}

func (o moduleOption) String() string {
	return fmt.Sprintf("module %s", o.name)
}

func (o moduleOption) apply(app *App) {
	// This is the top-level module's apply.
	// This module acts as the "root" module that connects all of its
	// its submodules to the rest of the App.
	// 1. Create a new module
	// 2. Apply child Options on the new module
	// 3. Append the new module to the App's modules.
	newModule := &module{
		name: o.name,
		app:  app,
	}
	for _, opt := range o.options {
		opt.applyModule(newModule)
	}
	app.modules = append(app.modules, newModule)
}

func (o moduleOption) applyModule(mod *module) {
	// This get called on any submodules' that are declared
	// as part of another module.

	// 1. Create a new module with the parent being the specified
	// module.
	// 2. Apply child Options on the new module.
	// 3. Append it to the parent module.
	newModule := &module{
		name:   o.name,
		parent: mod,
		app:    mod.app,
	}
	for _, opt := range o.options {
		opt.applyModule(newModule)
	}
	mod.modules = append(mod.modules, newModule)
}

type module struct {
	parent     *module
	name       string
	scope      *dig.Scope
	provides   []provide
	invokes    []invoke
	decorators []decorator
	modules    []*module
	app        *App
}

// builds the Scopes using the App's Container. Note that this happens
// after applyModules' are called because the App's Container needs to
// be built for any Scopes to be initialized, and applys' should be called
// before the Container can get initialized.
func (m *module) build(app *App, root *dig.Container) {
	if m.parent == nil {
		m.scope = root.Scope(m.name)
	} else {
		parentScope := m.parent.scope
		m.scope = parentScope.Scope(m.name)
	}

	for i := 0; i < len(m.modules); i++ {
		m.modules[i].build(app, root)
	}
}

func (m *module) provideAll() {
	for _, p := range m.provides {
		var info dig.ProvideInfo
		if err := runProvide(m.scope, p, dig.FillProvideInfo(&info), dig.Export(true)); err != nil {
			m.app.err = err
		}
		var ev fxevent.Event
		switch {
		case p.IsSupply:
			ev = &fxevent.Supplied{
				TypeName: p.SupplyType.String(),
				Err:      m.app.err,
			}

		default:
			outputNames := make([]string, len(info.Outputs))
			for i, o := range info.Outputs {
				outputNames[i] = o.String()
			}

			ev = &fxevent.Provided{
				ConstructorName: fxreflect.FuncName(p.Target),
				OutputTypeNames: outputNames,
				Err:             m.app.err,
			}
		}
		m.app.log.LogEvent(ev)
	}

	for _, m := range m.modules {
		m.provideAll()
	}
}

func (m *module) executeInvokes() error {
	for _, invoke := range m.invokes {
		if err := m.executeInvoke(invoke); err != nil {
			return err
		}
	}
	return nil
}

func (m *module) executeInvoke(i invoke) (err error) {
	fnName := fxreflect.FuncName(i.Target)
	m.app.log.LogEvent(&fxevent.Invoking{
		FunctionName: fnName,
	})
	err = runInvoke(m.scope, i)
	defer func() {
		m.app.log.LogEvent(&fxevent.Invoked{
			FunctionName: fnName,
			Err:          err,
			Trace:        fmt.Sprintf("%+v", i.Stack), // format stack trace as multi-line
		})
	}()
	return err
}

func (m *module) decorate() (err error) {
	for _, decorator := range m.decorators {
		if err := runDecorator(m.scope, decorator); err != nil {
			return err
		}
	}
	return nil
}
