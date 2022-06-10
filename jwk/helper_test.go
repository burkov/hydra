/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */

package jwk_test

import (
	"context"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"encoding/pem"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"

	"github.com/ory/hydra/internal"
	"github.com/ory/hydra/jwk"
	"github.com/ory/hydra/x"
	"github.com/ory/x/contextx"

	"gopkg.in/square/go-jose.v2/cryptosigner"

	"gopkg.in/square/go-jose.v2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerFindPublicKey(t *testing.T) {
	t.Run("Test_Helper/Run_FindPublicKey_With_RSA", func(t *testing.T) {
		RSIDKS, err := jwk.GenerateJWK(context.Background(), jose.RS256, "test-id-1", "sig")
		require.NoError(t, err)
		keys, err := jwk.FindPublicKey(RSIDKS)
		require.NoError(t, err)
		assert.Equal(t, keys.KeyID, "test-id-1")
		assert.IsType(t, keys.Key, new(rsa.PublicKey))
	})

	t.Run("Test_Helper/Run_FindPublicKey_With_ECDSA", func(t *testing.T) {
		ECDSAIDKS, err := jwk.GenerateJWK(context.Background(), jose.ES256, "test-id-2", "sig")
		require.NoError(t, err)
		keys, err := jwk.FindPublicKey(ECDSAIDKS)
		require.NoError(t, err)
		assert.Equal(t, keys.KeyID, "test-id-2")
		assert.IsType(t, keys.Key, new(ecdsa.PublicKey))
	})

	t.Run("Test_Helper/Run_FindPublicKey_With_EdDSA", func(t *testing.T) {
		EdDSAIDKS, err := jwk.GenerateJWK(context.Background(), jose.EdDSA, "test-id-3", "sig")
		require.NoError(t, err)
		keys, err := jwk.FindPublicKey(EdDSAIDKS)
		require.NoError(t, err)
		assert.Equal(t, keys.KeyID, "test-id-3")
		assert.IsType(t, keys.Key, ed25519.PublicKey{})
	})

	t.Run("Test_Helper/Run_FindPublicKey_With_KeyNotFound", func(t *testing.T) {
		keySet := &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{}}
		_, err := jwk.FindPublicKey(keySet)
		require.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "key not found"))
	})
}

func TestHandlerFindPrivateKey(t *testing.T) {
	t.Run("Test_Helper/Run_FindPrivateKey_With_RSA", func(t *testing.T) {
		RSIDKS, _ := jwk.GenerateJWK(context.Background(), jose.RS256, "test-id-1", "sig")
		keys, err := jwk.FindPrivateKey(RSIDKS)
		require.NoError(t, err)
		assert.Equal(t, keys.KeyID, "test-id-1")
		assert.IsType(t, keys.Key, new(rsa.PrivateKey))
	})

	t.Run("Test_Helper/Run_FindPrivateKey_With_ECDSA", func(t *testing.T) {
		ECDSAIDKS, err := jwk.GenerateJWK(context.Background(), jose.ES256, "test-id-2", "sig")
		require.NoError(t, err)
		keys, err := jwk.FindPrivateKey(ECDSAIDKS)
		require.NoError(t, err)
		assert.Equal(t, keys.KeyID, "test-id-2")
		assert.IsType(t, keys.Key, new(ecdsa.PrivateKey))
	})

	t.Run("Test_Helper/Run_FindPrivateKey_With_EdDSA", func(t *testing.T) {
		EdDSAIDKS, err := jwk.GenerateJWK(context.Background(), jose.EdDSA, "test-id-3", "sig")
		require.NoError(t, err)
		keys, err := jwk.FindPrivateKey(EdDSAIDKS)
		require.NoError(t, err)
		assert.Equal(t, keys.KeyID, "test-id-3")
		assert.IsType(t, keys.Key, ed25519.PrivateKey{})
	})

	t.Run("Test_Helper/Run_FindPrivateKey_With_KeyNotFound", func(t *testing.T) {
		keySet := &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{}}
		_, err := jwk.FindPublicKey(keySet)
		require.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "key not found"))
	})
}

func TestPEMBlockForKey(t *testing.T) {
	t.Run("Test_Helper/Run_PEMBlockForKey_With_RSA", func(t *testing.T) {
		RSIDKS, err := jwk.GenerateJWK(context.Background(), jose.RS256, "test-id-1", "sig")
		require.NoError(t, err)
		key, err := jwk.FindPrivateKey(RSIDKS)
		require.NoError(t, err)
		pemBlock, err := jwk.PEMBlockForKey(key.Key)
		require.NoError(t, err)
		assert.IsType(t, pem.Block{}, *pemBlock)
		assert.Equal(t, "RSA PRIVATE KEY", pemBlock.Type)
	})

	t.Run("Test_Helper/Run_PEMBlockForKey_With_ECDSA", func(t *testing.T) {
		ECDSAIDKS, err := jwk.GenerateJWK(context.Background(), jose.ES256, "test-id-2", "sig")
		require.NoError(t, err)
		key, err := jwk.FindPrivateKey(ECDSAIDKS)
		require.NoError(t, err)
		pemBlock, err := jwk.PEMBlockForKey(key.Key)
		require.NoError(t, err)
		assert.IsType(t, pem.Block{}, *pemBlock)
		assert.Equal(t, "EC PRIVATE KEY", pemBlock.Type)
	})

	t.Run("Test_Helper/Run_PEMBlockForKey_With_EdDSA", func(t *testing.T) {
		EdDSAIDKS, err := jwk.GenerateJWK(context.Background(), jose.EdDSA, "test-id-3", "sig")
		require.NoError(t, err)
		key, err := jwk.FindPrivateKey(EdDSAIDKS)
		require.NoError(t, err)
		pemBlock, err := jwk.PEMBlockForKey(key.Key)
		require.NoError(t, err)
		assert.IsType(t, pem.Block{}, *pemBlock)
		assert.Equal(t, "PRIVATE KEY", pemBlock.Type)
	})

	t.Run("Test_Helper/Run_PEMBlockForKey_With_InvalidKeyType", func(t *testing.T) {
		key := dsa.PrivateKey{}
		_, err := jwk.PEMBlockForKey(key)
		require.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "Invalid key type"))
	})
}

func TestExcludeOpaquePrivateKeys(t *testing.T) {
	opaqueKeys, err := jwk.GenerateJWK(context.Background(), jose.RS256, "test-id-1", "sig")
	assert.NoError(t, err)
	require.Len(t, opaqueKeys.Keys, 1)
	opaqueKeys.Keys[0].Key = cryptosigner.Opaque(opaqueKeys.Keys[0].Key.(*rsa.PrivateKey))
	keys := jwk.ExcludeOpaquePrivateKeys(opaqueKeys)
	require.Len(t, keys.Keys, 0)
}

func TestGetOrGenerateKeys(t *testing.T) {
	reg := internal.NewMockedRegistry(t, &contextx.Default{})

	setId := uuid.NewUUID().String()
	keyId := uuid.NewUUID().String()

	keySet, _ := jwk.GenerateJWK(context.Background(), jose.RS256, keyId, "sig")
	keySetWithoutPrivateKey := &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{keySet.Keys[0].Public()},
	}

	km := func(t *testing.T) *MockManager {
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)
		return NewMockManager(ctrl)
	}

	t.Run("Test_Helper/Run_GetOrGenerateKeys_With_GetKeySetError", func(t *testing.T) {
		keyManager := km(t)
		keyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq(setId)).Return(nil, errors.New("GetKeySetError"))
		privKey, err := jwk.GetOrGenerateKeys(context.TODO(), reg, keyManager, setId, keyId, "RS256")
		assert.Nil(t, privKey)
		assert.EqualError(t, err, "GetKeySetError")
	})

	t.Run("Test_Helper/Run_GetOrGenerateKeys_With_GenerateAndPersistKeySetError", func(t *testing.T) {
		keyManager := km(t)
		keyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq(setId)).Return(nil, errors.Wrap(x.ErrNotFound, ""))
		keyManager.EXPECT().GenerateAndPersistKeySet(gomock.Any(), gomock.Eq(setId), gomock.Eq(keyId), gomock.Eq("RS256"), gomock.Eq("sig")).Return(nil, errors.New("GetKeySetError"))
		privKey, err := jwk.GetOrGenerateKeys(context.TODO(), reg, keyManager, setId, keyId, "RS256")
		assert.Nil(t, privKey)
		assert.EqualError(t, err, "GetKeySetError")
	})

	t.Run("Test_Helper/Run_GetOrGenerateKeys_With_GenerateAndPersistKeySetError", func(t *testing.T) {
		keyManager := km(t)
		keyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq(setId)).Return(keySetWithoutPrivateKey, nil)
		keyManager.EXPECT().GenerateAndPersistKeySet(gomock.Any(), gomock.Eq(setId), gomock.Eq(keyId), gomock.Eq("RS256"), gomock.Eq("sig")).Return(nil, errors.New("GetKeySetError"))
		privKey, err := jwk.GetOrGenerateKeys(context.TODO(), reg, keyManager, setId, keyId, "RS256")
		assert.Nil(t, privKey)
		assert.EqualError(t, err, "GetKeySetError")
	})

	t.Run("Test_Helper/Run_GetOrGenerateKeys_With_GetKeySet_ContainsMissingPrivateKey", func(t *testing.T) {
		keyManager := km(t)
		keyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq(setId)).Return(keySetWithoutPrivateKey, nil)
		keyManager.EXPECT().GenerateAndPersistKeySet(gomock.Any(), gomock.Eq(setId), gomock.Eq(keyId), gomock.Eq("RS256"), gomock.Eq("sig")).Return(keySet, nil)
		privKey, err := jwk.GetOrGenerateKeys(context.TODO(), reg, keyManager, setId, keyId, "RS256")
		assert.NoError(t, err)
		assert.Equal(t, privKey, &keySet.Keys[0])
	})

	t.Run("Test_Helper/Run_GetOrGenerateKeys_With_GenerateAndPersistKeySet_ContainsMissingPrivateKey", func(t *testing.T) {
		keyManager := km(t)
		keyManager.EXPECT().GetKeySet(gomock.Any(), gomock.Eq(setId)).Return(keySetWithoutPrivateKey, nil)
		keyManager.EXPECT().GenerateAndPersistKeySet(gomock.Any(), gomock.Eq(setId), gomock.Eq(keyId), gomock.Eq("RS256"), gomock.Eq("sig")).Return(keySetWithoutPrivateKey, nil).Times(1)
		privKey, err := jwk.GetOrGenerateKeys(context.TODO(), reg, keyManager, setId, keyId, "RS256")
		assert.Nil(t, privKey)
		assert.EqualError(t, err, "key not found")
	})
}
