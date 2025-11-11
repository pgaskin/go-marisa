#pragma once

#ifdef __wasm__
#define MARISA_USE_SSE2
#define MARISA_USE_SSSE3

#if !__has_builtin(__builtin_ctzll)
#error expected support for __builtin_ctzll
#endif

#if !__has_builtin(__builtin_popcount)
#error expected support for __builtin_popcount
#endif

/*
 * Copyright 2020 The Emscripten Authors.  All rights reserved.
 * Emscripten is available under two separate licenses, the MIT license and the
 * University of Illinois/NCSA Open Source License.
 */

// https://emscripten.org/docs/porting/simd.html
// https://github.com/emscripten-core/emscripten/blob/4.0.18/system/include/compat/xmmintrin.h
// https://github.com/emscripten-core/emscripten/blob/4.0.18/system/include/compat/emmintrin.h
// https://github.com/emscripten-core/emscripten/blob/4.0.18/system/include/compat/tmmintrin.h

#include <wasm_simd128.h>
typedef v128_t __m128i;
[[clang::always_inline]] [[gnu::nodebug]] [[maybe_unused]] static inline __m128i _mm_add_epi8(__m128i __a, __m128i __b) { return (__m128i)wasm_i8x16_add((v128_t)__a, (v128_t)__b); }
[[clang::always_inline]] [[gnu::nodebug]] [[maybe_unused]] static inline __m128i _mm_add_epi32(__m128i __a, __m128i __b) { return (__m128i)wasm_i32x4_add((v128_t)__a, (v128_t)__b); }
[[clang::always_inline]] [[gnu::nodebug]] [[maybe_unused]] static inline __m128i _mm_and_si128(__m128i __a, __m128i __b) { return (__m128i)wasm_v128_and((v128_t)__a, (v128_t)__b); }
[[clang::always_inline]] [[gnu::nodebug]] [[maybe_unused]] static inline __m128i _mm_cmpgt_epi8(__m128i __a, __m128i __b) { return (__m128i)wasm_i8x16_gt((v128_t)__a, (v128_t)__b); }
[[clang::always_inline]] [[gnu::nodebug]] [[maybe_unused]] static inline __m128i _mm_cvtsi32_si128(int __a) { return (__m128i)wasm_i32x4_make(__a, 0, 0, 0); }
[[clang::always_inline]] [[gnu::nodebug]] [[maybe_unused]] static inline __m128i _mm_or_si128(__m128i __a, __m128i __b) { return (__m128i)wasm_v128_or((v128_t)__b, (v128_t)__a); }
[[clang::always_inline]] [[gnu::nodebug]] [[maybe_unused]] static inline __m128i _mm_set_epi8(char b15, char b14, char b13, char b12, char b11, char b10, char b9, char b8, char b7, char b6, char b5, char b4, char b3, char b2, char b1, char b0) { return (__m128i)wasm_i8x16_make(b0, b1, b2, b3, b4, b5, b6, b7, b8, b9, b10, b11, b12, b13, b14, b15); }
[[clang::always_inline]] [[gnu::nodebug]] [[maybe_unused]] static inline __m128i _mm_set_epi32(int i3, int i2, int i1, int i0) { return (__m128i)wasm_i32x4_make(i0, i1, i2, i3); }
[[clang::always_inline]] [[gnu::nodebug]] [[maybe_unused]] static inline __m128i _mm_set1_epi8(char __b) { return (__m128i)wasm_i8x16_splat(__b); }
[[clang::always_inline]] [[gnu::nodebug]] [[maybe_unused]] static inline __m128i _mm_setzero_si128(void) { return wasm_i64x2_const(0, 0); }
[[clang::always_inline]] [[gnu::nodebug]] [[maybe_unused]] static inline __m128i _mm_srli_epi32(__m128i __a, int __count) { return (__m128i)(((unsigned int)__count < 32) ? wasm_u32x4_shr((v128_t)__a, (unsigned int)__count) : wasm_i32x4_const(0,0,0,0)); }
[[clang::always_inline]] [[gnu::nodebug]] [[maybe_unused]] static inline __m128i _mm_sub_epi8(__m128i __a, __m128i __b) { return (__m128i)wasm_i8x16_sub((v128_t)__a, (v128_t)__b); }
[[clang::always_inline]] [[gnu::nodebug]] [[maybe_unused]] static inline __m128i _mm_shuffle_epi8(__m128i __a, __m128i __b) { return (__m128i)wasm_i8x16_swizzle((v128_t)__a, (v128_t)_mm_and_si128(__b, _mm_set1_epi8(0x8F))); }
[[clang::always_inline]] [[gnu::nodebug]] [[maybe_unused]] static inline int     _mm_movemask_epi8(__m128i __a) { return (int)wasm_i8x16_bitmask((v128_t)__a); }
[[clang::always_inline]] [[gnu::nodebug]] [[maybe_unused]] static inline void    _mm_store_si128(__m128i *__p, __m128i __b) { *__p = __b; }
#define _mm_slli_si128(__a, __imm) __extension__ ({ (__m128i)wasm_i8x16_shuffle(_mm_setzero_si128(), (__a), ((__imm)&0xF0) ? 0 : 16 - ((__imm)&0xF), ((__imm)&0xF0) ? 0 : 17 - ((__imm)&0xF), ((__imm)&0xF0) ? 0 : 18 - ((__imm)&0xF), ((__imm)&0xF0) ? 0 : 19 - ((__imm)&0xF), ((__imm)&0xF0) ? 0 : 20 - ((__imm)&0xF), ((__imm)&0xF0) ? 0 : 21 - ((__imm)&0xF), ((__imm)&0xF0) ? 0 : 22 - ((__imm)&0xF), ((__imm)&0xF0) ? 0 : 23 - ((__imm)&0xF), ((__imm)&0xF0) ? 0 : 24 - ((__imm)&0xF), ((__imm)&0xF0) ? 0 : 25 - ((__imm)&0xF), ((__imm)&0xF0) ? 0 : 26 - ((__imm)&0xF), ((__imm)&0xF0) ? 0 : 27 - ((__imm)&0xF), ((__imm)&0xF0) ? 0 : 28 - ((__imm)&0xF), ((__imm)&0xF0) ? 0 : 29 - ((__imm)&0xF), ((__imm)&0xF0) ? 0 : 30 - ((__imm)&0xF), ((__imm)&0xF0) ? 0 : 31 - ((__imm)&0xF)); })
#endif
