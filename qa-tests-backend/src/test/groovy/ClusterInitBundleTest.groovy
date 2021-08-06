import static java.util.UUID.randomUUID

import io.stackrox.proto.api.v1.ApiTokenService

import groups.BAT
import services.BaseService
import services.ClusterInitBundleService
import services.ClusterService

import org.junit.experimental.categories.Category
import spock.lang.Shared
import spock.lang.Unroll

@Category(BAT)
class ClusterInitBundleTest extends BaseSpecification {

    @Shared
    private ApiTokenService.GenerateTokenResponse adminToken

    def setupSpec() {
        disableAuthzPlugin()
        adminToken = services.ApiTokenService.generateToken(randomUUID().toString(), "Admin")
    }

    def cleanupSpec() {
        if (adminToken != null) {
            services.ApiTokenService.revokeToken(adminToken.metadata.id)
        }
    }

    @Unroll
    def "Test that revoke cluster init bundle requires impacted clusters"() {
        BaseService.useApiToken(adminToken.token)

        when:
        "making a request for the cluster init bundle"
        def bundles = ClusterInitBundleService.getInitBundles()
        def clusterId = ClusterService.getClusterId()

        then:
        "there is a bundle for current cluster"
        def bundle = bundles.find { b -> b.impactedClustersList.find { c -> c.id == clusterId } }
        assert bundle

        when:
        "try to delete used init bundle not confirming impacted clusters"
        def response = ClusterInitBundleService.revokeInitBundle(bundle.id)

        then:
        "no bundle is revoked"
        assert response.initBundleRevokedIdsCount == 0
        and:
        "impacted cluster is listed"
        assert response.initBundleRevocationErrorsList.first().impactedClustersList*.id.contains(clusterId)
    }

    def "Test that cluster init bundle can be revoked when it has no impacted clusters"() {
        BaseService.useApiToken(adminToken.token)

        given:
        "init bundle with no impacted cluster"
        def bundle = ClusterInitBundleService.generateInintBundle("qa-test").getMeta()
        when:
        "revoke it"
        def response = ClusterInitBundleService.revokeInitBundle(bundle.id)

        then:
        "no errors"
        assert response.initBundleRevocationErrorsList.empty
        and:
        "id is revoked"
        assert response.initBundleRevokedIdsList == [bundle.id]
        assert !ClusterInitBundleService.initBundles.find { it.id == bundle.id }
    }
}
