# go-marisa

Go wrapper for [marisa-trie](https://github.com/s-yata/marisa-trie) using [wazero](https://github.com/wazero/wazero).

This library supports little-endian MARISA dictionaries up to 4 GiB. On 32-bit systems, the size is limited to 2 GiB. Big-endian dictionaries (i.e., ones generated with the native tools on big-endian hosts) are not supported.

On platforms which support wazero's JIT compiler (windows/darwin/linux arm64/amd64), it's about 2-3x slower than the native library. Using the interpreter, it's about 70-150x slower.

Memory-mapped dictionaries are supported on unix-like platforms.

The wasm blob is fully reproducible and [verified](https://github.com/pgaskin/go-marisa/attestations).

The API is stable, type-safe, idiomatic, and does not leak implementation details of the marisa-trie library. All errors are handled appropriately and returned. Concurrent usage is not currently supported.

This module also includes drop-in replacements for the native command-line tools. The have compatible input/output and exit codes, but the error messages may differ.
