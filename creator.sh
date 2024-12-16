#!/bin/bash

# Create main project directory
mkdir -p textexpander

# Create all directories
mkdir -p cmd/textexpander
mkdir -p internal/{config,keyboard,matcher,expander,script}
mkdir -p pkg/utils
mkdir -p configs
mkdir -p scripts

# Create all files
touch cmd/main.go

touch internal/config/config.go
touch internal/config/defaults.go

touch internal/keyboard/listener.go
touch internal/keyboard/simulator.go

touch internal/matcher/trie.go
touch internal/matcher/matcher.go

touch internal/expander/expander.go

touch internal/script/executor.go

touch pkg/utils/clipboard.go
touch pkg/utils/keycode.go

touch configs/expander.yaml

touch go.mod
touch README.md

# Add comments to the main files using echo
echo "# Application entry point" > cmd/main.go
echo "# Configuration handling" > internal/config/config.go
echo "# Default settings" > internal/config/defaults.go
echo "# Keyboard event handling" > internal/keyboard/listener.go
echo "# Keystroke simulation" > internal/keyboard/simulator.go
echo "# Trie implementation for pattern matching" > internal/matcher/trie.go
echo "# Pattern matching logic" > internal/matcher/matcher.go
echo "# Text expansion logic" > internal/expander/expander.go
echo "# Custom script execution" > internal/script/executor.go
echo "# Clipboard operations" > pkg/utils/clipboard.go
echo "# Key code mappings" > pkg/utils/keycode.go
echo "# Expansion rules configuration" > configs/expander.yaml