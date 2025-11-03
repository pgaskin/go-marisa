# go-marisa

Go bindings for [marisa-trie](https://github.com/s-yata/marisa-trie) using [wazero](https://github.com/wazero/wazero).

This library supports little-endian in-memory dictionaries up to 4GB. Memory mapping, big-endian dictionaries, and concurrent usage are not supported. All other features are exposed with a high-level Go interface. All errors are properly handled.
