package log

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"github.com/sirupsen/logrus"
	"gopkg.in/gookit/color.v1"
)

func TestCliLoggerFormatter_Format(t *testing.T) {
	emptyEntry := &logrus.Entry{}
	entry := &logrus.Entry{
		Level:   logrus.ErrorLevel,
		Message: "hi log",
	}
	emptyFormatter := &CliLoggerFormatter{}
	formatter := &CliLoggerFormatter{
		prefix:   "pp",
		showType: "error",
	}
	tests := []struct {
		name      string
		formatter *CliLoggerFormatter
		entry     *logrus.Entry
		want      []byte
		wantErr   bool
	}{
		{"base", emptyFormatter, emptyEntry, createBufferForCliLoggerFormatter(t, emptyFormatter, emptyEntry).Bytes(), false},
		{"base debug", formatter, entry, createBufferForCliLoggerFormatter2(t, formatter, entry).Bytes(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// logrus.SetLevel(logrus.DebugLevel)
			got, err := tt.formatter.Format(tt.entry)
			if (err != nil) != tt.wantErr {
				t.Errorf("CliLoggerFormatter.Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CliLoggerFormatter.Format() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func TestCliLoggerFormatter_levelPrintRender(t *testing.T) {
	tests := []struct {
		name      string
		formatter *CliLoggerFormatter
		want      *CliLoggerFormatter
	}{
		// TODO: Add test cases.
		{"base", &CliLoggerFormatter{}, &CliLoggerFormatter{}},
		{"base debug",
			&CliLoggerFormatter{showType: "debug"},
			&CliLoggerFormatter{
				level:           logrus.DebugLevel,
				showType:        "debug",
				formatLevelName: color.Blue.Render(DEBUG),
				prefix:          color.Blue.Render(normal.Debug)},
		},
		{"base info",
			&CliLoggerFormatter{showType: "info"},
			&CliLoggerFormatter{
				level:           logrus.InfoLevel,
				showType:        "info",
				formatLevelName: color.FgLightBlue.Render(INFO),
				prefix:          color.FgLightBlue.Render(normal.Info)},
		},
		{"base warn",
			&CliLoggerFormatter{showType: "warn"},
			&CliLoggerFormatter{
				level:           logrus.WarnLevel,
				showType:        "warn",
				formatLevelName: color.Yellow.Render(WARN),
				prefix:          color.Yellow.Render(normal.Warn)},
		},
		{"base error",
			&CliLoggerFormatter{showType: "error"},
			&CliLoggerFormatter{
				level:           logrus.ErrorLevel,
				showType:        "error",
				formatLevelName: color.BgRed.Render(ERROR),
				prefix:          color.Red.Render(normal.Error)},
		},
		{"base fatal",
			&CliLoggerFormatter{showType: "fatal"},
			&CliLoggerFormatter{
				level:           logrus.FatalLevel,
				showType:        "fatal",
				formatLevelName: color.BgRed.Render(FATAL),
				prefix:          color.Red.Render(normal.Fatal)},
		},
		{"base success",
			&CliLoggerFormatter{showType: "success"},
			&CliLoggerFormatter{
				level:           logrus.InfoLevel,
				showType:        "success",
				formatLevelName: color.Green.Render(SUCCESS),
				prefix:          color.Green.Render(normal.Success)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.formatter.levelPrintRender()
			if !reflect.DeepEqual(tt.formatter, tt.want) {
				t.Errorf("levelPrintRender = \n%v", tt.formatter)
				t.Errorf("want = \n%v", tt.want)
			}
		})
	}
}

func TestSeparatorFormatter_Format(t *testing.T) {
	b := createBufferForSeparatorFormatter(t)
	tests := []struct {
		name    string
		s       *SeparatorFormatter
		entry   *logrus.Entry
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{"base", &SeparatorFormatter{}, &logrus.Entry{}, b.Bytes(), false},
		{"base Entry with buffer", &SeparatorFormatter{}, &logrus.Entry{Buffer: &bytes.Buffer{}}, b.Bytes(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SeparatorFormatter{}
			got, err := s.Format(tt.entry)
			if (err != nil) != tt.wantErr {
				t.Errorf("SeparatorFormatter.Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SeparatorFormatter.Format() = \n%v", string(got))
				t.Errorf("want = \n%v", tt.want)
			}
		})
	}
}

func createBufferForSeparatorFormatter(t *testing.T) *bytes.Buffer {
	var b *bytes.Buffer = &bytes.Buffer{}
	entry := &logrus.Entry{}
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	newLog := fmt.Sprintf("%s %s %s %s\n",
		timestamp,
		color.Blue.Render(normal.Info),
		color.Blue.Render(INFO),
		color.Blue.Render(fmt.Sprintf("%s %s %s", "-------------------- [ ", entry.Message, " ] --------------------")))

	_, err := b.WriteString(newLog)
	if err != nil {
		t.Error(err)
		return nil
	}
	return b
}

func createBufferForCliLoggerFormatter(t *testing.T, m *CliLoggerFormatter, entry *logrus.Entry) *bytes.Buffer {
	var b *bytes.Buffer = &bytes.Buffer{}
	m.levelPrintRender()

	timestamp := entry.Time.Format("2006-01-02 15:04:05")

	newLog := fmt.Sprintf("%s %s %s %s\n", timestamp, m.prefix, m.formatLevelName, entry.Message)

	_, err := b.WriteString(newLog)
	if err != nil {
		t.Error(err)
		return nil
	}
	return b
}

func createBufferForCliLoggerFormatter2(t *testing.T, m *CliLoggerFormatter, entry *logrus.Entry) *bytes.Buffer {
	var b *bytes.Buffer = &bytes.Buffer{}
	m.levelPrintRender()

	timestamp := entry.Time.Format("2006-01-02 15:04:05")

	entry.Message = addCallStackIgnoreLogrus(entry.Message)

	newLog := fmt.Sprintf("%s %s %s %s\n", timestamp, m.prefix, m.formatLevelName, entry.Message)

	_, err := b.WriteString(newLog)
	if err != nil {
		t.Error(err)
		return nil
	}
	return b
}

func Test_addCallStackIgnoreLogrus(t *testing.T) {
	rawMsg := "hi log"
	tests := []struct {
		name       string
		rawMessage string
		want       string
	}{
		// TODO: Add test cases.
		{"base", rawMsg, wantMsgOfAddCallStackIgnoreLogrus(rawMsg)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addCallStackIgnoreLogrus(tt.rawMessage); got != tt.want {
				t.Errorf("addCallStackIgnoreLogrus() = \n%v, \nwant = %v", got, tt.want)
			}
		})
	}
}

func wantMsgOfAddCallStackIgnoreLogrus(stackMessage string) string {
	retMsg := stackMessage
	i := 10
	_, file, line, _ := runtime.Caller(i)
	retMsg = retMsg + "\n  -- " + file + fmt.Sprintf(" %d", line)
	return retMsg
}
