package fx

import (
	"fmt"

	"go.uber.org/dig"
)

func runDecorator(c container, d decorator, opts ...dig.DecorateOption) error {
	decorator := d.Target

	switch decorator := decorator.(type) {
	case annotated:
		dcor, err := decorator.Build()
		if err != nil {
			return fmt.Errorf("fx.Decorate(%v) from:\n%+vFailed: %v", decorator, d.Stack, err)
		}

		if err := c.Decorate(dcor, opts...); err != nil {
			return fmt.Errorf("fx.Decorate(%v) from:\n%+vFailed: %v", decorator, d.Stack, err)
		}
	default:
		if err := c.Decorate(decorator, opts...); err != nil {
			return fmt.Errorf("fx.Decorate(%v) from:\n%+vFailed: %v", decorator, d.Stack, err)
		}
	}

	return nil
}
