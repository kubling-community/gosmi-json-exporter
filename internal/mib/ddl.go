package mib

import (
	"encoding/json"
	"strings"
)

const indent = "    " // 4 spaces

func GenerateSQLSchema(jsonData []byte) (string, error) {
	var payload struct {
		Module string            `json:"module"`
		Tables []TableDefinition `json:"tables"`
	}
	if err := json.Unmarshal(jsonData, &payload); err != nil {
		return "", err
	}

	var sb strings.Builder
	for _, table := range payload.Tables {
		sb.WriteString("CREATE FOREIGN TABLE ")
		sb.WriteString(table.Name)
		sb.WriteString("\n(\n")

		for i, col := range table.Columns {
			if i > 0 {
				sb.WriteString(",\n")
			}
			sb.WriteString(indent)
			sb.WriteString(col.Name)
			sb.WriteString(" ")
			sb.WriteString(toKublingType(col.Type))
			sb.WriteString(" OPTIONS (snmp_oid '")
			sb.WriteString(col.Oid)
			sb.WriteString("'")
			if col.Description != "" {
				sb.WriteString(", ANNOTATION '")
				sb.WriteString(formatDescription(col.Description))
				sb.WriteString("'")
			}
			sb.WriteString(")")
		}

		sb.WriteString("\n)\nOPTIONS (updatable 'false', snmp_type 'full_table'")
		if table.Description != "" {
			sb.WriteString(", ANNOTATION '")
			sb.WriteString(formatDescription(table.Description))
			sb.WriteString("'")
		}
		sb.WriteString(");\n\n")
	}

	return sb.String(), nil
}

func toKublingType(snmpType string) string {
	switch strings.ToLower(snmpType) {
	case "integer", "integer32", "gauge32", "counter32", "unsigned32", "int", "uint", "counter":
		return "integer"
	case "counter64":
		return "biginteger"
	case "octetstring", "displaystring", "physaddress", "macaddress", "string", "utf8string", "bitstring", "timeticks":
		return "string"
	case "ipaddress":
		return "ip"
	case "objectidentifier", "oid":
		return "string"
	case "truthvalue", "boolean":
		return "boolean"
	case "float", "float32", "float64", "double":
		return "float"
	default:
		return "string"
	}
}

func formatDescription(desc string) string {
	return strings.ReplaceAll(strings.ReplaceAll(desc, "'", "''"), "\n", " ")
}
