snek: req build
	./main 100 5

req:
	go get github.com/hajimehoshi/ebiten/v2

build:
	go build main.go

clean:
	rm main
