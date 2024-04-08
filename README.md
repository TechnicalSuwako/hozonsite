# 保存サイト
世界初FOSS系ウエブアーカイバーです。

# インストールする方法
## OpenBSD
```sh
cd hozonsite
make
doas make install
```

## FreeBSD
```sh
cd hozonsite
make
doas make install MANPREFIX=/usr/local/share/man
```

## NetBSD
```sh
cd hozonsite
make
doas make install PREFIX=/usr/pkg MANPREFIX=/usr/pkg/share/man
```

## Linux
```sh
cd hozonsite
bmake
doas bmake install PREFIX=/usr MANPREFIX=/usr/share/man
```
