package webp

/*
#cgo pkg-config: libwebp
#include <stdlib.h>
#include "webp_util.h"
*/
import "C"
import (
	"fmt"
)

type Config struct {
	Lossless bool
	Method   int
	Quality  float32
	Hint     Hint

	cConf *C.struct_WebPConfig
}

func (c *Config) toCStruct() (*C.struct_WebPConfig, error) {
	if c.cConf != nil {
		return c.cConf, nil
	}

	var lossless int
	if c.Lossless {
		lossless = 1
	}

	var config C.struct_WebPConfig

	i, err := C.WebPConfigPreset(&config, C.WEBP_PRESET_PHOTO, C.float(c.Quality))

	if err != nil {
		return nil, err
	}

	if i != 1 {
		return nil, fmt.Errorf("error occurred initializing config: %d", i)
	}

	config.lossless = C.int(lossless)
	config.quality = C.float(c.Quality)
	config.method = C.int(c.Method)

	c.cConf = &config
	return &config, nil
}
