package sublime

import (
	"code.google.com/p/log4go"
	"fmt"
	"lime/3rdparty/libs/gopy/lib"
	"lime/backend"
	"lime/backend/primitives"
)

var (
	_ = backend.View{}
	_ = primitives.Region{}
)

var (
	_onQueryContextGlueClass = py.Class{
		Name:    "sublime.OnQueryContextGlue",
		Pointer: (*OnQueryContextGlue)(nil),
	}
	_viewEventGlueClass = py.Class{
		Name:    "sublime.ViewEventGlue",
		Pointer: (*ViewEventGlue)(nil),
	}
)

type (
	OnQueryContextGlue struct {
		py.BaseObject
		inner py.Object
	}
	ViewEventGlue struct {
		py.BaseObject
		inner py.Object
	}
)

var evmap = map[string]*backend.ViewEvent{
	"on_modified":    &backend.OnModified,
	"on_activated":   &backend.OnActivated,
	"on_deactivated": &backend.OnDeactivated,
	"on_load":        &backend.OnLoad,
	"on_new":         &backend.OnNew,
	"on_pre_save":    &backend.OnPreSave,
	"on_post_save":   &backend.OnPostSave,
}

func (c *ViewEventGlue) PyInit(args *py.Tuple, kwds *py.Dict) error {
	if args.Size() != 2 {
		return fmt.Errorf("Expected 2 arguments not %d", args.Size())
	}
	if v, err := args.GetItem(0); err != nil {
		return err
	} else {
		c.inner = v
	}
	if v, err := args.GetItem(1); err != nil {
		return err
	} else if v2, ok := v.(*py.Unicode); !ok {
		return fmt.Errorf("Second argument not a string: %v", v)
	} else {
		ev := evmap[v2.String()]
		if ev == nil {
			return fmt.Errorf("Unknown event: %s", v2)
		}
		ev.Add(c.onEvent)
		c.inner.Incref()
		c.Incref()
	}
	return nil
}

func (c *ViewEventGlue) onEvent(v *backend.View) {
	if pv, err := toPython(v); err != nil {
		log4go.Error(err)
	} else {
		log4go.Debug("onEvent: %v, %v, %v", c, c.inner, pv)
		if ret, err := c.inner.Base().CallFunctionObjArgs(pv); err != nil {
			log4go.Error(err)
		} else if ret != nil {
			ret.Decref()
		}
	}
}

func (c *OnQueryContextGlue) PyInit(args *py.Tuple, kwds *py.Dict) error {
	if args.Size() != 1 {
		return fmt.Errorf("Expected only 1 argument not %d", args.Size())
	}
	if v, err := args.GetItem(0); err != nil {
		return err
	} else {
		c.inner = v
	}
	c.inner.Incref()
	c.Incref()

	backend.OnQueryContext.Add(c.onQueryContext)
	return nil
}

func (c *OnQueryContextGlue) onQueryContext(v *backend.View, key string, operator backend.Op, operand interface{}, match_all bool) backend.QueryContextReturn {
	if pv, err := toPython(v); err != nil {
		log4go.Error(err)
	} else if pk, err := toPython(key); err != nil {
		log4go.Error(err)
	} else if po, err := toPython(operator); err != nil {
		log4go.Error(err)
	} else if poa, err := toPython(operand); err != nil {
		log4go.Error(err)
	} else if pm, err := toPython(match_all); err != nil {
		log4go.Error(err)
	} else if ret, err := c.inner.Base().CallFunctionObjArgs(pv, pk, po, poa, pm); err != nil {
		log4go.Error(err)
	} else if ret != nil {
		defer ret.Decref()
		if r2, ok := ret.(*py.Bool); ok {
			if r2.Bool() {
				return backend.True
			} else {
				return backend.False
			}
		} else {
			log4go.Debug("other: %v", ret)
		}
	}
	return backend.Unknown
}