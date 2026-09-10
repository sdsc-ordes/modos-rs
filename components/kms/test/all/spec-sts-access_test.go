//go:build test && integration

package all

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	. "github.com/onsi/ginkgo/v2"
	mdJwt "github.com/sdsc-ordes/modos-rs/components/kms/internal/jwt"
	mdS3 "github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/s3"
	mdSt "github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"
	"github.com/stretchr/testify/require"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/auth"
)

// Run with
//
// ```bash
//
//	just quitsh exec-target \
//		--log-level debug
//		-K "test.showTestLog: true" \
//		-K 'test.testArgs: [ "-ginkgo.label-filter=storage-access" ]'
//		"kms::test-integration"
//
// ```
var _ = Describe("S3", func() {
	var t testing.TB

	BeforeEach(func() {
		t = GinkgoTB()
	})

	Describe("requesting credentials (STS)", Label("storage-access"), func() {
		It("should give read access to bucket-*", Label("read"), func() {
			tCtx := NewTestContext(t)
			defer tCtx.Close(t)

			_, signedTokenW := CreateToken(
				t, tCtx, mdSt.BucketPermissions{
					mdSt.BucketPermission{
						Path:        "bucket-a",
						Permissions: []mdSt.Permission{mdSt.PermissionRead},
					},
					mdSt.BucketPermission{
						Path:        "bucket-b",
						Permissions: []mdSt.Permission{mdSt.PermissionRead},
					},
				})

			// Test writing.
			cl := mdJwt.NewClaims(&tCtx.Cfg.OIDC.ClaimBucketPermissions)
			err := auth.ValidateJWT(tCtx.Ctx, tCtx.JWTVerifier, signedTokenW, nil, cl)
			require.NoError(t, err)

			err = testStorageAccess(t, tCtx, cl)
			require.NoError(t, err)
		})

		It("should give write access to bucket-*", Label("write"), func() {
			tCtx := NewTestContext(t)
			defer tCtx.Close(t)

			_, signedTokenW := CreateToken(
				t, tCtx, mdSt.BucketPermissions{
					mdSt.BucketPermission{
						Path:        "bucket-a",
						Permissions: []mdSt.Permission{mdSt.PermissionWrite},
					},
					mdSt.BucketPermission{
						Path:        "bucket-b",
						Permissions: []mdSt.Permission{mdSt.PermissionWrite},
					},
				})

			// Test writing.
			cl := mdJwt.NewClaims(&tCtx.Cfg.OIDC.ClaimBucketPermissions)
			err := auth.ValidateJWT(tCtx.Ctx, tCtx.JWTVerifier, signedTokenW, nil, cl)
			require.NoError(t, err)

			err = testStorageAccess(t, tCtx, cl)
			require.NoError(t, err)
		})
	})
})

func testStorageAccess(
	t testing.TB,
	tCtx *TestContext,
	cl *mdJwt.Claims,
) error {
	ctx := tCtx.Ctx

	creds, err := tCtx.Storage.NewCredentials(
		ctx,
		cl.BucketPermissions,
		1*time.Hour,
	)
	require.NoError(t, err)

	s3cred, ok := creds.(*mdS3.S3Credentials)
	require.Equal(t, creds.Type(), "s3")
	require.True(t, ok)

	clientS3, ok := tCtx.Storage.(*mdS3.Client)
	require.True(t, ok)

	opts := func(o *s3.Options) {
		o.Credentials =
			credentials.NewStaticCredentialsProvider(
				string(s3cred.AccessKeyID),
				string(s3cred.SecretAccessKey),
				string(s3cred.SessionToken))
	}

	for _, p := range cl.BucketPermissions {
		_, err = clientS3.Client.GetObject(ctx,
			&s3.GetObjectInput{
				Bucket: aws.String(p.Bucket()),
				Key:    aws.String("test.txt"),
			},
			opts,
		)

		if p.Permissions.Contains(mdSt.PermissionRead) {
			require.NoError(t, err,
				"permissions '%v' allow read, but does not work", p)
		} else {
			require.Error(t, err,
				"permissions '%v' dont allow read, but does work", p)
		}

		file := fmt.Sprintf("test-%v.txt", uuid.New())
		body := strings.NewReader("hello")

		_, err = clientS3.Client.PutObject(ctx,
			&s3.PutObjectInput{
				Bucket: aws.String(p.Bucket()),
				Key:    aws.String(file),
				Body:   body,
			},
			opts)

		if p.Permissions.Contains(mdSt.PermissionWrite) {
			require.NoError(t, err,
				"permissions '%v' allow write, but does not work", p)
		} else {
			require.Error(t, err,
				"permissions '%v' dont allow write, but does work", p)
		}
	}

	return nil
}
