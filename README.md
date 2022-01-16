# xlsx2csv

Convert XLSX file to CSV.

## Installation

go install github.com/gotidy/xlsx2csv

## Converting

### Examples

```sh
xlsx2csv  -b -c utf-16le test.xlsx test-{{num}}.csv
```

```sh
xlsx2csv test.xlsx > test.csv
```

### Usage

```sh
Usage: xlsx2csv convert <source> [<output>]

Convert file.

Arguments:
  <source>      Paths to XLSX file.
  [<output>]    Output CSV file path. It may include variables {{name}} and
                {{num}}

Flags:
  -h, --help                       Show context-sensitive help.

  -c, --encoding="utf-8"           Encoding of output CSV
                                   (http://www.w3.org/TR/encoding).
  -b, --bom                        Write BOM to CSV.
  -d, --delimiter=","              Cells delimiter.
  -l, --line_terminator="\\n"      line terminator - lines terminator in CSV.
                                   Default is \n
  -n, --sheetname=SHEETNAME,...    Sheet names to convert.
  -s, --sheet=SHEET,...            Sheet indexes to convert.
  -N, --exclude_sheet_name=EXCLUDE_SHEET_NAME,...
                                   Exclude sheets with names.
  -S, --exclude_sheet=EXCLUDE_SHEET,...
                                   Exclude sheets with ID.
  -I, --sheet-pattern=STRING       Only include sheets named matching given
                                   pattern.
  -E, --exclude-sheet-pattern=STRING
                                   Exclude sheets named matching given pattern.
  -C, --columns=COLUMNS,...        Columns names added to CSV as a zero row.
  -m, --max-columns=UINT           Columns maximum that will be exported.
  -e, --skip-empty                 Skip empty lines.
  -t, --trim                       Remove all leading and trailing white space.
```

## List sheets

```sh
Usage: xlsx2csv list <source>

List sheets.

Arguments:
  <source>    Path to XLSX file.

Flags:
  -h, --help    Show context-sensitive help.
```
