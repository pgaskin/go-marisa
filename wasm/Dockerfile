FROM debian:trixie AS base

# v0.3.1
ARG MARISA_REV=3e87d53b78e15f2f43783d5e376561a8c9722051
ARG MARISA_URL=https://github.com/s-yata/marisa-trie.git

ARG WASI_SDK_VERSION=27
ARG BINARYEN_VERSION=124

ARG WASI_SDK_CHECKSUM=sha256:b7d4d944c88503e4f21d84af07ac293e3440b1b6210bfd7fe78e0afd92c23bc2
ARG BINARYEN_CHECKSUM=sha256:0290c3779fedf592b8da0ded3032ff55c41a2b7bfa2d6bf7b7bac6f0e6e28963

ARG WASI_SDK_URL=https://github.com/WebAssembly/wasi-sdk/releases/download/wasi-sdk-${WASI_SDK_VERSION}/wasi-sdk-${WASI_SDK_VERSION}.0-x86_64-linux.tar.gz
ARG BINARYEN_URL=https://github.com/WebAssembly/binaryen/releases/download/version_${BINARYEN_VERSION}/binaryen-version_${BINARYEN_VERSION}-x86_64-linux.tar.gz

FROM base AS wasi-sdk
ADD --checksum=${WASI_SDK_CHECKSUM} ${WASI_SDK_URL} ./wasi-sdk.tar.gz
RUN mkdir /opt/wasi-sdk && tar -C /opt/wasi-sdk --strip-components=1 -xf /wasi-sdk.tar.gz

FROM base AS binaryen
ADD --checksum=${BINARYEN_CHECKSUM} ${BINARYEN_URL} ./binaryen.tar.gz
RUN mkdir /opt/binaryen && tar -C /opt/binaryen --strip-components=1 -xf binaryen.tar.gz

FROM base AS build
ADD ${MARISA_URL}#${MARISA_REV} /src
WORKDIR /src

COPY --link --from=wasi-sdk /opt /opt
COPY --link --from=binaryen /opt /opt
ENV PATH="/opt/wasi-sdk/bin:/opt/binaryen/bin:$PATH"

# we don't support stack unwinding, so ensure there are no catches
RUN ! grep -e MARISA_CATCH -e catch -Fr lib
RUN ! grep -e MARISA_CATCH -e catch -Fr include | grep -v '//'

# replace throws
# note: commented out for now since implmenting __cxa_throw works fine, but if it ever doesn't, this will work (and add -fno-exceptions to the marisa build)
#RUN sed -e '/#define MARISA_BASE_H_/a void wautil_throw[[noreturn]](const std::exception&); static bool wautil_throw(const std::exception& ex, bool b) { wautil_throw(ex); return b; }' -e 's/(throw /wautil_throw(/g' -i include/marisa/base.h

COPY --link ./src/wautil.cc ./src/wautil.h /src/misc/

RUN wasm32-wasi-clang++ -c -Wall -Wextra -Wconversion -O2 \
 -mmutable-globals -mbulk-memory -mmultivalue -mnontrapping-fptoint -mreference-types -msign-ext -msimd128 -fno-exceptions -Xclang -target-abi -Xclang experimental-mv \
 -std=c++20 -Imisc -isysteminclude \
 misc/wautil.cc

COPY --link ./src/io.cc /src/lib/marisa/grimoire/io/
COPY --link ./src/io.h ./src/intrin.h /src/lib/marisa/grimoire/

RUN wasm32-wasi-clang++ -c -Wall -Wextra -Wconversion -Wno-implicit-fallthrough -Wno-unused-function -O2 \
 -mmutable-globals -mbulk-memory -mmultivalue -mnontrapping-fptoint -mreference-types -msign-ext -msimd128 \
 -std=c++20 -Iinclude -Ilib -Ilib/marisa/grimoire -isystemmisc -DNDEBUG \
 lib/marisa/agent.cc \
 lib/marisa/grimoire/io/io.cc \
 lib/marisa/grimoire/trie/louds-trie.cc \
 lib/marisa/grimoire/trie/tail.cc \
 lib/marisa/grimoire/vector/bit-vector.cc \
 lib/marisa/keyset.cc \
 lib/marisa/trie.cc

COPY --link ./src/wrapper.cc /src/misc/

RUN wasm32-wasi-clang++ -c -Wall -Wextra -Wconversion -O2 \
 -mmutable-globals -mbulk-memory -mmultivalue -mnontrapping-fptoint -mreference-types -msign-ext -msimd128 -fno-exceptions -Xclang -target-abi -Xclang experimental-mv \
 -std=c++20 -isystemmisc -isysteminclude \
 misc/wrapper.cc

COPY --link ./exports.txt /src/

RUN wasm32-wasi-clang++ -o marisa.wasm -mexec-model=reactor $(printf " -Wl,--export=%s" $(cat exports.txt)) *.o

RUN wasm-ctor-eval -g -c _initialize marisa.wasm -o marisa.tmp

RUN wasm-opt -g marisa.tmp -o marisa.wasm -O2 \
 --enable-mutable-globals --enable-bulk-memory --enable-multivalue --enable-nontrapping-float-to-int --enable-reference-types --enable-sign-ext --enable-simd \
 --low-memory-unused --strip-debug --strip-dwarf

RUN chmod 644 marisa.wasm

RUN printf '%s\n' '-Wall' '-Wextra' '-Wconversion' '-isysteminclude' '-xc++' '-std=c++20' > compile_flags.txt

FROM scratch
COPY --link --from=build /src/include /include
COPY --link --from=build /src/compile_flags.txt /
COPY --link --from=build /src/COPYING.md /marisa.COPYING
COPY --link --from=build /src/marisa.wasm /
