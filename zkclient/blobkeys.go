// Copyright (c) 2016-2020 Company 0, LLC.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"path"
	"strconv"
	"time"

	"github.com/companyzero/zkc/inidb"
)

const (
	blobKeysDir  = "blobkeys"
	blobKeysFile = "blobkeys.ini"
	expiration   = 24 * 7 * time.Hour
)

var (
	blobKeysPath = path.Join(blobKeysDir, blobKeysFile)
)

func (z *ZKC) pruneKey(kdb *inidb.INIDB) {
	//r := kdb.Records("")
	//for k, v := range r {
	//	t, err := strconv.ParseInt(v, 10, 64)
	//	if err != nil {
	//		// token corrupt, remove from db and complain
	//		z.Error(idZKC, "corrupt token %v", k)
	//		_ = kdb.Del("", k)
	//		continue
	//	}
	//	ts := time.Unix(t, 0)
	//	if ts.Before(time.Now()) {
	//		// token expired, remove from db
	//		_ = kdb.Del("", k)
	//		continue
	//	}
	//}

	// save db back
	err := kdb.Save()
	if err != nil {
		z.Error(idZKC, "could not save keys database: %v", err)
		// falthrough
	}
}

func (z *ZKC) saveKey(key *[32]byte) error {
	// open db
	kdb, err := inidb.New(path.Join(z.settings.Root, blobKeysPath),
		true, 10)
	if err != nil && !errors.Is(err, inidb.ErrCreated) {
		return err
	}

	defer z.pruneKey(kdb) // kill all expired keys

	for {
		k := time.Now().Add(expiration).Unix()
		keyTS := strconv.FormatInt(k, 10)
		// see if it exists
		_, err := kdb.Get("", keyTS)
		if err == nil {
			// really shouldnt happen
			continue
		}

		err = kdb.Set("", keyTS, hex.EncodeToString(key[:]))
		if err != nil {
			return fmt.Errorf("saveKey %v", err)
		}

		break
	}

	return nil
}
