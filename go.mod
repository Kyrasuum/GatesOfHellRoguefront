module roguefront

go 1.22

replace github.com/gen2brain/raylib-go/raylib => ./include/raylib-go/raylib

require (
	github.com/gen2brain/raylib-go/raylib v0.0.0-20250215042252-db8e47f0e5c5
	github.com/go-gl/gl v0.0.0-20231021071112-07e5d0ea2e71
	github.com/mitchellh/go-homedir v1.1.0
	github.com/zyedidia/json5 v0.0.0-20200102012142-2da050b1a98d
)

require (
	github.com/ebitengine/purego v0.7.1 // indirect
	golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842 // indirect
	golang.org/x/sys v0.20.0 // indirect
)
