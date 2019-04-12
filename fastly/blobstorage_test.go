package fastly

import "testing"

func TestClient_BlobStorages(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "blobstorages/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var bs *BlobStorage
	record(t, "blobstorages/create", func(c *Client) {
		bs, err = c.CreateBlobStorage(&CreateBlobStorageInput{
			Service:         testServiceID,
			Version:         tv.Number,
			Name:            "test-blobstorage",
			Path:            "/logs",
			AccountName:     "test",
			Container:       "fastly",
			SASToken:        "sv=2015-04-05&ss=b&srt=sco&sp=rw&se=2019-07-21T18%3A00%3A00Z&sig=3ABdLOJZosCp0o491T%2BqZGKIhafF1nlM3MzESDDD3Gg%3D",
			Period:          12,
			TimestampFormat: "%Y-%m-%dT%H:%M:%S.000",
			GzipLevel:       9,
			PublicKey:       "-----BEGIN PGP PUBLIC KEY BLOCK-----\n\nmQENBFyUD8sBCACyFnB39AuuTygseek+eA4fo0cgwva6/FSjnWq7riouQee8GgQ/\nibXTRyv4iVlwI12GswvMTIy7zNvs1R54i0qvsLr+IZ4GVGJqs6ZJnvQcqe3xPoR4\n8AnBfw90o32r/LuHf6QCJXi+AEu35koNlNAvLJ2B+KACaNB7N0EeWmqpV/1V2k9p\nlDYk+th7LcCuaFNGqKS/PrMnnMqR6VDLCjHhNx4KR79b0Twm/2qp6an3hyNRu8Gn\ndwxpf1/BUu3JWf+LqkN4Y3mbOmSUL3MaJNvyQguUzTfS0P0uGuBDHrJCVkMZCzDB\n89ag55jCPHyGeHBTd02gHMWzsg3WMBWvCsrzABEBAAG0JXRlcnJhZm9ybSAodGVz\ndCkgPHRlc3RAdGVycmFmb3JtLmNvbT6JAU4EEwEIADgWIQSHYyc6Kj9l6HzQsau6\nvFFc9jxV/wUCXJQPywIbAwULCQgHAgYVCgkICwIEFgIDAQIeAQIXgAAKCRC6vFFc\n9jxV/815CAClb32OxV7wG01yF97TzlyTl8TnvjMtoG29Mw4nSyg+mjM3b8N7iXm9\nOLX59fbDAWtBSldSZE22RXd3CvlFOG/EnKBXSjBtEqfyxYSnyOPkMPBYWGL/ApkX\nSvPYJ4LKdvipYToKFh3y9kk2gk1DcDBDyaaHvR+3rv1u3aoy7/s2EltAfDS3ZQIq\n7/cWTLJml/lleeB/Y6rPj8xqeCYhE5ahw9gsV/Mdqatl24V9Tks30iijx0Hhw+Gx\nkATUikMGr2GDVqoIRga5kXI7CzYff4rkc0Twn47fMHHHe/KY9M2yVnMHUXmAZwbG\nM1cMI/NH1DjevCKdGBLcRJlhuLPKF/anuQENBFyUD8sBCADIpd7r7GuPd6n/Ikxe\nu6h7umV6IIPoAm88xCYpTbSZiaK30Svh6Ywra9jfE2KlU9o6Y/art8ip0VJ3m07L\n4RSfSpnzqgSwdjSq5hNour2Fo/BzYhK7yaz2AzVSbe33R0+RYhb4b/6N+bKbjwGF\nftCsqVFMH+PyvYkLbvxyQrHlA9woAZaNThI1ztO5rGSnGUR8xt84eup28WIFKg0K\nUEGUcTzz+8QGAwAra+0ewPXo/AkO+8BvZjDidP417u6gpBHOJ9qYIcO9FxHeqFyu\nYrjlrxowEgXn5wO8xuNz6Vu1vhHGDHGDsRbZF8pv1d5O+0F1G7ttZ2GRRgVBZPwi\nkiyRABEBAAGJATYEGAEIACAWIQSHYyc6Kj9l6HzQsau6vFFc9jxV/wUCXJQPywIb\nDAAKCRC6vFFc9jxV/9YOCACe8qmOSnKQpQfW+PqYOqo3dt7JyweTs3FkD6NT8Zml\ndYy/vkstbTjPpX6aTvUZjkb46BVi7AOneVHpD5GBqvRsZ9iVgDYHaehmLCdKiG5L\n3Tp90NN+QY5WDbsGmsyk6+6ZMYejb4qYfweQeduOj27aavCJdLkCYMoRKfcFYI8c\nFaNmEfKKy/r1PO20NXEG6t9t05K/frHy6ZG8bCNYdpagfFVot47r9JaQqWlTNtIR\n5+zkkSq/eG9BEtRij3a6cTdQbktdBzx2KBeI0PYc1vlZR0LpuFKZqY9vlE6vTGLR\nwMfrTEOvx0NxUM3rpaCgEmuWbB1G1Hu371oyr4srrr+N\n=28dr\n-----END PGP PUBLIC KEY BLOCK-----\n",
			Format:          "%h %l %u %{now}V %{req.method}V %{req.url}V %>s %{resp.http.Content-Length}V",
			FormatVersion:   2,
			MessageType:     "classic",
			Placement:       "waf_debug",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "blobstorages/cleanup", func(c *Client) {
			c.DeleteBlobStorage(&DeleteBlobStorageInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "test-blobstorage",
			})

			c.DeleteBlobStorage(&DeleteBlobStorageInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "new-test-blobstorage",
			})
		})
	}()

	if bs.Name != "test-blobstorage" {
		t.Errorf("bad name: %q", bs.Name)
	}
	if bs.Path != "/logs" {
		t.Errorf("bad path: %q", bs.Path)
	}
	if bs.AccountName != "test" {
		t.Errorf("bad account_name: %q", bs.AccountName)
	}
	if bs.Container != "fastly" {
		t.Errorf("bad container: %q", bs.Container)
	}
	if bs.SASToken != "sv=2015-04-05&ss=b&srt=sco&sp=rw&se=2019-07-21T18%3A00%3A00Z&sig=3ABdLOJZosCp0o491T%2BqZGKIhafF1nlM3MzESDDD3Gg%3D" {
		t.Errorf("bad sas_token: %q", bs.SASToken)
	}
	if bs.Period != 12 {
		t.Errorf("bad period: %q", bs.Period)
	}
	if bs.TimestampFormat != "%Y-%m-%dT%H:%M:%S.000" {
		t.Errorf("bad timestamp_format: %q", bs.TimestampFormat)
	}
	if bs.GzipLevel != 9 {
		t.Errorf("bad gzip_level: %q", bs.GzipLevel)
	}
	if bs.PublicKey != "-----BEGIN PGP PUBLIC KEY BLOCK-----\n\nmQENBFyUD8sBCACyFnB39AuuTygseek+eA4fo0cgwva6/FSjnWq7riouQee8GgQ/\nibXTRyv4iVlwI12GswvMTIy7zNvs1R54i0qvsLr+IZ4GVGJqs6ZJnvQcqe3xPoR4\n8AnBfw90o32r/LuHf6QCJXi+AEu35koNlNAvLJ2B+KACaNB7N0EeWmqpV/1V2k9p\nlDYk+th7LcCuaFNGqKS/PrMnnMqR6VDLCjHhNx4KR79b0Twm/2qp6an3hyNRu8Gn\ndwxpf1/BUu3JWf+LqkN4Y3mbOmSUL3MaJNvyQguUzTfS0P0uGuBDHrJCVkMZCzDB\n89ag55jCPHyGeHBTd02gHMWzsg3WMBWvCsrzABEBAAG0JXRlcnJhZm9ybSAodGVz\ndCkgPHRlc3RAdGVycmFmb3JtLmNvbT6JAU4EEwEIADgWIQSHYyc6Kj9l6HzQsau6\nvFFc9jxV/wUCXJQPywIbAwULCQgHAgYVCgkICwIEFgIDAQIeAQIXgAAKCRC6vFFc\n9jxV/815CAClb32OxV7wG01yF97TzlyTl8TnvjMtoG29Mw4nSyg+mjM3b8N7iXm9\nOLX59fbDAWtBSldSZE22RXd3CvlFOG/EnKBXSjBtEqfyxYSnyOPkMPBYWGL/ApkX\nSvPYJ4LKdvipYToKFh3y9kk2gk1DcDBDyaaHvR+3rv1u3aoy7/s2EltAfDS3ZQIq\n7/cWTLJml/lleeB/Y6rPj8xqeCYhE5ahw9gsV/Mdqatl24V9Tks30iijx0Hhw+Gx\nkATUikMGr2GDVqoIRga5kXI7CzYff4rkc0Twn47fMHHHe/KY9M2yVnMHUXmAZwbG\nM1cMI/NH1DjevCKdGBLcRJlhuLPKF/anuQENBFyUD8sBCADIpd7r7GuPd6n/Ikxe\nu6h7umV6IIPoAm88xCYpTbSZiaK30Svh6Ywra9jfE2KlU9o6Y/art8ip0VJ3m07L\n4RSfSpnzqgSwdjSq5hNour2Fo/BzYhK7yaz2AzVSbe33R0+RYhb4b/6N+bKbjwGF\nftCsqVFMH+PyvYkLbvxyQrHlA9woAZaNThI1ztO5rGSnGUR8xt84eup28WIFKg0K\nUEGUcTzz+8QGAwAra+0ewPXo/AkO+8BvZjDidP417u6gpBHOJ9qYIcO9FxHeqFyu\nYrjlrxowEgXn5wO8xuNz6Vu1vhHGDHGDsRbZF8pv1d5O+0F1G7ttZ2GRRgVBZPwi\nkiyRABEBAAGJATYEGAEIACAWIQSHYyc6Kj9l6HzQsau6vFFc9jxV/wUCXJQPywIb\nDAAKCRC6vFFc9jxV/9YOCACe8qmOSnKQpQfW+PqYOqo3dt7JyweTs3FkD6NT8Zml\ndYy/vkstbTjPpX6aTvUZjkb46BVi7AOneVHpD5GBqvRsZ9iVgDYHaehmLCdKiG5L\n3Tp90NN+QY5WDbsGmsyk6+6ZMYejb4qYfweQeduOj27aavCJdLkCYMoRKfcFYI8c\nFaNmEfKKy/r1PO20NXEG6t9t05K/frHy6ZG8bCNYdpagfFVot47r9JaQqWlTNtIR\n5+zkkSq/eG9BEtRij3a6cTdQbktdBzx2KBeI0PYc1vlZR0LpuFKZqY9vlE6vTGLR\nwMfrTEOvx0NxUM3rpaCgEmuWbB1G1Hu371oyr4srrr+N\n=28dr\n-----END PGP PUBLIC KEY BLOCK-----\n" {
		t.Errorf("bad public_key: %q", bs.PublicKey)
	}
	if bs.Format != "%h %l %u %{now}V %{req.method}V %{req.url}V %>s %{resp.http.Content-Length}V" {
		t.Errorf("bad format: %q", bs.Format)
	}
	if bs.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", bs.FormatVersion)
	}
	if bs.MessageType != "classic" {
		t.Errorf("bad message_type: %q", bs.MessageType)
	}
	if bs.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", bs.Placement)
	}

	// List
	var bsl []*BlobStorage
	record(t, "blobstorages/list", func(c *Client) {
		bsl, err = c.ListBlobStorages(&ListBlobStoragesInput{
			Service: testServiceID,
			Version: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(bsl) < 1 {
		t.Errorf("bad blob storages: %v", bsl)
	}

	// Get
	var nbs *BlobStorage
	record(t, "blobstorages/get", func(c *Client) {
		nbs, err = c.GetBlobStorage(&GetBlobStorageInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-blobstorage",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if bs.Name != nbs.Name {
		t.Errorf("bad name: %q", bs.Name)
	}
	if bs.Path != nbs.Path {
		t.Errorf("bad path: %q", bs.Path)
	}
	if bs.AccountName != nbs.AccountName {
		t.Errorf("bad account_name: %q", bs.AccountName)
	}
	if bs.Container != nbs.Container {
		t.Errorf("bad container: %q", bs.Container)
	}
	if bs.SASToken != nbs.SASToken {
		t.Errorf("bad sas_token: %q", bs.SASToken)
	}
	if bs.Period != nbs.Period {
		t.Errorf("bad period: %q", bs.Period)
	}
	if bs.TimestampFormat != nbs.TimestampFormat {
		t.Errorf("bad timestamp_format: %q", bs.TimestampFormat)
	}
	if bs.GzipLevel != nbs.GzipLevel {
		t.Errorf("bad gzip_level: %q", bs.GzipLevel)
	}
	if bs.PublicKey != nbs.PublicKey {
		t.Errorf("bad public_key: %q", bs.PublicKey)
	}
	if bs.Format != nbs.Format {
		t.Errorf("bad format: %q", bs.Format)
	}
	if bs.FormatVersion != nbs.FormatVersion {
		t.Errorf("bad format_version: %q", bs.FormatVersion)
	}
	if bs.MessageType != nbs.MessageType {
		t.Errorf("bad message_type: %q", bs.MessageType)
	}
	if bs.Placement != nbs.Placement {
		t.Errorf("bad placement: %q", bs.Placement)
	}

	// Update
	var ubs *BlobStorage
	record(t, "blobstorages/update", func(c *Client) {
		ubs, err = c.UpdateBlobStorage(&UpdateBlobStorageInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-blobstorage",
			NewName: "new-test-blobstorage",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ubs.Name != "new-test-blobstorage" {
		t.Errorf("bad name: %q", ubs.Name)
	}

	// Delete
	record(t, "blobstorages/delete", func(c *Client) {
		err = c.DeleteBlobStorage(&DeleteBlobStorageInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-test-blobstorage",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListBlobStorages_validation(t *testing.T) {
	var err error
	_, err = testClient.ListBlobStorages(&ListBlobStoragesInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListBlobStorages(&ListBlobStoragesInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateBlobStorage_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateBlobStorage(&CreateBlobStorageInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateBlobStorage(&CreateBlobStorageInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetBlobStorage_validation(t *testing.T) {
	var err error
	_, err = testClient.GetBlobStorage(&GetBlobStorageInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetBlobStorage(&GetBlobStorageInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetBlobStorage(&GetBlobStorageInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateBlobStorage_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateBlobStorage(&UpdateBlobStorageInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateBlobStorage(&UpdateBlobStorageInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateBlobStorage(&UpdateBlobStorageInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteBlobStorage_validation(t *testing.T) {
	var err error
	err = testClient.DeleteBlobStorage(&DeleteBlobStorageInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteBlobStorage(&DeleteBlobStorageInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteBlobStorage(&DeleteBlobStorageInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
