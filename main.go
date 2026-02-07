package gochain

import (
	"errors"
	"fmt"
	"reflect"
)

type Chain[T, U any] struct {
	ctx        T
	result     U
	handlers   []Handler[T, U]
	forcedStop bool
}

type StopFunc func()

type Handler[T, U any] func(chain *Chain[T, U], stop StopFunc) error

var (
	ErrContextIsUndefined         = errors.New("context is undefined")
	ErrContextShouldBeStruct      = errors.New("context should be a struct")
	ErrResultIsUndefined          = errors.New("result is undefined")
	ErrResultShouldBeStruct       = errors.New("result should be a struct")
	ErrStructFieldNotExported     = errors.New("struct field is not exported")
	ErrStructFieldCannotBeChanged = errors.New("struct field cannot be changed")
	ErrIncompatibleValueType      = errors.New("incompatible value type")
)

func NewChain[T, U any]() *Chain[T, U] {
	var (
		ctx    T
		result U
	)

	return &Chain[T, U]{
		ctx:      ctx,
		result:   result,
		handlers: make([]Handler[T, U], 0),
	}
}

func (c *Chain[T, U]) Add(handler Handler[T, U]) *Chain[T, U] {
	c.handlers = append(c.handlers, handler)
	return c
}

func stop[T, U any](chain *Chain[T, U]) func() {
	return func() {
		chain.forcedStop = true
	}
}

func (c *Chain[T, U]) Run() error {
	for _, handler := range c.handlers {
		if err := handler(c, stop(c)); err != nil {
			return err
		} else if c.forcedStop {
			return nil
		}
	}

	return nil
}

func (c *Chain[T, U]) GetContext() T {
	return c.ctx
}

func (c *Chain[T, U]) GetResult() U {
	return c.result
}

func (c *Chain[T, U]) UpdateContext(fieldName string, value any) error {
	pointer := reflect.ValueOf(&c.ctx)

	if pointer.Kind() != reflect.Pointer {
		return ErrContextIsUndefined
	}

	pointerValue := pointer.Elem()

	if pointerValue.Kind() != reflect.Struct {
		return ErrContextShouldBeStruct
	}

	field := pointerValue.FieldByName(fieldName)

	if !field.IsValid() {
		return errors.Join(
			ErrStructFieldNotExported,
			fmt.Errorf("the %s field does not exist", fieldName),
		)
	} else if !field.CanSet() {
		return errors.Join(
			ErrStructFieldCannotBeChanged,
			fmt.Errorf("the %s field cannot be changed", fieldName),
		)
	}

	fieldValue := reflect.ValueOf(value)

	if field.Type() != fieldValue.Type() {
		return errors.Join(
			ErrIncompatibleValueType,
			fmt.Errorf("The value type is incompatible with the %s field", fieldName),
		)
	}

	field.Set(fieldValue)
	return nil
}

func (c *Chain[T, U]) UpdateResult(fieldName string, value any) error {
	pointer := reflect.ValueOf(&c.result)

	if pointer.Kind() != reflect.Pointer {
		return ErrResultIsUndefined
	}

	pointerValue := pointer.Elem()

	if pointerValue.Kind() != reflect.Struct {
		return ErrResultShouldBeStruct
	}

	field := pointerValue.FieldByName(fieldName)

	if !field.IsValid() {
		return ErrStructFieldNotExported
	} else if !field.CanSet() {
		return errors.Join(
			ErrStructFieldCannotBeChanged,
			fmt.Errorf("the %s field cannot be changed", fieldName),
		)
	}

	fieldValue := reflect.ValueOf(value)

	if field.Type() != fieldValue.Type() {
		return errors.Join(
			ErrIncompatibleValueType,
			fmt.Errorf(
				"The value type is incompatible with the %s field", fieldName,
			),
		)
	}

	field.Set(fieldValue)
	return nil
}
