# gosmi-json-exporter

A simple CLI tool that parses SNMP MIB modules using [GoSMI](https://github.com/sleepinggenius2/gosmi) and exports their structure as JSON and Kubling Schema.

This project was created to provide a fully structured JSON representation of SNMP MIBs (including scalars, tables, and column types) with the specific goal of integrating into **Kubling**.  
More information at [docs.kubling.com](https://docs.kubling.com/engine/ds/industrial/snmp).

---

## Features

- Parses and loads SNMP MIB files
- Outputs:
  - Raw node definitions
  - SNMP table structures
  - Scalar groupings (separate, grouped, or all)
- Outputs to stdout or file
- Designed for CLI usage and automation

---

## Installation

```bash
git clone https://github.com/youruser/gosmi-json-exporter
cd gosmi-json-exporter
make build
```

The built binary will be placed in `./bin/gosmi-json-exporter`.

---

## Usage

```bash
./bin/gosmi-json-exporter --mibdir ./testdata --module IF-MIB --dump-tables --out tables.json
```

### CLI Flags

| Flag            | Description                                                    |
|-----------------|----------------------------------------------------------------|
| `--mibdir`      | Directory containing `.mib` files (default: `.`)               |
| `--module`      | Name of the MIB module to load (required)                      |
| `--dump-tables` | Dump SNMP table structures instead of raw nodes                |
| `--scalar-mode` | How to emit scalar nodes: `none`, `separate`, `grouped`, `all` |
| `--group-depth` | Grouping depth for OID segments (used with `grouped`)          |
| `--out`         | Output file (if not specified, defaults to stdout)             |
| `--version`     | Show version and exit                                          |

---

## Example Output

```json
{
  "module": "IF-MIB",
  "nodes": [
    {
      "name": "interfaces",
      "oid": "1.3.6.1.2.1.2"
    },
    {
      "name": "ifNumber",
      "oid": "1.3.6.1.2.1.2.1",
      "description": "The number of network interfaces (regardless of their\ncurrent state) present on this system."
    },
    {
      "name": "ifTable",
      "oid": "1.3.6.1.2.1.2.2",
      "description": "A list of interface entries.  The number of entries is\ngiven by the value of ifNumber."
    }
  ]
}
```

---

## Development

### Run Tests

```bash
make test
```

Tests use a minimal, dependency-controlled MIB set under `./testdata`. You can extend this directory with your own `.mib` files.

---

## Working with MIBs

To get meaningful and well-structured output, you'll need access to valid and complete MIB files, many real-world MIBs have complex dependency chains and vendor-specific extensions.

We recommend exploring the [pgmillon/observium](https://github.com/pgmillon/observium) repository (or other similar collections), which includes a rich collection of cleaned and organized MIBs sourced from the Observium monitoring platform.

This collection is extremely helpful for:
- Testing the exporter with realistic SNMP modules
- Understanding how enterprise-level MIBs are structured
- Generating consistent and representative JSON outputs

You can extract specific modules or load the entire tree using the `--mibdir` flag:

```bash
./bin/gosmi-json-exporter --mibdir /path/to/observium/mibs --module UPS-MIB
```

### About `snmp_derived_index` Support

This tool **does not currently generate `snmp_derived_index` options** in the output DDL.

While the [Kubling SNMP engine](https://docs.kubling.com/engine/ds/industrial/snmp) supports advanced modeling of table indexes using `snmp_derived_index`, inferring these relationships reliably from raw MIB data is complex and often ambiguous. In most MIBs, the underlying relationships between scalar columns and OID indexes are not explicitly defined in a structured or machine-readable way, making automated detection highly error-prone.

For now, the tool generates valid SQL schemas using only `snmp_oid` references, which are sufficient for the majority of SNMP table interactions.

If you require `snmp_derived_index` for more advanced or nested table modeling (e.g., `ifStackTable`), you can:

- **Manually edit the generated DDL** and replace specific `snmp_oid` entries with the appropriate `snmp_derived_index` references.
- Use external reference projects like [`observium`](https://github.com/pgmillon/observium) to help explore and understand complex index structures before editing.

If you’re interested in contributing or experimenting with automatic detection of derived indexes, feel free to open a discussion or pull request.

## Acknowledgments

- [gosmi](https://github.com/sleepinggenius2/gosmi) for the excellent Go SNMP parser
- MIB files sourced from various RFC and vendor public repositories

---

## License

Licensed under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0).