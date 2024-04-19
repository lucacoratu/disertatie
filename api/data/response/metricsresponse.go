package data

import (
	"encoding/json"
	"io"

	"github.com/lucacoratu/disertatie/api/data"
)

type MethodMetricsResponse struct {
	Metrics []data.MethodsMetrics `json:"metrics"` // The list of methods metrics
}

func (mm *MethodMetricsResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(mm)
}

func (mm *MethodMetricsResponse) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(mm)
}

type DaysMetricsResponse struct {
	Metrics []data.DayMetrics `json:"metrics"` // The list of day metrics
}

func (dm *DaysMetricsResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(dm)
}

func (dm *DaysMetricsResponse) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(dm)
}

type StatusCodesMetricsResponse struct {
	Metrics []data.StatusCodesMetrics `json:"metrics"` //The list of metrics
}

func (scm *StatusCodesMetricsResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(scm)
}

func (scm *StatusCodesMetricsResponse) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(scm)
}

type IPAddressMetricsResponse struct {
	Metrics []data.IPMetrics `json:"metrics"` //The list of metrics
}

func (ipm *IPAddressMetricsResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(ipm)
}

func (ipm *IPAddressMetricsResponse) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(ipm)
}

type FindingsMetricsResponse struct {
	Metrics []data.FindingsMetrics `json:"metrics"`
}

func (fmr *FindingsMetricsResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(fmr)
}

func (fmr *FindingsMetricsResponse) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(fmr)
}
