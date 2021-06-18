## Description

The Saper Go is an implementation of [Minesweeper]() game written in
[golang](https://golang.org)

## How to start?

### Download

- by using `git`

```sh
git clone https://github.com/gucio321/saper-go
cd saper-go
go get -d ./...
```

to build, use `go build .`
to run - `go run cmd/game/game.go`

- by `go` itself:

```sh
go get github.com/gucio321/saper-go
```

### usage

project provides a giu widget for saper game
usage is very simple e.g.:

<details><summary>code</summary>

```golang
package main

import (
	"github.com/AllenDang/giu"

	game "github.com/gucio321/saper-go/pkg/sgiu"
)

func loop() {
	giu.SingleWindow("game").Layout(
		game.Create(16, 30, 99),
	)
}

func main() {
	wnd := giu.NewMasterWindow("minesweeper", 640, 480, 0)
	wnd.Run(loop)
}
```

</details>
