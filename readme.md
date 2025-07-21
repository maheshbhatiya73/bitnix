# Bitnix

**Bitnix** is a fully terminal-based, real-time network speed monitor built in Go using a modern TUI (Text User Interface).  
It displays upload and download speeds in Mbps with a dynamic fullscreen UI.

## Features

- **Real-time Monitoring**: Displays live upload and download network speeds.
- **TUI-based**: A clean and modern text-based user interface that runs in your terminal.
- **Cross-platform**: Built with Go, it should run on any platform supported by the dependencies.
- **Lightweight**: Minimal resource usage.

## How it Works

Bitnix uses the following components:

- **`gopsutil`**: To fetch network I/O statistics from the underlying operating system.
- **`Bubble Tea`**: A Go framework for building terminal-based user interfaces.
- **`lipgloss`**: For styling the TUI with colors and layouts.

The application polls for network data every second and updates the UI with the calculated upload and download speeds in Mbps.

## Dependencies

- **[github.com/charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea)**: For the terminal user interface.
- **[github.com/shirou/gopsutil/v3](https://github.com/shirou/gopsutil/v3)**: For fetching network statistics.
- **[github.com/charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss)**: For styling the UI.

## Installation

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-username/bitnix.git
    cd bitnix
    ```

2.  **Build the application:**
    ```bash
    go build
    ```

## Usage

Run the compiled binary from your terminal:

```bash
./bitnix
```

Press `q` or `ctrl+c` to quit the application.