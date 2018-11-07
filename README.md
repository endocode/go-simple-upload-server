# go-simple-upload-server
Simple HTTP server to save artifacts. This is a fork from Mei Akizuru's go-simple-server https://github.com/mayth/go-simple-upload-server
We needed a testing web server and adapted this to our needs. Main feature is a fake path prefix. We have a very deep path structure and
we needed to virtualize this.

- Removed token authentication
- Added Testing
- Refactoring
- Different flag and dockerfile structure

Please see the security section.

# Build

```
make
```

# Test

```
make test
```

# Usage

## Start Server

```
$ mkdir $HOME/tmp
$ ./go-simple-upload-server -serverRoot $HOME/tmp
```

## Options

```
$ ./go_simple_upload_server --help
```

# Security

No security! Do not use this in production - you have been warned.

# Docker

```
$ make serve
```
