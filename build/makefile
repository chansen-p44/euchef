.PHONY: all

all: build package deploy	

# Build the lambdas in the source folder
build:
	$(MAKE) -C ../src/ build

# Package the lambdas
package: template.yaml makefile
	sam.cmd package --region eu-west-1 --profile=ias --template-file template.yaml --output-template-file packaged.yaml --s3-bucket euchef-lambda

# Deploy the stack
deploy: packaged.yaml makefile
	sam.cmd deploy --region eu-west-1 --profile=ias --template-file packaged.yaml --stack-name euchef --capabilities CAPABILITY_IAM