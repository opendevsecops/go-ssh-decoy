package main

import (
	"fmt"
	ssh "golang.org/x/crypto/ssh"
	"log"
	"net"
)

const port = 2022

const sshKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEArq1HMJ3p6wjfD5t/M2LWRsZmJ85niO0QmN9myBcsfoVPcUsR
1EmqZIxT/O8GYIeujpsgAKf+BiVxoRaK4kQdCZVcZ8NBhjJ4CsaRwWMSZAQVVTBZ
AxLSkWSlDhbWyP+wahmDc8nywPZbPsaZysExFcyvzX0K/6NrRsKzrytyDiDRGqk8
b6mxMcBCMrI//H16xSgrj09I+mx0Adrf/u4q3a+BbqLrYnMkGu1Ly5TOarvoc7bZ
dzSt3tFS/WmdX7FzeYSevEL4u+fDt+RgqScIht3YaWUwAIhLmOQsN7YYR1rogi+N
yHpVILgsCLaNJXcO05cyM/iO+CqFzPR1hTGdLQIDAQABAoIBAHasWNpp3tuEum9T
GJdjxeptrkfLDkJTtVv3I1A7lkoa7f8tnl5Y8f+/6uvDxNReOjS+pX5so8OvOsTT
mOXimsvBAveoC2NN9Ip1n98AHSWANLIc18yjxBVtiEnLMH98X3GmBc3r3RZGCXXb
3e6HVH7YAnScSJWnhCGd9A/Fd4aqm5aYzbz650UgM/F2HqW7sHOXva77WSawSUK9
xjuO40k7hAom9xs+rL2gpd2U37byKMxttpdGCrwg9yzO735x6zIxajARnhWtXY5m
euRrfEaDHnsX6Dx96Y565gNQSMtl4BsiMfe5dsIRumig6NVwTxncaWffQd9nxrPx
tL2Y0WECgYEA2Q/gg5PK5Cbu70zl8lSIoF9oQJns29siNmMD1DHSeCTqTqoWsWhI
sd2Xwm86OrOUAspIJVhQ7dcA4ZX21En4vLG+DRNulOdHQrnjJvXQMT2WW3KLR2N8
cKK9PaVkv8LmxjgjHPpYGlAHFEUCQiXw+x7zEUIDIfbhrdmlX5k4lpkCgYEAzgLz
bZS0PzMxugM08HZFvLR7uekUh7snH2CliWA3WsCkgqj6fwRWi6yKTalcA0CoYmmv
FN5wYQd8/HEUCszjpe1+wSpc75KIcMedEmE3P6/mthXJluXtV47zROT+t8ullsOg
oTSj9eFkOVWOySo2qpusawVmY6WKYGsU8SPVG7UCgYAbcm2CVcrfBKlL6x5cgSHx
nX7SRGR1/ISb+fM+/rnNZWWXYtyRvE0M6KdK98OWLqT1oVx0FHHPUVOUMuFOQLhK
K/OLNbzS6VfScSzu/UBBKbd8gsRn14WhvIJPbD2MHfoOcITIIkPHt/zdLEi30pJh
Pq2frgg1YEFzOUU3DGniaQKBgDTPTNeqZwpMdVLZv5hkuTvGiHD/7uNcdor0m3q7
z3TULVfROWWWFxl3AX0nDQ9IY+HWdatD2ksFQGT2F80s+K5wUy3xTiGbzp4ajYlI
ooEQ9nN24lZsWos3eeUPTryO18PuIh8w/1bokGiiJhgrWhgiD/DfUX/5z58n1BZ8
uQSBAoGBAMjfWYxtXapBQ8BDgdekxNtBi03Z8zOVX4uXV3/an5Xpe/nqr7MB8i/W
M3lqM1VpSvuvk8l5PIGu6VIm9XIh3NrQ8Ou+pGSGirpcJMlZzHjbKSqEvilud1wB
lUG13d7eoYbr/QBaQoC/Ra8SGQXB81gepjvdgsizH1OP7+76bieB
-----END RSA PRIVATE KEY-----`

func recordCredentials(user string, pass string) {
	fmt.Printf("ssh: %s <---> %s\n", user, pass)
}

func main() {
	config := &ssh.ServerConfig{
		PasswordCallback: func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			recordCredentials(conn.User(), string(password))

			return nil, fmt.Errorf("password rejected for %q", conn.User())
		},
	}

	key, err := ssh.ParsePrivateKey([]byte(sshKey))

	if err != nil {
		log.Fatal("failed to parse ssh key")

		return
	}

	config.AddHostKey(key)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))

	if err != nil {
		log.Fatal(fmt.Sprintf("failed to listen for ssh connections on port %d", port))

		return
	}

	for {
		nconn, err := listener.Accept()

		if err != nil {
			log.Print("failed to accept incoming ssh connection\n")

			continue
		}

		conn, _, _, err := ssh.NewServerConn(nconn, config)

		if err != nil {
			log.Fatal("failed to handshake: ", err)
		}

		conn.Close()
	}
}
