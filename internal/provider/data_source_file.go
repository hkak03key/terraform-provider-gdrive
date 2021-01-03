package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"google.golang.org/api/drive/v3"
)

func dataSourceFile() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: `Gets an existing file, folder or shortcut inside an google drive.  
See the [official API documentation](https://developers.google.com/drive/api/v3/reference).`,

		ReadContext: dataSourceFileRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
				Description: `The ID of the file, folder or shortcut.  
ID is found on url as follow:  
https://drive.google.com/file/d/{ID}  
https://drive.google.com/drive/u/0/folders/{ID}`,
			},
			"real_id": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The ID of the file, folder or shortcut.  
This value is got from API and may differ from ` + "`id`" + `.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the file, folder or shortcut.`,
			},
			"target_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of shortcut target of the the file or folder.`,
			},
			"drive_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the shared drive the file resides in.`,
			},
			"md5_checksum": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The MD5 checksum for the content of the file.  
This is only applicable to files with binary content in Google Drive.`,
			},
			"mime_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The MIME type of the file, folder or shortcut.`,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The type of google drive system files.  
folder: google drive folder  
shortcut: google drive shortcut`,
			},
			"parents": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: `The IDs of the parent folders which contain the file, folder or shortcut.`,
			},
			"real_parents": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
				Description: `The IDs of the parent folders which contain the file, folder or shortcut.  
This is same value as ` + "`parents`" + ` and created for compatibility with resource.`,
			},
		},
	}
}

func dataSourceFileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Get("id").(string)

	srv, err := drive.NewService(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	resDriveFile, err := srv.Files.Get(id).
		SupportsAllDrives(true).
		Fields("id, parents, mimeType, md5Checksum, name, driveId, shortcutDetails").
		Do()
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("real_id", resDriveFile.Id)
	d.Set("real_parents", resDriveFile.Parents)
	d.Set("mime_type", resDriveFile.MimeType)
	d.Set("md5_checksum", resDriveFile.Md5Checksum)

	d.Set("name", resDriveFile.Name)
	d.Set("drive_id", resDriveFile.DriveId)
	d.Set("parents", resDriveFile.Parents)
	switch resDriveFile.MimeType {
	case definedMimeTypes["folder"]:
		d.Set("type", "folder")
	case definedMimeTypes["shortcut"]:
		d.Set("type", "shortcut")
		d.Set("target_id", resDriveFile.ShortcutDetails.TargetId)
	default:
		// do nothing
	}

	d.SetId(resDriveFile.Id)

	return nil
}
