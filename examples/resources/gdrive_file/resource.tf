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

