# minisign-vanity-keygen

A utility to generate vanity [minisign](https://jedisct1.github.io/minisign/) key pairs where the public key matches a specified pattern.

## Building

```bash
go build
```

## Usage

```bash
./minisign-vanity-keygen [-alphanumeric] [-overwrite] <regexp>
```

The tool generates two files in the current directory:

- `minisign.pub` (public key)
- `minisign.key` (secret key)

**NOTE:** `minisign.key` is saved unencrypted. To protect it with a password, run:

```bash
minisign -C -s minisign.key
```

### Options

- `<regexp>` - A [Go regular expression](https://pkg.go.dev/regexp/syntax) to match against the public key.
- `-alphanumeric` - Exclude keys that contain symbols (`+`, `/`, and `=`) when Base64-encoded.
- `-overwrite` - Overwrite existing `./minisign.key` and `./minisign.pub` files.

## Examples

Generate a key starting with "RWRmini" (note that minisign public keys always start with "RWQ", "RWR", "RWS", or "RWT"):

```bash
./minisign-vanity-keygen "^RWRmini"
```

Generate an alphanumeric-only key containing "sign" or "key":

```bash
./minisign-vanity-keygen -alphanumeric "sign|key"
```

## Notes

- This tool relies on the [aead.dev/minisign](https://pkg.go.dev/aead.dev/minisign) library for key generation.
- More specific patterns (longer matches) will take exponentially longer to find. Try patterns of a few characters first to get a sense of how long it will take.
- The tool uses all available CPU cores for optimal performance.
- Progress is reported every 5 seconds during the search.
