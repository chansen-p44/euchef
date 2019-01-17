.PHONY: all

all: build package deploy	

# Build the lambdas in the source folder
build:
	$(MAKE) -C ../src/ build

# Package the lambdas
package: template.yaml makefile
	sam package --template-file template.yaml --output-template-file packaged.yaml --s3-bucket euchef-lambda
	lambda_upload.exe

# Deploy the stack
deploy: packaged.yaml makefile
	sam deploy --template-file packaged.yaml --stack-name euchef --capabilities CAPABILITY_IAM