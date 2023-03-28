package main

import (
	"errors"
	"flag"
	"fmt"
	_ "image/png"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

type PaddingType int

const (
	Pixel PaddingType = iota
	Percent
)

type PaddingArg struct {
	Value int
	Type  PaddingType
}

type PaddingArgs struct {
	Top    PaddingArg
	Right  PaddingArg
	Bottom PaddingArg
	Left   PaddingArg
}

type Args struct {
	Padding   PaddingArgs
	Tolerance int
	Radius    int
	// Indicates whether the program should read from stdin
	ReadStdin bool
	// If ReadString is false, these are the paths to the input files
	Files            []string
	OutFolder        string
	StopAtFirstError bool
	Force            bool
}

func parsePixelOrPercent(s string) (PaddingArg, error) {
	if strings.HasSuffix(s, "%") {
		v, err := strconv.Atoi(s[:len(s)-1])
		if err != nil {
			return PaddingArg{}, err
		}
		return PaddingArg{v, Percent}, err
	}

	v, err := strconv.Atoi(s)
	if err != nil {
		return PaddingArg{}, err
	}
	return PaddingArg{v, Pixel}, nil
}

func crash(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func main() {
	padding := flag.String("p", "", "Padding to add to all 4 sides of the image. Either an amount of pixels or a percentage relative to the output image")
	paddingX := flag.String("px", "", "Padding to add to the left and right sides of the image. Either an amount of pixels or a percentage relative to the output image")
	paddingY := flag.String("py", "", "Padding to add to the top and bottom sides of the image. Either an amount of pixels or a percentage relative to the output image")
	paddingT := flag.String("pt", "", "Padding to add to the top sides of the image. Either an amount of pixels or a percentage relative to the output image")
	paddingR := flag.String("pr", "", "Padding to add to the right sides of the image. Either an amount of pixels or a percentage relative to the output image")
	paddingB := flag.String("pb", "", "Padding to add to the bottom sides of the image. Either an amount of pixels or a percentage relative to the output image")
	paddingL := flag.String("pl", "", "Padding to add to the left sides of the image. Either an amount of pixels or a percentage relative to the output image")

	tolerance := flag.Int("t", 0, "Tolerance for detecting transparent pixels. 0-255, 0 being exact and 255 being anything")
	radius := flag.Int("r", 0, "Radius of non trasparent pixels around the current pixel for it to be considered an edge. Must be >= 0")

	outFolder := flag.String("o", "./centered", "Output folder. If not specified, the output files will be written to a \"centered\" folder in the current working directory.")

	stopAtFirstError := flag.Bool("s", false, "Stop at the first error encountered. If not specified, the program will continue processing the rest of the files.")

	force := flag.Bool("f", false, "Force overwrite of existing files")

	flag.Parse()

	paddingArgs := PaddingArgs{}

	if *padding != "" {
		p, err := parsePixelOrPercent(*padding)
		if err != nil {
			crash(err)
		}
		paddingArgs.Top = p
		paddingArgs.Right = p
		paddingArgs.Bottom = p
		paddingArgs.Left = p
	}

	if *paddingX != "" {
		p, err := parsePixelOrPercent(*paddingX)
		if err != nil {
			crash(err)
		}
		paddingArgs.Right = p
		paddingArgs.Left = p
	}

	if *paddingY != "" {
		p, err := parsePixelOrPercent(*paddingY)
		if err != nil {
			crash(err)
		}
		paddingArgs.Top = p
		paddingArgs.Bottom = p
	}

	if *paddingT != "" {
		p, err := parsePixelOrPercent(*paddingT)
		if err != nil {
			crash(err)
		}
		paddingArgs.Top = p
	}

	if *paddingR != "" {
		p, err := parsePixelOrPercent(*paddingR)
		if err != nil {
			crash(err)
		}
		paddingArgs.Right = p
	}

	if *paddingB != "" {
		p, err := parsePixelOrPercent(*paddingB)
		if err != nil {
			crash(err)
		}
		paddingArgs.Bottom = p
	}

	if *paddingL != "" {
		p, err := parsePixelOrPercent(*paddingL)
		if err != nil {
			crash(err)
		}
		paddingArgs.Left = p
	}

	args := Args{
		Padding:          paddingArgs,
		Tolerance:        *tolerance,
		Radius:           *radius,
		ReadStdin:        len(flag.Args()) == 0,
		Files:            flag.Args(),
		OutFolder:        *outFolder,
		StopAtFirstError: *stopAtFirstError,
		Force:            *force,
	}

	if args.ReadStdin {
		crash(errors.New("Empty input. Please specify at least one file to process."))
	}

	execute(args)
}

func execute(args Args) {

	err := os.MkdirAll(args.OutFolder, 0755)
	if err != nil {
		crash(err)
	}

	for _, file := range args.Files {

		outFile := path.Join(args.OutFolder, file)

		if _, err := os.Stat(outFile); err == nil {
			if !args.Force {
				if args.StopAtFirstError {
					crash(errors.New(fmt.Sprintf("File %s already exists.", outFile)))
				} else {
					log.Println("WARN: skipping", file, "because", outFile, "already exists.")
					continue
				}
			}
		} else if err != nil && !os.IsNotExist(err) {
			if args.StopAtFirstError {
				crash(err)
			} else {
				log.Println("WARN: skipping", file, "because", err)
				continue
			}
		}

		if stat, err := os.Stat(file); err == nil && stat.IsDir() {
			if args.StopAtFirstError {
				crash(errors.New(fmt.Sprintf("File %s is a directory.", outFile)))
			} else {
				log.Println("WARN: skipping", file, "because", outFile, "is a directory.")
				continue
			}
		} else if err != nil {
			if args.StopAtFirstError {
				crash(err)
			} else {
				log.Println("WARN: skipping", file, "because", err)
				continue
			}
		}

		err = Process(file, outFile, uint8(args.Tolerance), args.Padding)
		if err != nil {
			if args.StopAtFirstError {
				crash(fmt.Errorf("ERROR while processing %s: %w", file, err))
			} else {
				fmt.Fprintln(os.Stderr, err)
				log.Printf("WARN: encountered error while processing file %s: %s", file, err.Error())
				continue
			}
		}

		log.Printf("PASS: processed file %s -> %s", file, outFile)
	}

}
