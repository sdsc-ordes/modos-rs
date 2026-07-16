package ch.ordes.keycloak;

import org.keycloak.models.ClientSessionContext;
import org.keycloak.models.KeycloakSession;
import org.keycloak.models.ProtocolMapperModel;
import org.keycloak.models.UserModel;
import org.keycloak.models.UserSessionModel;
import org.keycloak.protocol.oidc.mappers.AbstractOIDCProtocolMapper;
import org.keycloak.protocol.oidc.mappers.OIDCAccessTokenMapper;
import org.keycloak.protocol.oidc.mappers.OIDCAttributeMapperHelper;
import org.keycloak.protocol.oidc.mappers.OIDCIDTokenMapper;
import org.keycloak.protocol.oidc.mappers.UserInfoTokenMapper;
import org.keycloak.provider.ProviderConfigProperty;
import org.keycloak.representations.IDToken;

import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * OIDC protocol mapper that emits a structured "permissions" claim built from the
 * user's group memberships:
 *
 * <pre>
 *   "permissions": [
 *     { "prefix": "mybucket",    "permissions": "write" },
 *     { "prefix": "otherbucket", "permissions": "read"  }
 *   ]
 * </pre>
 *
 * Each group carries two attributes (configurable): one holding the bucket prefix and
 * one holding the access level. A group's membership therefore encodes both dimensions
 * (which prefix, which access), which a plain group-membership claim cannot.
 *
 * Model your groups as one group per (prefix, access) pair, e.g.
 *   s3-otherbucket-read   -> prefix=otherbucket, access=read
 *   s3-otherbucket-write  -> prefix=otherbucket, access=write
 * and add users to the groups matching the access they should have.
 *
 * When a user ends up with several access levels for the same prefix (e.g. via nested
 * group inheritance), only the strongest is emitted (write beats read) so the claim is
 * unambiguous. If "Write implies read" is enabled, a prefix granted write additionally
 * emits a separate {prefix, read} entry, for resource servers that match exact strings.
 */
public class BucketPermissionsMapper extends AbstractOIDCProtocolMapper
        implements OIDCAccessTokenMapper, OIDCIDTokenMapper, UserInfoTokenMapper {

    public static final String PROVIDER_ID = "oidc-bucket-permissions-mapper";

    private static final String PREFIX_ATTR = "prefix.attribute";
    private static final String ACCESS_ATTR = "access.attribute";
    private static final String WRITE_IMPLIES_READ = "write.implies.read";

    private static final String ACCESS_READ = "read";
    private static final String ACCESS_WRITE = "write";

    private static final List<ProviderConfigProperty> CONFIG_PROPERTIES = new ArrayList<>();

    static {
        ProviderConfigProperty prefixAttr = new ProviderConfigProperty();
        prefixAttr.setName(PREFIX_ATTR);
        prefixAttr.setLabel("Prefix group attribute");
        prefixAttr.setType(ProviderConfigProperty.STRING_TYPE);
        prefixAttr.setDefaultValue("prefix");
        prefixAttr.setHelpText("Name of the group attribute holding the bucket prefix.");
        CONFIG_PROPERTIES.add(prefixAttr);

        ProviderConfigProperty accessAttr = new ProviderConfigProperty();
        accessAttr.setName(ACCESS_ATTR);
        accessAttr.setLabel("Access group attribute");
        accessAttr.setType(ProviderConfigProperty.STRING_TYPE);
        accessAttr.setDefaultValue("access");
        accessAttr.setHelpText("Name of the group attribute holding the access level (read/write).");
        CONFIG_PROPERTIES.add(accessAttr);

        ProviderConfigProperty writeImpliesRead = new ProviderConfigProperty();
        writeImpliesRead.setName(WRITE_IMPLIES_READ);
        writeImpliesRead.setLabel("Write implies read");
        writeImpliesRead.setType(ProviderConfigProperty.BOOLEAN_TYPE);
        writeImpliesRead.setDefaultValue("false");
        writeImpliesRead.setHelpText(
                "If enabled, a prefix granted write also emits a separate {prefix, read} entry.");
        CONFIG_PROPERTIES.add(writeImpliesRead);

        // Standard "Token Claim Name" + "Add to ID/access token/userinfo" config.
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
        return "Emits a structured 'permissions' claim: an array of {prefix, permissions} "
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
        String prefixAttr = orDefault(config.get(PREFIX_ATTR), "prefix");
        String accessAttr = orDefault(config.get(ACCESS_ATTR), "access");
        boolean writeImpliesRead = Boolean.parseBoolean(config.get(WRITE_IMPLIES_READ));

        // prefix -> strongest access seen so far (write beats read). LinkedHashMap keeps
        // a stable, insertion-ordered claim which makes tokens easier to eyeball/diff.
        Map<String, String> strongest = new LinkedHashMap<>();

        user.getGroupsStream().forEach(group -> {
            String prefix = group.getFirstAttribute(prefixAttr);
            String access = group.getFirstAttribute(accessAttr);
            if (prefix == null || access == null) {
                return;
            }
            prefix = prefix.trim();
            access = access.trim().toLowerCase();
            if (prefix.isEmpty() || access.isEmpty()) {
                return;
            }
            String current = strongest.get(prefix);
            if (current == null || rank(access) > rank(current)) {
                strongest.put(prefix, access);
            }
        });

        List<Map<String, String>> permissions = new ArrayList<>();
        for (Map.Entry<String, String> e : strongest.entrySet()) {
            String prefix = e.getKey();
            String access = e.getValue();
            if (writeImpliesRead && ACCESS_WRITE.equals(access)) {
                permissions.add(entry(prefix, ACCESS_READ));
            }
            permissions.add(entry(prefix, access));
        }

        // mapClaim honours the configured claim name (incl. nested "a.b.c" names) and
        // serialises the List<Map> as a JSON array of objects.
        OIDCAttributeMapperHelper.mapClaim(token, mappingModel, permissions);
    }

    private static Map<String, String> entry(String prefix, String access) {
        Map<String, String> m = new LinkedHashMap<>();
        m.put("prefix", prefix);
        m.put("permissions", access);
        return m;
    }

    private static int rank(String access) {
        if (ACCESS_WRITE.equals(access)) {
            return 2;
        }
        if (ACCESS_READ.equals(access)) {
            return 1;
        }
        return 0; // unknown access levels are kept but rank lowest
    }

    private static String orDefault(String value, String fallback) {
        return (value == null || value.isEmpty()) ? fallback : value;
    }
}
