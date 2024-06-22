# Version Library

A Go library for parsing, comparing, and managing software versions. This library supports semantic versioning and allows for detailed version part comparisons.

## Features

- Parse versions from strings.
- Compare versions (greater than, less than, equal to).
- Handle semantic versioning (major, minor, patch).
- Support for pre-release versions (alpha, beta, RC).
- Find the latest version from a list.

## Installation

To install the library, use `go get`:

```sh
go get github.com/GiGurra/version
```

## Usage

### Parsing Versions

```go
import "github.com/GiGurra/version"

v := version.ParseVersion("1.2.3")
fmt.Println(v) // Output: 1.2.3
```

### Creating Versions

```go
v := version.NewVersion(1, 2, 3)
fmt.Println(v) // Output: 1.2.3
```

### Comparing Versions

```go
v1 := version.ParseVersion("1.2.3")
v2 := version.ParseVersion("1.2.4")

if v1.IsLessThan(v2) {
    fmt.Println("v1 is less than v2")
}
```

### Finding the Latest Version

```go
versions := []string{"1.2.3", "1.2.4", "1.3.0"}
parsedVersions := lo.Map(versions, version.ParseVersion)

latest := version.FindLatestVersionBy(parsedVersions, func(v version.Version) version.Version {
    return v
})

fmt.Println(latest) // Output: 1.3.0
```

## Testing

Run the tests using:

```sh
go test ./...
```

## Dependencies

- [samber/lo](https://github.com/samber/lo) v1.39.0
- [golang.org/x/exp](https://pkg.go.dev/golang.org/x/exp) (indirect)

## License

This project is licensed under the MIT License.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## Contact

For any questions or issues, please open an issue on GitHub.

---

Happy versioning! 🚀
