package templater

import (
	"testing"
)

func TestRender(t *testing.T) {
	tests := []struct {
		name    string
		tmpl    string
		want    string
		wantErr bool
	}{
		{
			name:    "valid template",
			tmpl:    "Hello, {{ .Name }}",
			want:    "Hello, <no value>",
			wantErr: false,
		},
		{
			name:    "invalid template",
			tmpl:    "{{ if }}",
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty template",
			tmpl:    "",
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Render(tt.tmpl)
			if (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Render() = %v, want %v", got, tt.want)
			}
		})
	}
}
