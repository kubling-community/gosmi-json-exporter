package mib

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sleepinggenius2/gosmi"
	"github.com/sleepinggenius2/gosmi/types"
)

type Node struct {
	Name        string `json:"name"`
	Oid         string `json:"oid"`
	Description string `json:"description,omitempty"`
}

type TableColumn struct {
	Name        string `json:"name"`
	Oid         string `json:"oid"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

type TableDefinition struct {
	Name           string        `json:"table"`
	Description    string        `json:"description,omitempty"`
	OidPrefix      string        `json:"snmp_oid_prefix"`
	Columns        []TableColumn `json:"columns"`
	OriginalOid    string        `json:"oid,omitempty"`
	OriginalModule string        `json:"module,omitempty"`
}

type DumpOptions struct {
	DumpTables bool
	ScalarMode string
	GroupDepth int
}

func DumpModule(mod gosmi.SmiModule, opts DumpOptions) ([]byte, error) {
	validModes := map[string]bool{"none": true, "separate": true, "grouped": true, "all": true}
	if !validModes[opts.ScalarMode] {
		return nil, fmt.Errorf("invalid scalar mode: %s", opts.ScalarMode)
	}

	if opts.DumpTables {
		tables := extractTables(mod)
		return json.MarshalIndent(map[string]interface{}{
			"module": mod.Name,
			"tables": tables,
		}, "", "  ")
	}

	switch opts.ScalarMode {
	case "separate", "grouped", "all":
		scalars := extractScalars(mod, opts.ScalarMode, opts.GroupDepth)
		return json.MarshalIndent(map[string]interface{}{
			"module": mod.Name,
			"tables": scalars,
		}, "", "  ")
	default:
		nodes := extractNodes(mod)
		return json.MarshalIndent(map[string]interface{}{
			"module": mod.Name,
			"nodes":  nodes,
		}, "", "  ")
	}
}

func extractNodes(mod gosmi.SmiModule) []Node {
	nodes := []Node{}
	for _, n := range mod.GetNodes() {
		nodes = append(nodes, Node{
			Name:        n.Name,
			Oid:         n.Oid.String(),
			Description: n.Description,
		})
	}
	return nodes
}

func extractTables(mod gosmi.SmiModule) []TableDefinition {
	nodes := mod.GetNodes()
	tables := []TableDefinition{}

	for _, table := range nodes {
		if table.Kind != types.NodeTable {
			continue
		}

		tableOid := table.Oid.String()
		var rowOid string
		for _, maybeRow := range nodes {
			if maybeRow.Kind == types.NodeRow &&
				strings.HasPrefix(maybeRow.Oid.String(), tableOid+".") {
				rowOid = maybeRow.Oid.String()
				break
			}
		}
		if rowOid == "" {
			continue
		}

		columns := []TableColumn{}
		for _, col := range nodes {
			if col.Kind == types.NodeColumn &&
				strings.HasPrefix(col.Oid.String(), rowOid+".") &&
				col.Type != nil {
				columns = append(columns, TableColumn{
					Name:        col.Name,
					Oid:         col.Oid.String(),
					Type:        col.Type.Name,
					Description: col.Description,
				})
			}
		}

		tables = append(tables, TableDefinition{
			Name:           table.Name,
			Description:    table.Description,
			OidPrefix:      tableOid,
			Columns:        columns,
			OriginalOid:    tableOid,
			OriginalModule: mod.Name,
		})
	}

	return tables
}

func extractScalars(mod gosmi.SmiModule, mode string, groupDepth int) []TableDefinition {
	nodes := mod.GetNodes()
	grouped := map[string][]TableColumn{}
	var tables []TableDefinition

	for _, n := range nodes {
		if n.Kind != types.NodeScalar || n.Type == nil {
			continue
		}

		col := TableColumn{
			Name:        n.Name,
			Oid:         n.Oid.String(),
			Type:        n.Type.Name,
			Description: n.Description,
		}

		switch mode {
		case "separate":
			tables = append(tables, TableDefinition{
				Name:           "scalar_" + n.Name,
				Description:    n.Description,
				OidPrefix:      n.Oid.String(),
				Columns:        []TableColumn{col},
				OriginalModule: mod.Name,
			})
		case "grouped":
			parts := strings.Split(n.Oid.String(), ".")
			depth := groupDepth
			if depth <= 0 || depth > len(parts) {
				depth = len(parts)
			}
			prefix := strings.Join(parts[:depth], ".")
			grouped[prefix] = append(grouped[prefix], col)
		case "all":
			grouped["__all__"] = append(grouped["__all__"], col)
		}
	}

	if mode == "grouped" || mode == "all" {
		for prefix, cols := range grouped {
			tables = append(tables, TableDefinition{
				Name:           "scalar_group_" + strings.ReplaceAll(prefix, ".", "_"),
				Description:    "Grouped scalar metrics under " + prefix,
				OidPrefix:      prefix,
				Columns:        cols,
				OriginalModule: mod.Name,
			})
		}
	}

	return tables
}
