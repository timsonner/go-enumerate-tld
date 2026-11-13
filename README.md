# go-enumerate-tld

This program takes one name and checks it against every suffix in `tld-list.txt`. For each line it sends `GET http://<name>.<suffix>` and prints the URL when the response status is 200.

That is a live HTTP check, not a DNS or registry lookup. A printed line means some server answered 200 for that URL. Connection errors and every other status are ignored.

## Prerequisites

- Go. The program imports only the standard library: `bufio`, `fmt`, `net/http`, `os`, and `sync`.
- The file `tld-list.txt` in the current working directory. The program opens that relative path and does not search anywhere else.

Commands below were run with `go version go1.22.2 linux/amd64`. There is no `go.mod` and there are no third-party modules.

## Install, build, run, and test

Clone the repository and work from its root so `tld-list.txt` is in the working directory:

```sh
git clone https://github.com/timsonner/go-enumerate-tld.git
cd go-enumerate-tld
```

Print usage. This does not open the list and does not make network requests:

```sh
go run go-enum-tld.go
```

Build a binary named `go-enum-tld` in the current directory:

```sh
go build -o go-enum-tld go-enum-tld.go
```

`go build go-enum-tld.go` writes the same filename.

Run a check. Replace `example` with the label to put in front of each suffix. This reads all 3339 lines and starts one request per line:

```sh
go run go-enum-tld.go example
```

The built binary takes the same argument:

```sh
./go-enum-tld example
```

This repository has no test files, so there is no test command.

## Usage

```text
go run go-enum-tld.go <domain or IP>
```

The single argument is copied into the host as-is. It is not checked as a domain or as an IP address.

| Input | Request |
| --- | --- |
| `example` and list line `com` | `http://example.com` |
| `example` and list line `co.ae` | `http://example.co.ae` |
| `192.0.2.1` and list line `com` | `http://192.0.2.1.com` |

`192.0.2.1` is the documentation address from RFC 5737. A value shaped like an address is still just a hostname label. The program does not connect to that address.

Behavior from `go-enum-tld.go`:

- Each list line is one suffix. The list has 3339 lines and no blank lines. Some lines have more than one label, such as `co.ae` and `com.af`.
- Each suffix runs in its own goroutine. Nothing limits how many run at once, and the HTTP client has no timeout.
- The URL is HTTP only, with no path and no query string.
- Go's default client follows redirects. Status 200 is the status after those redirects.
- A hit is printed as `Valid tld: ` plus the URL. Anything else produces no output.
- Missing the argument prints `Usage: go run go-enum-tld.go <domain or IP>` and returns. A missing `tld-list.txt` prints `Error opening file:` and the error. Both paths return from `main`, so the process exit status is still 0.
- Run it from the directory that contains `tld-list.txt`.

## Layout

- `go-enum-tld.go`: the program.
- `tld-list.txt`: suffixes, one per line.
- `.gitignore`: Go binary and workspace ignore patterns. The built `go-enum-tld` binary is not listed there.

## License

This repository has no license file.
