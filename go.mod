module v

go 1.20

require (
	cloud.google.com/go/speech v1.21.0
	cloud.google.com/go/texttospeech v1.7.4
	github.com/gofiber/contrib/jwt v1.0.7
	github.com/gofiber/fiber/v2 v2.50.0
	github.com/golang-jwt/jwt/v5 v5.1.0
	github.com/google/uuid v1.4.0
	github.com/lithammer/shortuuid/v4 v4.0.0
	github.com/livekit/protocol v1.9.0
	github.com/pion/rtp v1.8.3
	github.com/pion/webrtc/v3 v3.2.22
	github.com/sashabaranov/go-openai v1.17.5
	github.com/stretchr/testify v1.8.4
	github.com/urfave/cli/v2 v2.25.7
	golang.org/x/exp v0.0.0-20231110203233-9a3e6036ecaa
	google.golang.org/api v0.150.0
	google.golang.org/grpc v1.59.0
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/driver/sqlite v1.5.4
	v/protocol v0.0.0-00010101000000-000000000000
)

replace v/protocol => ./protocol // When developin, this replace uses the local package instead of a pushed version.

require (
	cloud.google.com/go v0.110.10 // indirect
	cloud.google.com/go/compute v1.23.3 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/longrunning v0.5.4 // indirect
	github.com/MicahParks/keyfunc/v2 v2.1.0 // indirect
	github.com/andybalholm/brotli v1.0.6 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bep/debounce v1.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/eapache/channels v1.1.0 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/frostbyte73/core v0.0.9 // indirect
	github.com/gammazero/deque v0.2.1 // indirect
	github.com/go-jose/go-jose/v3 v3.0.1 // indirect
	github.com/go-logr/logr v1.3.0 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gofiber/adaptor/v2 v2.2.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.0 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.5 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.0 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jxskiss/base62 v1.1.0 // indirect
	github.com/klauspost/compress v1.17.2 // indirect
	github.com/klauspost/cpuid/v2 v2.2.6 // indirect
	github.com/livekit/mageutil v0.0.0-20230125210925-54e8a70427c1 // indirect
	github.com/livekit/mediatransportutil v0.0.0-20231017082622-43f077b4e60e // indirect
	github.com/livekit/psrpc v0.5.0 // indirect
	github.com/mackerelio/go-osstat v0.2.4 // indirect
	github.com/magefile/mage v1.15.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/mattn/go-sqlite3 v1.14.18 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/matttproud/golang_protobuf_extensions/v2 v2.0.0 // indirect
	github.com/nats-io/nats.go v1.31.0 // indirect
	github.com/nats-io/nkeys v0.4.6 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/philhofer/fwd v1.1.2 // indirect
	github.com/pion/datachannel v1.5.5 // indirect
	github.com/pion/dtls/v2 v2.2.7 // indirect
	github.com/pion/ice/v2 v2.3.11 // indirect
	github.com/pion/interceptor v0.1.25 // indirect
	github.com/pion/logging v0.2.2 // indirect
	github.com/pion/mdns v0.0.9 // indirect
	github.com/pion/randutil v0.1.0 // indirect
	github.com/pion/rtcp v1.2.12 // indirect
	github.com/pion/sctp v1.8.9 // indirect
	github.com/pion/sdp/v3 v3.0.6 // indirect
	github.com/pion/srtp/v2 v2.0.18 // indirect
	github.com/pion/stun v0.6.1 // indirect
	github.com/pion/transport/v2 v2.2.4 // indirect
	github.com/pion/turn/v2 v2.1.4 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.17.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.45.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/thoas/go-funk v0.9.3 // indirect
	github.com/tinylib/msgp v1.1.9 // indirect
	github.com/twitchtv/twirp v8.1.3+incompatible // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.50.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/crypto v0.15.0 // indirect
	golang.org/x/net v0.18.0 // indirect
	golang.org/x/oauth2 v0.14.0 // indirect
	golang.org/x/sync v0.5.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.4.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto v0.0.0-20231106174013-bbf56f31fb17 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20231106174013-bbf56f31fb17 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231106174013-bbf56f31fb17 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)

require (
	github.com/ansrivas/fiberprometheus/v2 v2.6.1
	github.com/cpuguy83/go-md2man/v2 v2.0.3 // indirect
	github.com/goccy/go-json v0.10.2
	github.com/gofiber/template v1.8.1
	github.com/joho/godotenv v1.5.1
	github.com/jordic/lti v0.0.0-20160211051708-2c756eacbab9
	github.com/livekit/server-sdk-go v1.1.1
	github.com/redis/go-redis/v9 v9.3.0
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0
	gorm.io/driver/postgres v1.5.4
	gorm.io/gorm v1.25.5
)
