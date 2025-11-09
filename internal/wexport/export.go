// Package wexport contains type-safe wrappers for exporting host functions.
package wexport

import (
	"context"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// inspired by github.com/ncruces/go-sqlite3@v0.29.1/internal/util/func.go, but has a more functional design and supports named parameters

type (
	i32 interface{ ~int32 | ~uint32 }
	i64 interface{ ~int64 | ~uint64 }
)

type Func func(wazero.HostModuleBuilder)

func Instantiate(ctx context.Context, runtime wazero.Runtime, name string, fn ...Func) (api.Module, error) {
	hmb := runtime.NewHostModuleBuilder(name)
	for _, fn := range fn {
		fn(hmb)
	}
	return hmb.Instantiate(ctx)
}

type funcIIII[TR, T0, T1, T2 i32] func(context.Context, api.Module, T0, T1, T2) TR

func (fn funcIIII[TR, T0, T1, T2]) Call(ctx context.Context, mod api.Module, stack []uint64) {
	_ = stack[2]
	stack[0] = uint64(fn(ctx, mod, T0(stack[0]), T1(stack[1]), T2(stack[2])))
}

func IIII[TR, T0, T1, T2 i32](exportName string, fn func(context.Context, api.Module, T0, T1, T2) TR, name ...string) Func {
	return func(mod wazero.HostModuleBuilder) {
		b := mod.NewFunctionBuilder()
		b.WithGoModuleFunction(funcIIII[TR, T0, T1, T2](fn), []api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{api.ValueTypeI32})
		if len(name) != 0 {
			if len(name) != 5 {
				panic("incorrect number of names")
			}
			b.WithName(name[0])
			b.WithParameterNames(name[1], name[2], name[3])
			b.WithResultNames(name[4])
		}
		b.Export(exportName)
	}
}

type funcVII[T0, T1 i32] func(context.Context, api.Module, T0, T1)

func (fn funcVII[T0, T1]) Call(ctx context.Context, mod api.Module, stack []uint64) {
	_ = stack[1]
	fn(ctx, mod, T0(stack[0]), T1(stack[1]))
}

func VII[T0, T1 i32](exportName string, fn func(context.Context, api.Module, T0, T1), name ...string) Func {
	return func(mod wazero.HostModuleBuilder) {
		b := mod.NewFunctionBuilder()
		b.WithGoModuleFunction(funcVII[T0, T1](fn), []api.ValueType{api.ValueTypeI32, api.ValueTypeI32}, nil)
		if len(name) != 0 {
			if len(name) != 3 {
				panic("incorrect number of names")
			}
			b.WithName(name[0])
			b.WithParameterNames(name[1], name[2])
		}
		b.Export(exportName)
	}
}

type funcVIII[T0, T1, T2 i32] func(context.Context, api.Module, T0, T1, T2)

func (fn funcVIII[T0, T1, T2]) Call(ctx context.Context, mod api.Module, stack []uint64) {
	_ = stack[2]
	fn(ctx, mod, T0(stack[0]), T1(stack[1]), T2(stack[2]))
}

func VIII[T0, T1, T2 i32](exportName string, fn func(context.Context, api.Module, T0, T1, T2), name ...string) Func {
	return func(mod wazero.HostModuleBuilder) {
		b := mod.NewFunctionBuilder()
		b.WithGoModuleFunction(funcVIII[T0, T1, T2](fn), []api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32}, nil)
		if len(name) != 0 {
			if len(name) != 4 {
				panic("incorrect number of names")
			}
			b.WithName(name[0])
			b.WithParameterNames(name[1], name[2], name[3])
		}
		b.Export(exportName)
	}
}
