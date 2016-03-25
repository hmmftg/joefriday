package platform

import "testing"

func TestSerializeDeserializeKernel(t *testing.T) {
	k, err := GetKernel()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	p := k.SerializeFlat()
	kD := DeserializeKernelFlat(p)
	if k.Version != kD.Version {
		t.Errorf("Version: got %s; want %s", kD.Version, k.Version)
	}
	if k.CompileUser != kD.CompileUser {
		t.Errorf("CompileUser: got %s; want %s", kD.CompileUser, k.CompileUser)
	}
	if k.GCC != kD.GCC {
		t.Errorf("GCC: got %s; want %s", kD.GCC, k.GCC)
	}
	if k.OSGCC != kD.OSGCC {
		t.Errorf("Version: got %s; want %s", kD.OSGCC, k.OSGCC)
	}
	if k.Type != kD.Type {
		t.Errorf("Version: got %s; want %s", kD.Type, k.Type)
	}
	if k.CompileDate != kD.CompileDate {
		t.Errorf("CompileDate: got %s; want %s", kD.CompileDate, k.CompileDate)
	}
	if k.Arch != kD.Arch {
		t.Errorf("Arch: got %s; want %s", kD.Arch, k.Arch)
	}
}

func TestSerializeDeserializeRelease(t *testing.T) {
	r, err := GetRelease()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	p := r.SerializeFlat()
	rD := DeserializeReleaseFlat(p)
	if r.ID != rD.ID {
		t.Errorf("ID: got %s; want %s", rD.ID, r.ID)
	}
	if r.IDLike != rD.IDLike {
		t.Errorf("IDLike: got %s; want %s", rD.IDLike, r.IDLike)
	}
	if r.PrettyName != rD.PrettyName {
		t.Errorf("PrettyName: got %s; want %s", rD.PrettyName, r.PrettyName)
	}
	if r.Version != rD.Version {
		t.Errorf("Version: got %s; want %s", rD.Version, r.Version)
	}
	if r.VersionID != rD.VersionID {
		t.Errorf("VersionID: got %s; want %s", rD.VersionID, r.VersionID)
	}
	if r.HomeURL != rD.HomeURL {
		t.Errorf("HomeURL: got %s; want %s", rD.HomeURL, r.HomeURL)
	}
	if r.BugReportURL != rD.BugReportURL {
		t.Errorf("BugReportURL: got %s; want %s", rD.BugReportURL, r.BugReportURL)
	}
}
