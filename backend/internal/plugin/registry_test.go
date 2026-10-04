package plugin

import "testing"

type staticPlugin struct{ manifest Manifest }

func (p staticPlugin) Manifest() Manifest { return p.manifest }

func validManifest() Manifest {
	return Manifest{
		ID:          "example-security",
		Name:        "Example Security",
		Version:     "0.1.0",
		Category:    CategorySecurity,
		Description: "Declares a permission boundary for tests.",
		Permissions: []Permission{PermClusterRead, PermAuditRead},
	}
}

func TestRegistryRejectsOutOfBoundsPermissions(t *testing.T) {
	registry := NewRegistry()
	manifest := validManifest()
	manifest.Permissions = []Permission{"secret.read"}
	if err := registry.Register(staticPlugin{manifest}); err != ErrInvalidManifest {
		t.Fatalf("secret permission err %v", err)
	}
	manifest = validManifest()
	manifest.Permissions = []Permission{"kubernetes.exec"}
	if err := registry.Register(staticPlugin{manifest}); err != ErrInvalidManifest {
		t.Fatalf("exec permission err %v", err)
	}
	if err := registry.Register(staticPlugin{validManifest()}); err != nil {
		t.Fatal(err)
	}
	if err := registry.Register(staticPlugin{validManifest()}); err != ErrDuplicate {
		t.Fatalf("duplicate err %v", err)
	}
	if got := registry.List(); len(got) != 1 || got[0].Permissions[0] != PermClusterRead {
		t.Fatalf("list %+v", got)
	}
}

func TestValidateManifestFields(t *testing.T) {
	manifest := validManifest()
	manifest.Version = "v1"
	if err := ValidateManifest(manifest); err != ErrInvalidManifest {
		t.Fatal("expected version rejection")
	}
	manifest = validManifest()
	manifest.Category = "billing"
	if err := ValidateManifest(manifest); err != ErrInvalidManifest {
		t.Fatal("expected category rejection")
	}
}
