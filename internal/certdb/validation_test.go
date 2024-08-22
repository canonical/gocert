package certdb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/canonical/gocert/internal/certdb"
)

const (
	AppleCSR string = `-----BEGIN CERTIFICATE REQUEST-----
MIICsTCCAZkCAQAwbDELMAkGA1UEBhMCQ0ExFDASBgNVBAgMC05vdmEgU2NvdGlh
MRAwDgYDVQQHDAdIYWxpZmF4MSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0
eSBMdGQxEjAQBgNVBAMMCWFwcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEP
ADCCAQoCggEBAOhDSpNbeFiXMQzQcobExHqYMEGzqpX8N9+AR6/HPZWBybgx1hr3
ejqsKornzpVph/dO9UC7O9aBlG071O9VQGHt3OU3rkZIk2009vYwLuSrAlJtnUne
p7KKn2lZGvh7jVyZE5RkS0X27vlT0soANsmcVq/82VneHrF/nbDcK6DOjQpS5o5l
EiNk2CIpYGUkw3WnQF4pBk8t4bNOl3nfpaAOfnmNuBX3mWyfPnaKMCENMpDqL9FR
V/O5bIPLmyH30OHUEJUkWOmFt9GFi+QfMoM0fR34KmRbDz79hZZb/yVPZZJl7l6i
FWXkNR3gxdEnwCZkTgWk5OqS9dCJOtsDE8ECAwEAAaAAMA0GCSqGSIb3DQEBCwUA
A4IBAQCqBX5WaNv/HjkzAyNXYuCToCb8GjmiMqL54t+1nEI1QTm6axQXivEbQT3x
GIh7uQYC06wHE23K6Znc1/G+o3y6lID07rvhBNal1qoXUiq6CsAqk+DXYdd8MEh5
joerEedFqcW+WTUDcqddfIyDAGPqrM9j6/E+aFYyZjJ/xRuMf1zlWMljRiwj1NI9
NxqjsYYQ3zxfUjv8gxXm0hN8Up1O9saoEF+zbuWNdiUWd6Ih3/3u5VBNSxgVOrDQ
CeXyyzkMx1pWTx0rWa7NSa+DMKVVzv46pck/9kLB4gPL8zqvIOMQsf74N0VcbVfd
9jQR8mPXQYPUERl1ZhNrkzkyA0kd
-----END CERTIFICATE REQUEST-----
`
	BananaCSR string = `-----BEGIN CERTIFICATE REQUEST-----
MIICrjCCAZYCAQAwaTELMAkGA1UEBhMCVFIxDjAMBgNVBAgMBUl6bWlyMRIwEAYD
VQQHDAlOYXJsaWRlcmUxITAfBgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5IEx0
ZDETMBEGA1UEAwwKYmFuYW5hLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC
AQoCggEBAK+vJMxO1GTty09/E4M/RbTCPABleCuYc/uzj72KWaIvoDaanuJ4NBWM
2aUiepxWdMNTR6oe31gLq4agLYT309tXwCeBLQnOxvBFWONmBG1qo0fQkvT5kSoq
AO29D7hkQ0gVwg7EF3qOd0JgbDm/yvexKpYLVvWMQAngHwZRnd5vHGk6M3P7G4oG
mIj/CL2bF6va7GWODYHb+a7jI1nkcsrk+vapc+doVszcoJ+2ryoK6JndOSGjt9SD
uxulWZHQO32XC0btyub63pom4QxRtRXmb1mjM37XEwXJSsQO1HOnmc6ycqUK53p0
jF8Qbs0m8y/p2NHFGTUfiyNYA3EdkjUCAwEAAaAAMA0GCSqGSIb3DQEBCwUAA4IB
AQA+hq8kS2Y1Y6D8qH97Mnnc6Ojm61Q5YJ4MghaTD+XXbueTCx4DfK7ujYzK3IEF
pH1AnSeJCsQeBdjT7p6nv5GcwqWXWztNKn9zibXiASK/yYKwqvQpjSjSeqGEh+Sa
9C9SHeaPhZrJRj0i3NkqmN8moWasF9onW6MNKBX0B+pvBB+igGPcjCIFIFGUUaky
upMXY9IG3LlWvlt+HTfuMZV+zSOZgD9oyqkh5K9XRKNq/mnNz/1llUCBZRmfeRBY
+sJ4M6MJRztiyX4/Fjb8UHQviH931rkiEGtG826IvWIyiRSnAeE8B/VzL0GlT9Zq
ge6lFRxB1FlDuU4Blef8FnOI
-----END CERTIFICATE REQUEST-----`

	StrawberryCSR string = `-----BEGIN CERTIFICATE REQUEST-----
MIICrzCCAZcCAQAwajELMAkGA1UEBhMCSVQxDzANBgNVBAgMBlBhZG92YTEOMAwG
A1UEBwwFUGFkdWExITAfBgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDEX
MBUGA1UEAwwOc3RyYXdiZXJyeS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw
ggEKAoIBAQDXXHpy+3LLRCImyEQitM9eUdgYkexLz2PcAf89tTpkpt3L1woJw0bv
+YR80UcR2Pg+7uUVm4XSKFvcdyWg8yADHIDDZkEmKFEbrOLUsWWTQEsCpFt5MU4u
6YnYXV0YflPXmRsJRd90NOen+wlM2ajK1gGTtLPdJ6axz15LdcT2uXXIvWhncjgL
CvVpd/x44AMxD/BPf/d27VO5hEjxR//DtcOmS/jA+Zf1+dyIAWs2LH+ctsaPLOcg
1rBiRrHtGL8wmPwgwK9b+QLiq9Ik+dx1Jl6BvC36LRk2CxTxfZ6e4UdYVhtnjMW2
VEUAVg9LtowvXTexESUv6Mh4uQF6pW5ZAgMBAAGgADANBgkqhkiG9w0BAQsFAAOC
AQEAW40HaxjVSDNKeWJ8StWGfstdvk3dwqjsfLgmnBBZSLcGppYEnnRlJxhMJ9Ks
x2IYw7wJ55kOJ7V+SunKPPoY+7PwNDV9Llxp58vvE8CFnOc3WcL9pA2V5LbTXwtT
R7jID5GZjOv0bn3x1WXuKVW5tkYdT6sW14rfGut1T+r1kYls+JQ5ap+BzfMtThZz
38PCnEMmSo0/KmgUu5/LakPoy3JPaFB0bCgViZSWlxiSR44YZPsVaRL8E7Zt/qjJ
glRL/48q/tORtxv18/Girl6oiQholkADaH3j2gB3t/fCLp8guAVLWB9DzhwrqWwP
GFl9zB5HDoij2l0kHrb44TuonQ==
-----END CERTIFICATE REQUEST-----
`

	BananaCert string = `-----BEGIN CERTIFICATE-----
MIIEUTCCAjkCFE8lmuBE85/RPw2M17Kzl93O+9IIMA0GCSqGSIb3DQEBCwUAMGEx
CzAJBgNVBAYTAlRSMQ4wDAYDVQQIDAVJem1pcjESMBAGA1UEBwwJTmFybGlkZXJl
MSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQxCzAJBgNVBAMMAm1l
MB4XDTI0MDYyODA4NDIyMFoXDTI1MDYyODA4NDIyMFowaTELMAkGA1UEBhMCVFIx
DjAMBgNVBAgMBUl6bWlyMRIwEAYDVQQHDAlOYXJsaWRlcmUxITAfBgNVBAoMGElu
dGVybmV0IFdpZGdpdHMgUHR5IEx0ZDETMBEGA1UEAwwKYmFuYW5hLmNvbTCCASIw
DQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAK+vJMxO1GTty09/E4M/RbTCPABl
eCuYc/uzj72KWaIvoDaanuJ4NBWM2aUiepxWdMNTR6oe31gLq4agLYT309tXwCeB
LQnOxvBFWONmBG1qo0fQkvT5kSoqAO29D7hkQ0gVwg7EF3qOd0JgbDm/yvexKpYL
VvWMQAngHwZRnd5vHGk6M3P7G4oGmIj/CL2bF6va7GWODYHb+a7jI1nkcsrk+vap
c+doVszcoJ+2ryoK6JndOSGjt9SDuxulWZHQO32XC0btyub63pom4QxRtRXmb1mj
M37XEwXJSsQO1HOnmc6ycqUK53p0jF8Qbs0m8y/p2NHFGTUfiyNYA3EdkjUCAwEA
ATANBgkqhkiG9w0BAQsFAAOCAgEAVZJZD0/ojZSOVIesZvrjLG0agSp0tsXY+hEt
I/knpYLvRcAd8b3Jx9gk+ug+FwDQ4IBIkTX18qhK2fgVUuMR/ubfpQeCMbp64N3Q
kmN/E1eu0bl6hhHAL7jEbi0DE3vAN9huQxAIu5pCyLvZIrPJtvuyj2jOpJBZwGoP
539lfEM++XALzI4qKQ6Z0a0rJZ4HoruKiYwEFZ7VkmRLD0uef6NMZRqa/Vx+o0uT
1TjH4AeDDmJmP/aHlHbpXkHQ9h9rfTa6Qbypo+T9pGDhd02O1tEqrHfiQyNWJxb0
rbR+owT32iCfayzKKqhmAYSF2d9XKWEhulgxWDaXgvUbq4Y+fgfU2qMVz5uusTDh
a9Mp9dsYWySWEUcEa4v2w6FfaaVXE1S9ubm+HoIVtotuutL5fn86q19pAAePYjLQ
ybiETp5LU3chuYmMlCiDRNGHYhN5nvGcttqRdWIBe454RRPNo4iGVl13l6aG8rmI
xDfk5lIwObalbELv+mEIGI1j/j4//nJFXByxlLHm5/BF8rmvHDj1aPtPRw9DLgSX
ejhjjec1xnkBR+JF0g474hLdPjCnA0aqLQInZbjJJm5iXzyXBg1cy7KvIBy3ZkrR
Pp7ObjaWxjCT3O6nEH3w6Ozsyg2cHXQIdVXLvNnV1bxUbPnfhQosKGKgU6s+lcLM
SRhHB2k=
-----END CERTIFICATE-----
`

	IssuerCert string = `-----BEGIN CERTIFICATE-----
MIIFozCCA4ugAwIBAgIUDjtO3bEluUX3tzvrckATlycRVfwwDQYJKoZIhvcNAQEL
BQAwYTELMAkGA1UEBhMCVFIxDjAMBgNVBAgMBUl6bWlyMRIwEAYDVQQHDAlOYXJs
aWRlcmUxITAfBgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDELMAkGA1UE
AwwCbWUwHhcNMjQwNjI4MDYwNTQ5WhcNMzQwNjI2MDYwNTQ5WjBhMQswCQYDVQQG
EwJUUjEOMAwGA1UECAwFSXptaXIxEjAQBgNVBAcMCU5hcmxpZGVyZTEhMB8GA1UE
CgwYSW50ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMQswCQYDVQQDDAJtZTCCAiIwDQYJ
KoZIhvcNAQEBBQADggIPADCCAgoCggIBAJU+5YaFlpn+bWvVri5L6EkmbAPuavsI
/KXY7ufRmc5qb08o1na9lLJ/7TuMD4K36Idnq20n1JohSlrdymBpNZ8O3m5fYYtk
hx5WADlBZsKnC5aZJIChEb4bYcOFLP+d3PooVsAKBxW0Q6TECviQcK7GxaxEZw0L
7FRhX2c9+CxbvRGP6OGVggXZxwkZik/JJ9aym+fltt9QvlxQVBq/GlFYZYC+H8jV
Z6RnUjugnWcTm9PAsQ6+EHEevAW+dWaDP+gr9AgKKz1EXbc1mVKAVOLHjb+Ue7RC
vFoar/YxYIszD58dOSB/GuAxn+JAjWbnOu7jeX3XeWlKOagUJF9L9TgMIUWdiuJG
8Uu/kK2MjyRFdT8opnPFAXrK7vSuMBzhRtswAlWc8xoZWeSQF+NpjU+swbg8ySYT
LfZxVB+s/ftxnGU3RM/RWdbZhb0DAuIBsFAGCbnj+Q61/cK4i58JVjUqzLk+XOwR
55LAyS0Y5pj9jDc5mqvS0z7ot7s2OBM1+o8e3KJgdMSXorYkv3toHMGEIUmPQZCX
JtRCjFNgnoWeLDc+oLiN6BlPx7bS4MDN9tMPCJwF6vnxFzLAzdRqY3D7uRS3chsx
7ClMR9MDsSxplC7tptXgv8UTzh1XZjWGCeZq0Gbe927Hmwy2q8k/BFwnR4PIVSiE
7YAZPb0CPmrfAgMBAAGjUzBRMB0GA1UdDgQWBBRgLXukRHTovOG6g9Z5eCaeh6Sx
aTAfBgNVHSMEGDAWgBRgLXukRHTovOG6g9Z5eCaeh6SxaTAPBgNVHRMBAf8EBTAD
AQH/MA0GCSqGSIb3DQEBCwUAA4ICAQA9TpgTrGmnyxKB2ne76LNQadiijVPpS6/U
OPFAX4EPJ0V5DhDreJjsZJC6Is2Q9+qsPpn/nlW7bvZUVHGodUKcE+TQWFiMtLvu
8ifzk8x1R46aqhTyxb7WBBFfvbvdmlEENKTmTS6A/C3nYgmkfk5N7x84iTowmsVl
Yzz9iRzxkqQ+mU3L2/Sp5nXPYWfzV9WXIJdxWcot7f4CJ79eVFu4D9hYfzcPQ9P9
0qCBRbH/01D2E/3uTHhZPPmK2Tp1ao5SuGLppjMPX8VWVL5CMTXOj+1LF0nJJc/J
9MrqXwtlLyKGP6HX8qALbaXwcv7db6bF+aEsgWmIEB+0ecGk9IR3XQn7I379CO3v
J2oUCZ++lV9e2tcRehUprE1v8i+DFhPtS1iNjrO7KnDYkXimR5zI+3sGFI9/9wY0
4PAV/roZFiEJHe5kA49vwIihJaDgy/SPIYgG/vhdj+WeIbi1ilEi12ou7VF0tyiE
j3eXaMAL8EAKxCUZbXcuwmK9qistAYXBFFEK9M08FwLH8HM4LoPjshMg3II9Ncs8
p3to8U99/ZeFbJRzEUF9poZ7VwxBEcgfWD1RV0+gNLC3Au2yuc4C3anknOv7Db/r
jdzVA8yTI8cZ/RtRohp5H/s+j2tcdfB3Zt+wfS4nLxqN/kf7qv2VSdPbXyTyz/ft
btZkbfdL5A==
-----END CERTIFICATE-----
`
	WrongPKIssuerCert string = `-----BEGIN CERTIFICATE-----
MIIDozCCAougAwIBAgIUE8WQeUw8YMVJlt37CjOBXJ+R7wQwDQYJKoZIhvcNAQEL
BQAwYTELMAkGA1UEBhMCVFIxDjAMBgNVBAgMBUl6bWlyMRIwEAYDVQQHDAlOYXJs
aWRlcmUxITAfBgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDELMAkGA1UE
AwwCbWUwHhcNMjQwODIyMDkyMDI2WhcNMzQwODIwMDkyMDI2WjBhMQswCQYDVQQG
EwJUUjEOMAwGA1UECAwFSXptaXIxEjAQBgNVBAcMCU5hcmxpZGVyZTEhMB8GA1UE
CgwYSW50ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMQswCQYDVQQDDAJtZTCCASIwDQYJ
KoZIhvcNAQEBBQADggEPADCCAQoCggEBAI5KK5H6JN8n3ZF8gp4DRLtIuE16MMZu
H6Br59in71lIE64TMiL6ScqVu7x+etUlvTDcLBX5yYpQ5gZLwB9MyqTRqctZDtP8
82Pa/XIkknFPhcfYN/njINKp2mm1P5zsSm8bznhiCnrfxsYZ13lrJBPjsceRgnD4
Z3207STUO9XIKb1qDUo2tRS1t49g4XiYhEaeATftXladO8AjM99ERXF41MRl8TOm
tRvhl0QrJnEn7CTOhbgN9HYdE9Bu6nOVWLM0zjyeqFJGFlWMTCRYwxYx1/jr6vwl
sF8N+8mkuMpQg13oQdFNpCK9YyoWRoC9zKJbh727VSPzqpyR2I1upg8CAwEAAaNT
MFEwHQYDVR0OBBYEFMi9/T/yZ5n6/PFIketk1fsDpXLgMB8GA1UdIwQYMBaAFMi9
/T/yZ5n6/PFIketk1fsDpXLgMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQEL
BQADggEBAB+pIJCiRzaWAht4pBMmbrDaymIhaHeBsAkFmleAZo0cKixAZp4cP2J6
zIN4pEHchsX259wRiAoy0oQ4D1B2fUE+4FYKdIUMQqXh3h8eXaOJAea/OOLHU+9q
nJoQ/4LqsLpwEGB0ZUJN8RO+LML3U1FyY+5Y7tNj5JlpWMtBebAEdhDS91fVdAp+
jALl5X1Wbx/dtBQnubm1YolBVYXnI2zYywa8IgpnguCu9NIp3uqSVf0xcBEnNIny
W5/mfOoXTnuKZKTEvButfrlkLsABQvVepitmZGv+q/f4crCkhms8B23WMRLdteiK
BqHOQR7Y7LSxxC+bAa1QdhgumR3PL8I=
-----END CERTIFICATE-----`
	WrongSubjectIssuerCert string = `-----BEGIN CERTIFICATE-----
MIIFqTCCA5GgAwIBAgIUWJY4vKnl3+kQ487QtMfzLDBTnAowDQYJKoZIhvcNAQEL
BQAwZDELMAkGA1UEBhMCVFIxDjAMBgNVBAgMBUl6bWlyMRIwEAYDVQQHDAlOYXJs
aWRlcmUxITAfBgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDEOMAwGA1UE
AwwFbm90bWUwHhcNMjQwODIyMDkxOTQ4WhcNMzQwODIwMDkxOTQ4WjBkMQswCQYD
VQQGEwJUUjEOMAwGA1UECAwFSXptaXIxEjAQBgNVBAcMCU5hcmxpZGVyZTEhMB8G
A1UECgwYSW50ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMQ4wDAYDVQQDDAVub3RtZTCC
AiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAJU+5YaFlpn+bWvVri5L6Ekm
bAPuavsI/KXY7ufRmc5qb08o1na9lLJ/7TuMD4K36Idnq20n1JohSlrdymBpNZ8O
3m5fYYtkhx5WADlBZsKnC5aZJIChEb4bYcOFLP+d3PooVsAKBxW0Q6TECviQcK7G
xaxEZw0L7FRhX2c9+CxbvRGP6OGVggXZxwkZik/JJ9aym+fltt9QvlxQVBq/GlFY
ZYC+H8jVZ6RnUjugnWcTm9PAsQ6+EHEevAW+dWaDP+gr9AgKKz1EXbc1mVKAVOLH
jb+Ue7RCvFoar/YxYIszD58dOSB/GuAxn+JAjWbnOu7jeX3XeWlKOagUJF9L9TgM
IUWdiuJG8Uu/kK2MjyRFdT8opnPFAXrK7vSuMBzhRtswAlWc8xoZWeSQF+NpjU+s
wbg8ySYTLfZxVB+s/ftxnGU3RM/RWdbZhb0DAuIBsFAGCbnj+Q61/cK4i58JVjUq
zLk+XOwR55LAyS0Y5pj9jDc5mqvS0z7ot7s2OBM1+o8e3KJgdMSXorYkv3toHMGE
IUmPQZCXJtRCjFNgnoWeLDc+oLiN6BlPx7bS4MDN9tMPCJwF6vnxFzLAzdRqY3D7
uRS3chsx7ClMR9MDsSxplC7tptXgv8UTzh1XZjWGCeZq0Gbe927Hmwy2q8k/BFwn
R4PIVSiE7YAZPb0CPmrfAgMBAAGjUzBRMB0GA1UdDgQWBBRgLXukRHTovOG6g9Z5
eCaeh6SxaTAfBgNVHSMEGDAWgBRgLXukRHTovOG6g9Z5eCaeh6SxaTAPBgNVHRMB
Af8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4ICAQAwNq0z26m13RBvZX/uOR5tIQ6j
l/JpSMhocr6GUTKx1NEmyaO9UEAdwHi7nFGocCbCeMNPBxpaJGkSTxe5HefhDJOI
QcnOo1yY9q5HXsp2SPvXjkZ2Palg1rV/u8BChVvULDDT+JtABJlll+cfggh1pkZv
Z3V7Zh7u7gWbnsSnM0X3zVpxGf/cqZNEoHesAaWJA4yYIH2wr5TwqksXGPFE/g/Z
fhUDeI7OP8kM/A8HnCXdxUok2Zf/wyuoPvrFUaPrcYkZK3omT6H24VdyejuBe2k5
+e0ij3nU8DxKEbKn6XaJFhBzAmP1APi8fLIwO6gig/XUWrfKrqO0ax4Vgl4r88Ht
y4hiHmP9kgWjYqUijLpK5ap5607tfbtZ0QIS54HAPAjE77ZdsEGfkAZPmyCTPg41
Q+YWZJS8HogVTZKY267x7u4lQ68jSVBxpeRHGYzd2HWxWGKVQq8pEa2bob9zby/N
QNRikyGkbp7ep5HgBrZeJJJ5zdaqNzVmXY0JIfhkUypSiCe5X1WgZ9GVCC9wi72D
y6MHDTAyVHrSouCqfh9XD6RDN58d+u9kLEg0WJD55wH4E4z+ZZhEMicCWfT/rn+b
b3dRTVslxdJ0dOApn/6zwfRMXgI7j2yRSkA7F39ekwlPhJy2bGrEDgTlDK33AwPU
wM1PZYERQJNOGMAI5Q==
-----END CERTIFICATE-----`
)

func TestCSRValidationSuccess(t *testing.T) {
	cases := []string{AppleCSR, BananaCSR, StrawberryCSR}

	for i, c := range cases {
		t.Run(fmt.Sprintf("ValidCSR%d", i), func(t *testing.T) {
			if err := certdb.ValidateCertificateRequest(c); err != nil {
				t.Errorf("Couldn't verify valid CSR: %s", err)
			}
		})
	}
}

func TestCSRValidationFail(t *testing.T) {
	var wrongString = "this is a real csr!!!"
	var wrongStringErr = "PEM Certificate Request string not found or malformed"
	var ValidCSRWithoutWhitespace = strings.ReplaceAll(AppleCSR, "\n", "")
	var ValidCSRWithoutWhitespaceErr = "PEM Certificate Request string not found or malformed"
	var wrongPemType = strings.ReplaceAll(AppleCSR, "CERTIFICATE REQUEST", "SOME RANDOM PEM TYPE")
	var wrongPemTypeErr = "given PEM string not a certificate request"
	var InvalidCSR = strings.ReplaceAll(AppleCSR, "s", "p")
	var InvalidCSRErr = "asn1: syntax error: data truncated"

	cases := []struct {
		input       string
		expectedErr string
	}{
		{
			input:       wrongString,
			expectedErr: wrongStringErr,
		},
		{
			input:       ValidCSRWithoutWhitespace,
			expectedErr: ValidCSRWithoutWhitespaceErr,
		},
		{
			input:       wrongPemType,
			expectedErr: wrongPemTypeErr,
		},
		{
			input:       InvalidCSR,
			expectedErr: InvalidCSRErr,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("InvalidCSR%d", i), func(t *testing.T) {
			err := certdb.ValidateCertificateRequest(c.input)
			if err == nil {
				t.Errorf("No error received. Expected: %s", c.expectedErr)
				return
			}
			if err.Error() != c.expectedErr {
				t.Errorf("Expected error not found:\nReceived: %s\nExpected: %s", err, c.expectedErr)
			}
		})
	}
}

func TestCertValidationSuccess(t *testing.T) {
	cases := []string{fmt.Sprintf("%s\n%s", BananaCert, IssuerCert)}

	for i, c := range cases {
		t.Run(fmt.Sprintf("ValidCert%d", i), func(t *testing.T) {
			if err := certdb.ValidateCertificate(c); err != nil {
				t.Errorf("Couldn't verify valid Cert: %s", err)
			}
		})
	}
}

func TestCertValidationFail(t *testing.T) {
	var wrongCertString = "this is a real cert!!!"
	var wrongCertStringErr = "less than 2 certificate PEM strings were found"
	var wrongPemType = strings.ReplaceAll(BananaCert, "CERTIFICATE", "SOME RANDOM PEM TYPE")
	var wrongPemTypeErr = "a given PEM string was not a certificate"
	var InvalidCert = strings.ReplaceAll(BananaCert, "M", "i")
	var InvalidCertErr = "x509: malformed certificate"
	var singleCert = BananaCert
	var singleCertErr = "less than 2 certificate PEM strings were found"
	var issuerCertPKDoesNotMatch = fmt.Sprintf("%s\n%s", BananaCert, WrongPKIssuerCert)
	var issuerCertPKDoesNotMatchErr = "invalid certificate chain: certificate 0, certificate 1: keys do not match"
	var issuerCertSubjectDoesNotMatch = fmt.Sprintf("%s\n%s", BananaCert, WrongSubjectIssuerCert)
	var issuerCertSubjectDoesNotMatchErr = "invalid certificate chain: certificate 0, certificate 1: subjects do not match"

	cases := []struct {
		inputCert   string
		expectedErr string
	}{
		{
			inputCert:   wrongCertString,
			expectedErr: wrongCertStringErr,
		},
		{
			inputCert:   wrongPemType,
			expectedErr: wrongPemTypeErr,
		},
		{
			inputCert:   InvalidCert,
			expectedErr: InvalidCertErr,
		},
		{
			inputCert:   singleCert,
			expectedErr: singleCertErr,
		},
		{
			inputCert:   issuerCertPKDoesNotMatch,
			expectedErr: issuerCertPKDoesNotMatchErr,
		},
		{
			inputCert:   issuerCertSubjectDoesNotMatch,
			expectedErr: issuerCertSubjectDoesNotMatchErr,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("InvalidCert%d", i), func(t *testing.T) {
			err := certdb.ValidateCertificate(c.inputCert)
			if err == nil {
				t.Errorf("No error received. Expected: %s", c.expectedErr)
				return
			}
			if !strings.HasPrefix(err.Error(), c.expectedErr) {
				t.Errorf("Expected error not found:\nReceived: %s\n Expected: %s", err, c.expectedErr)
			}
		})
	}
}

func TestCertificateMatchesCSRSuccess(t *testing.T) {
	cases := []struct {
		inputCSR  string
		inputCert string
	}{
		{
			inputCSR:  BananaCSR,
			inputCert: fmt.Sprintf("%s\n%s", BananaCert, IssuerCert),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("InvalidCert%d", i), func(t *testing.T) {
			err := certdb.CertificateMatchesCSR(c.inputCert, c.inputCSR)
			if err != nil {
				t.Errorf("Certificate did not match when it should have")
			}
		})
	}
}

func TestCertificateMatchesCSRFail(t *testing.T) {
	var certificateDoesNotMatchErr = "certificate does not match CSR"

	cases := []struct {
		inputCSR    string
		inputCert   string
		expectedErr string
	}{
		{
			inputCSR:    AppleCSR,
			inputCert:   fmt.Sprintf("%s\n%s", BananaCert, IssuerCert),
			expectedErr: certificateDoesNotMatchErr,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("InvalidCert%d", i), func(t *testing.T) {
			err := certdb.CertificateMatchesCSR(c.inputCert, c.inputCSR)
			if err == nil {
				t.Errorf("No error received. Expected: %s", c.expectedErr)
				return
			}
			if err.Error() != c.expectedErr {
				t.Errorf("Expected error not found:\nReceived: %s\n Expected: %s", err, c.expectedErr)
			}
		})
	}
}
