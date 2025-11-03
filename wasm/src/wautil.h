#pragma once

#include <cstdint>
#include <iostream>

void wautil_throw[[noreturn]](const std::exception& ex);

namespace wautil {

class gostream : public std::basic_streambuf<char> {
    using streamsize = std::streamsize;
    uint32_t handle_;
    char_type rbuf_; // single-byte read buffer (i.e. direct access to the io.Reader)
public:
    explicit gostream(uint32_t);
    gostream (const gostream&) = delete;
    gostream& operator= (const gostream&) = delete;

    int_type underflow() override;
    streamsize xsgetn(char_type* buf, streamsize buf_n) override;

    int_type overflow(int_type c = traits_type::eof()) override;
    streamsize xsputn(const char_type* buf, streamsize buf_n) override;
};

class writer : private gostream, public std::ostream {
public: explicit writer(uint32_t handle);
};

class reader : private gostream, public std::istream {
public: explicit reader(uint32_t handle);
};

}
