#pragma once

#include <cstdint>
#include <exception>

void wautil_throw[[noreturn]](const std::exception& ex);

namespace wautil {

void write(uintptr_t handle, const char *buf, size_t n);
void write(uintptr_t handle, size_t n);
void read(uintptr_t handle, char *buf, size_t n);
void read(uintptr_t handle, size_t n);

}
