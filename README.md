# ifaddr
ifaddr is command-line tool for getting ip address that assigned local interfaces.

Replaces `ip addr show eth0 | awk '/inet / {print $2}' | cut -d"/" -f1` to `ifaddr -i eth0 -1`

## Usage
```
Usage:
  ifaddr [OPTIONS] [Pattern]

Application Options:
  -i, --interface=         Set interface
  -m, --netmask=           Filter address by netmask
  -1, --first              Show only first address
  -E, --regexp             Enable regexp for pattern
      --include-ipv6       Include IPv6 address(default)
      --exclude-ipv6       Exclude IPv6 address
      --include-ipv4       Include IPv4 address(default)
      --exclude-ipv4       Exclude IPv4 address
      --include-linklocal  Include Link-Local address
      --exclude-linklocal  Exclude Link-Local address(default)
      --include-loopback   Include Loopback address
      --exclude-loopback   Exclude Loopback address(default)
  -a, --all                Show all addresses(--include-linklocal, --include-loopback)
  -6, --only-ipv6          Show only IPv6 address(--exclude-ipv4)
  -4, --only-ipv4          Show only IPv6 address(--exclude-ipv6)

Help Options:
  -h, --help               Show this help message

Arguments:
  Pattern
```

### Show all IP addresses
```
$ ifaddr -a
::1
127.0.0.1
fe80::1
192.0.2.1
2001:DB8::1
2001:DB8::2
198.51.100.1
203.0.113.1
```

### Exclude LinkLocal/Loopback

default behavior

```
$ ifaddr
192.0.2.1
2001:DB8::1
2001:DB8::2
198.51.100.1
203.0.113.1
```

### Show only IPv4 address

```
$ ifaddr -4
192.0.2.1
198.51.100.1
203.0.113.1
```

### Show only IPv4 address(include Loopback)

```
$ ifaddr -4 --include-loopback
192.0.2.1
198.51.100.1
203.0.113.1
```

### Show only IPv6 address

```
$ ifaddr -6
2001:DB8::1
2001:DB8::2
```

### Show only first IPv6 address

```
$ ifaddr -6 -1
2001:DB8::1
```

### Show IP address by interface name

```
$ ifaddr -i eth0
192.0.2.1
```

### Show IP address by netmask

```
$ ifaddr -m 203.0.113.0/24
203.0.113.1
```

### Show IP address by pattern

```
$ ifaddr 192
192.0.2.1
```

### Show IP address by regexp pattern

```
$ ifaddr -E '2$'
2001:DB8::2
```

## LICENSE
MIT
