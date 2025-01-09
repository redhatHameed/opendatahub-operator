package featurestore

import (
	"context"
	"fmt"

	componentApi "github.com/opendatahub-io/opendatahub-operator/v2/apis/components/v1alpha1"
	odhtypes "github.com/opendatahub-io/opendatahub-operator/v2/pkg/controller/types"
	odhdeploy "github.com/opendatahub-io/opendatahub-operator/v2/pkg/deploy"
)

func initialize(_ context.Context, rr *odhtypes.ReconciliationRequest) error {
	rr.Manifests = append(rr.Manifests, manifestsPath())

	return nil
}

func devFlags(ctx context.Context, rr *odhtypes.ReconciliationRequest) error {
	featureStore, ok := rr.Instance.(*componentApi.FeatureStore)
	if !ok {
		return fmt.Errorf("resource instance %v is not a componentApi.FeatureStore)", rr.Instance)
	}

	if featureStore.Spec.DevFlags == nil {
		return nil
	}

	// Implement devflags support logic
	// If dev flags are set, update default manifests path
	if len(featureStore.Spec.DevFlags.Manifests) != 0 {
		manifestConfig := featureStore.Spec.DevFlags.Manifests[0]
		if err := odhdeploy.DownloadManifests(ctx, ComponentName, manifestConfig); err != nil {
			return err
		}

		if manifestConfig.SourcePath != "" {
			rr.Manifests[0].Path = odhdeploy.DefaultManifestPath
			rr.Manifests[0].ContextDir = ComponentName
			rr.Manifests[0].SourcePath = manifestConfig.SourcePath
		}
	}

	return nil
}
