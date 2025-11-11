// AUTOMATICALLY GENERATED, DO NOT EDIT!
// marisa-trie v0.3.1 (3e87d53b78e15f2f43783d5e376561a8c9722051) (patched)
//
// ### COPYING
// 
// libmarisa and its command line tools are licensed under BSD-2-Clause OR LGPL-2.1-or-later.
// 
// #### The BSD 2-clause license
// 
// Copyright (c) 2010-2025, Susumu Yata
// All rights reserved.
// 
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
// 
// - Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
// - Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
// 
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
// 
// #### The LGPL 2.1 or any later version
// 
// marisa-trie - A static and space-efficient trie data structure.
// Copyright (C) 2010-2025  Susumu Yata
// 
// This library is free software; you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 2.1 of the License, or (at your option) any later version.
// 
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Lesser General Public License for more details.
// 
// You should have received a copy of the GNU Lesser General Public
// License along with this library; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301  USA
//
// clang-format off

#line 1 "include/marisa/base.h"
#ifndef MARISA_BASE_H_
#define MARISA_BASE_H_

#include <cstddef>
#include <cstdint>
#include <exception>
#include <system_error>
#include <utility>

// These aliases are left for backward compatibility.
using marisa_uint8 [[deprecated]] = std::uint8_t;
using marisa_uint16 [[deprecated]] = std::uint16_t;
using marisa_uint32 [[deprecated]] = std::uint32_t;
using marisa_uint64 [[deprecated]] = std::uint64_t;

#if UINTPTR_MAX == UINT64_MAX
 #define MARISA_WORD_SIZE 64
#elif UINTPTR_MAX == UINT32_MAX
 #define MARISA_WORD_SIZE 32
#else
 #error Failed to detect MARISA_WORD_SIZE
#endif

// These constant variables are left for backward compatibility.
[[deprecated]] constexpr auto MARISA_UINT8_MAX = UINT8_MAX;
[[deprecated]] constexpr auto MARISA_UINT16_MAX = UINT16_MAX;
[[deprecated]] constexpr auto MARISA_UINT32_MAX = UINT32_MAX;
[[deprecated]] constexpr auto MARISA_UINT64_MAX = UINT64_MAX;
[[deprecated]] constexpr auto MARISA_SIZE_MAX = SIZE_MAX;

#define MARISA_INVALID_LINK_ID UINT32_MAX
#define MARISA_INVALID_KEY_ID  UINT32_MAX
#define MARISA_INVALID_EXTRA   (UINT32_MAX >> 8)

// Error codes are defined as members of marisa_error_code. This library throws
// an exception with one of the error codes when an error occurs.
enum marisa_error_code {
  // MARISA_OK means that a requested operation has succeeded. In practice, an
  // exception never has MARISA_OK because it is not an error.
  MARISA_OK = 0,

  // MARISA_STATE_ERROR means that an object was not ready for a requested
  // operation. For example, an operation to modify a fixed vector throws an
  // exception with MARISA_STATE_ERROR.
  MARISA_STATE_ERROR = 1,

  // MARISA_NULL_ERROR means that an invalid nullptr has been given.
  MARISA_NULL_ERROR = 2,

  // MARISA_BOUND_ERROR means that an operation has tried to access an out of
  // range address.
  MARISA_BOUND_ERROR = 3,

  // MARISA_RANGE_ERROR means that an out of range value has appeared in
  // operation.
  MARISA_RANGE_ERROR = 4,

  // MARISA_CODE_ERROR means that an undefined code has appeared in operation.
  MARISA_CODE_ERROR = 5,

  // MARISA_RESET_ERROR means that a smart pointer has tried to reset itself.
  MARISA_RESET_ERROR = 6,

  // MARISA_SIZE_ERROR means that a size has exceeded a library limitation.
  MARISA_SIZE_ERROR = 7,

  // MARISA_MEMORY_ERROR means that a memory allocation has failed.
  MARISA_MEMORY_ERROR = 8,

  // MARISA_IO_ERROR means that an I/O operation has failed.
  MARISA_IO_ERROR = 9,

  // MARISA_FORMAT_ERROR means that input was in invalid format.
  MARISA_FORMAT_ERROR = 10,
};

// Flags for memory mapping are defined as members of marisa_map_flags.
// Trie::open() accepts a combination of these flags.
enum marisa_map_flags {
  // MARISA_MAP_POPULATE specifies MAP_POPULATE.
  MARISA_MAP_POPULATE = 1 << 0,
};

// Min/max values, flags and masks for dictionary settings are defined below.
// Please note that unspecified settings will be replaced with the default
// settings. For example, 0 is equivalent to (MARISA_DEFAULT_NUM_TRIES |
// MARISA_DEFAULT_TRIE | MARISA_DEFAULT_TAIL | MARISA_DEFAULT_ORDER).

// A dictionary consists of 3 tries in default. Usually more tries make a
// dictionary space-efficient but time-inefficient.
enum marisa_num_tries {
  MARISA_MIN_NUM_TRIES = 0x00001,
  MARISA_MAX_NUM_TRIES = 0x0007F,
  MARISA_DEFAULT_NUM_TRIES = 0x00003,
};

// This library uses a cache technique to accelerate search functions. The
// following enumerated type marisa_cache_level gives a list of available cache
// size options. A larger cache enables faster search but takes a more space.
enum marisa_cache_level {
  MARISA_HUGE_CACHE = 0x00080,
  MARISA_LARGE_CACHE = 0x00100,
  MARISA_NORMAL_CACHE = 0x00200,
  MARISA_SMALL_CACHE = 0x00400,
  MARISA_TINY_CACHE = 0x00800,
  MARISA_DEFAULT_CACHE = MARISA_NORMAL_CACHE
};

// This library provides 2 kinds of TAIL implementations.
enum marisa_tail_mode {
  // MARISA_TEXT_TAIL merges last labels as zero-terminated strings. So, it is
  // available if and only if the last labels do not contain a NULL character.
  // If MARISA_TEXT_TAIL is specified and a NULL character exists in the last
  // labels, the setting is automatically switched to MARISA_BINARY_TAIL.
  MARISA_TEXT_TAIL = 0x01000,

  // MARISA_BINARY_TAIL also merges last labels but as byte sequences. It uses
  // a bit vector to detect the end of a sequence, instead of NULL characters.
  // So, MARISA_BINARY_TAIL requires a larger space if the average length of
  // labels is greater than 8.
  MARISA_BINARY_TAIL = 0x02000,

  MARISA_DEFAULT_TAIL = MARISA_TEXT_TAIL,
};

// The arrangement of nodes affects the time cost of matching and the order of
// predictive search.
enum marisa_node_order {
  // MARISA_LABEL_ORDER arranges nodes in ascending label order.
  // MARISA_LABEL_ORDER is useful if an application needs to predict keys in
  // label order.
  MARISA_LABEL_ORDER = 0x10000,

  // MARISA_WEIGHT_ORDER arranges nodes in descending weight order.
  // MARISA_WEIGHT_ORDER is generally a better choice because it enables faster
  // matching.
  MARISA_WEIGHT_ORDER = 0x20000,

  MARISA_DEFAULT_ORDER = MARISA_WEIGHT_ORDER,
};

enum marisa_config_mask {
  MARISA_NUM_TRIES_MASK = 0x0007F,
  MARISA_CACHE_LEVEL_MASK = 0x00F80,
  MARISA_TAIL_MODE_MASK = 0x0F000,
  MARISA_NODE_ORDER_MASK = 0xF0000,
  MARISA_CONFIG_MASK = 0xFFFFF
};

namespace marisa {

// These aliases are left for backward compatibility.
using UInt8 [[deprecated]] = std::uint8_t;
using UInt16 [[deprecated]] = std::uint16_t;
using UInt32 [[deprecated]] = std::uint32_t;
using UInt64 [[deprecated]] = std::uint64_t;

using std::uint16_t;
using std::uint32_t;
using std::uint64_t;
using std::uint8_t;

using ErrorCode = marisa_error_code;
using CacheLevel = marisa_cache_level;
using TailMode = marisa_tail_mode;
using NodeOrder = marisa_node_order;

// This is left for backward compatibility.
using std::swap;

// This is left for backward compatibility.
using Exception = std::exception;

}  // namespace marisa

// These macros are used to convert a line number to a string constant.
#define MARISA_INT_TO_STR(value) #value
#define MARISA_LINE_TO_STR(line) MARISA_INT_TO_STR(line)
#define MARISA_LINE_STR          MARISA_LINE_TO_STR(__LINE__)

// MARISA_THROW throws an exception with a filename, a line number, an error
// code and an error message. The message format is as follows:
//  "__FILE__:__LINE__: error_code: error_message"
#define MARISA_THROW(error_type, error_message)                     \
  (throw (error_type)(__FILE__ ":" MARISA_LINE_STR ": " #error_type \
                               ": " error_message))

// MARISA_THROW_IF throws an exception if `condition' is true.
#define MARISA_THROW_IF(condition, error_type) \
  (void)((!(condition)) || (MARISA_THROW(error_type, #condition), 0))

// MARISA_THROW_SYSTEM_ERROR_IF throws an exception if `condition` is true.
// ::GetLastError() or errno should be passed as `error_value`.
#define MARISA_THROW_SYSTEM_ERROR_IF(condition, error_value, error_category,   \
                                     function_name)                            \
  (void)((!(condition)) ||                                                     \
         (throw std::system_error(                                             \
              std::error_code(error_value, error_category),                    \
              __FILE__ ":" MARISA_LINE_STR                                     \
                       ": std::system_error: " function_name ": " #condition), \
          false))

// #ifndef MARISA_USE_EXCEPTIONS
//  #if defined(__GNUC__) && !defined(__EXCEPTIONS)
//   #define MARISA_USE_EXCEPTIONS 0
//  #elif defined(__clang__) && !defined(__cpp_exceptions)
//   #define MARISA_USE_EXCEPTIONS 0
//  #elif defined(_MSC_VER) && !_HAS_EXCEPTIONS
//   #define MARISA_USE_EXCEPTIONS 0
//  #else
//   #define MARISA_USE_EXCEPTIONS 1
//  #endif
// #endif

// #if MARISA_USE_EXCEPTIONS
//  #define MARISA_TRY      try
//  #define MARISA_CATCH(x) catch (x)
// #else
//  #define MARISA_TRY      if (true)
//  #define MARISA_CATCH(x) if (false)
// #endif

#endif  // MARISA_BASE_H_
#line 1 "include/marisa/key.h"
#ifndef MARISA_KEY_H_
#define MARISA_KEY_H_

#include <cassert>
#include <string_view>

//amalgamate marisa/base.h

namespace marisa {

class Key {
 public:
  Key() = default;
  Key(const Key &key) = default;
  Key &operator=(const Key &key) = default;

  char operator[](std::size_t i) const {
    assert(i < length_);
    return ptr_[i];
  }

  void set_str(std::string_view str) {
    set_str(str.data(), str.length());
  }
  void set_str(const char *str) {
    assert(str != nullptr);
    std::size_t length = 0;
    while (str[length] != '\0') {
      ++length;
    }
    assert(length <= UINT32_MAX);
    ptr_ = str;
    length_ = static_cast<uint32_t>(length);
  }
  void set_str(const char *ptr, std::size_t length) {
    assert((ptr != nullptr) || (length == 0));
    assert(length <= UINT32_MAX);
    ptr_ = ptr;
    length_ = static_cast<uint32_t>(length);
  }
  void set_id(std::size_t id) {
    assert(id <= UINT32_MAX);
    union_.id = static_cast<uint32_t>(id);
  }
  void set_weight(float weight) {
    union_.weight = weight;
  }

  std::string_view str() const {
    return std::string_view(ptr_, length_);
  }
  const char *ptr() const {
    return ptr_;
  }
  std::size_t length() const {
    return length_;
  }
  std::size_t id() const {
    return union_.id;
  }
  float weight() const {
    return union_.weight;
  }

  void clear() noexcept {
    Key().swap(*this);
  }
  void swap(Key &rhs) noexcept {
    std::swap(ptr_, rhs.ptr_);
    std::swap(length_, rhs.length_);
    std::swap(union_.id, rhs.union_.id);
  }

 private:
  const char *ptr_ = nullptr;
  uint32_t length_ = 0;
  union Union {
    uint32_t id = 0;
    float weight;
  } union_;
};

}  // namespace marisa

#endif  // MARISA_KEY_H_
#line 1 "include/marisa/query.h"
#ifndef MARISA_QUERY_H_
#define MARISA_QUERY_H_

#include <cassert>
#include <string_view>

//amalgamate marisa/base.h

namespace marisa {

class Query {
 public:
  Query() = default;
  Query(const Query &query) = default;

  Query &operator=(const Query &query) = default;

  char operator[](std::size_t i) const {
    assert(i < length_);
    return ptr_[i];
  }

  void set_str(std::string_view str) {
    set_str(str.data(), str.length());
  }
  void set_str(const char *str) {
    assert(str != nullptr);
    std::size_t length = 0;
    while (str[length] != '\0') {
      ++length;
    }
    ptr_ = str;
    length_ = length;
  }
  void set_str(const char *ptr, std::size_t length) {
    assert((ptr != nullptr) || (length == 0));
    ptr_ = ptr;
    length_ = length;
  }
  void set_id(std::size_t id) {
    id_ = id;
  }

  std::string_view str() const {
    return std::string_view(ptr_, length_);
  }
  const char *ptr() const {
    return ptr_;
  }
  std::size_t length() const {
    return length_;
  }
  std::size_t id() const {
    return id_;
  }

  void clear() noexcept {
    Query().swap(*this);
  }
  void swap(Query &rhs) noexcept {
    std::swap(ptr_, rhs.ptr_);
    std::swap(length_, rhs.length_);
    std::swap(id_, rhs.id_);
  }

 private:
  const char *ptr_ = nullptr;
  std::size_t length_ = 0;
  std::size_t id_ = 0;
};

}  // namespace marisa

#endif  // MARISA_QUERY_H_
#line 1 "include/marisa/agent.h"
#ifndef MARISA_AGENT_H_
#define MARISA_AGENT_H_

#include <cassert>
#include <memory>
#include <string_view>

//amalgamate marisa/key.h
//amalgamate marisa/query.h

namespace marisa {
namespace grimoire::trie {

class State;

}  // namespace grimoire::trie

class Agent {
 public:
  Agent();
  ~Agent();

  Agent(const Agent &other);
  Agent &operator=(const Agent &other);
  Agent(Agent &&other) noexcept;
  Agent &operator=(Agent &&other) noexcept;

  const Query &query() const {
    return query_;
  }
  const Key &key() const {
    return key_;
  }

  void set_query(std::string_view str) {
    set_query(str.data(), str.length());
  }
  void set_query(const char *str);
  void set_query(const char *ptr, std::size_t length);
  void set_query(std::size_t key_id);

  const grimoire::trie::State &state() const {
    return *state_;
  }
  grimoire::trie::State &state() {
    return *state_;
  }

  void set_key(std::string_view str) {
    set_key(str.data(), str.length());
  }
  void set_key(const char *str) {
    assert(str != nullptr);
    key_.set_str(str);
  }
  void set_key(const char *ptr, std::size_t length) {
    assert((ptr != nullptr) || (length == 0));
    assert(length <= UINT32_MAX);
    key_.set_str(ptr, length);
  }
  void set_key(std::size_t id) {
    assert(id <= UINT32_MAX);
    key_.set_id(id);
  }

  bool has_state() const {
    return state_ != nullptr;
  }
  void init_state();

  void clear() noexcept;
  void swap(Agent &rhs) noexcept;

 private:
  Query query_;
  Key key_;
  std::unique_ptr<grimoire::trie::State> state_;
};

}  // namespace marisa

#endif  // MARISA_AGENT_H_
#line 1 "include/marisa/keyset.h"
#ifndef MARISA_KEYSET_H_
#define MARISA_KEYSET_H_

#include <cassert>
#include <memory>
#include <string_view>

//amalgamate marisa/key.h

namespace marisa {

class Keyset {
 public:
  enum {
    BASE_BLOCK_SIZE = 4096,
    EXTRA_BLOCK_SIZE = 1024,
    KEY_BLOCK_SIZE = 256
  };

  Keyset();

  Keyset(const Keyset &) = delete;
  Keyset &operator=(const Keyset &) = delete;

  void push_back(const Key &key);
  void push_back(const Key &key, char end_marker);

  void push_back(std::string_view str, float weight = 1.0) {
    push_back(str.data(), str.length(), weight);
  }
  void push_back(const char *str);
  void push_back(const char *ptr, std::size_t length, float weight = 1.0);

  const Key &operator[](std::size_t i) const {
    assert(i < size_);
    return key_blocks_[i / KEY_BLOCK_SIZE][i % KEY_BLOCK_SIZE];
  }
  Key &operator[](std::size_t i) {
    assert(i < size_);
    return key_blocks_[i / KEY_BLOCK_SIZE][i % KEY_BLOCK_SIZE];
  }

  std::size_t num_keys() const {
    return size_;
  }

  bool empty() const {
    return size_ == 0;
  }
  std::size_t size() const {
    return size_;
  }
  std::size_t total_length() const {
    return total_length_;
  }

  void reset();

  void clear() noexcept;
  void swap(Keyset &rhs) noexcept;

 private:
  std::unique_ptr<std::unique_ptr<char[]>[]> base_blocks_;
  std::size_t base_blocks_size_ = 0;
  std::size_t base_blocks_capacity_ = 0;
  std::unique_ptr<std::unique_ptr<char[]>[]> extra_blocks_;
  std::size_t extra_blocks_size_ = 0;
  std::size_t extra_blocks_capacity_ = 0;
  std::unique_ptr<std::unique_ptr<Key[]>[]> key_blocks_;
  std::size_t key_blocks_size_ = 0;
  std::size_t key_blocks_capacity_ = 0;
  char *ptr_ = nullptr;
  std::size_t avail_ = 0;
  std::size_t size_ = 0;
  std::size_t total_length_ = 0;

  char *reserve(std::size_t size);

  void append_base_block();
  void append_extra_block(std::size_t size);
  void append_key_block();
};

}  // namespace marisa

#endif  // MARISA_KEYSET_H_
#line 1 "include/marisa/trie.h"
#ifndef MARISA_TRIE_H_
#define MARISA_TRIE_H_

#include <memory>

//amalgamate marisa/agent.h
//amalgamate marisa/keyset.h

namespace marisa {
namespace grimoire::trie {

class LoudsTrie;

}  // namespace grimoire::trie

class Trie {
  friend class TrieIO;

 public:
  Trie();
  ~Trie();

  Trie(const Trie &) = delete;
  Trie &operator=(const Trie &) = delete;

  Trie(Trie &&) noexcept;
  Trie &operator=(Trie &&) noexcept;

  void build(Keyset &keyset, int config_flags = 0);

  void mmap(const char *filename, int flags = 0);
  void map(const void *ptr, std::size_t size);

  void load(const char *filename);
  void read(int fd);

  void save(const char *filename) const;
  void write(int fd) const;

  bool lookup(Agent &agent) const;
  void reverse_lookup(Agent &agent) const;
  bool common_prefix_search(Agent &agent) const;
  bool predictive_search(Agent &agent) const;

  std::size_t num_tries() const;
  std::size_t num_keys() const;
  std::size_t num_nodes() const;

  TailMode tail_mode() const;
  NodeOrder node_order() const;

  bool empty() const;
  std::size_t size() const;
  std::size_t total_size() const;
  std::size_t io_size() const;

  void clear() noexcept;
  void swap(Trie &rhs) noexcept;

 private:
  std::unique_ptr<grimoire::trie::LoudsTrie> trie_;
};

}  // namespace marisa

#endif  // MARISA_TRIE_H_
#line 1 "include/marisa.h"
#ifndef MARISA_H_
#define MARISA_H_

// "marisa/stdio.h" includes <cstdio> for I/O using std::FILE.
//amalgamate marisa/stdio.h

// "marisa/iostream.h" includes <iosfwd> for I/O using std::iostream.
//amalgamate marisa/iostream.h

// You can use <marisa/trie.h> instead of <marisa.h> if you don't need the
// above I/O interfaces and don't want to include the above I/O headers.
//amalgamate marisa/trie.h

#endif  // MARISA_H_
