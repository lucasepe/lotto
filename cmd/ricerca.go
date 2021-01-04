package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/lucasepe/lotto/collect"
	"github.com/lucasepe/lotto/data"
	"github.com/spf13/cobra"
)

const (
	optID  = "ruota"
	optMin = "min"
	optMax = "max"

	optLo  = "inf"
	optHi  = "sup"
	optMid = "mid"

	optDebug = "debug"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	DisableSuggestions:    true,
	DisableFlagsInUseLine: true,
	Args:                  cobra.NoArgs,
	Use:                   "ricerca [--limite / -l]",
	Short:                 "Aggiorna l'archivio delle estrazioni",
	Example:               searchCmdExample(),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetString(optID)
		max, _ := cmd.Flags().GetInt(optMax)
		min, _ := cmd.Flags().GetInt(optMin)
		lo, _ := cmd.Flags().GetInt(optLo)
		hi, _ := cmd.Flags().GetInt(optHi)
		mid, _ := cmd.Flags().GetInt(optMid)
		debug, _ := cmd.Flags().GetBool(optDebug)

		recs, err := data.Load("archivio.tsv")
		if err != nil {
			return err
		}

		if n := len(recs); n < hi {
			return fmt.Errorf("il numero massimo di estrazioni archiviate è %d", n)
		}

		sample := recs[lo:hi]

		// Contatori sortite nell' intervallo [mid;sup]
		counters := collect.Count(sample[mid:], id, debug)
		// Conserva solo quelli usciti almeno `tot` volte
		heroes := []int{}
		for k, v := range counters {
			if v >= max {
				heroes = append(heroes, k)
			}
		}
		sort.Ints(heroes)
		if debug {
			fmt.Fprintf(os.Stderr, "numeri sortiti almeno %d volte (intevallo %d - %d): %v\n", max, mid, len(sample), heroes)
		}

		// Contatori sortite nell' intervallo [inf;mid]
		counters = collect.Count(sample[:mid], id, debug)
		// Conserva solo quelli mai usciti
		ghosts := []int{}
		for k, v := range counters {
			if v <= min {
				ghosts = append(ghosts, k)
			}
		}
		sort.Ints(ghosts)
		if debug {
			fmt.Fprintf(os.Stderr, "numeri con sortite minore o uguale a %d (intevallo %d - %d): %v\n", min, lo, mid-1, ghosts)
		}

		if len(ghosts) == 0 {
			fmt.Fprintf(os.Stdout, "Non c'è alcun numero che soddisfi i requisiti\n")
			return nil
		}

		res := collect.Intersection(heroes, ghosts)
		if debug {
			fmt.Fprintf(os.Stderr, "intersezione risultati: %v\n", res)
		} else {
			fmt.Fprintf(os.Stdout, "%v\n", res)
		}

		return nil
	},
}

func init() {
	searchCmd.Flags().BoolP(optDebug, "d", false, "stampa anche le estrazioni esaminate")
	searchCmd.Flags().IntP(optLo, "i", 0, "limite inferiore intervallo estrazioni")
	searchCmd.Flags().IntP(optHi, "s", 36, "limite superiore intervallo estrazioni")
	searchCmd.Flags().IntP(optMid, "m", 19, "limite separazione intervallo estrazioni")
	searchCmd.Flags().IntP(optMax, "z", 3, "numero massimo di sortite")
	searchCmd.Flags().IntP(optMin, "a", 0, "numero minimo di sortite")
	searchCmd.Flags().StringP(optID, "r", "", "ruota da esaminare")
	searchCmd.MarkFlagRequired(optID)

	rootCmd.AddCommand(searchCmd)
}

func searchCmdExample() string {
	tpl := `{{APP}} ricerca
  {{APP}} ricerca -l 100`

	return strings.Replace(tpl, "{{APP}}", appName(), -1)
}
