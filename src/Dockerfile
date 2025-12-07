FROM debian:trixie AS base

ARG WASI_SDK_VERSION=29
ARG BINARYEN_VERSION=125

ARG WASI_SDK_CHECKSUM=sha256:87d1d1a2879d139cdc624b968efad3d4a97b8078cdff95e63ac88ecafd1a0171
ARG BINARYEN_CHECKSUM=sha256:7c3bc16599c8274a04d34a504fe4be2047884f900e0e2da2f6fb9cd667183be4

ARG WASI_SDK_URL=https://github.com/WebAssembly/wasi-sdk/releases/download/wasi-sdk-${WASI_SDK_VERSION}/wasi-sdk-${WASI_SDK_VERSION}.0-x86_64-linux.tar.gz
ARG BINARYEN_URL=https://github.com/WebAssembly/binaryen/releases/download/version_${BINARYEN_VERSION}/binaryen-version_${BINARYEN_VERSION}-x86_64-linux.tar.gz

FROM base AS wasi-sdk
ADD --checksum=${WASI_SDK_CHECKSUM} ${WASI_SDK_URL} ./wasi-sdk.tar.gz
RUN mkdir /opt/wasi-sdk && tar -C /opt/wasi-sdk --strip-components=1 -xf /wasi-sdk.tar.gz

FROM base AS binaryen
ADD --checksum=${BINARYEN_CHECKSUM} ${BINARYEN_URL} ./binaryen.tar.gz
RUN mkdir /opt/binaryen && tar -C /opt/binaryen --strip-components=1 -xf binaryen.tar.gz

FROM base AS sdk
COPY --link --from=wasi-sdk /opt /opt
COPY --link --from=binaryen /opt /opt
ENV PATH="/opt/wasi-sdk/bin:/opt/binaryen/bin:$PATH"

FROM golang:1.25 AS download
WORKDIR /build/src
COPY --link download.go /build/src/
RUN mkdir /build/lib
RUN go run download.go

FROM sdk AS build
WORKDIR /build/

COPY --from=download /build/lib/ /build/lib/
COPY --link io.h intrin.h /build/src/

RUN wasm32-wasi-clang++ -c -O3 \
 -mmutable-globals -mbulk-memory -mmultivalue -mnontrapping-fptoint -msign-ext -msimd128 \
 -std=c++20 -DNDEBUG \
 lib/marisa.cc

COPY --link io.cc wexcept.cc wrapper.cc /build/src/

RUN wasm32-wasi-clang++ -c -O3 \
 -mmutable-globals -mbulk-memory -mmultivalue -mnontrapping-fptoint -msign-ext -msimd128 -Xclang -target-abi -Xclang experimental-mv \
 -std=c++20 -DNDEBUG -Wall -Wextra -Wconversion \
 src/io.cc src/wexcept.cc src/wrapper.cc

COPY --link exports.txt /build/src/

RUN wasm32-wasi-clang++ -o marisa.wasm -mexec-model=reactor -Wl,--global-base=1024 $(printf " -Wl,--export=%s" $(cat src/exports.txt)) *.o
RUN wasm-ctor-eval -g -c _initialize marisa.wasm -o marisa.tmp
RUN wasm-opt -g marisa.tmp -o marisa.wasm -O3 \
 --enable-mutable-globals --enable-bulk-memory --enable-multivalue --enable-nontrapping-float-to-int --enable-sign-ext --enable-simd \
 --low-memory-unused --strip-debug --strip-dwarf --gufa --converge
RUN install -m644 marisa.wasm lib/marisa.wasm

RUN printf '%s\n' '-Wall' '-Wextra' '-Wconversion' '-isysteminclude' '-xc++' '-std=c++20' > src/compile_flags.txt

FROM scratch
COPY --link --from=build /build/src/compile_flags.txt /src/
COPY --link --from=build /build/lib/ /lib/
