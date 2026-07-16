# wireframe
A simple 3D terrain visualizer written in Go. You will quickly find out how mildly irritating it is to do any form of computer graphics programming in a language without rich support for classes, static fields and methods, and operator overloading. Nevertheless, I wanted to try out Golang, so I thought I'd rewrite something I'm familiar with even if it may not be the best use case for the language. Shoutout [raylib](https://github.com/raysan5/raylib) and [raylib-go](https://github.com/gen2brain/raylib-go). 

## How to build

1) Install cgo
```
apt-get install libgl1-mesa-dev libxi-dev libxcursor-dev libxrandr-dev libxinerama-dev libwayland-dev libxkbcommon-dev
```

2) Install raylib-go
```
go install -v github.com/gen2brain/raylib-go/raylib@latest
```

3) Build the project
```
go build .
```