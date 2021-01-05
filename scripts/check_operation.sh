#!/bin/sh
set -e

if [ "$GOOGLE_APPLICATION_CREDENTIALS" ] && [ ! -f $GOOGLE_APPLICATION_CREDENTIALS ] ; then
  printf "ERROR: %s is not found.\n" $GOOGLE_APPLICATION_CREDENTIALS
  exit 1
fi

ADC_JSON=${CLOUDSDK_CONFIG:-~/.config/gcloud}/application_default_credentials.json
if [ ! "$GOOGLE_APPLICATION_CREDENTIALS" ] && [ ! -f $ADC_JSON ] ; then
  printf "WARNING: %s is not found.\ncontinue?" $ADC_JSON
  while :
  do
    read -p "(y/n) > " TEMP
      case "$TEMP" in
        [yY]*) break;;
        [nN]*) printf "exit"; exit 1;;
        *) continue;;
      esac
  done
fi

make install

TERRAFORM_RUNNING_DIR=.check_operation

if [ ! -d $TERRAFORM_RUNNING_DIR ]; then mkdir $TERRAFORM_RUNNING_DIR; fi

find $TERRAFORM_RUNNING_DIR -maxdepth 1 -type f -not -name "terraform.tfstate" -not -name "terraform.tfstate.backup" | xargs rm -rf
find examples/ -type f -name "*.tf" | awk -v d=$TERRAFORM_RUNNING_DIR '{$2=$1; sub("examples/", "", $2); gsub("/", "_", $2); $2=d"/"$2; print}' | xargs -L 1 cp
find examples/ -type f -not -name "*.tf" | xargs -L 1 -I {} sh -c "cp {} $TERRAFORM_RUNNING_DIR"

cd $TERRAFORM_RUNNING_DIR
terraform init
terraform validate
# terraform plan
# terraform apply --auto-approve
# terraform destroy --auto-approve
