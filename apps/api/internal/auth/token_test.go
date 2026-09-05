package auth

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"
	"testing"
)

func TestGenerateRefreshCredentialUses32RawURLBytes(t *testing.T) {
	raw := bytes.Repeat([]byte{0xab}, RefreshCredentialBytes)
	encoded, err := GenerateRefreshCredentialFrom(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("GenerateRefreshCredentialFrom() error = %v", err)
	}

	decoded, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("decode generated credential: %v", err)
	}
	if len(decoded) != RefreshCredentialBytes {
		t.Fatalf("decoded credential length = %d, want %d", len(decoded), RefreshCredentialBytes)
	}
	if encoded != base64.RawURLEncoding.EncodeToString(raw) {
		t.Fatalf("credential = %q, want canonical RawURL encoding", encoded)
	}
	if strings.Contains(encoded, "=") {
		t.Fatal("RawURL credential contains padding")
	}
	t.Logf("raw_entropy_bytes=%d transport_encoding=RawURL padding=false", len(decoded))
}

func TestRefreshCredentialRoundTripAndDigest(t *testing.T) {
	raw := bytes.Repeat([]byte{0x10}, RefreshCredentialBytes)
	encoded, err := GenerateRefreshCredentialFrom(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("GenerateRefreshCredentialFrom() error = %v", err)
	}

	decoded, err := DecodeRefreshCredential(encoded)
	if err != nil {
		t.Fatalf("DecodeRefreshCredential() error = %v", err)
	}
	if !bytes.Equal(decoded, raw) {
		t.Fatalf("decoded credential = %x, want %x", decoded, raw)
	}

	want := sha256.Sum256(raw)
	got := HashRefreshCredential(decoded)
	if !bytes.Equal(got, want[:]) {
		t.Fatalf("digest = %x, want %x", got, want)
	}
	if bytes.Equal([]byte(encoded), got) {
		t.Fatal("encoded plaintext credential equals persisted digest")
	}
	t.Logf("round_trip=true digest_bytes=%d plaintext_persisted=false", len(got))
}

func TestDecodeRefreshCredentialRejectsMalformedAndWrongLength(t *testing.T) {
	validRaw := bytes.Repeat([]byte{0x01}, RefreshCredentialBytes)
	valid := base64.RawURLEncoding.EncodeToString(validRaw)
	tests := []struct {
		name  string
		value string
	}{
		{name: "malformed", value: "not a base64 credential"},
		{name: "empty", value: ""},
		{name: "wrong decoded length", value: base64.RawURLEncoding.EncodeToString([]byte{0x01})},
		{name: "padded non canonical", value: valid + "="},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := DecodeRefreshCredential(test.value)
			if !errors.Is(err, ErrInvalidCredential) {
				t.Fatalf("DecodeRefreshCredential() error = %v, want ErrInvalidCredential", err)
			}
		})
	}
}

func TestIndependentRefreshCredentialSourcesProduceDistinctCredentials(t *testing.T) {
	first, err := GenerateRefreshCredentialFrom(bytes.NewReader(bytes.Repeat([]byte{0x01}, RefreshCredentialBytes)))
	if err != nil {
		t.Fatalf("generate first credential: %v", err)
	}
	second, err := GenerateRefreshCredentialFrom(bytes.NewReader(bytes.Repeat([]byte{0x02}, RefreshCredentialBytes)))
	if err != nil {
		t.Fatalf("generate second credential: %v", err)
	}
	if first == second {
		t.Fatal("independent credential generations returned the same credential")
	}
	t.Log("independent_generations_distinct=true")
}

func TestGenerateRefreshCredentialPropagatesShortSource(t *testing.T) {
	_, err := GenerateRefreshCredentialFrom(bytes.NewReader([]byte{0x01}))
	if err == nil {
		t.Fatal("GenerateRefreshCredentialFrom() error = nil, want short source error")
	}
}

func TestAuthenticationErrorKeepsInternalReasonWithGenericMessage(t *testing.T) {
	err := newAuthenticationError(ErrRefreshReuse)
	if err.Error() != ErrAuthenticationFailed.Error() {
		t.Fatalf("authentication error string = %q, want generic %q", err.Error(), ErrAuthenticationFailed)
	}
	if !errors.Is(err, ErrAuthenticationFailed) || !errors.Is(err, ErrRefreshReuse) {
		t.Fatal("authentication error did not preserve public and internal classifications")
	}
	if strings.Contains(err.Error(), "reuse") {
		t.Fatal("authentication error leaked reuse classification")
	}
}
