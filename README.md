# incident CLI 🛡️

A modern and modular CLI tool written in Go for incident management directly from the terminal. `incident` helps track, declare, and isolate the state of technical occurrences quickly and safely.

## 🚀 Features

- **Hybrid Identification**: IDs generated in the `INC-YYYYMMDD-XXXX` format for easy sorting and reference.
- **State Isolation**: Each incident has its own dedicated folder in `.incidents/[ID]`.
- **Git Friendly**: Automatically creates a `.gitignore` file to prevent accidental versioning of local logs and states.
- **Flexible Configuration**: Manage tokens and URLs via `~/.incident.yaml` using Viper.
- **Stateless Architecture**: Designed to be lightweight and scalable.

## 📦 Installation

Ensure you have Go 1.21+ installed.

```powershell
# Clone the repository
git clone https://github.com/ESousa97/goincidentcli.git
cd goincidentcli

# Build the binary
go build -o incident.exe
```

## 🛠️ Usage

### Declare an Incident

```powershell
.\incident.exe declare --title "Authentication API latency spike"
```

**Expected results:**
- Created directory `.incidents/INC-20260405-b7x2/`.
- Configuration loaded from `~/.incident.yaml`.

## ⚙️ Configuration

On the first run, the CLI will automatically create a template in your home directory (`~/.incident.yaml`). Edit this file to fill in your credentials:

```yaml
api_token: "your_token_here"
base_url: "https://api.example.com"
```

## 🏗️ Project Structure

- `/cmd`: Cobra commands (CLI entrypoints).
- `/internal/config`: Configuration loading and typing logic.
- `/internal/incident`: Business domain and incident file manipulation.

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.
