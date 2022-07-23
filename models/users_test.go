package models

import (
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestAddUser(t *testing.T) {
	type args struct {
		val *User
		db  *sqlx.DB
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddUser(tt.args.val, tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("AddUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
