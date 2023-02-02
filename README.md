# INTRODUCTION

During vm migartion testing, understanding the changes of packet loss, timeout, etc. on the fly is important to know if the migration is working smoothly across all related phases. To support the test, a simple tool based on pro-bing is implemented which supports intermediate ping statistics as below:

```
# sudo ./icmp_stat -h pkg.go.dev -v -i 32    
PING pkg.go.dev (34.149.140.181):
2023-02-02T14:25:41.933438196+08:00: 2 packets transmitted(delta 0), 1 packets received(delta 0), 50.00% packet loss, 35 ms rtt(delta 0)
2023-02-02T14:25:41.998902515+08:00: 4 packets transmitted(delta 2), 2 packets received(delta 1), 50.00% packet loss, 35 ms rtt(delta 0)
2023-02-02T14:25:42.063075991+08:00: 6 packets transmitted(delta 2), 3 packets received(delta 1), 50.00% packet loss, 35 ms rtt(delta 0)
2023-02-02T14:25:42.126978735+08:00: 8 packets transmitted(delta 2), 4 packets received(delta 1), 50.00% packet loss, 35 ms rtt(delta 0)
2023-02-02T14:25:42.190598733+08:00: 10 packets transmitted(delta 2), 5 packets received(delta 1), 50.00% packet loss, 35 ms rtt(delta 0)

--- pkg.go.dev ping statistics ---
10 packets transmitted, 5 packets received, 50.00% packet loss
round-trip min/avg/max/stddev = 35.057512ms/35.440028ms/35.691217ms/225.584µs
```

# INSTALLATION

- Download the required binary from the release page directly;
- Or build binaries as below:

  ```
  make x86_64
  make all
  ```

# USAGE

```
sudo ./icmp_stat --help
sudo ./icmp_stat --host pkg.go.dev
sudo ./icmp_stat -v -h pkg.go.dev
```
