package shopify

// The iso8601 package encodes and decodes time.Time in JSON in
// ISO 8601 format, without subsecond resolution or time zone info.

import "time"

const Format = "2006-01-02"
const jsonFormat = `"` + Format + `"`

var fixedZone = time.FixedZone("", 0)

type IsoDate time.Time

func New(t time.Time) IsoDate {
	return IsoDate(time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		0,
		fixedZone,
	))
}

func (it IsoDate) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(it).Format(jsonFormat)), nil
}

func (it *IsoDate) UnmarshalJSON(data []byte) error {
	t, err := time.ParseInLocation(jsonFormat, string(data), fixedZone)
	if err == nil {
		*it = IsoDate(t)
	}

	return err
}

func (it IsoDate) String() string {
	return time.Time(it).String()
}
