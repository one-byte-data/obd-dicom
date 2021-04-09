package services

import (
	"testing"

	"git.onebytedata.com/odb/go-dicom/media"
	"git.onebytedata.com/odb/go-dicom/network"
	"git.onebytedata.com/odb/go-dicom/tags"
)

func Test_scu_EchoSCU(t *testing.T) {
	media.InitDict()

	type fields struct {
		destination *network.Destination
	}
	type args struct {
		timeout int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "C-Echo Success",
			fields: fields{
				destination: &network.Destination{
					Name:      "Test Destination",
					CalledAE:  "DICOM_SCP",
					CallingAE: "DICOM_SCU",
					HostName:  "cluster.k8.onebytedata.net",
					Port:      1040,
					IsCFind:   true,
					IsCMove:   true,
					IsCStore:  true,
					IsTLS:     false,
				},
			},
			args: args{
				timeout: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewSCU(tt.fields.destination)
			if err := d.EchoSCU(tt.args.timeout); (err != nil) != tt.wantErr {
				t.Errorf("scu.EchoSCU() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_scu_FindSCU(t *testing.T) {
	media.InitDict()

	type fields struct {
		destination *network.Destination
	}
	type args struct {
		Query   media.DcmObj
		timeout int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "C-Find All",
			fields: fields{
				destination: &network.Destination{
					Name:      "Test Destination",
					CalledAE:  "DICOM_SCP",
					CallingAE: "DICOM_SCU",
					HostName:  "cluster.k8.onebytedata.net",
					Port:      1040,
					IsCFind:   true,
					IsCMove:   true,
					IsCStore:  true,
					IsTLS:     false,
				},
			},
			args: args{
				Query:   media.DefaultCFindRequest(),
				timeout: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.Query.WriteString(tags.StudyDate, "20150617")
			d := NewSCU(tt.fields.destination)
			d.SetOnCFindResult(func(result media.DcmObj) {
				result.DumpTags()
			})
			_, status, err := d.FindSCU(tt.args.Query, tt.args.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("scu.FindSCU() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if status != tt.want {
				t.Errorf("scu.FindSCU() = %v, want %v", status, tt.want)
			}
		})
	}
}

func Test_scu_StoreSCU(t *testing.T) {
	media.InitDict()

	type fields struct {
		destination   *network.Destination
		onCFindResult func(result media.DcmObj)
		onCMoveResult func(result media.DcmObj)
	}
	type args struct {
		FileName string
		timeout  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "C-Store All",
			fields: fields{
				destination: &network.Destination{
					Name:      "Test Destination",
					CalledAE:  "DICOM_SCP",
					CallingAE: "DICOM_SCU",
					HostName:  "cluster.k8.onebytedata.net",
					Port:      1040,
					IsCFind:   true,
					IsCMove:   true,
					IsCStore:  true,
					IsTLS:     false,
				},
			},
			args: args{
				FileName: "../test/test2.dcm",
				timeout:  0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewSCU(tt.fields.destination)
			if err := d.StoreSCU(tt.args.FileName, tt.args.timeout); (err != nil) != tt.wantErr {
				t.Errorf("scu.StoreSCU() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
