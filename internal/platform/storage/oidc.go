package storage

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"go.uber.org/zap"
)

// oidcTokenKey is the context key for the Vercel OIDC token.
type oidcTokenKey struct{}

// WithOIDCToken stores the OIDC token in the request context so the
// S3 credential provider can retrieve it per-request.
func WithOIDCToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, oidcTokenKey{}, token)
}

// OIDCTokenFromContext retrieves the OIDC token stored by WithOIDCToken.
func OIDCTokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(oidcTokenKey{}).(string)
	return token, ok
}

// oidcCredentialProvider fetches AWS credentials via STS
// AssumeRoleWithWebIdentity using a per-request OIDC token.
type oidcCredentialProvider struct {
	roleARN   string
	stsClient *sts.Client
}

// Retrieve implements aws.CredentialsProvider.
func (p *oidcCredentialProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	token, ok := OIDCTokenFromContext(ctx)
	if !ok || token == "" {
		return aws.Credentials{}, errors.New("no OIDC token in request context; set x-vercel-oidc-token header")
	}

	log := zap.L().With(zap.String("step", "sts_assume_role_with_web_identity"))

	out, err := p.stsClient.AssumeRoleWithWebIdentity(ctx, &sts.AssumeRoleWithWebIdentityInput{
		RoleArn:          aws.String(p.roleARN),
		RoleSessionName:  aws.String("muse-api"),
		WebIdentityToken: aws.String(token),
	})
	if err != nil {
		log.Error("sts assume role with web identity failed", zap.Error(err))
		return aws.Credentials{}, err
	}

	if out.Credentials == nil {
		return aws.Credentials{}, errors.New("sts returned no credentials")
	}

	expiry := time.Now().Add(15 * time.Minute)
	if out.Credentials.Expiration != nil {
		expiry = *out.Credentials.Expiration
	}

	return aws.Credentials{
		AccessKeyID:     *out.Credentials.AccessKeyId,
		SecretAccessKey: *out.Credentials.SecretAccessKey,
		SessionToken:    *out.Credentials.SessionToken,
		CanExpire:       true,
		Expires:         expiry,
	}, nil
}

// Ensure oidcCredentialProvider implements aws.CredentialsProvider.
var _ aws.CredentialsProvider = (*oidcCredentialProvider)(nil)
