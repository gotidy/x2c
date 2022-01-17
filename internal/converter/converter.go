package converter

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/htmlindex"
)

var (
	BOM              = []byte("\xEF\xBB\xBF")
	EncodingsWithBOM = []string{"utf-8", "utf-16be", "utf-16le"}
)

type Writer interface {
	Write(cells []string) error
	WriteAll(cells []string) error
	Flush() error
}

type Options struct {
	Encoding       string
	BOM            bool
	Delimiter      string
	LineTerminator string

	SheetName        []string
	SheetID          []int
	ExcludeSheetName []string
	ExcludeSheetID   []int

	SheetPattern        string
	ExcludeSheetPattern string

	Columns    []string
	MaxColumns int
	Trim       bool
	SkipEmpty  bool

	Source      string
	Destination string

	Output io.Writer
}

type Converter struct {
	file    *excelize.File
	options Options
}

func New(options Options) (*Converter, error) {
	f, err := excelize.OpenFile(options.Source)
	if err != nil {
		return nil, fmt.Errorf("opening excel file \"%s\" failed: %w", options.Source, err)
	}

	if options.Destination == "" && options.Output == nil {
		return nil, errors.New("output is not defined")
	}

	return &Converter{file: f, options: options}, nil
}

func (c *Converter) Close() error {
	err := c.file.Close()
	if err != nil {
		return fmt.Errorf("closing excel file \"%s\" failed: %w", c.file.Path, err)
	}

	return nil
}

func (c *Converter) Sheets() []string {
	return c.file.GetSheetList()
}

func (c *Converter) Convert() error {
	var include, exclude *regexp.Regexp
	var err error
	if c.options.SheetPattern != "" {
		include, err = regexp.Compile(c.options.SheetPattern)
		if err != nil {
			return fmt.Errorf("include sheet pattern is invalid (%s): %w", c.options.SheetPattern, err)
		}
	}
	if c.options.ExcludeSheetPattern != "" {
		exclude, err = regexp.Compile(c.options.ExcludeSheetPattern)
		if err != nil {
			return fmt.Errorf("exclude sheet pattern is invalid (%s): %w", c.options.ExcludeSheetPattern, err)
		}
	}

	var bom []byte
	if c.options.BOM {
		bom = BOM
	}

	var encoder encoding.Encoding
	if c.options.Encoding != "" {
		encoder, err = htmlindex.Get(c.options.Encoding)
		if err != nil {
			return fmt.Errorf("getting encoder for \"%s\": %w", c.options.Encoding, err)
		}

		if c.options.Output != nil {
			c.options.Output = encoder.NewEncoder().Writer(c.options.Output)
		}

		switch encName, err := htmlindex.Name(encoder); {
		case err != nil:
			return fmt.Errorf("getting encoder name \"%s\": %w", c.options.Encoding, err)
		case inStrings(encName, EncodingsWithBOM):
		default:
			bom = nil
		}
	}

	if c.options.Output != nil && c.options.Destination == "" && bom != nil {
		if _, err := c.options.Output.Write(bom); err != nil {
			return fmt.Errorf("writing BOM: %w", err)
		}
	}

	path := ""
	addColumns := true
	for i := 0; i < c.file.SheetCount; i++ {
		name := c.file.GetSheetName(i)
		w := c.options.Output
		if c.options.Destination != "" {
			repl := strings.NewReplacer("{{name}}", name, "{{num}}", strconv.Itoa(i))
			dest := repl.Replace(c.options.Destination)
			if dest != path {
				path = dest
				f, err := os.Create(path)
				if err != nil {
					return fmt.Errorf("opening file %s: %w", path, err)
				}
				defer f.Close()
				w = f
				if encoder != nil {
					w = encoder.NewEncoder().Writer(w)
				}
				if bom != nil {
					if _, err := w.Write(bom); err != nil {
						return fmt.Errorf("writing BOM: %w", err)
					}
				}
				addColumns = true
			}
		}

		writer, err := NewCSV(w, c.options.Delimiter, c.options.Encoding, c.options.LineTerminator)
		if err != nil {
			return fmt.Errorf("creating writer: %w", err)
		}

		if addColumns && len(c.options.Columns) != 0 {
			if err = writer.Write(c.options.Columns); err != nil {
				return fmt.Errorf("writing columns row: %w", err)
			}
		}

		switch {
		case inInts(i, c.options.ExcludeSheetID) || inStrings(name, c.options.ExcludeSheetName) || (exclude != nil && exclude.MatchString(name)):
			continue
		case len(c.options.SheetID) == 0 && len(c.options.SheetName) == 0 && c.options.SheetPattern == "":
		case inInts(i, c.options.SheetID) || inStrings(name, c.options.SheetName) || (include != nil && include.MatchString(name)):
		default:
			continue
		}

		rows, err := c.file.GetRows(name)
		if err != nil {
			return fmt.Errorf("getting rows of \"%s\" sheet: %w", name, err)
		}
		for _, row := range rows {
			if c.options.MaxColumns != 0 && len(row) > c.options.MaxColumns {
				row = row[:c.options.MaxColumns]
			}

			if c.options.Trim {
				trim(row)
			}

			if c.options.SkipEmpty && empty(row) {
				continue
			}

			if err = writer.Write(row); err != nil {
				return fmt.Errorf("writing row: %w", err)
			}
		}

		if err = writer.Flush(); err != nil {
			return fmt.Errorf("flushing rows[%d]: %w", len(rows), err)
		}
	}

	return nil
}

func inInts(i int, list []int) bool {
	for _, v := range list {
		if v == i {
			return true
		}
	}
	return false
}

func inStrings(i string, list []string) bool {
	for _, v := range list {
		if v == i {
			return true
		}
	}
	return false
}

func empty(s []string) bool {
	for _, s := range s {
		if s != "" {
			return false
		}
	}
	return true
}

func trim(s []string) {
	for i, v := range s {
		s[i] = strings.TrimSpace(v)
	}
}
