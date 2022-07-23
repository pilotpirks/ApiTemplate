package dbsql

import (
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func Test_database_DbConnect(t *testing.T) {
	type fields struct {
		dbType string
	}

	type args struct {
		param string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"DbConnect", fields{dbType: "sqlite3"}, args{param: "./test_db.db"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &database{
				dbType: tt.fields.dbType,
			}

			_, err := d.DbConnect(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("database.DbConnect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetDatabase(t *testing.T) {

	want := &sqlite_{
		database: database{
			dbType: "sqlite3",
		},
	}

	type args struct {
		t string
	}

	tests := []struct {
		name    string
		args    args
		want    iDatabase
		wantErr bool
	}{
		{"GetDatabase", args{"sqlite3"}, want, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDatabase(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDatabase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDatabase() = %v, want %v", got, tt.want)
			}
		})
	}
}
