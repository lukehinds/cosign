// +build !pivkey !cgo

// Copyright 2021 The Sigstore Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pivkey

import (
	"context"
	"crypto"
	"errors"
	"io"

	"github.com/sigstore/sigstore/pkg/signature"
)

func NewPublicKeyProvider(slotName string) (signature.Verifier, error) {
	return nil, errors.New("unimplemented")
}

func NewSigner() (signature.Signer, error) {
	return nil, errors.New("unimplemented")
}

type PIVSigner struct {
	Priv crypto.PrivateKey
	Pub  crypto.PrivateKey
	signature.ECDSAVerifier
}

func (ps *PIVSigner) Sign(ctx context.Context, rawPayload []byte) ([]byte, []byte, error) {
	return nil, nil, errors.New("unimplemented")
}

func (ps *PIVSigner) SignMessage(rawPayload io.Reader, opts ...signature.SignOption) ([]byte, error) {
	return nil, errors.New("unimplemented")
}

func (ps *PIVSigner) PublicKey(opts ...signature.PublicKeyOption) (crypto.PublicKey, error) {
	return nil, errors.New("unimplemented")
}

var _ signature.Signer = &PIVSigner{}

func NewSignerVerifier(slotName string) (signature.SignerVerifier, error) {
	return nil, errors.New("unimplemented")
}
