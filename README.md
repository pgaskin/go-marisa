# go-marisa

Go bindings for [marisa-trie](https://github.com/s-yata/marisa-trie) using [wazero](https://github.com/wazero/wazero).

This library supports little-endian MARISA dictionaries up to 4 GiB. On 32-bit systems, the size is limited to 2 GiB. Big-endian dictionaries (i.e., ones generated with the native tools on big-endian hosts) are not supported. Concurrent usage is not currently supported. All other features are exposed with a high-level Go interface. All errors are properly handled.

On platforms which support wazero's JIT compiler (arm64, amd64), it's about 2-3x slower than the native library. Using the interpreter, it's about 100-200x slower.

This module also includes drop-in replacements for the native command-line tools. The have compatible input/output and exit codes, but the error messages may differ.
