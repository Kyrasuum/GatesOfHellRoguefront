package input

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
)

var (
	// literal max
	//maxcontrollers                  = int32(math.MaxInt32 - 1)
	// realistic max
	maxcontrollers                  = int32(16)
	controllers    map[string]int32 = map[string]int32{}

	modifiers map[string]int32 = map[string]int32{
		// Keyboard Function Keys
		"LeftShift":    raylib.KeyLeftShift,
		"LeftControl":  raylib.KeyLeftControl,
		"LeftAlt":      raylib.KeyLeftAlt,
		"LeftSuper":    raylib.KeyLeftSuper,
		"RightShift":   raylib.KeyRightShift,
		"RightControl": raylib.KeyRightControl,
		"RightAlt":     raylib.KeyRightAlt,
		"RightSuper":   raylib.KeyRightSuper,
	}

	keyboard map[string]int32 = map[string]int32{
		// Keyboard Function Keys
		"Space":        raylib.KeySpace,
		"Escape":       raylib.KeyEscape,
		"Enter":        raylib.KeyEnter,
		"Tab":          raylib.KeyTab,
		"Backspace":    raylib.KeyBackspace,
		"Insert":       raylib.KeyInsert,
		"Delete":       raylib.KeyDelete,
		"Right":        raylib.KeyRight,
		"Left":         raylib.KeyLeft,
		"Down":         raylib.KeyDown,
		"Up":           raylib.KeyUp,
		"PageUp":       raylib.KeyPageUp,
		"PageDown":     raylib.KeyPageDown,
		"Home":         raylib.KeyHome,
		"End":          raylib.KeyEnd,
		"CapsLock":     raylib.KeyCapsLock,
		"ScrollLock":   raylib.KeyScrollLock,
		"NumLock":      raylib.KeyNumLock,
		"PrintScreen":  raylib.KeyPrintScreen,
		"Pause":        raylib.KeyPause,
		"F1":           raylib.KeyF1,
		"F2":           raylib.KeyF2,
		"F3":           raylib.KeyF3,
		"F4":           raylib.KeyF4,
		"F5":           raylib.KeyF5,
		"F6":           raylib.KeyF6,
		"F7":           raylib.KeyF7,
		"F8":           raylib.KeyF8,
		"F9":           raylib.KeyF9,
		"F10":          raylib.KeyF10,
		"F11":          raylib.KeyF11,
		"F12":          raylib.KeyF12,
		"KbMenu":       raylib.KeyKbMenu,
		"LeftBracket":  raylib.KeyLeftBracket,
		"BackSlash":    raylib.KeyBackSlash,
		"RightBracket": raylib.KeyRightBracket,
		"Grave":        raylib.KeyGrave,

		// Keyboard Number Pad Keys
		"Kp0":        raylib.KeyKp0,
		"Kp1":        raylib.KeyKp1,
		"Kp2":        raylib.KeyKp2,
		"Kp3":        raylib.KeyKp3,
		"Kp4":        raylib.KeyKp4,
		"Kp5":        raylib.KeyKp5,
		"Kp6":        raylib.KeyKp6,
		"Kp7":        raylib.KeyKp7,
		"Kp8":        raylib.KeyKp8,
		"Kp9":        raylib.KeyKp9,
		"KpDecimal":  raylib.KeyKpDecimal,
		"KpDivide":   raylib.KeyKpDivide,
		"KpMultiply": raylib.KeyKpMultiply,
		"KpSubtract": raylib.KeyKpSubtract,
		"KpAdd":      raylib.KeyKpAdd,
		"KpEnter":    raylib.KeyKpEnter,
		"KpEqual":    raylib.KeyKpEqual,

		// Keyboard Alpha Numeric Keys
		"Apostrophe": raylib.KeyApostrophe,
		"Comma":      raylib.KeyComma,
		"Minus":      raylib.KeyMinus,
		"Period":     raylib.KeyPeriod,
		"Slash":      raylib.KeySlash,
		"Zero":       raylib.KeyZero,
		"One":        raylib.KeyOne,
		"Two":        raylib.KeyTwo,
		"Three":      raylib.KeyThree,
		"Four":       raylib.KeyFour,
		"Five":       raylib.KeyFive,
		"Six":        raylib.KeySix,
		"Seven":      raylib.KeySeven,
		"Eight":      raylib.KeyEight,
		"Nine":       raylib.KeyNine,
		"Semicolon":  raylib.KeySemicolon,
		"Equal":      raylib.KeyEqual,
		"A":          raylib.KeyA,
		"B":          raylib.KeyB,
		"C":          raylib.KeyC,
		"D":          raylib.KeyD,
		"E":          raylib.KeyE,
		"F":          raylib.KeyF,
		"G":          raylib.KeyG,
		"H":          raylib.KeyH,
		"I":          raylib.KeyI,
		"J":          raylib.KeyJ,
		"K":          raylib.KeyK,
		"L":          raylib.KeyL,
		"M":          raylib.KeyM,
		"N":          raylib.KeyN,
		"O":          raylib.KeyO,
		"P":          raylib.KeyP,
		"Q":          raylib.KeyQ,
		"R":          raylib.KeyR,
		"S":          raylib.KeyS,
		"T":          raylib.KeyT,
		"U":          raylib.KeyU,
		"V":          raylib.KeyV,
		"W":          raylib.KeyW,
		"X":          raylib.KeyX,
		"Y":          raylib.KeyY,
		"Z":          raylib.KeyZ,

		// Android keys
		"Back":       raylib.KeyBack,
		"Menu":       raylib.KeyMenu,
		"VolumeUp":   raylib.KeyVolumeUp,
		"VolumeDown": raylib.KeyVolumeDown,
	}

	mouse map[string]int32 = map[string]int32{
		// Mouse Buttons
		"LeftButton":    int32(raylib.MouseButtonLeft),
		"RightButton":   int32(raylib.MouseButtonRight),
		"MiddleButton":  int32(raylib.MouseButtonMiddle),
		"SideButton":    int32(raylib.MouseButtonSide),
		"ExtraButton":   int32(raylib.MouseButtonExtra),
		"ForwardButton": int32(raylib.MouseButtonForward),
		"BackButton":    int32(raylib.MouseButtonBack),
	}

	gamepad map[string]int32 = map[string]int32{
		"Unknown":        raylib.GamepadButtonUnknown,
		"LeftFaceUp":     raylib.GamepadButtonLeftFaceUp,
		"LeftFaceRight":  raylib.GamepadButtonLeftFaceRight,
		"LeftFaceDown":   raylib.GamepadButtonLeftFaceDown,
		"LeftFaceLeft":   raylib.GamepadButtonLeftFaceLeft,
		"RightFaceUp":    raylib.GamepadButtonRightFaceUp,
		"RightFaceRight": raylib.GamepadButtonRightFaceRight,
		"RightFaceDown":  raylib.GamepadButtonRightFaceDown,
		"RightFaceLeft":  raylib.GamepadButtonRightFaceLeft,
		"LeftTrigger1":   raylib.GamepadButtonLeftTrigger1,
		"LeftTrigger2":   raylib.GamepadButtonLeftTrigger2,
		"RightTrigger1":  raylib.GamepadButtonRightTrigger1,
		"RightTrigger2":  raylib.GamepadButtonRightTrigger2,
		"MiddleLeft":     raylib.GamepadButtonMiddleLeft,
		"Middle":         raylib.GamepadButtonMiddle,
		"MiddleRight":    raylib.GamepadButtonMiddleRight,
		"LeftThumb":      raylib.GamepadButtonLeftThumb,
		"RightThumb":     raylib.GamepadButtonRightThumb,

		"AxisLeftX":        raylib.GamepadAxisLeftX,
		"AxisLeftY":        raylib.GamepadAxisLeftY,
		"AxisRightX":       raylib.GamepadAxisRightX,
		"AxisRightY":       raylib.GamepadAxisRightY,
		"AxisLeftTrigger":  raylib.GamepadAxisLeftTrigger,
		"AxisRightTrigger": raylib.GamepadAxisRightTrigger,
	}

	gesture map[string]raylib.Gestures = map[string]raylib.Gestures{
		"GestureTap":        raylib.GestureTap,
		"GestureDoubletap":  raylib.GestureDoubletap,
		"GestureHold":       raylib.GestureHold,
		"GestureDrag":       raylib.GestureDrag,
		"GestureSwipeRight": raylib.GestureSwipeRight,
		"GestureSwipeLeft":  raylib.GestureSwipeLeft,
		"GestureSwipeUp":    raylib.GestureSwipeUp,
		"GestureSwipeDown":  raylib.GestureSwipeDown,
		"GesturePinchIn":    raylib.GesturePinchIn,
		"GesturePinchOut":   raylib.GesturePinchOut,
	}
)
