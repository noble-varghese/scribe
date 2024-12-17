# Scribe

A fast and flexible text expander for macOS. Scribe helps you create and use custom text snippets to speed up your typing and reduce repetitive text entry.

## Features

- Fast text expansion
- Support for both Intel and Apple Silicon Macs
- Simple snippet management
- Custom delay configuration
- Lightweight and efficient

## Installation

### Using Homebrew (Recommended)

```bash
# Add the tap repository
brew tap noble-varghese/scribe

# Install scribe
brew install scribe
```

### Manual Installation

1. Download the latest release from the [releases page](https://github.com/noble-varghese/scribe/releases)
2. Extract the downloaded archive
3. Move the `scribe` binary to your PATH

## Usage

```bash
# Start scribe
scribe

# Check version
scribe --version

# Get help
scribe --help
```

## Requirements

- macOS 10.15 or later
- Intel or Apple Silicon processor

## Building from Source

```bash
# Clone the repository
git clone https://github.com/noble-varghese/scribe.git

# Change to project directory
cd scribe

# Build
go build -o scribe cmd/scribe/main.go
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you encounter any issues or have questions, please file an issue on the [GitHub issues page](https://github.com/noble-varghese/scribe/issues).
