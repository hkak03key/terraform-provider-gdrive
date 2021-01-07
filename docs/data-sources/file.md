---
page_title: "gdrive_file Data Source - terraform-provider-gdrive"
subcategory: ""
description: |-
  Gets an existing file, folder or shortcut inside an google drive.See the official API documentation https://developers.google.com/drive/api/v3/reference.
---

# Data Source `gdrive_file`

Gets an existing file, folder or shortcut inside an google drive.  
See the [official API documentation](https://developers.google.com/drive/api/v3/reference).

## Example Usage

```terraform
data "gdrive_file" "file" {
  id = "1mpEgTsEZhkTrqQ0iCBDmzoOyiDraGoX6"
}

data "gdrive_file" "folder" {
  id = "1_Jjmep6x5qoS83AUaGwwgYmT4IxFpZR5"
}

data "gdrive_file" "shortcut" {
  id = "1Hv-fH7WZdpAmVm2u0NjDOzLA-bk604E4"
}
```

## Schema

### Required

- **id** (String, Required) The ID of the file, folder or shortcut.  
ID is found on url as follow:  
https://drive.google.com/file/d/{ID}  
https://drive.google.com/drive/u/0/folders/{ID}

### Read-only

- **drive_id** (String, Read-only) The ID of the shared drive the file resides in.
- **md5_checksum** (String, Read-only) The MD5 checksum for the content of the file.  
This is only applicable to files with binary content in Google Drive.
- **mime_type** (String, Read-only) The MIME type of the file, folder or shortcut.
- **name** (String, Read-only) The name of the file, folder or shortcut.
- **parents** (List of String, Read-only) The IDs of the parent folders which contain the file, folder or shortcut.
- **real_id** (String, Read-only) The ID of the file, folder or shortcut.  
This value is got from API and may differ from `id`.
- **real_parents** (List of String, Read-only) The IDs of the parent folders which contain the file, folder or shortcut.  
This is same value as `parents` and created for compatibility with resource.
- **target_id** (String, Read-only) The ID of shortcut target of the the file or folder.
- **type** (String, Read-only) The type of google drive system files.  
folder: google drive folder  
shortcut: google drive shortcut


