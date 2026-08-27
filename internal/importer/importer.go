package importer

import (
	"bufio"
	"encoding/json"
	"io"
	"paymentconsole/internal/model"
)

func ReadJSON(r io.Reader) ([]model.Record, error) {
	sc := bufio.NewScanner(r)
	out := []model.Record{}
	for sc.Scan() {
		var x model.Record
		if json.Unmarshal(sc.Bytes(), &x) != nil {
			return nil, sc.Err()
		}
		out = append(out, x)
	}
	return out, sc.Err()
}
func Validate(records []model.Record) error {
	seen := map[string]bool{}
	for _, r := range records {
		if r.ID == "" || seen[r.ID] {
			return io.ErrUnexpectedEOF
		}
		seen[r.ID] = true
	}
	return nil
}
