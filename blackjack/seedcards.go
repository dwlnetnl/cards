// +build IGNORE

package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/dwlnetnl/cards/card"
)

func main() {
	min, max, err := parseInput()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	d := card.NewStandardDeck()
	tw := tabwriter.NewWriter(os.Stdout, 5, 0, 1, ' ', 0)

	for i := min; i <= max; i++ {
		s := card.NewSeededShuffler(d, 6, rand.NewSource(i))

		fmt.Fprintf(tw, "%d\t", i)
		for i := 0; i < 10; i++ {
			fmt.Fprintf(tw, "%v\t", s.MustDraw())
		}
		fmt.Fprintln(tw)

	}

	if err := tw.Flush(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func parseInput() (min, max int64, err error) {
	if len(os.Args) == 1 {
		err = errors.New("provide seed range (e.g. 34-219 or -219)")

	} else if strings.HasSuffix(os.Args[1], "help") {
		err = errors.New("generate shuffler output for a seed range (e.g. 34-219)")

	} else if os.Args[1][0] == '-' {
		max, err = strconv.ParseInt(os.Args[1][1:], 10, 64)

	} else {
		_, err = fmt.Sscanf(os.Args[1], "%d-%d", &min, &max)
		if err == io.EOF {
			err = errors.New("provide upper range bound")
		}
	}
	return
}
