# keycloak-bucket-permissions-mapper

A Keycloak **OIDC protocol mapper** (Java SPI) that emits a structured
`permissions` claim from a user's group memberships:

```json
"permissions": [
  { "prefix": "mybucket",    "permissions": "write" },
  { "prefix": "otherbucket", "permissions": "read"  }
]
```

Built for Keycloak **26** (`quay.io/keycloak/keycloak:26`).

## Why a custom mapper

A group membership only encodes **one** dimension ("which group"). A per-bucket
grant needs **two** — the prefix _and_ the access level (`read`/`write`). And no
built-in mapper can emit an array-of-objects claim. This mapper reads two
attributes off each of the user's groups and assembles the JSON array.

## Group model

Create one group per `(prefix, access)` pair and give each group two attributes
(**Groups → select group → Attributes** in the admin console):

| Group name             | `prefix`      | `access` |
| ---------------------- | ------------- | -------- |
| `s3-otherbucket-read`  | `otherbucket` | `read`   |
| `s3-otherbucket-write` | `otherbucket` | `write`  |
| `s3-mybucket-write`    | `mybucket`    | `write`  |

Add each user to the groups matching the access they should have. A user in
`s3-otherbucket-read` gets `{prefix: otherbucket, permissions: read}` and
**cannot** get `write`, because `write` only comes from a group whose `access`
attribute is `write`.

Values come from the _attributes_, not the group name — so groups can be renamed
freely.

## Behaviour

- **Dedup**: if a user has several access levels for the same prefix (e.g. via
  nested group inheritance), only the strongest is emitted (`write` beats
  `read`).
- **Write implies read** (config toggle, default off): when on, a prefix granted
  `write` additionally emits a separate `{prefix, read}` entry.
- Unknown access strings are passed through as-is (ranked lowest for dedup).

## Build (Nix)

The build fetches its Maven dependencies through a fixed-output derivation, so
it needs network access on first run.

```sh
# From the repo root. First build fails on purpose with the correct hash:
nix-build -E 'with import <nixpkgs> {}; callPackage ./java-mapper {}'
# -> error: hash mismatch ... got: sha256-XXXX...
```

Copy that `got:` value into `default.nix` (`mvnHash = "sha256-XXXX...";`), then
rebuild:

```sh
nix-build -E 'with import <nixpkgs> {}; callPackage ./java-mapper {}'
ls result/share/keycloak/providers/   # keycloak-bucket-permissions-mapper.jar
```

Re-run this two-step whenever `pom.xml` dependencies change. Keep
`keycloak.version` in `pom.xml` aligned with the running server.

Integrating into a larger Nix expression:

```nix
pkgs.callPackage ./java-mapper { }
```

## Deploy

Drop the JAR into Keycloak's providers directory and let the server pick it up.

With the docker-compose stack from `../keycloak.md`, mount the built JAR:

```yaml
keycloak:
  image: quay.io/keycloak/keycloak:26
  command: start
  volumes:
    - ./result/share/keycloak/providers/keycloak-bucket-permissions-mapper.jar:/opt/keycloak/providers/keycloak-bucket-permissions-mapper.jar:ro
```

Because the command is `start` (not `--optimized`), Keycloak re-runs its build
step automatically when it detects a new provider. Restart the container after
adding/updating the JAR. (For an image built with `kc.sh build --optimized`,
rebuild the image instead.)

## Configure the mapper

1. Admin console → **Clients → your client → Client scopes → <client>-dedicated
   → Add mapper → By configuration**.
2. Pick **Bucket Permissions Mapper**.
3. Settings:
   - **Token Claim Name**: `permissions`
   - **Prefix group attribute**: `prefix`
   - **Access group attribute**: `access`
   - **Write implies read**: off (unless your resource server matches exact
     strings)
   - **Add to access token / ID token / userinfo**: as needed.

Verify with a token: decode the access token and confirm the `permissions` array
matches the user's group memberships.

## Files

```
java-mapper/
├── default.nix                        # pkgs.callPackage build (maven.buildMavenPackage)
├── pom.xml                            # Java 17, Keycloak 26 SPI deps (provided scope)
├── README.md
└── src/main/
    ├── java/ch/ordes/keycloak/BucketPermissionsMapper.java
    └── resources/META-INF/services/org.keycloak.protocol.ProtocolMapper   # SPI registration
```
