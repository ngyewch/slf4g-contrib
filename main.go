package main

import (
	"fmt"
	slf4g "github.com/echocat/slf4g"
	"github.com/echocat/slf4g/level"
	"github.com/echocat/slf4g/native"
	"github.com/echocat/slf4g/native/color"
	"github.com/echocat/slf4g/native/consumer"
	"github.com/echocat/slf4g/native/formatter"
	consumer2 "github.com/ngyewch/slf4g-contrib/native/consumer"
)

func main() {
	native.DefaultProvider.SetLevel(level.Debug)
	native.DefaultProvider.Consumer = consumer2.NewMultiConsumer(consumer.Default, consumer2.NewSentryConsumer(level.Error))
	formatter.Default = formatter.NewText(func(v *formatter.Text) {
		v.ColorMode = color.ModeNever
	})

	logger := slf4g.GetLoggerForCurrentPackage()
	logger.Tracef("trace1")
	logger.With("boo", 123).WithError(fmt.Errorf("myErr")).Debugf("debug1")
	logger.Infof("info1 %s", "hoo")
	logger.With("someValue", 123).WithError(fmt.Errorf("zimbu1")).Errorf("error1 %s", "zoo")
	logger.With("someValue", 456).WithError(fmt.Errorf("zimbu2")).Fatalf("fatal1 %s", "gjoac")
}
