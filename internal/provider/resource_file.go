package provider

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/api/drive/v3"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFile() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: `Creates a new file, folder or shorcut inside an existing google drive.
See the [official API documentation](https://developers.google.com/drive/api/v3/reference).`,

		CreateContext: resourceFileCreate,
		ReadContext:   resourceFileRead,
		// UpdateContext: resourceFileUpdate,
		DeleteContext: resourceFileDelete,

		Schema: map[string]*schema.Schema{
			"source": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"type", "target_id"},
				Description:   `The file source path.`,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The ID of the file, folder or shortcut.  
ID is found on url as follow:  
https://drive.google.com/file/d/{ID}  
https://drive.google.com/drive/u/0/folders/{ID}`,
			},
			"real_id": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The ID of the file, folder or shortcut.  
This is same value as ` + "`id`" + ` and created for compatibility with data source.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The name of the file, folder or shortcut.`,
			},
			"target_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source"},
				Description:   `The ID of shortcut target of the the file or folder.`,
			},
			"drive_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The ID of the shared drive the file resides in.`,
			},
			"md5_checksum": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The MD5 checksum for the content of the file.  
This is only applicable to files with binary content in Google Drive.`,
			},
			"md5_checksum_for_diff": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: `For finding file source diff.  
**Please do not set value.**`,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					source := d.Get("source").(string)
					if source == "" {
						return old == ""
					}
					local_md5, err := getFileMd5Checksum(source)
					if err != nil {
						log.Print(err)
						return old == ""
					}
					return local_md5 == old
				},
			},
			"mime_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The MIME type of the file, folder or shortcut.`,
			},
			"type": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source"},
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					enum := []string{"folder", "shortcut"}
					if val == nil {
						return
					}
					v := val.(string)
					
					for _, e := range enum {
						if v == e {
							return
						}
					}
					errs = append(errs, fmt.Errorf("%q must be any one of %q, got: %q", key, enum, v))
					return
				},
				Description: `The type of google drive system files.  
folder: google drive folder  
shortcut: google drive shortcut`,
			},
			"parents": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				ForceNew:    true,
				Description: `The IDs of the parent folders which contain the file, folder or shortcut.`,
			},
			"real_parents": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
				Description: `The IDs of the parent folders which contain the file, folder or shortcut.  
This value is got from API and may differ from ` + "`parents`" + `.`,
			},
		},
	}
}

func resourceFileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var driveFile drive.File
	var f *os.File

	driveFile.Name = d.Get("name").(string)
	driveFile.DriveId = d.Get("drive_id").(string)

	if v, ok := d.GetOk("source"); ok {
		var err error
		f, err = os.Open(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		defer f.Close()
	}

	if v, ok := d.GetOk("type"); ok {
		driveFile.MimeType = definedMimeTypes[v.(string)]
		if v.(string) == "shortcut" {
			if v, ok := d.GetOk("target_id"); ok {
				driveFile.ShortcutDetails = &drive.FileShortcutDetails{
					TargetId: v.(string),
				}
			} else {
				return diag.FromErr(fmt.Errorf("\"target_id\" must be specified if \"type\" is \"shortcut\""))
			}
		}
	} else {
		mimeType, err := getMimeType(f)
		if err != nil {
			return diag.FromErr(err)
		}
		driveFile.MimeType = mimeType
	}

	if v, ok := d.GetOk("parents"); ok {
		parentsRaw := v.([]interface{})
		parents := make([]string, len(parentsRaw))
		for i, v := range parentsRaw {
			parents[i] = v.(string)
		}
		driveFile.Parents = parents
	}

	srv, err := drive.NewService(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	createCall := srv.Files.Create(&driveFile)
	if f != nil {
		createCall.Media(f)
	}
	resDriveFile, err := createCall.
		SupportsAllDrives(true).
		Fields("id").
		Do()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resDriveFile.Id)
	d.Set("id", resDriveFile.Id)

	return resourceFileRead(ctx, d, meta)
}

func resourceFileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Get("id").(string)

	srv, err := drive.NewService(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	resDriveFile, err := srv.Files.Get(id).
		SupportsAllDrives(true).
		Fields("id, parents, mimeType, md5Checksum").
		Do()
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("real_id", resDriveFile.Id)
	d.Set("real_parents", resDriveFile.Parents)
	d.Set("mime_type", resDriveFile.MimeType)
	d.Set("md5_checksum", resDriveFile.Md5Checksum)

	d.Set("md5_checksum_for_diff", resDriveFile.Md5Checksum)

	return nil
}

// func resourceFileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
// 	// use the meta value to retrieve your client from the provider configure method
// 	// client := meta.(*apiClient)
//
// 	return diag.Errorf("not implemented")
// }

func resourceFileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Get("id").(string)

	srv, err := drive.NewService(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	err = srv.Files.Delete(id).
		SupportsAllDrives(true).
		Do()
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
