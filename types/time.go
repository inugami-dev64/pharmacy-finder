package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Time time.Time

func (t Time) Value() (driver.Value, error) {
	return time.Time(t), nil
}

func (t *Time) Scan(src interface{}) error {
	if v, ok := src.(time.Time); ok {
		*t = Time(v)
		return nil
	}

	return fmt.Errorf("cannot scan type %T as types.Time", src)
}

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
