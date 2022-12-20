package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"reward/internal/util"
)

type CryptoTestSuite struct {
	suite.Suite
}

func (suite *CryptoTestSuite) SetupTest() {
	util.FS = &afero.Afero{afero.NewMemMapFs()}
}

func TestUtilTestSuite(t *testing.T) {
	suite.Run(t, new(CryptoTestSuite))
}

func (suite *CryptoTestSuite) TestEncodeRSAPrivateKeyToPEM() {
	testKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	type args struct {
		privateKey *rsa.PrivateKey
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "invalid key",
			args: args{
				privateKey: &rsa.PrivateKey{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid key",
			args: args{
				privateKey: testKey,
			},
			want:    []byte("-----BEGIN RSA PRIVATE KEY-----"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			got, err := EncodeRSAPrivateKeyToPEM(tt.args.privateKey)
			if (err != nil) != tt.wantErr {
				assert.Errorf(t,
					err,
					"EncodeRSAPrivateKeyToPEM() error = %v, wantErr %v",
					err,
					tt.wantErr)
				return
			} else if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.Contains(t, string(got), string(tt.want))
		})
	}
}

func (suite *CryptoTestSuite) TestGenerateSSHKeys() {
	type args struct {
		bitSize int
		path    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "invalid size",
			args: args{
				bitSize: 0,
				path:    "/testkey",
			},
			wantErr: true,
		},
		{
			name: "valid size",
			args: args{
				bitSize: 2048,
				path:    "/testkey",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			err := GenerateSSHKeys(tt.args.bitSize, tt.args.path)
			if (err != nil) != tt.wantErr {
				assert.Errorf(t, err, "GenerateSSHKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			content, _ := util.FS.ReadFile(tt.args.path)
			assert.Contains(t, string(content), "-----BEGIN RSA PRIVATE KEY-----")
		})
	}
}

func (suite *CryptoTestSuite) TestGenerateSSHPublicKey() {
	testKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	type args struct {
		publicKey *rsa.PublicKey
	}
	tests := []struct {
		name      string
		args      args
		want      []byte
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "invalid public key",
			args: args{
				publicKey: &rsa.PublicKey{},
			},
			want:      nil,
			wantErr:   true,
			wantPanic: true,
		},
		{
			name: "invalid public key",
			args: args{
				publicKey: &testKey.PublicKey,
			},
			want:      []byte("ssh-rsa"),
			wantErr:   false,
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				assert.Panics(t, func() {
					GenerateSSHPublicKey(tt.args.publicKey)
				})
				return
			}

			got, err := GenerateSSHPublicKey(tt.args.publicKey)
			if (err != nil) != tt.wantErr {
				assert.Errorf(t, err, "GenerateSSHPublicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.Contains(t, string(got), string(tt.want))
		})
	}
}

func (suite *CryptoTestSuite) Test_generateRSAPrivateKey() {
	type args struct {
		bitSize int
	}
	tests := []struct {
		name    string
		args    args
		want    *rsa.PrivateKey
		wantErr bool
	}{
		{
			name: "invalid size",
			args: args{
				bitSize: 0,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid size",
			args: args{
				bitSize: 2048,
			},
			want:    &rsa.PrivateKey{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			got, err := generateRSAPrivateKey(tt.args.bitSize)
			if (err != nil) != tt.wantErr {
				assert.Errorf(t, err, "generateRSAPrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.IsType(t, tt.want, got)
		})
	}
}
