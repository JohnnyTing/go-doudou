package es

import (
	"github.com/Jeffail/gabs/v2"
	"testing"
)

func TestNewMapping(t *testing.T) {
	const index = "team3_voice_analysis_wb"

	type args struct {
		mp MappingPayload
	}

	parsed, err := gabs.ParseJSON([]byte(`{
    "settings": {
        "refresh_interval": "60s",
        "number_of_replicas": "1",
        "number_of_shards": "15"
    },
    "mappings": {
        "` + index + `": {
            "properties": {
				"createAt": {
					"type": "date"
				}
			}
        }
    }
}`))
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				mp: MappingPayload{
					Base{
						Index: index,
						Type:  index,
					},
					[]Field{
						{
							Name: "createAt",
							Type: DATE,
						},
					},
				},
			},
			want: parsed.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMapping(tt.args.mp); got != tt.want {
				t.Errorf("NewMapping() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckTypeExists(t *testing.T) {

	const index = "team3_voice_analysis_wb"

	teardownSubTest := SetupSubTest(index, t)
	defer teardownSubTest(t)

	type args struct {
		esindex string
		estype  string
	}
	tests := []struct {
		name    string
		args    args
		wantB   bool
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				esindex: index,
				estype:  index,
			},
			wantB:   true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotB, err := CheckTypeExists(tt.args.esindex, tt.args.estype)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckTypeExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotB != tt.wantB {
				t.Errorf("CheckTypeExists() = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}

func TestCheckTypeExists2(t *testing.T) {

	const index = "team3_voice_analysis_wb"

	teardownSubTest := SetupSubTest(index, t)
	defer teardownSubTest(t)

	type args struct {
		esindex string
		estype  string
	}
	tests := []struct {
		name    string
		args    args
		wantB   bool
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				esindex: index,
				estype:  "not_exists",
			},
			wantB:   false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotB, err := CheckTypeExists(tt.args.esindex, tt.args.estype)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckTypeExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotB != tt.wantB {
				t.Errorf("CheckTypeExists() = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}

func TestPutMapping(t *testing.T) {
	const index = "team3_voice_analysis_wb"

	type args struct {
		mp MappingPayload
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				mp: MappingPayload{
					Base{
						Index: index,
						Type:  index,
					},
					[]Field{
						{
							Name: "orderPhrase",
							Type: SHORT,
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PutMapping(tt.args.mp); (err != nil) != tt.wantErr {
				t.Errorf("PutMapping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}