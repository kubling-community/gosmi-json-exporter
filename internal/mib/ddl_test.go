package mib

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestGenerateSQLSchema_Basic(t *testing.T) {
	data := map[string]interface{}{
		"module": "testModule",
		"tables": []TableDefinition{
			{
				Name:        "testTable",
				Description: "This is a test table.",
				OidPrefix:   "1.3.6.1.4.1.9999.1",
				Columns: []TableColumn{
					{
						Name:        "testColumn1",
						Oid:         "1.3.6.1.4.1.9999.1.1",
						Type:        "integer",
						Description: "First column",
					},
					{
						Name:        "testColumn2",
						Oid:         "1.3.6.1.4.1.9999.1.2",
						Type:        "string",
						Description: "Second column",
					},
				},
			},
		},
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("failed to marshal test data: %v", err)
	}

	ddl, err := GenerateSQLSchema(jsonData)
	if err != nil {
		t.Fatalf("GenerateSQLSchema returned error: %v", err)
	}

	if !strings.Contains(ddl, "snmp_oid") {
		t.Error("Expected snmp_oid OPTION to be included")
	}
}
