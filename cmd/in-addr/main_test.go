package main

import (
	"net/netip"
	"reflect"
	"testing"
)

func TestIpv6(t *testing.T) {
	prefixOne, _ := netip.ParsePrefix("2001:db8:abcd:1234::1/128")
	prefixTwo, _ := netip.ParsePrefix("2001:db8:abcd:1234::1/64")
	prefixThree, _ := netip.ParsePrefix("2001:db8:abcd:1234::1/63")
	type args struct {
		prefix netip.Prefix
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "ipv6 /128",
			args: args{
				prefix: prefixOne,
			},
			want: []string{
				"4.3.2.1.d.c.b.a.8.b.d.0.1.0.0.2.in-addr.arpa.",
			},
			wantErr: false,
		},
		{
			name: "ipv6 /64",
			args: args{
				prefix: prefixTwo,
			},
			want: []string{
				"4.3.2.1.d.c.b.a.8.b.d.0.1.0.0.2.in-addr.arpa.",
			},
			wantErr: false,
		},
		{
			name: "ipv6 /63",
			args: args{
				prefix: prefixThree,
			},
			want: []string{
				"4.3.2.1.d.c.b.a.8.b.d.0.1.0.0.2.in-addr.arpa.",
				"5.3.2.1.d.c.b.a.8.b.d.0.1.0.0.2.in-addr.arpa.",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ipv6(tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("ipv6() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ipv6() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIpv4(t *testing.T) {
	v4prefixOne, _ := netip.ParsePrefix("192.168.0.1/24")
	v4prefixTwo, _ := netip.ParsePrefix("192.168.0.1/23")
	v4prefixThree, _ := netip.ParsePrefix("192.168.0.1/28")
	type args struct {
		prefix netip.Prefix
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "ipv4 /24",
			args: args{
				prefix: v4prefixOne,
			},
			want: []string{
				"0.168.192.in-addr.arpa.",
			},
			wantErr: false,
		},
		{
			name: "ipv4 /23",
			args: args{
				prefix: v4prefixTwo,
			},
			want: []string{
				"0.168.192.in-addr.arpa.",
				"1.168.192.in-addr.arpa.",
			},
			wantErr: false,
		},
		{
			name: "ipv4 /28",
			args: args{
				prefix: v4prefixThree,
			},
			want: []string{
				"0.168.192.in-addr.arpa.",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ipv4(tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("ipv4() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ipv4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fatal(t *testing.T) {
	type args struct {
		m []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "simple text",
			args: args{
				m: []string{"this", "is", "a", "test"},
			},
			want: "[FATAL] this is a test ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fatal(tt.args.m...); got != tt.want {
				t.Errorf("fatal() = %v, want %v", got, tt.want)
			}
		})
	}
}
