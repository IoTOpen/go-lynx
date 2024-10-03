package lynx

import "testing"

func TestFunction_FormatValue(t *testing.T) {
	type fields struct {
		ID             int64
		Type           string
		InstallationID int64
		Meta           Meta
		ProtectedMeta  Meta
		Created        int64
		Updated        int64
	}
	type args struct {
		value    float64
		topicKey string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "FormatMeta",
			fields: fields{
				Meta: map[string]string{
					"format": "%.1f°C",
				},
			},
			args: args{
				value:    21.22222,
				topicKey: "",
			},
			want: "21.2°C",
		},
		{
			name: "UnitMeta",
			fields: fields{
				Meta: map[string]string{
					"unit": "°C",
				},
			},
			args: args{
				value:    21.222,
				topicKey: "",
			},
			want: "21.222000°C",
		},
		{
			name: "StateOn",
			fields: fields{
				Meta: map[string]string{
					"state_on":  "1",
					"state_off": "0",
				},
			},
			args: args{
				value:    1,
				topicKey: "",
			},
			want: "on",
		},
		{
			name: "StateOff",
			fields: fields{
				Meta: map[string]string{
					"state_on":  "1",
					"state_off": "0",
					"text_off":  "turned off",
				},
			},
			args: args{
				value:    0,
				topicKey: "",
			},
			want: "turned off",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Function{
				ID:             tt.fields.ID,
				Type:           tt.fields.Type,
				InstallationID: tt.fields.InstallationID,
				Meta:           tt.fields.Meta,
				ProtectedMeta:  tt.fields.ProtectedMeta,
				Created:        tt.fields.Created,
				Updated:        tt.fields.Updated,
			}
			if got := f.FormatValue(tt.args.value, tt.args.topicKey); got != tt.want {
				t.Errorf("FormatValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
