---
page_title: "gdrive_file Resource - terraform-provider-gdrive"
subcategory: ""
description: |-
  Creates a new file, folder or shorcut inside an existing google drive.
  See the official API documentation https://developers.google.com/drive/api/v3/reference.
---

# Resource `gdrive_file`

Creates a new file, folder or shorcut inside an existing google drive.
See the [official API documentation](https://developers.google.com/drive/api/v3/reference).

## Example Usage

```terraform
resource "gdrive_file" "file" {
  source = "sample.txt"
  name   = "file_created_by_terraform"
}

resource "gdrive_file" "folder" {
  type = "folder"
  name = "folder_created_by_terraform"
}

resource "gdrive_file" "shortcut" {
  type      = "shortcut"
  name      = "shortcut_created_by_terraform"
  target_id = gdrive_file.file.id
  parents   = [gdrive_file.folder.id]
}
```

## Schema

### Optional

- **drive_id** (String, Optional) The ID of the shared drive the file resides in.
- **md5_checksum_for_diff** (String, Optional) For finding file source diff.  
**Please do not set value.**
- **name** (String, Optional) The name of the file, folder or shortcut.
- **parents** (List of String, Optional) The IDs of the parent folders which contain the file, folder or shortcut.
- **source** (String, Optional) The file source path.
- **target_id** (String, Optional) The ID of shortcut target of the the file or folder.
- **type** (String, Optional) The type of google drive system files.  
folder: google drive folder  
shortcut: google drive shortcut

### Read-only

- **id** (String, Read-only) The ID of the file, folder or shortcut.  
ID is found on url as follow:  
https://drive.google.com/file/d/{ID}  
https://drive.google.com/drive/u/0/folders/{ID}
- **md5_checksum** (String, Read-only) The MD5 checksum for the content of the file.  
This is only applicable to files with binary content in Google Drive.
- **mime_type** (String, Read-only) The MIME type of the file, folder or shortcut.
- **real_id** (String, Read-only) The ID of the file, folder or shortcut.  
This is same value as `id` and created for compatibility with data source.
- **real_parents** (List of String, Read-only) The IDs of the parent folders which contain the file, folder or shortcut.  
This value is got from API and may differ from `parents`.


