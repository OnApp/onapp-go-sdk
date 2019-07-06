package onappgo

import (
  "fmt"
)

// Backup represent VirtualMachine backup
type Backup struct {
  AllowResizeWithoutReboot bool        `json:"allow_resize_without_reboot,bool"`
  AllowedHotMigrate        bool        `json:"allowed_hot_migrate,bool"`
  AllowedSwap              bool        `json:"allowed_swap,bool"`
  BackupServerID           int         `json:"backup_server_id,omitempty"`
  BackupSize               int         `json:"backup_size,omitempty"`
  BackupType               string      `json:"backup_type,omitempty"`
  Built                    bool        `json:"built,bool"`
  BuiltAt                  string      `json:"built_at,omitempty"`
  CreatedAt                string      `json:"created_at,omitempty"`
  DataStoreType            string      `json:"data_store_type,omitempty"`
  DiskID                   int         `json:"disk_id,omitempty"`
  ID                       int         `json:"id,omitempty"`
  Identifier               string      `json:"identifier,omitempty"`
  Initiated                string      `json:"initiated,omitempty"`
  Iqn                      string      `json:"iqn,omitempty"`
  Locked                   bool        `json:"locked,bool"`
  MarkedForDelete          bool        `json:"marked_for_delete,bool"`
  MinDiskSize              int         `json:"min_disk_size,omitempty"`
  MinMemorySize            int         `json:"min_memory_size,omitempty"`
  Note                     string      `json:"note,omitempty"`
  OperatingSystem          string      `json:"operating_system,omitempty"`
  OperatingSystemDistro    string      `json:"operating_system_distro,omitempty"`
  TargetID                 int         `json:"target_id,omitempty"`
  TargetType               string      `json:"target_type,omitempty"`
  TemplateID               int         `json:"template_id,omitempty"`
  UpdatedAt                string      `json:"updated_at,omitempty"`
  UserID                   int         `json:"user_id,omitempty"`
  VolumeID                 int         `json:"volume_id,omitempty"`
}

// Debug - print formatted Backup structure
func (obj Backup) Debug() {
  fmt.Printf("                   ID: %d\n",   obj.ID)
  fmt.Printf("               UserID: %d\n",   obj.UserID)
  fmt.Printf("           Identifier: %s\n",   obj.Identifier)
  fmt.Printf("          MinDiskSize: %dGB\n", obj.MinDiskSize)
  fmt.Printf("      OperatingSystem: %s\n",   obj.OperatingSystem)
  fmt.Printf("OperatingSystemDistro: %s\n",   obj.OperatingSystemDistro)
  fmt.Printf("            CreatedAt: %s\n",   obj.CreatedAt)
  fmt.Printf("           TemplateID: %d\n",   obj.TemplateID)
  fmt.Printf("        DataStoreType: %s\n",   obj.DataStoreType)
  fmt.Printf("           BackupSize: %d\n",   obj.BackupSize)
  fmt.Printf("            Initiated: %s\n",   obj.Initiated)
}
