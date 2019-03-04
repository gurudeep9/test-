package api4

import (
	"net/http"
	"testing"
)

func TestGetOldClientLicense(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()
	Client := th.Client

	license, resp := Client.GetOldClientLicense("")
	CheckNoError(t, resp)

	if len(license["IsLicensed"]) == 0 {
		t.Fatal("license not returned correctly")
	}

	Client.Logout()

	_, resp = Client.GetOldClientLicense("")
	CheckNoError(t, resp)

	if _, err := Client.DoApiGet("/license/client", ""); err == nil || err.StatusCode != http.StatusNotImplemented {
		t.Fatal("should have errored with 501")
	}

	if _, err := Client.DoApiGet("/license/client?format=junk", ""); err == nil || err.StatusCode != http.StatusBadRequest {
		t.Fatal("should have errored with 400")
	}

	license, resp = th.SystemAdminClient.GetOldClientLicense("")
	CheckNoError(t, resp)

	if len(license["IsLicensed"]) == 0 {
		t.Fatal("license not returned correctly")
	}
}

func TestUploadLicenseFile(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()
	Client := th.Client

	ok, resp := Client.UploadLicenseFile([]byte{})
	CheckForbiddenStatus(t, resp)
	if ok {
		t.Fatal("should fail")
	}

	ok, resp = th.SystemAdminClient.UploadLicenseFile([]byte{})
	CheckBadRequestStatus(t, resp)
	if ok {
		t.Fatal("should fail")
	}
}

func TestRemoveLicenseFile(t *testing.T) {
	th := Setup().InitBasic()
	defer th.TearDown()
	Client := th.Client

	ok, resp := Client.RemoveLicenseFile()
	CheckForbiddenStatus(t, resp)
	if ok {
		t.Fatal("should fail")
	}

	ok, resp = th.SystemAdminClient.RemoveLicenseFile()
	CheckNoError(t, resp)
	if !ok {
		t.Fatal("should pass")
	}
}
