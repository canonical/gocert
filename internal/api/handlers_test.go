package server_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	server "github.com/canonical/gocert/internal/api"
	"github.com/canonical/gocert/internal/certdb"
	"github.com/golang-jwt/jwt"
)

const (
	validCSR1 = `-----BEGIN CERTIFICATE REQUEST-----
MIICszCCAZsCAQAwFjEUMBIGA1UEAwwLZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3
DQEBAQUAA4IBDwAwggEKAoIBAQDC5KgrADpuOUPwSh0YLmpWF66VTcciIGC2HcGn
oJknL7pm5q9qhfWGIdvKKlIA6cBB32jPd0QcYDsx7+AvzEvBuO7mq7v2Q1sPU4Q+
L0s2pLJges6/cnDWvk/p5eBjDLOqHhUNzpMUga9SgIod8yymTZm3eqQvt1ABdwTg
FzBs5QdSm2Ny1fEbbcRE+Rv5rqXyJb2isXSujzSuS22VqslDIyqnY5WaLg+pjZyR
+0j13ecJsdh6/MJMUZWheimV2Yv7SFtxzFwbzBMO9YFS098sy4F896eBHLNe9cUC
+d1JDtLaewlMogjHBHAxmP54dhe6vvc78anElKKP4hm5N5nlAgMBAAGgWDBWBgkq
hkiG9w0BCQ4xSTBHMA4GA1UdDwEB/wQEAwIFoDAdBgNVHSUEFjAUBggrBgEFBQcD
AQYIKwYBBQUHAwIwFgYDVR0RBA8wDYILZXhhbXBsZS5jb20wDQYJKoZIhvcNAQEL
BQADggEBACP1VKEGVYKoVLMDJS+EZ0CPwIYWsO4xBXgK6atHe8WIChVn/8I7eo60
cuMDiy4LR70G++xL1tpmYGRbx21r9d/shL2ehp9VdClX06qxlcGxiC/F8eThRuS5
zHcdNqSVyMoLJ0c7yWHJahN5u2bn1Lov34yOEqGGpWCGF/gT1nEvM+p/v30s89f2
Y/uPl4g3jpGqLCKTASWJDGnZLroLICOzYTVs5P3oj+VueSUwYhGK5tBnS2x5FHID
uMNMgwl0fxGMQZjrlXyCBhXBm1k6PmwcJGJF5LQ31c+5aTTMFU7SyZhlymctB8mS
y+ErBQsRpcQho6Ok+HTXQQUcx7WNcwI=
-----END CERTIFICATE REQUEST-----`
	validCSR2 = `-----BEGIN CERTIFICATE REQUEST-----
MIIC5zCCAc8CAQAwRzEWMBQGA1UEAwwNMTAuMTUyLjE4My41MzEtMCsGA1UELQwk
MzlhY2UxOTUtZGM1YS00MzJiLTgwOTAtYWZlNmFiNGI0OWNmMIIBIjANBgkqhkiG
9w0BAQEFAAOCAQ8AMIIBCgKCAQEAjM5Wz+HRtDveRzeDkEDM4ornIaefe8d8nmFi
pUat9qCU3U9798FR460DHjCLGxFxxmoRitzHtaR4ew5H036HlGB20yas/CMDgSUI
69DyAsyPwEJqOWBGO1LL50qXdl5/jOkO2voA9j5UsD1CtWSklyhbNhWMpYqj2ObW
XcaYj9Gx/TwYhw8xsJ/QRWyCrvjjVzH8+4frfDhBVOyywN7sq+I3WwCbyBBcN8uO
yae0b/q5+UJUiqgpeOAh/4Y7qI3YarMj4cm7dwmiCVjedUwh65zVyHtQUfLd8nFW
Kl9775mNBc1yicvKDU3ZB5hZ1MZtpbMBwaA1yMSErs/fh5KaXwIDAQABoFswWQYJ
KoZIhvcNAQkOMUwwSjBIBgNVHREEQTA/hwQKmLc1gjd2YXVsdC1rOHMtMC52YXVs
dC1rOHMtZW5kcG9pbnRzLnZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsMA0GCSqGSIb3
DQEBCwUAA4IBAQCJt8oVDbiuCsik4N5AOJIT7jKsMb+j0mizwjahKMoCHdx+zv0V
FGkhlf0VWPAdEu3gHdJfduX88WwzJ2wBBUK38UuprAyvfaZfaYUgFJQNC6DH1fIa
uHYEhvNJBdFJHaBvW7lrSFi57fTA9IEPrB3m/XN3r2F4eoHnaJJqHZmMwqVHck87
cAQXk3fvTWuikHiCHqqdSdjDYj/8cyiwCrQWpV245VSbOE0WesWoEnSdFXVUfE1+
RSKeTRuuJMcdGqBkDnDI22myj0bjt7q8eqBIjTiLQLnAFnQYpcCrhc8dKU9IJlv1
H9Hay4ZO9LRew3pEtlx2WrExw/gpUcWM8rTI
-----END CERTIFICATE REQUEST-----`
	validCSR3 = `-----BEGIN CERTIFICATE REQUEST-----
MIICszCCAZsCAQAwFjEUMBIGA1UEAwwLZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3
DQEBAQUAA4IBDwAwggEKAoIBAQDN7tHggWTtxiT5Sh5Npoif8J2BdpJjtMdpZ7Vu
NVzMxW/eojSRlq0p3nafmpjnSdSH1k/XMmPsgmv9txxEHMw1LIUJUef2QVrQTI6J
4ueu9NvexZWXZ+UxFip63PKyn/CkZRFiHCRIGzDDPxM2aApjghXy9ISMtGqDVSnr
5hQDu2U1CEiUWKMoTpyk/KlBZliDDOzaGm3cQuzKWs6Stjzpq+uX4ecJAXZg5Cj+
+JUETH93A/VOfsiiHXoKeTnFMCsmJgEHz2DZixw8EN8XgpOp5BA2n8Y/xS+Ren5R
ZH7uNJI/SmQ0yrR+2bYR6hm+4bCzspyCfzbiuI5IS9+2eXA/AgMBAAGgWDBWBgkq
hkiG9w0BCQ4xSTBHMA4GA1UdDwEB/wQEAwIFoDAdBgNVHSUEFjAUBggrBgEFBQcD
AQYIKwYBBQUHAwIwFgYDVR0RBA8wDYILZXhhbXBsZS5jb20wDQYJKoZIhvcNAQEL
BQADggEBAB/aPfYLbnCubYyKnxLRipoLr3TBSYFnRfcxiZR1o+L3/tuv2NlrXJjY
K13xzzPhwuZwd6iKfX3xC33sKgnUNFawyE8IuAmyhJ2cl97iA2lwoYcyuWP9TOEx
LT60zxp7PHsKo53gqaqRJ5B9RZtiv1jYdUZvynHP4J5JG7Zwaa0VNi/Cx5cwGW8K
rfvNABPUAU6xIqqYgd2heDPF6kjvpoNiOl056qIAbk0dbmpqOJf/lxKBRfqlHhSC
0qRScGu70l2Oxl89YSsfGtUyQuzTkLshI2VkEUM+W/ZauXbxLd8SyWveH3/7mDC+
Sgi7T+lz+c1Tw+XFgkqryUwMeG2wxt8=
-----END CERTIFICATE REQUEST-----`
	validCert2 = `-----BEGIN CERTIFICATE-----
MIIDrDCCApSgAwIBAgIURKr+jf7hj60SyAryIeN++9wDdtkwDQYJKoZIhvcNAQEL
BQAwOTELMAkGA1UEBhMCVVMxKjAoBgNVBAMMIXNlbGYtc2lnbmVkLWNlcnRpZmlj
YXRlcy1vcGVyYXRvcjAeFw0yNDAzMjcxMjQ4MDRaFw0yNTAzMjcxMjQ4MDRaMEcx
FjAUBgNVBAMMDTEwLjE1Mi4xODMuNTMxLTArBgNVBC0MJDM5YWNlMTk1LWRjNWEt
NDMyYi04MDkwLWFmZTZhYjRiNDljZjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC
AQoCggEBAIzOVs/h0bQ73kc3g5BAzOKK5yGnn3vHfJ5hYqVGrfaglN1Pe/fBUeOt
Ax4wixsRccZqEYrcx7WkeHsOR9N+h5RgdtMmrPwjA4ElCOvQ8gLMj8BCajlgRjtS
y+dKl3Zef4zpDtr6APY+VLA9QrVkpJcoWzYVjKWKo9jm1l3GmI/Rsf08GIcPMbCf
0EVsgq7441cx/PuH63w4QVTsssDe7KviN1sAm8gQXDfLjsmntG/6uflCVIqoKXjg
If+GO6iN2GqzI+HJu3cJoglY3nVMIeuc1ch7UFHy3fJxVipfe++ZjQXNconLyg1N
2QeYWdTGbaWzAcGgNcjEhK7P34eSml8CAwEAAaOBnTCBmjAhBgNVHSMEGjAYgBYE
FN/vgl9cAapV7hH9lEyM7qYS958aMB0GA1UdDgQWBBRJJDZkHr64VqTC24DPQVld
Ba3iPDAMBgNVHRMBAf8EAjAAMEgGA1UdEQRBMD+CN3ZhdWx0LWs4cy0wLnZhdWx0
LWs4cy1lbmRwb2ludHMudmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWyHBAqYtzUwDQYJ
KoZIhvcNAQELBQADggEBAEH9NTwDiSsoQt/QXkWPMBrB830K0dlwKl5WBNgVxFP+
hSfQ86xN77jNSp2VxOksgzF9J9u/ubAXvSFsou4xdP8MevBXoFJXeqMERq5RW3gc
WyhXkzguv3dwH+n43GJFP6MQ+n9W/nPZCUQ0Iy7ueAvj0HFhGyZzAE2wxNFZdvCs
gCX3nqYpp70oZIFDrhmYwE5ij5KXlHD4/1IOfNUKCDmQDgGPLI1tVtwQLjeRq7Hg
XVelpl/LXTQawmJyvDaVT/Q9P+WqoDiMjrqF6Sy7DzNeeccWVqvqX5TVS6Ky56iS
Mvo/+PAJHkBciR5Xn+Wg2a+7vrZvT6CBoRSOTozlLSM=
-----END CERTIFICATE-----`
)

const (
	expectedGetAllCertsResponseBody1 = "[{\"id\":1,\"csr\":\"-----BEGIN CERTIFICATE REQUEST-----\\nMIICszCCAZsCAQAwFjEUMBIGA1UEAwwLZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3\\nDQEBAQUAA4IBDwAwggEKAoIBAQDC5KgrADpuOUPwSh0YLmpWF66VTcciIGC2HcGn\\noJknL7pm5q9qhfWGIdvKKlIA6cBB32jPd0QcYDsx7+AvzEvBuO7mq7v2Q1sPU4Q+\\nL0s2pLJges6/cnDWvk/p5eBjDLOqHhUNzpMUga9SgIod8yymTZm3eqQvt1ABdwTg\\nFzBs5QdSm2Ny1fEbbcRE+Rv5rqXyJb2isXSujzSuS22VqslDIyqnY5WaLg+pjZyR\\n+0j13ecJsdh6/MJMUZWheimV2Yv7SFtxzFwbzBMO9YFS098sy4F896eBHLNe9cUC\\n+d1JDtLaewlMogjHBHAxmP54dhe6vvc78anElKKP4hm5N5nlAgMBAAGgWDBWBgkq\\nhkiG9w0BCQ4xSTBHMA4GA1UdDwEB/wQEAwIFoDAdBgNVHSUEFjAUBggrBgEFBQcD\\nAQYIKwYBBQUHAwIwFgYDVR0RBA8wDYILZXhhbXBsZS5jb20wDQYJKoZIhvcNAQEL\\nBQADggEBACP1VKEGVYKoVLMDJS+EZ0CPwIYWsO4xBXgK6atHe8WIChVn/8I7eo60\\ncuMDiy4LR70G++xL1tpmYGRbx21r9d/shL2ehp9VdClX06qxlcGxiC/F8eThRuS5\\nzHcdNqSVyMoLJ0c7yWHJahN5u2bn1Lov34yOEqGGpWCGF/gT1nEvM+p/v30s89f2\\nY/uPl4g3jpGqLCKTASWJDGnZLroLICOzYTVs5P3oj+VueSUwYhGK5tBnS2x5FHID\\nuMNMgwl0fxGMQZjrlXyCBhXBm1k6PmwcJGJF5LQ31c+5aTTMFU7SyZhlymctB8mS\\ny+ErBQsRpcQho6Ok+HTXQQUcx7WNcwI=\\n-----END CERTIFICATE REQUEST-----\",\"certificate\":\"\"}]"
	expectedGetAllCertsResponseBody2 = "[{\"id\":1,\"csr\":\"-----BEGIN CERTIFICATE REQUEST-----\\nMIICszCCAZsCAQAwFjEUMBIGA1UEAwwLZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3\\nDQEBAQUAA4IBDwAwggEKAoIBAQDC5KgrADpuOUPwSh0YLmpWF66VTcciIGC2HcGn\\noJknL7pm5q9qhfWGIdvKKlIA6cBB32jPd0QcYDsx7+AvzEvBuO7mq7v2Q1sPU4Q+\\nL0s2pLJges6/cnDWvk/p5eBjDLOqHhUNzpMUga9SgIod8yymTZm3eqQvt1ABdwTg\\nFzBs5QdSm2Ny1fEbbcRE+Rv5rqXyJb2isXSujzSuS22VqslDIyqnY5WaLg+pjZyR\\n+0j13ecJsdh6/MJMUZWheimV2Yv7SFtxzFwbzBMO9YFS098sy4F896eBHLNe9cUC\\n+d1JDtLaewlMogjHBHAxmP54dhe6vvc78anElKKP4hm5N5nlAgMBAAGgWDBWBgkq\\nhkiG9w0BCQ4xSTBHMA4GA1UdDwEB/wQEAwIFoDAdBgNVHSUEFjAUBggrBgEFBQcD\\nAQYIKwYBBQUHAwIwFgYDVR0RBA8wDYILZXhhbXBsZS5jb20wDQYJKoZIhvcNAQEL\\nBQADggEBACP1VKEGVYKoVLMDJS+EZ0CPwIYWsO4xBXgK6atHe8WIChVn/8I7eo60\\ncuMDiy4LR70G++xL1tpmYGRbx21r9d/shL2ehp9VdClX06qxlcGxiC/F8eThRuS5\\nzHcdNqSVyMoLJ0c7yWHJahN5u2bn1Lov34yOEqGGpWCGF/gT1nEvM+p/v30s89f2\\nY/uPl4g3jpGqLCKTASWJDGnZLroLICOzYTVs5P3oj+VueSUwYhGK5tBnS2x5FHID\\nuMNMgwl0fxGMQZjrlXyCBhXBm1k6PmwcJGJF5LQ31c+5aTTMFU7SyZhlymctB8mS\\ny+ErBQsRpcQho6Ok+HTXQQUcx7WNcwI=\\n-----END CERTIFICATE REQUEST-----\",\"certificate\":\"\"},{\"id\":2,\"csr\":\"-----BEGIN CERTIFICATE REQUEST-----\\nMIIC5zCCAc8CAQAwRzEWMBQGA1UEAwwNMTAuMTUyLjE4My41MzEtMCsGA1UELQwk\\nMzlhY2UxOTUtZGM1YS00MzJiLTgwOTAtYWZlNmFiNGI0OWNmMIIBIjANBgkqhkiG\\n9w0BAQEFAAOCAQ8AMIIBCgKCAQEAjM5Wz+HRtDveRzeDkEDM4ornIaefe8d8nmFi\\npUat9qCU3U9798FR460DHjCLGxFxxmoRitzHtaR4ew5H036HlGB20yas/CMDgSUI\\n69DyAsyPwEJqOWBGO1LL50qXdl5/jOkO2voA9j5UsD1CtWSklyhbNhWMpYqj2ObW\\nXcaYj9Gx/TwYhw8xsJ/QRWyCrvjjVzH8+4frfDhBVOyywN7sq+I3WwCbyBBcN8uO\\nyae0b/q5+UJUiqgpeOAh/4Y7qI3YarMj4cm7dwmiCVjedUwh65zVyHtQUfLd8nFW\\nKl9775mNBc1yicvKDU3ZB5hZ1MZtpbMBwaA1yMSErs/fh5KaXwIDAQABoFswWQYJ\\nKoZIhvcNAQkOMUwwSjBIBgNVHREEQTA/hwQKmLc1gjd2YXVsdC1rOHMtMC52YXVs\\ndC1rOHMtZW5kcG9pbnRzLnZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsMA0GCSqGSIb3\\nDQEBCwUAA4IBAQCJt8oVDbiuCsik4N5AOJIT7jKsMb+j0mizwjahKMoCHdx+zv0V\\nFGkhlf0VWPAdEu3gHdJfduX88WwzJ2wBBUK38UuprAyvfaZfaYUgFJQNC6DH1fIa\\nuHYEhvNJBdFJHaBvW7lrSFi57fTA9IEPrB3m/XN3r2F4eoHnaJJqHZmMwqVHck87\\ncAQXk3fvTWuikHiCHqqdSdjDYj/8cyiwCrQWpV245VSbOE0WesWoEnSdFXVUfE1+\\nRSKeTRuuJMcdGqBkDnDI22myj0bjt7q8eqBIjTiLQLnAFnQYpcCrhc8dKU9IJlv1\\nH9Hay4ZO9LRew3pEtlx2WrExw/gpUcWM8rTI\\n-----END CERTIFICATE REQUEST-----\",\"certificate\":\"\"}]"
	expectedGetAllCertsResponseBody3 = "[{\"id\":2,\"csr\":\"-----BEGIN CERTIFICATE REQUEST-----\\nMIIC5zCCAc8CAQAwRzEWMBQGA1UEAwwNMTAuMTUyLjE4My41MzEtMCsGA1UELQwk\\nMzlhY2UxOTUtZGM1YS00MzJiLTgwOTAtYWZlNmFiNGI0OWNmMIIBIjANBgkqhkiG\\n9w0BAQEFAAOCAQ8AMIIBCgKCAQEAjM5Wz+HRtDveRzeDkEDM4ornIaefe8d8nmFi\\npUat9qCU3U9798FR460DHjCLGxFxxmoRitzHtaR4ew5H036HlGB20yas/CMDgSUI\\n69DyAsyPwEJqOWBGO1LL50qXdl5/jOkO2voA9j5UsD1CtWSklyhbNhWMpYqj2ObW\\nXcaYj9Gx/TwYhw8xsJ/QRWyCrvjjVzH8+4frfDhBVOyywN7sq+I3WwCbyBBcN8uO\\nyae0b/q5+UJUiqgpeOAh/4Y7qI3YarMj4cm7dwmiCVjedUwh65zVyHtQUfLd8nFW\\nKl9775mNBc1yicvKDU3ZB5hZ1MZtpbMBwaA1yMSErs/fh5KaXwIDAQABoFswWQYJ\\nKoZIhvcNAQkOMUwwSjBIBgNVHREEQTA/hwQKmLc1gjd2YXVsdC1rOHMtMC52YXVs\\ndC1rOHMtZW5kcG9pbnRzLnZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsMA0GCSqGSIb3\\nDQEBCwUAA4IBAQCJt8oVDbiuCsik4N5AOJIT7jKsMb+j0mizwjahKMoCHdx+zv0V\\nFGkhlf0VWPAdEu3gHdJfduX88WwzJ2wBBUK38UuprAyvfaZfaYUgFJQNC6DH1fIa\\nuHYEhvNJBdFJHaBvW7lrSFi57fTA9IEPrB3m/XN3r2F4eoHnaJJqHZmMwqVHck87\\ncAQXk3fvTWuikHiCHqqdSdjDYj/8cyiwCrQWpV245VSbOE0WesWoEnSdFXVUfE1+\\nRSKeTRuuJMcdGqBkDnDI22myj0bjt7q8eqBIjTiLQLnAFnQYpcCrhc8dKU9IJlv1\\nH9Hay4ZO9LRew3pEtlx2WrExw/gpUcWM8rTI\\n-----END CERTIFICATE REQUEST-----\",\"certificate\":\"-----BEGIN CERTIFICATE-----\\nMIIDrDCCApSgAwIBAgIURKr+jf7hj60SyAryIeN++9wDdtkwDQYJKoZIhvcNAQEL\\nBQAwOTELMAkGA1UEBhMCVVMxKjAoBgNVBAMMIXNlbGYtc2lnbmVkLWNlcnRpZmlj\\nYXRlcy1vcGVyYXRvcjAeFw0yNDAzMjcxMjQ4MDRaFw0yNTAzMjcxMjQ4MDRaMEcx\\nFjAUBgNVBAMMDTEwLjE1Mi4xODMuNTMxLTArBgNVBC0MJDM5YWNlMTk1LWRjNWEt\\nNDMyYi04MDkwLWFmZTZhYjRiNDljZjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC\\nAQoCggEBAIzOVs/h0bQ73kc3g5BAzOKK5yGnn3vHfJ5hYqVGrfaglN1Pe/fBUeOt\\nAx4wixsRccZqEYrcx7WkeHsOR9N+h5RgdtMmrPwjA4ElCOvQ8gLMj8BCajlgRjtS\\ny+dKl3Zef4zpDtr6APY+VLA9QrVkpJcoWzYVjKWKo9jm1l3GmI/Rsf08GIcPMbCf\\n0EVsgq7441cx/PuH63w4QVTsssDe7KviN1sAm8gQXDfLjsmntG/6uflCVIqoKXjg\\nIf+GO6iN2GqzI+HJu3cJoglY3nVMIeuc1ch7UFHy3fJxVipfe++ZjQXNconLyg1N\\n2QeYWdTGbaWzAcGgNcjEhK7P34eSml8CAwEAAaOBnTCBmjAhBgNVHSMEGjAYgBYE\\nFN/vgl9cAapV7hH9lEyM7qYS958aMB0GA1UdDgQWBBRJJDZkHr64VqTC24DPQVld\\nBa3iPDAMBgNVHRMBAf8EAjAAMEgGA1UdEQRBMD+CN3ZhdWx0LWs4cy0wLnZhdWx0\\nLWs4cy1lbmRwb2ludHMudmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWyHBAqYtzUwDQYJ\\nKoZIhvcNAQELBQADggEBAEH9NTwDiSsoQt/QXkWPMBrB830K0dlwKl5WBNgVxFP+\\nhSfQ86xN77jNSp2VxOksgzF9J9u/ubAXvSFsou4xdP8MevBXoFJXeqMERq5RW3gc\\nWyhXkzguv3dwH+n43GJFP6MQ+n9W/nPZCUQ0Iy7ueAvj0HFhGyZzAE2wxNFZdvCs\\ngCX3nqYpp70oZIFDrhmYwE5ij5KXlHD4/1IOfNUKCDmQDgGPLI1tVtwQLjeRq7Hg\\nXVelpl/LXTQawmJyvDaVT/Q9P+WqoDiMjrqF6Sy7DzNeeccWVqvqX5TVS6Ky56iS\\nMvo/+PAJHkBciR5Xn+Wg2a+7vrZvT6CBoRSOTozlLSM=\\n-----END CERTIFICATE-----\"},{\"id\":3,\"csr\":\"-----BEGIN CERTIFICATE REQUEST-----\\nMIICszCCAZsCAQAwFjEUMBIGA1UEAwwLZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3\\nDQEBAQUAA4IBDwAwggEKAoIBAQDN7tHggWTtxiT5Sh5Npoif8J2BdpJjtMdpZ7Vu\\nNVzMxW/eojSRlq0p3nafmpjnSdSH1k/XMmPsgmv9txxEHMw1LIUJUef2QVrQTI6J\\n4ueu9NvexZWXZ+UxFip63PKyn/CkZRFiHCRIGzDDPxM2aApjghXy9ISMtGqDVSnr\\n5hQDu2U1CEiUWKMoTpyk/KlBZliDDOzaGm3cQuzKWs6Stjzpq+uX4ecJAXZg5Cj+\\n+JUETH93A/VOfsiiHXoKeTnFMCsmJgEHz2DZixw8EN8XgpOp5BA2n8Y/xS+Ren5R\\nZH7uNJI/SmQ0yrR+2bYR6hm+4bCzspyCfzbiuI5IS9+2eXA/AgMBAAGgWDBWBgkq\\nhkiG9w0BCQ4xSTBHMA4GA1UdDwEB/wQEAwIFoDAdBgNVHSUEFjAUBggrBgEFBQcD\\nAQYIKwYBBQUHAwIwFgYDVR0RBA8wDYILZXhhbXBsZS5jb20wDQYJKoZIhvcNAQEL\\nBQADggEBAB/aPfYLbnCubYyKnxLRipoLr3TBSYFnRfcxiZR1o+L3/tuv2NlrXJjY\\nK13xzzPhwuZwd6iKfX3xC33sKgnUNFawyE8IuAmyhJ2cl97iA2lwoYcyuWP9TOEx\\nLT60zxp7PHsKo53gqaqRJ5B9RZtiv1jYdUZvynHP4J5JG7Zwaa0VNi/Cx5cwGW8K\\nrfvNABPUAU6xIqqYgd2heDPF6kjvpoNiOl056qIAbk0dbmpqOJf/lxKBRfqlHhSC\\n0qRScGu70l2Oxl89YSsfGtUyQuzTkLshI2VkEUM+W/ZauXbxLd8SyWveH3/7mDC+\\nSgi7T+lz+c1Tw+XFgkqryUwMeG2wxt8=\\n-----END CERTIFICATE REQUEST-----\",\"certificate\":\"\"},{\"id\":4,\"csr\":\"-----BEGIN CERTIFICATE REQUEST-----\\nMIICszCCAZsCAQAwFjEUMBIGA1UEAwwLZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3\\nDQEBAQUAA4IBDwAwggEKAoIBAQDC5KgrADpuOUPwSh0YLmpWF66VTcciIGC2HcGn\\noJknL7pm5q9qhfWGIdvKKlIA6cBB32jPd0QcYDsx7+AvzEvBuO7mq7v2Q1sPU4Q+\\nL0s2pLJges6/cnDWvk/p5eBjDLOqHhUNzpMUga9SgIod8yymTZm3eqQvt1ABdwTg\\nFzBs5QdSm2Ny1fEbbcRE+Rv5rqXyJb2isXSujzSuS22VqslDIyqnY5WaLg+pjZyR\\n+0j13ecJsdh6/MJMUZWheimV2Yv7SFtxzFwbzBMO9YFS098sy4F896eBHLNe9cUC\\n+d1JDtLaewlMogjHBHAxmP54dhe6vvc78anElKKP4hm5N5nlAgMBAAGgWDBWBgkq\\nhkiG9w0BCQ4xSTBHMA4GA1UdDwEB/wQEAwIFoDAdBgNVHSUEFjAUBggrBgEFBQcD\\nAQYIKwYBBQUHAwIwFgYDVR0RBA8wDYILZXhhbXBsZS5jb20wDQYJKoZIhvcNAQEL\\nBQADggEBACP1VKEGVYKoVLMDJS+EZ0CPwIYWsO4xBXgK6atHe8WIChVn/8I7eo60\\ncuMDiy4LR70G++xL1tpmYGRbx21r9d/shL2ehp9VdClX06qxlcGxiC/F8eThRuS5\\nzHcdNqSVyMoLJ0c7yWHJahN5u2bn1Lov34yOEqGGpWCGF/gT1nEvM+p/v30s89f2\\nY/uPl4g3jpGqLCKTASWJDGnZLroLICOzYTVs5P3oj+VueSUwYhGK5tBnS2x5FHID\\nuMNMgwl0fxGMQZjrlXyCBhXBm1k6PmwcJGJF5LQ31c+5aTTMFU7SyZhlymctB8mS\\ny+ErBQsRpcQho6Ok+HTXQQUcx7WNcwI=\\n-----END CERTIFICATE REQUEST-----\",\"certificate\":\"rejected\"}]"
	expectedGetAllCertsResponseBody4 = "[{\"id\":2,\"csr\":\"-----BEGIN CERTIFICATE REQUEST-----\\nMIIC5zCCAc8CAQAwRzEWMBQGA1UEAwwNMTAuMTUyLjE4My41MzEtMCsGA1UELQwk\\nMzlhY2UxOTUtZGM1YS00MzJiLTgwOTAtYWZlNmFiNGI0OWNmMIIBIjANBgkqhkiG\\n9w0BAQEFAAOCAQ8AMIIBCgKCAQEAjM5Wz+HRtDveRzeDkEDM4ornIaefe8d8nmFi\\npUat9qCU3U9798FR460DHjCLGxFxxmoRitzHtaR4ew5H036HlGB20yas/CMDgSUI\\n69DyAsyPwEJqOWBGO1LL50qXdl5/jOkO2voA9j5UsD1CtWSklyhbNhWMpYqj2ObW\\nXcaYj9Gx/TwYhw8xsJ/QRWyCrvjjVzH8+4frfDhBVOyywN7sq+I3WwCbyBBcN8uO\\nyae0b/q5+UJUiqgpeOAh/4Y7qI3YarMj4cm7dwmiCVjedUwh65zVyHtQUfLd8nFW\\nKl9775mNBc1yicvKDU3ZB5hZ1MZtpbMBwaA1yMSErs/fh5KaXwIDAQABoFswWQYJ\\nKoZIhvcNAQkOMUwwSjBIBgNVHREEQTA/hwQKmLc1gjd2YXVsdC1rOHMtMC52YXVs\\ndC1rOHMtZW5kcG9pbnRzLnZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsMA0GCSqGSIb3\\nDQEBCwUAA4IBAQCJt8oVDbiuCsik4N5AOJIT7jKsMb+j0mizwjahKMoCHdx+zv0V\\nFGkhlf0VWPAdEu3gHdJfduX88WwzJ2wBBUK38UuprAyvfaZfaYUgFJQNC6DH1fIa\\nuHYEhvNJBdFJHaBvW7lrSFi57fTA9IEPrB3m/XN3r2F4eoHnaJJqHZmMwqVHck87\\ncAQXk3fvTWuikHiCHqqdSdjDYj/8cyiwCrQWpV245VSbOE0WesWoEnSdFXVUfE1+\\nRSKeTRuuJMcdGqBkDnDI22myj0bjt7q8eqBIjTiLQLnAFnQYpcCrhc8dKU9IJlv1\\nH9Hay4ZO9LRew3pEtlx2WrExw/gpUcWM8rTI\\n-----END CERTIFICATE REQUEST-----\",\"certificate\":\"\"},{\"id\":3,\"csr\":\"-----BEGIN CERTIFICATE REQUEST-----\\nMIICszCCAZsCAQAwFjEUMBIGA1UEAwwLZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3\\nDQEBAQUAA4IBDwAwggEKAoIBAQDN7tHggWTtxiT5Sh5Npoif8J2BdpJjtMdpZ7Vu\\nNVzMxW/eojSRlq0p3nafmpjnSdSH1k/XMmPsgmv9txxEHMw1LIUJUef2QVrQTI6J\\n4ueu9NvexZWXZ+UxFip63PKyn/CkZRFiHCRIGzDDPxM2aApjghXy9ISMtGqDVSnr\\n5hQDu2U1CEiUWKMoTpyk/KlBZliDDOzaGm3cQuzKWs6Stjzpq+uX4ecJAXZg5Cj+\\n+JUETH93A/VOfsiiHXoKeTnFMCsmJgEHz2DZixw8EN8XgpOp5BA2n8Y/xS+Ren5R\\nZH7uNJI/SmQ0yrR+2bYR6hm+4bCzspyCfzbiuI5IS9+2eXA/AgMBAAGgWDBWBgkq\\nhkiG9w0BCQ4xSTBHMA4GA1UdDwEB/wQEAwIFoDAdBgNVHSUEFjAUBggrBgEFBQcD\\nAQYIKwYBBQUHAwIwFgYDVR0RBA8wDYILZXhhbXBsZS5jb20wDQYJKoZIhvcNAQEL\\nBQADggEBAB/aPfYLbnCubYyKnxLRipoLr3TBSYFnRfcxiZR1o+L3/tuv2NlrXJjY\\nK13xzzPhwuZwd6iKfX3xC33sKgnUNFawyE8IuAmyhJ2cl97iA2lwoYcyuWP9TOEx\\nLT60zxp7PHsKo53gqaqRJ5B9RZtiv1jYdUZvynHP4J5JG7Zwaa0VNi/Cx5cwGW8K\\nrfvNABPUAU6xIqqYgd2heDPF6kjvpoNiOl056qIAbk0dbmpqOJf/lxKBRfqlHhSC\\n0qRScGu70l2Oxl89YSsfGtUyQuzTkLshI2VkEUM+W/ZauXbxLd8SyWveH3/7mDC+\\nSgi7T+lz+c1Tw+XFgkqryUwMeG2wxt8=\\n-----END CERTIFICATE REQUEST-----\",\"certificate\":\"\"},{\"id\":4,\"csr\":\"-----BEGIN CERTIFICATE REQUEST-----\\nMIICszCCAZsCAQAwFjEUMBIGA1UEAwwLZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3\\nDQEBAQUAA4IBDwAwggEKAoIBAQDC5KgrADpuOUPwSh0YLmpWF66VTcciIGC2HcGn\\noJknL7pm5q9qhfWGIdvKKlIA6cBB32jPd0QcYDsx7+AvzEvBuO7mq7v2Q1sPU4Q+\\nL0s2pLJges6/cnDWvk/p5eBjDLOqHhUNzpMUga9SgIod8yymTZm3eqQvt1ABdwTg\\nFzBs5QdSm2Ny1fEbbcRE+Rv5rqXyJb2isXSujzSuS22VqslDIyqnY5WaLg+pjZyR\\n+0j13ecJsdh6/MJMUZWheimV2Yv7SFtxzFwbzBMO9YFS098sy4F896eBHLNe9cUC\\n+d1JDtLaewlMogjHBHAxmP54dhe6vvc78anElKKP4hm5N5nlAgMBAAGgWDBWBgkq\\nhkiG9w0BCQ4xSTBHMA4GA1UdDwEB/wQEAwIFoDAdBgNVHSUEFjAUBggrBgEFBQcD\\nAQYIKwYBBQUHAwIwFgYDVR0RBA8wDYILZXhhbXBsZS5jb20wDQYJKoZIhvcNAQEL\\nBQADggEBACP1VKEGVYKoVLMDJS+EZ0CPwIYWsO4xBXgK6atHe8WIChVn/8I7eo60\\ncuMDiy4LR70G++xL1tpmYGRbx21r9d/shL2ehp9VdClX06qxlcGxiC/F8eThRuS5\\nzHcdNqSVyMoLJ0c7yWHJahN5u2bn1Lov34yOEqGGpWCGF/gT1nEvM+p/v30s89f2\\nY/uPl4g3jpGqLCKTASWJDGnZLroLICOzYTVs5P3oj+VueSUwYhGK5tBnS2x5FHID\\nuMNMgwl0fxGMQZjrlXyCBhXBm1k6PmwcJGJF5LQ31c+5aTTMFU7SyZhlymctB8mS\\ny+ErBQsRpcQho6Ok+HTXQQUcx7WNcwI=\\n-----END CERTIFICATE REQUEST-----\",\"certificate\":\"rejected\"}]"
	expectedGetCertReqResponseBody1  = "{\"id\":2,\"csr\":\"-----BEGIN CERTIFICATE REQUEST-----\\nMIIC5zCCAc8CAQAwRzEWMBQGA1UEAwwNMTAuMTUyLjE4My41MzEtMCsGA1UELQwk\\nMzlhY2UxOTUtZGM1YS00MzJiLTgwOTAtYWZlNmFiNGI0OWNmMIIBIjANBgkqhkiG\\n9w0BAQEFAAOCAQ8AMIIBCgKCAQEAjM5Wz+HRtDveRzeDkEDM4ornIaefe8d8nmFi\\npUat9qCU3U9798FR460DHjCLGxFxxmoRitzHtaR4ew5H036HlGB20yas/CMDgSUI\\n69DyAsyPwEJqOWBGO1LL50qXdl5/jOkO2voA9j5UsD1CtWSklyhbNhWMpYqj2ObW\\nXcaYj9Gx/TwYhw8xsJ/QRWyCrvjjVzH8+4frfDhBVOyywN7sq+I3WwCbyBBcN8uO\\nyae0b/q5+UJUiqgpeOAh/4Y7qI3YarMj4cm7dwmiCVjedUwh65zVyHtQUfLd8nFW\\nKl9775mNBc1yicvKDU3ZB5hZ1MZtpbMBwaA1yMSErs/fh5KaXwIDAQABoFswWQYJ\\nKoZIhvcNAQkOMUwwSjBIBgNVHREEQTA/hwQKmLc1gjd2YXVsdC1rOHMtMC52YXVs\\ndC1rOHMtZW5kcG9pbnRzLnZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsMA0GCSqGSIb3\\nDQEBCwUAA4IBAQCJt8oVDbiuCsik4N5AOJIT7jKsMb+j0mizwjahKMoCHdx+zv0V\\nFGkhlf0VWPAdEu3gHdJfduX88WwzJ2wBBUK38UuprAyvfaZfaYUgFJQNC6DH1fIa\\nuHYEhvNJBdFJHaBvW7lrSFi57fTA9IEPrB3m/XN3r2F4eoHnaJJqHZmMwqVHck87\\ncAQXk3fvTWuikHiCHqqdSdjDYj/8cyiwCrQWpV245VSbOE0WesWoEnSdFXVUfE1+\\nRSKeTRuuJMcdGqBkDnDI22myj0bjt7q8eqBIjTiLQLnAFnQYpcCrhc8dKU9IJlv1\\nH9Hay4ZO9LRew3pEtlx2WrExw/gpUcWM8rTI\\n-----END CERTIFICATE REQUEST-----\",\"certificate\":\"\"}"
	expectedGetCertReqResponseBody2  = "{\"id\":4,\"csr\":\"-----BEGIN CERTIFICATE REQUEST-----\\nMIICszCCAZsCAQAwFjEUMBIGA1UEAwwLZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3\\nDQEBAQUAA4IBDwAwggEKAoIBAQDC5KgrADpuOUPwSh0YLmpWF66VTcciIGC2HcGn\\noJknL7pm5q9qhfWGIdvKKlIA6cBB32jPd0QcYDsx7+AvzEvBuO7mq7v2Q1sPU4Q+\\nL0s2pLJges6/cnDWvk/p5eBjDLOqHhUNzpMUga9SgIod8yymTZm3eqQvt1ABdwTg\\nFzBs5QdSm2Ny1fEbbcRE+Rv5rqXyJb2isXSujzSuS22VqslDIyqnY5WaLg+pjZyR\\n+0j13ecJsdh6/MJMUZWheimV2Yv7SFtxzFwbzBMO9YFS098sy4F896eBHLNe9cUC\\n+d1JDtLaewlMogjHBHAxmP54dhe6vvc78anElKKP4hm5N5nlAgMBAAGgWDBWBgkq\\nhkiG9w0BCQ4xSTBHMA4GA1UdDwEB/wQEAwIFoDAdBgNVHSUEFjAUBggrBgEFBQcD\\nAQYIKwYBBQUHAwIwFgYDVR0RBA8wDYILZXhhbXBsZS5jb20wDQYJKoZIhvcNAQEL\\nBQADggEBACP1VKEGVYKoVLMDJS+EZ0CPwIYWsO4xBXgK6atHe8WIChVn/8I7eo60\\ncuMDiy4LR70G++xL1tpmYGRbx21r9d/shL2ehp9VdClX06qxlcGxiC/F8eThRuS5\\nzHcdNqSVyMoLJ0c7yWHJahN5u2bn1Lov34yOEqGGpWCGF/gT1nEvM+p/v30s89f2\\nY/uPl4g3jpGqLCKTASWJDGnZLroLICOzYTVs5P3oj+VueSUwYhGK5tBnS2x5FHID\\nuMNMgwl0fxGMQZjrlXyCBhXBm1k6PmwcJGJF5LQ31c+5aTTMFU7SyZhlymctB8mS\\ny+ErBQsRpcQho6Ok+HTXQQUcx7WNcwI=\\n-----END CERTIFICATE REQUEST-----\",\"certificate\":\"\"}"
	expectedGetCertReqResponseBody3  = "{\"id\":2,\"csr\":\"-----BEGIN CERTIFICATE REQUEST-----\\nMIIC5zCCAc8CAQAwRzEWMBQGA1UEAwwNMTAuMTUyLjE4My41MzEtMCsGA1UELQwk\\nMzlhY2UxOTUtZGM1YS00MzJiLTgwOTAtYWZlNmFiNGI0OWNmMIIBIjANBgkqhkiG\\n9w0BAQEFAAOCAQ8AMIIBCgKCAQEAjM5Wz+HRtDveRzeDkEDM4ornIaefe8d8nmFi\\npUat9qCU3U9798FR460DHjCLGxFxxmoRitzHtaR4ew5H036HlGB20yas/CMDgSUI\\n69DyAsyPwEJqOWBGO1LL50qXdl5/jOkO2voA9j5UsD1CtWSklyhbNhWMpYqj2ObW\\nXcaYj9Gx/TwYhw8xsJ/QRWyCrvjjVzH8+4frfDhBVOyywN7sq+I3WwCbyBBcN8uO\\nyae0b/q5+UJUiqgpeOAh/4Y7qI3YarMj4cm7dwmiCVjedUwh65zVyHtQUfLd8nFW\\nKl9775mNBc1yicvKDU3ZB5hZ1MZtpbMBwaA1yMSErs/fh5KaXwIDAQABoFswWQYJ\\nKoZIhvcNAQkOMUwwSjBIBgNVHREEQTA/hwQKmLc1gjd2YXVsdC1rOHMtMC52YXVs\\ndC1rOHMtZW5kcG9pbnRzLnZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsMA0GCSqGSIb3\\nDQEBCwUAA4IBAQCJt8oVDbiuCsik4N5AOJIT7jKsMb+j0mizwjahKMoCHdx+zv0V\\nFGkhlf0VWPAdEu3gHdJfduX88WwzJ2wBBUK38UuprAyvfaZfaYUgFJQNC6DH1fIa\\nuHYEhvNJBdFJHaBvW7lrSFi57fTA9IEPrB3m/XN3r2F4eoHnaJJqHZmMwqVHck87\\ncAQXk3fvTWuikHiCHqqdSdjDYj/8cyiwCrQWpV245VSbOE0WesWoEnSdFXVUfE1+\\nRSKeTRuuJMcdGqBkDnDI22myj0bjt7q8eqBIjTiLQLnAFnQYpcCrhc8dKU9IJlv1\\nH9Hay4ZO9LRew3pEtlx2WrExw/gpUcWM8rTI\\n-----END CERTIFICATE REQUEST-----\",\"certificate\":\"-----BEGIN CERTIFICATE-----\\nMIIDrDCCApSgAwIBAgIURKr+jf7hj60SyAryIeN++9wDdtkwDQYJKoZIhvcNAQEL\\nBQAwOTELMAkGA1UEBhMCVVMxKjAoBgNVBAMMIXNlbGYtc2lnbmVkLWNlcnRpZmlj\\nYXRlcy1vcGVyYXRvcjAeFw0yNDAzMjcxMjQ4MDRaFw0yNTAzMjcxMjQ4MDRaMEcx\\nFjAUBgNVBAMMDTEwLjE1Mi4xODMuNTMxLTArBgNVBC0MJDM5YWNlMTk1LWRjNWEt\\nNDMyYi04MDkwLWFmZTZhYjRiNDljZjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC\\nAQoCggEBAIzOVs/h0bQ73kc3g5BAzOKK5yGnn3vHfJ5hYqVGrfaglN1Pe/fBUeOt\\nAx4wixsRccZqEYrcx7WkeHsOR9N+h5RgdtMmrPwjA4ElCOvQ8gLMj8BCajlgRjtS\\ny+dKl3Zef4zpDtr6APY+VLA9QrVkpJcoWzYVjKWKo9jm1l3GmI/Rsf08GIcPMbCf\\n0EVsgq7441cx/PuH63w4QVTsssDe7KviN1sAm8gQXDfLjsmntG/6uflCVIqoKXjg\\nIf+GO6iN2GqzI+HJu3cJoglY3nVMIeuc1ch7UFHy3fJxVipfe++ZjQXNconLyg1N\\n2QeYWdTGbaWzAcGgNcjEhK7P34eSml8CAwEAAaOBnTCBmjAhBgNVHSMEGjAYgBYE\\nFN/vgl9cAapV7hH9lEyM7qYS958aMB0GA1UdDgQWBBRJJDZkHr64VqTC24DPQVld\\nBa3iPDAMBgNVHRMBAf8EAjAAMEgGA1UdEQRBMD+CN3ZhdWx0LWs4cy0wLnZhdWx0\\nLWs4cy1lbmRwb2ludHMudmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWyHBAqYtzUwDQYJ\\nKoZIhvcNAQELBQADggEBAEH9NTwDiSsoQt/QXkWPMBrB830K0dlwKl5WBNgVxFP+\\nhSfQ86xN77jNSp2VxOksgzF9J9u/ubAXvSFsou4xdP8MevBXoFJXeqMERq5RW3gc\\nWyhXkzguv3dwH+n43GJFP6MQ+n9W/nPZCUQ0Iy7ueAvj0HFhGyZzAE2wxNFZdvCs\\ngCX3nqYpp70oZIFDrhmYwE5ij5KXlHD4/1IOfNUKCDmQDgGPLI1tVtwQLjeRq7Hg\\nXVelpl/LXTQawmJyvDaVT/Q9P+WqoDiMjrqF6Sy7DzNeeccWVqvqX5TVS6Ky56iS\\nMvo/+PAJHkBciR5Xn+Wg2a+7vrZvT6CBoRSOTozlLSM=\\n-----END CERTIFICATE-----\"}"
	expectedGetCertReqResponseBody4  = "{\"id\":2,\"csr\":\"-----BEGIN CERTIFICATE REQUEST-----\\nMIIC5zCCAc8CAQAwRzEWMBQGA1UEAwwNMTAuMTUyLjE4My41MzEtMCsGA1UELQwk\\nMzlhY2UxOTUtZGM1YS00MzJiLTgwOTAtYWZlNmFiNGI0OWNmMIIBIjANBgkqhkiG\\n9w0BAQEFAAOCAQ8AMIIBCgKCAQEAjM5Wz+HRtDveRzeDkEDM4ornIaefe8d8nmFi\\npUat9qCU3U9798FR460DHjCLGxFxxmoRitzHtaR4ew5H036HlGB20yas/CMDgSUI\\n69DyAsyPwEJqOWBGO1LL50qXdl5/jOkO2voA9j5UsD1CtWSklyhbNhWMpYqj2ObW\\nXcaYj9Gx/TwYhw8xsJ/QRWyCrvjjVzH8+4frfDhBVOyywN7sq+I3WwCbyBBcN8uO\\nyae0b/q5+UJUiqgpeOAh/4Y7qI3YarMj4cm7dwmiCVjedUwh65zVyHtQUfLd8nFW\\nKl9775mNBc1yicvKDU3ZB5hZ1MZtpbMBwaA1yMSErs/fh5KaXwIDAQABoFswWQYJ\\nKoZIhvcNAQkOMUwwSjBIBgNVHREEQTA/hwQKmLc1gjd2YXVsdC1rOHMtMC52YXVs\\ndC1rOHMtZW5kcG9pbnRzLnZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsMA0GCSqGSIb3\\nDQEBCwUAA4IBAQCJt8oVDbiuCsik4N5AOJIT7jKsMb+j0mizwjahKMoCHdx+zv0V\\nFGkhlf0VWPAdEu3gHdJfduX88WwzJ2wBBUK38UuprAyvfaZfaYUgFJQNC6DH1fIa\\nuHYEhvNJBdFJHaBvW7lrSFi57fTA9IEPrB3m/XN3r2F4eoHnaJJqHZmMwqVHck87\\ncAQXk3fvTWuikHiCHqqdSdjDYj/8cyiwCrQWpV245VSbOE0WesWoEnSdFXVUfE1+\\nRSKeTRuuJMcdGqBkDnDI22myj0bjt7q8eqBIjTiLQLnAFnQYpcCrhc8dKU9IJlv1\\nH9Hay4ZO9LRew3pEtlx2WrExw/gpUcWM8rTI\\n-----END CERTIFICATE REQUEST-----\",\"certificate\":\"\"}"
)

const (
	adminUser              = `{"username": "testadmin", "password": "Admin123"}`
	validUser              = `{"username": "testuser", "password": "userPass!"}`
	invalidUser            = `{"username": "", "password": ""}`
	noPasswordUser         = `{"username": "nopass"}`
	adminUserNewPassword   = `{"id": 1, "password": "newPassword1"}`
	userNewInvalidPassword = `{"id": 1, "password": "password"}`
	userMissingPassword    = `{"id": 1}`
	adminUserWrongPass     = `{"username": "testadmin", "password": "wrongpass"}`
	notExistingUser        = `{"username": "not_existing", "password": "user"}`
)

func TestGoCertCertificatesHandlers(t *testing.T) {
	testdb, err := certdb.NewCertificateRequestsRepository(":memory:", "CertificateRequests")
	if err != nil {
		log.Fatalf("couldn't create test sqlite db: %s", err)
	}
	env := &server.Environment{}
	env.DB = testdb
	ts := httptest.NewTLSServer(server.NewGoCertRouter(env))
	defer ts.Close()

	client := ts.Client()

	var adminToken string
	var nonAdminToken string
	t.Run("prepare user accounts and tokens", prepareUserAccounts(ts.URL, client, &adminToken, &nonAdminToken))

	testCases := []struct {
		desc     string
		method   string
		path     string
		data     string
		response string
		status   int
	}{
		{
			desc:     "healthcheck success",
			method:   "GET",
			path:     "/status",
			data:     "",
			response: "",
			status:   http.StatusOK,
		},
		{
			desc:     "empty get csrs success",
			method:   "GET",
			path:     "/api/v1/certificate_requests",
			data:     "",
			response: "null",
			status:   http.StatusOK,
		},
		{
			desc:     "post csr1 fail",
			method:   "POST",
			path:     "/api/v1/certificate_requests",
			data:     "this is very clearly not a csr",
			response: "error: csr validation failed: PEM Certificate Request string not found or malformed",
			status:   http.StatusBadRequest,
		},
		{
			desc:     "post csr1 success",
			method:   "POST",
			path:     "/api/v1/certificate_requests",
			data:     validCSR1,
			response: "1",
			status:   http.StatusCreated,
		},
		{
			desc:     "get csrs 1 success",
			method:   "GET",
			path:     "/api/v1/certificate_requests",
			data:     "",
			response: expectedGetAllCertsResponseBody1,
			status:   http.StatusOK,
		},
		{
			desc:     "post csr2 success",
			method:   "POST",
			path:     "/api/v1/certificate_requests",
			data:     validCSR2,
			response: "2",
			status:   http.StatusCreated,
		},
		{
			desc:     "get csrs 2 success",
			method:   "GET",
			path:     "/api/v1/certificate_requests",
			data:     "",
			response: expectedGetAllCertsResponseBody2,
			status:   http.StatusOK,
		},
		{
			desc:     "post csr2 fail",
			method:   "POST",
			path:     "/api/v1/certificate_requests",
			data:     validCSR2,
			response: "error: given csr already recorded",
			status:   http.StatusBadRequest,
		},
		{
			desc:     "post csr3 success",
			method:   "POST",
			path:     "/api/v1/certificate_requests",
			data:     validCSR3,
			response: "3",
			status:   http.StatusCreated,
		},
		{
			desc:     "delete csr1 success",
			method:   "DELETE",
			path:     "/api/v1/certificate_requests/1",
			data:     "",
			response: "1",
			status:   http.StatusAccepted,
		},
		{
			desc:     "delete csr5 fail",
			method:   "DELETE",
			path:     "/api/v1/certificate_requests/5",
			data:     "",
			response: "error: id not found",
			status:   http.StatusNotFound,
		},
		{
			desc:     "get csr1 fail",
			method:   "GET",
			path:     "/api/v1/certificate_requests/1",
			data:     "",
			response: "error: id not found",
			status:   http.StatusNotFound,
		},
		{
			desc:     "get csr2 success",
			method:   "GET",
			path:     "/api/v1/certificate_requests/2",
			data:     "",
			response: expectedGetCertReqResponseBody1,
			status:   http.StatusOK,
		},
		{
			desc:     "post csr4 success",
			method:   "POST",
			path:     "/api/v1/certificate_requests",
			data:     validCSR1,
			response: "4",
			status:   http.StatusCreated,
		},
		{
			desc:     "get csr4 success",
			method:   "GET",
			path:     "/api/v1/certificate_requests/4",
			data:     "",
			response: expectedGetCertReqResponseBody2,
			status:   http.StatusOK,
		},
		{
			desc:     "post cert2 fail 1",
			method:   "POST",
			path:     "/api/v1/certificate_requests/4/certificate",
			data:     validCert2,
			response: "error: cert validation failed: certificate does not match CSR",
			status:   http.StatusBadRequest,
		},
		{
			desc:     "post cert2 fail 2",
			method:   "POST",
			path:     "/api/v1/certificate_requests/4/certificate",
			data:     "some random data that's clearly not a cert",
			response: "error: cert validation failed: PEM Certificate string not found or malformed",
			status:   http.StatusBadRequest,
		},
		{
			desc:     "post cert2 success",
			method:   "POST",
			path:     "/api/v1/certificate_requests/2/certificate",
			data:     validCert2,
			response: "1",
			status:   http.StatusCreated,
		},
		{
			desc:     "get csr2 success",
			method:   "GET",
			path:     "/api/v1/certificate_requests/2",
			data:     "",
			response: expectedGetCertReqResponseBody3,
			status:   http.StatusOK,
		},
		{
			desc:     "reject csr4 success",
			method:   "POST",
			path:     "/api/v1/certificate_requests/4/certificate/reject",
			data:     "",
			response: "1",
			status:   http.StatusAccepted,
		},
		{
			desc:     "get all csrs success",
			method:   "GET",
			path:     "/api/v1/certificate_requests",
			data:     "",
			response: expectedGetAllCertsResponseBody3,
			status:   http.StatusOK,
		},
		{
			desc:     "delete csr2 cert success",
			method:   "DELETE",
			path:     "/api/v1/certificate_requests/2/certificate",
			data:     "",
			response: "1",
			status:   http.StatusAccepted,
		},
		{
			desc:     "get csr2 success",
			method:   "GET",
			path:     "/api/v1/certificate_requests/2",
			data:     "",
			response: expectedGetCertReqResponseBody4,
			status:   http.StatusOK,
		},
		{
			desc:     "get csrs 3 success",
			method:   "GET",
			path:     "/api/v1/certificate_requests",
			data:     "",
			response: expectedGetAllCertsResponseBody4,
			status:   http.StatusOK,
		},
		{
			desc:     "healthcheck success",
			method:   "GET",
			path:     "/status",
			data:     "",
			response: "",
			status:   http.StatusOK,
		},
		{
			desc:     "metrics endpoint success",
			method:   "GET",
			path:     "/metrics",
			data:     "",
			response: "",
			status:   http.StatusOK,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			req, err := http.NewRequest(tC.method, ts.URL+tC.path, strings.NewReader(tC.data))
			req.Header.Set("Authorization", "Bearer "+adminToken)
			if err != nil {
				t.Fatal(err)
			}
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			resBody, err := io.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			if res.StatusCode != tC.status || !strings.Contains(string(resBody), tC.response) {
				t.Errorf("expected response did not match.\nExpected vs Received status code: %d vs %d\nExpected vs Received body: \n%s\nvs\n%s\n", tC.status, res.StatusCode, tC.response, string(resBody))
			}
		})
	}

}

func TestGoCertUsersHandlers(t *testing.T) {
	testdb, err := certdb.NewCertificateRequestsRepository(":memory:", "CertificateRequests")
	if err != nil {
		log.Fatalf("couldn't create test sqlite db: %s", err)
	}
	env := &server.Environment{}
	env.DB = testdb
	ts := httptest.NewTLSServer(server.NewGoCertRouter(env))
	defer ts.Close()

	client := ts.Client()

	var adminToken string
	var nonAdminToken string
	t.Run("prepare user accounts and tokens", prepareUserAccounts(ts.URL, client, &adminToken, &nonAdminToken))

	testCases := []struct {
		desc     string
		method   string
		path     string
		data     string
		auth     string
		response string
		status   int
	}{
		{
			desc:     "Retrieve admin user success",
			method:   "GET",
			path:     "/api/v1/accounts/1",
			data:     "",
			auth:     adminToken,
			response: "{\"id\":1,\"username\":\"testadmin\",\"permissions\":1}",
			status:   http.StatusOK,
		},
		{
			desc:     "Retrieve admin user fail",
			method:   "GET",
			path:     "/api/v1/accounts/1",
			data:     "",
			auth:     nonAdminToken,
			response: "error: forbidden",
			status:   http.StatusForbidden,
		},
		{
			desc:     "Create no password user success",
			method:   "POST",
			path:     "/api/v1/accounts",
			data:     noPasswordUser,
			auth:     adminToken,
			response: "{\"id\":3,\"password\":",
			status:   http.StatusCreated,
		},
		{
			desc:     "Retrieve normal user success",
			method:   "GET",
			path:     "/api/v1/accounts/2",
			data:     "",
			auth:     adminToken,
			response: "{\"id\":2,\"username\":\"testuser\",\"permissions\":0}",
			status:   http.StatusOK,
		},
		{
			desc:     "Retrieve user failure",
			method:   "GET",
			path:     "/api/v1/accounts/300",
			data:     "",
			auth:     adminToken,
			response: "error: id not found",
			status:   http.StatusNotFound,
		},
		{
			desc:     "Create user failure",
			method:   "POST",
			path:     "/api/v1/accounts",
			data:     invalidUser,
			auth:     adminToken,
			response: "error: Username is required",
			status:   http.StatusBadRequest,
		},
		{
			desc:     "Change password success",
			method:   "POST",
			path:     "/api/v1/accounts/1/change_password",
			data:     adminUserNewPassword,
			auth:     adminToken,
			response: "1",
			status:   http.StatusOK,
		},
		{
			desc:     "Change password failure no user",
			method:   "POST",
			path:     "/api/v1/accounts/100/change_password",
			data:     adminUserNewPassword,
			auth:     adminToken,
			response: "id not found",
			status:   http.StatusNotFound,
		},
		{
			desc:     "Change password failure missing password",
			method:   "POST",
			path:     "/api/v1/accounts/1/change_password",
			data:     userMissingPassword,
			auth:     adminToken,
			response: "Password is required",
			status:   http.StatusBadRequest,
		},
		{
			desc:     "Change password failure bad password",
			method:   "POST",
			path:     "/api/v1/accounts/1/change_password",
			data:     userNewInvalidPassword,
			auth:     adminToken,
			response: "Password must have 8 or more characters, must include at least one capital letter, one lowercase letter, and either a number or a symbol.",
			status:   http.StatusBadRequest,
		},
		{
			desc:     "Delete user success",
			method:   "DELETE",
			path:     "/api/v1/accounts/2",
			data:     invalidUser,
			auth:     adminToken,
			response: "1",
			status:   http.StatusAccepted,
		},
		{
			desc:     "Delete user failure",
			method:   "DELETE",
			path:     "/api/v1/accounts/2",
			data:     invalidUser,
			auth:     adminToken,
			response: "error: id not found",
			status:   http.StatusNotFound,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			req, err := http.NewRequest(tC.method, ts.URL+tC.path, strings.NewReader(tC.data))
			req.Header.Add("Authorization", "Bearer "+tC.auth)
			if err != nil {
				t.Fatal(err)
			}
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			resBody, err := io.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			if res.StatusCode != tC.status || !strings.Contains(string(resBody), tC.response) {
				t.Errorf("expected response did not match.\nExpected vs Received status code: %d vs %d\nExpected vs Received body: \n%s\nvs\n%s\n", tC.status, res.StatusCode, tC.response, string(resBody))
			}
			if tC.desc == "Create no password user success" {
				match, _ := regexp.MatchString(`"password":"[!-~]{16}"`, string(resBody))
				if !match {
					t.Errorf("password does not match expected format or length: got %s", string(resBody))
				}
			}
		})
	}
}

func TestLogin(t *testing.T) {
	testdb, err := certdb.NewCertificateRequestsRepository(":memory:", "CertificateRequests")
	if err != nil {
		log.Fatalf("couldn't create test sqlite db: %s", err)
	}
	env := &server.Environment{}
	env.DB = testdb
	env.JWTSecret = []byte("secret")
	ts := httptest.NewTLSServer(server.NewGoCertRouter(env))
	defer ts.Close()

	client := ts.Client()

	testCases := []struct {
		desc     string
		method   string
		path     string
		data     string
		response string
		status   int
	}{
		{
			desc:     "Create admin user",
			method:   "POST",
			path:     "/api/v1/accounts",
			data:     adminUser,
			response: "{\"id\":1}",
			status:   http.StatusCreated,
		},
		{
			desc:     "Login success",
			method:   "POST",
			path:     "/login",
			data:     adminUser,
			response: "",
			status:   http.StatusOK,
		},
		{
			desc:     "Login failure missing username",
			method:   "POST",
			path:     "/login",
			data:     invalidUser,
			response: "Username is required",
			status:   http.StatusBadRequest,
		},
		{
			desc:     "Login failure missing password",
			method:   "POST",
			path:     "/login",
			data:     noPasswordUser,
			response: "Password is required",
			status:   http.StatusBadRequest,
		},
		{
			desc:     "Login failure invalid password",
			method:   "POST",
			path:     "/login",
			data:     adminUserWrongPass,
			response: "error: The username or password is incorrect. Try again.",
			status:   http.StatusUnauthorized,
		},
		{
			desc:     "Login failure invalid username",
			method:   "POST",
			path:     "/login",
			data:     notExistingUser,
			response: "error: The username or password is incorrect. Try again.",
			status:   http.StatusUnauthorized,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			req, err := http.NewRequest(tC.method, ts.URL+tC.path, strings.NewReader(tC.data))
			if err != nil {
				t.Fatal(err)
			}
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			resBody, err := io.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			if res.StatusCode != tC.status || !strings.Contains(string(resBody), tC.response) {
				t.Errorf("expected response did not match.\nExpected vs Received status code: %d vs %d\nExpected vs Received body: \n%s\nvs\n%s\n", tC.status, res.StatusCode, tC.response, string(resBody))
			}
			if tC.desc == "Login success" && res.StatusCode == http.StatusOK {
				token, parseErr := jwt.Parse(string(resBody), func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
					}
					return []byte(env.JWTSecret), nil
				})
				if parseErr != nil {
					t.Errorf("Error parsing JWT: %v", parseErr)
					return
				}

				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					if claims["username"] != "testadmin" {
						t.Errorf("Username found in JWT does not match expected value.")
					} else if int(claims["permissions"].(float64)) != 1 {
						t.Errorf("Permissions found in JWT does not match expected value.")
					}
				} else {
					t.Errorf("Invalid JWT token or JWT claims are not readable")
				}
			}
		})
	}
}

func TestAuthorization(t *testing.T) {
	testdb, err := certdb.NewCertificateRequestsRepository(":memory:", "CertificateRequests")
	if err != nil {
		log.Fatalf("couldn't create test sqlite db: %s", err)
	}
	env := &server.Environment{}
	env.DB = testdb
	env.JWTSecret = []byte("secret")
	ts := httptest.NewTLSServer(server.NewGoCertRouter(env))
	defer ts.Close()

	client := ts.Client()
	var adminToken string
	var nonAdminToken string
	t.Run("prepare user accounts and tokens", prepareUserAccounts(ts.URL, client, &adminToken, &nonAdminToken))

	testCases := []struct {
		desc     string
		method   string
		path     string
		data     string
		auth     string
		response string
		status   int
	}{
		{
			desc:     "metrics reachable without auth",
			method:   "GET",
			path:     "/metrics",
			data:     "",
			auth:     "",
			response: "# HELP certificate_requests Total number of certificate requests",
			status:   http.StatusOK,
		},
		{
			desc:     "status reachable without auth",
			method:   "GET",
			path:     "/status",
			data:     "",
			auth:     "",
			response: "",
			status:   http.StatusOK,
		},
		{
			desc:     "missing endpoints produce 404",
			method:   "GET",
			path:     "/this/path/does/not/exist",
			data:     "",
			auth:     nonAdminToken,
			response: "",
			status:   http.StatusNotFound,
		},
		{
			desc:     "nonadmin can't see accounts",
			method:   "GET",
			path:     "/api/v1/accounts",
			data:     "",
			auth:     nonAdminToken,
			response: "",
			status:   http.StatusForbidden,
		},
		{
			desc:     "admin can see accounts",
			method:   "GET",
			path:     "/api/v1/accounts",
			data:     "",
			auth:     adminToken,
			response: `[{"id":1,"username":"testadmin","permissions":1},{"id":2,"username":"testuser","permissions":0}]`,
			status:   http.StatusOK,
		},
		{
			desc:     "nonadmin can't delete admin account",
			method:   "DELETE",
			path:     "/api/v1/accounts/1",
			data:     "",
			auth:     nonAdminToken,
			response: "",
			status:   http.StatusForbidden,
		},
		{
			desc:     "user can't change admin password",
			method:   "POST",
			path:     "/api/v1/accounts/1/change_password",
			data:     `{"password":"Pwnd123!"}`,
			auth:     nonAdminToken,
			response: "",
			status:   http.StatusForbidden,
		},
		{
			desc:     "user can change self password with /me",
			method:   "POST",
			path:     "/api/v1/accounts/me/change_password",
			data:     `{"password":"BetterPW1!"}`,
			auth:     nonAdminToken,
			response: "",
			status:   http.StatusOK,
		},
		{
			desc:     "user can login with new password",
			method:   "POST",
			path:     "/login",
			data:     `{"username":"testuser","password":"BetterPW1!"}`,
			auth:     nonAdminToken,
			response: "",
			status:   http.StatusOK,
		},
		{
			desc:     "admin can't delete itself",
			method:   "DELETE",
			path:     "/api/v1/accounts/1",
			data:     "",
			auth:     adminToken,
			response: "error: deleting an Admin account is not allowed.",
			status:   http.StatusBadRequest,
		},
		{
			desc:     "admin can delete nonuser",
			method:   "DELETE",
			path:     "/api/v1/accounts/2",
			data:     "",
			auth:     adminToken,
			response: "1",
			status:   http.StatusAccepted,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			req, err := http.NewRequest(tC.method, ts.URL+tC.path, strings.NewReader(tC.data))
			req.Header.Add("Authorization", "Bearer "+tC.auth)
			if err != nil {
				t.Fatal(err)
			}
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			resBody, err := io.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			if res.StatusCode != tC.status || !strings.Contains(string(resBody), tC.response) {
				t.Errorf("expected response did not match.\nExpected vs Received status code: %d vs %d\nExpected vs Received body: \n%s\nvs\n%s\n", tC.status, res.StatusCode, tC.response, string(resBody))
			}
			if tC.desc == "Create no password user success" {
				match, _ := regexp.MatchString(`"password":"[!-~]{16}"`, string(resBody))
				if !match {
					t.Errorf("password does not match expected format or length: got %s", string(resBody))
				}
			}
		})
	}
}

func prepareUserAccounts(url string, client *http.Client, adminToken, nonAdminToken *string) func(*testing.T) {
	return func(t *testing.T) {
		req, err := http.NewRequest("POST", url+"/api/v1/accounts", strings.NewReader(adminUser))
		if err != nil {
			t.Fatal(err)
		}
		res, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		_, err = io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != http.StatusCreated {
			t.Fatalf("creating the first request should succeed when unauthorized. status code received: %d", res.StatusCode)
		}
		req, err = http.NewRequest("POST", url+"/api/v1/accounts", strings.NewReader(validUser))
		if err != nil {
			t.Fatal(err)
		}
		res, err = client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		_, err = io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != http.StatusUnauthorized {
			t.Fatalf("the second request should have been rejected. status code received: %d", res.StatusCode)
		}
		req, err = http.NewRequest("POST", url+"/login", strings.NewReader(adminUser))
		if err != nil {
			t.Fatal(err)
		}
		res, err = client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		resBody, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("the admin login request should have succeeded. status code received: %d", res.StatusCode)
		}
		*adminToken = string(resBody)
		req, err = http.NewRequest("POST", url+"/api/v1/accounts", strings.NewReader(validUser))
		req.Header.Set("Authorization", "Bearer "+*adminToken)
		if err != nil {
			t.Fatal(err)
		}
		res, err = client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		_, err = io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != http.StatusCreated {
			t.Fatalf("creating the second request should have succeeded when given the admin auth header. status code received: %d", res.StatusCode)
		}
		req, err = http.NewRequest("POST", url+"/login", strings.NewReader(validUser))
		if err != nil {
			t.Fatal(err)
		}
		res, err = client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		resBody, err = io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != http.StatusOK {
			t.Errorf("the admin login request should have succeeded. status code received: %d", res.StatusCode)
		}
		*nonAdminToken = string(resBody)
	}
}
