package vegeta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"vegeta-server/models"

	vegeta "github.com/tsenart/vegeta/lib"
)

type Format string

const (
	JSONFormat      Format = "json"
	TextFormat      Format = "text"
	HistogramFormat Format = "histogram"
)

func CreateReportFromReader(reader io.Reader, id string, format Format) (string, error) {
	dec := vegeta.DecoderFor(reader)

	m := vegeta.Metrics{}

	var report vegeta.Report = &m

decode:
	for {
		var r vegeta.Result
		err := dec.Decode(&r)
		if err != nil {
			if err == io.EOF {
				break decode
			}
			return "", err
		}

		report.Add(&r)
	}

	rc := report.(vegeta.Closer)
	rc.Close()

	var rep vegeta.Reporter

	switch format {
	case JSONFormat:
		// Create a new reporter with the metrics
		rep = vegeta.NewJSONReporter(&m)
	case TextFormat:
		rep = vegeta.NewTextReporter(&m)
	// TODO: Figure out how to provide historgram report
	//case HistogramFormat:
	//	var hist vegeta.Histogram
	//	if err := hist.Buckets.UnmarshalText([]byte(typ[4:])); err != nil {
	//		return err
	//	}
	default:
		return "", fmt.Errorf("format %s not supported", format)
	}

	var b []byte
	buf := bytes.NewBuffer(b)
	err := rep.Report(buf)
	if err != nil {
		return "", err
	}

	if format == JSONFormat {
		// Add ID to the report struct
		var jsonReportResponse models.JSONReportResponse
		err = json.Unmarshal(buf.Bytes(), &jsonReportResponse)
		if err != nil {
			return "", err
		}
		jsonReportResponse.ID = id
		bReport, err := json.Marshal(jsonReportResponse)
		if err != nil {
			return "", err
		}
		return string(bReport), nil
	}

	return buf.String(), nil
}
