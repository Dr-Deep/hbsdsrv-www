```
# /usr/local/etc/poudriere.d/make.conf hsbdsrv by @Dr.Deep

# CCACHE
WITH_CCACHE_BUILD=yes
CCACHE_DIR=/usr/ccache

# Nur ccache NICHT für sich selbst verwenden
.if ${.CURDIR:M*/devel/ccache*}
NO_CCACHE= yes
.endif

# Poudriere
USE_PACKAGE_DEPENDS=yes
BATCH=yes
FORCE_PACKAGE=yes
PACKAGE_BUILDING=yes
PACKAGE_BUILDING_FLAVORS=yes
DISABLE_MAKE_JOBS=poudriere
#FETCH_CMD=aria2c
#FETCH_BEFORE_ARGS=-x 10 -s 10

NO_CPU_CFLAGS=yes		# Don't add -march=<cpu> to CFLAGS automatically

# SIMD
CFLAGS_COMMON            = -O3 -pipe -march=x86-64 -falign-functions=32 -fno-semantic-interposition -g1 -gz=zstd -gno-column-info
# -funroll-loops -finline-functions -ftree-vectorize

# -march -mtune
# -O3				Aggressive Opt
# -falign-functions=32		Fetch-Performance
# -fno-lto			Kein Link-Time-Optimization for PGO verwendet wird.
# -fno-semantic-interposition	better cpu-cache use, inlining, branch prediction
# -g1 
# -gz=zstd			: compress dbg infos
# -gno-column-info		: no loc info

LDFLAGS_COMMON		= -Wl,--gc-sections -Wl,-z,max-page-size=0x1000
# -Wl,--gc-sections		: Entfernt nicht genutzte .text/.data Sections
# -Wl,-z,max-page-size=0x1000 	: page-size=4KiB, better for Cache/TLB (default:64KiB align)


#-fno-strict-overflow	# 
#-fcf-protection=full	# Intel CET (Control-Flow-Protetion)

# ports that do not compile with CFLAGS_COMMON
COMMON_BLACKLIST = graphics/wayland

.for PORT in ${COMMON_BLACKLIST}
.if !empty(.CURDIR:M*/${PORT}*)

CFLAGS_COMMON = -O3 -pipe -march=x86-64
LDFLAGS_COMMON = 

.endif
.endfor

CFLAGS			+= ${CFLAGS_COMMON}
CXXFLAGS 		+= ${CFLAGS_COMMON}
COPTFLAGS		+= ${CFLAGS}
LDFLAGS  		+= ${LDFLAGS_COMMON}
RUSTFLAGS               +=  -C target-cpu=x86-64

#HARDENING_ALL=		cfi fortifysource pie relro retpoline safestack slh zeroreg
# PIE, RELRO+BIND_NOW
#USE_HARDENING		+= pie relro
USE_HARDENING += auto

LICENSES_ACCEPTED = *
DISABLE_LICENSES=yes
```
