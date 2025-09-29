# T2Stock

**T2Stock** is a terminal-based real-time stock price tracker built with Go.
The name comes from Time Tracker (2 t's -> [T2](https://github.com/franpfeiffer/t2), other tui I've made) and the Stock part is because
I'm a big fan of stocks.


## Prerequisites

### on macOS/Linux
You'll need to install the package manager [`Homebrew`](https://brew.sh/). That's it.

### on Windows
You'll need to install the package manager [`Scoop`](https://scoop.sh/).
Then, you go to the Environment Variables and under the System Variables,
go to Path and add your Scoop installation directory, usually C:\Users\YourUser\scoop\shims

## Installation

### Install via Homebrew (macOS/Linux)
```bash
brew tap franpfeiffer/t2stock # create the tap
brew install t2stock          # install
```
Or `brew install franpfeiffer/t2stock/t2stock`.

Or, in a [`brew bundle`](https://github.com/Homebrew/homebrew-bundle) `Brewfile`:
```bash
tap "franpfeiffer/t2stock"
brew "t2stock"
```

### Install via Scoop (Windows)
```bash
scoop bucket add t2stock https://github.com/franpfeiffer/scoop-t2stock # create the bucket
scoop install t2stock                                                  # install
```

### Install via Package Managers (Linux)

#### Debian/Ubuntu (.deb)
```bash
wget https://github.com/franpfeiffer/t2stock/releases/latest/download/t2stock_1.0.0_linux_amd64.deb
sudo dpkg -i t2stock_1.0.0_linux_amd64.deb
```

#### RHEL/Fedora (.rpm)
```bash
sudo rpm -i https://github.com/franpfeiffer/t2stock/releases/latest/download/t2stock_1.0.0_linux_amd64.rpm
```

#### Arch Linux (AUR)
```bash
yay -S t2stock-bin
```

### Direct Download
Download the appropriate binary for your platform from the [releases page](https://github.com/franpfeiffer/t2stock/releases).

## Usage

### Start T2Stock
Open your terminal and run:
```bash
t2stock
```

The application will launch with a terminal UI showing real-time stock prices and market data.

### Features
- Real-time stock price tracking
- Terminal-based user interface
- Cross-platform support (macOS, Windows, Linux, FreeBSD)
- Multiple installation methods


### License
MIT
