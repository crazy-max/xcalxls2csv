package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/crazy-max/xcalxls2csv/pkg/xcal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	version = "dev"
	c       cli
)

type cli struct {
	Version kong.VersionFlag
	XcalXls string `kong:"name='xcalxls',arg,required,help='XCalibur XLS file to convert.'"`
	Output  string `kong:"name='output',arg,required,help='Output CSV filename.'"`
}

func main() {
	// parse command line
	kctx := kong.Parse(&c,
		kong.Name("xcalxls2csv"),
		kong.Description(`XCaliburXLS data frames to CSV`),
		kong.UsageOnError(),
		kong.Vars{
			"version": fmt.Sprintf("%s", version),
		},
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))

	// logger
	output := zerolog.ConsoleWriter{
		Out: os.Stdout,
	}
	output.FormatTimestamp = func(i interface{}) string {
		return kctx.Model.Name
	}
	log.Logger = zerolog.New(output).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// handle os signals
	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-channel
		log.Warn().Msgf("caught signal %v", sig)
		os.Exit(0)
	}()

	// start
	if dt, err := xcal.ConvertToCSV(c.XcalXls); err != nil {
		log.Fatal().Err(err).Msg("cannot convert XCalibur XLS file")
	} else if err := os.WriteFile(c.Output, dt, 0644); err != nil {
		log.Fatal().Err(err).Msg("cannot write output file")
	}
}
