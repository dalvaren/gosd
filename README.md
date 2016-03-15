# gosd

A simple service discovery client for small architectures made in GO

## Installation

Just add the library to your project.

```
$ go get github.com/dalvaren/gosd
```

## Usage

Encode
```
$ cat foo.png | gosr
```

Decode

```
$ cat foo.drcs | gosd > foo.png
```

Use as library

```go
img, _, _ := image.Decode(filename)
sixel.NewEncoder(os.Stdout).Encode(img)
```

## Example

## Contribute creating drivers for it!

Right now, I've only implemented a driver for Redis.
Fork it and feel free to develop other drivers (etcd, memcache, custom apis...).
The driver documentation can be found [here](https://jinzhu.github.io/gorm)!

## Author

Daniel Campos
