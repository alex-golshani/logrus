package logrus

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestErrorNotLost(t *testing.T) {
	formatter := &JSONFormatter{}

	b, err := formatter.Format(WithField("error", errors.New("wild walrus")))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["error"] != "wild walrus" {
		t.Fatal("Error field not set")
	}
}

func TestErrorNotLostOnFieldNotNamedError(t *testing.T) {
	formatter := &JSONFormatter{}

	b, err := formatter.Format(WithField("omg", errors.New("wild walrus")))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["omg"] != "wild walrus" {
		t.Fatal("Error field not set")
	}
}

func TestFieldClashWithTime(t *testing.T) {
	formatter := &JSONFormatter{}

	b, err := formatter.Format(WithField("time", "right now!"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["fields.time"] != "right now!" {
		t.Fatal("fields.time not set to original time field")
	}

	if entry[timeKey] != "0001-01-01T00:00:00Z" {
		t.Fatal("time field not set to current time, was: ", entry[timeKey])
	}
}

func TestFieldClashWithMsg(t *testing.T) {
	formatter := &JSONFormatter{}

	b, err := formatter.Format(WithField(messageKey, "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["fields.msg"] != "something" {
		t.Fatal("fields.msg not set to original msg field")
	}
}

func TestFieldClashWithLevel(t *testing.T) {
	formatter := &JSONFormatter{}

	b, err := formatter.Format(WithField(levelKey, "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["fields.level"] != "something" {
		t.Fatal("fields.level not set to original level field")
	}
}

func TestJSONEntryEndsWithNewline(t *testing.T) {
	formatter := &JSONFormatter{}

	b, err := formatter.Format(WithField(levelKey, "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	if b[len(b)-1] != '\n' {
		t.Fatal("Expected JSON log entry to end with a newline")
	}
}

func TestJSONMessageKey(t *testing.T) {
	formatter := &JSONFormatter{
		FieldMap: FieldMap{
			messageKey: "Message",
		},
	}

	b, err := formatter.Format(&Entry{Message: "oh hai"})
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !(strings.Contains(s, "Message") && strings.Contains(s, "oh hai")) {
		t.Fatal("Expected JSON to format Message key")
	}
}

func TestJSONLevelKey(t *testing.T) {
	formatter := &JSONFormatter{
		FieldMap: FieldMap{
			levelKey: "somelevel",
		},
	}

	b, err := formatter.Format(WithField(levelKey, "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !strings.Contains(s, "somelevel") {
		t.Fatal("Expected JSON to format level key")
	}
}

func TestJSONTimeKey(t *testing.T) {
	formatter := &JSONFormatter{
		FieldMap: FieldMap{
			timeKey: "timeywimey",
		},
	}

	b, err := formatter.Format(WithField(levelKey, "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !strings.Contains(s, "timeywimey") {
		t.Fatal("Expected JSON to format time key")
	}
}

func TestJSONDisableTimestamp(t *testing.T) {
	formatter := &JSONFormatter{
		DisableTimestamp: true,
	}

	b, err := formatter.Format(WithField(levelKey, "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if strings.Contains(s, timeKey) {
		t.Error("Did not prevent timestamp", s)
	}
}

func TestJSONEnableTimestamp(t *testing.T) {
	formatter := &JSONFormatter{}

	b, err := formatter.Format(WithField(levelKey, "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !strings.Contains(s, timeKey) {
		t.Error("Timestamp not present", s)
	}
}
