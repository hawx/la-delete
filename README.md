# la-delete

Delets all scrobbles after a duration. (Use
<https://secure.last.fm/settings/account> to clear listening data before first
use, or you will a) be waiting forever, and b) probably get your api key
blocked.)

Put credentials in a file (or pass on command line, see `--help`) like,

``` toml
apiKey = "..."
apiSecret = "..."
username = "..."
password = "..."
```

Password can also be set as the md5 hash of your password instead.

``` bash
$ go get hawx.me/code/la-delete
$ la-delete --auth auth.conf --after 72h --save ./someplace
```

If `--save` is given the deleted scrobbles are saved in the folder as `.json`
documents named with the time of the scrobble.
