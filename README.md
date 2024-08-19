# Go-Git-Finder

Go Git Finder is a tool designed to scan exposed **.git** URLs.

## Preview
![Image](https://raw.githubusercontent.com/faizdotid/go-git-finder/main/images.png)

## Usage
> All results save into folder **/results**

### Using source code
```bash
go run main.go -f yourfile.txt -t 10
```

### Compile your own binary file
> Make sure you have installed golang
```bash
go build cmd/main.go
```

### Options
- `-f`: Yourfile include URLs
- `-t`: Threads
