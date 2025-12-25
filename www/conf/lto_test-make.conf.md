```
# /etc/main_lto-make.conf hsbdsrv by @Dr.Deep

# -flto / -flto=full -flto=thin	# LTO
LTO_FLAGS		= -flto=thin -g

# ports that do not compile with LTO
PORT_BLACKLIST = audio/libogg devel/libmtdev devel/libedit textproc/libunibreak lang/tcl audio/speexdsp audio/libvorbis devel/pcre devel/pkgconf lang/perl devel/libffi devel/gettext-runtime devel/libdatrie devel/libtextstyle libinotify multimedia/libdvdread devel/boehm-gc www/libnghttp devel/libevent devel/popt devel/libunwind net/libngtcp2 audio/flac devel/libuv textproc/libyaml devel/libltdl x11/libxshmfence audio/alsa-lib x11/libXau x11/libICE devel/libthai textproc/libxml2 devel/gettext-tools security/cyrus-sasl2 audio/opus multimedia/libvpx security/libsodium multimedia/libdvdnav database/gdbm security/libgpg-error lang/python graphics/libexif textproc/hunspell dns/libidn math/gmp filesystems/fusefs-libs print/libpaper math/fftw3 textproc/libucl multimedia/libx264

.for PORT in ${PORT_BLACKLIST}

.if !empty(.CURDIR:M*/${PORT}*)
CFLAGS			+= -fno-lto
CXXFLAGS		+= -fno-lto

.else
CFLAGS			+= ${LTO_FLAGS}
CXXFLAGS		+= ${LTO_FLAGS}

.endif
.endfor

OPTIONS_SET		+= LTO
```
