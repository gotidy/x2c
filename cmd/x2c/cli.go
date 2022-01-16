package main

import (
	"fmt"
	"os"

	"github.com/gotidy/x2c/internal/converter"
	"github.com/rs/zerolog"
)

type Context struct {
	Logger zerolog.Logger
}

type ConvertCmd struct {
	Encoding       string `short:"c" name:"encoding" help:"Encoding of output CSV (http://www.w3.org/TR/encoding)." default:"utf-8"`
	BOM            bool   `short:"b" name:"bom" help:"Write BOM to CSV."`
	Delimiter      string `short:"d" name:"delimiter" help:"Cells delimiter." default:","`
	LineTerminator string `short:"l" name:"line_terminator" help:"line terminator - lines terminator in CSV. Default is \n" default:"\n" enum:"\n,\r\n,CRLF,LF"`

	SheetName        []string `short:"n" name:"sheetname" help:"Sheet names to convert."`
	SheetID          []int    `short:"s" name:"sheet" help:"Sheet indexes to convert."`
	ExcludeSheetName []string `short:"N" name:"exclude_sheet_name" help:"Exclude sheets with names."`
	ExcludeSheetID   []int    `short:"S" name:"exclude_sheet" help:"Exclude sheets with ID."`

	SheetPattern        string `short:"I" help:"Only include sheets named matching given pattern."`
	ExcludeSheetPattern string `short:"E" help:"Exclude sheets named matching given pattern."`

	Columns    []string `short:"C" name:"columns" help:"Columns names added to CSV as a zero row."`
	MaxColumns uint     `short:"m" help:"Columns maximum that will be exported."`
	SkipEmpty  bool     `short:"e" help:"Skip empty lines."`
	Trim       bool     `short:"t" help:"Remove all leading and trailing white space."`

	Source string `arg:"" name:"source" help:"Paths to XLSX file." type:"path"`
	Output string `arg:"" optional:"" name:"output" help:"Output CSV file path. It may include variables {{name}} and {{num}}" type:"path"`
}

func (c *ConvertCmd) Run(ctx *Context) error {
	conv, err := converter.New(converter.Options{
		Source:              c.Source,
		Output:              os.Stdout,
		Destination:         c.Output,
		Encoding:            c.Encoding,
		BOM:                 c.BOM,
		Delimiter:           c.Delimiter,
		LineTerminator:      c.LineTerminator,
		SheetName:           c.SheetName,
		SheetID:             c.SheetID,
		ExcludeSheetName:    c.ExcludeSheetName,
		ExcludeSheetID:      c.ExcludeSheetID,
		SheetPattern:        c.SheetPattern,
		ExcludeSheetPattern: c.ExcludeSheetPattern,
		MaxColumns:          int(c.MaxColumns),
		Trim:                c.Trim,
		SkipEmpty:           c.SkipEmpty,
	})
	if err != nil {
		ctx.Logger.Fatal().Msg(err.Error())
	}

	err = conv.Convert()
	if err != nil {
		ctx.Logger.Fatal().Msg(err.Error())
	}

	return nil
}

type ListCmd struct {
	Source string `arg:"" required:"" name:"source" help:"Path to XLSX file." type:"path"`
}

func (l *ListCmd) Run(ctx *Context) error {
	c, err := converter.New(converter.Options{Source: l.Source, Output: os.Stdout})
	if err != nil {
		ctx.Logger.Fatal().Msg(err.Error())
	}
	defer c.Close()

	for _, l := range c.Sheets() {
		fmt.Println(l)
	}
	return nil
}

type VersionCmd struct {
}

func (l *VersionCmd) Run(ctx *Context) error {
	fmt.Printf("Version %s\n\n", Version)
	fmt.Println(description)
	return nil
}

var cli struct {
	Convert ConvertCmd `cmd:"" help:"Convert XLSX file CSV." default:"withargs"`
	List    ListCmd    `cmd:"" help:"List sheets of XLSX file."`
	Version VersionCmd `cmd:"" help:"Version."`
}
