package mib

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/sleepinggenius2/gosmi"
)

func TestDumpRawNodes(t *testing.T) {
	gosmi.Init()
	gosmi.AppendPath("../../testdata")

	moduleName, err := gosmi.LoadModule("IF-MIB")
	if err != nil {
		t.Skipf("Skipping test: failed to load MIB module: %v", err)
	}

	mod, err := gosmi.GetModule(moduleName)
	if err != nil {
		t.Fatalf("Failed to get module: %v", err)
	}

	data, err := DumpModule(mod, DumpOptions{
		DumpTables: false,
		ScalarMode: "none",
	})
	if err != nil {
		t.Fatalf("DumpModule failed: %v", err)
	}

	//t.Logf("JSON output \n%s", string(data))

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Invalid JSON output: %v", err)
	}

	if parsed["module"] != "IF-MIB" {
		t.Errorf("Expected module IF-MIB, got %v", parsed["module"])
	}
}

func TestInvalidScalarMode(t *testing.T) {
	gosmi.Init()
	gosmi.AppendPath("../../testdata")

	moduleName, err := gosmi.LoadModule("IF-MIB")
	if err != nil {
		t.Skipf("Skipping test: failed to load MIB module: %v", err)
	}
	mod, err := gosmi.GetModule(moduleName)
	if err != nil {
		t.Fatalf("Failed to get module: %v", err)
	}

	_, err = DumpModule(mod, DumpOptions{
		ScalarMode: "invalid-mode",
	})
	if err == nil || !strings.Contains(err.Error(), "invalid scalar mode") {
		t.Errorf("Expected error for invalid scalar mode, got: %v", err)
	}
}
