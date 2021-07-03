package memkv

import (
	"github.com/xp/shorttext-db/errors"
)

const (
	idLen     = 8
	prefixLen = 1 + idLen /*tableID*/ + 2
	// RecordRowKeyLen is public for calculating avgerage row size.
	RecordRowKeyLen       = prefixLen + idLen /*handle*/
	tablePrefixLength     = 1
	recordPrefixSepLength = 2
	metaPrefixLength      = 1
	// MaxOldEncodeValueLen is the maximum len of the old encoding of index value.
	MaxOldEncodeValueLen = 9
)

var (
	tablePrefix     = []byte{'t'}
	recordPrefixSep = []byte("_r")
	indexPrefixSep  = []byte("_i")
	metaPrefix      = []byte{'m'}
)

func hasTablePrefix(key Key) bool {
	return key[0] == tablePrefix[0]
}

func hasRecordPrefixSep(key Key) bool {
	return key[0] == recordPrefixSep[0] && key[1] == recordPrefixSep[1]
}

// DecodeRecordKey decodes the key and gets the tableID, handle.
func DecodeRecordKey(key Key) (tableID int64, handle int64, err error) {
	if len(key) <= prefixLen {
		return 0, 0, errors.Errorf("invalid record key - %v", key)
	}

	if !hasTablePrefix(key) {
		return 0, 0, errors.Errorf("invalid record key - %v", key)
	}

	key = key[tablePrefixLength:]
	key, tableID, err = DecodeInt(key)
	if err != nil {
		return 0, 0, errors.Trace(err)
	}

	if !hasRecordPrefixSep(key) {
		return 0, 0, errors.Errorf("invalid record key - %v", key)
	}

	key = key[recordPrefixSepLength:]
	key, handle, err = DecodeInt(key)
	if err != nil {
		return 0, 0, errors.Trace(err)
	}
	return
}
