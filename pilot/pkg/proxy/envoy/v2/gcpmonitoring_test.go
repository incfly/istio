// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v2

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gm "istio.io/istio/pilot/pkg/gcpmonitoring"
	"istio.io/istio/pkg/test/util/retry"
	"istio.io/pkg/monitoring"
)

var (
	successTestTag = tag.MustNewKey("success")
	typeTestTag    = tag.MustNewKey("type")
)

func TestGCPMonitoringPilotXDSPushMetrics(t *testing.T) {
	os.Setenv("ENABLE_STACKDRIVER_MONITORING", "true")
	defer os.Unsetenv("ENABLE_STACKDRIVER_MONITORING")
	exp := &gm.TestExporter{Rows: make(map[string][]*view.Row)}
	view.RegisterExporter(exp)
	view.SetReportingPeriod(1 * time.Millisecond)

	var cases = []struct {
		name       string
		increment  func()
		wantMetric string
		wantVal    *view.Row
	}{
		{"cdsPushes", incrementCDSPush, "config_push_count", &view.Row{
			Tags: []tag.Tag{{Key: successTestTag, Value: "true"}, {Key: typeTestTag, Value: "CDS"}}, Data: &view.SumData{1.0}}},
		{"edsPushes", incrementEDSPush, "config_push_count", &view.Row{
			Tags: []tag.Tag{{Key: successTestTag, Value: "true"}, {Key: typeTestTag, Value: "EDS"}}, Data: &view.SumData{1.0}}},
		{"ldsPushes", incrementLDSPush, "config_push_count", &view.Row{
			Tags: []tag.Tag{{Key: successTestTag, Value: "true"}, {Key: typeTestTag, Value: "LDS"}}, Data: &view.SumData{1.0}}},
		{"rdsPushes", incrementRDSPush, "config_push_count", &view.Row{
			Tags: []tag.Tag{{Key: successTestTag, Value: "true"}, {Key: typeTestTag, Value: "RDS"}}, Data: &view.SumData{1.0}}},
		{"apiPushes", incrementAPIPush, "config_push_count", &view.Row{
			Tags: []tag.Tag{{Key: successTestTag, Value: "true"}, {Key: typeTestTag, Value: "API"}}, Data: &view.SumData{1.0}}},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			exp.Lock()
			exp.Rows = make(map[string][]*view.Row)
			exp.Unlock()

			tt.increment()
			verifyMetric(t, exp, tt.wantMetric, tt.wantVal)
		})
	}
}

func TestGCPMonitoringPilotXDSPushErrorMetrics(t *testing.T) {
	os.Setenv("ENABLE_STACKDRIVER_MONITORING", "true")
	defer os.Unsetenv("ENABLE_STACKDRIVER_MONITORING")
	exp := &gm.TestExporter{Rows: make(map[string][]*view.Row)}
	view.RegisterExporter(exp)
	view.SetReportingPeriod(1 * time.Millisecond)

	var cases = []struct {
		name       string
		xdsType    string
		metric     monitoring.Metric
		wantMetric string
		wantVal    *view.Row
	}{
		{"cdsSendErrPushes", "CDS", cdsSendErrPushes, "config_push_count", &view.Row{
			Tags: []tag.Tag{{Key: successTestTag, Value: "false"}, {Key: typeTestTag, Value: "CDS"}}, Data: &view.SumData{1.0}}},
		{"edsSendErrPushes", "EDS", edsSendErrPushes, "config_push_count", &view.Row{
			Tags: []tag.Tag{{Key: successTestTag, Value: "false"}, {Key: typeTestTag, Value: "EDS"}}, Data: &view.SumData{1.0}}},
		{"ldsSendErrPushes", "LDS", ldsSendErrPushes, "config_push_count", &view.Row{
			Tags: []tag.Tag{{Key: successTestTag, Value: "false"}, {Key: typeTestTag, Value: "LDS"}}, Data: &view.SumData{1.0}}},
		{"rdsSendErrPushes", "RDS", rdsSendErrPushes, "config_push_count", &view.Row{
			Tags: []tag.Tag{{Key: successTestTag, Value: "false"}, {Key: typeTestTag, Value: "RDS"}}, Data: &view.SumData{1.0}}},
		{"apiSendErrPushes", "API", apiSendErrPushes, "config_push_count", &view.Row{
			Tags: []tag.Tag{{Key: successTestTag, Value: "false"}, {Key: typeTestTag, Value: "API"}}, Data: &view.SumData{1.0}}},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			exp.Lock()
			exp.Rows = make(map[string][]*view.Row)
			exp.Unlock()

			recordSendError(tt.xdsType, tt.metric, status.Error(codes.DeadlineExceeded, "dummy"))
			verifyMetric(t, exp, tt.wantMetric, tt.wantVal)
		})
	}
}

func TestGCPMonitoringPilotXDSRejectMetrics(t *testing.T) {
	os.Setenv("ENABLE_STACKDRIVER_MONITORING", "true")
	defer os.Unsetenv("ENABLE_STACKDRIVER_MONITORING")
	exp := &gm.TestExporter{Rows: make(map[string][]*view.Row)}
	view.RegisterExporter(exp)
	view.SetReportingPeriod(1 * time.Millisecond)

	var cases = []struct {
		name       string
		xdsType    string
		metric     monitoring.Metric
		wantMetric string
		wantVal    *view.Row
	}{
		{"cdsReject", "CDS", cdsReject, "rejected_config_count", &view.Row{
			Tags: []tag.Tag{{Key: typeTestTag, Value: "CDS"}}, Data: &view.SumData{1.0}}},
		{"edsReject", "EDS", edsReject, "rejected_config_count", &view.Row{
			Tags: []tag.Tag{{Key: typeTestTag, Value: "EDS"}}, Data: &view.SumData{1.0}}},
		{"ldsReject", "LDS", ldsReject, "rejected_config_count", &view.Row{
			Tags: []tag.Tag{{Key: typeTestTag, Value: "LDS"}}, Data: &view.SumData{1.0}}},
		{"rdsReject", "RDS", rdsReject, "rejected_config_count", &view.Row{
			Tags: []tag.Tag{{Key: typeTestTag, Value: "RDS"}}, Data: &view.SumData{1.0}}},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			exp.Lock()
			exp.Rows = make(map[string][]*view.Row)
			exp.Unlock()

			incrementXDSRejects(tt.xdsType, tt.metric, "dummy", "dummy")
			verifyMetric(t, exp, tt.wantMetric, tt.wantVal)
		})
	}
}

func TestGCPMonitoringXDSClientsMetrics(t *testing.T) {
	os.Setenv("ENABLE_STACKDRIVER_MONITORING", "true")
	defer os.Unsetenv("ENABLE_STACKDRIVER_MONITORING")
	exp := &gm.TestExporter{Rows: make(map[string][]*view.Row)}
	view.RegisterExporter(exp)
	view.SetReportingPeriod(1 * time.Millisecond)

	exp.Lock()
	exp.Rows = make(map[string][]*view.Row)
	exp.Unlock()

	recordProxyClients(10)
	wantMetric := "proxy_clients"
	wantVal := &view.Row{Tags: []tag.Tag{}, Data: &view.LastValueData{Value: 10.0}}
	verifyMetric(t, exp, wantMetric, wantVal)
}

func TestGCPMonitoringConvergencyLatencyMetrics(t *testing.T) {
	os.Setenv("ENABLE_STACKDRIVER_MONITORING", "true")
	defer os.Unsetenv("ENABLE_STACKDRIVER_MONITORING")
	exp := &gm.TestExporter{Rows: make(map[string][]*view.Row)}
	view.RegisterExporter(exp)
	view.SetReportingPeriod(1 * time.Millisecond)

	exp.Lock()
	exp.Rows = make(map[string][]*view.Row)
	exp.Unlock()

	recordConvergencyDeley(0.4)
	wantMetric := "config_convergence_latencies"
	wantVal := &view.Row{
		Tags: []tag.Tag{},
		Data: &view.DistributionData{
			CountPerBucket: []int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
		},
	}
	verifyMetric(t, exp, wantMetric, wantVal)
}

func verifyMetric(t *testing.T, exp *gm.TestExporter, wantMetric string, wantVal *view.Row) {
	if err := retry.UntilSuccess(func() error {
		exp.Lock()
		defer exp.Unlock()
		if len(exp.Rows[wantMetric]) < 1 {
			return fmt.Errorf("wanted metrics %v not received", wantMetric)
		}
		for _, got := range exp.Rows[wantMetric] {
			if len(got.Tags) != len(wantVal.Tags) ||
				(len(wantVal.Tags) != 0 && !reflect.DeepEqual(got.Tags, wantVal.Tags)) {
				continue
			}
			switch v := wantVal.Data.(type) {
			case *view.SumData:
				if int64(v.Value) == int64(got.Data.(*view.SumData).Value) {
					return nil
				}
			case *view.LastValueData:
				if int64(v.Value) == int64(got.Data.(*view.LastValueData).Value) {
					return nil
				}
			case *view.DistributionData:
				gotDist := got.Data.(*view.DistributionData)
				if reflect.DeepEqual(gotDist.CountPerBucket, v.CountPerBucket) {
					return nil
				}
			}
		}
		return fmt.Errorf("metrics %v does not have expected values, want %+v", wantMetric, wantVal)
	}); err != nil {
		t.Fatalf("failed to get expected metric: %v", err)
	}
}
