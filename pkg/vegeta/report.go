package vegeta

import (
	"bytes"
	"fmt"
	"github.com/tsenart/vegeta/lib"
	"io"
)

type Format string

const (
	JSONFormat      Format = "json"
	HistogramFormat        = "histogram"
	TextFormat             = "text"
)

func CreateReportFromReader(reader io.Reader, format Format) (string, error) {
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
		break
	case TextFormat:
		rep = vegeta.NewTextReporter(&m)
		break
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

	return buf.String(), nil
}
