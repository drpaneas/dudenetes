/*Package cmd test
Copyright © 2019 Panagiotis Georgiadis <drpaneas@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"reflect"
	"testing"
)

func Test_unmarshalOutputResult(t *testing.T) {
	tests := []struct {
		name    string
		out     []byte
		want    []string
		want1   []string
		want2   []string
		wantErr bool
	}{
		{"vmware no error", []byte(vmwareOutput), []string{"11.45.155.189"}, []string{"11.45.72.50", "11.45.73.33", "11.45.72.44"}, []string{"11.45.72.183", "11.45.72.86", "11.45.72.175"}, false},
		{"openstack no error", []byte(openstackOutput), []string{"11.45.155.189"}, []string{"11.45.72.50", "11.45.73.33", "11.45.72.44"}, []string{"11.45.72.183", "11.45.72.86", "11.45.72.175"}, false},
		{"error in unmarshall", []byte("!"), nil, nil, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := unmarshalOutputResult(tt.out)
			if (err != nil) != tt.wantErr {
				t.Errorf("unmarshalOutputResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unmarshalOutputResult() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("unmarshalOutputResult() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("unmarshalOutputResult() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_formatValues(t *testing.T) {
	tests := []struct {
		name          string
		loadbalancers []string
		master        []string
		worker        []string
		want          []string
	}{
		{"format", []string{"11.45.155.189"}, []string{"11.45.72.50", "11.45.73.33", "11.45.72.44"}, []string{"11.45.72.183", "11.45.72.86", "11.45.72.175"}, formatedResult},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatValues(tt.loadbalancers, tt.master, tt.worker); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("formatValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

var vmwareOutput = `{
    "ip_load_balancer": {
        "sensitive": false,
        "type": "list",
        "value": [
            "11.45.155.189"
        ]
    },
    "ip_masters": {
        "sensitive": false,
        "type": "list",
        "value": [
            "11.45.72.50",
            "11.45.73.33",
            "11.45.72.44"
        ]
    },
    "ip_workers": {
        "sensitive": false,
        "type": "list",
        "value": [
            "11.45.72.183",
            "11.45.72.86",
            "11.45.72.175"
        ]
    }
}
`

var openstackOutput = `{
    "hostnames_masters": {
        "sensitive": false,
        "type": "list",
        "value": []
    },
    "ip_internal_load_balancer": {
        "sensitive": false,
        "type": "string",
        "value": "172.28.0.19"
    },
    "ip_load_balancer": {
        "sensitive": false,
        "type": "string",
        "value": "11.45.155.189"
    },
    "ip_masters": {
        "sensitive": false,
        "type": "list",
        "value": [
            "11.45.72.50",
            "11.45.73.33",
            "11.45.72.44"
        ]
    },
    "ip_workers": {
        "sensitive": false,
        "type": "list",
        "value": [
            "11.45.72.183",
            "11.45.72.86",
            "11.45.72.175"
        ]
    }
}
`

var formatedResult = []string{
	"loadbalancer=11.45.155.189",
	"master1=11.45.72.50",
	"master2=11.45.73.33",
	"master3=11.45.72.44",
	"worker1=11.45.72.183",
	"worker2=11.45.72.86",
	"worker3=11.45.72.175"}
