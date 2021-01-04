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

		debug, _ := cmd.Flags().GetBool(optDebug)

		recs, err := data.Load("archivio.tsv", id)
		if err != nil {
			return err
		}

		sort.Slice(recs, func(i, j int) bool {
			x, y := recs[i], recs[j]
			return (x.Day > y.Day)
		})

		sample := make([]data.Record, 2*cycle+1)
		copy(sample, recs)

		sort.Slice(sample, func(i, j int) bool {
			x, y := sample[i], sample[j]
			return (x.Day < y.Day)
		})

		// Contatori sortite nell' intervallo [mid;sup]
		counters := data.Count(sample[0:cycle], id, debug)
		// Conserva solo quelli usciti almeno `tot` volte
		heroes := []int{}
		for k, v := range counters {
			if v >= max {
				heroes = append(heroes, k)
			}
		}
		sort.Ints(heroes)
		if debug {
			fmt.Fprintf(os.Stderr, "numeri sortiti almeno %d volte (intevallo %d - %d): %v\n", max, 0, cycle, heroes)
		}

		if debug {
			fmt.Fprintln(os.Stderr)
		}

		// Contatori sortite nell' intervallo [inf;mid]
		counters = data.Count(sample[cycle:], id, debug)
		// Conserva solo quelli mai usciti
		ghosts := []int{}
		for k, v := range counters {
			if v <= min {
				ghosts = append(ghosts, k)
			}
		}
		sort.Ints(ghosts)
		if debug {
			fmt.Fprintf(os.Stderr, "numeri mai sortiti (intevallo %d - %d): %v\n", cycle+1, 2*cycle, ghosts)
		}

		if len(ghosts) == 0 {
			fmt.Fprintf(os.Stdout, "Non c'è alcun numero che soddisfi i requisiti\n")
			return nil
		}

		res := collect.Intersection(heroes, ghosts)
		if debug {
			fmt.Fprintf(os.Stderr, "\nintersezione risultati: %v\n", res)
		} else {
			fmt.Fprintf(os.Stdout, "%v\n", res)
		}

		return nil
	},
}

func init() {
	searchCmd.Flags().BoolP(optDebug, "d", false, "stampa le estrazioni esaminate")
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
