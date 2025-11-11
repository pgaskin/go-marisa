#include "io.h"

#include <cstdint>

namespace marisa::grimoire::io {

namespace env {
#ifdef __wasm__
__attribute__((__import_module__("marisa"),__import_name__(("read"))))
extern void read(void *buf, size_t n);
void marisa_read([[maybe_unused]] uintptr_t handle, void *buf, size_t n) { read(buf, n); }

__attribute__((__import_module__("marisa"),__import_name__(("write"))))
extern void write(const void *buf, size_t n);
void marisa_write([[maybe_unused]] uintptr_t handle, const void *buf, size_t n) { write(buf, n); }
#else
extern "C" void marisa_read(uintptr_t handle, void *buf, size_t n);
extern "C" void marisa_write(uintptr_t handle, const void *buf, size_t n);
#endif
}

void Reader::seek(size_t size) {
    env::marisa_read(handle_, nullptr, size);
}

void Writer::seek(size_t size) {
    env::marisa_write(handle_, nullptr, size);
}

void Reader::data(void *buf, size_t size) {
    MARISA_THROW_IF(size != 0 && buf == nullptr, std::logic_error);
    env::marisa_read(handle_, buf, size);
}

void Writer::data(const void *data, size_t size) {
    MARISA_THROW_IF(size != 0 && data == nullptr, std::logic_error);
    env::marisa_write(handle_, const_cast<void*>(data), size);
}

}
