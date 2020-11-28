package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuildUrlQuery(t *testing.T) {
	type args struct {
		url string
		m   map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name:    "test empty m",
			args:    args{url: "http://m.com?a=5", m: nil},
			want:    "http://m.com?a=5",
			wantErr: nil,
		},
		{
			name:    "test empty m",
			args:    args{url: "http://m.com?a=5", m: map[string]interface{}{"b": "6"}},
			want:    "http://m.com?a=5&b=6",
			wantErr: nil,
		},
		{
			name:    "test empty m",
			args:    args{url: "http://m.com", m: map[string]interface{}{"b": "6"}},
			want:    "http://m.com?b=6",
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		u, err := BuildUrlQuery(tt.args.url, tt.args.m)
		require.Equal(t, tt.wantErr, err)
		require.Equal(t, tt.want, u)
	}
}
