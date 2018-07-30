PROMBENCH_CMD        = ./prombench
DOCKER_TAG = docker.io/sipian/prombench:v2.0.0

#Prow config has the following args set in it's configuration during deployment
#	ZONE
#	PROJECT_ID
#	CLUSTER_NAME
#
#When the start-benchmark/cancel-benchmark prow-job is created, the benchmark plugin adds the following args
#	PR_NUMBER
#	RELEASE - which release version to benchmark PR with (Default: master)

deploy:
	$(PROMBENCH_CMD) gke nodepool create -a /etc/serviceaccount/service-account.json \
		-v ZONE:${ZONE} -v PROJECT_ID:${PROJECT_ID} -v CLUSTER_NAME:${CLUSTER_NAME} -v PR_NUMBER:${PR_NUMBER} \
		-f  components/prombench/nodepools.yaml

	$(PROMBENCH_CMD) gke resource apply -a /etc/serviceaccount/service-account.json \
		-v ZONE:${ZONE} -v PROJECT_ID:${PROJECT_ID} -v CLUSTER_NAME:${CLUSTER_NAME} \
		-v PR_NUMBER:${PR_NUMBER} -v RELEASE:${RELEASE} \
		-f components/prombench/manifests/benchmark

clean:
	$(PROMBENCH_CMD) gke resource delete -a /etc/serviceaccount/service-account.json \
		-v ZONE:${ZONE} -v PROJECT_ID:${PROJECT_ID} -v CLUSTER_NAME:${CLUSTER_NAME} -v PR_NUMBER:${PR_NUMBER} \
		-f components/prombench/manifests/benchmark/1a_namespace.yaml -f components/prombench/manifests/benchmark/1c_cluster-role-binding.yaml

	$(PROMBENCH_CMD) gke nodepool delete -a /etc/serviceaccount/service-account.json \
		-v ZONE:${ZONE} -v PROJECT_ID:${PROJECT_ID} -v CLUSTER_NAME:${CLUSTER_NAME} -v PR_NUMBER:${PR_NUMBER} \
		-f components/prombench/nodepools.yaml

build:
	@vgo build -o prombench cmd/prombench/main.go

docker: build
	@docker build -t $(DOCKER_TAG) .
	@docker push $(DOCKER_TAG)

.PHONY: deploy clean build docker