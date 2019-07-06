package onappgo

import (
  "context"
  "net/http"
  "fmt"

  "github.com/digitalocean/godo"
)

const diskBasePath = "settings/disks"

// DisksService is an interface for interfacing with the Disk
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/disks
type DisksService interface {
  List(context.Context, *ListOptions) ([]Disk, *Response, error)
  Get(context.Context, int) (*Disk, *Response, error)
  Create(context.Context, *DiskCreateRequest) (*Disk, *Response, error)
  // Delete(context.Context, int) (*Response, error)
  Delete(context.Context, int/*, interface{}*/) (*Transaction, *Response, error)
  // Edit(context.Context, int, *ListOptions) ([]Disk, *Response, error)
}

// DisksServiceOp handles communication with the Disk related methods of the
// OnApp API.
type DisksServiceOp struct {
  client *Client
}

var _ DisksService = &DisksServiceOp{}

// Disk - represent disk from Virtual Machine
type Disk struct {
  AddToFreebsdFstab              string                         `json:"add_to_freebsd_fstab,omitempty"`
  AddToLinuxFstab                string                         `json:"add_to_linux_fstab,omitempty"`
  Built                          bool                           `json:"built,bool"`
  BurstBw                        int                            `json:"burst_bw,omitempty"`
  BurstIops                      int                            `json:"burst_iops,omitempty"`
  CreatedAt                      string                         `json:"created_at,omitempty"`
  DataStoreID                    int                            `json:"data_store_id,omitempty"`
  DiskSize                       int                            `json:"disk_size,omitempty"`
  DiskVMNumber                   int                            `json:"disk_vm_number,omitempty"`
  FileSystem                     string                         `json:"file_system,omitempty"`
  HasAutobackups                 bool                           `json:"has_autobackups"`
  ID                             int                            `json:"id,omitempty"`
  Identifier                     string                         `json:"identifier,omitempty"`
  IntegratedStorageCacheEnabled  bool                           `json:"integrated_storage_cache_enabled,bool"`
  IntegratedStorageCacheOverride bool                           `json:"integrated_storage_cache_override,bool"`
  IntegratedStorageCacheSettings IntegratedStorageCacheSettings `json:"integrated_storage_cache_settings,omitempty"`
  IoLimits                       IoLimits                       `json:"io_limits,omitempty"`
  IoLimitsOverride               bool                           `json:"io_limits_override"`
  Iqn                            string                         `json:"iqn,omitempty"`
  IsSwap                         bool                           `json:"is_swap,bool"`
  Label                          string                         `json:"label,omitempty"`
  Locked                         bool                           `json:"locked,bool"`
  MaxBw                          int                            `json:"max_bw,omitempty"`
  MaxIops                        int                            `json:"max_iops,omitempty"`
  MinIops                        int                            `json:"min_iops,omitempty"`
  MountPoint                     string                         `json:"mount_point,omitempty"`
  Mounted                        bool                           `json:"mounted,bool"`
  OpenstackID                    int                            `json:"openstack_id,omitempty"`
  Primary                        bool                           `json:"primary,bool"`
  TemporaryVirtualMachineID      int                            `json:"temporary_virtual_machine_id,omitempty"`
  UpdatedAt                      string                         `json:"updated_at,omitempty"`
  VirtualMachineID               int                            `json:"virtual_machine_id,omitempty"`
  VolumeID                       int                            `json:"volume_id,omitempty"`
}

// DiskCreateRequest - data for creating Disk
type DiskCreateRequest struct {
  Primary           bool   `json:"primary,bool"`
  DiskSize          int    `json:"disk_size,omitempty"`
  // "ext3","ext4"
  FileSystem        string `json:"file_system,omitempty"`
  DataStoreID       int    `json:"data_store_id,omitempty"`
  Label             string `json:"label,omitempty"`
  RequireFormatDisk bool   `json:"require_format_disk,bool"`
  MountPoint        string `json:"mount_point,omitempty"`
  HotAttach         bool   `json:"hot_attach,bool"`
  MinIops           int    `json:"min_iops,omitempty"`
  Mounted           bool   `json:"mounted,bool"`

  // Additional field to determine Virtual Machine to create disk
  VirtualMachineID  int    /*`json:"virtual_machine_id,omitempty"`*/
}

type diskCreateRequestRoot struct {
  DiskCreateRequest  *DiskCreateRequest  `json:"disk"`
}

type diskRoot struct {
  Disk  *Disk  `json:"disk"`
}

func (d DiskCreateRequest) String() string {
  return godo.Stringify(d)
}

// List all Disks in the cloud.
func (s *DisksServiceOp) List(ctx context.Context, opt *ListOptions) ([]Disk, *Response, error) {
  path := diskBasePath + apiFormat
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]Disk
  resp, err := s.client.Do(ctx, req, &out)

  if err != nil {
    return nil, resp, err
  }

  arr := make([]Disk, len(out))
  for i := range arr {
    arr[i] = out[i]["disk"]
  }

  return arr, resp, err
}

// Get individual Disk.
func (s *DisksServiceOp) Get(ctx context.Context, id int) (*Disk, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", diskBasePath, id, apiFormat)

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  root := new(diskRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.Disk, resp, err
}

// Create Disk.
func (s *DisksServiceOp) Create(ctx context.Context, createRequest *DiskCreateRequest) (*Disk, *Response, error) {
  if createRequest == nil {
    return nil, nil, godo.NewArgError("createRequest", "cannot be nil")
  }

  path := fmt.Sprintf("%s/%d/disks%s", virtualMachineBasePath, createRequest.VirtualMachineID, apiFormat)

  rootRequest := &diskCreateRequestRoot{
    DiskCreateRequest : createRequest,
  }

  req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
  if err != nil {
    return nil, nil, err
  }

  fmt.Println("\nDisk [Create]  req: ", req)

  root := new(diskRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, nil, err
  }

  return root.Disk, resp, err
}

// Delete Disk.
func (s *DisksServiceOp) Delete(ctx context.Context, id int/*, meta interface{}*/) (*Transaction, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", diskBasePath, id, apiFormat)

  // path, err := addOptions(path, meta)
  // if err != nil {
  //   return nil, nil, err
  // }

  req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
  if err != nil {
    return nil, nil, err
  }

  resp, err := s.client.Do(ctx, req, nil)
  if err != nil {
    return nil, resp, err
  }

  filter := struct{
    ParentID    int
    ParentType  string
  }{
    ParentID    : id,
    ParentType  : "Disk",
  }

  return lastTransaction(ctx, s.client, filter)
  // return lastTransaction(ctx, s.client, id, "Disk")
}

// Debug - print formatted Disk structure
func (obj Disk) Debug() {
  fmt.Printf("              ID: %d\n", obj.ID)
  fmt.Printf("      Identifier: %s\n", obj.Identifier)
  fmt.Printf("VirtualMachineID: %d\n", obj.VirtualMachineID)
  fmt.Printf("     DataStoreID: %d\n", obj.DataStoreID)
  fmt.Printf("           Built: %t\n", obj.Built)
  fmt.Printf("           Label: %s\n", obj.Label)
  fmt.Printf("      FileSystem: %s\n", obj.FileSystem)
  fmt.Printf("       CreatedAt: %s\n", obj.CreatedAt)
  fmt.Printf("          Locked: %t\n", obj.Locked)
  fmt.Printf("        DiskSize: %d\n", obj.DiskSize)
  fmt.Printf("    DiskVMNumber: %d\n", obj.DiskVMNumber)
  fmt.Printf("      MountPoint: %s\n", obj.MountPoint)
  fmt.Printf("         Mounted: %t\n", obj.Mounted)
}
