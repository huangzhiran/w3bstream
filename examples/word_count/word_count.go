package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed testdata/word_count.wasm
var code []byte

func main() {
	ctx := context.Background()

	// new wasm runtime.
	r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig().
		WithFeatureBulkMemoryOperations(true).
		WithFeatureNonTrappingFloatToIntConversion(true).
		WithFeatureSignExtensionOps(true).WithFeatureMultiValue(true))
	defer r.Close(ctx)

	// exports
	{
		_, err := r.NewModuleBuilder("env").
			ExportFunction("log", log).
			ExportFunction("add", add).
			ExportFunction("get", get).
			Instantiate(ctx, r)
		if err != nil {
			panic(err)
		}
	}

	if _, err := wasi_snapshot_preview1.Instantiate(ctx, r); err != nil {
		panic(err)
	}

	mod, err := r.InstantiateModuleFromBinary(ctx, code)
	if err != nil {
		panic(err)
	}

	// fns
	var (
		counter = mod.ExportedFunction("countWords")
		_       = mod.ExportedFunction("greet")
		malloc  = mod.ExportedFunction("malloc")
		free    = mod.ExportedFunction("free")
	)

	str := os.Args[1]
	strlen := uint64(len(str))

	results, err := malloc.Call(ctx, strlen)
	if err != nil {
		panic(err)
	}
	ptr := results[0]
	defer free.Call(ctx, ptr)

	if !mod.Memory().Write(ctx, uint32(ptr), []byte(str)) {
		panic(fmt.Sprintf("Memory.Write(%d, %d) out of range of memory size %d",
			ptr, strlen, mod.Memory().Size(ctx)))
	}

	_, err = counter.Call(ctx, ptr, strlen)
	if err != nil {
		panic(err)
	}

	msg, _ := json.Marshal(words)
	fmt.Println("host >> " + string(msg))
}

func log(ctx context.Context, m api.Module, offset, size uint32) {
	buf, ok := m.Memory().Read(ctx, offset, size)
	if !ok {
		panic(fmt.Sprintf("Memory.Read(%d,%d) out of range)", offset, size))
	}
	fmt.Println(string(buf))
}

var words = make(map[string]int32)

func add(ctx context.Context, m api.Module, offset, size uint32) (code int32) {
	buf, ok := m.Memory().Read(ctx, offset, size)
	if !ok {
		return 1
	}
	str := string(buf)
	if _, ok := words[str]; !ok {
		words[str] = 1
	} else {
		words[str]++
	}
	return 0
}

func get(ctx context.Context, m api.Module, offset, size uint32) (count int32) {
	buf, ok := m.Memory().Read(ctx, offset, size)
	if !ok {
		return 1
	}
	str := string(buf)
	if _, ok := words[str]; !ok {
		return 0
	}
	return words[str]
}
