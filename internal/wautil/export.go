package wautil

import (
	"context"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// inspired by github.com/ncruces/go-sqlite3@v0.29.1/internal/util/func.go, but has a more functional design and supports named parameters

type funcIIII[TR, T0, T1, T2 i32] func(context.Context, api.Module, T0, T1, T2) TR

func (fn funcIIII[TR, T0, T1, T2]) Call(ctx context.Context, mod api.Module, stack []uint64) {
	_ = stack[2]
	stack[0] = uint64(fn(ctx, mod, T0(stack[0]), T1(stack[1]), T2(stack[2])))
}

func ExportFuncIIII[TR, T0, T1, T2 i32](exportName string, fn func(context.Context, api.Module, T0, T1, T2) TR, name ...string) func(wazero.HostModuleBuilder) {
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

type funcVIIII[T0, T1, T2, T3 i32] func(context.Context, api.Module, T0, T1, T2, T3)

func (fn funcVIIII[T0, T1, T2, T3]) Call(ctx context.Context, mod api.Module, stack []uint64) {
	_ = stack[3]
	fn(ctx, mod, T0(stack[0]), T1(stack[1]), T2(stack[2]), T3(stack[3]))
}

func ExportFuncVIIII[T0, T1, T2, T3 i32](exportName string, fn func(context.Context, api.Module, T0, T1, T2, T3), name ...string) func(wazero.HostModuleBuilder) {
	return func(mod wazero.HostModuleBuilder) {
		b := mod.NewFunctionBuilder()
		b.WithGoModuleFunction(funcVIIII[T0, T1, T2, T3](fn), []api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32}, nil)
		if len(name) != 0 {
			if len(name) != 5 {
				panic("incorrect number of names")
			}
			b.WithName(name[0])
			b.WithParameterNames(name[1], name[2], name[3], name[4])
		}
		b.Export(exportName)
	}
}

type funcVIIIIII[T0, T1, T2, T3, T4, T5 i32] func(context.Context, api.Module, T0, T1, T2, T3, T4, T5)

func (fn funcVIIIIII[T0, T1, T2, T3, T4, T5]) Call(ctx context.Context, mod api.Module, stack []uint64) {
	_ = stack[5]
	fn(ctx, mod, T0(stack[0]), T1(stack[1]), T2(stack[2]), T3(stack[3]), T4(stack[4]), T5(stack[5]))
}

func ExportFuncVIIIIII[T0, T1, T2, T3, T4, T5 i32](exportName string, fn func(context.Context, api.Module, T0, T1, T2, T3, T4, T5), name ...string) func(wazero.HostModuleBuilder) {
	return func(mod wazero.HostModuleBuilder) {
		b := mod.NewFunctionBuilder()
		b.WithGoModuleFunction(funcVIIIIII[T0, T1, T2, T3, T4, T5](fn), []api.ValueType{api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32, api.ValueTypeI32}, nil)
		if len(name) != 0 {
			if len(name) != 7 {
				panic("incorrect number of names")
			}
			b.WithName(name[0])
			b.WithParameterNames(name[1], name[2], name[3], name[4], name[5], name[6])
		}
		b.Export(exportName)
	}
}
