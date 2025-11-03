#include "wautil.h"

#include <cstdint>
#include <cstddef>
#include <cstring>
#include <functional>
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

wasm_import(wautil, cxx_throw) void cxx_throw[[noreturn]](const char *typ, const char *std, const char *msg);
wasm_import(wautil, cxx_write) void cxx_write(uintptr_t handle, const char *buf, size_t n);
wasm_import(wautil, cxx_read_full) void cxx_read_full(uintptr_t handle, char *buf, size_t n);
wasm_import(wautil, cxx_write_zeros) void cxx_write_zeros(uintptr_t handle, size_t n);
wasm_import(wautil, cxx_read_skip) void cxx_read_skip(uintptr_t handle, size_t n);
wasm_export(wautil_alloc) void *wautil_alloc(size_t n) { return ::operator new(n); }
wasm_export(wautil_free) void wautil_free(void *ptr) { ::operator delete(ptr); }

void write(uintptr_t handle, const char *buf, size_t n) { cxx_write(handle, buf, n); }
void write(uintptr_t handle, size_t n) { cxx_write_zeros(handle, n); }
void read(uintptr_t handle, char *buf, size_t n) { cxx_read_full(handle, buf, n); }
void read(uintptr_t handle, size_t n) { cxx_read_skip(handle, n); }

void wthrow[[noreturn]](const std::exception& ex) {
    // see https://en.cppreference.com/w/cpp/error/exception.html for the hierachy
    // search for "MARISA_THROW" to see what's used

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

    wautil::cxx_throw(typeid(ex).name(), std, ex.what());
}

void wthrow[[noreturn]](const char *typ, const char *what) {
    wautil::cxx_throw(typ, nullptr, what);
}

void wthrow[[noreturn]](const std::type_info &tinfo, const char *what) {
    wautil::wthrow(tinfo.name(), what);
}

void wthrow[[noreturn]](const char *what) {
    wautil::wthrow(nullptr, what);
}

static void *wthrow_destruct_obj = nullptr;
static void *(*wthrow_destruct_fn)(void*) = nullptr;

static void wthrow_set_destructor(void *obj, void *(*dest)(void*)) {
    wthrow_destruct_obj = obj;
    wthrow_destruct_fn = dest;
}

wasm_export(wautil_post_throw) void wautil_post_throw() {
    if (wthrow_destruct_fn && wthrow_destruct_obj) {
        wthrow_destruct_fn(wthrow_destruct_obj);
        wthrow_destruct_fn = nullptr;
        wthrow_destruct_obj = nullptr;
    }
}

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

namespace __cxxabiv1 {
class __shim_type_info : public std::type_info {
public:
    virtual ~__shim_type_info();
    virtual void noop1() const;
    virtual void noop2() const;
    virtual bool can_catch(const __shim_type_info *thrown_type, void *&adjustedPtr) const = 0;
};
}

extern "C" void __cxa_throw(void *thrown_exception, std::type_info *tinfo, void *(*dest)(void*)) {
    // this took a while to figure out, but it's really simple in the end
    //
    // see the following links:
    //  - https://github.com/llvm/llvm-project/blob/564c3de67d20d578d05678b49045378fdcf5ccaa/libcxxabi/include/cxxabi.h
    //  - https://github.com/llvm/llvm-project/blob/564c3de67d20d578d05678b49045378fdcf5ccaa/libcxxabi/src/private_typeinfo.cpp
    //  - https://github.com/llvm/llvm-project/blob/564c3de67d20d578d05678b49045378fdcf5ccaa/libcxx/include/typeinfo#L189
    //  - https://github.com/llvm/llvm-project/blob/564c3de67d20d578d05678b49045378fdcf5ccaa/libcxxabi/src/cxa_default_handlers.cpp#L65-L78
    //  - https://github.com/llvm/llvm-project/blob/564c3de67d20d578d05678b49045378fdcf5ccaa/libcxxabi/src/aix_state_tab_eh.inc#L579-L628
    //  - https://github.com/llvm/llvm-project/blob/564c3de67d20d578d05678b49045378fdcf5ccaa/libcxxabi/src/private_typeinfo.h#L17-L55
    //
    // limitations:
    //   - we don't support c++ catch blocks (all exceptions will go straight to go)
    //   - we don't support unwinding the stack, so destructors won't get called

    wautil::wthrow_set_destructor(thrown_exception, dest);
    if (tinfo) {
        auto throw_type = static_cast<const __cxxabiv1::__shim_type_info*>(tinfo);
        auto catch_type = static_cast<const __cxxabiv1::__shim_type_info*>(&typeid(std::exception));
        if (catch_type->can_catch(throw_type, thrown_exception)) {
            wautil::wthrow(*static_cast<std::exception*>(thrown_exception));
        }
        wautil::wthrow(*tinfo, "unknown error type");
    }
    wautil::wthrow("unknown error type");
}

// printf_core (and write/writev/close) is brought in by __abort_message, which
// is only used internally in libcxx in places we won't hit in practice (in
// particular, we override operator new's error handling), so just replace
// __abort_message (and since the messages are mostly static, don't bother with
// the format string)
//  - run `twiggy top wasm/marisa.wasm --retained`
//  - run `twiggy paths wasm/marisa.wasm printf_core`
//  - run `wasm-objdump -x wasm/marisa.wasm`
//  - https://github.com/llvm/llvm-project/blob/f3b407f8f4624461eedfe5a2da540469a0f69dc9/libcxxabi/src/stdlib_new_delete.cpp#L31C13-L43
extern "C" void __abort_message[[noreturn]](const char* fmt, ...) {
    wautil::wthrow("abort", fmt);
}

wasm_constructor static void wautil_new_handler_init() {
    static auto ex = std::bad_alloc(); // preallocate
    std::set_new_handler([]() { wautil::wthrow(ex); });
}
