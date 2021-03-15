package main

import (
	"github.com/d2r2/go-i2c"
	logger "github.com/d2r2/go-logger"
)

var lg = logger.NewPackageLogger("main",
	logger.DebugLevel,
)

func main() {
	defer logger.FinalizeLogger()
	i2c, err := i2c.NewI2C(0x76, 3)
	if err != nil {
		lg.Fatal(err)
	}
	defer i2c.Close()

	t := NewTSYS01()

	err = t.Init(i2c)
	if err != nil {
		lg.Fatal(err)
	}

	temp, err := t.Read(i2c)
	if err != nil {
		lg.Fatal(err)
	}

	logger.ChangePackageLogLevel("i2c", logger.InfoLevel)

	lg.Infof("Temp = %v C", temp)
}