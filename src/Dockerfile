FROM debian:trixie AS base

ARG WASI_SDK_VERSION=30
ARG BINARYEN_VERSION=127

ARG WASI_SDK_CHECKSUM=sha256:0507679dff16814b74516cd969a9b16d2ced1347388024bc7966264648c78bfb
ARG BINARYEN_CHECKSUM=sha256:c8ebe5d00c978601086bcad33b2c80fcfe33d6a8b87b754ba4ea86a9a16cc316

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
RUN go run download.go

FROM golang:1.25 AS wasm2go
RUN GOBIN=/ go install -trimpath -v github.com/ncruces/wasm2go@278a5ce072bdf2bb1e90264cbc6132f281827f75

FROM sdk AS build
WORKDIR /build/

COPY --from=download /build/src/marisa.cc /build/src/marisa.h /build/src/
COPY --link io.h /build/src/

RUN wasm32-wasi-clang++ -c -O3 \
 -mmutable-globals -mbulk-memory -mmultivalue -mnontrapping-fptoint -msign-ext \
 -std=c++20 -DNDEBUG \
 src/marisa.cc

COPY --link io.cc wexcept.cc wrapper.cc /build/src/

RUN wasm32-wasi-clang++ -c -O3 \
 -mmutable-globals -mbulk-memory -mmultivalue -mnontrapping-fptoint -msign-ext -Xclang -target-abi -Xclang experimental-mv \
 -std=c++20 -DNDEBUG -Wall -Wextra -Wconversion \
 src/io.cc src/wexcept.cc src/wrapper.cc

COPY --link exports.txt /build/src/

RUN wasm32-wasi-clang++ -o marisa.wasm -mexec-model=reactor -Wl,--global-base=1024 $(printf " -Wl,--export=%s" $(cat src/exports.txt)) *.o
RUN wasm-ctor-eval -g -c _initialize marisa.wasm -o marisa.tmp
RUN wasm-opt -g marisa.tmp -o marisa.wasm -O3 \
 --enable-mutable-globals --enable-bulk-memory --enable-multivalue --enable-nontrapping-float-to-int --enable-sign-ext \
 --low-memory-unused --strip-debug --strip-dwarf --gufa --converge

RUN printf '%s\n' '-Wall' '-Wextra' '-Wconversion' '-isysteminclude' '-xc++' '-std=c++20' > src/compile_flags.txt

COPY --link --from=wasm2go /wasm2go /usr/local/bin/
RUN mkdir -p internal/marisa_wasm
RUN wasm2go -endian=little < marisa.wasm > internal/marisa_wasm/marisa.go
RUN wasm2go -endian=big < marisa.wasm > internal/marisa_wasm/marisa_big.go

FROM scratch
COPY --link --from=build /build/src/compile_flags.txt /src/
COPY --link --from=build /build/internal/ /internal/
COPY --link --from=build /build/src/marisa.cc /build/src/marisa.h /src/
