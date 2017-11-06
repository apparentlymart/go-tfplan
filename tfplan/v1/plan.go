package v1

import (
	"bytes"
	"encoding/gob"
)

type AtlasConfig struct {
	Name    string
	Include []string
	Exclude []string
}

type Backend struct {
	Type      string
	RawConfig *RawConfig

	Hash uint64
}

type BackendState struct {
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`

	Hash uint64 `json:"hash"`
}

type Config struct {
	Dir string

	Terraform       *Terraform
	Atlas           *AtlasConfig
	Modules         []*Module
	ProviderConfigs []*ProviderConfig
	Resources       []*Resource
	Variables       []*Variable
	Outputs         []*Output
}

type Diff struct {
	Modules []*ModuleDiff
}

type DiffAttrType byte

const (
	DiffAttrInput   DiffAttrType = 1
	DiffAttrOutput  DiffAttrType = 2
	DiffAttrUnknown DiffAttrType = 0
)

type InstanceDiff struct {
	Attributes     map[string]*ResourceAttrDiff
	Destroy        bool
	DestroyDeposed bool
	DestroyTainted bool

	Meta map[string]interface{}
}

type InstanceState struct {
	ID         string                 `json:"id"`
	Attributes map[string]string      `json:"attributes"`
	Meta       map[string]interface{} `json:"meta"`
	Tainted    bool                   `json:"tainted"`
}

type InterpolatedVariable interface {
	FullKey() string
}

type Module struct {
	Name      string
	Source    string
	RawConfig *RawConfig
}

type ModuleDiff struct {
	Path      []string
	Resources map[string]*InstanceDiff
	Destroy   bool
}

type ModuleState struct {
	Path []string `json:"path"`

	Outputs map[string]*OutputState `json:"outputs"`

	Resources map[string]*ResourceState `json:"resources"`

	Dependencies []string `json:"depends_on"`
}

type Output struct {
	Name        string
	DependsOn   []string
	Description string
	Sensitive   bool
	RawConfig   *RawConfig
}

type OutputState struct {
	Sensitive bool `json:"sensitive"`

	Type string `json:"type"`

	Value interface{} `json:"value"`
}

type Plan struct {
	Diff    *Diff
	Module  *Tree
	State   *State
	Vars    map[string]interface{}
	Targets []string
	Backend *BackendState
}

type ProviderConfig struct {
	Name      string
	Alias     string
	RawConfig *RawConfig
}

type Provisioner struct {
	Type      string
	RawConfig *RawConfig
	ConnInfo  *RawConfig

	When      ProvisionerWhen
	OnFailure ProvisionerOnFailure
}

type ProvisionerOnFailure int

const (
	ProvisionerOnFailureContinue ProvisionerOnFailure = 1
	ProvisionerOnFailureFail     ProvisionerOnFailure = 2
	ProvisionerOnFailureInvalid  ProvisionerOnFailure = 0
)

type ProvisionerWhen int

const (
	ProvisionerWhenCreate  ProvisionerWhen = 1
	ProvisionerWhenDestroy ProvisionerWhen = 2
	ProvisionerWhenInvalid ProvisionerWhen = 0
)

type RawConfig struct {
	Key string
	//Raw map[string]interface{}
}

func (r *RawConfig) GobDecode(b []byte) error {
	var data gobRawConfig
	err := gob.NewDecoder(bytes.NewReader(b)).Decode(&data)
	if err != nil {
		return err
	}

	r.Key = data.Key
	//r.Raw = data.Raw

	return nil
}

func (r *RawConfig) GobEncode() ([]byte, error) {
	data := gobRawConfig{
		Key: r.Key,
		//Raw: r.Raw,
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type gobRawConfig struct {
	Key string
	//Raw map[string]interface{}
}

type RemoteState struct {
	Type   string            `json:"type"`
	Config map[string]string `json:"config"`
}

type Resource struct {
	Mode         ResourceMode
	Name         string
	Type         string
	RawCount     *RawConfig
	RawConfig    *RawConfig
	Provisioners []*Provisioner
	Provider     string
	DependsOn    []string
	Lifecycle    ResourceLifecycle
}

type ResourceAttrDiff struct {
	Old         string
	New         string
	NewComputed bool
	NewRemoved  bool
	NewExtra    interface{}
	RequiresNew bool
	Sensitive   bool
	Type        DiffAttrType
}

type ResourceLifecycle struct {
	CreateBeforeDestroy bool     `mapstructure:"create_before_destroy"`
	PreventDestroy      bool     `mapstructure:"prevent_destroy"`
	IgnoreChanges       []string `mapstructure:"ignore_changes"`
}

type ResourceMode int

const (
	DataResourceMode    ResourceMode = 1
	ManagedResourceMode ResourceMode = 0
)

type ResourceState struct {
	Type         string           `json:"type"`
	Dependencies []string         `json:"depends_on"`
	Primary      *InstanceState   `json:"primary"`
	Deposed      []*InstanceState `json:"deposed"`
	Provider     string           `json:"provider"`
}

type State struct {
	Version   int            `json:"version"`
	TFVersion string         `json:"terraform_version,omitempty"`
	Serial    int64          `json:"serial"`
	Lineage   string         `json:"lineage"`
	Remote    *RemoteState   `json:"remote,omitempty"`
	Backend   *BackendState  `json:"backend,omitempty"`
	Modules   []*ModuleState `json:"modules"`
}

type Terraform struct {
	RequiredVersion string `hcl:"required_version"`
	Backend         *Backend
}

type Tree struct {
	name     string
	config   *Config
	children map[string]*Tree
	path     []string
}

func (t *Tree) GobDecode(bs []byte) error {
	// Decode the gob data
	var data treeGob
	dec := gob.NewDecoder(bytes.NewReader(bs))
	if err := dec.Decode(&data); err != nil {
		return err
	}

	// Set the fields
	t.name = data.Name
	t.config = data.Config
	t.children = data.Children
	t.path = data.Path

	return nil
}

func (t *Tree) GobEncode() ([]byte, error) {
	data := &treeGob{
		Config:   t.config,
		Children: t.children,
		Name:     t.name,
		Path:     t.path,
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// treeGob is used as a structure to Gob encode a tree.
//
// This structure is private so it can't be referenced but the fields are
// public, allowing Gob to properly encode this. When we decode this, we are
// able to turn it into a Tree.
type treeGob struct {
	Config   *Config
	Children map[string]*Tree
	Name     string
	Path     []string
}

type Variable struct {
	Name         string
	DeclaredType string `mapstructure:"type"`
	Default      interface{}
	Description  string
}
