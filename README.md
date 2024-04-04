# in-addr

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=jimmystewpot_in-addr&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=jimmystewpot_in-addr)

A simple tool to take a subnet mask and print out the in-addr reverse order to stdout.

This simplifies creating entire network ranges of files for DNS reverse lookup etc.


## IPv4

```

$ ./ptr generate 192.168.15.99/20
0.168.192.in-addr.arpa.
1.168.192.in-addr.arpa.
2.168.192.in-addr.arpa.
3.168.192.in-addr.arpa.
4.168.192.in-addr.arpa.
5.168.192.in-addr.arpa.
6.168.192.in-addr.arpa.
7.168.192.in-addr.arpa.
8.168.192.in-addr.arpa.
9.168.192.in-addr.arpa.
10.168.192.in-addr.arpa.
11.168.192.in-addr.arpa.
12.168.192.in-addr.arpa.
13.168.192.in-addr.arpa.
14.168.192.in-addr.arpa.
15.168.192.in-addr.arpa.
```

## IPv6

```
$ ./ptr generate 2001:db8:abcd:1234::1/56
0.0.2.1.d.c.b.a.8.b.d.0.1.0.0.2.in-addr.arpa.
1.0.2.1.d.c.b.a.8.b.d.0.1.0.0.2.in-addr.arpa.
2.0.2.1.d.c.b.a.8.b.d.0.1.0.0.2.in-addr.arpa.
3.0.2.1.d.c.b.a.8.b.d.0.1.0.0.2.in-addr.arpa.
4.0.2.1.d.c.b.a.8.b.d.0.1.0.0.2.in-addr.arpa.
5.0.2.1.d.c.b.a.8.b.d.0.1.0.0.2.in-addr.arpa.
6.0.2.1.d.c.b.a.8.b.d.0.1.0.0.2.in-addr.arpa.
....
```

## Building

```
go build -o ptr cmd/in-addr/main.go
```