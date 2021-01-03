package data

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type Record struct {
	Day     int
	Numbers map[string][5]int
}

func List(src []Record, id string) {
	rows := [][]string{}

	for _, el := range src {
		rows = append(rows, asRow(el, id))
	}

	dump(os.Stdout, rows)
}

func asRow(el Record, id string) []string {
	for key, vals := range el.Numbers {
		if strings.EqualFold(key, id) {
			dst := make([]int, len(vals))
			copy(dst, []int{vals[0], vals[1], vals[2], vals[3], vals[4]})
			sort.Ints(dst)

			return []string{
				fmt.Sprintf("%d", el.Day), id,
				fmt.Sprintf("%d", dst[0]),
				fmt.Sprintf("%d", dst[1]),
				fmt.Sprintf("%d", dst[2]),
				fmt.Sprintf("%d", dst[3]),
				fmt.Sprintf("%d", dst[4]),
			}
		}
	}

	return []string{}
}

func dump(wri io.Writer, rows [][]string) {
	table := tablewriter.NewWriter(wri)
	table.SetBorder(false)
	table.SetRowLine(true)
	//table.SetAutoMergeCells(true)
	table.SetCaption(true, fmt.Sprintf("Nr. estrazioni: %d", len(rows)))
	table.AppendBulk(rows)
	table.Render()
}

func Find(src []Record, day int) int {
	for i, el := range src {
		if el.Day == day {
			return i
		}
	}

	return -1
}

func Load(filename string) ([]Record, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	return Parse(fd)
}

func Parse(r io.Reader) ([]Record, error) {
	records := []Record{}

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())

		el, err := ParseLine(fields)
		if err != nil {
			return nil, err
		}

		if idx := Find(records, el.Day); idx == -1 {
			records = append(records, Record{
				Day: el.Day,
				Numbers: map[string][5]int{
					el.Id: el.Numbers,
				},
			})
		} else {
			records[idx].Numbers[el.Id] = el.Numbers
		}

	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Day > records[j].Day
	})

	return records, nil
}

func Save(recs []Record, filename string, limit int) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	file.Truncate(0)

	wri := csv.NewWriter(file)
	wri.Comma = '\t'
	for i, el := range recs {
		for key, vals := range el.Numbers {
			wri.Write([]string{
				fmt.Sprintf("%d", el.Day),
				key,
				fmt.Sprintf("%d", vals[0]),
				fmt.Sprintf("%d", vals[1]),
				fmt.Sprintf("%d", vals[2]),
				fmt.Sprintf("%d", vals[3]),
				fmt.Sprintf("%d", vals[4]),
			})
		}

		wri.Flush()

		if i >= limit {
			break
		}
	}

	return nil
}

type Row struct {
	Day     int
	Id      string
	Numbers [5]int
}

func ParseLine(fields []string) (Row, error) {
	res := Row{}

	date := strings.ReplaceAll(fields[0], "/", "")
	day, err := strconv.Atoi(date)
	if err != nil {
		return res, err
	}

	n1, err := strconv.Atoi(fields[2])
	if err != nil {
		return res, err
	}

	n2, err := strconv.Atoi(fields[3])
	if err != nil {
		return res, err
	}

	n3, err := strconv.Atoi(fields[4])
	if err != nil {
		return res, err
	}

	n4, err := strconv.Atoi(fields[5])
	if err != nil {
		return res, err
	}

	n5, err := strconv.Atoi(fields[6])
	if err != nil {
		return res, err
	}

	res.Day = day
	res.Id = fields[1]
	res.Numbers = [5]int{n1, n2, n3, n4, n5}

	return res, nil
}
