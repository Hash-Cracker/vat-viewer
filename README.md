# VAT Viewer

This is a file viewer which is pretty much a fork of [bat](https://github.com/sharkdp/bat) viewer but with only Nord highlighting. There are still some to of the highlighting bugs and issues to be fixed.

## Installation

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/Hash-Cracker/vat-viewer
   ```

2. **Navigate to the Project Directory**:

   ```bash
   cd vat-viewer
   ```

3. **Build the Project**:

   Ensure you have [Go](https://golang.org/dl/) installed. Then, run:

   ```bash
   go mod tidy
   go mod download
   go build main.go
   ```

## Usage

After building the project, execute the following command to run VAT Viewer:

```bash
./vat-viewer /path/to/file
```

Replace `/path/to/file` with the actual path to your VAT data file.


## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
