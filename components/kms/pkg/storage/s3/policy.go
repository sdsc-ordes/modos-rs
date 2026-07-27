package s3

import (
	"context"
	"strings"

	"github.com/sdsc-ordes/modos-rs/components/kms/pkg/storage/types"
	"gitlab.com/data-custodian/custodian/components/lib-common/pkg/errors"
	clog "gitlab.com/data-custodian/custodian/components/lib-common/pkg/log/context"
)

const s3GetObject = "s3:GetObject"
const s3PutObject = "s3:PutObject"
const s3ListBucket = "s3:ListBucket"

type (
	policyDoc struct {
		Version   string            `json:"Version"`   //nolint: tagliatelle // AWS spec.
		Statement []policyStatement `json:"Statement"` //nolint: tagliatelle // AWS spec.

	}

	policyStatement struct {
		Effect   string   `json:"Effect"`   //nolint: tagliatelle // AWS spec.
		Action   []string `json:"Action"`   //nolint: tagliatelle // AWS spec.
		Resource []string `json:"Resource"` //nolint: tagliatelle // AWS spec.
	}
)

func toAction(ctx context.Context, perms []types.Permission) (actions []string) {
	known := make(map[types.Permission]struct{})

	for _, p := range perms {
		if _, exists := known[p]; exists {
			continue
		}
		known[p] = struct{}{}
		switch p {
		case types.PermissionRead:
			actions = append(actions, s3GetObject, s3ListBucket)
		case types.PermissionWrite:
			actions = append(actions, s3PutObject)
		default:
			clog.Error(ctx, "Cannot create action for permissions '%v'.", p)

			continue
		}
	}

	return
}

func sanitizeResourcePath(resource string) (string, error) {
	p := strings.Trim(resource, "/")
	parts := strings.Split(p, "/")

	if len(parts) == 0 {
		return "", errors.New("Bucket path '%v' results in an empty resource path -> Ignore.", p)
	} else if strings.ContainsAny(parts[0], "*?") {
		return "", errors.New("Bucket path '%v' contains '?' or '*' which is not supported yet.", p)
	}

	return p, nil
}

// NewScopedPolicy returns a scoped policy based on bucket permissions [types.BucketPermissions].
func NewScopedPolicy(
	ctx context.Context,
	permissions types.BucketPermissions,
) (*policyDoc, error) {
	doc := &policyDoc{
		Version:   "2012-10-17",
		Statement: nil,
	}

	for i := range permissions {
		perm := &permissions[i]

		resourcePath, err := sanitizeResourcePath(perm.Path)
		if err != nil {
			clog.ErrorEf(ctx, err, "Skipping permissions for '%v' due to path errors.", perm)

			continue
		}

		doc.Statement = append(doc.Statement, policyStatement{
			Effect: "Allow",
			Action: toAction(ctx, perm.Permissions),
			Resource: []string{
				"arn:aws:s3:::" + resourcePath + "/*",
			},
		})
	}

	return doc, nil
}
