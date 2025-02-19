package sel

import (
	"encoding/hex"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/percona/percona-backup-mongodb/pbm/archive"
)

func IsSelective(ids []string) bool {
	for _, ns := range ids {
		if ns != "" && ns != "*.*" {
			return true
		}
	}

	return false
}

func MakeSelectedPred(nss []string) archive.NSFilterFn {
	if len(nss) == 0 {
		return func(string) bool { return true }
	}

	m := make(map[string]map[string]bool)

	for _, ns := range nss {
		db, coll, _ := strings.Cut(ns, ".")
		if db == "*" {
			db = ""
		}
		if coll == "*" {
			coll = ""
		}

		if m[db] == nil {
			m[db] = make(map[string]bool)
		}
		if !m[db][coll] {
			m[db][coll] = true
		}
	}

	return func(ns string) bool {
		db, coll, ok := strings.Cut(ns, ".")
		return (m[""] != nil || m[db][""]) || (ok && m[db][coll])
	}
}

type ChunkSelector interface {
	Add(bson.Raw)
	Selected(bson.Raw) bool

	BuildFilter() bson.D
}

type nsChunkMap map[string]struct{}

func NewNSChunkSelector() ChunkSelector {
	return make(nsChunkMap)
}

func (s nsChunkMap) Add(d bson.Raw) {
	ns := d.Lookup("_id").StringValue()
	s[ns] = struct{}{}
}

func (s nsChunkMap) Selected(d bson.Raw) bool {
	ns := d.Lookup("ns").StringValue()
	_, ok := s[ns]
	return ok
}

func (s nsChunkMap) BuildFilter() bson.D {
	nss := make([]string, 0, len(s))
	for ns := range s {
		nss = append(nss, ns)
	}

	return bson.D{{"ns", bson.M{"$in": nss}}}
}

type uuidChunkMap map[string]struct{}

func NewUUIDChunkSelector() ChunkSelector {
	return make(uuidChunkMap)
}

func (s uuidChunkMap) Add(d bson.Raw) {
	_, data := d.Lookup("uuid").Binary()
	s[hex.EncodeToString(data)] = struct{}{}
}

func (s uuidChunkMap) Selected(d bson.Raw) bool {
	_, data := d.Lookup("uuid").Binary()
	_, ok := s[hex.EncodeToString(data)]
	return ok
}

func (s uuidChunkMap) BuildFilter() bson.D {
	uuids := make([]primitive.Binary, 0, len(s))
	for ns := range s {
		data, _ := hex.DecodeString(ns)
		uuids = append(uuids, primitive.Binary{Subtype: 0x4, Data: data})
	}

	return bson.D{{"uuid", bson.M{"$in": uuids}}}
}
