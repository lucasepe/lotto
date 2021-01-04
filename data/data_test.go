package data

import (
	"testing"
)

func TestParseLine(t *testing.T) {
	cases := []struct {
		line string
		want string
	}{
		{
			"2001/01/03	BA	26	59	60	67	17",
			"20010103 ............... .7........6.... ............... .............90 ......7........ ............... [26 59 60 67 17]",
		},
		{
			"2001/01/03	CA	13	83	16	41	68",
			"20010103 ............3.. 6.............. ..........1.... ............... .......8....... .......3....... [13 83 16 41 68]",
		},
		{
			"2001/01/03	FI	33	87	76	7	10",
			"20010103 ......7..0..... ............... ..3............ ............... ............... 6..........7... [33 87 76 7 10]",
		},
		{
			"2001/01/03	GE	5	38	3	65	72",
			"20010103 ..3.5.......... ............... .......8....... ............... ....5......2... ............... [5 38 3 65 72]",
		},
		{
			"2001/01/03	MI	62	30	7	61	52",
			"20010103 ......7........ ..............0 ............... ......2........ 12............. ............... [62 30 7 61 52]",
		},
		{
			"2001/01/03	NA	25	38	15	31	86",
			"20010103 ..............5 .........5..... 1......8....... ............... ............... ..........6.... [25 38 15 31 86]",
		},
		{
			"2001/01/03	PA	15	56	89	85	63",
			"20010103 ..............5 ............... ............... ..........6.... ..3............ .........5...9. [15 56 89 85 63]",
		},
		{
			"2001/01/03	RM	41	17	80	48	9",
			"20010103 ........9...... .7............. ..........1.... ..8............ ............... ....0.......... [41 17 80 48 9]",
		},
		{
			"2001/01/03	TO	42	16	12	60	20",
			"20010103 ...........2... 6...0.......... ...........2... ..............0 ............... ............... [42 16 12 60 20]",
		},
		{
			"2001/01/03	VE	51	84	76	42	90",
			"20010103 ............... ............... ...........2... .....1......... ............... 6.......4.....0 [51 84 76 42 90]",
		},
	}

	for _, tt := range cases {
		t.Run(tt.line, func(t *testing.T) {
			got, err := parseLine(tt.line)
			if err != nil {
				t.Fatal(err)
			}

			if got.String() != tt.want {
				t.Errorf("got [%v] want [%v]", got, tt.want)
			}
		})
	}
}