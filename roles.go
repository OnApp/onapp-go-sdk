package onappgo

import (
  "context"
  "net/http"
  "fmt"

  "github.com/digitalocean/godo"
)

const roleBasePath = "roles"

// RolesService is an interface for interfacing with the Role
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/roles
type RolesService interface {
  List(context.Context, *ListOptions) ([]Role, *Response, error)
  Get(context.Context, int) (*Role, *Response, error)
  Create(context.Context, *RoleCreateRequest) (*Role, *Response, error)
  Delete(context.Context, int, interface{}) (*Response, error)
  // Edit(context.Context, int, *ListOptions) ([]Role, *Response, error)
}

// RolesServiceOp handles communication with the Roles related methods of the
// OnApp API.
type RolesServiceOp struct {
  client *Client
}

var _ RolesService = &RolesServiceOp{}

type Permission struct {
  ID         int        `json:"id,omitempty"`
  Identifier string     `json:"identifier,omitempty"`
  CreatedAt  string     `json:"created_at,omitempty"`
  UpdatedAt  string     `json:"updated_at,omitempty"`
  Label      string     `json:"label,omitempty"`
}

type Permissions struct {
  Permission Permission `json:"permission,omitempty"`
}

type Role struct {
  ID          int           `json:"id,omitempty"`
  Label       string        `json:"label,omitempty"`
  Identifier  string        `json:"identifier,omitempty"`
  CreatedAt   string        `json:"created_at,omitempty"`
  UpdatedAt   string        `json:"updated_at,omitempty"`
  UsersCount  int           `json:"users_count,omitempty"`
  System      bool          `json:"system,bool"`
  Permissions []Permissions `json:"permissions,omitempty"`
}

// RoleCreateRequest represents a request to create a Role
type RoleCreateRequest struct {
  Label         string  `json:"label,omitempty"`
  PermissionIds []int   `json:"permission_ids,omitempty"`
}

type roleCreateRequestRoot struct {
  RoleCreateRequest  *RoleCreateRequest  `json:"role"`
}

type roleRoot struct {
  Role  *Role  `json:"role"`
}

func (d RoleCreateRequest) String() string {
  return godo.Stringify(d)
}

// List all Roles.
func (s *RolesServiceOp) List(ctx context.Context, opt *ListOptions) ([]Role, *Response, error) {
  path := roleBasePath + apiFormat
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]Role
  resp, err := s.client.Do(ctx, req, &out)

  if err != nil {
    return nil, resp, err
  }

  arr := make([]Role, len(out))
  for i := range arr {
    arr[i] = out[i]["role"]
  }

  return arr, resp, err
}

// Get individual Role.
func (s *RolesServiceOp) Get(ctx context.Context, id int) (*Role, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", roleBasePath, id, apiFormat)
  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  root := new(roleRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.Role, resp, err
}

// Create Role.
func (s *RolesServiceOp) Create(ctx context.Context, createRequest *RoleCreateRequest) (*Role, *Response, error) {
  if createRequest == nil {
    return nil, nil, godo.NewArgError("Role createRequest", "cannot be nil")
  }

  path := roleBasePath + apiFormat
  rootRequest := &roleCreateRequestRoot{
    RoleCreateRequest: createRequest,
  }

  req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
  if err != nil {
    return nil, nil, err
  }
  fmt.Println("Role [Create] req: ", req)

  root := new(roleRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, nil, err
  }

  return root.Role, resp, err
}

// Delete Role.
func (s *RolesServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Response, error) {
  if id < 1 {
    return nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", roleBasePath, id, apiFormat)
  path, err := addOptions(path, meta)
  if err != nil {
    return nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
  if err != nil {
    return nil, err
  }
  fmt.Println("Role [Delete] req: ", req)

  resp, err := s.client.Do(ctx, req, nil)
  if err != nil {
    return resp, err
  }

  return resp, err
}

// Debug - print formatted Role structure
func (obj Role) Debug() {
  fmt.Printf("        ID: %d\n", obj.ID)
  fmt.Printf("     Label: %s\n", obj.Label)
  fmt.Printf("Identifier: %s\n", obj.Identifier)
  fmt.Printf("UsersCount: %d\n", obj.UsersCount)
  fmt.Printf("   Enabled: %t\n", obj.System)

  if len(obj.Permissions) > 0 {
    for i := range obj.Permissions {
      p := obj.Permissions[i].Permission
      fmt.Printf("\t\tPermission: [%d] -> ", i)
      p.Debug()
    }
  }
}

// Debug - print formatted Permission structure
func (obj Permission) Debug() {
  fmt.Printf("ID: %d,\tIdentifier: %s\n", obj.ID, obj.Label)
  fmt.Printf("ID: %d,\tLabel: %s\n", obj.ID, obj.Identifier)
}