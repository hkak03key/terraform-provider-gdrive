---
page_title: "gdrive Provider"
subcategory: ""
description: |-
  Controls or gets information google drive.
---
# gdrive Provider
Controls or gets information google drive.  
See the [official API documentation](https://developers.google.com/drive/api/v3/reference).

## Example Usage
```terraform
terraform {
  required_providers {
    gdrive = {
      source = "hkak03key/gdrive"
    }
  }
}

provider "gdrive" {
/*
this provider authenticates using Application Default Credentials (ADC).
see https://cloud.google.com/docs/authentication/production

if using ADC as `gcloud auth application-default login`, please do as follows:
- due to OAuth 2.0 scope, set `--scopes=https://www.googleapis.com/auth/drive,https://www.googleapis.com/auth/cloud-platform` as args
- due to using google drive api, prepare the gcp project that is enabled google drive api and created OAuth 2.0 Client IDs and json file,
  and set `--client-id-file=<downloaded_client_id_file>` as args
*/
}
```
