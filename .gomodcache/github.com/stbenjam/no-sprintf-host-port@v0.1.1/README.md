# no-sprintf-host-port 

The Go linter no-sprintf-host-port checks that sprintf is not used to
construct a host:port combination in a URL.    A frequent pattern is for a
developer to construct a URL like this:

```go
fmt.Sprintf("http://%s:%d/foo", host, port)
```

However, if "host" is an IPv6 address like `2001:4860:4860::8888`, the
URL constructed will be invalid. IPv6 addresses must be bracketed, like this:

```
http://[2001:4860:4860::8888]:9443
```

The linter  is naive, and really only looks for the most obvious cases, but where
it's possible to infer that a URL is being constructed with  Sprintf containing a `:`,
this informs the user to use `net.JoinHostPort` instead.

## Thanks

Based on the [`go-printf-func-name`](https://github.com/jirfag/go-printf-func-name) linter,
and this [article](https://disaev.me/p/writing-useful-go-analysis-linter/).
