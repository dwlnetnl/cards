// +build IGNORE

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"

	"github.com/dwlnetnl/cards/card"
)

func main() {
	min, max, err := parseInput()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	d := card.NewStandardDeck()
	buf := new(bytes.Buffer)

	for i := min; i <= max; i++ {
		s := card.NewSeededShuffler(d, 6, rand.NewSource(i))

		fmt.Fprintf(buf, "%d %v", i, s.MustDraw())
		for i := 0; i < 10; i++ {
			fmt.Fprintf(buf, ", %v", s.MustDraw())
		}
		fmt.Fprintln(buf)

		if _, err := buf.WriteTo(os.Stdout); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func parseInput() (min, max int64, err error) {
	if len(os.Args) == 1 {
		err = errors.New("provide seed range (e.g. 34-219 or -219)")

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
