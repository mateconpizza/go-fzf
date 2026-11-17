package menu

import (
	"errors"
	"fmt"
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
	log := m.log().With("fn", "build-keybind-string")

	keys := m.keymaps.List()
	if len(keys) == 0 {
		log.Debug("no keybindings configured")
		return ""
	}

	keybinds := make([]string, 0, len(keys))
	skipped := 0

	for _, k := range keys {
		if k == nil {
			skipped++
			continue
		}

		if k.Action == "" {
			log.Warn("skipping keybinding without action", "bind", k.Bind)
			skipped++
			continue
		}

		if k.Enabled {
			keybinds = append(keybinds, fmt.Sprintf("%s:%s", k.Bind, k.Action))
		}
	}

	result := strings.Join(keybinds, ",")

	log.Debug("keybind string built",
		"count", len(keybinds),
		"skipped", skipped,
	)

	return result
}

// buildKeybindArgs appends keybinding arguments to the menu.
func (m *Menu[T]) buildKeybindArgs() error {
	log := m.log().With("fn", "build-keybind-args")

	keybindStr := m.buildKeybindString()
	if keybindStr == "" {
		log.Debug("no keybind arguments to build")
		return nil
	}

	bindArg := fmt.Sprintf("%s=%q", m.argsBuilder.bind, keybindStr)

	binds, err := shellwords.Parse(bindArg)
	if err != nil {
		log.Error("failed to parse keybind arguments", "error", err)
		return fmt.Errorf("parse keybinds %q: %w", keybindStr, err)
	}

	m.argsBuilder.Add(binds...)

	log.Debug("keybind arguments built", "count", len(binds))

	return nil
}

// buildHeaderStrings returns the formatted header strings from enabled keymaps.
func (m *Menu[T]) buildHeader() (string, error) {
	log := m.log().With("fn", "build-header")

	if m.noHeader {
		log.Debug("header disabled")
		return "", nil
	}

	if !m.showKeybindHeader {
		header := strings.Join(m.header, "")
		m.argsBuilder.WithHeader(header)

		log.Debug("using custom header", "parts", len(m.header))

		return header, nil
	}

	keybinds := m.keymaps.List()
	added := 0

	for _, k := range keybinds {
		if k == nil || !k.Enabled || k.Hidden {
			continue
		}

		if k.Desc == "" {
			k.Desc = "?"
		}

		m.header = append(m.header, m.headerKeymapFmt(k.BindString(), ":", k.Desc))
		added++
	}

	if len(m.header) == 0 {
		log.Debug("no header content")
		return "", nil
	}

	if m.headerSeparator == "" {
		m.headerSeparator = " "
	}

	header := strings.Join(m.header, m.headerSeparatorFmt(m.headerSeparator))
	if _, err := shellwords.Parse(header); err != nil {
		log.Error("failed to parse header", "error", err)
		return "", err
	}

	m.argsBuilder.WithHeader(header)

	log.Debug("header built",
		"parts", len(m.header),
		"keybinds", added,
	)

	return header, nil
}
