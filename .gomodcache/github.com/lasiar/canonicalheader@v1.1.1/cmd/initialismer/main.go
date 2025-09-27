package main

import (
	"cmp"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"slices"
)

var mustBeIgnore = [...]string{
	"A-IM",
	"Accept",
	"Accept-Additions",
	"Accept-CH",
	"Accept-Charset",
	"Accept-Datetime",
	"Accept-Encoding",
	"Accept-Features",
	"Accept-Language",
	"Accept-Patch",
	"Accept-Post",
	"Accept-Ranges",
	"Accept-Signature",
	"Access-Control",
	"Access-Control-Allow-Credentials",
	"Access-Control-Allow-Headers",
	"Access-Control-Allow-Methods",
	"Access-Control-Allow-Origin",
	"Access-Control-Expose-Headers",
	"Access-Control-Max-Age",
	"Access-Control-Request-Headers",
	"Access-Control-Request-Method",
	"Age",
	"Allow",
	"ALPN",
	"Alt-Svc",
	"Alt-Used",
	"Alternates",
	"AMP-Cache-Transform",
	"Apply-To-Redirect-Ref",
	"Authentication-Control",
	"Authentication-Info",
	"Authorization",
	"C-Ext",
	"C-Man",
	"C-Opt",
	"C-PEP",
	"C-PEP-Info",
	"Cache-Control",
	"Cache-Status",
	"Cal-Managed-ID",
	"CalDAV-Timezones",
	"Capsule-Protocol",
	"CDN-Cache-Control",
	"CDN-Loop",
	"Cert-Not-After",
	"Cert-Not-Before",
	"Clear-Site-Data",
	"Client-Cert",
	"Client-Cert-Chain",
	"Close",
	"Configuration-Context",
	"Connection",
	"Content-Base",
	"Content-Digest",
	"Content-Disposition",
	"Content-Encoding",
	"Content-ID",
	"Content-Language",
	"Content-Length",
	"Content-Location",
	"Content-MD5",
	"Content-Range",
	"Content-Script-Type",
	"Content-Security-Policy",
	"Content-Security-Policy-Report-Only",
	"Content-Style-Type",
	"Content-Type",
	"Content-Version",
	"Cookie",
	"Cookie2",
	"Cross-Origin-Embedder-Policy",
	"Cross-Origin-Embedder-Policy-Report-Only",
	"Cross-Origin-Opener-Policy",
	"Cross-Origin-Opener-Policy-Report-Only",
	"Cross-Origin-Resource-Policy",
	"DASL",
	"Date",
	"DAV",
	"Default-Style",
	"Delta-Base",
	"Depth",
	"Derived-From",
	"Destination",
	"Differential-ID",
	"Digest",
	"DPoP",
	"DPoP-Nonce",
	"Early-Data",
	"EDIINT-Features",
	"Expect",
	"Expect-CT",
	"X-Correlation-ID",
	"X-UA-Compatible",
	"X-XSS-Protection",
	"Expires",
	"Ext",
	"Forwarded",
	"From",
	"GetProfile",
	"Hobareg",
	"Host",
	"HTTP2-Settings",
	"If",
	"If-Match",
	"If-Modified-Since",
	"If-None-Match",
	"If-Range",
	"If-Schedule-Tag-Match",
	"If-Unmodified-Since",
	"IM",
	"Include-Referred-Token-Binding-ID",
	"Isolation",
	"Keep-Alive",
	"Label",
	"Last-Event-ID",
	"Last-Modified",
	"Link",
	"Location",
	"Lock-Token",
	"Man",
	"Max-Forwards",
	"Memento-Datetime",
	"Meter",
	"Method-Check",
	"Method-Check-Expires",
	"MIME-Version",
	"Negotiate",
	"NEL",
	"OData-EntityId",
	"OData-Isolation",
	"OData-MaxVersion",
	"OData-Version",
	"Opt",
	"Optional-WWW-Authenticate",
	"Ordering-Type",
	"Origin",
	"Origin-Agent-Cluster",
	"OSCORE",
	"OSLC-Core-Version",
	"Overwrite",
	"P3P",
	"PEP",
	"PEP-Info",
	"Permissions-Policy",
	"PICS-Label",
	"Ping-From",
	"Ping-To",
	"Position",
	"Pragma",
	"Prefer",
	"Preference-Applied",
	"Priority",
	"ProfileObject",
	"Protocol",
	"Protocol-Info",
	"Protocol-Query",
	"Protocol-Request",
	"Proxy-Authenticate",
	"Proxy-Authentication-Info",
	"Proxy-Authorization",
	"Proxy-Features",
	"Proxy-Instruction",
	"Proxy-Status",
	"Public",
	"Public-Key-Pins",
	"Public-Key-Pins-Report-Only",
	"Range",
	"Redirect-Ref",
	"Referer",
	"Referer-Root",
	"Refresh",
	"Repeatability-Client-ID",
	"Repeatability-First-Sent",
	"Repeatability-Request-ID",
	"Repeatability-Result",
	"Replay-Nonce",
	"Reporting-Endpoints",
	"Repr-Digest",
	"Retry-After",
	"Safe",
	"Schedule-Reply",
	"Schedule-Tag",
	"Sec-GPC",
	"Sec-Purpose",
	"Sec-Token-Binding",
	"Sec-WebSocket-Accept",
	"Sec-WebSocket-Extensions",
	"Sec-WebSocket-Key",
	"Sec-WebSocket-Protocol",
	"Sec-WebSocket-Version",
	"Security-Scheme",
	"Server",
	"Server-Timing",
	"Set-Cookie",
	"Set-Cookie2",
	"SetProfile",
	"Signature",
	"Signature-Input",
	"SLUG",
	"SoapAction",
	"Status-URI",
	"Strict-Transport-Security",
	"Sunset",
	"Surrogate-Capability",
	"Surrogate-Control",
	"TCN",
	"TE",
	"Timeout",
	"Timing-Allow-Origin",
	"Topic",
	"Traceparent",
	"Tracestate",
	"Trailer",
	"Transfer-Encoding",
	"TTL",
	"Upgrade",
	"Urgency",
	"URI",
	"User-Agent",
	"Variant-Vary",
	"Vary",
	"Via",
	"Want-Content-Digest",
	"Want-Digest",
	"Want-Repr-Digest",
	"Warning",
	"X-Content-Type-Options",
	"X-Frame-Options",
	"ETag",
	"DNT",
	"X-Request-ID",
	"X-XSS",
	"X-DNS-Prefetch-Control",
	"WWW-Authenticate",
	"X-WebKit-CSP",
	"X-Real-IP",
}

type generateTarget uint8

func (t generateTarget) String() string {
	switch t {
	case generateTest:
		return "test"
	case generateTestGolden:
		return "test-golden"
	case generateMapping:
		return "mapping"
	default:
		return "unknown"
	}
}

const (
	generateUnknown generateTarget = iota
	generateTest
	generateTestGolden
	generateMapping
)

func parseTarget() (generateTarget, error) {
	var t string
	flag.StringVar(&t, "target", "", "test or mapping")
	flag.Parse()

	switch t {
	case generateTest.String():
		return generateTest, nil

	case generateMapping.String():
		return generateMapping, nil

	case generateTestGolden.String():
		return generateTestGolden, nil

	default:
		return generateUnknown, fmt.Errorf("unknown target %q", t)
	}
}

func main() {
	filtered := make(map[string]struct{}, len(mustBeIgnore))

	type ignore struct {
		Canonical,
		Original string
	}

	results := make([]ignore, 0, len(mustBeIgnore))

	for _, s := range mustBeIgnore {
		_, ok := filtered[s]
		if ok {
			slog.Error("has duplicate:", slog.String("duplicate", s))
			os.Exit(1)
		}
		filtered[s] = struct{}{}

		canonical := http.CanonicalHeaderKey(s)
		if canonical != s {
			results = append(results, ignore{
				Canonical: canonical,
				Original:  s,
			})
		}
	}
	slices.SortFunc(results, func(a, b ignore) int {
		return cmp.Compare(a.Canonical, b.Canonical)
	})

	genTarget, err := parseTarget()
	if err != nil {
		slog.Error("parse target:", slog.Any("error", err))
		os.Exit(1)
	}

	var tmpl string
	switch genTarget {
	case generateTest:
		tmpl = tmplTest

	case generateMapping:
		tmpl = tmplMapping

	case generateTestGolden:
		tmpl = tmplTestGolden
	}

	t := template.Must(template.New("").Parse(tmpl))

	err = t.Execute(os.Stdout, results)
	if err == nil {
		return
	}

	slog.Error("execute template:", slog.Any("error", err))
	os.Exit(1)
}

const tmplMapping = `// Code generated by initialismer; DO NOT EDIT.
package canonicalheader

// initialism mapping of not canonical headers from
// https://en.wikipedia.org/wiki/List_of_HTTP_header_fields
// https://www.iana.org/assignments/http-fields/http-fields.xhtml.
func initialism() map[string]string {
    return map[string]string{
        {{range .}}"{{.Canonical}}": "{{.Original}}",
        {{end}}
    }
}
`

const tmplTest = `// Code generated by initialismer; DO NOT EDIT.
package initialism

import "net/http"

func _() {
    h := http.Header{}
    {{range .}}
    h.Get("{{.Original}}"){{end}}
}
`

const tmplTestGolden = `// Code generated by initialismer; DO NOT EDIT.
package initialism

import "net/http"

func _() {
    h := http.Header{}
    {{range .}}
    h.Get("{{.Original}}"){{end}}
}
`
