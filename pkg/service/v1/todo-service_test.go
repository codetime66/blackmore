package v1

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/stelo/blackmore/pkg/api/v1"
)

func Test_linksellerServiceServer_Create(t *testing.T) {
	ctx := context.Background()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	s := NewLinksellerServiceServer(db)
	tm := time.Now().In(time.UTC)

	type args struct {
		ctx context.Context
		req *v1.CreateRequest
	}
	tests := []struct {
		name    string
		s       v1.LinksellerServiceServer
		args    args
		mock    func()
		want    *v1.CreateResponse
		wantErr bool
	}{
		{
			name: "OK",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.CreateRequest{
					Api: "v1",
					Linkseller: &v1.Linkseller{
						Person:  &v1.Person{Type: "PF", Document: "0203939"},
						Machine: &v1.Machine{Modelcode: 123, Seriesnumber: "123", Value: 123.45, Model: "123", Chip: "123"},
						Order:   &v1.Order{Ordercode: 123},
					},
				},
			},
			mock: func() {
				mock.ExpectExec("INSERT INTO linkseller").WithArgs("title", "description", tm).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want: &v1.CreateResponse{
				Api: "v1",
				Id:  1,
			},
		},
		{
			name: "Unsupported API",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.CreateRequest{
					Api: "v1000",
					Linkseller: &v1.Linkseller{
						Person:  &v1.Person{Type: "PF", Document: "0203939"},
						Machine: &v1.Machine{Modelcode: 123, Seriesnumber: "123", Value: 123.45, Model: "123", Chip: "123"},
						Order:   &v1.Order{Ordercode: 123},
					},
				},
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "Invalid Reminder field format",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.CreateRequest{
					Api: "v1",
					Linkseller: &v1.Linkseller{
						Person:  &v1.Person{Type: "PF", Document: "0203939"},
						Machine: &v1.Machine{Modelcode: 123, Seriesnumber: "123", Value: 123.45, Model: "123", Chip: "123"},
						Order:   &v1.Order{Ordercode: 123},
					},
				},
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "INSERT failed",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.CreateRequest{
					Api: "v1",
					Linkseller: &v1.Linkseller{
						Person:  &v1.Person{Type: "PF", Document: "0203939"},
						Machine: &v1.Machine{Modelcode: 123, Seriesnumber: "123", Value: 123.45, Model: "123", Chip: "123"},
						Order:   &v1.Order{Ordercode: 123},
					},
				},
			},
			mock: func() {
				mock.ExpectExec("INSERT INTO linkseller").WithArgs("title", "description", tm).
					WillReturnError(errors.New("INSERT failed"))
			},
			wantErr: true,
		},
		{
			name: "LastInsertId failed",
			s:    s,
			args: args{
				ctx: ctx,
				req: &v1.CreateRequest{
					Api: "v1",
					Linkseller: &v1.Linkseller{
						Person:  &v1.Person{Type: "PF", Document: "0203939"},
						Machine: &v1.Machine{Modelcode: 123, Seriesnumber: "123", Value: 123.45, Model: "123", Chip: "123"},
						Order:   &v1.Order{Ordercode: 123},
					},
				},
			},
			mock: func() {
				mock.ExpectExec("INSERT INTO linkseller").WithArgs("title", "description", tm).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("LastInsertId failed")))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.Create(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("linksellerServiceServer.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("linksellerServiceServer.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
