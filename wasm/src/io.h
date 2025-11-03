#pragma once

#include <cstdio>
#include <iostream>
#include <stdexcept>

#include "marisa/base.h"

namespace marisa::grimoire {

namespace io {

class Mapper {
public:
    Mapper();
    Mapper(const Mapper &) = delete;
    Mapper &operator=(const Mapper &) = delete;

    void open(const char *filename, int flags = 0);
    void open(const void *ptr, std::size_t size);
    void seek(std::size_t size);
    void swap(Mapper &rhs) noexcept;

    template <typename T>
    void map(T *obj) {
        MARISA_THROW_IF(obj == nullptr, std::invalid_argument);
        *obj = *static_cast<const T *>(map_data(sizeof(T)));
    }

    template <typename T>
    void map(const T **objs, std::size_t num_objs) {
        MARISA_THROW_IF((objs == nullptr) && (num_objs != 0), std::invalid_argument);
        MARISA_THROW_IF(num_objs > (SIZE_MAX / sizeof(T)), std::invalid_argument);
        *objs = static_cast<const T *>(map_data(sizeof(T) * num_objs));
    }

private:
    const void *ptr_ = nullptr;
    std::size_t avail_ = 0;

    const void *map_data(std::size_t size);
};

class Reader {
public:
    Reader();
    Reader(const Reader &) = delete;
    Reader &operator=(const Reader &) = delete;

    void open(const char *filename);
    void open(std::FILE *file);
    void open(int fd);
    void open(std::istream &stream);
    void seek(std::size_t size);

    template <typename T>
    void read(T *obj) {
        MARISA_THROW_IF(obj == nullptr, std::invalid_argument);
        read_data(obj, sizeof(T));
    }

    template <typename T>
    void read(T *objs, std::size_t num_objs) {
        MARISA_THROW_IF((objs == nullptr) && (num_objs != 0), std::invalid_argument);
        MARISA_THROW_IF(num_objs > (SIZE_MAX / sizeof(T)), std::invalid_argument);
        read_data(objs, sizeof(T) * num_objs);
    }

private:
    uint32_t handle_;
    void read_data(void *buf, std::size_t size);
};

class Writer {
public:
    Writer();
    Writer(const Writer &) = delete;
    Writer &operator=(const Writer &) = delete;

    void open(const char *filename);
    void open(std::FILE *file);
    void open(int fd);
    void open(std::ostream &stream);
    void seek(std::size_t size);

    template <typename T>
    void write(const T &obj) {
        write_data(&obj, sizeof(T));
    }

    template <typename T>
    void write(const T *objs, std::size_t num_objs) {
        MARISA_THROW_IF((objs == nullptr) && (num_objs != 0), std::invalid_argument);
        MARISA_THROW_IF(num_objs > (SIZE_MAX / sizeof(T)), std::invalid_argument);
        write_data(objs, sizeof(T) * num_objs);
    }

private:
    uint32_t handle_;
    void write_data(const void *data, std::size_t size);
};

}

using io::Mapper;
using io::Reader;
using io::Writer;

}