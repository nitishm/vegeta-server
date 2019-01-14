package vegeta

import (
	"bytes"
	"io"

	log "github.com/sirupsen/logrus"
	vegetalib "github.com/tsenart/vegeta/lib"
)

// Result stores the attack uuid and the results buffer
// returned from the attack Run
type Result struct {
	uuid   string
	result io.Reader
}

// ReportInterface is the Reporter Interface
type ReportInterface interface {
	// Get the JSON report by the attack UUID
	Get(uuid string) (string, error)
	// Add the JSON report for the attack by UUID
	Add(uuid string, report string) error
	// List all the attack JSON reports
	List() map[string]string
}

// Reporter implements the ReporterIntf
type Reporter struct {
	store     StorageInterface
	resultsCh <-chan *Result
}

func NewReporter(resCh <-chan *Result) *Reporter {
	rp := &Reporter{
		DefaultStorage,
		resCh,
	}

	// Read the results after an attack completes and convert it into
	// a JSON report.
	go rp.startResultHandler()

	return rp
}

func (rp *Reporter) Get(uuid string) (res string, err error) {
	return rp.store.Get(uuid)
}

func (rp *Reporter) Add(uuid string, report string) error {
	log.WithField("UUID", uuid).Debugf("Adding report to storage - %v", report)
	return rp.store.Add(uuid, report)
}

func (rp *Reporter) List() map[string]string {
	reports := make(map[string]string)
	for uuid, report := range rp.store.List() {
		reports[uuid] = report
	}
	return reports
}

func (rp *Reporter) startResultHandler() {
	for {
		select {
		case r := <-rp.resultsCh:
			log.Debugf("Got results from attack %s", r.uuid)
			// Convert the io.Reader results to a JSON report
			report, err := createReportFromReader(r.result)
			if err != nil {
				log.Error("Failed to create Report from Reader")
				continue
			}

			log.Debugf("Adding report to storage for attack %s", r.uuid)
			err = rp.Add(r.uuid, report)
			if err != nil {
				log.Error("Failed to Add Report to store")
				continue
			}
		default:
			break
		}
	}
}

func createReportFromReader(reader io.Reader) (string, error) {
	dec := vegetalib.DecoderFor(reader)

	m := vegetalib.Metrics{}

	var report vegetalib.Report = &m

decode:
	for {
		var r vegetalib.Result
		err := dec.Decode(&r)
		if err != nil {
			if err == io.EOF {
				break decode
			}
			return "", err
		}

		report.Add(&r)
	}

	rc := report.(vegetalib.Closer)
	rc.Close()

	// Create a new reporter with the metrics
	rep := vegetalib.NewJSONReporter(&m)

	var b []byte
	buf := bytes.NewBuffer(b)
	err := rep.Report(buf)
	if err != nil {
		log.Error("Failed to write report")
		return "", err
	}

	return buf.String(), nil
}
