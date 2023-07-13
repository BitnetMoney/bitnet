// Copyright 2023 Bitnet
// This file is part of the Bitnet library.
//
// This software is provided "as is", without warranty of any kind,
// express or implied, including but not limited to the warranties
// of merchantability, fitness for a particular purpose and
// noninfringement. In no even shall the authors or copyright
// holders be liable for any claim, damages, or other liability,
// whether in an action of contract, tort or otherwise, arising
// from, out of or in connection with the software or the use or
// other dealings in the software.

package jsre

import (
	"os"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/dop251/goja"
)

type testNativeObjectBinding struct {
	vm *goja.Runtime
}

type msg struct {
	Msg string
}

func (no *testNativeObjectBinding) TestMethod(call goja.FunctionCall) goja.Value {
	m := call.Argument(0).ToString().String()
	return no.vm.ToValue(&msg{m})
}

func newWithTestJS(t *testing.T, testjs string) *JSRE {
	dir := t.TempDir()
	if testjs != "" {
		if err := os.WriteFile(path.Join(dir, "test.js"), []byte(testjs), os.ModePerm); err != nil {
			t.Fatal("cannot create test.js:", err)
		}
	}
	jsre := New(dir, os.Stdout)
	return jsre
}

func TestExec(t *testing.T) {
	jsre := newWithTestJS(t, `msg = "testMsg"`)

	err := jsre.Exec("test.js")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	val, err := jsre.Run("msg")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if val.ExportType().Kind() != reflect.String {
		t.Errorf("expected string value, got %v", val)
	}
	exp := "testMsg"
	got := val.ToString().String()
	if exp != got {
		t.Errorf("expected '%v', got '%v'", exp, got)
	}
	jsre.Stop(false)
}

func TestNatto(t *testing.T) {
	jsre := newWithTestJS(t, `setTimeout(function(){msg = "testMsg"}, 1);`)

	err := jsre.Exec("test.js")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	time.Sleep(100 * time.Millisecond)
	val, err := jsre.Run("msg")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if val.ExportType().Kind() != reflect.String {
		t.Fatalf("expected string value, got %v", val)
	}
	exp := "testMsg"
	got := val.ToString().String()
	if exp != got {
		t.Fatalf("expected '%v', got '%v'", exp, got)
	}
	jsre.Stop(false)
}

func TestBind(t *testing.T) {
	jsre := New("", os.Stdout)
	defer jsre.Stop(false)

	jsre.Set("no", &testNativeObjectBinding{vm: jsre.vm})

	_, err := jsre.Run(`no.TestMethod("testMsg")`)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestLoadScript(t *testing.T) {
	jsre := newWithTestJS(t, `msg = "testMsg"`)

	_, err := jsre.Run(`loadScript("test.js")`)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	val, err := jsre.Run("msg")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if val.ExportType().Kind() != reflect.String {
		t.Errorf("expected string value, got %v", val)
	}
	exp := "testMsg"
	got := val.ToString().String()
	if exp != got {
		t.Errorf("expected '%v', got '%v'", exp, got)
	}
	jsre.Stop(false)
}
