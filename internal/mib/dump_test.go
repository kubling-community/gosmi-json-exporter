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

func TestDumpModule_GroupedScalars(t *testing.T) {
	gosmi.Init()
	gosmi.AppendPath("../../testdata")

	modName, err := gosmi.LoadModule("IF-MIB")
	if err != nil {
		t.Fatalf("Failed to load IF-MIB: %v", err)
	}
	mod, err := gosmi.GetModule(modName)
	if err != nil {
		t.Fatalf("Failed to get module: %v", err)
	}

	data, err := DumpModule(mod, DumpOptions{
		ScalarMode: "grouped",
		GroupDepth: 5,
	})
	if err != nil {
		t.Fatalf("DumpModule failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if _, ok := result["tables"]; !ok {
		t.Errorf("Expected 'tables' key in output for grouped scalar mode")
	}
}

func TestDumpModule_AllScalars(t *testing.T) {
	gosmi.Init()
	gosmi.AppendPath("../../testdata")

	modName, err := gosmi.LoadModule("IF-MIB")
	if err != nil {
		t.Fatalf("Failed to load IF-MIB: %v", err)
	}
	mod, err := gosmi.GetModule(modName)
	if err != nil {
		t.Fatalf("Failed to get module: %v", err)
	}

	data, err := DumpModule(mod, DumpOptions{
		ScalarMode: "all",
	})
	if err != nil {
		t.Fatalf("DumpModule failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if _, ok := result["tables"]; !ok {
		t.Errorf("Expected 'tables' key in output for scalar-mode=all")
	}
}
