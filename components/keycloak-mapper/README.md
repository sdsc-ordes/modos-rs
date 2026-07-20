# Custom Keycloak Permission Mapper

A Keycloak **OIDC protocol mapper** (Java SPI) that emits a structured
`<claim-name>` claim encoding the bucket paths and permissions from a user's
group memberships and group attributes:

```json
"<claim-name>": [
  { "p": "mybucket",    "perm": "write" },
  { "p": "otherbucket", "perm": "read"  }
]
```

where `p` is the bucket path and `perm` is the permission.

## Why a Custom Mapper

A group membership only encodes **one** dimension ("which group"). A per-bucket
grant needs **two** — the path _and_ the access level (`read`/`write`). And no
built-in mapper can emit an array-of-objects claim. This mapper reads two
attributes off each of the user's groups and assembles the JSON array.

## Group model

Create one group per `(prefix, access)` pair and give each group two attributes
(**Groups → select group → Attributes** in the admin console):

| Group name              | `bucket-path`   | `bucket-permission ` |
| ----------------------- | --------------- | -------------------- |
| `s3-group-a`            | `bucket-a`      | `read`               |
| ...-> `s3-subgroup-a.1` | `bucket-a/a/b/` | `write`              |
| `s3-group-b`            | `bucket-b/a/b`  | `write`              |

Add each user to the groups matching the access they should have. A user in
`s3-group-a` gets `{"p": "bucket-a", "perm": "read"}` and **cannot** get
`write`, because `write` only comes from a group whose `access` attribute is
`write`.

Values come from the _attributes_, not the group name — so groups can be renamed
freely.
