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
)

// Record defines a specific draw
type Record struct {
	Day     int
	Wheel   string
	Numbers []int
}

func (r Record) String() string {
	nums := make([]int, len(r.Numbers))
	copy(nums, r.Numbers)
	sort.Ints(nums)

	s := strings.Repeat(".", 15)
	columns := []string{s, s, s, s, s, s}

	for _, n := range nums {
		x := (n + 14) - ((n + 14) % 15)

		col := (x - 15) / 15
		//fmt.Printf("n = %d, col = %d, row := %d\n", v, col, row)

		tmp := strings.Split(columns[col], "")

		steps := 15
		offset := col * steps
		for i := offset + 1; i <= offset+steps; i++ {
			if i == n {
				if n < 10 {
					tmp[i-offset-1] = fmt.Sprintf("%d", n)
				} else {
					tmp[i-offset-1] = fmt.Sprintf("%d", (n % 10))
				}
			}
		}

		columns[col] = strings.Join(tmp, "")
	}

	return fmt.Sprintf("%d %s %v", r.Day, strings.Join(columns, " "), r.Numbers)
}

func Count(recs []Record, id string, debug bool) map[int]int {
	res := map[int]int{}
	for i := 1; i <= 90; i++ {
		res[i] = 0
	}

	if debug {
		fmt.Fprintln(os.Stderr, "         000000000111111 111122222222223 333333333444444 444455555555556 666666666777777 777788888888889")
		fmt.Fprintf(os.Stderr, "  %s     123456789012345 678901234567890 123456789012345 678901234567890 123456789012345 678901234567890\n", id)
	}

	for _, el := range recs {
		if strings.EqualFold(el.Wheel, id) {
			if debug {
				fmt.Fprintf(os.Stderr, "%s\n", el)
			}
			for _, n := range el.Numbers {
				res[n] = res[n] + 1
			}
		}
	}

	return res
}

func Load(filename string, id string) ([]Record, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	records, err := Parse(fd)
	if err != nil {
		return nil, err
	}

	res := []Record{}
	for _, rec := range records {
		if strings.EqualFold(rec.Wheel, id) {
			res = append(res, rec)
		}
	}

	return res, nil
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
		if i > limit-1 {
			break
		}

		wri.Write([]string{
			fmt.Sprintf("%d", el.Day),
			el.Wheel,
			fmt.Sprintf("%d", el.Numbers[0]),
			fmt.Sprintf("%d", el.Numbers[1]),
			fmt.Sprintf("%d", el.Numbers[2]),
			fmt.Sprintf("%d", el.Numbers[3]),
			fmt.Sprintf("%d", el.Numbers[4]),
		})

		wri.Flush()
	}

	return nil
}

// Parse a draws data file
func Parse(r io.Reader) ([]Record, error) {
	records := []Record{}

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {

		el, err := parseLine(scanner.Text())
		if err != nil {
			return nil, err
		}

		records = append(records, el)
	}

	return records, nil
}

func parseLine(line string) (Record, error) {
	fields := strings.Fields(line)

	res := Record{}

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
	res.Wheel = fields[1]
	res.Numbers = []int{n1, n2, n3, n4, n5}

	return res, nil
}
