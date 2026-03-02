.PHONY: s3-init s3-validate s3-plan s3-apply s3-output s3-destroy
s3-init:
	cd ./infrastructures/terraform/s3 && terraform init
s3-validate:
	cd ./infrastructures/terraform/s3 && terraform validate
s3-plan:
	cd ./infrastructures/terraform/s3 && terraform plan
s3-apply:
	cd ./infrastructures/terraform/s3 && terraform apply
s3-output:
	cd ./infrastructures/terraform/s3 && terraform output
s3-destroy:
	cd ./infrastructures/terraform/s3 && terraform destroy
