package helpers

import (
	"github.com/spf13/pflag"
)

type flags struct {
	Version *pflag.Flag
}

var flagVar *flags

func setFlags() *flags {
	var flagVersion string

	pflag.StringVarP(&flagVersion, VERSION_FLAG_FULL_NAME, VERSION_FLAG_SHORT_NAME, VERSION_FLAG_DEAFULT_VALUE, VERSION_FLAG_USAGE)
	pflag.Lookup("version").NoOptDefVal = VERSION_FLAG_DEAFULT_VALUE

	pflag.Parse()
	flagVar = &flags{}

	flagVar.Version = pflag.Lookup(VERSION_FLAG_FULL_NAME)
	return flagVar
}

func GetFlags() *flags {
	if flagVar != nil {
		return flagVar
	}

	return setFlags()
}
