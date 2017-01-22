package service

import (
	"database/sql"
	"net/http"
	"reflect"
	"testing"

	"github.com/rancher/rancher-domain-validaiton-service/model"
)

func TestGetDomains(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetDomains(tt.args.w, tt.args.r)
		})
	}
}

func TestCreateDomain(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CreateDomain(tt.args.w, tt.args.r)
		})
	}
}

func TestActivateDomain(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ActivateDomain(tt.args.w, tt.args.r)
		})
	}
}

func TestDeleteDomain(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DeleteDomain(tt.args.w, tt.args.r)
		})
	}
}

func TestValidateDomian(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ValidateDomian(tt.args.w, tt.args.r)
		})
	}
}

func Test_checkErr(t *testing.T) {
	type args struct {
		errMasg error
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkErr(tt.args.errMasg)
		})
	}
}

func Test_getAccountID(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getAccountID(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAccountID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getAccountID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReturnHTTPError(t *testing.T) {
	type args struct {
		w            http.ResponseWriter
		r            *http.Request
		status       string
		httpStatus   int
		errorMessage string
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReturnHTTPError(tt.args.w, tt.args.r, tt.args.status, tt.args.httpStatus, tt.args.errorMessage)
		})
	}
}

func Test_writeError(t *testing.T) {
	type args struct {
		w        http.ResponseWriter
		svcError model.DomainValidationErr
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeError(tt.args.w, tt.args.svcError)
		})
	}
}

func TestReturnHTTPSuccess(t *testing.T) {
	type args struct {
		w            http.ResponseWriter
		r            *http.Request
		status       string
		httpStatus   int
		errorMessage string
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReturnHTTPSuccess(tt.args.w, tt.args.r, tt.args.status, tt.args.httpStatus, tt.args.errorMessage)
		})
	}
}

func Test_praseQueryResult(t *testing.T) {
	type args struct {
		query *sql.Rows
	}
	tests := []struct {
		name    string
		args    args
		want    []DomainList
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := praseQueryResult(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("praseQueryResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("praseQueryResult() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getValue(t *testing.T) {
	type args struct {
		host  string
		path  string
		token string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getValue(tt.args.host, tt.args.path, tt.args.token); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_randToken(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := randToken(); got != tt.want {
				t.Errorf("randToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decodeDomainID(t *testing.T) {
	type args struct {
		cattleid string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := decodeDomainID(tt.args.cattleid); got != tt.want {
				t.Errorf("decodeDomainID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encodeDomainID(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encodeDomainID(tt.args.id); got != tt.want {
				t.Errorf("encodeDomainID() = %v, want %v", got, tt.want)
			}
		})
	}
}
