package main

import (
	"reflect"
	"testing"
)

func TestGenerateInsertSQL(t *testing.T) {

	tests := []struct {
		name    string
		oplog   string
		want    []string
		wantErr bool
	}{
		{name: "Empty Operation",
			oplog:   "",
			want:    []string{},
			wantErr: true,
		},
		{name: "Insert Operation", oplog: `{
			"op" : "i",
			"ns" : "test.student",
			"o" : {
			  "_id" : "635b79e231d82a8ab1de863b",
			  "name" : "Selena Miller",
			  "roll_no" : 51,
			  "is_graduated" : false,
			  "date_of_birth" : "2000-01-30"
			}
		  }`,
			want: []string{
				"CREATE SCHEMA test;",
				"CREATE TABLE test.student (_id VARCHAR(255) PRIMARY KEY, date_of_birth VARCHAR(255), is_graduated BOOLEAN, name VARCHAR(255), roll_no FLOAT);",
				"INSERT INTO test.student (_id, date_of_birth, is_graduated, name, roll_no) VALUES ('635b79e231d82a8ab1de863b', '2000-01-30', false, 'Selena Miller', 51);"},
			wantErr: false},
		{name: "Update Operation -set", oplog: `{
				"op": "u",
				"ns": "test.student",
				"o": {
				   "$v": 2,
				   "diff": {
					  "u": {
						 "is_graduated": true
					  }
				   }
				},
				 "o2": {
				   "_id": "635b79e231d82a8ab1de863b"
				}
			 }`,
			want:    []string{"UPDATE test.student SET is_graduated = true WHERE _id = '635b79e231d82a8ab1de863b';"},
			wantErr: false},
		{name: "Update Operation -set with multiple option", oplog: `{
				"op": "u",
				"ns": "test.student",
				"o": {
				   "$v": 2,
				   "diff": {
					  "u": {
						 "is_graduated": true,
						 "roll_no" : 51
					  }
				   }
				},
				 "o2": {
				   "_id": "635b79e231d82a8ab1de863b"
				}
			 }`,
			want:    []string{"UPDATE test.student SET is_graduated = true, roll_no = 51 WHERE _id = '635b79e231d82a8ab1de863b';"},
			wantErr: false},
		{name: "Update Operation -unset", oplog: `{
				"op": "u",
				"ns": "test.student",
				"o": {
				   "$v": 2,
				   "diff": {
					  "d": {
						 "roll_no": false
					  }
				   }
				},
				"o2": {
				   "_id": "635b79e231d82a8ab1de863b"
				}
			 }`,
			want:    []string{"UPDATE test.student SET roll_no = NULL WHERE _id = '635b79e231d82a8ab1de863b';"},
			wantErr: false},
		{name: "Delete Operation", oplog: `{
				"op": "d",
				"ns": "test.student",
				"o": {
				  "_id": "635b79e231d82a8ab1de863b"
				}
			  }`,
			want:    []string{"DELETE FROM test.student WHERE _id = '635b79e231d82a8ab1de863b';"},
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateSQL(tt.oplog)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateSQL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}
