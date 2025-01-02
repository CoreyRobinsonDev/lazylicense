# lazylicense
Grab a license from the internet and get back to coding.
<br>
[Usage](#Usage) <span>&nbsp;•&nbsp;</span> [Install](#Install)

## Usage
Run **lazylicense** from the command line at the root of your project to get a list of open source licenses to add.

## Install
Download pre-built binary for your system here [Releases](https://github.com/CoreyRobinsonDev/lazylicense/releases).

### Compiling from Source
- Clone this repository
```bash
git clone https://github.com/CoreyRobinsonDev/lazylicense.git
```
- Create **lazylicense** binary
```bash
cd lazylicense
go build
```
- Move binary to <code>/usr/local/bin</code> to call it from anywhere in the terminal
```bash
sudo mv ./lazylicense /usr/local/bin
```
- Confirm that the program was built successfully
```bash
lazylicense
```

## License
[Apache License 2.0](./LICENSE)
