# wireframe
A simple 3D terrain visualizer written in Go. When developing this project, I quickly found out how mildly inconvenient it is to do any form of computer graphics programming in a language without support for classes, static fields/methods, operator overloading, and even implicit narrowing conversions. Nevertheless, I wanted to try out Golang, so I thought I'd rewrite something I'm familiar with even if it may not be the best use case for the language. I used [raylib](https://github.com/raysan5/raylib) and [raylib-go](https://github.com/gen2brain/raylib-go) for creating windows and drawing lines/vertices, and did the rest of the rendering on the cpu. 

[![Terrain](https://cdn.discordapp.com/attachments/735900498212814865/1527867620341321799/image.png?ex=6a5c391b&is=6a5ae79b&hm=7d01265ccfe65017794ce90b340beba827fbe83d5a10d7589b3746f1e7b89003)](https://cdn.discordapp.com/attachments/735900498212814865/1527867620341321799/image.png?ex=6a5c391b&is=6a5ae79b&hm=7d01265ccfe65017794ce90b340beba827fbe83d5a10d7589b3746f1e7b89003)

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