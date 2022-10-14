package wasmtime_test

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/iotexproject/w3bstream/pkg/modules/vm"
	"github.com/iotexproject/w3bstream/pkg/modules/vm/wasmtime"
	"github.com/iotexproject/w3bstream/pkg/types/wasm"
)

// //go:embed ../../../../examples/log/log.wasm
var wasmLogCode []byte

// //go:embed ../../../../examples/word_count/word_count.wasm
var wasmWordCountCode []byte

// //go:embed ../../../../examples/word_count_v2/word_count_v2.wasm
var wasmWordCountV2Code []byte

func init() {
	wd, _ := os.Getwd()
	fmt.Println(wd)
	root := filepath.Join(wd, "../../../../examples")
	fmt.Println(root)

	var err error
	wasmLogCode, err = os.ReadFile(filepath.Join(root, "log/log.wasm"))
	if err != nil {
		panic(err)
	}

	wasmWordCountCode, err = os.ReadFile(filepath.Join(root, "word_count/word_count.wasm"))
	if err != nil {
		panic(err)
	}
	wasmWordCountV2Code, err = os.ReadFile(filepath.Join(root, "word_count_v2/word_count_v2.wasm"))
	if err != nil {
		panic(err)
	}
}

func TestInstance_LogWASM(t *testing.T) {
	require := require.New(t)
	i, err := wasmtime.NewInstanceByCode(wasmLogCode)
	require.NoError(err)
	id := vm.AddInstance(i)

	err = vm.StartInstance(id)
	require.NoError(err)
	defer vm.StopInstance(id)

	_, code := i.HandleEvent("start", []byte("IoTeX"))
	NewWithT(t).Expect(code).To(Equal(wasm.ResultStatusCode_OK))

	_, code = i.HandleEvent("not_exported", []byte("IoTeX"))
	NewWithT(t).Expect(code).To(Equal(wasm.ResultStatusCode_UnexportedHandler))
}

func TestInstance_WordCount(t *testing.T) {
	require := require.New(t)
	i, err := wasmtime.NewInstanceByCode(wasmWordCountCode)
	require.NoError(err)
	id := vm.AddInstance(i)

	err = vm.StartInstance(id)
	require.NoError(err)
	defer vm.StopInstance(id)

	_, code := i.HandleEvent("start", []byte("a b c d a"))
	require.Equal(wasm.ResultStatusCode_OK, code)

	require.Equal(int32(2), i.Get("a"))
	require.Equal(int32(1), i.Get("b"))
	require.Equal(int32(1), i.Get("c"))
	require.Equal(int32(1), i.Get("d"))

	_, code = i.HandleEvent("start", []byte("a b c d a"))
	require.Equal(wasm.ResultStatusCode_OK, code)

	require.Equal(int32(4), i.Get("a"))
	require.Equal(int32(2), i.Get("b"))
	require.Equal(int32(2), i.Get("c"))
	require.Equal(int32(2), i.Get("d"))
}

func TestInstance_WordCountV2(t *testing.T) {
	require := require.New(t)
	i, err := wasmtime.NewInstanceByCode(wasmWordCountV2Code)
	require.NoError(err)
	id := vm.AddInstance(i)

	err = vm.StartInstance(id)
	require.NoError(err)
	defer vm.StopInstance(id)

	_, code := i.HandleEvent("count", []byte("a b c d a"))
	NewWithT(t).Expect(code).To(Equal(wasm.ResultStatusCode_OK))

	NewWithT(t).Expect(i.Get("a")).To(Equal(int32(2)))
	NewWithT(t).Expect(i.Get("b")).To(Equal(int32(1)))
	NewWithT(t).Expect(i.Get("c")).To(Equal(int32(1)))
	NewWithT(t).Expect(i.Get("d")).To(Equal(int32(1)))

	_, code = i.HandleEvent("count", []byte("a b c d a"))
	NewWithT(t).Expect(code).To(Equal(wasm.ResultStatusCode_OK))

	NewWithT(t).Expect(i.Get("a")).To(Equal(int32(4)))
	NewWithT(t).Expect(i.Get("b")).To(Equal(int32(2)))
	NewWithT(t).Expect(i.Get("c")).To(Equal(int32(2)))
	NewWithT(t).Expect(i.Get("d")).To(Equal(int32(2)))

	_, unique := i.HandleEvent("unique", nil)
	NewWithT(t).Expect(unique).To(Equal(wasm.ResultStatusCode(4)))
}