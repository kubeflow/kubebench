package app

import (
	"github.com/ghodss/yaml"
	mpiJobv1alpha1 "github.com/kubeflow/mpi-operator/pkg/apis/kubeflow/v1alpha1"
	log "github.com/sirupsen/logrus"
)

type MPIJobManifestModifier struct {
	BaseManifestModifier
}

func (mm *MPIJobManifestModifier) ModifyManifest(manifest []byte, modSpec ManifestModSpec) ([]byte, error) {
	manifest, err := mm.BaseManifestModifier.ModifyManifest(manifest, modSpec)
	if err != nil {
		return nil, err
	}
	var mpiJob mpiJobv1alpha1.MPIJob
	if err := yaml.Unmarshal(manifest, &mpiJob); err != nil {
		log.Errorf("Failed to unmarshal manifest: %s", manifest)
		return nil, err
	}

	mpiJob.Spec.Template.Spec.Volumes = append(mpiJob.Spec.Template.Spec.Volumes, modSpec.Volumes...)
	for i, container := range mpiJob.Spec.Template.Spec.Containers {
		mpiJob.Spec.Template.Spec.Containers[i].VolumeMounts = append(container.VolumeMounts, modSpec.VolumeMounts...)
		mpiJob.Spec.Template.Spec.Containers[i].Env = append(container.Env, modSpec.EnvVars...)
	}

	manifest, err = yaml.Marshal(mpiJob)
	if err != nil {
		log.Errorf("Failed to create modified manifest: %s", err)
		return nil, err
	}
	return manifest, nil
}
