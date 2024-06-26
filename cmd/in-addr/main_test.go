package main

import (
	"net/netip"
	"reflect"
	"testing"
)

const (
	rev1921680 string = "0.168.192.in-addr.arpa."
	rev1921681 string = "1.168.192.in-addr.arpa."
	revIPv64   string = "4.3.2.1.d.c.b.a.8.b.d.0.1.0.0.2.in-addr.arpa."
	revIPv65   string = "5.3.2.1.d.c.b.a.8.b.d.0.1.0.0.2.in-addr.arpa."
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
				revIPv64,
			},
			wantErr: false,
		},
		{
			name: "ipv6 /64",
			args: args{
				prefix: prefixTwo,
			},
			want: []string{
				revIPv64,
			},
			wantErr: false,
		},
		{
			name: "ipv6 /63",
			args: args{
				prefix: prefixThree,
			},
			want: []string{
				revIPv64,
				revIPv65,
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
				rev1921680,
			},
			wantErr: false,
		},
		{
			name: "ipv4 /23",
			args: args{
				prefix: v4prefixTwo,
			},
			want: []string{
				rev1921680,
				rev1921681,
			},
			wantErr: false,
		},
		{
			name: "ipv4 /28",
			args: args{
				prefix: v4prefixThree,
			},
			want: []string{
				rev1921680,
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

func TestFatal(t *testing.T) {
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

func TestCheckPrefixes(t *testing.T) {
	v4prefixOne, _ := netip.ParsePrefix("192.168.0.1/23")
	v4prefixTwo, _ := netip.ParsePrefix("0.0.0.0")
	v6prefixOne, _ := netip.ParsePrefix("2001:db8:abcd:1234::1/63")
	v6prefixTwo, _ := netip.ParsePrefix("::")

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
			name: "v4 prefix",
			args: args{
				prefix: v4prefixOne,
			},
			want: []string{
				rev1921680,
				rev1921681,
			},
			wantErr: false,
		},
		{
			name: "v6 prefix",
			args: args{
				prefix: v6prefixOne,
			},
			want: []string{
				revIPv64,
				revIPv65,
			},
			wantErr: false,
		},
		{
			name: "v4 0/0",
			args: args{
				prefix: v4prefixTwo,
			},
			want:    []string{},
			wantErr: true,
		},
		{
			name: "v6 0/0",
			args: args{
				prefix: v6prefixTwo,
			},
			want:    []string{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkPrefixes(tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkPrefixes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkPrefixes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateRun(t *testing.T) {
	type fields struct {
		Subnet string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "v4 integration test",
			fields: fields{
				Subnet: "192.168.0.99/20",
			},
			wantErr: false,
		},
		{
			name: "v6 integration test",
			fields: fields{
				Subnet: "2001:db8:abcd:1234::1/63",
			},
			wantErr: false,
		},
		{
			name: "v4 no prefix test",
			fields: fields{
				Subnet: "192.168.0.99",
			},
			wantErr: true,
		},
		{
			name: "v6 no prefix test",
			fields: fields{
				Subnet: "2001:db8:abcd:1234::1",
			},
			wantErr: true,
		},
		{
			name: "v6 invalid prefix test",
			fields: fields{
				Subnet: "2001:db8:abcd:1234::1/",
			},
			wantErr: true,
		},
		{
			name: "v4 zero ip test",
			fields: fields{
				Subnet: "0.0.0.0",
			},
			wantErr: true,
		},
		{
			name: "v6 zero ip test",
			fields: fields{
				Subnet: "::",
			},
			wantErr: true,
		},
		{
			name: "nil",
			fields: fields{
				Subnet: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generate{
				Subnet: tt.fields.Subnet,
			}
			if err := g.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Generate.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
