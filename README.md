# AutoK3s GEO

AutoK3s GEO collects metrics about locates remote IP-address and exposes metrics to InfluxDB.

Thanks to https://freegeoip.live/ which provides a public HTTP(S) API for software developers to search the geolocation of IP addresses. It uses a database of IP addresses that are associated to cities along with other relevant information like time zone, latitude and longitude.

Please be aware that API has a limit of 50K requests per hour.

# Quick Start
```bash
# InfluxDB server base URL and an authentication token.
$ export ENDPOINT=http://x.x.x.x:x
$ export TOKEN=xxxxxxxxxxxxxxxxxxx

# Run serve command.
$ autok3s-geo serve
```

# Data Model / Line Protocol

```shell
+-----------+--------+-+---------+-+---------+
|measurement|,tag_set| |field_set| |timestamp|
+-----------+--------+-+---------+-+---------+

geo_locations,repository=autok3s,ip=192.168.1.2,country=China,city=ShenYang,latitude=35.6882,longitude=139.7532, active=14 1465839830100400200
```

# Retention Policy

- Retention Policy: 365d
- Shard Group Duration: 24h

# License

Copyright (c) 2022 [Rancher Labs, Inc.](http://rancher.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
