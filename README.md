# tetromnom
Bite-sized Tetromino game

## Developer setup
To install and run from source the following are required:
1. Download, install, and setup Golang https://golang.org/dl/
2. Clone/download this project locally.
3. From the project folder use the following command to download the Go module dependencies of this project:
    * `[path/to/tetromnom]$ go mod download`
4. The Ebiten game library may have [additional dependencies to install](https://ebiten.org/documents/install.html),
   depending on the OS.
5. Now you can use the `go run` command to run `tetromnom.go`:
    * `[path/to/tetromnom]$ go run tetromnom.go`

## Controls
* Move and rotate current piece using Arrow Keys
* Drop current piece straight down using Spacebar
