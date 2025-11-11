#ifndef MARISA_GRIMOIRE_IO_H_
#define MARISA_GRIMOIRE_IO_H_

#include <cstddef>
#include <cstdint>
#include <cstdio>
#include <iostream>
#include <stdexcept>
#include <utility>

#include "../lib/marisa.h"

namespace marisa::grimoire {

namespace io {

class Mapper {
public:
    Mapper() = default;
    Mapper(const Mapper &) = delete;
    Mapper &operator=(const Mapper &) = delete;

    void open([[maybe_unused]] const char *filename, [[maybe_unused]] int flags = 0) {
        MARISA_THROW(std::runtime_error, "not supported");
    }

    void open(const void *ptr, size_t size) {
        ptr_ = ptr;
        avail_ = size;
    }
    void seek(size_t size) {
        data(size);
    }
    void swap(Mapper &rhs) noexcept {
        std::swap(ptr_, rhs.ptr_);
        std::swap(avail_, rhs.avail_);
    }

    template <typename T>
    void map(T *obj) {
        MARISA_THROW_IF(obj == nullptr, std::invalid_argument);
        *obj = *static_cast<const T *>(data(sizeof(T)));
    }

    template <typename T>
    void map(const T **objs, size_t num_objs) {
        MARISA_THROW_IF((objs == nullptr) && (num_objs != 0), std::invalid_argument);
        MARISA_THROW_IF(num_objs > (SIZE_MAX / sizeof(T)), std::invalid_argument);
        *objs = static_cast<const T *>(data(sizeof(T) * num_objs));
    }

private:
    const void *ptr_ = nullptr;
    size_t avail_ = 0;

    const void *data(size_t size) {
        MARISA_THROW_IF(size > avail_, std::runtime_error);
        const char *const data = static_cast<const char*>(ptr_);
        ptr_ = data + size;
        avail_ -= size;
        return data;
    }
};

class Reader {
public:
    Reader() = default;
    Reader(const Reader &) = delete;
    Reader &operator=(const Reader &) = delete;

    void open([[maybe_unused]] std::FILE *file) {
        MARISA_THROW(std::runtime_error, "not supported");
    }
    void open([[maybe_unused]] std::istream &stream) {
        MARISA_THROW(std::runtime_error, "not supported");
    }

    void open(int fd) {
        handle_ = static_cast<uintptr_t>(fd);
    }
    void open(const char *filename) {
        handle_ = reinterpret_cast<uintptr_t>(filename);
    };
    void seek(size_t size);

    template <typename T>
    void read(T *obj) {
        MARISA_THROW_IF(obj == nullptr, std::invalid_argument);
        data(obj, sizeof(T));
    }

    template <typename T>
    void read(T *objs, size_t num_objs) {
        MARISA_THROW_IF((objs == nullptr) && (num_objs != 0), std::invalid_argument);
        MARISA_THROW_IF(num_objs > (SIZE_MAX / sizeof(T)), std::invalid_argument);
        data(objs, sizeof(T) * num_objs);
    }

private:
    void data(void *buf, size_t size);
    uintptr_t handle_;
};

class Writer {
public:
    Writer() = default;
    Writer(const Writer &) = delete;
    Writer &operator=(const Writer &) = delete;

    void open([[maybe_unused]] std::FILE *file) {
        MARISA_THROW(std::runtime_error, "not supported");
    }
    void open([[maybe_unused]] std::ostream &stream) {
        MARISA_THROW(std::runtime_error, "not supported");
    }

    void open(int fd) {
        handle_ = static_cast<uintptr_t>(fd);
    }
    void open(const char *filename) {
        handle_ = reinterpret_cast<uintptr_t>(filename);
    };
    void seek(size_t size);

    template <typename T>
    void write(const T &obj) {
        data(&obj, sizeof(T));
    }

    template <typename T>
    void write(const T *objs, size_t num_objs) {
        MARISA_THROW_IF((objs == nullptr) && (num_objs != 0), std::invalid_argument);
        MARISA_THROW_IF(num_objs > (SIZE_MAX / sizeof(T)), std::invalid_argument);
        data(objs, sizeof(T) * num_objs);
    }

private:
    void data(const void *data, size_t size);
    uintptr_t handle_;
};

}

using io::Mapper;
using io::Reader;
using io::Writer;

}

#endif
