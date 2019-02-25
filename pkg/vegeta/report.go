package vegeta

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"vegeta-server/models"

	vegeta "github.com/tsenart/vegeta/lib"
)

// Format defines a type for the format query param
type Format string

const (
	// JSONFormat typedef for query param "json"
	JSONFormat Format = "json"
	// TextFormat typedef for query param "text"
	TextFormat Format = "text"
	// HistogramFormat typedef for query param "histogram"
	HistogramFormat Format = "histogram"
	// BinaryFormat typedef for query param "binary"
	BinaryFormat Format = "binary"
	// DefaultBucketString Default Bucket String
	DefaultBucketString string = "0,500ms,1s,1.5s,2s,2.5s,3s"
)

// SplitFormat splits the Format into sub-Formats separated by '-' using StringsToFormat
func (f Format) SplitFormat() (fs []Format) {
	sf := string(f)
	sli := strings.Split(sf, "-")
	for _, s := range sli {
		fs = append(fs, Format(s))
	}
	return
}

// GetFormat returns the actual Format string. It is useful in case of Clubbed Format string
func (f Format) GetFormat() Format {
	fs := f.SplitFormat()
	return fs[0]
}

// GetTimeBucketOfHistogram returns time duration Bucket for the histogram format string
func (f Format) GetTimeBucketOfHistogram() (buk []byte, err error) {
	fs := f.SplitFormat()
	if fs[0] != HistogramFormat {
		return buk, fmt.Errorf("time bucket can be obtained by histogram format only")
	}
	return []byte(fmt.Sprintf("[%s]", fs[1])), err
}

// StringsToFormat takes string slice and returns a new Format string by separating them using '-'
func StringsToFormat(sli ...string) Format {
	if len(sli) == 0 {
		return JSONFormat // Default Format
	}
	str := sli[0]
	for _, s := range sli[1:] {
		str = fmt.Sprintf("%s-%s", str, s)
	}
	return Format(str)
}

// CreateReportFromReader takes in an io.Reader with the vegeta gob, encoded result and
// returns the decoded result as a byte array
func CreateReportFromReader(reader io.Reader, id string, format Format) ([]byte, error) {
	dec := vegeta.DecoderFor(reader)

	m := vegeta.Metrics{}

	var report vegeta.Report = &m

	var rep vegeta.Reporter

	fs := format.GetFormat()

	switch fs {
	case JSONFormat:
		// Create a new reporter with the metrics
		rep = vegeta.NewJSONReporter(&m)
	case TextFormat:
		rep = vegeta.NewTextReporter(&m)
	case HistogramFormat:
		var hist vegeta.Histogram
		buck, err := format.GetTimeBucketOfHistogram()
		if err != nil {
			return nil, err
		}
		if err = hist.Buckets.UnmarshalText(buck); err != nil { // Default bucket = "[0,500ms,1s,1.5s,2s,2.5s,3s]"
			return nil, err
		}
		rep, report = vegeta.NewHistogramReporter(&hist), &hist
	default:
		return nil, fmt.Errorf("format %s not supported", format)
	}

	rc, _ := report.(vegeta.Closer)
decode:
	for {
		var r vegeta.Result
		err := dec.Decode(&r)
		if err != nil {
			if err == io.EOF {
				break decode
			}
			return nil, err
		}

		report.Add(&r)
	}
	if rc != nil {
		rc.Close()
	}

	var b []byte
	buf := bytes.NewBuffer(b)
	err := rep.Report(buf)
	if err != nil {
		return nil, err
	}

	// Add ID to the report
	switch fs {
	case JSONFormat:
		var jsonReportResponse models.JSONReportResponse
		err = json.Unmarshal(buf.Bytes(), &jsonReportResponse)
		if err != nil {
			return nil, err
		}
		jsonReportResponse.ID = id
		return json.Marshal(jsonReportResponse)
	case TextFormat, HistogramFormat:
		return addID(buf, id), nil
	}

	return buf.Bytes(), nil
}

func addID(report *bytes.Buffer, id string) []byte {
	return append([]byte(fmt.Sprintf("ID %s\n", id)), report.Bytes()...)
}
