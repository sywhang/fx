// Copyright (c) 2020 Uber Technologies, Inc.
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
	"strings"

	"go.uber.org/fx/internal/fxreflect"
)

// Annotated annotates a constructor provided to Fx with additional options.
//
// For example,
//
//   func NewReadOnlyConnection(...) (*Connection, error)
//
//   fx.Provide(fx.Annotated{
//     Name: "ro",
//     Target: NewReadOnlyConnection,
//   })
//
// Is equivalent to,
//
//   type result struct {
//     fx.Out
//
//     Connection *Connection `name:"ro"`
//   }
//
//   fx.Provide(func(...) (result, error) {
//     conn, err := NewReadOnlyConnection(...)
//     return result{Connection: conn}, err
//   })
//
// Annotated cannot be used with constructors which produce fx.Out objects.
//
// When used with fx.Supply, the target is a value rather than a constructor function.
type Annotated struct {
	// If specified, this will be used as the name for all non-error values returned
	// by the constructor. For more information on named values, see the documentation
	// for the fx.Out type.
	//
	// A name option may not be provided if a group option is provided.
	Name string

	// If specified, this will be used as the group name for all non-error values returned
	// by the constructor. For more information on value groups, see the package documentation.
	//
	// A group option may not be provided if a name option is provided.
	//
	// Similar to group tags, the group name may be followed by a `,flatten`
	// option to indicate that each element in the slice returned by the
	// constructor should be injected into the value group individually.
	Group string

	// Target is the constructor or value being annotated with fx.Annotated.
	Target interface{}
}

func (a Annotated) String() string {
	var fields []string
	if len(a.Name) > 0 {
		fields = append(fields, fmt.Sprintf("Name: %q", a.Name))
	}
	if len(a.Group) > 0 {
		fields = append(fields, fmt.Sprintf("Group: %q", a.Group))
	}
	if a.Target != nil {
		fields = append(fields, fmt.Sprintf("Target: %v", fxreflect.FuncName(a.Target)))
	}
	return fmt.Sprintf("fx.Annotated{%v}", strings.Join(fields, ", "))
}

type Annotation interface {
	apply(ctor interface{})
}

type paramTag struct {
	tags []string
}

type resultTag struct {
	tags []string
}

func (p *paramTag) apply(fn interface{}) (interface{}, error) {
	// Verify all tags can be applied to the inputs
	fnVal := reflect.ValueOf(fn)
	fnType := fnVal.Type()
	numIn := fnType.NumIn()
	numOut := fnType.NumOut()
	if numIn < len(p.tags) {
		// Error out
		return nil, errors.New("cannot apply function because this is sad")
	}
	digInStructFields := []reflect.StructField{{
		Name:      "In",
		Anonymous: true,
		Type:      reflect.TypeOf(In{}),
	}}
	for i := 0; i < numIn(); i++ {
		name := fmt.Sprintf("Field%d", i)
		field := reflect.StructField{
			Name: name,
			Type: fnType.In(i)
			Tag: reflect.StructTag(p.tags[i])
		}
		digInStructFields = append(digInStructFields, field)
	}
	outs := make([]reflect.Type, numOut)
	for i := 0; i < numOut; i++ {
		outs[i] = fnType.Out(i)
	}
	inStructType := reflect.StructOf(digInStructFields)
	newFuncType := reflect.FuncOf([]reflect.Type{inStructType}, outs, false)
	newF := reflect.MakeFunc(newFuncType, func(args []reflect.Value) []reflect.Value {
		fnArgs := make([]reflect.Value, numIns)
		params := args[0]
		for i := 0; i < numIns; i++ {
			fnArgs[i] = params.Field(i+1)
		}
		return fnVal.Call(fnArgs)
	})
	return newF, nil
}

func ParamTags(tags ...string) Annotation {
	t := paramTag{}
	t.tags := make([]string, len(tags))
	// TODO: validation?
	for i, tag := range tags {
		t.tags[i] = tag
	}
	return t
}

func ResultTags(tags ...string) Annotation {
	t := resultTag{}
	t.tags := make([]string, len(tags))
	// TODO: validation?
	for i, tag := range tags {
		t.tags[i] = tag
	}
	return t
}

struct annotatedFunc type {
	target interface{}
	paramTag Annotation
	resultTag Annotation
}

func Annotate(fn interface{}, anns ...Annotation) interface{} {
	for _, ann := range anns {
		fn = ann.apply(fn)
	}
	digInStructFields := []reflect.StructField{{
	Name:      "In",
	Anonymous: true,
	Type:      reflect.TypeOf(In{}),
		}}
		for i := 0; i < numArgs; i++ {
			name := fmt.Sprintf("Field%d", i)
			field := reflect.StructField{
				Name: name,
				Type: userFuncType.In(i),
			}
			if i < numNames { // namedArguments
				tag := ""
				annos[i].isAnnotation()
				switch anno := annos[i].(type) {
				case groupAnnotation:
					tag = fmt.Sprintf(`group:"%s"`, anno.group)
				case nameAnnotation:
					tag = fmt.Sprintf(`name:"%s"`, anno.name)
				}

				field.Tag = reflect.StructTag(tag)
			}
			digInStructFields = append(digInStructFields, field)
		}
}
