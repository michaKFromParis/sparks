package logger

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/michaKFromParis/sparks/utils"
)

func Init() {

	log.SetFormatter(&Formatter{
		HideKeys:        true,
		ShowFullLevel:   true,
		FieldsOrder:     []string{"component", "category"},
		TimestampFormat: "15:04:05",
	})
	log.SetLevel(log.TraceLevel)

	log.Infof("Initialized logger at level %s", log.GetLevel().String())
	// log.Trace("Trace Sample")
	// log.Debug("Debug Sample")
	// log.Info("Info Sample")
	// log.Warn("Warning Sample")
	// log.Error("Error Sample")
	// log.Fatal("Critical Sample")
	// log.Panic("Panic Sample")
}

// modified FROM https://github.com/antonfisher/nested-logrus-formatter/blob/master/formatter.go

// Formatter - logrus formatter, implements logrus.Formatter
type Formatter struct {
	FieldsOrder     []string // default: fields sorted alphabetically
	TimestampFormat string   // default: time.StampMilli = "Jan _2 15:04:05.000"
	HideKeys        bool     // show [fieldValue] instead of [fieldKey:fieldValue]
	NoColors        bool     // disable colors
	NoFieldsColors  bool     // color only level, default is level + fields
	ShowFullLevel   bool     // true to show full level [WARNING] instead [WARN]
}

// Color numbers for stdout
const (
	Black = (iota + 30)
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

func ColorSeq(color int) string {
	return fmt.Sprintf("\033[%dm", color)
}

func ColorSeqFaint(color int) string {
	return fmt.Sprintf("\033[2m\033[%dm", color)
}

func ColorSeqBold(color int) string {
	return fmt.Sprintf("\033[%d;1m", color)
}

func getColorByLevel(level log.Level) (string, string) {
	switch level {
	case log.TraceLevel:
		return ColorSeq(Blue), ColorSeqFaint(White)
	case log.DebugLevel:
		return ColorSeq(Cyan), ColorSeqFaint(White)
	case log.InfoLevel:
		return ColorSeq(Green), ColorSeq(White)
	case log.WarnLevel:
		return ColorSeq(Yellow), ColorSeq(White)
	case log.ErrorLevel:
		return ColorSeq(Red), ColorSeqBold(White)
	case log.FatalLevel, log.PanicLevel:
		return ColorSeqBold(Red), ColorSeqBold(White)
	default:
		return ColorSeq(White), ColorSeq(White)
	}
}

// Format an log entry
func (f *Formatter) Format(entry *log.Entry) ([]byte, error) {
	timeColor := ColorSeq(Yellow)
	resetColor := ColorSeq(0)
	levelColor, messageColor := getColorByLevel(entry.Level)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.StampMilli
	}

	// output buffer
	b := &bytes.Buffer{}

	b.WriteString(fmt.Sprintf("[%s%s%s]", timeColor, entry.Time.Format(timestampFormat), resetColor))

	// write level
	level := strings.Title(entry.Level.String())
	if f.ShowFullLevel {
		level = fmt.Sprintf("%-7s", level)
	} else {
		level = level[:4]
	}

	b.WriteString(fmt.Sprintf("[%s%s%s] ", levelColor, level, resetColor))

	// write fields
	if f.FieldsOrder == nil {
		f.writeFields(b, entry)
	} else {
		f.writeOrderedFields(b, entry)
	}
	if !f.NoColors && !f.NoFieldsColors {
		b.WriteString(ColorSeq(0))
	}

	// write message
	b.WriteString(fmt.Sprintf("%s%s%s", messageColor, entry.Message, resetColor))
	b.WriteString(utils.LineBreak)
	return b.Bytes(), nil
}

func (f *Formatter) writeFields(b *bytes.Buffer, entry *log.Entry) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}
		sort.Strings(fields)
		for _, field := range fields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *Formatter) writeOrderedFields(b *bytes.Buffer, entry *log.Entry) {
	length := len(entry.Data)
	foundFieldsMap := map[string]bool{}
	for _, field := range f.FieldsOrder {
		if _, ok := entry.Data[field]; ok {
			foundFieldsMap[field] = true
			length--
			f.writeField(b, entry, field)
		}
	}

	if length > 0 {
		notFoundFields := make([]string, 0, length)
		for field := range entry.Data {
			if foundFieldsMap[field] == false {
				notFoundFields = append(notFoundFields, field)
			}
		}

		sort.Strings(notFoundFields)

		for _, field := range notFoundFields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *Formatter) writeField(b *bytes.Buffer, entry *log.Entry, field string) {
	if f.HideKeys {
		fmt.Fprintf(b, "[%v] ", entry.Data[field])
	} else {
		fmt.Fprintf(b, "[%s:%v] ", field, entry.Data[field])
	}
}
