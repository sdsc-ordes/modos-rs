# Glossary

| Term              | Meaning                                                                                                                                    |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------ |
| Authentik         | An other IdP [provider](https://goauthentik.io/).                                                                                          |
| **ARN**           | Amazon Resource Name. The identifier of a resource in an IAM policy, for example `arn:aws:s3:::bucket-a/*`                                 |
| **buf**           | The tool used to lint Protocol Buffer files and to generate the Go bindings from them                                                      |
| **Claim**         | One field inside a JWT, for example `aud` (audience) or the bucket permissions claim `bps`                                                 |
| **CLI**           | Command line interface. Here `modos-cli`, the program that logs in and calls the KMS                                                       |
| **gRPC**          | The remote procedure call framework the KMS serves. Messages are Protocol Buffers over HTTP/2                                              |
| **IAM**           | Identity and Access Management. The part of the object store that decides who may do what                                                  |
| **IdP**           | Identity Provider. The service that authenticates users and issues the JWT. Here Keycloak or Authentik                                     |
| **JWKS**          | JSON Web Key Set. The set of public keys an IdP publishes so that anybody can verify its tokens                                            |
| **JWT**           | JSON Web Token. A signed token carrying the user's identity, scopes and claims                                                             |
| Keycloak          | An IdP [provider](https://keycloak.org).                                                                                                   |
| **KMS**           | Key Management Service. This component                                                                                                     |
| **OIDC**          | OpenID Connect. The identity layer on top of OAuth 2.0 that defines how the JWT is obtained and how the issuer publishes its keys          |
| **Policy**        | A JSON document listing allowed actions on resources. A _session policy_ is passed along with `AssumeRole` and narrows the session further |
| **protovalidate** | The library that enforces the validation rules declared in the `.proto` file                                                               |
| **RPC**           | Remote Procedure Call. One method of a gRPC service, for example `GetCredential`                                                           |
| **RustFS**        | The S3 compatible object store used in the development setup                                                                               |
| **S3**            | Simple Storage Service. The object storage API used for buckets and objects                                                                |
| **Scope**         | A permission label inside the JWT. `GetCredential` requires the scope `permissions`                                                        |
| **STS**           | Security Token Service. The S3 API that issues temporary credentials through `AssumeRole`                                                  |
