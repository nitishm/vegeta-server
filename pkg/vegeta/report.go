package vegeta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"vegeta-server/models"

	vegeta "github.com/tsenart/vegeta/lib"
)

const (
	// JSONFormatString typedef for query param "json"
	JSONFormatString string = "json"
	// TextFormatString typedef for query param "text"
	TextFormatString string = "text"
	// HistogramFormatString typedef for query param "histogram"
	HistogramFormatString string = "histogram"
	// BinaryFormatString typedef for query param "binary"
	BinaryFormatString string = "binary"
	// DefaultBucketString Default Bucket String
	DefaultBucketString string = "0,500ms,1s,1.5s,2s,2.5s,3s"
)

// CreateReportFromReader takes in an io.Reader with the vegeta gob, encoded result and
// returns the decoded result as a byte array
func CreateReportFromReader(reader io.Reader, id string, format Format) ([]byte, error) {
	dec := vegeta.DecoderFor(reader)

	m := vegeta.Metrics{}

	var report vegeta.Report = &m

	var rep vegeta.Reporter

	fs := format.GetFormat()

	switch fs {
	case JSONFormatString:
		// Create a new reporter with the metrics
		rep = vegeta.NewJSONReporter(&m)
	case TextFormatString:
		rep = vegeta.NewTextReporter(&m)
	case HistogramFormatString:
		var hist vegeta.Histogram
		meta := format.GetMetaInfo()
		// Default bucket = "[0,500ms,1s,1.5s,2s,2.5s,3s]"
		if err := hist.Buckets.UnmarshalText([]byte("[" + meta[0] + "]")); err != nil {
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
	case JSONFormatString:
		var jsonReportResponse models.JSONReportResponse
		err = json.Unmarshal(buf.Bytes(), &jsonReportResponse)
		if err != nil {
			return nil, err
		}
		jsonReportResponse.ID = id
		return json.Marshal(jsonReportResponse)
	case TextFormatString, HistogramFormatString:
		return addID(buf, id), nil
	}

	return buf.Bytes(), nil
}

func addID(report *bytes.Buffer, id string) []byte {
	return append([]byte(fmt.Sprintf("ID %s\n", id)), report.Bytes()...)
}
