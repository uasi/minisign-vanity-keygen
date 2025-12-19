# minisign-vanity-keygen

A utility to generate vanity [minisign](https://jedisct1.github.io/minisign/) key pairs where the public key matches a specified pattern.

## Building

```bash
go build
```

## Usage

```bash
./minisign-vanity-keygen [-overwrite] <regexp> [<regexp> ...]
```

The tool generates two files in the current directory:

- `minisign.pub` (public key)
- `minisign.key` (secret key)

**NOTE:** `minisign.key` is saved unencrypted. To protect it with a password, run:

```bash
minisign -C -s minisign.key
```

### Options

- `<regexp>` - One or more [Go regular expressions](https://pkg.go.dev/regexp/syntax) that must all match the public key.
- `-overwrite` - Overwrite existing `./minisign.key` and `./minisign.pub` files.

## Examples

Generate a key starting with "RWRmini" (note that minisign public keys always start with "RWQ", "RWR", "RWS", or "RWT"):

```bash
./minisign-vanity-keygen '^RWRmini'
```

Generate a key containing either "sign" or "key", with alphanumeric characters only:

```bash
./minisign-vanity-keygen 'sign|key' '^[A-Za-z0-9]+$'
```

## Notes

- This tool relies on the [aead.dev/minisign](https://pkg.go.dev/aead.dev/minisign) library for key generation.
- More specific patterns (longer matches) will take exponentially longer to find. Try patterns of a few characters first to get a sense of how long it will take.
- The tool uses all available CPU cores for optimal performance.
- Progress is reported every 5 seconds during the search.
