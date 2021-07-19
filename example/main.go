package main

import (
	"flag"

	"github.com/shotis/webp-go"
)

var (
	InputFlag    = flag.String("input", "example.png", "Path to the input file")
	OutputFlag   = flag.String("output", "example.webp", "Path to the output file")
	DecodeFlag   = flag.Bool("decode", false, "If this flag is set this will decode the input file. Input file should be WEBP")
	QualityFlag  = flag.Float64("quality", 80.0, "Sets the quality to be used in the config")
	MethodFlag   = flag.Int("method", 0, "Sets the method to be used. 0 is the fastest, 6 is the slowest. Higher = slower, but better")
	LosslessFlag = flag.Bool("lossless", false, "Sets the encoding to lossless compression")
	HintFlag     = flag.String("hint", "picture", "Sets the type of image hint to use for encoding.")

	config *webp.Config
)

func init() {
	// parse our flags that we defined
	flag.Parse()

	// construct the flag
	config = &webp.Config{
		Quality:  float32(*QualityFlag),
		Method:   *MethodFlag,
		Lossless: *LosslessFlag,
		Hint:     webp.GetHint(*HintFlag),
	}
}

func main() {
	if !*DecodeFlag {
		if err := encode(*InputFlag, *OutputFlag, config); err != nil {
			panic(err)
		}
	}
}
