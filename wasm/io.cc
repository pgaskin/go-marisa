#include <utility>

#include "io.h"

namespace marisa::grimoire::io {

namespace cb {
#ifdef __wasm__
__attribute__((__import_module__("marisa"),__import_name__(("read"))))
#endif
extern void read(void *buf, size_t n);

#ifdef __wasm__
__attribute__((__import_module__("marisa"),__import_name__(("write"))))
#endif
extern void write(const void *buf, size_t n);
}

Mapper::Mapper() = default;
Reader::Reader() = default;
Writer::Writer() = default;

void Mapper::open(const void *ptr, size_t size) {
    ptr_ = ptr;
    avail_ = size;
}

void Mapper::seek(size_t size) {
    data(size);
}

void Mapper::swap(Mapper &rhs) noexcept {
    std::swap(ptr_, rhs.ptr_);
    std::swap(avail_, rhs.avail_);
}

const void *Mapper::data(size_t size) {
    MARISA_THROW_IF(size > avail_, std::runtime_error);
    const char *const data = static_cast<const char*>(ptr_);
    ptr_ = data + size;
    avail_ -= size;
    return data;
}

void Reader::open(int) {}
void Writer::open(int) {}

void Reader::seek(size_t size) {
    cb::read(nullptr, size);
}

void Writer::seek(size_t size) {
    cb::write(nullptr, size);
}

void Reader::data(void *buf, size_t size) {
    MARISA_THROW_IF(size && !buf, std::logic_error);
    cb::read(buf, size);
}

void Writer::data(const void *data, size_t size) {
    MARISA_THROW_IF(size && !data, std::logic_error);
    cb::write(const_cast<void*>(data), size);
}

}
