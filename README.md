# wireframe
A simple 3D terrain visualizer written in Go. You will quickly find out how mildly inconvenient it is to do any form of computer graphics programming in a language without support for classes, static fields and methods, operator overloading, and even implicit narrowing conversions. Nevertheless, I wanted to try out Golang, so I thought I'd rewrite something I'm familiar with even if it may not be the best use case for the language. I used [raylib](https://github.com/raysan5/raylib) and [raylib-go](https://github.com/gen2brain/raylib-go) for creating windows and drawing lines/vertices, and did the rest of the rendering on the cpu. 

https://cdn.discordapp.com/attachments/735900498212814865/1527866875076280431/image.png?ex=6a5c3869&is=6a5ae6e9&hm=2be9c9beb6aaadc28808de9df79b4165c0cfdcd3dfd5603f24250193e6f97eef

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