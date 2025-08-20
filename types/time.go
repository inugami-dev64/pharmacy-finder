package types

import (
	"encoding/json"
	"strconv"
	"time"
)

type Time time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).UnixMilli())
}

func (t *Time) UnmarshalJSON(b []byte) error {
	ts, err := strconv.ParseInt(string(b), 10, 64)

	if err != nil {
		return err
	}

	*t = Time(time.UnixMilli(ts))
	return nil
}
