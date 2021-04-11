package httphandler

import (
	"testing"
)

func TestAPIUrlBuilder(t *testing.T) {
	t.Log("start testing api url builder... ")

	t.Log("try to build correct url")

	correctTestcases := []struct {
		base, path string
		parameters map[string]string
		want       string
	}{
		{
			base:       "http://example.com",
			path:       "user",
			parameters: map[string]string{"q1": "v1", "q2": "v2"},
			want:       "http://example.com/user?q1=v1&q2=v2",
		},
		{
			base:       "http://example.com",
			path:       "product",
			parameters: map[string]string{"q1": "v1", "q2": "v2"},
			want:       "http://example.com/product?q1=v1&q2=v2",
		},
		{ // url encoded
			base:       "http://example.com",
			path:       "this will get automatically encoded",
			parameters: map[string]string{"q1": "v1", "q2": "v2", "q3": "this will get encoded as well"},
			want:       "http://example.com/this%20will%20get%20automatically%20encoded?q1=v1&q2=v2&q3=this+will+get+encoded+as+well",
		},
	}

	for i, tt := range correctTestcases {
		result, err := apiEndpointBuilder(tt.base, tt.path, tt.parameters)
		if err != nil {
			t.Fatalf("case %d, got err: %+v", i, err)
		}

		if result != tt.want {
			t.Fatalf("case %d, should be: %s, but got: %s", i, tt.want, result)
		}
	}

	t.Log("... Passed")

}
