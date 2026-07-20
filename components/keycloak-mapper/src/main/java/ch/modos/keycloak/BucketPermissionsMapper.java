package ch.modos.keycloak;

import org.keycloak.models.ClientSessionContext;
import org.keycloak.models.KeycloakSession;
import org.keycloak.models.ProtocolMapperModel;
import org.keycloak.models.UserModel;
import org.keycloak.models.UserSessionModel;
import org.keycloak.protocol.ProtocolMapperUtils;
import org.keycloak.protocol.oidc.mappers.AbstractOIDCProtocolMapper;
import org.keycloak.protocol.oidc.mappers.OIDCAccessTokenMapper;
import org.keycloak.protocol.oidc.mappers.OIDCAttributeMapperHelper;
import org.keycloak.protocol.oidc.mappers.OIDCIDTokenMapper;
import org.keycloak.protocol.oidc.mappers.UserInfoTokenMapper;
import org.keycloak.provider.ProviderConfigProperty;
import org.keycloak.representations.IDToken;
import org.keycloak.utils.JsonUtils;

import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * OIDC protocol mapper that emits a structured "permissions" claim built from the
 * user's group memberships:
 *
 * <pre>
 *   "${claim-name}": [
 *     { "p": "mybucket",    "perm": "write" },
 *     { "p": "otherbucket", "perm": "read"  }
 *   ]
 * </pre>
 *
 * where `p` stands for the bucket path and `pm` for the bucket permissions.
 *
 * Each group carries two attributes (configurable): one holding the bucket path and
 * one holding the permission level. A group's membership therefore encodes both dimensions
 * (which path, which permission), which a plain group-membership claim cannot.
 *
 * Model your groups/subgroups like
 *   s3-group-a  -> bucket-path=bucket-a, bucket-permission=read
 *     |-  s3-subgroup-a  -> bucket-path=bucket-a/b/c, bucket-permission=write
 * and add users to the groups matching the permission they should have.
 *
 * When a user ends up with several permission levels for the same path (e.g. via nested
 * group inheritance), only the strongest is emitted (write beats read) so the claim is
 * unambiguous. If "Write implies read" is enabled, a path granted write additionally
 * emits a separate {path, read} entry, for resource servers that match exact strings.
 */
public class BucketPermissionsMapper extends AbstractOIDCProtocolMapper
        implements OIDCAccessTokenMapper, OIDCIDTokenMapper, UserInfoTokenMapper {

    public static final String PROVIDER_ID = "oidc-bucket-permissions-mapper";

    private static final String PREFIX_ATTR = "path.attribute";
    private static final String ACCESS_ATTR = "permission.attribute";
    private static final String WRITE_IMPLIES_READ = "write.implies.read";

    private static final String ACCESS_READ = "read";
    private static final String ACCESS_WRITE = "write";

    private static final List<ProviderConfigProperty> CONFIG_PROPERTIES = new ArrayList<>();

    static {
        ProviderConfigProperty pathAttr = new ProviderConfigProperty();
        pathAttr.setName(PREFIX_ATTR);
        pathAttr.setLabel("Bucket path group attribute");
        pathAttr.setType(ProviderConfigProperty.STRING_TYPE);
        pathAttr.setDefaultValue("bucket-path");
        pathAttr.setHelpText("Name of the group attribute holding the bucket path.");
        CONFIG_PROPERTIES.add(pathAttr);

        ProviderConfigProperty permissionAttr = new ProviderConfigProperty();
        permissionAttr.setName(ACCESS_ATTR);
        permissionAttr.setLabel("Bucket permission group attribute");
        permissionAttr.setType(ProviderConfigProperty.STRING_TYPE);
        permissionAttr.setDefaultValue("bucket-permission");
        permissionAttr.setHelpText("Name of the group attribute holding the permission level (read/write).");
        CONFIG_PROPERTIES.add(permissionAttr);

        ProviderConfigProperty writeImpliesRead = new ProviderConfigProperty();
        writeImpliesRead.setName(WRITE_IMPLIES_READ);
        writeImpliesRead.setLabel("Write implies read");
        writeImpliesRead.setType(ProviderConfigProperty.BOOLEAN_TYPE);
        writeImpliesRead.setDefaultValue("false");
        writeImpliesRead.setHelpText(
                "If enabled, a path granted write also emits a separate {path, read} entry.");
        CONFIG_PROPERTIES.add(writeImpliesRead);

        // Standard "Token Claim Name" + "Add to ID/permission token/userinfo" config.
        OIDCAttributeMapperHelper.addTokenClaimNameConfig(CONFIG_PROPERTIES);
        OIDCAttributeMapperHelper.addIncludeInTokensConfig(CONFIG_PROPERTIES, BucketPermissionsMapper.class);
    }

    @Override
    public String getId() {
        return PROVIDER_ID;
    }

    @Override
    public String getDisplayType() {
        return "Bucket Permissions Mapper";
    }

    @Override
    public String getDisplayCategory() {
        return TOKEN_MAPPER_CATEGORY;
    }

    @Override
    public String getHelpText() {
        return "Emits a structured claim about bucket permissions: an array of `{p (path), perm (permissions)}` "
                + "objects derived from the user's group attributes.";
    }

    @Override
    public List<ProviderConfigProperty> getConfigProperties() {
        return CONFIG_PROPERTIES;
    }

    @Override
    protected void setClaim(IDToken token,
                            ProtocolMapperModel mappingModel,
                            UserSessionModel userSession,
                            KeycloakSession keycloakSession,
                            ClientSessionContext clientSessionCtx) {

        UserModel user = userSession.getUser();
        if (user == null) {
            return;
        }

        Map<String, String> config = mappingModel.getConfig();
        String pathAttr = orDefault(config.get(PREFIX_ATTR), "path");
        String permissionAttr = orDefault(config.get(ACCESS_ATTR), "permission");
        boolean writeImpliesRead = Boolean.parseBoolean(config.get(WRITE_IMPLIES_READ));

        // path -> strongest permission seen so far (write beats read). LinkedHashMap keeps
        // a stable, insertion-ordered claim which makes tokens easier to eyeball/diff.
        Map<String, String> strongest = new LinkedHashMap<>();

        user.getGroupsStream().forEach(group -> {
            String path = group.getFirstAttribute(pathAttr);
            String permission = group.getFirstAttribute(permissionAttr);
            if (path == null || permission == null) {
                return;
            }
            path = path.trim();
            permission = permission.trim().toLowerCase();
            if (path.isEmpty() || permission.isEmpty()) {
                return;
            }
            String current = strongest.get(path);
            if (current == null || rank(permission) > rank(current)) {
                strongest.put(path, permission);
            }
        });

        List<Map<String, String>> permissions = new ArrayList<>();
        for (Map.Entry<String, String> e : strongest.entrySet()) {
            String path = e.getKey();
            String permission = e.getValue();

            if (writeImpliesRead && ACCESS_WRITE.equals(permission)) {
                permissions.add(entry(path, ACCESS_READ));
            }

            permissions.add(entry(path, permission));
        }

        String claimName = mappingModel.getConfig().get(OIDCAttributeMapperHelper.TOKEN_CLAIM_NAME);
        if (claimName == null || claimName.isEmpty()) {
            return;
        }

        List<String> claimPath = JsonUtils.splitClaimPath(claimName);
        JsonUtils.mapClaim(claimPath, permissions, token.getOtherClaims(), true);
    }

    private static Map<String, String> entry(String path, String permission) {
        Map<String, String> m = new LinkedHashMap<>();
        m.put("path", path);
        m.put("permissions", permission);
        return m;
    }

    private static int rank(String permission) {
        if (ACCESS_WRITE.equals(permission)) {
            return 2;
        }
        if (ACCESS_READ.equals(permission)) {
            return 1;
        }
        return 0; // unknown permission levels are kept but rank lowest
    }

    private static String orDefault(String value, String fallback) {
        return (value == null || value.isEmpty()) ? fallback : value;
    }
}
