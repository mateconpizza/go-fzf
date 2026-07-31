package menu

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	shellwords "github.com/junegunn/go-shellwords"
)

var ErrInvalidHeaderArg = errors.New("invalid header argument")

// buildArgs loads header, prompt, keybind and args from Options.
func (m *Menu[T]) buildArgs() error {
	if _, err := m.buildHeader(); err != nil {
		return err
	}
	return m.buildKeybindArgs()
}

// buildKeybindString builds the keybind string for FZF.
func (m *Menu[T]) buildKeybindString() string {
	keys := m.keymaps.List()
	if len(keys) == 0 {
		return ""
	}

	keybinds := make([]string, 0, len(keys))
	for _, k := range keys {
		if k.Action == "" {
			slog.Warn("build keybind ignore", "bind", k.Bind, "action", k.Action)
			continue
		}

		if k.Enabled {
			keybinds = append(keybinds, fmt.Sprintf("%s:%s", k.Bind, k.Action))
		}
	}

	return strings.Join(keybinds, ",")
}

// buildKeybindArgs appends keybinding arguments to the menu.
func (m *Menu[T]) buildKeybindArgs() error {
	keybindStr := m.buildKeybindString()
	if keybindStr == "" {
		return nil
	}

	bindArg := fmt.Sprintf("%s=%q", m.argsBuilder.bind, keybindStr)
	binds, err := shellwords.Parse(bindArg)
	if err != nil {
		return fmt.Errorf("parse keybinds %q: %w", keybindStr, err)
	}

	m.argsBuilder.Add(binds...)

	return nil
}

// buildHeaderStrings returns the formatted header strings from enabled keymaps.
func (m *Menu[T]) buildHeader() (string, error) {
	if m.noHeader {
		return "", nil
	}

	if !m.showKeybindHeader {
		header := strings.Join(m.header, "")
		m.argsBuilder.WithHeader(header)
		return header, nil
	}

	keybinds := m.keymaps.List()
	for _, k := range keybinds {
		if k == nil || !k.Enabled || k.Hidden {
			continue
		}

		if k.Desc == "" {
			k.Desc = "?"
		}

		m.header = append(m.header, m.headerKeymapFmt(k.BindString(), ":", k.Desc))
	}

	if len(m.header) == 0 {
		return "", nil
	}

	if m.headerSeparator == "" {
		m.headerSeparator = " "
	}

	header := strings.Join(m.header, m.headerSeparatorFmt(m.headerSeparator))
	_, err := shellwords.Parse(header)
	if err != nil {
		return "", err
	}

	m.argsBuilder.WithHeader(header)

	return header, nil
}
