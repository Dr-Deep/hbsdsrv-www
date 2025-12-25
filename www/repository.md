# Repositories

## About hbsdsrv's Repos
[repo list](/poudriere/packages/)


* goldmont-plus: x86-64 architecture optimised for goldmont-plus processor 
`(SSE,SSE2,SSE3,SSSE3,SSE4.1,SSE4.2,POPCNT,CMPXCHG16B,AES-NI,CLMUL,SHA,RDRAND,RDSEED,LAHF/SAHF,MOBBE)`
* skylake: x86-64-v3 architecture optimised for skylake processor
`(SSE, SSE2, SSE3, SSSE3, SSE4.1, SSE4.2, AVX, AVX2, FMA, AES-NI, PCLMULQDQ, BMI1, BMI2, POPCNT, MOVBE, RDRAND, RDSEED, LAHF/SAHF, CMPXCHG16B, PREFETCHW, PREFETCHWT1, HLE, RTM)`
* lto: here i test which ports can support LTO and PGO, in the future we merge them in the main repos
* _test: please do not add them

The flag -march=x86-64-v3 enables all instruction sets required by the x86-64-v3 ABI level (like AVX2, FMA, BMI1/2),
making the binary incompatible with older CPUs.
The -mtune=skylake flag tells the compiler to optimize code generation (e.g. instruction scheduling, cache usage) for Skylake CPUs
without using any instructions beyond what's allowed by -march.
Together, this produces code that runs on all x86-64-v3 CPUs but performs best on Skylake.


## How to add

edit /etc/pkg/hbsdsrv.conf
and paste this snippet:
```ucl
hbsdsrv-pkgbase: {
        url: "https://hbsdsrv.1337.cx/poudriere/packages/XXX",
        mirror_type: "none",
        signature_type: "pubkey",
        pubkey: "/etc/pkg/hbsdsrv.cert",
        enabled: yes,
        priority: 0,
}

hbsdsrv-ports: {
	url: "https://hbsdsrv.1337.cx/poudriere/packages/XXX",
	mirror_type: "none",
	signature_type: "pubkey",
	pubkey: "/etc/pkg/hbsdsrv.cert",
	enabled: yes,
	priority: 0,
}
```
copy&paste the package repo of your choice in the `url` param.

- *.conf options described [here](https://man.freebsd.org/cgi/man.cgi?pkg#CONFIGURATION)
- [Public Key of the Repo](/hbsdsrv.cert)
