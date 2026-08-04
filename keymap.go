package menu

import (
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strings"
)

var (
	ErrInvalidConfigKeymap   = errors.New("invalid keymap")
	ErrInvalidConfigSettings = errors.New("invalid settings")
)

func KeymapToggleAll() *Keymap {
	return NewKeymap().
		WithBind(KeyCtrlA).
		WithDesc("toggle-all").
		WithBuiltinAction(KeybindActionToggleAll).
		WithArgs(func(b *ArgsBuilder) *ArgsBuilder {
			return b.WithMultiSelection().WithHighlightLine()
		})
}

func KeymapTogglePreview() *Keymap {
	return NewKeymap().
		WithBind(KeyCtrlSlash).
		WithBuiltinAction(KeybindActionTogglePreview).
		WithDesc("toggle-preview")
}

type (
	KeybindAction string // (e.g., "toggle-preview", "toggle-all")
	Keybind       string // (e.g., "ctrl-a", "ctrl-e", etc...)
)

const (
	// General.
	KeybindActionAccept KeybindAction = "accept"
	KeybindActionAbort  KeybindAction = "abort"

	// Selection.
	KeybindActionToggle      KeybindAction = "toggle"
	KeybindActionToggleAll   KeybindAction = "toggle-all"
	KeybindActionSelect      KeybindAction = "select"
	KeybindActionSelectAll   KeybindAction = "select-all"
	KeybindActionDeselectAll KeybindAction = "deselect-all"

	// Navigation.
	KeybindActionUp           KeybindAction = "up"
	KeybindActionDown         KeybindAction = "down"
	KeybindActionFirst        KeybindAction = "first"
	KeybindActionLast         KeybindAction = "last"
	KeybindActionPageUp       KeybindAction = "page-up"
	KeybindActionPageDown     KeybindAction = "page-down"
	KeybindActionHalfPageUp   KeybindAction = "half-page-up"
	KeybindActionHalfPageDown KeybindAction = "half-page-down"

	// Preview navigation.
	KeybindActionPreviewUp           KeybindAction = "preview-up"
	KeybindActionPreviewDown         KeybindAction = "preview-down"
	KeybindActionPreviewPageUp       KeybindAction = "preview-page-up"
	KeybindActionPreviewPageDown     KeybindAction = "preview-page-down"
	KeybindActionPreviewHalfPageUp   KeybindAction = "preview-half-page-up"
	KeybindActionPreviewHalfPageDown KeybindAction = "preview-half-page-down"

	// Preview.
	KeybindActionTogglePreview       KeybindAction = "toggle-preview"
	KeybindActionChangePreviewWindow KeybindAction = "change-preview-window"
)

const (
	KeyEnter Keybind = "enter"
	KeyTab   Keybind = "tab"
	KeyBTab  Keybind = "btab"
	KeyEsc   Keybind = "esc"
	KeySpace Keybind = "space"

	KeyUp    Keybind = "up"
	KeyDown  Keybind = "down"
	KeyLeft  Keybind = "left"
	KeyRight Keybind = "right"
	KeyHome  Keybind = "home"
	KeyEnd   Keybind = "end"
	KeyPgUp  Keybind = "pgup"
	KeyPgDn  Keybind = "pgdn"

	KeyCtrlA     Keybind = "ctrl-a"
	KeyCtrlB     Keybind = "ctrl-b"
	KeyCtrlC     Keybind = "ctrl-c"
	KeyCtrlD     Keybind = "ctrl-d"
	KeyCtrlE     Keybind = "ctrl-e"
	KeyCtrlF     Keybind = "ctrl-f"
	KeyCtrlG     Keybind = "ctrl-g"
	KeyCtrlH     Keybind = "ctrl-h"
	KeyCtrlI     Keybind = "ctrl-i"
	KeyCtrlJ     Keybind = "ctrl-j"
	KeyCtrlK     Keybind = "ctrl-k"
	KeyCtrlL     Keybind = "ctrl-l"
	KeyCtrlM     Keybind = "ctrl-m"
	KeyCtrlN     Keybind = "ctrl-n"
	KeyCtrlO     Keybind = "ctrl-o"
	KeyCtrlP     Keybind = "ctrl-p"
	KeyCtrlQ     Keybind = "ctrl-q"
	KeyCtrlR     Keybind = "ctrl-r"
	KeyCtrlS     Keybind = "ctrl-s"
	KeyCtrlT     Keybind = "ctrl-t"
	KeyCtrlU     Keybind = "ctrl-u"
	KeyCtrlV     Keybind = "ctrl-v"
	KeyCtrlW     Keybind = "ctrl-w"
	KeyCtrlX     Keybind = "ctrl-x"
	KeyCtrlY     Keybind = "ctrl-y"
	KeyCtrlZ     Keybind = "ctrl-z"
	KeyCtrlSlash Keybind = "ctrl-/"

	KeyAltA Keybind = "alt-a"
	KeyAltB Keybind = "alt-b"
	KeyAltC Keybind = "alt-c"
	KeyAltD Keybind = "alt-d"
	KeyAltE Keybind = "alt-e"
	KeyAltF Keybind = "alt-f"
	KeyAltG Keybind = "alt-g"
	KeyAltH Keybind = "alt-h"
	KeyAltI Keybind = "alt-i"
	KeyAltJ Keybind = "alt-j"
	KeyAltK Keybind = "alt-k"
	KeyAltL Keybind = "alt-l"
	KeyAltM Keybind = "alt-m"
	KeyAltN Keybind = "alt-n"
	KeyAltO Keybind = "alt-o"
	KeyAltP Keybind = "alt-p"
	KeyAltQ Keybind = "alt-q"
	KeyAltR Keybind = "alt-r"
	KeyAltS Keybind = "alt-s"
	KeyAltT Keybind = "alt-t"
	KeyAltU Keybind = "alt-u"
	KeyAltV Keybind = "alt-v"
	KeyAltW Keybind = "alt-w"
	KeyAltX Keybind = "alt-x"
	KeyAltY Keybind = "alt-y"
	KeyAltZ Keybind = "alt-z"

	KeyF1  Keybind = "f1"
	KeyF2  Keybind = "f2"
	KeyF3  Keybind = "f3"
	KeyF4  Keybind = "f4"
	KeyF5  Keybind = "f5"
	KeyF6  Keybind = "f6"
	KeyF7  Keybind = "f7"
	KeyF8  Keybind = "f8"
	KeyF9  Keybind = "f9"
	KeyF10 Keybind = "f10"
	KeyF11 Keybind = "f11"
	KeyF12 Keybind = "f12"

	KeyShiftUp    Keybind = "shift-up"
	KeyShiftDown  Keybind = "shift-down"
	KeyShiftLeft  Keybind = "shift-left"
	KeyShiftRight Keybind = "shift-right"
)

// Keymap holds the keymap configuration.
type Keymap struct {
	Bind    Keybind       `json:"bind"           yaml:"bind"`           // keybind combination
	Action  KeybindAction `json:"-"              yaml:"-"`              // action to execute
	Desc    string        `json:"description"    yaml:"description"`    // keybind description
	Enabled bool          `json:"enabled"        yaml:"enabled"`        // keybind enabled
	Hidden  bool          `json:"hidden"         yaml:"hidden"`         // keybind hidden
	Args    Args          `json:"args,omitempty" yaml:"args,omitempty"` // arguments added to the menu
}

func NewKeymap() *Keymap {
	return &Keymap{Enabled: true}
}

func (k *Keymap) String() string {
	return fmt.Sprintf("%s:%s", k.Bind, k.Action)
}

func (k *Keymap) Hide() *Keymap {
	k.Hidden = true
	return k
}

func (k *Keymap) IsEnabled() bool {
	return k.Enabled
}

func (k *Keymap) Arguments() []string {
	return k.Args
}

func (k *Keymap) BindString() string {
	return string(k.Bind)
}

func (k *Keymap) WithBind(b Keybind) *Keymap {
	k.Bind = b
	return k
}

func (k *Keymap) WithDesc(s string) *Keymap {
	k.Desc = s
	return k
}

func (k *Keymap) WithEnabled(b bool) *Keymap {
	k.Enabled = b
	return k
}

// WithExecute returns a new Keymap with the given action command.
func (k *Keymap) WithExecute(cmd string) *Keymap {
	k.Action = KeybindAction(fmt.Sprintf("execute(%s)", cmd))
	return k
}

func (k *Keymap) WithBuiltinAction(cmd KeybindAction) *Keymap {
	k.Action = cmd
	return k
}

func (k *Keymap) WithCommand(cmd string) *Keymap {
	k.Action = KeybindAction(cmd)
	return k
}

func (k *Keymap) WithBecome(cmd string) *Keymap {
	k.Action = KeybindAction(fmt.Sprintf("become(%s)", cmd))
	return k
}

// WithSilentExecute returns a new Keymap with the given action command.
func (k *Keymap) WithSilentExecute(cmd string) *Keymap {
	k.Action = KeybindAction(fmt.Sprintf("execute-silent(%s)", cmd))
	return k
}

func (k *Keymap) WithArgs(fn func(b *ArgsBuilder) *ArgsBuilder) *Keymap {
	k.Args = fn(NewArgsBuilder()).Build()
	return k
}

type KeymapManager struct {
	keymaps map[KeybindAction]*Keymap // keymaps action:keymap
}

func NewKeymapManager() *KeymapManager {
	return &KeymapManager{
		keymaps: make(map[KeybindAction]*Keymap),
	}
}

func (km *KeymapManager) NewKeymap() *Keymap {
	return NewKeymap()
}

func (km *KeymapManager) Register(keys ...*Keymap) {
	for i := range keys {
		k := keys[i]
		km.keymaps[k.Action] = k
		slog.Debug("keybind register", "bind", k.Bind, "action", k.Action)
	}
}

// Remove removes bind from keymaps map.
func (km *KeymapManager) Remove(k *Keymap) {
	slog.Debug("keybind remove", "bind", k.Bind, "action", k.Action)
	delete(km.keymaps, k.Action)
}

func (km *KeymapManager) Len() int {
	return len(km.keymaps)
}

// List returns a sorted keymap slice.
func (km *KeymapManager) List() []*Keymap {
	// extract keys and sort them
	keys := make([]string, 0, len(km.keymaps))
	for k := range km.keymaps {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)

	// build sorted result
	result := make([]*Keymap, 0, len(keys))
	for _, k := range keys {
		result = append(result, km.keymaps[KeybindAction(k)])
	}

	return result
}

// Find returns the first Keymap that matches the given action or bind.
// It returns nil if no keymap is found.
func (km *KeymapManager) Find(bind *Keymap) *Keymap {
	for _, k := range km.keymaps {
		if strings.EqualFold(string(k.Action), string(bind.Action)) ||
			strings.EqualFold(string(k.Bind), string(bind.Bind)) {
			return k
		}
	}
	return nil
}
