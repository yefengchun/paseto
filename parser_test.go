package paseto

import (
	"crypto"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ed25519"
)

func TestParse(t *testing.T) {
	type tokenInfo struct {
		version Version
		payload interface{}
		footer  interface{}
	}

	symmetricKey, _ := hex.DecodeString("707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f")

	cases := map[string]tokenInfo{
		"v1.local.rElw-WywOuwAqKC9Yao3YokSp7vx0YiUB9hLTnsVOYYTojmVaYumJSQt8aggtCaFKWyaodw5k-CUWhYKATopiabAl4OAmTxHCfm2E4NSPvrmMcmi8n-JcZ93HpcxC6rx_ps22vutv7iP7wf8QcSD1Mwx.Q3VvbiBBbHBpbnVz": {
			version: V1,
			payload: []byte("Love is stronger than hate or fear"),
			footer:  []byte("Cuon Alpinus"),
		},
		"v1.public.TG9yZW0gSXBzdW1684wbBiSvpwhED_5bdFnF2ithKoKDyzEyTOLUlFnz83IibTKCOw3LPOEp8xKM67EYOw1xU6OBBOdLQT-XO5mKMg51JJ4J91IBDwDazDex0D2UQphr7i8gPGP_5FyjlNincP_rToVbYOOzfk9cmnH1-iLmOxxbrsa7-v08Gx12ib-Z-KxKBXBHbxI8uvauVWUVS6A7rl0eAlb6SecSPPQpxQnD1zakA-nGFUbWq5Zx8XqgVZ-VidcGcd7kmhZ-bMy4Z1uGOWmAXHC793v8sbXuRdroZM8kmO0pqQMoE_wmlriIxflFABCa1PPWi5YB87aVF3oIWHYXawZXxRwxevgK.Zm9vdGVy": {
			version: V1,
			payload: []byte("Lorem Ipsum"),
			footer:  []byte("footer"),
		},
		"v2.local.FGVEQLywggpvH0AzKtLXz0QRmGYuC6yvl05z9GIX0cnol6UK94cfV77AXnShlUcNgpDR12FrQiurS8jxBRmvoIKmeMWC5wY9Y6w.Q3VvbiBBbHBpbnVz": {
			version: V2,
			payload: []byte("Love is stronger than hate or fear"),
			footer:  []byte("Cuon Alpinus"),
		},
		"v2.public.RnJhbmsgRGVuaXMgcm9ja3O7MPuu90WKNyvBUUhAGFmi4PiPOr2bN2ytUSU-QWlj8eNefki2MubssfN1b8figynnY0WusRPwIQ-o0HSZOS0F.Q3VvbiBBbHBpbnVz": {
			version: V2,
			payload: []byte("Frank Denis rocks"),
			footer:  []byte("Cuon Alpinus"),
		},
	}

	b, _ := hex.DecodeString("1eb9dbbbbc047c03fd70604e0071f0987e16b28b757225c11f00415d0e20b1a2")
	v2PublicKey := ed25519.PublicKey(b)

	for token, info := range cases {
		var payload []byte
		var footer []byte
		if ver, err := Parse(token, &payload, &footer, symmetricKey, map[Version]crypto.PublicKey{V1: rsaPublicKey, V2: v2PublicKey}); assert.NoError(t, err) {
			assert.Equal(t, info.version, ver)
			assert.Equal(t, info.payload, payload)
			assert.Equal(t, info.footer, footer)
		}
	}
}

func TestParse_Err(t *testing.T) {
	cases := map[string]error{
		"v1.publiceqreqqereqrqerq":              ErrIncorrectTokenFormat,
		"v0.local.rElw-WywOuwAqKC9Yao3YokSp7vx": ErrUnsupportedTokenVersion,
		"v1.private.rElw":                       ErrUnsupportedTokenType,
		"v1.public.rElw":                        ErrPublicKeyNotFound,
	}

	for token, e := range cases {
		_, err := Parse(token, nil, nil, nil, nil)
		assert.Equal(t, e, err)
	}

}

func TestParseFooter(t *testing.T) {
	type tokenInfo struct {
		footer  []byte
		version Version
		err     error
	}

	cases := map[string]tokenInfo{
		"v1.local.rElw-WywOuwAqKC9Yao3YokSp7vx0YiUB9hLTnsVOYYTojmVaYumJSQt8aggtCaFKWyaodw5k-CUWhYKATopiabAl4OAmTxHCfm2E4NSPvrmMcmi8n-JcZ93HpcxC6rx_ps22vutv7iP7wf8QcSD1Mwx.Q3VvbiBBbHBpbnVz": {
			footer:  []byte("Cuon Alpinus"),
			version: V1,
		},
		"v1.local.rElw-WywOuwAqKC9Yao3YokSp7vx0YiUB9hLTnsVOYYTojmVaYumJSQt8aggtCaFKWyaodw5k-CUWhYKATopiabAl4OAmTxHCfm2E4NSPvrmMcmi8n-JcZ93HpcxC6rx_ps22vutv7iP7wf8QcSD1Mwx": {
			version: V1,
		},
		"v1.rElw-WywOuwAqK": {
			err: ErrIncorrectTokenFormat,
		},
	}
	for token, info := range cases {
		var footer []byte
		err := ParseFooter(token, &footer)
		assert.Equal(t, info.err, err)
		assert.Equal(t, info.footer, footer)
	}
}

func TestGetTokenInfo(t *testing.T) {
	cases := map[string]struct {
		version Version
		purpose Purpose
		err     error
	}{
		"v1.local.rElw-WywOuwAqKC9Yao3YokSp7vx0YiUB9hLTnsVOYYTojmVaYumJSQt8aggtCaFKWyaodw5k-CUWhYKATopiabAl4OAmTxHCfm2E4NSPvrmMcmi8n-JcZ93HpcxC6rx_ps22vutv7iP7wf8QcSD1Mwx.Q3VvbiBBbHBpbnVz": {
			version: V1,
			purpose: LOCAL,
		},
		"v2.local.driRNhM20GQPvlWfJCepzh6HdijAq-yNUtKpdy5KXjKfpSKrOlqQvQ": {
			version: V2,
			purpose: LOCAL,
		},
		"v1.public.rElw-WywOuwAqKC9Yao3YokSp7vx0YiUB9hLTnsVOYYTojmVaYumJSQt8aggtCaFKWyaodw5k-CUWhYKATopiabAl4OAmTxHCfm2E4NSPvrmMcmi8n-JcZ93HpcxC6rx_ps22vutv7iP7wf8QcSD1Mwx.Q3VvbiBBbHBpbnVh": {
			version: V1,
			purpose: PUBLIC,
		},
		"v0.public.rElw-WywOuwAqKC9Yao3YokSp7vx0YiUB9hLTnsVOYYTojmVaYumJSQt8aggtCaFKWyaodw5k-CUWhYKATopiabAl4OAmTxHCfm2E4NSPvrmMcmi8n-JcZ93HpcxC6rx_ps22vutv7iP7wf8QcSD1Mwx.Q3VvbiBBbHBpbnVh": {
			err: ErrUnsupportedTokenVersion,
		},
		"v1.private.rElw-WywOuwAqKC9Yao3YokSp7vx0YiUB9hLTnsVOYYTojmVaYumJSQt8aggtCaFKWyaodw5k-CUWhYKATopiabAl4OAmTxHCfm2E4NSPvrmMcmi8n-JcZ93HpcxC6rx_ps22vutv7iP7wf8QcSD1Mwx.Q3VvbiBBbHBpbnVh": {
			err: ErrUnsupportedTokenType,
		},
		"v1.private": {
			err: ErrIncorrectTokenFormat,
		},
	}

	for token, info := range cases {
		version, purpose, err := GetTokenInfo(token)
		assert.Equal(t, info.err, err)
		assert.Equal(t, info.version, version)
		assert.Equal(t, info.purpose, purpose)
	}
}
