package cmd

import (
	"flag"
	"io"
	"strconv"
)

type AutoincrementSubcommand struct {
	name    string
	seed    int
	prepend bool
}

func (sub *AutoincrementSubcommand) Name() string {
	return "autoincrement"
}
func (sub *AutoincrementSubcommand) Aliases() []string {
	return []string{"autoinc"}
}
func (sub *AutoincrementSubcommand) Description() string {
	return "Add a column of incrementing integers to a CSV."
}
func (sub *AutoincrementSubcommand) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&sub.name, "name", "ID", "Name of autoincrementing column")
	fs.IntVar(&sub.seed, "seed", 1, "Initial value of autoincrementing column")
	fs.BoolVar(&sub.prepend, "prepend", false, "Prepend the autoincrementing column (defaults to append)")
}

func (sub *AutoincrementSubcommand) Run(args []string) {
	inputCsvs := GetInputCsvsOrPanic(args, 1)
	outputCsv := NewOutputCsvFromInputCsv(inputCsvs[0])
	sub.RunAutoincrement(inputCsvs[0], outputCsv)
}

func (sub *AutoincrementSubcommand) RunAutoincrement(inputCsv *InputCsv, outputCsvWriter OutputCsvWriter) {
	AutoIncrement(inputCsv, outputCsvWriter, sub.name, sub.seed, sub.prepend)
	err := inputCsv.Close()
	if err != nil {
		ExitWithError(err)
	}
}

func AutoIncrement(inputCsv *InputCsv, outputCsvWriter OutputCsvWriter, name string, seed int, prepend bool) {
	// Read and write header.
	header, err := inputCsv.Read()
	if err != nil {
		ExitWithError(err)
	}
	numInputColumns := len(header)
	shellRow := make([]string, numInputColumns+1)
	if prepend {
		shellRow[0] = name
		for i, elem := range header {
			shellRow[i+1] = elem
		}
	} else {
		copy(shellRow, header)
		shellRow[numInputColumns] = name
	}
	outputCsvWriter.Write(shellRow)

	// Write rows with autoincrement.
	inc := seed
	for {
		row, err := inputCsv.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				ExitWithError(err)
			}
		}
		incStr := strconv.Itoa(inc)
		if prepend {
			shellRow[0] = incStr
			for i, elem := range row {
				shellRow[i+1] = elem
			}
		} else {
			copy(shellRow, row)
			shellRow[numInputColumns] = incStr
		}
		inc++
		outputCsvWriter.Write(shellRow)
	}
}
