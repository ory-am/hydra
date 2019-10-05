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

package jwk

import (
	"context"
	"crypto/rand"
	"io"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/square/go-jose.v2"
)

func RandomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return []byte{}, errors.WithStack(err)
	}
	return bytes, nil
}

func TestHelperManagerKey(m Manager, keys *jose.JSONWebKeySet, suffix string) func(t *testing.T) {
	pub := keys.Key("public:" + suffix)
	priv := keys.Key("private:" + suffix)

	return func(t *testing.T) {
		_, err := m.GetKey(context.TODO(), "faz", "baz")
		assert.NotNil(t, err)

		err = m.AddKey(context.TODO(), "faz", First(priv))
		require.NoError(t, err)

		got, err := m.GetKey(context.TODO(), "faz", "private:"+suffix)
		require.NoError(t, err)
		assert.Equal(t, priv, got.Keys)

		err = m.AddKey(context.TODO(), "faz", First(pub))
		require.NoError(t, err)

		got, err = m.GetKey(context.TODO(), "faz", "private:"+suffix)
		require.NoError(t, err)
		assert.Equal(t, priv, got.Keys)

		got, err = m.GetKey(context.TODO(), "faz", "public:"+suffix)
		require.NoError(t, err)
		assert.Equal(t, pub, got.Keys)

		// Because MySQL
		time.Sleep(time.Second * 2)

		First(pub).KeyID = "new-key-id:" + suffix
		err = m.AddKey(context.TODO(), "faz", First(pub))
		require.NoError(t, err)

		_, err = m.GetKey(context.TODO(), "faz", "new-key-id:"+suffix)
		require.NoError(t, err)

		keys, err = m.GetKeySet(context.TODO(), "faz")
		require.NoError(t, err)
		assert.EqualValues(t, "new-key-id:"+suffix, First(keys.Keys).KeyID)

		beforeDeleteKeysCount := len(keys.Keys)
		err = m.DeleteKey(context.TODO(), "faz", "public:"+suffix)
		require.NoError(t, err)

		_, err = m.GetKey(context.TODO(), "faz", "public:"+suffix)
		require.Error(t, err)

		keys, err = m.GetKeySet(context.TODO(), "faz")
		require.NoError(t, err)
		assert.EqualValues(t, beforeDeleteKeysCount-1, len(keys.Keys))
	}
}

type keyGenerator interface {
	Generate(id, use string) (*jose.JSONWebKeySet, error)
}

func TestHelperManagerKeySet(m Manager, generator keyGenerator) func(t *testing.T) {
	return func(t *testing.T) {
		_, err := m.GetKeySet(context.TODO(), "foo")
		require.Error(t, err)

		oldKs, err := generator.Generate("OldTestManagerKeySet", "sig")
		require.NoError(t, err)
		err = m.AddKeySet(context.TODO(), "bar", oldKs)
		require.NoError(t, err)

		// To make difference timestamp of second
		time.Sleep(1 * time.Second)
		afterCreatedOldKeys := time.Now().UTC()

		// To delay creation timestamp of new keys for DeleteOldKeys()
		time.Sleep(1 * time.Second)

		newKs, err := generator.Generate("NewTestManagerKeySet", "sig")
		require.NoError(t, err)
		err = m.AddKeySet(context.TODO(), "bar", newKs)
		require.NoError(t, err)

		got, err := m.GetKeySet(context.TODO(), "bar")
		require.NoError(t, err)
		assert.Equal(t, oldKs.Key("public:OldTestManagerKeySet"), got.Key("public:OldTestManagerKeySet"))
		assert.Equal(t, oldKs.Key("private:OldTestManagerKeySet"), got.Key("private:OldTestManagerKeySet"))

		err = m.DeleteOldKeys(context.TODO(), "bar", afterCreatedOldKeys)
		require.NoError(t, err)
		got2, err := m.GetKeySet(context.TODO(), "bar")
		require.NoError(t, err)
		assert.Equal(t, 0, len(got2.Key("public:OldTestManagerKeySet")))
		assert.Equal(t, 1, len(got2.Key("public:NewTestManagerKeySet")))

		err = m.DeleteKeySet(context.TODO(), "bar")
		require.NoError(t, err)

		_, err = m.GetKeySet(context.TODO(), "bar")
		require.Error(t, err)
	}
}
