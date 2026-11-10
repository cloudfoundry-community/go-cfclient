package client_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/cloudfoundry/go-cfclient/v3/client"
	"github.com/stretchr/testify/require"
)

func TestTimestampFilterSerializeUTC(t *testing.T) {
	tests := []struct {
		name     string
		operator client.FilterModifier
		inputs   []string
		key      string
		want     string
	}{
		{"UTC unchanged", client.FilterModifierNone, []string{"2025-05-10T09:54:11Z"}, "updated_ats", "2025-05-10T09:54:11Z"},
		{"positive offset before", client.FilterModifierLessThan, []string{"2025-05-10T11:54:11+02:00"}, "updated_ats[lt]", "2025-05-10T09:54:11Z"},
		{"negative offset after", client.FilterModifierGreaterThan, []string{"2025-05-10T02:54:11-07:00"}, "updated_ats[gt]", "2025-05-10T09:54:11Z"},
		{"previous day before or equal", client.FilterModifierLessThanOrEqual, []string{"2025-01-01T00:30:00+02:00"}, "updated_ats[lte]", "2024-12-31T22:30:00Z"},
		{"next day after or equal", client.FilterModifierGreaterThanOrEqual, []string{"2024-12-31T23:30:00-07:00"}, "updated_ats[gte]", "2025-01-01T06:30:00Z"},
		{"multiple equality timestamps", client.FilterModifierNone, []string{"2025-05-10T11:54:11+02:00", "2025-05-11T02:54:11-07:00", "2025-05-12T09:54:11Z"}, "updated_ats", "2025-05-10T09:54:11Z,2025-05-11T09:54:11Z,2025-05-12T09:54:11Z"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timestamps := make([]time.Time, len(tt.inputs))
			for i, input := range tt.inputs {
				var err error
				timestamps[i], err = time.Parse(time.RFC3339, input)
				require.NoError(t, err)
			}
			original := append([]time.Time(nil), timestamps...)
			filter := client.TimestampFilter{Timestamp: timestamps, Operator: tt.operator}
			values := url.Values{}
			require.NoError(t, filter.Serialize(values, "updated_ats"))
			require.Equal(t, url.Values{tt.key: []string{tt.want}}, values)
			require.Equal(t, original, timestamps, "serialization must preserve caller timestamps")
		})
	}
}

func TestServiceInstanceTimestampFilterUTC(t *testing.T) {
	timestamp, err := time.Parse(time.RFC3339, "2025-05-10T11:54:11+02:00")
	require.NoError(t, err)
	filters := client.TimestampFilterList{}
	filters.Before(timestamp)
	opts := client.ServiceInstanceListOptions{ListOptions: &client.ListOptions{UpdatedAts: filters}}
	values, err := opts.ToQueryString()
	require.NoError(t, err)
	require.Equal(t, "updated_ats%5Blt%5D=2025-05-10T09%3A54%3A11Z", values.Encode())
}
