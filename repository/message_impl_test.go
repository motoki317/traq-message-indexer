package repository

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"sync"
	"testing"
)

func TestRepositoryImpl_SearchMessage(t *testing.T) {
	db := sqlx.MustConnect("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/%s?parseTime=true",
		"root",
		"password",
		"localhost",
		"traq",
	))

	type fields struct {
		db          *sqlx.DB
		lock        sync.Mutex
		channelLock sync.Mutex
	}
	type args struct {
		keywords   []string
		channelIDs []string
		userIDs    []string
		limit      int
		offset     int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Message
		wantErr bool
	}{
		// Need to load some messages beforehand
		{
			name: "basic",
			fields: fields{
				db:          db,
				lock:        sync.Mutex{},
				channelLock: sync.Mutex{},
			},
			args: args{
				keywords:   []string{"じゃんけんしよう"},
				channelIDs: nil,
				userIDs:    nil,
				limit:      5,
				offset:     0,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "channel filtering",
			fields: fields{
				db:          db,
				lock:        sync.Mutex{},
				channelLock: sync.Mutex{},
			},
			args: args{
				keywords:   []string{"じゃんけんしよう"},
				channelIDs: []string{"ec513c2c-b105-4a0e-8ccd-86b12ef307d8", "ec513c2c-b105-4a0e-8ccd-86b12ef307d8"},
				userIDs:    nil,
				limit:      5,
				offset:     0,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "user filtering",
			fields: fields{
				db:          db,
				lock:        sync.Mutex{},
				channelLock: sync.Mutex{},
			},
			args: args{
				keywords:   []string{"じゃんけんしよう"},
				channelIDs: nil,
				userIDs:    []string{"fffbbced-0fe6-4d48-b5ad-5ba607e6d74e"},
				limit:      5,
				offset:     0,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RepositoryImpl{
				db:          tt.fields.db,
				lock:        tt.fields.lock,
				channelLock: tt.fields.channelLock,
			}
			got, err := r.SearchMessage(tt.args.keywords, tt.args.channelIDs, tt.args.userIDs, tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("got len(got) == %v", len(got))
			if len(got) == 0 {
				t.Errorf("SearchMessage() got len(got) == 0, want non zero")
				return
			}
		})
	}
}
