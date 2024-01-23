package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

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
	XcalXls string `kong:"name='xcalxls',arg,required,help='Xcalibur XLS file to convert.'"`
	Output  string `kong:"name='output',help='Custom output filename.'"`
}

func main() {
	// parse command line
	kctx := kong.Parse(&c,
		kong.Name("xcalxls2csv"),
		kong.Description(`Xcalibur XLS data frames to CSV`),
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
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt, SIGTERM)
	go func() {
		sig := <-channel
		log.Warn().Msgf("Caught signal %v", sig)
		os.Exit(0)
	}()

	// start
	_, err := os.Stat(c.XcalXls)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot stat Xcalibur XLS file")
	}

	if c.Output == "" {
		c.Output = fmt.Sprintf("%s.csv", strings.TrimSuffix(c.XcalXls, filepath.Ext(c.XcalXls)))
	}

	log.Info().Msgf("Converting Xcalibur XLS file %s to %s", c.XcalXls, c.Output)
	if dt, err := xcal.ConvertToCSV(c.XcalXls); err != nil {
		log.Fatal().Err(err).Msg("Cannot convert Xcalibur XLS file")
	} else if err := os.WriteFile(c.Output, dt, 0644); err != nil {
		log.Fatal().Err(err).Msg("Cannot write output file")
	} else {
		log.Info().Msgf("Xcalibur XLS file converted successfully to %s", c.Output)
	}
}
