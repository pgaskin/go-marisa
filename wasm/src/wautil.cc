#include "wautil.h"

#include <cstdint>
#include <cstddef>
#include <cstring>
#include <cxxabi.h>
#include <functional>
#include <iostream>
#include <memory>
#include <new>
#include <stdexcept>
#include <typeinfo>

#ifdef __wasm__
#define wasm_import(module, name) extern "C" __attribute__((__import_module__(#module),__import_name__((#name))))
#define wasm_export(name) extern "C" __attribute__((export_name(#name)))
#define wasm_constructor __attribute__((constructor))
#else
#define wasm_import(module, name) extern "C"
#define wasm_export(name) extern "C"
#define wasm_constructor __attribute__((constructor))
#endif

namespace wautil {

/** cxx_write calls io.Writer.Write(buf[:n]), throwing on error */
wasm_import(wautil, cxx_write) size_t cxx_write(uintptr_t handle, const char *buf, size_t n);

/** cxx_read calls io.Reader.Read(buf[:n]), throwing on error, and returning 0 on eof */
wasm_import(wautil, cxx_read) size_t cxx_read(uintptr_t handle, const char *buf, size_t n);

/** cxx_throw throws an exception, stopping execution immeidately without running destructors */
wasm_import(wautil, cxx_throw) void cxx_throw[[noreturn]](const char *typ, size_t typlen, const char *std, size_t stdlen, const char *msg, size_t msglen);

/** alloc allocates memory */
wasm_export(wautil_alloc) void *wautil_alloc(size_t n) { return ::operator new(n); }

/** free frees memory */
wasm_export(wautil_free) void wautil_free(void *ptr) { ::operator delete(ptr); }

gostream::gostream(uint32_t handle) : handle_(handle) {
    this->setg(&this->rbuf_, &this->rbuf_ + 1, &this->rbuf_ + 1); // set the buffer, but at the end to force the next read to underflow
}

gostream::int_type gostream::underflow() {
    streamsize n = this->xsgetn(&this->rbuf_, sizeof(char_type));
    if (n == traits_type::eof()) return traits_type::eof();
    return traits_type::to_int_type(this->rbuf_);
}

gostream::streamsize gostream::xsgetn(char_type* buf, streamsize buf_n) {
    streamsize n = 0;
    while (n != buf_n) {
        size_t x = cxx_read(handle_, buf+n, static_cast<size_t>(buf_n-n));
        if (x == 0) break;
        n += x;
    }
    this->rbuf_ = n>0 ? buf[n-1] : 0;
    this->setg(&this->rbuf_, &this->rbuf_, &this->rbuf_ + (n>0 ? 1 : 0));
    if (this->gptr() == this->egptr()) return traits_type::eof();
    return n;
}

gostream::int_type gostream::overflow(int_type c) {
    if (traits_type::eq_int_type(c, traits_type::eof())) return 0;
    this->xsputn(reinterpret_cast<traits_type::char_type*>(&c), 1);
    return c;
}

gostream::streamsize gostream::xsputn(const char_type* buf, streamsize buf_n) {
    return static_cast<streamsize>(cxx_write(handle_, buf, static_cast<size_t>(buf_n)));
}

reader::reader(uint32_t handle) : gostream(handle), std::istream(this) {}

writer::writer(uint32_t handle) : gostream(handle), std::ostream(this) {}

}

void wautil_throw[[noreturn]](const std::exception& ex) {
    // see https://en.cppreference.com/w/cpp/error/exception.html for the hierachy
    // search for "MARISA_THROW" to see what's used
    const char *typ = typeid(ex).name();
    size_t typlen = std::strlen(typ);

    char dtyp[128];
    size_t dtyplen = sizeof(dtyp);
    int dstatus = -1;
    abi::__cxa_demangle(typ, dtyp, &dtyplen, &dstatus);
    if (dstatus == 0) {
        typ = dtyp;
        typlen = dtyplen-1;
    }

    // these must be a subset of the ones defined in internal/wautil/error.go for unwrapping to work correctly
    const char *std = "exception";
    if (false) {}
    #define std_(exception) else if (dynamic_cast<const exception*>(&ex)) std = #exception;
    std_(std::invalid_argument)               //   logic-error
    std_(std::domain_error)                   //   logic-error
    std_(std::length_error)                   //   logic-error
    std_(std::out_of_range)                   //   logic-error
    //std_(std::future_error)                 //   logic-error
    std_(std::logic_error)                    // exception
    std_(std::range_error)                    //   runtime_error
    std_(std::overflow_error)                 //   runtime_error
    std_(std::underflow_error)                //   runtime_error
    //std_(std::regex_error)                  //   runtime_error
    std_(std::ios_base::failure)              //     system_error
    //std_(std::filesystem::filesystem_error  //     system_error
    std_(std::system_error)                   //   runtime_error
    //std_(std::nonexistent_local_time)       //   runtime_error
    //std_(std::format_error)                 //   runtime_error
    std_(std::runtime_error)                  // exception
    std_(std::bad_typeid)                     // exception
    //std_(std::bad_any_cast)                 //   bad_cast
    std_(std::bad_cast)                       // exception
    //std_(std::bad_optional_access)          // exception
    //std_(std::bad_expected_access)          // exception
    std_(std::bad_weak_ptr)                   // exception
    std_(std::bad_function_call)              // exception
    std_(std::bad_array_new_length)           //   bad_alloc
    std_(std::bad_alloc)                      // exception
    std_(std::bad_exception)                  // exception
    //std_(std::bad_variant_access)           // exception
    size_t stdlen = std::strlen(std);

    const char *msg = ex.what();
    size_t msglen = std::strlen(msg);

    wautil::cxx_throw(typ, typlen, std, stdlen, msg, msglen);
}

wasm_constructor static void wautil_new_handler_init() {
    static auto ex = std::bad_alloc(); // preallocate
    std::set_new_handler([]() { wautil_throw(ex); });
}

static uint8_t fallback_exception_buf[4096];

extern "C" void *__cxa_allocate_exception(size_t size) {
    if (size < sizeof(fallback_exception_buf)) {
        return fallback_exception_buf;
    }
    //wautil::cxx_throw(nullptr, 0, nullptr, 0, nullptr, 0);
    const auto mem = ::operator new(size);
    return mem;
}

extern "C" void __cxa_free_exception(void *thrown_exception) {
    if (thrown_exception == fallback_exception_buf) return;
    ::operator delete(thrown_exception);
}

extern "C" void __cxa_throw(void *thrown_exception, [[maybe_unused]] std::type_info *tinfo, [[maybe_unused]] void *(*dest)(void *)) {
    // this is a terrible hack which will work as long as stuff is derived from
    // std::exception* (since the vtable prefix will be the same and
    // dynamic_cast will still be able to figure out the base type since
    // reinterpret_cast doesn't overwrite the vtable), but it will fail badly if
    // anything throws a non-exception (but at least wazero will give us a
    // useful stack trace)
    wautil_throw(*reinterpret_cast<std::exception*>(thrown_exception));
}
