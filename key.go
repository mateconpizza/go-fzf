//go:build ignore

package menu

import "strings"

type Modifier uint8

const (
	ModCtrl Modifier = 1 << iota
	ModShift
	ModAlt
	ModMeta
)

type Key string

const (
	KeyA Key = "a"
	KeyE Key = "e"
	KeyK Key = "k"
	KeyL Key = "l"
	KeyO Key = "o"
	KeyP Key = "p"
	KeyW Key = "w"
	KeyY Key = "y"

	KeyEnterr Key = "enter"
	KeySlash  Key = "/"
)

var (
	KeyEnter = Shortcut{Key: KeyEnterr}

	KeyCtrlA     = Shortcut{Mods: ModCtrl, Key: KeyA}
	KeyCtrlE     = Shortcut{Mods: ModCtrl, Key: KeyE}
	KeyCtrlK     = Shortcut{Mods: ModCtrl, Key: KeyK}
	KeyCtrlL     = Shortcut{Mods: ModCtrl, Key: KeyL}
	KeyCtrlO     = Shortcut{Mods: ModCtrl, Key: KeyO}
	KeyCtrlW     = Shortcut{Mods: ModCtrl, Key: KeyW}
	KeyCtrlY     = Shortcut{Mods: ModCtrl, Key: KeyY}
	KeyCtrlSlash = Shortcut{Mods: ModCtrl, Key: KeySlash}
)

type Shortcut struct {
	Mods Modifier
	Key  Key
}

func (s Shortcut) String() string {
	if s.Mods == 0 {
		return string(s.Key)
	}
	return formatMods(s.Mods) + "-" + string(s.Key)
}

func formatMods(m Modifier) string {
	var parts []string
	if m&ModCtrl != 0 {
		parts = append(parts, "ctrl")
	}
	if m&ModShift != 0 {
		parts = append(parts, "shift")
	}
	if m&ModAlt != 0 {
		parts = append(parts, "alt")
	}
	if m&ModMeta != 0 {
		parts = append(parts, "meta")
	}
	return strings.Join(parts, "-")
}
