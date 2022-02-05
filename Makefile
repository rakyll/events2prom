ecr = public.ecr.aws/q1p8v8z2
repo = events2prom
region = us-east-1

publish:
	KO_DOCKER_REPO=$(ecr)/$(repo) ko publish --bare ./cmd/events2prom
	KO_DOCKER_REPO=$(ecr)/$(repo)-example ko publish --bare ./examples/helloworld

login:
	aws ecr-public get-login-password --region $(region) | docker login --username AWS --password-stdin $(ecr)
