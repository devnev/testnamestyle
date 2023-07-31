package main

import (
	"github.com/devnev/testnamestyle"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(testnamestyle.TestNameStyle)
}
