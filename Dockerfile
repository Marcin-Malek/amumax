# Build frontend with bun
FROM docker.io/oven/bun:1.1.27 AS frontend
WORKDIR /frontend
COPY frontend .
RUN bun install && bun run build

# CUDA 12.4 for H100 (sm_90) support — cuRAND requires matching architecture kernels
FROM docker.io/nvidia/cuda:12.4.0-devel-ubuntu22.04 AS builder
RUN apt-get update \
&& apt-get install -y --no-install-recommends wget && rm -rf /var/lib/apt/lists/*

# Installing go
ENV GO_VERSION=1.23.6
RUN wget https://go.dev/dl/go$GO_VERSION.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go$GO_VERSION.linux-amd64.tar.gz
RUN rm go$GO_VERSION.linux-amd64.tar.gz
ENV PATH=/usr/local/go/bin:$PATH

# Copy source code and frontend build artifacts for building
COPY . .
COPY --from=frontend /frontend/dist /src/api/static

ENV GOPATH=/.go/path
ENV GOCACHE=/.go/cache
ENV CGO_CFLAGS="-I/usr/local/cuda/include/"  
ENV CGO_LDFLAGS="-lcufft -lcuda -lcurand -L/usr/local/cuda/lib64/stubs/ -Wl,-rpath -Wl,\$ORIGIN" 
ENV CGO_CFLAGS_ALLOW='(-fno-schedule-insns|-malign-double|-ffast-math)'

# Build amumax
RUN go build -v -o /build/amumax .

# Move the built binary to a fresh final runtime image
FROM docker.io/nvidia/cuda:12.4.0-runtime-ubuntu22.04
COPY --from=builder /build/amumax /usr/local/bin/amumax

ENTRYPOINT ["/usr/local/bin/amumax"]