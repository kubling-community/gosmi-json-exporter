package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sleepinggenius2/gosmi"
	"kubling.com/gosmi-json-exporter/internal/mib"
)

const version = "25.1"

func main() {

	dir := flag.String("mibdir", ".", "Directory containing MIB files")
	moduleName := flag.String("module", "", "Module name to process (required)")
	dumpTables := flag.Bool("dump-tables", false, "Dump SNMP table definitions instead of raw nodes")
	scalarMode := flag.String("scalar-mode", "none", "How to emit scalar nodes: none, separate, grouped, all")
	groupDepth := flag.Int("group-depth", 10, "When scalar-mode=grouped, how many OID parts to group by")
	outFile := flag.String("out", "", "Write output to a file (default: stdout)")
	showVersion := flag.Bool("version", false, "Print version and exit")
	sqlSchema := flag.Bool("sql-schema", false, "Emit SQL CREATE FOREIGN TABLE schema for Kubling")

	flag.Parse()

	if *showVersion {
		fmt.Println("gosmi-json-exporter version", version)
		os.Exit(0)
	}

	if *moduleName == "" {
		log.Fatal("missing required --module flag")
	}

	gosmi.Init()
	gosmi.AppendPath(*dir)

	loadedName, err := gosmi.LoadModule(*moduleName)
	if err != nil {
		log.Fatalf("Failed to load module %q: %v", *moduleName, err)
	}

	mod, err := gosmi.GetModule(loadedName)
	if err != nil {
		log.Fatalf("Failed to get module: %v", err)
	}

	opts := mib.DumpOptions{
		DumpTables: *dumpTables || *sqlSchema,
		ScalarMode: *scalarMode,
		GroupDepth: *groupDepth,
	}

	jsonData, err := mib.DumpModule(mod, opts)
	if err != nil {
		log.Fatalf("Failed to dump module: %v", err)
	}

	var content []byte
	if *sqlSchema {
		sqlText, err := mib.GenerateSQLSchema(jsonData)
		if err != nil {
			log.Fatalf("Failed to generate SQL schema: %v", err)
		}
		content = []byte(sqlText)
	} else {
		content = jsonData
	}

	// Write to file or stdout
	if *outFile != "" {
		if err := os.WriteFile(*outFile, content, 0644); err != nil {
			log.Fatalf("Failed to write to output file: %v", err)
		}
	} else {
		if _, err := os.Stdout.Write(content); err != nil {
			log.Fatalf("Failed to write to stdout: %v", err)
		}
	}
}
