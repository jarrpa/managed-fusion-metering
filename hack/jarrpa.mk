#KUBECTL ?= $(OCP_DIR)/bin/oc
#KUBECONFIG ?= $(OCP_CLUSTER_CONFIG_DIR)/auth/kubeconfig
#KUBECONFIG ?= /home/jrivera/projects/github.com/ocs-operator/eks.kubeconfig

##@ Hax

docker-rmi: ## Remove all dangling docker images
	docker rmi --force $$(docker images -a --filter=dangling=true -q)
