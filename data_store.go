package onappgo

import (
  "context"
  "net/http"
  "fmt"
  "time"

  "github.com/digitalocean/godo"
)

const dataStoreBasePath = "settings/data_stores"

// DataStoresService is an interface for interfacing with the DataStore
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/data-stores
type DataStoresService interface {
  List(context.Context, *ListOptions) ([]DataStore, *Response, error)
  Get(context.Context, int) (*DataStore, *Response, error)
  Create(context.Context, *DataStoreCreateRequest) (*DataStore, *Response, error)
  // Delete(context.Context, int) (*Response, error)
  Delete(context.Context, int, interface{}) (*Transaction, *Response, error)
  // Edit(context.Context, int, *ListOptions) ([]DataStore, *Response, error)
}

// DataStoresServiceOp handles communication with the Data Store related methods of the
// OnApp API.
type DataStoresServiceOp struct {
  client *Client
}

var _ DataStoresService = &DataStoresServiceOp{}

type DataStore struct {
  ID                             int                            `json:"id,omitempty"`
  Label                          string                         `json:"label,omitempty"`
  Identifier                     string                         `json:"identifier,omitempty"`
  CreatedAt                      time.Time                      `json:"created_at,omitempty"`
  UpdatedAt                      time.Time                      `json:"updated_at,omitempty"`
  LocalHypervisorID              int                            `json:"local_hypervisor_id,omitempty"`
  DataStoreSize                  int                            `json:"data_store_size,omitempty"`
  ZombieDisksSize                int                            `json:"zombie_disks_size,omitempty"`
  IP                             string                         `json:"ip,omitempty"`
  DataStoreGroupID               int                            `json:"data_store_group_id,omitempty"`
  Enabled                        bool                           `json:"enabled,bool"`
  DataStoreType                  string                         `json:"data_store_type,omitempty"`
  IscsiIP                        string                         `json:"iscsi_ip,omitempty"`
  HypervisorGroupID              int                            `json:"hypervisor_group_id,omitempty"`
  VdcID                          int                            `json:"vdc_id,omitempty"`
  IntegratedStorageCacheEnabled  bool                           `json:"integrated_storage_cache_enabled,bool"`
  IntegratedStorageCacheSettings IntegratedStorageCacheSettings `json:"integrated_storage_cache_settings,omitempty"`
  AutoHealing                    bool                           `json:"auto_healing,bool"`
  IoLimits                       IoLimits                       `json:"io_limits,omitempty"`
  Epoch                          bool                           `json:"epoch,bool"`
  Default                        bool                           `json:"default,bool"`
  Usage                          int                            `json:"usage,omitempty"`
}

// DataStoreCreateRequest represents a request to create a DataStore
type DataStoreCreateRequest struct {
  Label             string `json:"label,omitempty"`
  DataStoreGroupID  string `json:"data_store_group_id,omitempty"`
  LocalHypervisorID string `json:"local_hypervisor_id,omitempty"`
  IP                string `json:"ip,omitempty"`
  Enabled           bool   `json:"enabled,bool"`
  DataStoreSize     string `json:"data_store_size,omitempty"`
  DataStoreType     string `json:"data_store_type,omitempty"`
}

type dataStoreCreateRequestRoot struct {
  DataStoreCreateRequest  *DataStoreCreateRequest  `json:"data_store"`
}

type dataStoreRoot struct {
  DataStore  *DataStore  `json:"data_store"`
}

func (d DataStoreCreateRequest) String() string {
  return godo.Stringify(d)
}

// List all Hypervisors.
func (s *DataStoresServiceOp) List(ctx context.Context, opt *ListOptions) ([]DataStore, *Response, error) {
  path := dataStoreBasePath + apiFormat
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]DataStore
  resp, err := s.client.Do(ctx, req, &out)

  if err != nil {
    return nil, resp, err
  }

  arr := make([]DataStore, len(out))
  for i := range arr {
    arr[i] = out[i]["data_store"]
  }

  return arr, resp, err
}

// Get individual DataStore.
func (s *DataStoresServiceOp) Get(ctx context.Context, id int) (*DataStore, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", dataStoreBasePath, id, apiFormat)
  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  root := new(dataStoreRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.DataStore, resp, err
}

// Create DataStore.
func (s *DataStoresServiceOp) Create(ctx context.Context, createRequest *DataStoreCreateRequest) (*DataStore, *Response, error) {
  if createRequest == nil {
    return nil, nil, godo.NewArgError("DataStore createRequest", "cannot be nil")
  }

  path := dataStoreBasePath + apiFormat
  rootRequest := &dataStoreCreateRequestRoot{
    DataStoreCreateRequest: createRequest,
  }

  req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
  if err != nil {
    return nil, nil, err
  }

  fmt.Println("\nHypervisor [Create] req: ", req)

  root := new(dataStoreRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, nil, err
  }

  return root.DataStore, resp, err
}

// Delete DataStore.
func (s *DataStoresServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Transaction, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", dataStoreBasePath, id, apiFormat)
  path, err := addOptions(path, meta)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
  if err != nil {
    return nil, nil, err
  }

  resp, err := s.client.Do(ctx, req, nil)

  opt := &ListOptions{
    PerPage : searchTransactions,
  }

  trx, resp, err := s.client.Transactions.ListByGroup(ctx, id, "DataStore", opt)

  var root *Transaction
  e := trx.Front()
  if e != nil {
    val := e.Value.(Transaction)
    root = &val
    return root, resp, err
  }

  return nil, nil, err
}

// Debug - print formatted DataStore structure
func (h DataStore) Debug() {
  fmt.Printf("           ID: %d\n", h.ID)
  fmt.Printf("        Label: %s\n", h.Label)
  fmt.Printf("   Identifier: %s\n", h.Identifier)
  fmt.Printf("           IP: %s\n", h.IP)
  fmt.Printf("DataStoreType: %s\n", h.DataStoreType)
  fmt.Printf("      Default: %T\n", h.Default)
  fmt.Printf("      Enabled: %T\n", h.Enabled)
  fmt.Printf("DataStoreSize: %d\n", h.DataStoreSize)
}