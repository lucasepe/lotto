package cmd

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lucasepe/lotto/data"
	"github.com/spf13/cobra"
)

const (
	lottomaticaURL = "https://www.lottomaticaitalia.it/STORICO_ESTRAZIONI_LOTTO/storico01-oggi.zip"

	cycle        = 18
	wheels       = 11
	defaultLimit = 3 * cycle * wheels

	optLimit = "limite"
)

// archiveCmd represents the archive command
var archiveCmd = &cobra.Command{
	DisableSuggestions:    true,
	DisableFlagsInUseLine: true,
	Args:                  cobra.NoArgs,
	Use:                   "aggiorna [--limite / -l]",
	Short:                 "Aggiorna l'archivio delle estrazioni",
	Example:               archiveCmdExample(),
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt(optLimit)

		buf, err := fetchFromURI(lottomaticaURL, -1)
		if err != nil {
			return err
		}

		recs, err := unzipAndParse(buf)
		if err != nil {
			return err
		}

		sort.Slice(recs, func(i, j int) bool {
			x, y := recs[i], recs[j]
			return (x.Day > y.Day) //&& (x.Wheel < y.Wheel)
		})

		return data.Save(recs, "archivio.tsv", limit)
	},
}

func init() {
	archiveCmd.Flags().IntP(optLimit, "l", defaultLimit, "numero massimo di estrazioni da archiviare")

	rootCmd.AddCommand(archiveCmd)
}

func archiveCmdExample() string {
	tpl := `{{APP}} aggiorna
  {{APP}} aggiorna -l 100`

	return strings.Replace(tpl, "{{APP}}", appName(), -1)
}

// fetchFromURI fetch data (with limit) from an HTTP URL.
// if 'limit' is greater then zero, fetch stops
// with EOF after 'limit' bytes.
func fetchFromURI(uri string, limit int64) ([]byte, error) {
	res, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if limit > 0 {
		return ioutil.ReadAll(io.LimitReader(res.Body, limit))
	}

	return ioutil.ReadAll(res.Body)
}

func unzipAndParse(buf []byte) ([]data.Record, error) {
	r, err := zip.NewReader(bytes.NewReader(buf), int64(len(buf)))
	if err != nil {
		return nil, err
	}

	if len(r.File) == 0 {
		return nil, fmt.Errorf("l'archivio non contiene dati")
	}

	fd, err := r.File[0].Open()
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	return data.Parse(fd)
}

func dateToInt(s string) (int, error) {
	date := strings.ReplaceAll(s, "/", "")
	return strconv.Atoi(date)
}

func defaultThreshold() (int, error) {
	ct := time.Now()
	nt := ct.AddDate(0, -8, 0)
	s := nt.Format("20060102")
	return dateToInt(s)
}
