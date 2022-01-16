package converter

import (
	"encoding/csv"
	"fmt"
	"io"
)

type CSV struct {
	w *csv.Writer
}

func NewCSV(w io.Writer, delimiter string, encoding string, lineTerminator string) (*CSV, error) {
	csv := csv.NewWriter(w)

	if delimiter != "" {
		r := []rune(delimiter)
		if len(r) > 1 {
			return nil, fmt.Errorf("delimiter must be a single character")
		}
		csv.Comma = r[0]
	}

	switch lineTerminator {
	case "\\r\\n", "\r\n", "CRLF", "crlf":
		csv.UseCRLF = true
	case "\\n", "\n", "LF", "lf":
		csv.UseCRLF = false
	case "":
		csv.UseCRLF = false
	default:
		return nil, fmt.Errorf(`line delimiter must be one of "\r\n", "CRLF", "crlf", "\n", "LF", "lf"`)
	}

	return &CSV{w: csv}, nil
}

func (c *CSV) Write(cells []string) error {
	if err := c.w.Write(cells); err != nil {
		return fmt.Errorf("writing cells to CSV: %w", err)
	}

	return nil
}

func (c *CSV) WriteAll(cells [][]string) error {
	if err := c.w.WriteAll(cells); err != nil {
		return fmt.Errorf("writing cells to CSV: %w", err)
	}

	return nil
}

func (c *CSV) Flush() error {
	c.w.Flush()
	if err := c.w.Error(); err != nil {
		return fmt.Errorf("flushing data: %w", err)
	}

	return nil
}
