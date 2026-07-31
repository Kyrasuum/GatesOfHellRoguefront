module roguefront

go 1.25.0

replace github.com/gen2brain/raylib-go/raylib => ./include/raylib-go/raylib

require (
	github.com/gen2brain/raylib-go/raygui v0.0.0-20260730054012-0e0ad64f1120
	github.com/gen2brain/raylib-go/raylib v0.56.0-dev.0.20260513185948-c427d7332954
	github.com/go-gl/gl v0.0.0-20231021071112-07e5d0ea2e71
	github.com/iancoleman/strcase v0.3.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/sqweek/dialog v0.0.0-20260123140253-64c163d53aac
	github.com/zyedidia/json5 v0.0.0-20200102012142-2da050b1a98d
)

require (
	github.com/TheTitanrain/w32 v0.0.0-20180517000239-4f5cfb03fabf // indirect
	github.com/ebitengine/purego v0.10.0 // indirect
	golang.org/x/exp v0.0.0-20260508232706-74f9aab9d74a // indirect
	golang.org/x/sys v0.20.0 // indirect
)
