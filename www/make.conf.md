```conf
# /etc/make.conf hsbdsrv by @Dr.Deep

# CCACHE
WITH_CCACHE_BUILD=yes
CCACHE_DIR=/usr/ccache

# SIMD
CFLAGS_COMMON            = -O3 -pipe -march=x86-64 -falign-functions=32 -fno-semantic-interposition -g1 -gz=zstd -gno-column-info

CFLAGS                  += ${CFLAGS_COMMON}
CXXFLAGS                += ${CFLAGS_COMMON}
COPTFLAGS               += ${CFLAGS}
LDFLAGS                 += ${LDFLAGS_COMMON}
RUSTFLAGS               +=  -C target-cpu=x86-64
```