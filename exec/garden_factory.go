package exec

import (
	"crypto/sha1"
	"fmt"
	"path/filepath"

	"code.cloudfoundry.org/clock"
	"code.cloudfoundry.org/lager"

	"github.com/concourse/atc"
	"github.com/concourse/atc/dbng"
	"github.com/concourse/atc/resource"
	"github.com/concourse/atc/worker"
)

type gardenFactory struct {
	workerClient           worker.Client
	resourceFetcher        resource.Fetcher
	resourceFactory        resource.ResourceFactory
	dbResourceCacheFactory dbng.ResourceCacheFactory
}

func NewGardenFactory(
	workerClient worker.Client,
	resourceFetcher resource.Fetcher,
	resourceFactory resource.ResourceFactory,
	dbResourceCacheFactory dbng.ResourceCacheFactory,
) Factory {
	return &gardenFactory{
		workerClient:           workerClient,
		resourceFetcher:        resourceFetcher,
		resourceFactory:        resourceFactory,
		dbResourceCacheFactory: dbResourceCacheFactory,
	}
}

func (factory *gardenFactory) DependentGet(
	logger lager.Logger,
	teamID int,
	buildID int,
	planID atc.PlanID,
	stepMetadata StepMetadata,
	sourceName worker.ArtifactName,
	workerMetadata dbng.ContainerMetadata,
	delegate GetDelegate,
	resourceConfig atc.ResourceConfig,
	tags atc.Tags,
	params atc.Params,
	resourceTypes atc.VersionedResourceTypes,
) StepFactory {
	return newDependentGetStep(
		logger,
		sourceName,
		resourceConfig,
		params,
		stepMetadata,
		resource.Session{
			Metadata: workerMetadata,
		},
		tags,
		teamID,
		buildID,
		delegate,
		factory.resourceFetcher,
		resourceTypes,
		factory.dbResourceCacheFactory,
	)
}

func (factory *gardenFactory) Get(
	logger lager.Logger,
	teamID int,
	buildID int,
	plan atc.Plan,
	stepMetadata StepMetadata,
	workerMetadata dbng.ContainerMetadata,
	delegate GetDelegate,
) StepFactory {
	workerMetadata.WorkingDirectory = resource.ResourcesDir("get")

	actions := []Action{
		&GetAction{
			Type:     plan.Get.Type,
			Name:     plan.Get.Name,
			Resource: plan.Get.Resource,
			Source:   plan.Get.Source,
			Params:   plan.Get.Params,
			Version:  plan.Get.Version,
			Tags:     plan.Get.Tags,
			Outputs:  []string{plan.Get.Name},

			// TODO: can we remove these dependencies?
			delegate:          delegate,
			resourceFetcher:   factory.resourceFetcher,
			teamID:            teamID,
			containerMetadata: workerMetadata,
			resourceInstance: resource.NewResourceInstance(
				resource.ResourceType(plan.Get.Type),
				plan.Get.Version,
				plan.Get.Source,
				plan.Get.Params,
				dbng.ForBuild(buildID),
				plan.Get.VersionedResourceTypes,
				factory.dbResourceCacheFactory,
			),
			stepMetadata: stepMetadata,

			// TODO: remove after all actions are introduced
			resourceTypes: plan.Get.VersionedResourceTypes,
		},
	}

	return newActionsStep(actions, logger)
}

func (factory *gardenFactory) Put(
	logger lager.Logger,
	teamID int,
	buildID int,
	planID atc.PlanID,
	stepMetadata StepMetadata,
	workerMetadata dbng.ContainerMetadata,
	delegate PutDelegate,
	resourceConfig atc.ResourceConfig,
	tags atc.Tags,
	params atc.Params,
	resourceTypes atc.VersionedResourceTypes,
) StepFactory {
	workerMetadata.WorkingDirectory = resource.ResourcesDir("put")
	return newPutStep(
		logger,
		resourceConfig,
		params,
		stepMetadata,
		resource.Session{
			Metadata: workerMetadata,
		},
		tags,
		teamID,
		buildID,
		planID,
		delegate,
		factory.resourceFactory,
		resourceTypes,
	)
}

func (factory *gardenFactory) Task(
	logger lager.Logger,
	teamID int,
	buildID int,
	planID atc.PlanID,
	sourceName worker.ArtifactName,
	workerMetadata dbng.ContainerMetadata,
	delegate TaskDelegate,
	privileged Privileged,
	tags atc.Tags,
	configSource TaskConfigSource,
	resourceTypes atc.VersionedResourceTypes,
	inputMapping map[string]string,
	outputMapping map[string]string,
	imageArtifactName string,
	clock clock.Clock,
) StepFactory {
	workingDirectory := factory.taskWorkingDirectory(sourceName)
	workerMetadata.WorkingDirectory = workingDirectory
	return newTaskStep(
		logger,
		workerMetadata,
		tags,
		teamID,
		buildID,
		planID,
		delegate,
		privileged,
		configSource,
		factory.workerClient,
		workingDirectory,
		resourceTypes,
		inputMapping,
		outputMapping,
		imageArtifactName,
		clock,
	)
}

func (factory *gardenFactory) taskWorkingDirectory(sourceName worker.ArtifactName) string {
	sum := sha1.Sum([]byte(sourceName))
	return filepath.Join("/tmp", "build", fmt.Sprintf("%x", sum[:4]))
}
