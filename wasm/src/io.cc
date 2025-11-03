#include <utility>

#include "io.h"
#include "wautil.h"

namespace marisa::grimoire::io {

static_assert(sizeof(uint32_t) == sizeof(int));

Mapper::Mapper() = default;
Reader::Reader() = default;
Writer::Writer() = default;

void Mapper::open(const void *ptr, std::size_t size) {
    ptr_ = ptr;
    avail_ = size;
}

void Mapper::seek(std::size_t size) {
    map_data(size);
}

void Mapper::swap(Mapper &rhs) noexcept {
    std::swap(ptr_, rhs.ptr_);
    std::swap(avail_, rhs.avail_);
}

const void *Mapper::map_data(std::size_t size) {
    MARISA_THROW_IF(size > avail_, std::runtime_error);

    const char *const data = static_cast<const char*>(ptr_);
    ptr_ = data + size;
    avail_ -= size;
    return data;
}

void Reader::open(int fd) {
    handle_ = static_cast<uint32_t>(fd);
}

void Reader::seek(std::size_t size) {
    MARISA_THROW_IF(!handle_, std::logic_error);
    wautil::read(handle_, size);
}

void Reader::read_data(void *buf, std::size_t size) {
    MARISA_THROW_IF(!handle_, std::logic_error);
    wautil::read(handle_, static_cast<char*>(buf), size);
}

void Writer::open(int fd) {
    handle_ = static_cast<uint32_t>(fd);
}

void Writer::seek(std::size_t size) {
    MARISA_THROW_IF(!handle_, std::logic_error);
    wautil::write(handle_, size);
}

void Writer::write_data(const void *data, std::size_t size) {
    MARISA_THROW_IF(!handle_, std::logic_error);
    wautil::write(handle_, static_cast<const char*>(data), size);
}

}
