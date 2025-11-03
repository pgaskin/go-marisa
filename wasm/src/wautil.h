#pragma once

#include <cstdint>
#include <exception>

namespace wautil {

void write(uintptr_t handle, const char *buf, size_t n);
void write(uintptr_t handle, size_t n);
void read(uintptr_t handle, char *buf, size_t n);
void read(uintptr_t handle, size_t n);

void wthrow[[noreturn]](const std::exception& ex);
void wthrow[[noreturn]](const std::type_info &tinfo, const char *what);
void wthrow[[noreturn]](const char *typ, const char *what);
void wthrow[[noreturn]](const char *what);

}
