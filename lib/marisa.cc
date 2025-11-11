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

#ifdef __clang__
#pragma clang diagnostic warning "-Wall"
#pragma clang diagnostic warning "-Wextra"
#pragma clang diagnostic warning "-Wconversion"
#pragma clang diagnostic ignored "-Wimplicit-fallthrough"
#pragma clang diagnostic ignored "-Wunused-function"
#pragma clang diagnostic ignored "-Wunused-const-variable"
#elif __GNUC__
#pragma GCC diagnostic warning "-Wall"
#pragma GCC diagnostic warning "-Wextra"
#pragma GCC diagnostic warning "-Wconversion"
#pragma GCC diagnostic ignored "-Wimplicit-fallthrough"
#pragma GCC diagnostic ignored "-Wunused-function"
#pragma GCC diagnostic ignored "-Wunused-const-variable"
#endif

#line 1 "lib/marisa/grimoire/trie/entry.h"
#ifndef MARISA_GRIMOIRE_TRIE_ENTRY_H_
#define MARISA_GRIMOIRE_TRIE_ENTRY_H_

#include <cassert>

#include "./../lib/marisa.h" //amalgamate marisa/base.h

namespace marisa::grimoire::trie {

class Entry {
 public:
  Entry() = default;
  Entry(const Entry &entry) = default;
  Entry &operator=(const Entry &entry) = default;

  char operator[](std::size_t i) const {
    assert(i < length_);
    return *(ptr_ - i);
  }

  void set_str(const char *ptr, std::size_t length) {
    assert((ptr != nullptr) || (length == 0));
    assert(length <= UINT32_MAX);
    ptr_ = ptr + length - 1;
    length_ = static_cast<uint32_t>(length);
  }
  void set_id(std::size_t id) {
    assert(id <= UINT32_MAX);
    id_ = static_cast<uint32_t>(id);
  }

  const char *ptr() const {
    return ptr_ - length_ + 1;
  }
  std::size_t length() const {
    return length_;
  }
  std::size_t id() const {
    return id_;
  }

  class StringComparer {
   public:
    bool operator()(const Entry &lhs, const Entry &rhs) const {
      for (std::size_t i = 0; i < lhs.length(); ++i) {
        if (i == rhs.length()) {
          return true;
        }
        if (lhs[i] != rhs[i]) {
          return static_cast<uint8_t>(lhs[i]) > static_cast<uint8_t>(rhs[i]);
        }
      }
      return lhs.length() > rhs.length();
    }
  };

  class IDComparer {
   public:
    bool operator()(const Entry &lhs, const Entry &rhs) const {
      return lhs.id_ < rhs.id_;
    }
  };

 private:
  const char *ptr_ = nullptr;
  uint32_t length_ = 0;
  uint32_t id_ = 0;
};

}  // namespace marisa::grimoire::trie

#endif  // MARISA_GRIMOIRE_TRIE_ENTRY_H_
#line 1 "lib/marisa/grimoire/vector/rank-index.h"
#ifndef MARISA_GRIMOIRE_VECTOR_RANK_INDEX_H_
#define MARISA_GRIMOIRE_VECTOR_RANK_INDEX_H_

#include <cassert>

#include "./../lib/marisa.h" //amalgamate marisa/base.h

namespace marisa::grimoire::vector {

class RankIndex {
 public:
  RankIndex() = default;

  void set_abs(std::size_t value) {
    assert(value <= UINT32_MAX);
    abs_ = static_cast<uint32_t>(value);
  }
  void set_rel1(std::size_t value) {
    assert(value <= 64);
    rel_lo_ = static_cast<uint32_t>((rel_lo_ & ~0x7FU) | (value & 0x7FU));
  }
  void set_rel2(std::size_t value) {
    assert(value <= 128);
    rel_lo_ = static_cast<uint32_t>((rel_lo_ & ~(0xFFU << 7)) |
                                    ((value & 0xFFU) << 7));
  }
  void set_rel3(std::size_t value) {
    assert(value <= 192);
    rel_lo_ = static_cast<uint32_t>((rel_lo_ & ~(0xFFU << 15)) |
                                    ((value & 0xFFU) << 15));
  }
  void set_rel4(std::size_t value) {
    assert(value <= 256);
    rel_lo_ = static_cast<uint32_t>((rel_lo_ & ~(0x1FFU << 23)) |
                                    ((value & 0x1FFU) << 23));
  }
  void set_rel5(std::size_t value) {
    assert(value <= 320);
    rel_hi_ = static_cast<uint32_t>((rel_hi_ & ~0x1FFU) | (value & 0x1FFU));
  }
  void set_rel6(std::size_t value) {
    assert(value <= 384);
    rel_hi_ = static_cast<uint32_t>((rel_hi_ & ~(0x1FFU << 9)) |
                                    ((value & 0x1FFU) << 9));
  }
  void set_rel7(std::size_t value) {
    assert(value <= 448);
    rel_hi_ = static_cast<uint32_t>((rel_hi_ & ~(0x1FFU << 18)) |
                                    ((value & 0x1FFU) << 18));
  }

  std::size_t abs() const {
    return abs_;
  }
  std::size_t rel1() const {
    return rel_lo_ & 0x7FU;
  }
  std::size_t rel2() const {
    return (rel_lo_ >> 7) & 0xFFU;
  }
  std::size_t rel3() const {
    return (rel_lo_ >> 15) & 0xFFU;
  }
  std::size_t rel4() const {
    return (rel_lo_ >> 23) & 0x1FFU;
  }
  std::size_t rel5() const {
    return rel_hi_ & 0x1FFU;
  }
  std::size_t rel6() const {
    return (rel_hi_ >> 9) & 0x1FFU;
  }
  std::size_t rel7() const {
    return (rel_hi_ >> 18) & 0x1FFU;
  }

 private:
  uint32_t abs_ = 0;
  uint32_t rel_lo_ = 0;
  uint32_t rel_hi_ = 0;
};

}  // namespace marisa::grimoire::vector

#endif  // MARISA_GRIMOIRE_VECTOR_RANK_INDEX_H_
#line 1 "lib/marisa/grimoire/vector/vector.h"
#ifndef MARISA_GRIMOIRE_VECTOR_VECTOR_H_
#define MARISA_GRIMOIRE_VECTOR_VECTOR_H_

#include <cassert>
#include <cstring>
#include <memory>
#include <new>
#include <stdexcept>
#include <type_traits>
#include <utility>

#include "../src/io.h" //amalgamate marisa/grimoire/io.h

namespace marisa::grimoire::vector {

template <typename T>
class Vector {
 public:
  // These assertions are repeated for clarity/robustness where the property
  // is used.
  static_assert(std::is_trivially_copyable_v<T>);
  static_assert(std::is_trivially_destructible_v<T>);

  Vector() = default;
  // `T` is trivially destructible, so default destructor is ok.
  ~Vector() = default;

  Vector(const Vector<T> &other) : fixed_(other.fixed_) {
    if (other.buf_ == nullptr) {
      objs_ = other.objs_;
      const_objs_ = other.const_objs_;
      size_ = other.size_;
      capacity_ = other.capacity_;
    } else {
      copyInit(other.const_objs_, other.size_, other.capacity_);
    }
  }

  Vector &operator=(const Vector<T> &other) {
    clear();
    fixed_ = other.fixed_;
    if (other.buf_ == nullptr) {
      objs_ = other.objs_;
      const_objs_ = other.const_objs_;
      size_ = other.size_;
      capacity_ = other.capacity_;
    } else {
      copyInit(other.const_objs_, other.size_, other.capacity_);
    }
    return *this;
  }

  Vector(Vector &&) noexcept = default;
  Vector &operator=(Vector<T> &&) noexcept = default;

  __attribute__((noinline)) void map(Mapper &mapper) {
    Vector temp;
    temp.map_(mapper);
    swap(temp);
  }

  __attribute__((noinline)) void read(Reader &reader) {
    Vector temp;
    temp.read_(reader);
    swap(temp);
  }

  __attribute__((noinline)) void write(Writer &writer) const {
    write_(writer);
  }

  void push_back(const T &x) {
    assert(!fixed_);
    assert(size_ < max_size());
    reserve(size_ + 1);
    new (&objs_[size_]) T(x);
    ++size_;
  }

  void pop_back() {
    assert(!fixed_);
    assert(size_ != 0);
    --size_;
    static_assert(std::is_trivially_destructible_v<T>);
  }

  // resize() assumes that T's placement new does not throw an exception.
  void resize(std::size_t size) {
    assert(!fixed_);
    reserve(size);
    for (std::size_t i = size_; i < size; ++i) {
      new (&objs_[i]) T;
    }
    static_assert(std::is_trivially_destructible_v<T>);
    size_ = size;
  }

  // resize() assumes that T's placement new does not throw an exception.
  void resize(std::size_t size, const T &x) {
    assert(!fixed_);
    reserve(size);
    if (size > size_) {
      std::fill_n(&objs_[size_], size - size_, x);
    }
    // No need to destroy old elements.
    static_assert(std::is_trivially_destructible_v<T>);
    size_ = size;
  }

  void reserve(std::size_t capacity) {
    assert(!fixed_);
    if (capacity <= capacity_) {
      return;
    }
    assert(capacity <= max_size());
    std::size_t new_capacity = capacity;
    if (capacity_ > (capacity / 2)) {
      if (capacity_ > (max_size() / 2)) {
        new_capacity = max_size();
      } else {
        new_capacity = capacity_ * 2;
      }
    }
    realloc(new_capacity);
  }

  void shrink() {
    MARISA_THROW_IF(fixed_, std::logic_error);
    if (size_ != capacity_) {
      realloc(size_);
    }
  }

  void fix() {
    MARISA_THROW_IF(fixed_, std::logic_error);
    fixed_ = true;
  }

  const T *begin() const {
    return const_objs_;
  }
  const T *end() const {
    return const_objs_ + size_;
  }
  const T &operator[](std::size_t i) const {
    assert(i < size_);
    return const_objs_[i];
  }
  const T &front() const {
    assert(size_ != 0);
    return const_objs_[0];
  }
  const T &back() const {
    assert(size_ != 0);
    return const_objs_[size_ - 1];
  }

  T *begin() {
    assert(!fixed_);
    return objs_;
  }
  T *end() {
    assert(!fixed_);
    return objs_ + size_;
  }
  T &operator[](std::size_t i) {
    assert(!fixed_);
    assert(i < size_);
    return objs_[i];
  }
  T &front() {
    assert(!fixed_);
    assert(size_ != 0);
    return objs_[0];
  }
  T &back() {
    assert(!fixed_);
    assert(size_ != 0);
    return objs_[size_ - 1];
  }

  std::size_t size() const {
    return size_;
  }
  std::size_t capacity() const {
    return capacity_;
  }
  bool fixed() const {
    return fixed_;
  }

  bool empty() const {
    return size_ == 0;
  }
  std::size_t total_size() const {
    return sizeof(T) * size_;
  }
  std::size_t io_size() const {
    return sizeof(uint64_t) + ((total_size() + 7) & ~0x07U);
  }

  void clear() noexcept {
    Vector().swap(*this);
  }
  void swap(Vector &rhs) noexcept {
    buf_.swap(rhs.buf_);
    std::swap(objs_, rhs.objs_);
    std::swap(const_objs_, rhs.const_objs_);
    std::swap(size_, rhs.size_);
    std::swap(capacity_, rhs.capacity_);
    std::swap(fixed_, rhs.fixed_);
  }

  static std::size_t max_size() {
    return SIZE_MAX / sizeof(T);
  }

 private:
  std::unique_ptr<char[]> buf_;
  T *objs_ = nullptr;
  const T *const_objs_ = nullptr;
  std::size_t size_ = 0;
  std::size_t capacity_ = 0;
  bool fixed_ = false;

  void map_(Mapper &mapper) {
    uint64_t total_size;
    mapper.map(&total_size);
    MARISA_THROW_IF(total_size > SIZE_MAX, std::runtime_error);
    MARISA_THROW_IF((total_size % sizeof(T)) != 0, std::runtime_error);
    const std::size_t size = static_cast<std::size_t>(total_size / sizeof(T));
    mapper.map(&const_objs_, size);
    mapper.seek(static_cast<std::size_t>((8 - (total_size % 8)) % 8));
    size_ = size;
    fix();
  }
  void read_(Reader &reader) {
    uint64_t total_size;
    reader.read(&total_size);
    MARISA_THROW_IF(total_size > SIZE_MAX, std::runtime_error);
    MARISA_THROW_IF((total_size % sizeof(T)) != 0, std::runtime_error);
    const std::size_t size = static_cast<std::size_t>(total_size / sizeof(T));
    resize(size);
    reader.read(objs_, size);
    reader.seek(static_cast<std::size_t>((8 - (total_size % 8)) % 8));
  }
  void write_(Writer &writer) const {
    writer.write(static_cast<uint64_t>(total_size()));
    writer.write(const_objs_, size_);
    writer.seek((8 - (total_size() % 8)) % 8);
  }

  // Copies current elements to new buffer of size `new_capacity`.
  // Requires `new_capacity >= size_`.
  void realloc(std::size_t new_capacity) {
    assert(new_capacity >= size_);
    assert(new_capacity <= max_size());

    std::unique_ptr<char[]> new_buf(new char[sizeof(T) * new_capacity]());
    T *new_objs = reinterpret_cast<T *>(new_buf.get());

    static_assert(std::is_trivially_copyable_v<T>);
    std::memcpy(new_objs, objs_, sizeof(T) * size_);

    buf_ = std::move(new_buf);
    objs_ = new_objs;
    const_objs_ = new_objs;
    capacity_ = new_capacity;
  }

  // copyInit() assumes that T's placement new does not throw an exception.
  // Requires the vector to be empty.
  void copyInit(const T *src, std::size_t size, std::size_t capacity) {
    assert(size_ == 0);

    std::unique_ptr<char[]> new_buf(new char[sizeof(T) * capacity]());
    T *new_objs = reinterpret_cast<T *>(new_buf.get());

    static_assert(std::is_trivially_copyable_v<T>);
    std::memcpy(new_objs, src, sizeof(T) * size);

    buf_ = std::move(new_buf);
    objs_ = new_objs;
    const_objs_ = new_objs;
    size_ = size;
    capacity_ = capacity;
  }
};

}  // namespace marisa::grimoire::vector

#endif  // MARISA_GRIMOIRE_VECTOR_VECTOR_H_
#line 1 "lib/marisa/grimoire/vector/bit-vector.h"
#ifndef MARISA_GRIMOIRE_VECTOR_BIT_VECTOR_H_
#define MARISA_GRIMOIRE_VECTOR_BIT_VECTOR_H_

#include <cassert>
#include <stdexcept>

//amalgamate marisa/grimoire/vector/rank-index.h
//amalgamate marisa/grimoire/vector/vector.h

namespace marisa::grimoire::vector {

class BitVector {
 public:
#if MARISA_WORD_SIZE == 64
  using Unit = uint64_t;
#else   // MARISA_WORD_SIZE == 64
  using Unit = uint32_t;
#endif  // MARISA_WORD_SIZE == 64

  BitVector() = default;

  BitVector(const BitVector &) = delete;
  BitVector &operator=(const BitVector &) = delete;

  void build(bool enables_select0, bool enables_select1) {
    BitVector temp;
    temp.build_index(*this, enables_select0, enables_select1);
    units_.shrink();
    temp.units_.swap(units_);
    swap(temp);
  }

  __attribute__((noinline)) void map(Mapper &mapper) {
    BitVector temp;
    temp.map_(mapper);
    swap(temp);
  }
  __attribute__((noinline)) void read(Reader &reader) {
    BitVector temp;
    temp.read_(reader);
    swap(temp);
  }
  __attribute__((noinline)) void write(Writer &writer) const {
    write_(writer);
  }

  void disable_select0() {
    select0s_.clear();
  }
  void disable_select1() {
    select1s_.clear();
  }

  void push_back(bool bit) {
    MARISA_THROW_IF(size_ == UINT32_MAX, std::length_error);
    if (size_ == (MARISA_WORD_SIZE * units_.size())) {
      units_.resize(units_.size() + (64 / MARISA_WORD_SIZE), 0);
    }
    if (bit) {
      units_[size_ / MARISA_WORD_SIZE] |= Unit{1} << (size_ % MARISA_WORD_SIZE);
      ++num_1s_;
    }
    ++size_;
  }

  bool operator[](std::size_t i) const {
    assert(i < size_);
    return (units_[i / MARISA_WORD_SIZE] &
            (Unit{1} << (i % MARISA_WORD_SIZE))) != 0;
  }

  std::size_t rank0(std::size_t i) const {
    assert(!ranks_.empty());
    assert(i <= size_);
    return i - rank1(i);
  }
  std::size_t rank1(std::size_t i) const;

  std::size_t select0(std::size_t i) const;
  std::size_t select1(std::size_t i) const;

  std::size_t num_0s() const {
    return size_ - num_1s_;
  }
  std::size_t num_1s() const {
    return num_1s_;
  }

  bool empty() const {
    return size_ == 0;
  }
  std::size_t size() const {
    return size_;
  }
  std::size_t total_size() const {
    return units_.total_size() + ranks_.total_size() + select0s_.total_size() +
           select1s_.total_size();
  }
  std::size_t io_size() const {
    return units_.io_size() + (sizeof(uint32_t) * 2) + ranks_.io_size() +
           select0s_.io_size() + select1s_.io_size();
  }

  void clear() noexcept {
    BitVector().swap(*this);
  }
  void swap(BitVector &rhs) noexcept {
    units_.swap(rhs.units_);
    std::swap(size_, rhs.size_);
    std::swap(num_1s_, rhs.num_1s_);
    ranks_.swap(rhs.ranks_);
    select0s_.swap(rhs.select0s_);
    select1s_.swap(rhs.select1s_);
  }

 private:
  Vector<Unit> units_;
  std::size_t size_ = 0;
  std::size_t num_1s_ = 0;
  Vector<RankIndex> ranks_;
  Vector<uint32_t> select0s_;
  Vector<uint32_t> select1s_;

  void build_index(const BitVector &bv, bool enables_select0,
                   bool enables_select1);

  void map_(Mapper &mapper) {
    units_.map(mapper);
    {
      uint32_t temp_size;
      mapper.map(&temp_size);
      size_ = temp_size;
    }
    {
      uint32_t temp_num_1s;
      mapper.map(&temp_num_1s);
      MARISA_THROW_IF(temp_num_1s > size_, std::runtime_error);
      num_1s_ = temp_num_1s;
    }
    ranks_.map(mapper);
    select0s_.map(mapper);
    select1s_.map(mapper);
  }

  void read_(Reader &reader) {
    units_.read(reader);
    {
      uint32_t temp_size;
      reader.read(&temp_size);
      size_ = temp_size;
    }
    {
      uint32_t temp_num_1s;
      reader.read(&temp_num_1s);
      MARISA_THROW_IF(temp_num_1s > size_, std::runtime_error);
      num_1s_ = temp_num_1s;
    }
    ranks_.read(reader);
    select0s_.read(reader);
    select1s_.read(reader);
  }

  void write_(Writer &writer) const {
    units_.write(writer);
    writer.write(static_cast<uint32_t>(size_));
    writer.write(static_cast<uint32_t>(num_1s_));
    ranks_.write(writer);
    select0s_.write(writer);
    select1s_.write(writer);
  }
};

}  // namespace marisa::grimoire::vector

#endif  // MARISA_GRIMOIRE_VECTOR_BIT_VECTOR_H_
#line 1 "lib/marisa/grimoire/vector/flat-vector.h"
#ifndef MARISA_GRIMOIRE_VECTOR_FLAT_VECTOR_H_
#define MARISA_GRIMOIRE_VECTOR_FLAT_VECTOR_H_

#include <cassert>
#include <stdexcept>

//amalgamate marisa/grimoire/vector/vector.h

namespace marisa::grimoire::vector {

class FlatVector {
 public:
#if MARISA_WORD_SIZE == 64
  using Unit = uint64_t;
#else   // MARISA_WORD_SIZE == 64
  using Unit = uint32_t;
#endif  // MARISA_WORD_SIZE == 64

  FlatVector() = default;

  FlatVector(const FlatVector &) = delete;
  FlatVector &operator=(const FlatVector &) = delete;

  void build(const Vector<uint32_t> &values) {
    FlatVector temp;
    temp.build_(values);
    swap(temp);
  }

  __attribute__((noinline)) void map(Mapper &mapper) {
    FlatVector temp;
    temp.map_(mapper);
    swap(temp);
  }
  __attribute__((noinline)) void read(Reader &reader) {
    FlatVector temp;
    temp.read_(reader);
    swap(temp);
  }
  __attribute__((noinline)) void write(Writer &writer) const {
    write_(writer);
  }

  uint32_t operator[](std::size_t i) const {
    assert(i < size_);

    const std::size_t pos = i * value_size_;
    const std::size_t unit_id = pos / MARISA_WORD_SIZE;
    const std::size_t unit_offset = pos % MARISA_WORD_SIZE;

    if ((unit_offset + value_size_) <= MARISA_WORD_SIZE) {
      return static_cast<uint32_t>(units_[unit_id] >> unit_offset) & mask_;
    } else {
      return static_cast<uint32_t>(
                 (units_[unit_id] >> unit_offset) |
                 (units_[unit_id + 1] << (MARISA_WORD_SIZE - unit_offset))) &
             mask_;
    }
  }

  std::size_t value_size() const {
    return value_size_;
  }
  uint32_t mask() const {
    return mask_;
  }

  bool empty() const {
    return size_ == 0;
  }
  std::size_t size() const {
    return size_;
  }
  std::size_t total_size() const {
    return units_.total_size();
  }
  std::size_t io_size() const {
    return units_.io_size() + (sizeof(uint32_t) * 2) + sizeof(uint64_t);
  }

  void clear() noexcept {
    FlatVector().swap(*this);
  }
  void swap(FlatVector &rhs) noexcept {
    units_.swap(rhs.units_);
    std::swap(value_size_, rhs.value_size_);
    std::swap(mask_, rhs.mask_);
    std::swap(size_, rhs.size_);
  }

 private:
  Vector<Unit> units_;
  std::size_t value_size_ = 0;
  uint32_t mask_ = 0;
  std::size_t size_ = 0;

  void build_(const Vector<uint32_t> &values) {
    uint32_t max_value = 0;
    for (std::size_t i = 0; i < values.size(); ++i) {
      if (values[i] > max_value) {
        max_value = values[i];
      }
    }

    std::size_t value_size = 0;
    while (max_value != 0) {
      ++value_size;
      max_value >>= 1;
    }

    std::size_t num_units = values.empty() ? 0 : (64 / MARISA_WORD_SIZE);
    if (value_size != 0) {
      num_units = static_cast<std::size_t>(
          ((static_cast<uint64_t>(value_size) * values.size()) +
           (MARISA_WORD_SIZE - 1)) /
          MARISA_WORD_SIZE);
      num_units += num_units % (64 / MARISA_WORD_SIZE);
    }

    units_.resize(num_units);
    if (num_units > 0) {
      units_.back() = 0;
    }

    value_size_ = value_size;
    if (value_size != 0) {
      mask_ = UINT32_MAX >> (32 - value_size);
    }
    size_ = values.size();

    for (std::size_t i = 0; i < values.size(); ++i) {
      set(i, values[i]);
    }
  }

  void map_(Mapper &mapper) {
    units_.map(mapper);
    {
      uint32_t temp_value_size;
      mapper.map(&temp_value_size);
      MARISA_THROW_IF(temp_value_size > 32, std::runtime_error);
      value_size_ = temp_value_size;
    }
    {
      uint32_t temp_mask;
      mapper.map(&temp_mask);
      mask_ = temp_mask;
    }
    {
      uint64_t temp_size;
      mapper.map(&temp_size);
      MARISA_THROW_IF(temp_size > SIZE_MAX, std::runtime_error);
      size_ = static_cast<std::size_t>(temp_size);
    }
  }

  void read_(Reader &reader) {
    units_.read(reader);
    {
      uint32_t temp_value_size;
      reader.read(&temp_value_size);
      MARISA_THROW_IF(temp_value_size > 32, std::runtime_error);
      value_size_ = temp_value_size;
    }
    {
      uint32_t temp_mask;
      reader.read(&temp_mask);
      mask_ = temp_mask;
    }
    {
      uint64_t temp_size;
      reader.read(&temp_size);
      MARISA_THROW_IF(temp_size > SIZE_MAX, std::runtime_error);
      size_ = static_cast<std::size_t>(temp_size);
    }
  }

  void write_(Writer &writer) const {
    units_.write(writer);
    writer.write(static_cast<uint32_t>(value_size_));
    writer.write(static_cast<uint32_t>(mask_));
    writer.write(static_cast<uint64_t>(size_));
  }

  void set(std::size_t i, uint32_t value) {
    assert(i < size_);
    assert(value <= mask_);

    const std::size_t pos = i * value_size_;
    const std::size_t unit_id = pos / MARISA_WORD_SIZE;
    const std::size_t unit_offset = pos % MARISA_WORD_SIZE;

    units_[unit_id] &= ~(static_cast<Unit>(mask_) << unit_offset);
    units_[unit_id] |= static_cast<Unit>(value & mask_) << unit_offset;
    if ((unit_offset + value_size_) > MARISA_WORD_SIZE) {
      units_[unit_id + 1] &=
          ~(static_cast<Unit>(mask_) >> (MARISA_WORD_SIZE - unit_offset));
      units_[unit_id + 1] |=
          static_cast<Unit>(value & mask_) >> (MARISA_WORD_SIZE - unit_offset);
    }
  }
};

}  // namespace marisa::grimoire::vector

#endif  // MARISA_GRIMOIRE_VECTOR_FLAT_VECTOR_H_
#line 1 "lib/marisa/grimoire/vector.h"
#ifndef MARISA_GRIMOIRE_VECTOR_H_
#define MARISA_GRIMOIRE_VECTOR_H_

//amalgamate marisa/grimoire/vector/bit-vector.h
//amalgamate marisa/grimoire/vector/flat-vector.h
//amalgamate marisa/grimoire/vector/vector.h

namespace marisa::grimoire {

using vector::BitVector;
using vector::FlatVector;
using vector::Vector;

}  // namespace marisa::grimoire

#endif  // MARISA_GRIMOIRE_VECTOR_H_
#line 1 "lib/marisa/grimoire/trie/tail.h"
#ifndef MARISA_GRIMOIRE_TRIE_TAIL_H_
#define MARISA_GRIMOIRE_TRIE_TAIL_H_

#include <cassert>

#include "./../lib/marisa.h" //amalgamate marisa/agent.h
//amalgamate marisa/grimoire/trie/entry.h
//amalgamate marisa/grimoire/vector.h

namespace marisa::grimoire::trie {

class Tail {
 public:
  Tail();

  Tail(const Tail &) = delete;
  Tail &operator=(const Tail &) = delete;

  void build(Vector<Entry> &entries, Vector<uint32_t> *offsets, TailMode mode);

  __attribute__((noinline)) void map(Mapper &mapper);
  __attribute__((noinline)) void read(Reader &reader);
  __attribute__((noinline)) void write(Writer &writer) const;

  void restore(Agent &agent, std::size_t offset) const;
  bool match(Agent &agent, std::size_t offset) const;
  bool prefix_match(Agent &agent, std::size_t offset) const;

  const char &operator[](std::size_t offset) const {
    assert(offset < buf_.size());
    return buf_[offset];
  }

  TailMode mode() const {
    return end_flags_.empty() ? MARISA_TEXT_TAIL : MARISA_BINARY_TAIL;
  }

  bool empty() const {
    return buf_.empty();
  }
  std::size_t size() const {
    return buf_.size();
  }
  std::size_t total_size() const {
    return buf_.total_size() + end_flags_.total_size();
  }
  std::size_t io_size() const {
    return buf_.io_size() + end_flags_.io_size();
  }

  void clear() noexcept;
  void swap(Tail &rhs) noexcept;

 private:
  Vector<char> buf_;
  BitVector end_flags_;

  void build_(Vector<Entry> &entries, Vector<uint32_t> *offsets, TailMode mode);

  void map_(Mapper &mapper);
  void read_(Reader &reader);
  void write_(Writer &writer) const;
};

}  // namespace marisa::grimoire::trie

#endif  // MARISA_GRIMOIRE_TRIE_TAIL_H_
#line 1 "lib/marisa/grimoire/algorithm/sort.h"
#ifndef MARISA_GRIMOIRE_ALGORITHM_SORT_H_
#define MARISA_GRIMOIRE_ALGORITHM_SORT_H_

#include <cassert>

#include "./../lib/marisa.h" //amalgamate marisa/base.h

namespace marisa::grimoire::algorithm {
namespace details {

enum {
  MARISA_INSERTION_SORT_THRESHOLD = 10
};

template <typename T>
int get_label(const T &unit, std::size_t depth) {
  assert(depth <= unit.length());

  return (depth < unit.length()) ? int{static_cast<uint8_t>(unit[depth])} : -1;
}

template <typename T>
int median(const T &a, const T &b, const T &c, std::size_t depth) {
  const int x = get_label(a, depth);
  const int y = get_label(b, depth);
  const int z = get_label(c, depth);
  if (x < y) {
    if (y < z) {
      return y;
    } else if (x < z) {
      return z;
    }
    return x;
  } else if (x < z) {
    return x;
  } else if (y < z) {
    return z;
  }
  return y;
}

template <typename T>
int compare(const T &lhs, const T &rhs, std::size_t depth) {
  for (std::size_t i = depth; i < lhs.length(); ++i) {
    if (i == rhs.length()) {
      return 1;
    }
    if (lhs[i] != rhs[i]) {
      return static_cast<uint8_t>(lhs[i]) - static_cast<uint8_t>(rhs[i]);
    }
  }
  if (lhs.length() == rhs.length()) {
    return 0;
  }
  return (lhs.length() < rhs.length()) ? -1 : 1;
}

template <typename Iterator>
std::size_t insertion_sort(Iterator l, Iterator r, std::size_t depth) {
  assert(l <= r);

  std::size_t count = 1;
  for (Iterator i = l + 1; i < r; ++i) {
    int result = 0;
    for (Iterator j = i; j > l; --j) {
      result = compare(*(j - 1), *j, depth);
      if (result <= 0) {
        break;
      }
      std::swap(*(j - 1), *j);
    }
    if (result != 0) {
      ++count;
    }
  }
  return count;
}

template <typename Iterator>
std::size_t sort(Iterator l, Iterator r, std::size_t depth) {
  assert(l <= r);

  std::size_t count = 0;
  while ((r - l) > MARISA_INSERTION_SORT_THRESHOLD) {
    Iterator pl = l;
    Iterator pr = r;
    Iterator pivot_l = l;
    Iterator pivot_r = r;

    const int pivot = median(*l, *(l + (r - l) / 2), *(r - 1), depth);
    for (;;) {
      while (pl < pr) {
        const int label = get_label(*pl, depth);
        if (label > pivot) {
          break;
        } else if (label == pivot) {
          std::swap(*pl, *pivot_l);
          ++pivot_l;
        }
        ++pl;
      }
      while (pl < pr) {
        const int label = get_label(*--pr, depth);
        if (label < pivot) {
          break;
        } else if (label == pivot) {
          std::swap(*pr, *--pivot_r);
        }
      }
      if (pl >= pr) {
        break;
      }
      std::swap(*pl, *pr);
      ++pl;
    }
    while (pivot_l > l) {
      std::swap(*--pivot_l, *--pl);
    }
    while (pivot_r < r) {
      std::swap(*pivot_r, *pr);
      ++pivot_r;
      ++pr;
    }

    if (((pl - l) > (pr - pl)) || ((r - pr) > (pr - pl))) {
      if ((pr - pl) == 1) {
        ++count;
      } else if ((pr - pl) > 1) {
        if (pivot == -1) {
          ++count;
        } else {
          count += sort(pl, pr, depth + 1);
        }
      }

      if ((pl - l) < (r - pr)) {
        if ((pl - l) == 1) {
          ++count;
        } else if ((pl - l) > 1) {
          count += sort(l, pl, depth);
        }
        l = pr;
      } else {
        if ((r - pr) == 1) {
          ++count;
        } else if ((r - pr) > 1) {
          count += sort(pr, r, depth);
        }
        r = pl;
      }
    } else {
      if ((pl - l) == 1) {
        ++count;
      } else if ((pl - l) > 1) {
        count += sort(l, pl, depth);
      }

      if ((r - pr) == 1) {
        ++count;
      } else if ((r - pr) > 1) {
        count += sort(pr, r, depth);
      }

      l = pl, r = pr;
      if ((pr - pl) == 1) {
        ++count;
      } else if ((pr - pl) > 1) {
        if (pivot == -1) {
          l = r;
          ++count;
        } else {
          ++depth;
        }
      }
    }
  }

  if ((r - l) > 1) {
    count += insertion_sort(l, r, depth);
  }
  return count;
}

}  // namespace details

template <typename Iterator>
std::size_t sort(Iterator begin, Iterator end) {
  assert(begin <= end);
  return details::sort(begin, end, 0);
}

}  // namespace marisa::grimoire::algorithm

#endif  // MARISA_GRIMOIRE_ALGORITHM_SORT_H_
#line 1 "lib/marisa/grimoire/trie/history.h"
#ifndef MARISA_GRIMOIRE_TRIE_STATE_HISTORY_H_
#define MARISA_GRIMOIRE_TRIE_STATE_HISTORY_H_

#include <cassert>

#include "./../lib/marisa.h" //amalgamate marisa/base.h

namespace marisa::grimoire::trie {

class History {
 public:
  History() = default;

  void set_node_id(std::size_t node_id) {
    assert(node_id <= UINT32_MAX);
    node_id_ = static_cast<uint32_t>(node_id);
  }
  void set_louds_pos(std::size_t louds_pos) {
    assert(louds_pos <= UINT32_MAX);
    louds_pos_ = static_cast<uint32_t>(louds_pos);
  }
  void set_key_pos(std::size_t key_pos) {
    assert(key_pos <= UINT32_MAX);
    key_pos_ = static_cast<uint32_t>(key_pos);
  }
  void set_link_id(std::size_t link_id) {
    assert(link_id <= UINT32_MAX);
    link_id_ = static_cast<uint32_t>(link_id);
  }
  void set_key_id(std::size_t key_id) {
    assert(key_id <= UINT32_MAX);
    key_id_ = static_cast<uint32_t>(key_id);
  }

  std::size_t node_id() const {
    return node_id_;
  }
  std::size_t louds_pos() const {
    return louds_pos_;
  }
  std::size_t key_pos() const {
    return key_pos_;
  }
  std::size_t link_id() const {
    return link_id_;
  }
  std::size_t key_id() const {
    return key_id_;
  }

 private:
  uint32_t node_id_ = 0;
  uint32_t louds_pos_ = 0;
  uint32_t key_pos_ = 0;
  uint32_t link_id_ = MARISA_INVALID_LINK_ID;
  uint32_t key_id_ = MARISA_INVALID_KEY_ID;
};

}  // namespace marisa::grimoire::trie

#endif  // MARISA_GRIMOIRE_TRIE_STATE_HISTORY_H_
#line 1 "lib/marisa/grimoire/trie/state.h"
#ifndef MARISA_GRIMOIRE_TRIE_STATE_H_
#define MARISA_GRIMOIRE_TRIE_STATE_H_

#include <cassert>
#include <vector>

//amalgamate marisa/grimoire/trie/history.h

namespace marisa::grimoire::trie {

// A search agent has its internal state and the status codes are defined
// below.
enum StatusCode {
  MARISA_READY_TO_ALL,
  MARISA_READY_TO_COMMON_PREFIX_SEARCH,
  MARISA_READY_TO_PREDICTIVE_SEARCH,
  MARISA_END_OF_COMMON_PREFIX_SEARCH,
  MARISA_END_OF_PREDICTIVE_SEARCH,
};

class State {
 public:
  State() = default;

  State(const State &) = default;
  State &operator=(const State &) = default;
  State(State &&) noexcept = default;
  State &operator=(State &&) noexcept = default;

  void set_node_id(std::size_t node_id) {
    assert(node_id <= UINT32_MAX);
    node_id_ = static_cast<uint32_t>(node_id);
  }
  void set_query_pos(std::size_t query_pos) {
    assert(query_pos <= UINT32_MAX);
    query_pos_ = static_cast<uint32_t>(query_pos);
  }
  void set_history_pos(std::size_t history_pos) {
    assert(history_pos <= UINT32_MAX);
    history_pos_ = static_cast<uint32_t>(history_pos);
  }
  void set_status_code(StatusCode status_code) {
    status_code_ = status_code;
  }

  std::size_t node_id() const {
    return node_id_;
  }
  std::size_t query_pos() const {
    return query_pos_;
  }
  std::size_t history_pos() const {
    return history_pos_;
  }
  StatusCode status_code() const {
    return status_code_;
  }

  const std::vector<char> &key_buf() const {
    return key_buf_;
  }
  const std::vector<History> &history() const {
    return history_;
  }

  std::vector<char> &key_buf() {
    return key_buf_;
  }
  std::vector<History> &history() {
    return history_;
  }

  void reset() {
    status_code_ = MARISA_READY_TO_ALL;
  }

  void lookup_init() {
    node_id_ = 0;
    query_pos_ = 0;
    status_code_ = MARISA_READY_TO_ALL;
  }
  void reverse_lookup_init() {
    key_buf_.resize(0);
    key_buf_.reserve(32);
    status_code_ = MARISA_READY_TO_ALL;
  }
  void common_prefix_search_init() {
    node_id_ = 0;
    query_pos_ = 0;
    status_code_ = MARISA_READY_TO_COMMON_PREFIX_SEARCH;
  }
  void predictive_search_init() {
    key_buf_.resize(0);
    key_buf_.reserve(64);
    history_.resize(0);
    history_.reserve(4);
    node_id_ = 0;
    query_pos_ = 0;
    history_pos_ = 0;
    status_code_ = MARISA_READY_TO_PREDICTIVE_SEARCH;
  }

 private:
  std::vector<char> key_buf_;
  std::vector<History> history_;
  uint32_t node_id_ = 0;
  uint32_t query_pos_ = 0;
  uint32_t history_pos_ = 0;
  StatusCode status_code_ = MARISA_READY_TO_ALL;
};

}  // namespace marisa::grimoire::trie

#endif  // MARISA_GRIMOIRE_TRIE_STATE_H_
#line 1 "lib/marisa/grimoire/trie/cache.h"
#ifndef MARISA_GRIMOIRE_TRIE_CACHE_H_
#define MARISA_GRIMOIRE_TRIE_CACHE_H_

#include <cassert>
#include <cfloat>

#include "./../lib/marisa.h" //amalgamate marisa/base.h

namespace marisa::grimoire::trie {

class Cache {
 public:
  Cache() = default;
  Cache(const Cache &cache) = default;
  Cache &operator=(const Cache &cache) = default;

  void set_parent(std::size_t parent) {
    assert(parent <= UINT32_MAX);
    parent_ = static_cast<uint32_t>(parent);
  }
  void set_child(std::size_t child) {
    assert(child <= UINT32_MAX);
    child_ = static_cast<uint32_t>(child);
  }
  void set_base(uint8_t base) {
    union_.link = (union_.link & ~0xFFU) | base;
  }
  void set_extra(std::size_t extra) {
    assert(extra <= (UINT32_MAX >> 8));
    union_.link = static_cast<uint32_t>((union_.link & 0xFFU) | (extra << 8));
  }
  void set_weight(float weight) {
    union_.weight = weight;
  }

  std::size_t parent() const {
    return parent_;
  }
  std::size_t child() const {
    return child_;
  }
  uint8_t base() const {
    return static_cast<uint8_t>(union_.link & 0xFFU);
  }
  std::size_t extra() const {
    return union_.link >> 8;
  }
  char label() const {
    return static_cast<char>(base());
  }
  std::size_t link() const {
    return union_.link;
  }
  float weight() const {
    return union_.weight;
  }

 private:
  uint32_t parent_ = 0;
  uint32_t child_ = 0;
  union Union {
    uint32_t link;
    float weight = FLT_MIN;
  } union_;
};

}  // namespace marisa::grimoire::trie

#endif  // MARISA_GRIMOIRE_TRIE_CACHE_H_
#line 1 "lib/marisa/grimoire/trie/config.h"
#ifndef MARISA_GRIMOIRE_TRIE_CONFIG_H_
#define MARISA_GRIMOIRE_TRIE_CONFIG_H_

#include <stdexcept>

#include "./../lib/marisa.h" //amalgamate marisa/base.h

namespace marisa::grimoire::trie {

class Config {
 public:
  Config() = default;

  Config(const Config &) = delete;
  Config &operator=(const Config &) = delete;

  void parse(int config_flags) {
    Config temp;
    temp.parse_(config_flags);
    swap(temp);
  }

  int flags() const {
    return static_cast<int>(num_tries_) | tail_mode_ | node_order_;
  }

  std::size_t num_tries() const {
    return num_tries_;
  }
  CacheLevel cache_level() const {
    return cache_level_;
  }
  TailMode tail_mode() const {
    return tail_mode_;
  }
  NodeOrder node_order() const {
    return node_order_;
  }

  void clear() noexcept {
    Config().swap(*this);
  }
  void swap(Config &rhs) noexcept {
    std::swap(num_tries_, rhs.num_tries_);
    std::swap(cache_level_, rhs.cache_level_);
    std::swap(tail_mode_, rhs.tail_mode_);
    std::swap(node_order_, rhs.node_order_);
  }

 private:
  std::size_t num_tries_ = MARISA_DEFAULT_NUM_TRIES;
  CacheLevel cache_level_ = MARISA_DEFAULT_CACHE;
  TailMode tail_mode_ = MARISA_DEFAULT_TAIL;
  NodeOrder node_order_ = MARISA_DEFAULT_ORDER;

  void parse_(int config_flags) {
    MARISA_THROW_IF((config_flags & ~MARISA_CONFIG_MASK) != 0,
                    std::invalid_argument);

    parse_num_tries(config_flags);
    parse_cache_level(config_flags);
    parse_tail_mode(config_flags);
    parse_node_order(config_flags);
  }

  void parse_num_tries(int config_flags) {
    const int num_tries = config_flags & MARISA_NUM_TRIES_MASK;
    if (num_tries != 0) {
      num_tries_ = static_cast<std::size_t>(num_tries);
    }
  }

  void parse_cache_level(int config_flags) {
    switch (config_flags & MARISA_CACHE_LEVEL_MASK) {
      case 0: {
        cache_level_ = MARISA_DEFAULT_CACHE;
        break;
      }
      case MARISA_HUGE_CACHE: {
        cache_level_ = MARISA_HUGE_CACHE;
        break;
      }
      case MARISA_LARGE_CACHE: {
        cache_level_ = MARISA_LARGE_CACHE;
        break;
      }
      case MARISA_NORMAL_CACHE: {
        cache_level_ = MARISA_NORMAL_CACHE;
        break;
      }
      case MARISA_SMALL_CACHE: {
        cache_level_ = MARISA_SMALL_CACHE;
        break;
      }
      case MARISA_TINY_CACHE: {
        cache_level_ = MARISA_TINY_CACHE;
        break;
      }
      default: {
        MARISA_THROW(std::invalid_argument, "undefined cache level");
      }
    }
  }

  void parse_tail_mode(int config_flags) {
    switch (config_flags & MARISA_TAIL_MODE_MASK) {
      case 0: {
        tail_mode_ = MARISA_DEFAULT_TAIL;
        break;
      }
      case MARISA_TEXT_TAIL: {
        tail_mode_ = MARISA_TEXT_TAIL;
        break;
      }
      case MARISA_BINARY_TAIL: {
        tail_mode_ = MARISA_BINARY_TAIL;
        break;
      }
      default: {
        MARISA_THROW(std::invalid_argument, "undefined tail mode");
      }
    }
  }

  void parse_node_order(int config_flags) {
    switch (config_flags & MARISA_NODE_ORDER_MASK) {
      case 0: {
        node_order_ = MARISA_DEFAULT_ORDER;
        break;
      }
      case MARISA_LABEL_ORDER: {
        node_order_ = MARISA_LABEL_ORDER;
        break;
      }
      case MARISA_WEIGHT_ORDER: {
        node_order_ = MARISA_WEIGHT_ORDER;
        break;
      }
      default: {
        MARISA_THROW(std::invalid_argument, "undefined node order");
      }
    }
  }
};

}  // namespace marisa::grimoire::trie

#endif  // MARISA_GRIMOIRE_TRIE_CONFIG_H_
#line 1 "lib/marisa/grimoire/trie/key.h"
#ifndef MARISA_GRIMOIRE_TRIE_KEY_H_
#define MARISA_GRIMOIRE_TRIE_KEY_H_

#include <cassert>

#include "./../lib/marisa.h" //amalgamate marisa/base.h

namespace marisa::grimoire::trie {

class Key {
 public:
  Key() = default;
  Key(const Key &entry) = default;
  Key &operator=(const Key &entry) = default;

  char operator[](std::size_t i) const {
    assert(i < length_);
    return ptr_[i];
  }

  void substr(std::size_t pos, std::size_t length) {
    assert(pos <= length_);
    assert(length <= length_);
    assert(pos <= (length_ - length));
    ptr_ += pos;
    length_ = static_cast<uint32_t>(length);
  }

  void set_str(const char *ptr, std::size_t length) {
    assert((ptr != nullptr) || (length == 0));
    assert(length <= UINT32_MAX);
    ptr_ = ptr;
    length_ = static_cast<uint32_t>(length);
  }
  void set_weight(float weight) {
    union_.weight = weight;
  }
  void set_terminal(std::size_t terminal) {
    assert(terminal <= UINT32_MAX);
    union_.terminal = static_cast<uint32_t>(terminal);
  }
  void set_id(std::size_t id) {
    assert(id <= UINT32_MAX);
    id_ = static_cast<uint32_t>(id);
  }

  const char *ptr() const {
    return ptr_;
  }
  std::size_t length() const {
    return length_;
  }
  float weight() const {
    return union_.weight;
  }
  std::size_t terminal() const {
    return union_.terminal;
  }
  std::size_t id() const {
    return id_;
  }

 private:
  const char *ptr_ = nullptr;
  uint32_t length_ = 0;
  union Union {
    float weight;
    uint32_t terminal = 0;
  } union_;
  uint32_t id_ = 0;
};

inline bool operator==(const Key &lhs, const Key &rhs) {
  if (lhs.length() != rhs.length()) {
    return false;
  }
  for (std::size_t i = 0; i < lhs.length(); ++i) {
    if (lhs[i] != rhs[i]) {
      return false;
    }
  }
  return true;
}

inline bool operator!=(const Key &lhs, const Key &rhs) {
  return !(lhs == rhs);
}

inline bool operator<(const Key &lhs, const Key &rhs) {
  for (std::size_t i = 0; i < lhs.length(); ++i) {
    if (i == rhs.length()) {
      return false;
    }
    if (lhs[i] != rhs[i]) {
      return static_cast<uint8_t>(lhs[i]) < static_cast<uint8_t>(rhs[i]);
    }
  }
  return lhs.length() < rhs.length();
}

inline bool operator>(const Key &lhs, const Key &rhs) {
  return rhs < lhs;
}

class ReverseKey {
 public:
  ReverseKey() = default;
  ReverseKey(const ReverseKey &entry) = default;
  ReverseKey &operator=(const ReverseKey &entry) = default;

  char operator[](std::size_t i) const {
    assert(i < length_);
    return *(ptr_ - i - 1);
  }

  void substr(std::size_t pos, std::size_t length) {
    assert(pos <= length_);
    assert(length <= length_);
    assert(pos <= (length_ - length));
    ptr_ -= pos;
    length_ = static_cast<uint32_t>(length);
  }

  void set_str(const char *ptr, std::size_t length) {
    assert((ptr != nullptr) || (length == 0));
    assert(length <= UINT32_MAX);
    ptr_ = ptr + length;
    length_ = static_cast<uint32_t>(length);
  }
  void set_weight(float weight) {
    union_.weight = weight;
  }
  void set_terminal(std::size_t terminal) {
    assert(terminal <= UINT32_MAX);
    union_.terminal = static_cast<uint32_t>(terminal);
  }
  void set_id(std::size_t id) {
    assert(id <= UINT32_MAX);
    id_ = static_cast<uint32_t>(id);
  }

  const char *ptr() const {
    return ptr_ - length_;
  }
  std::size_t length() const {
    return length_;
  }
  float weight() const {
    return union_.weight;
  }
  std::size_t terminal() const {
    return union_.terminal;
  }
  std::size_t id() const {
    return id_;
  }

 private:
  const char *ptr_ = nullptr;
  uint32_t length_ = 0;
  union Union {
    float weight;
    uint32_t terminal = 0;
  } union_;
  uint32_t id_ = 0;
};

inline bool operator==(const ReverseKey &lhs, const ReverseKey &rhs) {
  if (lhs.length() != rhs.length()) {
    return false;
  }
  for (std::size_t i = 0; i < lhs.length(); ++i) {
    if (lhs[i] != rhs[i]) {
      return false;
    }
  }
  return true;
}

inline bool operator!=(const ReverseKey &lhs, const ReverseKey &rhs) {
  return !(lhs == rhs);
}

inline bool operator<(const ReverseKey &lhs, const ReverseKey &rhs) {
  for (std::size_t i = 0; i < lhs.length(); ++i) {
    if (i == rhs.length()) {
      return false;
    }
    if (lhs[i] != rhs[i]) {
      return static_cast<uint8_t>(lhs[i]) < static_cast<uint8_t>(rhs[i]);
    }
  }
  return lhs.length() < rhs.length();
}

inline bool operator>(const ReverseKey &lhs, const ReverseKey &rhs) {
  return rhs < lhs;
}

}  // namespace marisa::grimoire::trie

#endif  // MARISA_GRIMOIRE_TRIE_KEY_H_
#line 1 "lib/marisa/grimoire/trie/louds-trie.h"
#ifndef MARISA_GRIMOIRE_TRIE_LOUDS_TRIE_H_
#define MARISA_GRIMOIRE_TRIE_LOUDS_TRIE_H_

#include <memory>

#include "./../lib/marisa.h" //amalgamate marisa/agent.h
//amalgamate marisa/grimoire/trie/cache.h
//amalgamate marisa/grimoire/trie/config.h
//amalgamate marisa/grimoire/trie/key.h
//amalgamate marisa/grimoire/trie/tail.h
//amalgamate marisa/grimoire/vector.h
#include "./../lib/marisa.h" //amalgamate marisa/keyset.h

namespace marisa::grimoire::trie {

class LoudsTrie {
 public:
  LoudsTrie();
  ~LoudsTrie();

  LoudsTrie(const LoudsTrie &) = delete;
  LoudsTrie &operator=(const LoudsTrie &) = delete;

  void build(Keyset &keyset, int flags);

  __attribute__((noinline)) void map(Mapper &mapper);
  __attribute__((noinline)) void read(Reader &reader);
  __attribute__((noinline)) void write(Writer &writer) const;

  bool lookup(Agent &agent) const;
  void reverse_lookup(Agent &agent) const;
  bool common_prefix_search(Agent &agent) const;
  bool predictive_search(Agent &agent) const;

  std::size_t num_tries() const {
    return config_.num_tries();
  }
  std::size_t num_keys() const {
    return size();
  }
  std::size_t num_nodes() const {
    return (louds_.size() / 2) - 1;
  }

  CacheLevel cache_level() const {
    return config_.cache_level();
  }
  TailMode tail_mode() const {
    return config_.tail_mode();
  }
  NodeOrder node_order() const {
    return config_.node_order();
  }

  bool empty() const {
    return size() == 0;
  }
  std::size_t size() const {
    return terminal_flags_.num_1s();
  }
  std::size_t total_size() const;
  std::size_t io_size() const;

  void clear() noexcept;
  void swap(LoudsTrie &rhs) noexcept;

 private:
  BitVector louds_;
  BitVector terminal_flags_;
  BitVector link_flags_;
  Vector<uint8_t> bases_;
  FlatVector extras_;
  Tail tail_;
  std::unique_ptr<LoudsTrie> next_trie_;
  Vector<Cache> cache_;
  std::size_t cache_mask_ = 0;
  std::size_t num_l1_nodes_ = 0;
  Config config_;
  Mapper mapper_;

  void build_(Keyset &keyset, const Config &config);

  template <typename T>
  void build_trie(Vector<T> &keys, Vector<uint32_t> *terminals,
                  const Config &config, std::size_t trie_id);
  template <typename T>
  void build_current_trie(Vector<T> &keys, Vector<uint32_t> *terminals,
                          const Config &config, std::size_t trie_id);
  template <typename T>
  void build_next_trie(Vector<T> &keys, Vector<uint32_t> *terminals,
                       const Config &config, std::size_t trie_id);
  template <typename T>
  void build_terminals(const Vector<T> &keys,
                       Vector<uint32_t> *terminals) const;

  void reserve_cache(const Config &config, std::size_t trie_id,
                     std::size_t num_keys);
  template <typename T>
  void cache(std::size_t parent, std::size_t child, float weight, char label);
  void fill_cache();

  void map_(Mapper &mapper);
  void read_(Reader &reader);
  void write_(Writer &writer) const;

  inline bool find_child(Agent &agent) const;
  inline bool predictive_find_child(Agent &agent) const;

  inline void restore(Agent &agent, std::size_t node_id) const;
  inline bool match(Agent &agent, std::size_t node_id) const;
  inline bool prefix_match(Agent &agent, std::size_t node_id) const;

  void restore_(Agent &agent, std::size_t node_id) const;
  bool match_(Agent &agent, std::size_t node_id) const;
  bool prefix_match_(Agent &agent, std::size_t node_id) const;

  inline std::size_t get_cache_id(std::size_t node_id, char label) const;
  inline std::size_t get_cache_id(std::size_t node_id) const;

  inline std::size_t get_link(std::size_t node_id) const;
  inline std::size_t get_link(std::size_t node_id, std::size_t link_id) const;

  inline std::size_t update_link_id(std::size_t link_id,
                                    std::size_t node_id) const;
};

}  // namespace marisa::grimoire::trie

#endif  // MARISA_GRIMOIRE_TRIE_LOUDS_TRIE_H_
#line 1 "lib/marisa/grimoire/trie/header.h"
#ifndef MARISA_GRIMOIRE_TRIE_HEADER_H_
#define MARISA_GRIMOIRE_TRIE_HEADER_H_

#include <stdexcept>

#include "../src/io.h" //amalgamate marisa/grimoire/io.h

namespace marisa::grimoire::trie {

class Header {
 public:
  enum {
    HEADER_SIZE = 16
  };

  Header() = default;

  Header(const Header &) = delete;
  Header &operator=(const Header &) = delete;

  __attribute__((noinline)) void map(Mapper &mapper) {
    const char *ptr;
    mapper.map(&ptr, HEADER_SIZE);
    MARISA_THROW_IF(!test_header(ptr), std::runtime_error);
  }
  __attribute__((noinline)) void read(Reader &reader) {
    char buf[HEADER_SIZE];
    reader.read(buf, HEADER_SIZE);
    MARISA_THROW_IF(!test_header(buf), std::runtime_error);
  }
  __attribute__((noinline)) void write(Writer &writer) const {
    writer.write(get_header(), HEADER_SIZE);
  }

  std::size_t io_size() const {
    return HEADER_SIZE;
  }

 private:
  static const char *get_header() {
    static const char buf[HEADER_SIZE] = "We love Marisa.";
    return buf;
  }

  static bool test_header(const char *ptr) {
    for (std::size_t i = 0; i < HEADER_SIZE; ++i) {
      if (ptr[i] != get_header()[i]) {
        return false;
      }
    }
    return true;
  }
};

}  // namespace marisa::grimoire::trie

#endif  // MARISA_GRIMOIRE_TRIE_HEADER_H_
#line 1 "lib/marisa/grimoire/trie/range.h"
#ifndef MARISA_GRIMOIRE_TRIE_RANGE_H_
#define MARISA_GRIMOIRE_TRIE_RANGE_H_

#include <cassert>

#include "./../lib/marisa.h" //amalgamate marisa/base.h

namespace marisa::grimoire::trie {

class Range {
 public:
  Range() = default;

  void set_begin(std::size_t begin) {
    assert(begin <= UINT32_MAX);
    begin_ = static_cast<uint32_t>(begin);
  }
  void set_end(std::size_t end) {
    assert(end <= UINT32_MAX);
    end_ = static_cast<uint32_t>(end);
  }
  void set_key_pos(std::size_t key_pos) {
    assert(key_pos <= UINT32_MAX);
    key_pos_ = static_cast<uint32_t>(key_pos);
  }

  std::size_t begin() const {
    return begin_;
  }
  std::size_t end() const {
    return end_;
  }
  std::size_t key_pos() const {
    return key_pos_;
  }

 private:
  uint32_t begin_ = 0;
  uint32_t end_ = 0;
  uint32_t key_pos_ = 0;
};

inline Range make_range(std::size_t begin, std::size_t end,
                        std::size_t key_pos) {
  Range range;
  range.set_begin(begin);
  range.set_end(end);
  range.set_key_pos(key_pos);
  return range;
}

class WeightedRange {
 public:
  WeightedRange() = default;

  void set_range(const Range &range) {
    range_ = range;
  }
  void set_begin(std::size_t begin) {
    range_.set_begin(begin);
  }
  void set_end(std::size_t end) {
    range_.set_end(end);
  }
  void set_key_pos(std::size_t key_pos) {
    range_.set_key_pos(key_pos);
  }
  void set_weight(float weight) {
    weight_ = weight;
  }

  const Range &range() const {
    return range_;
  }
  std::size_t begin() const {
    return range_.begin();
  }
  std::size_t end() const {
    return range_.end();
  }
  std::size_t key_pos() const {
    return range_.key_pos();
  }
  float weight() const {
    return weight_;
  }

 private:
  Range range_;
  float weight_ = 0.0F;
};

inline bool operator<(const WeightedRange &lhs, const WeightedRange &rhs) {
  return lhs.weight() < rhs.weight();
}

inline bool operator>(const WeightedRange &lhs, const WeightedRange &rhs) {
  return lhs.weight() > rhs.weight();
}

inline WeightedRange make_weighted_range(std::size_t begin, std::size_t end,
                                         std::size_t key_pos, float weight) {
  WeightedRange range;
  range.set_begin(begin);
  range.set_end(end);
  range.set_key_pos(key_pos);
  range.set_weight(weight);
  return range;
}

}  // namespace marisa::grimoire::trie

#endif  // MARISA_GRIMOIRE_TRIE_RANGE_H_
#line 1 "lib/marisa/grimoire/vector/pop-count.h"
#ifndef MARISA_GRIMOIRE_VECTOR_POP_COUNT_H_
#define MARISA_GRIMOIRE_VECTOR_POP_COUNT_H_

#if __cplusplus >= 202002L
 #include <bit>
#endif

#include "../src/intrin.h" //amalgamate marisa/grimoire/intrin.h

namespace marisa::grimoire::vector {

#if defined(__cpp_lib_bitops) && __cpp_lib_bitops >= 201907L

inline std::size_t popcount(uint64_t x) {
  return static_cast<std::size_t>(std::popcount(x));
}

#else  // c++17

 #ifdef __has_builtin
  #define MARISA_HAS_BUILTIN(x) __has_builtin(x)
 #else
  #define MARISA_HAS_BUILTIN(x) 0
 #endif

 #if MARISA_WORD_SIZE == 64

inline std::size_t popcount(uint64_t x) {
  #if MARISA_HAS_BUILTIN(__builtin_popcountll)
  static_assert(sizeof(x) == sizeof(unsigned long long),
                "__builtin_popcountll does not take 64-bit arg");
  return __builtin_popcountll(x);
  #elif defined(MARISA_X64) && defined(MARISA_USE_POPCNT)
   #ifdef _MSC_VER
  return __popcnt64(x);
   #else   // _MSC_VER
  return static_cast<std::size_t>(_mm_popcnt_u64(x));
   #endif  // _MSC_VER
  #elif defined(MARISA_AARCH64)
  // Byte-wise popcount followed by horizontal add.
  return vaddv_u8(vcnt_u8(vcreate_u8(x)));
  #else   // defined(MARISA_AARCH64)
  x = (x & 0x5555555555555555ULL) + ((x & 0xAAAAAAAAAAAAAAAAULL) >> 1);
  x = (x & 0x3333333333333333ULL) + ((x & 0xCCCCCCCCCCCCCCCCULL) >> 2);
  x = (x & 0x0F0F0F0F0F0F0F0FULL) + ((x & 0xF0F0F0F0F0F0F0F0ULL) >> 4);
  x *= 0x0101010101010101ULL;
  return x >> 56;
  #endif  // defined(MARISA_AARCH64)
}

 #else  // MARISA_WORD_SIZE == 64

inline std::size_t popcount(uint32_t x) {
  #if MARISA_HAS_BUILTIN(__builtin_popcount)
  static_assert(sizeof(x) == sizeof(unsigned int),
                "__builtin_popcount does not take 32-bit arg");
  return __builtin_popcount(x);
  #elif defined(MARISA_USE_POPCNT)
   #ifdef _MSC_VER
  return __popcnt(x);
   #else   // _MSC_VER
  return _mm_popcnt_u32(x);
   #endif  // _MSC_VER
  #else    // MARISA_USE_POPCNT
  x = (x & 0x55555555U) + ((x & 0xAAAAAAAAU) >> 1);
  x = (x & 0x33333333U) + ((x & 0xCCCCCCCCU) >> 2);
  x = (x & 0x0F0F0F0FU) + ((x & 0xF0F0F0F0U) >> 4);
  x *= 0x01010101U;
  return x >> 24;
  #endif   // MARISA_USE_POPCNT
}

 #endif  // MARISA_WORD_SIZE == 64

 #undef MARISA_HAS_BUILTIN

#endif  // c++17

}  // namespace marisa::grimoire::vector

#endif  // MARISA_GRIMOIRE_VECTOR_POP_COUNT_H_
#line 1 "lib/marisa/grimoire/trie.h"
#ifndef MARISA_GRIMOIRE_TRIE_H_
#define MARISA_GRIMOIRE_TRIE_H_

//amalgamate marisa/grimoire/trie/louds-trie.h
//amalgamate marisa/grimoire/trie/state.h

namespace marisa::grimoire {

using trie::LoudsTrie;
using trie::State;

}  // namespace marisa::grimoire

#endif  // MARISA_GRIMOIRE_TRIE_H_
#line 1 "lib/marisa/grimoire/trie/tail.cc"
//amalgamate marisa/grimoire/trie/tail.h

#include <cassert>
#include <stdexcept>

//amalgamate marisa/grimoire/algorithm/sort.h
//amalgamate marisa/grimoire/trie/state.h

namespace marisa::grimoire::trie {

Tail::Tail() = default;

void Tail::build(Vector<Entry> &entries, Vector<uint32_t> *offsets,
                 TailMode mode) {
  MARISA_THROW_IF(offsets == nullptr, std::invalid_argument);

  switch (mode) {
    case MARISA_TEXT_TAIL: {
      for (std::size_t i = 0; i < entries.size(); ++i) {
        const char *const ptr = entries[i].ptr();
        const std::size_t length = entries[i].length();
        for (std::size_t j = 0; j < length; ++j) {
          if (ptr[j] == '\0') {
            mode = MARISA_BINARY_TAIL;
            break;
          }
        }
        if (mode == MARISA_BINARY_TAIL) {
          break;
        }
      }
      break;
    }
    case MARISA_BINARY_TAIL: {
      break;
    }
    default: {
      MARISA_THROW(std::invalid_argument, "undefined tail mode");
    }
  }

  Tail temp;
  temp.build_(entries, offsets, mode);
  swap(temp);
}

void Tail::map(Mapper &mapper) {
  Tail temp;
  temp.map_(mapper);
  swap(temp);
}

void Tail::read(Reader &reader) {
  Tail temp;
  temp.read_(reader);
  swap(temp);
}

void Tail::write(Writer &writer) const {
  write_(writer);
}

void Tail::restore(Agent &agent, std::size_t offset) const {
  assert(!buf_.empty());

  State &state = agent.state();
  if (end_flags_.empty()) {
    for (const char *ptr = &buf_[offset]; *ptr != '\0'; ++ptr) {
      state.key_buf().push_back(*ptr);
    }
  } else {
    do {
      state.key_buf().push_back(buf_[offset]);
    } while (!end_flags_[offset++]);
  }
}

bool Tail::match(Agent &agent, std::size_t offset) const {
  assert(!buf_.empty());
  assert(agent.state().query_pos() < agent.query().length());

  State &state = agent.state();
  if (end_flags_.empty()) {
    const char *const ptr = &buf_[offset] - state.query_pos();
    do {
      if (ptr[state.query_pos()] != agent.query()[state.query_pos()]) {
        return false;
      }
      state.set_query_pos(state.query_pos() + 1);
      if (ptr[state.query_pos()] == '\0') {
        return true;
      }
    } while (state.query_pos() < agent.query().length());
    return false;
  }

  do {
    if (buf_[offset] != agent.query()[state.query_pos()]) {
      return false;
    }
    state.set_query_pos(state.query_pos() + 1);
    if (end_flags_[offset++]) {
      return true;
    }
  } while (state.query_pos() < agent.query().length());
  return false;
}

bool Tail::prefix_match(Agent &agent, std::size_t offset) const {
  assert(!buf_.empty());

  State &state = agent.state();
  if (end_flags_.empty()) {
    const char *ptr = &buf_[offset] - state.query_pos();
    do {
      if (ptr[state.query_pos()] != agent.query()[state.query_pos()]) {
        return false;
      }
      state.key_buf().push_back(ptr[state.query_pos()]);
      state.set_query_pos(state.query_pos() + 1);
      if (ptr[state.query_pos()] == '\0') {
        return true;
      }
    } while (state.query_pos() < agent.query().length());
    ptr += state.query_pos();
    do {
      state.key_buf().push_back(*ptr);
    } while (*++ptr != '\0');
    return true;
  }

  do {
    if (buf_[offset] != agent.query()[state.query_pos()]) {
      return false;
    }
    state.key_buf().push_back(buf_[offset]);
    state.set_query_pos(state.query_pos() + 1);
    if (end_flags_[offset++]) {
      return true;
    }
  } while (state.query_pos() < agent.query().length());
  do {
    state.key_buf().push_back(buf_[offset]);
  } while (!end_flags_[offset++]);
  return true;
}

void Tail::clear() noexcept {
  Tail().swap(*this);
}

void Tail::swap(Tail &rhs) noexcept {
  buf_.swap(rhs.buf_);
  end_flags_.swap(rhs.end_flags_);
}

void Tail::build_(Vector<Entry> &entries, Vector<uint32_t> *offsets,
                  TailMode mode) {
  for (std::size_t i = 0; i < entries.size(); ++i) {
    entries[i].set_id(i);
  }
  algorithm::sort(entries.begin(), entries.end());

  Vector<uint32_t> temp_offsets;
  temp_offsets.resize(entries.size(), 0);

  const Entry dummy;
  const Entry *last = &dummy;
  for (std::size_t i = entries.size(); i > 0; --i) {
    const Entry &current = entries[i - 1];
    MARISA_THROW_IF(current.length() == 0, std::out_of_range);
    std::size_t match = 0;
    while ((match < current.length()) && (match < last->length()) &&
           ((*last)[match] == current[match])) {
      ++match;
    }
    if ((match == current.length()) && (last->length() != 0)) {
      temp_offsets[current.id()] = static_cast<uint32_t>(
          temp_offsets[last->id()] + (last->length() - match));
    } else {
      temp_offsets[current.id()] = static_cast<uint32_t>(buf_.size());
      for (std::size_t j = 1; j <= current.length(); ++j) {
        buf_.push_back(current[current.length() - j]);
      }
      if (mode == MARISA_TEXT_TAIL) {
        buf_.push_back('\0');
      } else {
        for (std::size_t j = 1; j < current.length(); ++j) {
          end_flags_.push_back(false);
        }
        end_flags_.push_back(true);
      }
      MARISA_THROW_IF(buf_.size() > UINT32_MAX, std::length_error);
    }
    last = &current;
  }
  buf_.shrink();

  offsets->swap(temp_offsets);
}

void Tail::map_(Mapper &mapper) {
  buf_.map(mapper);
  end_flags_.map(mapper);
}

void Tail::read_(Reader &reader) {
  buf_.read(reader);
  end_flags_.read(reader);
}

void Tail::write_(Writer &writer) const {
  buf_.write(writer);
  end_flags_.write(writer);
}

}  // namespace marisa::grimoire::trie
#line 1 "lib/marisa/grimoire/trie/louds-trie.cc"
//amalgamate marisa/grimoire/trie/louds-trie.h

#include <algorithm>
#include <cassert>
#include <functional>
#include <queue>
#include <stdexcept>

//amalgamate marisa/grimoire/algorithm/sort.h
//amalgamate marisa/grimoire/trie/header.h
//amalgamate marisa/grimoire/trie/range.h
//amalgamate marisa/grimoire/trie/state.h

namespace marisa::grimoire::trie {

LoudsTrie::LoudsTrie() = default;

LoudsTrie::~LoudsTrie() = default;

void LoudsTrie::build(Keyset &keyset, int flags) {
  Config config;
  config.parse(flags);

  LoudsTrie temp;
  temp.build_(keyset, config);
  swap(temp);
}

void LoudsTrie::map(Mapper &mapper) {
  Header().map(mapper);

  LoudsTrie temp;
  temp.map_(mapper);
  temp.mapper_.swap(mapper);
  swap(temp);
}

void LoudsTrie::read(Reader &reader) {
  Header().read(reader);

  LoudsTrie temp;
  temp.read_(reader);
  swap(temp);
}

void LoudsTrie::write(Writer &writer) const {
  Header().write(writer);

  write_(writer);
}

bool LoudsTrie::lookup(Agent &agent) const {
  assert(agent.has_state());

  State &state = agent.state();
  state.lookup_init();
  while (state.query_pos() < agent.query().length()) {
    if (!find_child(agent)) {
      return false;
    }
  }
  if (!terminal_flags_[state.node_id()]) {
    return false;
  }
  agent.set_key(agent.query().ptr(), agent.query().length());
  agent.set_key(terminal_flags_.rank1(state.node_id()));
  return true;
}

void LoudsTrie::reverse_lookup(Agent &agent) const {
  assert(agent.has_state());
  MARISA_THROW_IF(agent.query().id() >= size(), std::out_of_range);

  State &state = agent.state();
  state.reverse_lookup_init();

  state.set_node_id(terminal_flags_.select1(agent.query().id()));
  if (state.node_id() == 0) {
    agent.set_key(state.key_buf().data(), state.key_buf().size());
    agent.set_key(agent.query().id());
    return;
  }
  for (;;) {
    if (link_flags_[state.node_id()]) {
      const std::size_t prev_key_pos = state.key_buf().size();
      restore(agent, get_link(state.node_id()));
      std::reverse(
          state.key_buf().begin() + static_cast<ptrdiff_t>(prev_key_pos),
          state.key_buf().end());
    } else {
      state.key_buf().push_back(static_cast<char>(bases_[state.node_id()]));
    }

    if (state.node_id() <= num_l1_nodes_) {
      std::reverse(state.key_buf().begin(), state.key_buf().end());
      agent.set_key(state.key_buf().data(), state.key_buf().size());
      agent.set_key(agent.query().id());
      return;
    }
    state.set_node_id(louds_.select1(state.node_id()) - state.node_id() - 1);
  }
}

bool LoudsTrie::common_prefix_search(Agent &agent) const {
  assert(agent.has_state());

  State &state = agent.state();
  if (state.status_code() == MARISA_END_OF_COMMON_PREFIX_SEARCH) {
    return false;
  }

  if (state.status_code() != MARISA_READY_TO_COMMON_PREFIX_SEARCH) {
    state.common_prefix_search_init();
    if (terminal_flags_[state.node_id()]) {
      agent.set_key(agent.query().ptr(), state.query_pos());
      agent.set_key(terminal_flags_.rank1(state.node_id()));
      return true;
    }
  }

  while (state.query_pos() < agent.query().length()) {
    if (!find_child(agent)) {
      state.set_status_code(MARISA_END_OF_COMMON_PREFIX_SEARCH);
      return false;
    }
    if (terminal_flags_[state.node_id()]) {
      agent.set_key(agent.query().ptr(), state.query_pos());
      agent.set_key(terminal_flags_.rank1(state.node_id()));
      return true;
    }
  }
  state.set_status_code(MARISA_END_OF_COMMON_PREFIX_SEARCH);
  return false;
}

bool LoudsTrie::predictive_search(Agent &agent) const {
  assert(agent.has_state());

  State &state = agent.state();
  if (state.status_code() == MARISA_END_OF_PREDICTIVE_SEARCH) {
    return false;
  }

  if (state.status_code() != MARISA_READY_TO_PREDICTIVE_SEARCH) {
    state.predictive_search_init();
    while (state.query_pos() < agent.query().length()) {
      if (!predictive_find_child(agent)) {
        state.set_status_code(MARISA_END_OF_PREDICTIVE_SEARCH);
        return false;
      }
    }

    History history;
    history.set_node_id(state.node_id());
    history.set_key_pos(state.key_buf().size());
    state.history().push_back(history);
    state.set_history_pos(1);

    if (terminal_flags_[state.node_id()]) {
      agent.set_key(state.key_buf().data(), state.key_buf().size());
      agent.set_key(terminal_flags_.rank1(state.node_id()));
      return true;
    }
  }

  for (;;) {
    if (state.history_pos() == state.history().size()) {
      const History &current = state.history().back();
      History next;
      next.set_louds_pos(louds_.select0(current.node_id()) + 1);
      next.set_node_id(next.louds_pos() - current.node_id() - 1);
      state.history().push_back(next);
    }

    History &next = state.history()[state.history_pos()];
    const bool link_flag = louds_[next.louds_pos()];
    next.set_louds_pos(next.louds_pos() + 1);
    if (link_flag) {
      state.set_history_pos(state.history_pos() + 1);
      if (link_flags_[next.node_id()]) {
        next.set_link_id(update_link_id(next.link_id(), next.node_id()));
        restore(agent, get_link(next.node_id(), next.link_id()));
      } else {
        state.key_buf().push_back(static_cast<char>(bases_[next.node_id()]));
      }
      next.set_key_pos(state.key_buf().size());

      if (terminal_flags_[next.node_id()]) {
        if (next.key_id() == MARISA_INVALID_KEY_ID) {
          next.set_key_id(terminal_flags_.rank1(next.node_id()));
        } else {
          next.set_key_id(next.key_id() + 1);
        }
        agent.set_key(state.key_buf().data(), state.key_buf().size());
        agent.set_key(next.key_id());
        return true;
      }
    } else if (state.history_pos() != 1) {
      History &current = state.history()[state.history_pos() - 1];
      current.set_node_id(current.node_id() + 1);
      const History &prev = state.history()[state.history_pos() - 2];
      state.key_buf().resize(prev.key_pos());
      state.set_history_pos(state.history_pos() - 1);
    } else {
      state.set_status_code(MARISA_END_OF_PREDICTIVE_SEARCH);
      return false;
    }
  }
}

std::size_t LoudsTrie::total_size() const {
  return louds_.total_size() + terminal_flags_.total_size() +
         link_flags_.total_size() + bases_.total_size() + extras_.total_size() +
         tail_.total_size() +
         ((next_trie_ != nullptr) ? next_trie_->total_size() : 0) +
         cache_.total_size();
}

std::size_t LoudsTrie::io_size() const {
  return Header().io_size() + louds_.io_size() + terminal_flags_.io_size() +
         link_flags_.io_size() + bases_.io_size() + extras_.io_size() +
         tail_.io_size() +
         ((next_trie_ != nullptr) ? (next_trie_->io_size() - Header().io_size())
                                  : 0) +
         cache_.io_size() + (sizeof(uint32_t) * 2);
}

void LoudsTrie::clear() noexcept {
  LoudsTrie().swap(*this);
}

void LoudsTrie::swap(LoudsTrie &rhs) noexcept {
  louds_.swap(rhs.louds_);
  terminal_flags_.swap(rhs.terminal_flags_);
  link_flags_.swap(rhs.link_flags_);
  bases_.swap(rhs.bases_);
  extras_.swap(rhs.extras_);
  tail_.swap(rhs.tail_);
  next_trie_.swap(rhs.next_trie_);
  cache_.swap(rhs.cache_);
  std::swap(cache_mask_, rhs.cache_mask_);
  std::swap(num_l1_nodes_, rhs.num_l1_nodes_);
  config_.swap(rhs.config_);
  mapper_.swap(rhs.mapper_);
}

void LoudsTrie::build_(Keyset &keyset, const Config &config) {
  Vector<Key> keys;
  keys.resize(keyset.size());
  for (std::size_t i = 0; i < keyset.size(); ++i) {
    keys[i].set_str(keyset[i].ptr(), keyset[i].length());
    keys[i].set_weight(keyset[i].weight());
  }

  Vector<uint32_t> terminals;
  build_trie(keys, &terminals, config, 1);

  using TerminalIdPair = std::pair<uint32_t, uint32_t>;
  const std::size_t pairs_size = terminals.size();
  std::unique_ptr<TerminalIdPair[]> pairs(new TerminalIdPair[pairs_size]);
  for (std::size_t i = 0; i < pairs_size; ++i) {
    pairs[i].first = terminals[i];
    pairs[i].second = static_cast<uint32_t>(i);
  }
  terminals.clear();
  std::sort(pairs.get(), pairs.get() + pairs_size);

  std::size_t node_id = 0;
  for (std::size_t i = 0; i < pairs_size; ++i) {
    while (node_id < pairs[i].first) {
      terminal_flags_.push_back(false);
      ++node_id;
    }
    if (node_id == pairs[i].first) {
      terminal_flags_.push_back(true);
      ++node_id;
    }
  }
  while (node_id < bases_.size()) {
    terminal_flags_.push_back(false);
    ++node_id;
  }
  terminal_flags_.push_back(false);
  terminal_flags_.build(false, true);

  for (std::size_t i = 0; i < keyset.size(); ++i) {
    keyset[pairs[i].second].set_id(terminal_flags_.rank1(pairs[i].first));
  }
}

template <typename T>
void LoudsTrie::build_trie(Vector<T> &keys, Vector<uint32_t> *terminals,
                           const Config &config, std::size_t trie_id) {
  build_current_trie(keys, terminals, config, trie_id);

  Vector<uint32_t> next_terminals;
  if (!keys.empty()) {
    build_next_trie(keys, &next_terminals, config, trie_id);
  }

  if (next_trie_ != nullptr) {
    config_.parse(static_cast<int>((next_trie_->num_tries() + 1)) |
                  next_trie_->tail_mode() | next_trie_->node_order());
  } else {
    config_.parse(1 | tail_.mode() | config.node_order() |
                  config.cache_level());
  }

  link_flags_.build(false, false);
  std::size_t node_id = 0;
  for (std::size_t i = 0; i < next_terminals.size(); ++i) {
    while (!link_flags_[node_id]) {
      ++node_id;
    }
    bases_[node_id] = static_cast<uint8_t>(next_terminals[i] % 256);
    next_terminals[i] /= 256;
    ++node_id;
  }
  extras_.build(next_terminals);
  fill_cache();
}

template <typename T>
void LoudsTrie::build_current_trie(Vector<T> &keys, Vector<uint32_t> *terminals,
                                   const Config &config, std::size_t trie_id) {
  for (std::size_t i = 0; i < keys.size(); ++i) {
    keys[i].set_id(i);
  }
  const std::size_t num_keys = algorithm::sort(keys.begin(), keys.end());
  reserve_cache(config, trie_id, num_keys);

  louds_.push_back(true);
  louds_.push_back(false);
  bases_.push_back('\0');
  link_flags_.push_back(false);

  Vector<T> next_keys;
  std::queue<Range> queue;
  Vector<WeightedRange> w_ranges;

  queue.push(make_range(0, keys.size(), 0));
  while (!queue.empty()) {
    const std::size_t node_id = link_flags_.size() - queue.size();

    Range range = queue.front();
    queue.pop();

    while ((range.begin() < range.end()) &&
           (keys[range.begin()].length() == range.key_pos())) {
      keys[range.begin()].set_terminal(node_id);
      range.set_begin(range.begin() + 1);
    }

    if (range.begin() == range.end()) {
      louds_.push_back(false);
      continue;
    }

    w_ranges.clear();
    double weight = double{keys[range.begin()].weight()};
    for (std::size_t i = range.begin() + 1; i < range.end(); ++i) {
      if (keys[i - 1][range.key_pos()] != keys[i][range.key_pos()]) {
        w_ranges.push_back(make_weighted_range(
            range.begin(), i, range.key_pos(), static_cast<float>(weight)));
        range.set_begin(i);
        weight = 0.0;
      }
      weight += double{keys[i].weight()};
    }
    w_ranges.push_back(make_weighted_range(range.begin(), range.end(),
                                           range.key_pos(),
                                           static_cast<float>(weight)));
    if (config.node_order() == MARISA_WEIGHT_ORDER) {
      std::stable_sort(w_ranges.begin(), w_ranges.end(),
                       std::greater<WeightedRange>());
    }

    if (node_id == 0) {
      num_l1_nodes_ = w_ranges.size();
    }

    for (std::size_t i = 0; i < w_ranges.size(); ++i) {
      WeightedRange &w_range = w_ranges[i];
      std::size_t key_pos = w_range.key_pos() + 1;
      while (key_pos < keys[w_range.begin()].length()) {
        std::size_t j;
        for (j = w_range.begin() + 1; j < w_range.end(); ++j) {
          if (keys[j - 1][key_pos] != keys[j][key_pos]) {
            break;
          }
        }
        if (j < w_range.end()) {
          break;
        }
        ++key_pos;
      }
      cache<T>(node_id, bases_.size(), w_range.weight(),
               keys[w_range.begin()][w_range.key_pos()]);

      if (key_pos == w_range.key_pos() + 1) {
        bases_.push_back(static_cast<unsigned char>(
            keys[w_range.begin()][w_range.key_pos()]));
        link_flags_.push_back(false);
      } else {
        bases_.push_back('\0');
        link_flags_.push_back(true);
        T next_key;
        next_key.set_str(keys[w_range.begin()].ptr(),
                         keys[w_range.begin()].length());
        next_key.substr(w_range.key_pos(), key_pos - w_range.key_pos());
        next_key.set_weight(w_range.weight());
        next_keys.push_back(next_key);
      }
      w_range.set_key_pos(key_pos);
      queue.push(w_range.range());
      louds_.push_back(true);
    }
    louds_.push_back(false);
  }

  louds_.push_back(false);
  louds_.build(trie_id == 1, true);
  bases_.shrink();

  build_terminals(keys, terminals);
  keys.swap(next_keys);
}

template <>
void LoudsTrie::build_next_trie(Vector<Key> &keys, Vector<uint32_t> *terminals,
                                const Config &config, std::size_t trie_id) {
  if (trie_id == config.num_tries()) {
    Vector<Entry> entries;
    entries.resize(keys.size());
    for (std::size_t i = 0; i < keys.size(); ++i) {
      entries[i].set_str(keys[i].ptr(), keys[i].length());
    }
    tail_.build(entries, terminals, config.tail_mode());
    return;
  }
  Vector<ReverseKey> reverse_keys;
  reverse_keys.resize(keys.size());
  for (std::size_t i = 0; i < keys.size(); ++i) {
    reverse_keys[i].set_str(keys[i].ptr(), keys[i].length());
    reverse_keys[i].set_weight(keys[i].weight());
  }
  keys.clear();
  next_trie_.reset(new LoudsTrie);
  next_trie_->build_trie(reverse_keys, terminals, config, trie_id + 1);
}

template <>
void LoudsTrie::build_next_trie(Vector<ReverseKey> &keys,
                                Vector<uint32_t> *terminals,
                                const Config &config, std::size_t trie_id) {
  if (trie_id == config.num_tries()) {
    Vector<Entry> entries;
    entries.resize(keys.size());
    for (std::size_t i = 0; i < keys.size(); ++i) {
      entries[i].set_str(keys[i].ptr(), keys[i].length());
    }
    tail_.build(entries, terminals, config.tail_mode());
    return;
  }
  next_trie_.reset(new LoudsTrie);
  next_trie_->build_trie(keys, terminals, config, trie_id + 1);
}

template <typename T>
void LoudsTrie::build_terminals(const Vector<T> &keys,
                                Vector<uint32_t> *terminals) const {
  Vector<uint32_t> temp;
  temp.resize(keys.size());
  for (std::size_t i = 0; i < keys.size(); ++i) {
    temp[keys[i].id()] = static_cast<uint32_t>(keys[i].terminal());
  }
  terminals->swap(temp);
}

template <>
void LoudsTrie::cache<Key>(std::size_t parent, std::size_t child, float weight,
                           char label) {
  assert(parent < child);

  const std::size_t cache_id = get_cache_id(parent, label);
  if (weight > cache_[cache_id].weight()) {
    cache_[cache_id].set_parent(parent);
    cache_[cache_id].set_child(child);
    cache_[cache_id].set_weight(weight);
  }
}

void LoudsTrie::reserve_cache(const Config &config, std::size_t trie_id,
                              std::size_t num_keys) {
  std::size_t cache_size = (trie_id == 1) ? 256 : 1;
  while (cache_size < (num_keys / config.cache_level())) {
    cache_size *= 2;
  }
  cache_.resize(cache_size);
  cache_mask_ = cache_size - 1;
}

template <>
void LoudsTrie::cache<ReverseKey>(std::size_t parent, std::size_t child,
                                  float weight, char) {
  assert(parent < child);

  const std::size_t cache_id = get_cache_id(child);
  if (weight > cache_[cache_id].weight()) {
    cache_[cache_id].set_parent(parent);
    cache_[cache_id].set_child(child);
    cache_[cache_id].set_weight(weight);
  }
}

void LoudsTrie::fill_cache() {
  for (std::size_t i = 0; i < cache_.size(); ++i) {
    const std::size_t node_id = cache_[i].child();
    if (node_id != 0) {
      cache_[i].set_base(bases_[node_id]);
      cache_[i].set_extra(!link_flags_[node_id]
                              ? MARISA_INVALID_EXTRA
                              : extras_[link_flags_.rank1(node_id)]);
    } else {
      cache_[i].set_parent(UINT32_MAX);
      cache_[i].set_child(UINT32_MAX);
    }
  }
}

void LoudsTrie::map_(Mapper &mapper) {
  louds_.map(mapper);
  terminal_flags_.map(mapper);
  link_flags_.map(mapper);
  bases_.map(mapper);
  extras_.map(mapper);
  tail_.map(mapper);
  if ((link_flags_.num_1s() != 0) && tail_.empty()) {
    next_trie_.reset(new LoudsTrie);
    next_trie_->map_(mapper);
  }
  cache_.map(mapper);
  cache_mask_ = cache_.size() - 1;
  {
    uint32_t temp_num_l1_nodes;
    mapper.map(&temp_num_l1_nodes);
    num_l1_nodes_ = temp_num_l1_nodes;
  }
  {
    uint32_t temp_config_flags;
    mapper.map(&temp_config_flags);
    config_.parse(static_cast<int>(temp_config_flags));
  }
}

void LoudsTrie::read_(Reader &reader) {
  louds_.read(reader);
  terminal_flags_.read(reader);
  link_flags_.read(reader);
  bases_.read(reader);
  extras_.read(reader);
  tail_.read(reader);
  if ((link_flags_.num_1s() != 0) && tail_.empty()) {
    next_trie_.reset(new LoudsTrie);
    next_trie_->read_(reader);
  }
  cache_.read(reader);
  cache_mask_ = cache_.size() - 1;
  {
    uint32_t temp_num_l1_nodes;
    reader.read(&temp_num_l1_nodes);
    num_l1_nodes_ = temp_num_l1_nodes;
  }
  {
    uint32_t temp_config_flags;
    reader.read(&temp_config_flags);
    config_.parse(static_cast<int>(temp_config_flags));
  }
}

void LoudsTrie::write_(Writer &writer) const {
  louds_.write(writer);
  terminal_flags_.write(writer);
  link_flags_.write(writer);
  bases_.write(writer);
  extras_.write(writer);
  tail_.write(writer);
  if (next_trie_ != nullptr) {
    next_trie_->write_(writer);
  }
  cache_.write(writer);
  writer.write(static_cast<uint32_t>(num_l1_nodes_));
  writer.write(static_cast<uint32_t>(config_.flags()));
}

bool LoudsTrie::find_child(Agent &agent) const {
  assert(agent.state().query_pos() < agent.query().length());

  State &state = agent.state();
  const std::size_t cache_id =
      get_cache_id(state.node_id(), agent.query()[state.query_pos()]);
  if (state.node_id() == cache_[cache_id].parent()) {
    if (cache_[cache_id].extra() != MARISA_INVALID_EXTRA) {
      if (!match(agent, cache_[cache_id].link())) {
        return false;
      }
    } else {
      state.set_query_pos(state.query_pos() + 1);
    }
    state.set_node_id(cache_[cache_id].child());
    return true;
  }

  std::size_t louds_pos = louds_.select0(state.node_id()) + 1;
  if (!louds_[louds_pos]) {
    return false;
  }
  state.set_node_id(louds_pos - state.node_id() - 1);
  std::size_t link_id = MARISA_INVALID_LINK_ID;
  do {
    if (link_flags_[state.node_id()]) {
      link_id = update_link_id(link_id, state.node_id());
      const std::size_t prev_query_pos = state.query_pos();
      if (match(agent, get_link(state.node_id(), link_id))) {
        return true;
      }
      if (state.query_pos() != prev_query_pos) {
        return false;
      }
    } else if (bases_[state.node_id()] ==
               static_cast<uint8_t>(agent.query()[state.query_pos()])) {
      state.set_query_pos(state.query_pos() + 1);
      return true;
    }
    state.set_node_id(state.node_id() + 1);
    ++louds_pos;
  } while (louds_[louds_pos]);
  return false;
}

bool LoudsTrie::predictive_find_child(Agent &agent) const {
  assert(agent.state().query_pos() < agent.query().length());

  State &state = agent.state();
  const std::size_t cache_id =
      get_cache_id(state.node_id(), agent.query()[state.query_pos()]);
  if (state.node_id() == cache_[cache_id].parent()) {
    if (cache_[cache_id].extra() != MARISA_INVALID_EXTRA) {
      if (!prefix_match(agent, cache_[cache_id].link())) {
        return false;
      }
    } else {
      state.key_buf().push_back(cache_[cache_id].label());
      state.set_query_pos(state.query_pos() + 1);
    }
    state.set_node_id(cache_[cache_id].child());
    return true;
  }

  std::size_t louds_pos = louds_.select0(state.node_id()) + 1;
  if (!louds_[louds_pos]) {
    return false;
  }
  state.set_node_id(louds_pos - state.node_id() - 1);
  std::size_t link_id = MARISA_INVALID_LINK_ID;
  do {
    if (link_flags_[state.node_id()]) {
      link_id = update_link_id(link_id, state.node_id());
      const std::size_t prev_query_pos = state.query_pos();
      if (prefix_match(agent, get_link(state.node_id(), link_id))) {
        return true;
      }
      if (state.query_pos() != prev_query_pos) {
        return false;
      }
    } else if (bases_[state.node_id()] ==
               static_cast<uint8_t>(agent.query()[state.query_pos()])) {
      state.key_buf().push_back(static_cast<char>(bases_[state.node_id()]));
      state.set_query_pos(state.query_pos() + 1);
      return true;
    }
    state.set_node_id(state.node_id() + 1);
    ++louds_pos;
  } while (louds_[louds_pos]);
  return false;
}

void LoudsTrie::restore(Agent &agent, std::size_t link) const {
  if (next_trie_ != nullptr) {
    next_trie_->restore_(agent, link);
  } else {
    tail_.restore(agent, link);
  }
}

bool LoudsTrie::match(Agent &agent, std::size_t link) const {
  if (next_trie_ != nullptr) {
    return next_trie_->match_(agent, link);
  }
  return tail_.match(agent, link);
}

bool LoudsTrie::prefix_match(Agent &agent, std::size_t link) const {
  if (next_trie_ != nullptr) {
    return next_trie_->prefix_match_(agent, link);
  }
  return tail_.prefix_match(agent, link);
}

void LoudsTrie::restore_(Agent &agent, std::size_t node_id) const {
  assert(node_id != 0);

  State &state = agent.state();
  for (;;) {
    const std::size_t cache_id = get_cache_id(node_id);
    if (node_id == cache_[cache_id].child()) {
      if (cache_[cache_id].extra() != MARISA_INVALID_EXTRA) {
        restore(agent, cache_[cache_id].link());
      } else {
        state.key_buf().push_back(cache_[cache_id].label());
      }

      node_id = cache_[cache_id].parent();
      if (node_id == 0) {
        return;
      }
      continue;
    }

    if (link_flags_[node_id]) {
      restore(agent, get_link(node_id));
    } else {
      state.key_buf().push_back(static_cast<char>(bases_[node_id]));
    }

    if (node_id <= num_l1_nodes_) {
      return;
    }
    node_id = louds_.select1(node_id) - node_id - 1;
  }
}

bool LoudsTrie::match_(Agent &agent, std::size_t node_id) const {
  assert(agent.state().query_pos() < agent.query().length());
  assert(node_id != 0);

  State &state = agent.state();
  for (;;) {
    const std::size_t cache_id = get_cache_id(node_id);
    if (node_id == cache_[cache_id].child()) {
      if (cache_[cache_id].extra() != MARISA_INVALID_EXTRA) {
        if (!match(agent, cache_[cache_id].link())) {
          return false;
        }
      } else if (cache_[cache_id].label() == agent.query()[state.query_pos()]) {
        state.set_query_pos(state.query_pos() + 1);
      } else {
        return false;
      }

      node_id = cache_[cache_id].parent();
      if (node_id == 0) {
        return true;
      }
      if (state.query_pos() >= agent.query().length()) {
        return false;
      }
      continue;
    }

    if (link_flags_[node_id]) {
      if (next_trie_ != nullptr) {
        if (!match(agent, get_link(node_id))) {
          return false;
        }
      } else if (!tail_.match(agent, get_link(node_id))) {
        return false;
      }
    } else if (bases_[node_id] ==
               static_cast<uint8_t>(agent.query()[state.query_pos()])) {
      state.set_query_pos(state.query_pos() + 1);
    } else {
      return false;
    }

    if (node_id <= num_l1_nodes_) {
      return true;
    }
    if (state.query_pos() >= agent.query().length()) {
      return false;
    }
    node_id = louds_.select1(node_id) - node_id - 1;
  }
}

bool LoudsTrie::prefix_match_(Agent &agent, std::size_t node_id) const {
  assert(agent.state().query_pos() < agent.query().length());
  assert(node_id != 0);

  State &state = agent.state();
  for (;;) {
    const std::size_t cache_id = get_cache_id(node_id);
    if (node_id == cache_[cache_id].child()) {
      if (cache_[cache_id].extra() != MARISA_INVALID_EXTRA) {
        if (!prefix_match(agent, cache_[cache_id].link())) {
          return false;
        }
      } else if (cache_[cache_id].label() == agent.query()[state.query_pos()]) {
        state.key_buf().push_back(cache_[cache_id].label());
        state.set_query_pos(state.query_pos() + 1);
      } else {
        return false;
      }

      node_id = cache_[cache_id].parent();
      if (node_id == 0) {
        return true;
      }
    } else {
      if (link_flags_[node_id]) {
        if (!prefix_match(agent, get_link(node_id))) {
          return false;
        }
      } else if (bases_[node_id] ==
                 static_cast<uint8_t>(agent.query()[state.query_pos()])) {
        state.key_buf().push_back(static_cast<char>(bases_[node_id]));
        state.set_query_pos(state.query_pos() + 1);
      } else {
        return false;
      }

      if (node_id <= num_l1_nodes_) {
        return true;
      }
      node_id = louds_.select1(node_id) - node_id - 1;
    }

    if (state.query_pos() >= agent.query().length()) {
      restore_(agent, node_id);
      return true;
    }
  }
}

std::size_t LoudsTrie::get_cache_id(std::size_t node_id, char label) const {
  return (node_id ^ (node_id << 5) ^ static_cast<uint8_t>(label)) & cache_mask_;
}

std::size_t LoudsTrie::get_cache_id(std::size_t node_id) const {
  return node_id & cache_mask_;
}

std::size_t LoudsTrie::get_link(std::size_t node_id) const {
  return bases_[node_id] | (extras_[link_flags_.rank1(node_id)] * 256);
}

std::size_t LoudsTrie::get_link(std::size_t node_id,
                                std::size_t link_id) const {
  return bases_[node_id] | (extras_[link_id] * 256);
}

std::size_t LoudsTrie::update_link_id(std::size_t link_id,
                                      std::size_t node_id) const {
  return (link_id == MARISA_INVALID_LINK_ID) ? link_flags_.rank1(node_id)
                                             : (link_id + 1);
}

}  // namespace marisa::grimoire::trie
#line 1 "lib/marisa/grimoire/vector/bit-vector.cc"
//amalgamate marisa/grimoire/vector/bit-vector.h

#include <algorithm>
#if __cplusplus >= 202002L
 #include <bit>
#endif
#include <cassert>

//amalgamate marisa/grimoire/vector/pop-count.h

namespace marisa::grimoire::vector {
namespace {

#if defined(__cpp_lib_bitops) && __cpp_lib_bitops >= 201907L

inline std::size_t countr_zero(uint64_t x) {
  return static_cast<std::size_t>(std::countr_zero(x));
}

#else  // c++17

inline std::size_t countr_zero(uint64_t x) {
 #ifdef _MSC_VER
  unsigned long pos;
  ::_BitScanForward64(&pos, x);
  return pos;
 #else   // _MSC_VER
  return __builtin_ctzll(x);
 #endif  // _MSC_VER
}

#endif  // c++17

#ifdef MARISA_USE_BMI2
inline std::size_t select_bit(std::size_t i, std::size_t bit_id,
                              uint64_t unit) {
  return bit_id + countr_zero(_pdep_u64(1ULL << i, unit));
}
#else  // MARISA_USE_BMI2
// clang-format off
const uint8_t SELECT_TABLE[8][256] = {
  {
    7, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0,
    4, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0,
    5, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0,
    4, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0,
    6, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0,
    4, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0,
    5, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0,
    4, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0,
    7, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0,
    4, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0,
    5, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0,
    4, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0,
    6, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0,
    4, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0,
    5, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0,
    4, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0
  },
  {
    7, 7, 7, 1, 7, 2, 2, 1, 7, 3, 3, 1, 3, 2, 2, 1,
    7, 4, 4, 1, 4, 2, 2, 1, 4, 3, 3, 1, 3, 2, 2, 1,
    7, 5, 5, 1, 5, 2, 2, 1, 5, 3, 3, 1, 3, 2, 2, 1,
    5, 4, 4, 1, 4, 2, 2, 1, 4, 3, 3, 1, 3, 2, 2, 1,
    7, 6, 6, 1, 6, 2, 2, 1, 6, 3, 3, 1, 3, 2, 2, 1,
    6, 4, 4, 1, 4, 2, 2, 1, 4, 3, 3, 1, 3, 2, 2, 1,
    6, 5, 5, 1, 5, 2, 2, 1, 5, 3, 3, 1, 3, 2, 2, 1,
    5, 4, 4, 1, 4, 2, 2, 1, 4, 3, 3, 1, 3, 2, 2, 1,
    7, 7, 7, 1, 7, 2, 2, 1, 7, 3, 3, 1, 3, 2, 2, 1,
    7, 4, 4, 1, 4, 2, 2, 1, 4, 3, 3, 1, 3, 2, 2, 1,
    7, 5, 5, 1, 5, 2, 2, 1, 5, 3, 3, 1, 3, 2, 2, 1,
    5, 4, 4, 1, 4, 2, 2, 1, 4, 3, 3, 1, 3, 2, 2, 1,
    7, 6, 6, 1, 6, 2, 2, 1, 6, 3, 3, 1, 3, 2, 2, 1,
    6, 4, 4, 1, 4, 2, 2, 1, 4, 3, 3, 1, 3, 2, 2, 1,
    6, 5, 5, 1, 5, 2, 2, 1, 5, 3, 3, 1, 3, 2, 2, 1,
    5, 4, 4, 1, 4, 2, 2, 1, 4, 3, 3, 1, 3, 2, 2, 1
  },
  {
    7, 7, 7, 7, 7, 7, 7, 2, 7, 7, 7, 3, 7, 3, 3, 2,
    7, 7, 7, 4, 7, 4, 4, 2, 7, 4, 4, 3, 4, 3, 3, 2,
    7, 7, 7, 5, 7, 5, 5, 2, 7, 5, 5, 3, 5, 3, 3, 2,
    7, 5, 5, 4, 5, 4, 4, 2, 5, 4, 4, 3, 4, 3, 3, 2,
    7, 7, 7, 6, 7, 6, 6, 2, 7, 6, 6, 3, 6, 3, 3, 2,
    7, 6, 6, 4, 6, 4, 4, 2, 6, 4, 4, 3, 4, 3, 3, 2,
    7, 6, 6, 5, 6, 5, 5, 2, 6, 5, 5, 3, 5, 3, 3, 2,
    6, 5, 5, 4, 5, 4, 4, 2, 5, 4, 4, 3, 4, 3, 3, 2,
    7, 7, 7, 7, 7, 7, 7, 2, 7, 7, 7, 3, 7, 3, 3, 2,
    7, 7, 7, 4, 7, 4, 4, 2, 7, 4, 4, 3, 4, 3, 3, 2,
    7, 7, 7, 5, 7, 5, 5, 2, 7, 5, 5, 3, 5, 3, 3, 2,
    7, 5, 5, 4, 5, 4, 4, 2, 5, 4, 4, 3, 4, 3, 3, 2,
    7, 7, 7, 6, 7, 6, 6, 2, 7, 6, 6, 3, 6, 3, 3, 2,
    7, 6, 6, 4, 6, 4, 4, 2, 6, 4, 4, 3, 4, 3, 3, 2,
    7, 6, 6, 5, 6, 5, 5, 2, 6, 5, 5, 3, 5, 3, 3, 2,
    6, 5, 5, 4, 5, 4, 4, 2, 5, 4, 4, 3, 4, 3, 3, 2
  },
  {
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 3,
    7, 7, 7, 7, 7, 7, 7, 4, 7, 7, 7, 4, 7, 4, 4, 3,
    7, 7, 7, 7, 7, 7, 7, 5, 7, 7, 7, 5, 7, 5, 5, 3,
    7, 7, 7, 5, 7, 5, 5, 4, 7, 5, 5, 4, 5, 4, 4, 3,
    7, 7, 7, 7, 7, 7, 7, 6, 7, 7, 7, 6, 7, 6, 6, 3,
    7, 7, 7, 6, 7, 6, 6, 4, 7, 6, 6, 4, 6, 4, 4, 3,
    7, 7, 7, 6, 7, 6, 6, 5, 7, 6, 6, 5, 6, 5, 5, 3,
    7, 6, 6, 5, 6, 5, 5, 4, 6, 5, 5, 4, 5, 4, 4, 3,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 3,
    7, 7, 7, 7, 7, 7, 7, 4, 7, 7, 7, 4, 7, 4, 4, 3,
    7, 7, 7, 7, 7, 7, 7, 5, 7, 7, 7, 5, 7, 5, 5, 3,
    7, 7, 7, 5, 7, 5, 5, 4, 7, 5, 5, 4, 5, 4, 4, 3,
    7, 7, 7, 7, 7, 7, 7, 6, 7, 7, 7, 6, 7, 6, 6, 3,
    7, 7, 7, 6, 7, 6, 6, 4, 7, 6, 6, 4, 6, 4, 4, 3,
    7, 7, 7, 6, 7, 6, 6, 5, 7, 6, 6, 5, 6, 5, 5, 3,
    7, 6, 6, 5, 6, 5, 5, 4, 6, 5, 5, 4, 5, 4, 4, 3
  },
  {
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 4,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 5,
    7, 7, 7, 7, 7, 7, 7, 5, 7, 7, 7, 5, 7, 5, 5, 4,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 6,
    7, 7, 7, 7, 7, 7, 7, 6, 7, 7, 7, 6, 7, 6, 6, 4,
    7, 7, 7, 7, 7, 7, 7, 6, 7, 7, 7, 6, 7, 6, 6, 5,
    7, 7, 7, 6, 7, 6, 6, 5, 7, 6, 6, 5, 6, 5, 5, 4,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 4,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 5,
    7, 7, 7, 7, 7, 7, 7, 5, 7, 7, 7, 5, 7, 5, 5, 4,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 6,
    7, 7, 7, 7, 7, 7, 7, 6, 7, 7, 7, 6, 7, 6, 6, 4,
    7, 7, 7, 7, 7, 7, 7, 6, 7, 7, 7, 6, 7, 6, 6, 5,
    7, 7, 7, 6, 7, 6, 6, 5, 7, 6, 6, 5, 6, 5, 5, 4
  },
  {
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 5,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 6,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 6,
    7, 7, 7, 7, 7, 7, 7, 6, 7, 7, 7, 6, 7, 6, 6, 5,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 5,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 6,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 6,
    7, 7, 7, 7, 7, 7, 7, 6, 7, 7, 7, 6, 7, 6, 6, 5
  },
  {
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 6,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 6
  },
  {
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
    7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7
  }
};
// clang-format on

 #if MARISA_WORD_SIZE == 64
constexpr uint64_t MASK_01 = 0x0101010101010101ULL;

// Pre-computed lookup table trick from Gog, Simon and Matthias Petri.
// "Optimized succinct data structures for massive data."  Software:
// Practice and Experience 44 (2014): 1287 - 1314.
// PREFIX_SUM_OVERFLOW[i] = (0x7F - i) * MASK_01.
const uint64_t PREFIX_SUM_OVERFLOW[64] = {
    // clang-format off
  0x7F * MASK_01, 0x7E * MASK_01, 0x7D * MASK_01, 0x7C * MASK_01,
  0x7B * MASK_01, 0x7A * MASK_01, 0x79 * MASK_01, 0x78 * MASK_01,
  0x77 * MASK_01, 0x76 * MASK_01, 0x75 * MASK_01, 0x74 * MASK_01,
  0x73 * MASK_01, 0x72 * MASK_01, 0x71 * MASK_01, 0x70 * MASK_01,

  0x6F * MASK_01, 0x6E * MASK_01, 0x6D * MASK_01, 0x6C * MASK_01,
  0x6B * MASK_01, 0x6A * MASK_01, 0x69 * MASK_01, 0x68 * MASK_01,
  0x67 * MASK_01, 0x66 * MASK_01, 0x65 * MASK_01, 0x64 * MASK_01,
  0x63 * MASK_01, 0x62 * MASK_01, 0x61 * MASK_01, 0x60 * MASK_01,

  0x5F * MASK_01, 0x5E * MASK_01, 0x5D * MASK_01, 0x5C * MASK_01,
  0x5B * MASK_01, 0x5A * MASK_01, 0x59 * MASK_01, 0x58 * MASK_01,
  0x57 * MASK_01, 0x56 * MASK_01, 0x55 * MASK_01, 0x54 * MASK_01,
  0x53 * MASK_01, 0x52 * MASK_01, 0x51 * MASK_01, 0x50 * MASK_01,

  0x4F * MASK_01, 0x4E * MASK_01, 0x4D * MASK_01, 0x4C * MASK_01,
  0x4B * MASK_01, 0x4A * MASK_01, 0x49 * MASK_01, 0x48 * MASK_01,
  0x47 * MASK_01, 0x46 * MASK_01, 0x45 * MASK_01, 0x44 * MASK_01,
  0x43 * MASK_01, 0x42 * MASK_01, 0x41 * MASK_01, 0x40 * MASK_01
    // clang-format on
};

std::size_t select_bit(std::size_t i, std::size_t bit_id, uint64_t unit) {
  uint64_t counts;
  {
  #if defined(MARISA_X64) && defined(MARISA_USE_SSSE3)
    __m128i lower_nibbles =
        _mm_cvtsi64_si128(static_cast<long long>(unit & 0x0F0F0F0F0F0F0F0FULL));
    __m128i upper_nibbles =
        _mm_cvtsi64_si128(static_cast<long long>(unit & 0xF0F0F0F0F0F0F0F0ULL));
    upper_nibbles = _mm_srli_epi32(upper_nibbles, 4);

    __m128i lower_counts =
        _mm_set_epi8(4, 3, 3, 2, 3, 2, 2, 1, 3, 2, 2, 1, 2, 1, 1, 0);
    lower_counts = _mm_shuffle_epi8(lower_counts, lower_nibbles);
    __m128i upper_counts =
        _mm_set_epi8(4, 3, 3, 2, 3, 2, 2, 1, 3, 2, 2, 1, 2, 1, 1, 0);
    upper_counts = _mm_shuffle_epi8(upper_counts, upper_nibbles);

    counts = static_cast<uint64_t>(
        _mm_cvtsi128_si64(_mm_add_epi8(lower_counts, upper_counts)));
  #elif defined(MARISA_AARCH64)
    // Byte-wise popcount using CNT (plus a lot of conversion noise).
    // This actually only requires NEON, not AArch64, but we are already
    // in a 64-bit `#ifdef`.
    counts = vget_lane_u64(vreinterpret_u64_u8(vcnt_u8(vcreate_u8(unit))), 0);
  #else   // defined(MARISA_AARCH64)
    constexpr uint64_t MASK_0F = 0x0F0F0F0F0F0F0F0FULL;
    constexpr uint64_t MASK_33 = 0x3333333333333333ULL;
    constexpr uint64_t MASK_55 = 0x5555555555555555ULL;
    counts = unit - ((unit >> 1) & MASK_55);
    counts = (counts & MASK_33) + ((counts >> 2) & MASK_33);
    counts = (counts + (counts >> 4)) & MASK_0F;
  #endif  // defined(MARISA_AARCH64)
    counts *= MASK_01;
  }

  #if defined(MARISA_X64) && defined(MARISA_USE_POPCNT)
  uint8_t skip;
  {
    __m128i x = _mm_cvtsi64_si128(static_cast<long long>((i + 1) * MASK_01));
    __m128i y = _mm_cvtsi64_si128(static_cast<long long>(counts));
    x = _mm_cmpgt_epi8(x, y);
    skip = (uint8_t)popcount(static_cast<uint64_t>(_mm_cvtsi128_si64(x)));
  }
  #else   // defined(MARISA_X64) && defined(MARISA_USE_POPCNT)
  constexpr uint64_t MASK_80 = 0x8080808080808080ULL;
  const uint64_t x = (counts + PREFIX_SUM_OVERFLOW[i]) & MASK_80;
  // We masked with `MASK_80`, so the first bit set is the high bit in the
  // byte, therefore `num_trailing_zeros == 8 * byte_nr + 7` and the byte
  // number is the number of trailing zeros divided by 8.  We just shift off
  // the low 7 bits, so `CTZ` gives us the `skip` value we want for the
  // number of bits of `counts` to shift.
  const int skip = countr_zero(x >> 7);
  #endif  // defined(MARISA_X64) && defined(MARISA_USE_POPCNT)

  bit_id += static_cast<std::size_t>(skip);
  unit >>= skip;
  i -= ((counts << 8) >> skip) & 0xFF;

  return bit_id + SELECT_TABLE[i][unit & 0xFF];
}
 #else    // MARISA_WORD_SIZE == 64
  #ifdef MARISA_USE_SSE2
// Popcount of the byte times eight.
const uint8_t POPCNT_X8_TABLE[256] = {
    // clang-format off
   0,  8,  8, 16,  8, 16, 16, 24,  8, 16, 16, 24, 16, 24, 24, 32,
   8, 16, 16, 24, 16, 24, 24, 32, 16, 24, 24, 32, 24, 32, 32, 40,
   8, 16, 16, 24, 16, 24, 24, 32, 16, 24, 24, 32, 24, 32, 32, 40,
  16, 24, 24, 32, 24, 32, 32, 40, 24, 32, 32, 40, 32, 40, 40, 48,
   8, 16, 16, 24, 16, 24, 24, 32, 16, 24, 24, 32, 24, 32, 32, 40,
  16, 24, 24, 32, 24, 32, 32, 40, 24, 32, 32, 40, 32, 40, 40, 48,
  16, 24, 24, 32, 24, 32, 32, 40, 24, 32, 32, 40, 32, 40, 40, 48,
  24, 32, 32, 40, 32, 40, 40, 48, 32, 40, 40, 48, 40, 48, 48, 56,
   8, 16, 16, 24, 16, 24, 24, 32, 16, 24, 24, 32, 24, 32, 32, 40,
  16, 24, 24, 32, 24, 32, 32, 40, 24, 32, 32, 40, 32, 40, 40, 48,
  16, 24, 24, 32, 24, 32, 32, 40, 24, 32, 32, 40, 32, 40, 40, 48,
  24, 32, 32, 40, 32, 40, 40, 48, 32, 40, 40, 48, 40, 48, 48, 56,
  16, 24, 24, 32, 24, 32, 32, 40, 24, 32, 32, 40, 32, 40, 40, 48,
  24, 32, 32, 40, 32, 40, 40, 48, 32, 40, 40, 48, 40, 48, 48, 56,
  24, 32, 32, 40, 32, 40, 40, 48, 32, 40, 40, 48, 40, 48, 48, 56,
  32, 40, 40, 48, 40, 48, 48, 56, 40, 48, 48, 56, 48, 56, 56, 64
    // clang-format on
};

std::size_t select_bit(std::size_t i, std::size_t bit_id, uint32_t unit_lo,
                       uint32_t unit_hi) {
  __m128i unit;
  {
    __m128i lower_dword = _mm_cvtsi32_si128(unit_lo);
    __m128i upper_dword = _mm_cvtsi32_si128(unit_hi);
    upper_dword = _mm_slli_si128(upper_dword, 4);
    unit = _mm_or_si128(lower_dword, upper_dword);
  }

  __m128i counts;
  {
   #ifdef MARISA_USE_SSSE3
    __m128i lower_nibbles = _mm_set1_epi8(0x0F);
    lower_nibbles = _mm_and_si128(lower_nibbles, unit);
    __m128i upper_nibbles = _mm_set1_epi8((uint8_t)0xF0);
    upper_nibbles = _mm_and_si128(upper_nibbles, unit);
    upper_nibbles = _mm_srli_epi32(upper_nibbles, 4);

    __m128i lower_counts =
        _mm_set_epi8(4, 3, 3, 2, 3, 2, 2, 1, 3, 2, 2, 1, 2, 1, 1, 0);
    lower_counts = _mm_shuffle_epi8(lower_counts, lower_nibbles);
    __m128i upper_counts =
        _mm_set_epi8(4, 3, 3, 2, 3, 2, 2, 1, 3, 2, 2, 1, 2, 1, 1, 0);
    upper_counts = _mm_shuffle_epi8(upper_counts, upper_nibbles);

    counts = _mm_add_epi8(lower_counts, upper_counts);
   #else   // MARISA_USE_SSSE3
    __m128i x = _mm_srli_epi32(unit, 1);
    x = _mm_and_si128(x, _mm_set1_epi8(0x55));
    x = _mm_sub_epi8(unit, x);

    __m128i y = _mm_srli_epi32(x, 2);
    y = _mm_and_si128(y, _mm_set1_epi8(0x33));
    x = _mm_and_si128(x, _mm_set1_epi8(0x33));
    x = _mm_add_epi8(x, y);

    y = _mm_srli_epi32(x, 4);
    x = _mm_add_epi8(x, y);
    counts = _mm_and_si128(x, _mm_set1_epi8(0x0F));
   #endif  // MARISA_USE_SSSE3
  }

  __m128i accumulated_counts;
  {
    __m128i x = counts;
    x = _mm_slli_si128(x, 1);
    __m128i y = counts;
    y = _mm_add_epi32(y, x);

    x = y;
    y = _mm_slli_si128(y, 2);
    x = _mm_add_epi32(x, y);

    y = x;
    x = _mm_slli_si128(x, 4);
    y = _mm_add_epi32(y, x);

    accumulated_counts = _mm_set_epi32(0x7F7F7F7FU, 0x7F7F7F7FU, 0, 0);
    accumulated_counts = _mm_or_si128(accumulated_counts, y);
  }

  uint8_t skip;
  {
    __m128i x = _mm_set1_epi8((uint8_t)(i + 1));
    x = _mm_cmpgt_epi8(x, accumulated_counts);
    // Since we use `_mm_movemask_epi8`, to move the top bit of every byte,
    // popcount times eight gives the original popcount of `x` before the
    // movemask.  (`_mm_cmpgt_epi8` sets all bits in a byte to 0 or 1.)
    skip = POPCNT_X8_TABLE[_mm_movemask_epi8(x)];
  }

  uint8_t byte;
  {
    alignas(16) uint8_t unit_bytes[16];
    alignas(16) uint8_t accumulated_counts_bytes[16];
    accumulated_counts = _mm_slli_si128(accumulated_counts, 1);
    _mm_store_si128(reinterpret_cast<__m128i *>(unit_bytes), unit);
    _mm_store_si128(reinterpret_cast<__m128i *>(accumulated_counts_bytes),
                    accumulated_counts);

    bit_id += skip;
    byte = unit_bytes[skip / 8];
    i -= accumulated_counts_bytes[skip / 8];
  }

  return bit_id + SELECT_TABLE[i][byte];
}
  #else    // MARISA_USE_SSE2
const uint8_t POPCNT_TABLE[256] = {
    // clang-format off
  0, 1, 1, 2, 1, 2, 2, 3, 1, 2, 2, 3, 2, 3, 3, 4,
  1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5,
  1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5,
  2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
  1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5,
  2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
  2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
  3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7,
  1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5,
  2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
  2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
  3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7,
  2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
  3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7,
  3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7,
  4, 5, 5, 6, 5, 6, 6, 7, 5, 6, 6, 7, 6, 7, 7, 8
    // clang-format on
};

std::size_t select_bit(std::size_t i, std::size_t bit_id, uint32_t unit_lo,
                       uint32_t unit_hi) {
  uint32_t next_byte = unit_lo & 0xFF;
  uint32_t byte_popcount = POPCNT_TABLE[next_byte];
  // Assuming the desired bit is in a random byte, branches are not
  // taken 7/8 of the time, so this is branch-predictor friendly,
  // unlike binary search.
  if (i < byte_popcount) return bit_id + SELECT_TABLE[i][next_byte];
  i -= byte_popcount;
  next_byte = (unit_lo >> 8) & 0xFF;
  byte_popcount = POPCNT_TABLE[next_byte];
  if (i < byte_popcount) return bit_id + 8 + SELECT_TABLE[i][next_byte];
  i -= byte_popcount;
  next_byte = (unit_lo >> 16) & 0xFF;
  byte_popcount = POPCNT_TABLE[next_byte];
  if (i < byte_popcount) return bit_id + 16 + SELECT_TABLE[i][next_byte];
  i -= byte_popcount;
  next_byte = unit_lo >> 24;
  byte_popcount = POPCNT_TABLE[next_byte];
  if (i < byte_popcount) return bit_id + 24 + SELECT_TABLE[i][next_byte];
  i -= byte_popcount;

  next_byte = unit_hi & 0xFF;
  byte_popcount = POPCNT_TABLE[next_byte];
  if (i < byte_popcount) return bit_id + 32 + SELECT_TABLE[i][next_byte];
  i -= byte_popcount;
  next_byte = (unit_hi >> 8) & 0xFF;
  byte_popcount = POPCNT_TABLE[next_byte];
  if (i < byte_popcount) return bit_id + 40 + SELECT_TABLE[i][next_byte];
  i -= byte_popcount;
  next_byte = (unit_hi >> 16) & 0xFF;
  byte_popcount = POPCNT_TABLE[next_byte];
  if (i < byte_popcount) return bit_id + 48 + SELECT_TABLE[i][next_byte];
  i -= byte_popcount;
  next_byte = unit_hi >> 24;
  // Assume `i < POPCNT_TABLE[next_byte]`.
  return bit_id + 56 + SELECT_TABLE[i][next_byte];
}
  #endif   // MARISA_USE_SSE2

// This is only used by build_index, so don't worry about the small performance
// penalty from not having version taking only a uint32_t.
inline std::size_t select_bit(std::size_t i, std::size_t bit_id,
                              uint32_t unit) {
  return select_bit(i, bit_id, /*unit_lo=*/unit, /*unit_hi=*/0);
}

 #endif  // MARISA_WORD_SIZE == 64
#endif   // MARISA_USE_BMI2

}  // namespace

#if MARISA_WORD_SIZE == 64

std::size_t BitVector::rank1(std::size_t i) const {
  assert(!ranks_.empty());
  assert(i <= size_);

  const RankIndex &rank = ranks_[i / 512];
  std::size_t offset = rank.abs();
  switch ((i / 64) % 8) {
    case 1: {
      offset += rank.rel1();
      break;
    }
    case 2: {
      offset += rank.rel2();
      break;
    }
    case 3: {
      offset += rank.rel3();
      break;
    }
    case 4: {
      offset += rank.rel4();
      break;
    }
    case 5: {
      offset += rank.rel5();
      break;
    }
    case 6: {
      offset += rank.rel6();
      break;
    }
    case 7: {
      offset += rank.rel7();
      break;
    }
  }
  offset += popcount(units_[i / 64] & ((1ULL << (i % 64)) - 1));
  return offset;
}

std::size_t BitVector::select0(std::size_t i) const {
  assert(!select0s_.empty());
  assert(i < num_0s());

  const std::size_t select_id = i / 512;
  assert((select_id + 1) < select0s_.size());
  if ((i % 512) == 0) {
    return select0s_[select_id];
  }
  std::size_t begin = select0s_[select_id] / 512;
  std::size_t end = (select0s_[select_id + 1] + 511) / 512;
  if (begin + 10 >= end) {
    while (i >= ((begin + 1) * 512) - ranks_[begin + 1].abs()) {
      ++begin;
    }
  } else {
    while (begin + 1 < end) {
      const std::size_t middle = (begin + end) / 2;
      if (i < (middle * 512) - ranks_[middle].abs()) {
        end = middle;
      } else {
        begin = middle;
      }
    }
  }
  const std::size_t rank_id = begin;
  i -= (rank_id * 512) - ranks_[rank_id].abs();

  const RankIndex &rank = ranks_[rank_id];
  std::size_t unit_id = rank_id * 8;
  if (i < (256U - rank.rel4())) {
    if (i < (128U - rank.rel2())) {
      if (i >= (64U - rank.rel1())) {
        unit_id += 1;
        i -= 64 - rank.rel1();
      }
    } else if (i < (192U - rank.rel3())) {
      unit_id += 2;
      i -= 128 - rank.rel2();
    } else {
      unit_id += 3;
      i -= 192 - rank.rel3();
    }
  } else if (i < (384U - rank.rel6())) {
    if (i < (320U - rank.rel5())) {
      unit_id += 4;
      i -= 256 - rank.rel4();
    } else {
      unit_id += 5;
      i -= 320 - rank.rel5();
    }
  } else if (i < (448U - rank.rel7())) {
    unit_id += 6;
    i -= 384 - rank.rel6();
  } else {
    unit_id += 7;
    i -= 448 - rank.rel7();
  }

  return select_bit(i, unit_id * 64, ~units_[unit_id]);
}

std::size_t BitVector::select1(std::size_t i) const {
  assert(!select1s_.empty());
  assert(i < num_1s());

  const std::size_t select_id = i / 512;
  assert((select_id + 1) < select1s_.size());
  if ((i % 512) == 0) {
    return select1s_[select_id];
  }
  std::size_t begin = select1s_[select_id] / 512;
  std::size_t end = (select1s_[select_id + 1] + 511) / 512;
  if (begin + 10 >= end) {
    while (i >= ranks_[begin + 1].abs()) {
      ++begin;
    }
  } else {
    while (begin + 1 < end) {
      const std::size_t middle = (begin + end) / 2;
      if (i < ranks_[middle].abs()) {
        end = middle;
      } else {
        begin = middle;
      }
    }
  }
  const std::size_t rank_id = begin;
  i -= ranks_[rank_id].abs();

  const RankIndex &rank = ranks_[rank_id];
  std::size_t unit_id = rank_id * 8;
  if (i < rank.rel4()) {
    if (i < rank.rel2()) {
      if (i >= rank.rel1()) {
        unit_id += 1;
        i -= rank.rel1();
      }
    } else if (i < rank.rel3()) {
      unit_id += 2;
      i -= rank.rel2();
    } else {
      unit_id += 3;
      i -= rank.rel3();
    }
  } else if (i < rank.rel6()) {
    if (i < rank.rel5()) {
      unit_id += 4;
      i -= rank.rel4();
    } else {
      unit_id += 5;
      i -= rank.rel5();
    }
  } else if (i < rank.rel7()) {
    unit_id += 6;
    i -= rank.rel6();
  } else {
    unit_id += 7;
    i -= rank.rel7();
  }

  return select_bit(i, unit_id * 64, units_[unit_id]);
}

#else  // MARISA_WORD_SIZE == 64

std::size_t BitVector::rank1(std::size_t i) const {
  assert(!ranks_.empty());
  assert(i <= size_);

  const RankIndex &rank = ranks_[i / 512];
  std::size_t offset = rank.abs();
  switch ((i / 64) % 8) {
    case 1: {
      offset += rank.rel1();
      break;
    }
    case 2: {
      offset += rank.rel2();
      break;
    }
    case 3: {
      offset += rank.rel3();
      break;
    }
    case 4: {
      offset += rank.rel4();
      break;
    }
    case 5: {
      offset += rank.rel5();
      break;
    }
    case 6: {
      offset += rank.rel6();
      break;
    }
    case 7: {
      offset += rank.rel7();
      break;
    }
  }
  if (((i / 32) & 1) == 1) {
    offset += popcount(units_[(i / 32) - 1]);
  }
  offset += popcount(units_[i / 32] & ((1U << (i % 32)) - 1));
  return offset;
}

std::size_t BitVector::select0(std::size_t i) const {
  assert(!select0s_.empty());
  assert(i < num_0s());

  const std::size_t select_id = i / 512;
  assert((select_id + 1) < select0s_.size());
  if ((i % 512) == 0) {
    return select0s_[select_id];
  }
  std::size_t begin = select0s_[select_id] / 512;
  std::size_t end = (select0s_[select_id + 1] + 511) / 512;
  if (begin + 10 >= end) {
    while (i >= ((begin + 1) * 512) - ranks_[begin + 1].abs()) {
      ++begin;
    }
  } else {
    while (begin + 1 < end) {
      const std::size_t middle = (begin + end) / 2;
      if (i < (middle * 512) - ranks_[middle].abs()) {
        end = middle;
      } else {
        begin = middle;
      }
    }
  }
  const std::size_t rank_id = begin;
  i -= (rank_id * 512) - ranks_[rank_id].abs();

  const RankIndex &rank = ranks_[rank_id];
  std::size_t unit_id = rank_id * 16;
  if (i < (256U - rank.rel4())) {
    if (i < (128U - rank.rel2())) {
      if (i >= (64U - rank.rel1())) {
        unit_id += 2;
        i -= 64 - rank.rel1();
      }
    } else if (i < (192U - rank.rel3())) {
      unit_id += 4;
      i -= 128 - rank.rel2();
    } else {
      unit_id += 6;
      i -= 192 - rank.rel3();
    }
  } else if (i < (384U - rank.rel6())) {
    if (i < (320U - rank.rel5())) {
      unit_id += 8;
      i -= 256 - rank.rel4();
    } else {
      unit_id += 10;
      i -= 320 - rank.rel5();
    }
  } else if (i < (448U - rank.rel7())) {
    unit_id += 12;
    i -= 384 - rank.rel6();
  } else {
    unit_id += 14;
    i -= 448 - rank.rel7();
  }

  return select_bit(i, unit_id * 32, ~units_[unit_id], ~units_[unit_id + 1]);
}

std::size_t BitVector::select1(std::size_t i) const {
  assert(!select1s_.empty());
  assert(i < num_1s());

  const std::size_t select_id = i / 512;
  assert((select_id + 1) < select1s_.size());
  if ((i % 512) == 0) {
    return select1s_[select_id];
  }
  std::size_t begin = select1s_[select_id] / 512;
  std::size_t end = (select1s_[select_id + 1] + 511) / 512;
  if (begin + 10 >= end) {
    while (i >= ranks_[begin + 1].abs()) {
      ++begin;
    }
  } else {
    while (begin + 1 < end) {
      const std::size_t middle = (begin + end) / 2;
      if (i < ranks_[middle].abs()) {
        end = middle;
      } else {
        begin = middle;
      }
    }
  }
  const std::size_t rank_id = begin;
  i -= ranks_[rank_id].abs();

  const RankIndex &rank = ranks_[rank_id];
  std::size_t unit_id = rank_id * 16;
  if (i < rank.rel4()) {
    if (i < rank.rel2()) {
      if (i >= rank.rel1()) {
        unit_id += 2;
        i -= rank.rel1();
      }
    } else if (i < rank.rel3()) {
      unit_id += 4;
      i -= rank.rel2();
    } else {
      unit_id += 6;
      i -= rank.rel3();
    }
  } else if (i < rank.rel6()) {
    if (i < rank.rel5()) {
      unit_id += 8;
      i -= rank.rel4();
    } else {
      unit_id += 10;
      i -= rank.rel5();
    }
  } else if (i < rank.rel7()) {
    unit_id += 12;
    i -= rank.rel6();
  } else {
    unit_id += 14;
    i -= rank.rel7();
  }

  return select_bit(i, unit_id * 32, units_[unit_id], units_[unit_id + 1]);
}

#endif  // MARISA_WORD_SIZE == 64

void BitVector::build_index(const BitVector &bv, bool enables_select0,
                            bool enables_select1) {
  const std::size_t num_bits = bv.size();
  ranks_.resize((num_bits / 512) + (((num_bits % 512) != 0) ? 1 : 0) + 1);

  std::size_t num_0s = 0;  // Only updated if enables_select0 is true.
  std::size_t num_1s = 0;

  const std::size_t num_units = bv.units_.size();
  for (std::size_t unit_id = 0; unit_id < num_units; ++unit_id) {
    const std::size_t bit_id = unit_id * MARISA_WORD_SIZE;

    if ((bit_id % 64) == 0) {
      const std::size_t rank_id = bit_id / 512;
      switch ((bit_id / 64) % 8) {
        case 0: {
          ranks_[rank_id].set_abs(num_1s);
          break;
        }
        case 1: {
          ranks_[rank_id].set_rel1(num_1s - ranks_[rank_id].abs());
          break;
        }
        case 2: {
          ranks_[rank_id].set_rel2(num_1s - ranks_[rank_id].abs());
          break;
        }
        case 3: {
          ranks_[rank_id].set_rel3(num_1s - ranks_[rank_id].abs());
          break;
        }
        case 4: {
          ranks_[rank_id].set_rel4(num_1s - ranks_[rank_id].abs());
          break;
        }
        case 5: {
          ranks_[rank_id].set_rel5(num_1s - ranks_[rank_id].abs());
          break;
        }
        case 6: {
          ranks_[rank_id].set_rel6(num_1s - ranks_[rank_id].abs());
          break;
        }
        case 7: {
          ranks_[rank_id].set_rel7(num_1s - ranks_[rank_id].abs());
          break;
        }
      }
    }

    const Unit unit = bv.units_[unit_id];
    // push_back resizes with 0, so the high bits of the last unit are 0 and
    // do not affect the 1s count.
    const std::size_t unit_num_1s = popcount(unit);

    if (enables_select0) {
      // num_0s is somewhat move involved to compute, so only do it if we
      // need it.  The last word has zeros in the high bits, so that needs
      // to be accounted for when computing the unit_num_0s from unit_num_1s.
      const std::size_t bits_remaining = num_bits - bit_id;
      const std::size_t unit_num_0s =
          std::min<std::size_t>(bits_remaining, MARISA_WORD_SIZE) - unit_num_1s;

      // Note: MSVC rejects unary minus operator applied to unsigned type.
      const std::size_t zero_bit_id = (0 - num_0s) % 512;
      if (unit_num_0s > zero_bit_id) {
        // select0s_ is uint32_t, but select_bit returns size_t, so cast to
        // suppress narrowing conversion warning.  push_back checks the
        // size, so there is no truncation here.
        select0s_.push_back(
            static_cast<uint32_t>(select_bit(zero_bit_id, bit_id, ~unit)));
      }

      num_0s += unit_num_0s;
    }

    if (enables_select1) {
      // Note: MSVC rejects unary minus operator applied to unsigned type.
      const std::size_t one_bit_id = (0 - num_1s) % 512;
      if (unit_num_1s > one_bit_id) {
        select1s_.push_back(
            static_cast<uint32_t>(select_bit(one_bit_id, bit_id, unit)));
      }
    }

    num_1s += unit_num_1s;
  }

  if ((num_bits % 512) != 0) {
    const std::size_t rank_id = (num_bits - 1) / 512;
    switch (((num_bits - 1) / 64) % 8) {
      case 0: {
        ranks_[rank_id].set_rel1(num_1s - ranks_[rank_id].abs());
      }
        [[fallthrough]];
      case 1: {
        ranks_[rank_id].set_rel2(num_1s - ranks_[rank_id].abs());
      }
        [[fallthrough]];
      case 2: {
        ranks_[rank_id].set_rel3(num_1s - ranks_[rank_id].abs());
      }
        [[fallthrough]];
      case 3: {
        ranks_[rank_id].set_rel4(num_1s - ranks_[rank_id].abs());
      }
        [[fallthrough]];
      case 4: {
        ranks_[rank_id].set_rel5(num_1s - ranks_[rank_id].abs());
      }
        [[fallthrough]];
      case 5: {
        ranks_[rank_id].set_rel6(num_1s - ranks_[rank_id].abs());
      }
        [[fallthrough]];
      case 6: {
        ranks_[rank_id].set_rel7(num_1s - ranks_[rank_id].abs());
        break;
      }
    }
  }

  size_ = num_bits;
  num_1s_ = bv.num_1s();

  ranks_.back().set_abs(num_1s);
  if (enables_select0) {
    select0s_.push_back(static_cast<uint32_t>(num_bits));
    select0s_.shrink();
  }
  if (enables_select1) {
    select1s_.push_back(static_cast<uint32_t>(num_bits));
    select1s_.shrink();
  }
}

}  // namespace marisa::grimoire::vector
#line 1 "lib/marisa/agent.cc"
#include "./../lib/marisa.h" //amalgamate marisa/agent.h

#include <new>
#include <stdexcept>
#include <utility>

//amalgamate marisa/grimoire/trie.h
//amalgamate marisa/grimoire/trie/state.h
#include "./../lib/marisa.h" //amalgamate marisa/key.h

namespace marisa {
namespace {
void UpdateAgentAfterCopyingState(const grimoire::trie::State &state,
                                  Agent &agent) {
  // Point the agent's key to the newly copied buffer if necessary.
  switch (state.status_code()) {
    case grimoire::trie::MARISA_READY_TO_PREDICTIVE_SEARCH:
    case grimoire::trie::MARISA_END_OF_PREDICTIVE_SEARCH:
      // In states corresponding to predictive_search, the agent's
      // key points into the state key buffer. We need to repoint
      // after copying the state.
      agent.set_key(state.key_buf().data(), state.key_buf().size());
      break;
    default:
      // In other states, they key is either null, or points to the
      // query, so we do not need to repoint it.
      break;
  }
}
}  // namespace

Agent::Agent() = default;

Agent::~Agent() = default;

Agent::Agent(const Agent &other)
    : query_(other.query_), key_(other.key_),
      state_(other.has_state() ? new grimoire::trie::State(other.state())
                               : nullptr) {
  if (other.has_state()) {
    UpdateAgentAfterCopyingState(*state_, *this);
  }
}

Agent &Agent::operator=(const Agent &other) {
  query_ = other.query_;
  key_ = other.key_;
  if (other.has_state()) {
    state_.reset(new grimoire::trie::State(other.state()));
    UpdateAgentAfterCopyingState(*state_, *this);
  } else {
    state_ = nullptr;
  }
  return *this;
}

Agent::Agent(Agent &&other) noexcept = default;
Agent &Agent::operator=(Agent &&other) noexcept = default;

void Agent::set_query(const char *str) {
  MARISA_THROW_IF(str == nullptr, std::invalid_argument);
  if (state_ != nullptr) {
    state_->reset();
  }
  query_.set_str(str);
}

void Agent::set_query(const char *ptr, std::size_t length) {
  MARISA_THROW_IF((ptr == nullptr) && (length != 0), std::invalid_argument);
  if (state_ != nullptr) {
    state_->reset();
  }
  query_.set_str(ptr, length);
}

void Agent::set_query(std::size_t key_id) {
  if (state_ != nullptr) {
    state_->reset();
  }
  query_.set_id(key_id);
}

void Agent::init_state() {
  MARISA_THROW_IF(state_ != nullptr, std::logic_error);
  state_.reset(new grimoire::State);
}

void Agent::clear() noexcept {
  Agent().swap(*this);
}

void Agent::swap(Agent &rhs) noexcept {
  query_.swap(rhs.query_);
  key_.swap(rhs.key_);
  state_.swap(rhs.state_);
}

}  // namespace marisa
#line 1 "lib/marisa/keyset.cc"
#include "./../lib/marisa.h" //amalgamate marisa/keyset.h

#include <cassert>
#include <cstring>
#include <memory>
#include <new>
#include <stdexcept>

namespace marisa {

Keyset::Keyset() = default;

void Keyset::push_back(const Key &key) {
  assert(size_ < SIZE_MAX);

  char *const key_ptr = reserve(key.length());
  std::memcpy(key_ptr, key.ptr(), key.length());

  Key &new_key = key_blocks_[size_ / KEY_BLOCK_SIZE][size_ % KEY_BLOCK_SIZE];
  new_key.set_str(key_ptr, key.length());
  new_key.set_id(key.id());
  ++size_;
  total_length_ += new_key.length();
}

void Keyset::push_back(const Key &key, char end_marker) {
  assert(size_ < SIZE_MAX);

  if ((size_ / KEY_BLOCK_SIZE) == key_blocks_size_) {
    append_key_block();
  }

  char *const key_ptr = reserve(key.length() + 1);
  std::memcpy(key_ptr, key.ptr(), key.length());
  key_ptr[key.length()] = end_marker;

  Key &new_key = key_blocks_[size_ / KEY_BLOCK_SIZE][size_ % KEY_BLOCK_SIZE];
  new_key.set_str(key_ptr, key.length());
  new_key.set_id(key.id());
  ++size_;
  total_length_ += new_key.length();
}

void Keyset::push_back(const char *str) {
  assert(size_ < SIZE_MAX);
  MARISA_THROW_IF(str == nullptr, std::invalid_argument);

  std::size_t length = 0;
  while (str[length] != '\0') {
    ++length;
  }
  push_back(str, length);
}

void Keyset::push_back(const char *ptr, std::size_t length, float weight) {
  assert(size_ < SIZE_MAX);
  MARISA_THROW_IF((ptr == nullptr) && (length != 0), std::invalid_argument);
  MARISA_THROW_IF(length > UINT32_MAX, std::invalid_argument);

  char *const key_ptr = reserve(length);
  std::memcpy(key_ptr, ptr, length);

  Key &key = key_blocks_[size_ / KEY_BLOCK_SIZE][size_ % KEY_BLOCK_SIZE];
  key.set_str(key_ptr, length);
  key.set_weight(weight);
  ++size_;
  total_length_ += length;
}

void Keyset::reset() {
  base_blocks_size_ = 0;
  extra_blocks_size_ = 0;
  ptr_ = nullptr;
  avail_ = 0;
  size_ = 0;
  total_length_ = 0;
}

void Keyset::clear() noexcept {
  Keyset().swap(*this);
}

void Keyset::swap(Keyset &rhs) noexcept {
  base_blocks_.swap(rhs.base_blocks_);
  std::swap(base_blocks_size_, rhs.base_blocks_size_);
  std::swap(base_blocks_capacity_, rhs.base_blocks_capacity_);
  extra_blocks_.swap(rhs.extra_blocks_);
  std::swap(extra_blocks_size_, rhs.extra_blocks_size_);
  std::swap(extra_blocks_capacity_, rhs.extra_blocks_capacity_);
  key_blocks_.swap(rhs.key_blocks_);
  std::swap(key_blocks_size_, rhs.key_blocks_size_);
  std::swap(key_blocks_capacity_, rhs.key_blocks_capacity_);
  std::swap(ptr_, rhs.ptr_);
  std::swap(avail_, rhs.avail_);
  std::swap(size_, rhs.size_);
  std::swap(total_length_, rhs.total_length_);
}

char *Keyset::reserve(std::size_t size) {
  if ((size_ / KEY_BLOCK_SIZE) == key_blocks_size_) {
    append_key_block();
  }

  if (size > EXTRA_BLOCK_SIZE) {
    append_extra_block(size);
    return extra_blocks_[extra_blocks_size_ - 1].get();
  }
  if (size > avail_) {
    append_base_block();
  }
  ptr_ += size;
  avail_ -= size;
  return ptr_ - size;
}

void Keyset::append_base_block() {
  if (base_blocks_size_ == base_blocks_capacity_) {
    const std::size_t new_capacity =
        (base_blocks_size_ != 0) ? (base_blocks_size_ * 2) : 1;
    std::unique_ptr<std::unique_ptr<char[]>[]> new_blocks(
        new std::unique_ptr<char[]>[new_capacity]);
    for (std::size_t i = 0; i < base_blocks_size_; ++i) {
      base_blocks_[i].swap(new_blocks[i]);
    }
    base_blocks_.swap(new_blocks);
    base_blocks_capacity_ = new_capacity;
  }
  if (base_blocks_[base_blocks_size_] == nullptr) {
    std::unique_ptr<char[]> new_block(new char[BASE_BLOCK_SIZE]);
    base_blocks_[base_blocks_size_].swap(new_block);
  }
  ptr_ = base_blocks_[base_blocks_size_++].get();
  avail_ = BASE_BLOCK_SIZE;
}

void Keyset::append_extra_block(std::size_t size) {
  if (extra_blocks_size_ == extra_blocks_capacity_) {
    const std::size_t new_capacity =
        (extra_blocks_size_ != 0) ? (extra_blocks_size_ * 2) : 1;
    std::unique_ptr<std::unique_ptr<char[]>[]> new_blocks(
        new std::unique_ptr<char[]>[new_capacity]);
    for (std::size_t i = 0; i < extra_blocks_size_; ++i) {
      extra_blocks_[i].swap(new_blocks[i]);
    }
    extra_blocks_.swap(new_blocks);
    extra_blocks_capacity_ = new_capacity;
  }
  std::unique_ptr<char[]> new_block(new char[size]);
  extra_blocks_[extra_blocks_size_++].swap(new_block);
}

void Keyset::append_key_block() {
  if (key_blocks_size_ == key_blocks_capacity_) {
    const std::size_t new_capacity =
        (key_blocks_size_ != 0) ? (key_blocks_size_ * 2) : 1;
    std::unique_ptr<std::unique_ptr<Key[]>[]> new_blocks(
        new std::unique_ptr<Key[]>[new_capacity]);
    for (std::size_t i = 0; i < key_blocks_size_; ++i) {
      key_blocks_[i].swap(new_blocks[i]);
    }
    key_blocks_.swap(new_blocks);
    key_blocks_capacity_ = new_capacity;
  }
  std::unique_ptr<Key[]> new_block(new Key[KEY_BLOCK_SIZE]);
  key_blocks_[key_blocks_size_++].swap(new_block);
}

}  // namespace marisa
#line 1 "lib/marisa/trie.cc"
#include "./../lib/marisa.h" //amalgamate marisa/trie.h

#include <memory>
#include <stdexcept>

//amalgamate marisa/grimoire/trie.h
//amalgamate marisa/iostream.h
//amalgamate marisa/stdio.h

namespace marisa {

Trie::Trie() = default;

Trie::~Trie() = default;

Trie::Trie(Trie &&other) noexcept = default;

Trie &Trie::operator=(Trie &&other) noexcept = default;

void Trie::build(Keyset &keyset, int config_flags) {
  std::unique_ptr<grimoire::LoudsTrie> temp(new grimoire::LoudsTrie);

  temp->build(keyset, config_flags);
  trie_.swap(temp);
}

void Trie::mmap(const char *filename, int flags) {
  MARISA_THROW_IF(filename == nullptr, std::invalid_argument);

  std::unique_ptr<grimoire::LoudsTrie> temp(new grimoire::LoudsTrie);

  grimoire::Mapper mapper;
  mapper.open(filename, flags);
  temp->map(mapper);
  trie_.swap(temp);
}

void Trie::map(const void *ptr, std::size_t size) {
  MARISA_THROW_IF((ptr == nullptr) && (size != 0), std::invalid_argument);

  std::unique_ptr<grimoire::LoudsTrie> temp(new grimoire::LoudsTrie);

  grimoire::Mapper mapper;
  mapper.open(ptr, size);
  temp->map(mapper);
  trie_.swap(temp);
}

void Trie::load(const char *filename) {
  MARISA_THROW_IF(filename == nullptr, std::invalid_argument);

  std::unique_ptr<grimoire::LoudsTrie> temp(new grimoire::LoudsTrie);

  grimoire::Reader reader;
  reader.open(filename);
  temp->read(reader);
  trie_.swap(temp);
}

void Trie::read(int fd) {
  MARISA_THROW_IF(fd == -1, std::invalid_argument);

  std::unique_ptr<grimoire::LoudsTrie> temp(new grimoire::LoudsTrie);

  grimoire::Reader reader;
  reader.open(fd);
  temp->read(reader);
  trie_.swap(temp);
}

void Trie::save(const char *filename) const {
  MARISA_THROW_IF(trie_ == nullptr, std::logic_error);
  MARISA_THROW_IF(filename == nullptr, std::invalid_argument);

  grimoire::Writer writer;
  writer.open(filename);
  trie_->write(writer);
}

void Trie::write(int fd) const {
  MARISA_THROW_IF(trie_ == nullptr, std::logic_error);
  MARISA_THROW_IF(fd == -1, std::invalid_argument);

  grimoire::Writer writer;
  writer.open(fd);
  trie_->write(writer);
}

bool Trie::lookup(Agent &agent) const {
  MARISA_THROW_IF(trie_ == nullptr, std::logic_error);
  if (!agent.has_state()) {
    agent.init_state();
  }
  return trie_->lookup(agent);
}

void Trie::reverse_lookup(Agent &agent) const {
  MARISA_THROW_IF(trie_ == nullptr, std::logic_error);
  if (!agent.has_state()) {
    agent.init_state();
  }
  trie_->reverse_lookup(agent);
}

bool Trie::common_prefix_search(Agent &agent) const {
  MARISA_THROW_IF(trie_ == nullptr, std::logic_error);
  if (!agent.has_state()) {
    agent.init_state();
  }
  return trie_->common_prefix_search(agent);
}

bool Trie::predictive_search(Agent &agent) const {
  MARISA_THROW_IF(trie_ == nullptr, std::logic_error);
  if (!agent.has_state()) {
    agent.init_state();
  }
  return trie_->predictive_search(agent);
}

std::size_t Trie::num_tries() const {
  MARISA_THROW_IF(trie_ == nullptr, std::logic_error);
  return trie_->num_tries();
}

std::size_t Trie::num_keys() const {
  MARISA_THROW_IF(trie_ == nullptr, std::logic_error);
  return trie_->num_keys();
}

std::size_t Trie::num_nodes() const {
  MARISA_THROW_IF(trie_ == nullptr, std::logic_error);
  return trie_->num_nodes();
}

TailMode Trie::tail_mode() const {
  MARISA_THROW_IF(trie_ == nullptr, std::logic_error);
  return trie_->tail_mode();
}

NodeOrder Trie::node_order() const {
  MARISA_THROW_IF(trie_ == nullptr, std::logic_error);
  return trie_->node_order();
}

bool Trie::empty() const {
  MARISA_THROW_IF(trie_ == nullptr, std::logic_error);
  return trie_->empty();
}

std::size_t Trie::size() const {
  MARISA_THROW_IF(trie_ == nullptr, std::logic_error);
  return trie_->size();
}

std::size_t Trie::total_size() const {
  MARISA_THROW_IF(trie_ == nullptr, std::logic_error);
  return trie_->total_size();
}

std::size_t Trie::io_size() const {
  MARISA_THROW_IF(trie_ == nullptr, std::logic_error);
  return trie_->io_size();
}

void Trie::clear() noexcept {
  Trie().swap(*this);
}

void Trie::swap(Trie &rhs) noexcept {
  trie_.swap(rhs.trie_);
}

}  // namespace marisa

#include <iostream>

namespace marisa {

class TrieIO {
 public:
  static void fread(std::FILE *file, Trie *trie) {
    MARISA_THROW_IF(trie == nullptr, std::invalid_argument);

    std::unique_ptr<grimoire::LoudsTrie> temp(new grimoire::LoudsTrie);

    grimoire::Reader reader;
    reader.open(file);
    temp->read(reader);
    trie->trie_.swap(temp);
  }
  static void fwrite(std::FILE *file, const Trie &trie) {
    MARISA_THROW_IF(file == nullptr, std::invalid_argument);
    MARISA_THROW_IF(trie.trie_ == nullptr, std::logic_error);
    grimoire::Writer writer;
    writer.open(file);
    trie.trie_->write(writer);
  }

  static std::istream &read(std::istream &stream, Trie *trie) {
    MARISA_THROW_IF(trie == nullptr, std::invalid_argument);

    std::unique_ptr<grimoire::LoudsTrie> temp(new grimoire::LoudsTrie);

    grimoire::Reader reader;
    reader.open(stream);
    temp->read(reader);
    trie->trie_.swap(temp);
    return stream;
  }
  static std::ostream &write(std::ostream &stream, const Trie &trie) {
    MARISA_THROW_IF(trie.trie_ == nullptr, std::logic_error);
    grimoire::Writer writer;
    writer.open(stream);
    trie.trie_->write(writer);
    return stream;
  }
};

void fread(std::FILE *file, Trie *trie) {
  MARISA_THROW_IF(file == nullptr, std::invalid_argument);
  MARISA_THROW_IF(trie == nullptr, std::invalid_argument);
  TrieIO::fread(file, trie);
}

void fwrite(std::FILE *file, const Trie &trie) {
  MARISA_THROW_IF(file == nullptr, std::invalid_argument);
  TrieIO::fwrite(file, trie);
}

std::istream &read(std::istream &stream, Trie *trie) {
  MARISA_THROW_IF(trie == nullptr, std::invalid_argument);
  return TrieIO::read(stream, trie);
}

std::ostream &write(std::ostream &stream, const Trie &trie) {
  return TrieIO::write(stream, trie);
}

std::istream &operator>>(std::istream &stream, Trie &trie) {
  return read(stream, &trie);
}

std::ostream &operator<<(std::ostream &stream, const Trie &trie) {
  return write(stream, trie);
}

}  // namespace marisa
