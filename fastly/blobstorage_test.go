package fastly

import (
	"testing"
)

const (
	MiB = 1048576
)

func TestClient_BlobStorages(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "blobstorages/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var bsCreateResp1, bsCreateResp2, bsCreateResp3 *BlobStorage
	record(t, "blobstorages/create", func(c *Client) {
		bsCreateResp1, err = c.CreateBlobStorage(&CreateBlobStorageInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-blobstorage",
			Path:             "/logs",
			AccountName:      "test",
			Container:        "fastly",
			SASToken:         "sv=2015-04-05&ss=b&srt=sco&sp=rw&se=2030-07-21T18%3A00%3A00Z&sig=3ABdLOJZosCp0o491T%2BqZGKIhafF1nlM3MzESDDD3Gg%3D",
			Period:           12,
			TimestampFormat:  "%Y-%m-%dT%H:%M:%S.000",
			CompressionCodec: "snappy",
			PublicKey:        "-----BEGIN PGP PUBLIC KEY BLOCK-----\n\nmQENBFyUD8sBCACyFnB39AuuTygseek+eA4fo0cgwva6/FSjnWq7riouQee8GgQ/\nibXTRyv4iVlwI12GswvMTIy7zNvs1R54i0qvsLr+IZ4GVGJqs6ZJnvQcqe3xPoR4\n8AnBfw90o32r/LuHf6QCJXi+AEu35koNlNAvLJ2B+KACaNB7N0EeWmqpV/1V2k9p\nlDYk+th7LcCuaFNGqKS/PrMnnMqR6VDLCjHhNx4KR79b0Twm/2qp6an3hyNRu8Gn\ndwxpf1/BUu3JWf+LqkN4Y3mbOmSUL3MaJNvyQguUzTfS0P0uGuBDHrJCVkMZCzDB\n89ag55jCPHyGeHBTd02gHMWzsg3WMBWvCsrzABEBAAG0JXRlcnJhZm9ybSAodGVz\ndCkgPHRlc3RAdGVycmFmb3JtLmNvbT6JAU4EEwEIADgWIQSHYyc6Kj9l6HzQsau6\nvFFc9jxV/wUCXJQPywIbAwULCQgHAgYVCgkICwIEFgIDAQIeAQIXgAAKCRC6vFFc\n9jxV/815CAClb32OxV7wG01yF97TzlyTl8TnvjMtoG29Mw4nSyg+mjM3b8N7iXm9\nOLX59fbDAWtBSldSZE22RXd3CvlFOG/EnKBXSjBtEqfyxYSnyOPkMPBYWGL/ApkX\nSvPYJ4LKdvipYToKFh3y9kk2gk1DcDBDyaaHvR+3rv1u3aoy7/s2EltAfDS3ZQIq\n7/cWTLJml/lleeB/Y6rPj8xqeCYhE5ahw9gsV/Mdqatl24V9Tks30iijx0Hhw+Gx\nkATUikMGr2GDVqoIRga5kXI7CzYff4rkc0Twn47fMHHHe/KY9M2yVnMHUXmAZwbG\nM1cMI/NH1DjevCKdGBLcRJlhuLPKF/anuQENBFyUD8sBCADIpd7r7GuPd6n/Ikxe\nu6h7umV6IIPoAm88xCYpTbSZiaK30Svh6Ywra9jfE2KlU9o6Y/art8ip0VJ3m07L\n4RSfSpnzqgSwdjSq5hNour2Fo/BzYhK7yaz2AzVSbe33R0+RYhb4b/6N+bKbjwGF\nftCsqVFMH+PyvYkLbvxyQrHlA9woAZaNThI1ztO5rGSnGUR8xt84eup28WIFKg0K\nUEGUcTzz+8QGAwAra+0ewPXo/AkO+8BvZjDidP417u6gpBHOJ9qYIcO9FxHeqFyu\nYrjlrxowEgXn5wO8xuNz6Vu1vhHGDHGDsRbZF8pv1d5O+0F1G7ttZ2GRRgVBZPwi\nkiyRABEBAAGJATYEGAEIACAWIQSHYyc6Kj9l6HzQsau6vFFc9jxV/wUCXJQPywIb\nDAAKCRC6vFFc9jxV/9YOCACe8qmOSnKQpQfW+PqYOqo3dt7JyweTs3FkD6NT8Zml\ndYy/vkstbTjPpX6aTvUZjkb46BVi7AOneVHpD5GBqvRsZ9iVgDYHaehmLCdKiG5L\n3Tp90NN+QY5WDbsGmsyk6+6ZMYejb4qYfweQeduOj27aavCJdLkCYMoRKfcFYI8c\nFaNmEfKKy/r1PO20NXEG6t9t05K/frHy6ZG8bCNYdpagfFVot47r9JaQqWlTNtIR\n5+zkkSq/eG9BEtRij3a6cTdQbktdBzx2KBeI0PYc1vlZR0LpuFKZqY9vlE6vTGLR\nwMfrTEOvx0NxUM3rpaCgEmuWbB1G1Hu371oyr4srrr+N\n=28dr\n-----END PGP PUBLIC KEY BLOCK-----\n",
			Format:           "%h %l %u %{now}V %{req.method}V %{req.url}V %>s %{resp.http.Content-Length}V",
			FormatVersion:    2,
			MessageType:      "classic",
			Placement:        "waf_debug",
			FileMaxBytes:     MiB,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	record(t, "blobstorages/create2", func(c *Client) {
		bsCreateResp2, err = c.CreateBlobStorage(&CreateBlobStorageInput{
			ServiceID:       testServiceID,
			ServiceVersion:  tv.Number,
			Name:            "test-blobstorage-2",
			Path:            "/logs",
			AccountName:     "test",
			Container:       "fastly",
			SASToken:        "sv=2015-04-05&ss=b&srt=sco&sp=rw&se=2030-07-21T18%3A00%3A00Z&sig=3ABdLOJZosCp0o491T%2BqZGKIhafF1nlM3MzESDDD3Gg%3D",
			Period:          12,
			TimestampFormat: "%Y-%m-%dT%H:%M:%S.000",
			GzipLevel:       8,
			PublicKey:       "-----BEGIN PGP PUBLIC KEY BLOCK-----\n\nmQENBFyUD8sBCACyFnB39AuuTygseek+eA4fo0cgwva6/FSjnWq7riouQee8GgQ/\nibXTRyv4iVlwI12GswvMTIy7zNvs1R54i0qvsLr+IZ4GVGJqs6ZJnvQcqe3xPoR4\n8AnBfw90o32r/LuHf6QCJXi+AEu35koNlNAvLJ2B+KACaNB7N0EeWmqpV/1V2k9p\nlDYk+th7LcCuaFNGqKS/PrMnnMqR6VDLCjHhNx4KR79b0Twm/2qp6an3hyNRu8Gn\ndwxpf1/BUu3JWf+LqkN4Y3mbOmSUL3MaJNvyQguUzTfS0P0uGuBDHrJCVkMZCzDB\n89ag55jCPHyGeHBTd02gHMWzsg3WMBWvCsrzABEBAAG0JXRlcnJhZm9ybSAodGVz\ndCkgPHRlc3RAdGVycmFmb3JtLmNvbT6JAU4EEwEIADgWIQSHYyc6Kj9l6HzQsau6\nvFFc9jxV/wUCXJQPywIbAwULCQgHAgYVCgkICwIEFgIDAQIeAQIXgAAKCRC6vFFc\n9jxV/815CAClb32OxV7wG01yF97TzlyTl8TnvjMtoG29Mw4nSyg+mjM3b8N7iXm9\nOLX59fbDAWtBSldSZE22RXd3CvlFOG/EnKBXSjBtEqfyxYSnyOPkMPBYWGL/ApkX\nSvPYJ4LKdvipYToKFh3y9kk2gk1DcDBDyaaHvR+3rv1u3aoy7/s2EltAfDS3ZQIq\n7/cWTLJml/lleeB/Y6rPj8xqeCYhE5ahw9gsV/Mdqatl24V9Tks30iijx0Hhw+Gx\nkATUikMGr2GDVqoIRga5kXI7CzYff4rkc0Twn47fMHHHe/KY9M2yVnMHUXmAZwbG\nM1cMI/NH1DjevCKdGBLcRJlhuLPKF/anuQENBFyUD8sBCADIpd7r7GuPd6n/Ikxe\nu6h7umV6IIPoAm88xCYpTbSZiaK30Svh6Ywra9jfE2KlU9o6Y/art8ip0VJ3m07L\n4RSfSpnzqgSwdjSq5hNour2Fo/BzYhK7yaz2AzVSbe33R0+RYhb4b/6N+bKbjwGF\nftCsqVFMH+PyvYkLbvxyQrHlA9woAZaNThI1ztO5rGSnGUR8xt84eup28WIFKg0K\nUEGUcTzz+8QGAwAra+0ewPXo/AkO+8BvZjDidP417u6gpBHOJ9qYIcO9FxHeqFyu\nYrjlrxowEgXn5wO8xuNz6Vu1vhHGDHGDsRbZF8pv1d5O+0F1G7ttZ2GRRgVBZPwi\nkiyRABEBAAGJATYEGAEIACAWIQSHYyc6Kj9l6HzQsau6vFFc9jxV/wUCXJQPywIb\nDAAKCRC6vFFc9jxV/9YOCACe8qmOSnKQpQfW+PqYOqo3dt7JyweTs3FkD6NT8Zml\ndYy/vkstbTjPpX6aTvUZjkb46BVi7AOneVHpD5GBqvRsZ9iVgDYHaehmLCdKiG5L\n3Tp90NN+QY5WDbsGmsyk6+6ZMYejb4qYfweQeduOj27aavCJdLkCYMoRKfcFYI8c\nFaNmEfKKy/r1PO20NXEG6t9t05K/frHy6ZG8bCNYdpagfFVot47r9JaQqWlTNtIR\n5+zkkSq/eG9BEtRij3a6cTdQbktdBzx2KBeI0PYc1vlZR0LpuFKZqY9vlE6vTGLR\nwMfrTEOvx0NxUM3rpaCgEmuWbB1G1Hu371oyr4srrr+N\n=28dr\n-----END PGP PUBLIC KEY BLOCK-----\n",
			Format:          "%h %l %u %{now}V %{req.method}V %{req.url}V %>s %{resp.http.Content-Length}V",
			FormatVersion:   2,
			MessageType:     "classic",
			Placement:       "waf_debug",
			FileMaxBytes:    10 * MiB,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	record(t, "blobstorages/create3", func(c *Client) {
		bsCreateResp3, err = c.CreateBlobStorage(&CreateBlobStorageInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-blobstorage-3",
			Path:             "/logs",
			AccountName:      "test",
			Container:        "fastly",
			SASToken:         "sv=2015-04-05&ss=b&srt=sco&sp=rw&se=2030-07-21T18%3A00%3A00Z&sig=3ABdLOJZosCp0o491T%2BqZGKIhafF1nlM3MzESDDD3Gg%3D",
			Period:           12,
			TimestampFormat:  "%Y-%m-%dT%H:%M:%S.000",
			CompressionCodec: "snappy",
			PublicKey:        "-----BEGIN PGP PUBLIC KEY BLOCK-----\n\nmQENBFyUD8sBCACyFnB39AuuTygseek+eA4fo0cgwva6/FSjnWq7riouQee8GgQ/\nibXTRyv4iVlwI12GswvMTIy7zNvs1R54i0qvsLr+IZ4GVGJqs6ZJnvQcqe3xPoR4\n8AnBfw90o32r/LuHf6QCJXi+AEu35koNlNAvLJ2B+KACaNB7N0EeWmqpV/1V2k9p\nlDYk+th7LcCuaFNGqKS/PrMnnMqR6VDLCjHhNx4KR79b0Twm/2qp6an3hyNRu8Gn\ndwxpf1/BUu3JWf+LqkN4Y3mbOmSUL3MaJNvyQguUzTfS0P0uGuBDHrJCVkMZCzDB\n89ag55jCPHyGeHBTd02gHMWzsg3WMBWvCsrzABEBAAG0JXRlcnJhZm9ybSAodGVz\ndCkgPHRlc3RAdGVycmFmb3JtLmNvbT6JAU4EEwEIADgWIQSHYyc6Kj9l6HzQsau6\nvFFc9jxV/wUCXJQPywIbAwULCQgHAgYVCgkICwIEFgIDAQIeAQIXgAAKCRC6vFFc\n9jxV/815CAClb32OxV7wG01yF97TzlyTl8TnvjMtoG29Mw4nSyg+mjM3b8N7iXm9\nOLX59fbDAWtBSldSZE22RXd3CvlFOG/EnKBXSjBtEqfyxYSnyOPkMPBYWGL/ApkX\nSvPYJ4LKdvipYToKFh3y9kk2gk1DcDBDyaaHvR+3rv1u3aoy7/s2EltAfDS3ZQIq\n7/cWTLJml/lleeB/Y6rPj8xqeCYhE5ahw9gsV/Mdqatl24V9Tks30iijx0Hhw+Gx\nkATUikMGr2GDVqoIRga5kXI7CzYff4rkc0Twn47fMHHHe/KY9M2yVnMHUXmAZwbG\nM1cMI/NH1DjevCKdGBLcRJlhuLPKF/anuQENBFyUD8sBCADIpd7r7GuPd6n/Ikxe\nu6h7umV6IIPoAm88xCYpTbSZiaK30Svh6Ywra9jfE2KlU9o6Y/art8ip0VJ3m07L\n4RSfSpnzqgSwdjSq5hNour2Fo/BzYhK7yaz2AzVSbe33R0+RYhb4b/6N+bKbjwGF\nftCsqVFMH+PyvYkLbvxyQrHlA9woAZaNThI1ztO5rGSnGUR8xt84eup28WIFKg0K\nUEGUcTzz+8QGAwAra+0ewPXo/AkO+8BvZjDidP417u6gpBHOJ9qYIcO9FxHeqFyu\nYrjlrxowEgXn5wO8xuNz6Vu1vhHGDHGDsRbZF8pv1d5O+0F1G7ttZ2GRRgVBZPwi\nkiyRABEBAAGJATYEGAEIACAWIQSHYyc6Kj9l6HzQsau6vFFc9jxV/wUCXJQPywIb\nDAAKCRC6vFFc9jxV/9YOCACe8qmOSnKQpQfW+PqYOqo3dt7JyweTs3FkD6NT8Zml\ndYy/vkstbTjPpX6aTvUZjkb46BVi7AOneVHpD5GBqvRsZ9iVgDYHaehmLCdKiG5L\n3Tp90NN+QY5WDbsGmsyk6+6ZMYejb4qYfweQeduOj27aavCJdLkCYMoRKfcFYI8c\nFaNmEfKKy/r1PO20NXEG6t9t05K/frHy6ZG8bCNYdpagfFVot47r9JaQqWlTNtIR\n5+zkkSq/eG9BEtRij3a6cTdQbktdBzx2KBeI0PYc1vlZR0LpuFKZqY9vlE6vTGLR\nwMfrTEOvx0NxUM3rpaCgEmuWbB1G1Hu371oyr4srrr+N\n=28dr\n-----END PGP PUBLIC KEY BLOCK-----\n",
			Format:           "%h %l %u %{now}V %{req.method}V %{req.url}V %>s %{resp.http.Content-Length}V",
			FormatVersion:    2,
			MessageType:      "classic",
			Placement:        "waf_debug",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "blobstorages/cleanup", func(c *Client) {
			c.DeleteBlobStorage(&DeleteBlobStorageInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-blobstorage",
			})

			c.DeleteBlobStorage(&DeleteBlobStorageInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-blobstorage-2",
			})

			c.DeleteBlobStorage(&DeleteBlobStorageInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-blobstorage-3",
			})

			c.DeleteBlobStorage(&DeleteBlobStorageInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "new-test-blobstorage",
			})
		})
	}()

	if bsCreateResp1.Name != "test-blobstorage" {
		t.Errorf("bad name: %q", bsCreateResp1.Name)
	}
	if bsCreateResp1.Path != "/logs" {
		t.Errorf("bad path: %q", bsCreateResp1.Path)
	}
	if bsCreateResp1.AccountName != "test" {
		t.Errorf("bad account_name: %q", bsCreateResp1.AccountName)
	}
	if bsCreateResp1.Container != "fastly" {
		t.Errorf("bad container: %q", bsCreateResp1.Container)
	}
	if bsCreateResp1.SASToken != "sv=2015-04-05&ss=b&srt=sco&sp=rw&se=2030-07-21T18%3A00%3A00Z&sig=3ABdLOJZosCp0o491T%2BqZGKIhafF1nlM3MzESDDD3Gg%3D" {
		t.Errorf("bad sas_token: %q", bsCreateResp1.SASToken)
	}
	if bsCreateResp1.Period != 12 {
		t.Errorf("bad period: %q", bsCreateResp1.Period)
	}
	if bsCreateResp1.TimestampFormat != "%Y-%m-%dT%H:%M:%S.000" {
		t.Errorf("bad timestamp_format: %q", bsCreateResp1.TimestampFormat)
	}
	if bsCreateResp1.CompressionCodec != "snappy" {
		t.Errorf("bad compression_codec: %q", bsCreateResp1.CompressionCodec)
	}
	if bsCreateResp1.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", bsCreateResp1.GzipLevel)
	}
	if bsCreateResp1.PublicKey != "-----BEGIN PGP PUBLIC KEY BLOCK-----\n\nmQENBFyUD8sBCACyFnB39AuuTygseek+eA4fo0cgwva6/FSjnWq7riouQee8GgQ/\nibXTRyv4iVlwI12GswvMTIy7zNvs1R54i0qvsLr+IZ4GVGJqs6ZJnvQcqe3xPoR4\n8AnBfw90o32r/LuHf6QCJXi+AEu35koNlNAvLJ2B+KACaNB7N0EeWmqpV/1V2k9p\nlDYk+th7LcCuaFNGqKS/PrMnnMqR6VDLCjHhNx4KR79b0Twm/2qp6an3hyNRu8Gn\ndwxpf1/BUu3JWf+LqkN4Y3mbOmSUL3MaJNvyQguUzTfS0P0uGuBDHrJCVkMZCzDB\n89ag55jCPHyGeHBTd02gHMWzsg3WMBWvCsrzABEBAAG0JXRlcnJhZm9ybSAodGVz\ndCkgPHRlc3RAdGVycmFmb3JtLmNvbT6JAU4EEwEIADgWIQSHYyc6Kj9l6HzQsau6\nvFFc9jxV/wUCXJQPywIbAwULCQgHAgYVCgkICwIEFgIDAQIeAQIXgAAKCRC6vFFc\n9jxV/815CAClb32OxV7wG01yF97TzlyTl8TnvjMtoG29Mw4nSyg+mjM3b8N7iXm9\nOLX59fbDAWtBSldSZE22RXd3CvlFOG/EnKBXSjBtEqfyxYSnyOPkMPBYWGL/ApkX\nSvPYJ4LKdvipYToKFh3y9kk2gk1DcDBDyaaHvR+3rv1u3aoy7/s2EltAfDS3ZQIq\n7/cWTLJml/lleeB/Y6rPj8xqeCYhE5ahw9gsV/Mdqatl24V9Tks30iijx0Hhw+Gx\nkATUikMGr2GDVqoIRga5kXI7CzYff4rkc0Twn47fMHHHe/KY9M2yVnMHUXmAZwbG\nM1cMI/NH1DjevCKdGBLcRJlhuLPKF/anuQENBFyUD8sBCADIpd7r7GuPd6n/Ikxe\nu6h7umV6IIPoAm88xCYpTbSZiaK30Svh6Ywra9jfE2KlU9o6Y/art8ip0VJ3m07L\n4RSfSpnzqgSwdjSq5hNour2Fo/BzYhK7yaz2AzVSbe33R0+RYhb4b/6N+bKbjwGF\nftCsqVFMH+PyvYkLbvxyQrHlA9woAZaNThI1ztO5rGSnGUR8xt84eup28WIFKg0K\nUEGUcTzz+8QGAwAra+0ewPXo/AkO+8BvZjDidP417u6gpBHOJ9qYIcO9FxHeqFyu\nYrjlrxowEgXn5wO8xuNz6Vu1vhHGDHGDsRbZF8pv1d5O+0F1G7ttZ2GRRgVBZPwi\nkiyRABEBAAGJATYEGAEIACAWIQSHYyc6Kj9l6HzQsau6vFFc9jxV/wUCXJQPywIb\nDAAKCRC6vFFc9jxV/9YOCACe8qmOSnKQpQfW+PqYOqo3dt7JyweTs3FkD6NT8Zml\ndYy/vkstbTjPpX6aTvUZjkb46BVi7AOneVHpD5GBqvRsZ9iVgDYHaehmLCdKiG5L\n3Tp90NN+QY5WDbsGmsyk6+6ZMYejb4qYfweQeduOj27aavCJdLkCYMoRKfcFYI8c\nFaNmEfKKy/r1PO20NXEG6t9t05K/frHy6ZG8bCNYdpagfFVot47r9JaQqWlTNtIR\n5+zkkSq/eG9BEtRij3a6cTdQbktdBzx2KBeI0PYc1vlZR0LpuFKZqY9vlE6vTGLR\nwMfrTEOvx0NxUM3rpaCgEmuWbB1G1Hu371oyr4srrr+N\n=28dr\n-----END PGP PUBLIC KEY BLOCK-----\n" {
		t.Errorf("bad public_key: %q", bsCreateResp1.PublicKey)
	}
	if bsCreateResp1.Format != "%h %l %u %{now}V %{req.method}V %{req.url}V %>s %{resp.http.Content-Length}V" {
		t.Errorf("bad format: %q", bsCreateResp1.Format)
	}
	if bsCreateResp1.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", bsCreateResp1.FormatVersion)
	}
	if bsCreateResp1.MessageType != "classic" {
		t.Errorf("bad message_type: %q", bsCreateResp1.MessageType)
	}
	if bsCreateResp1.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", bsCreateResp1.Placement)
	}
	if bsCreateResp1.FileMaxBytes != MiB {
		t.Errorf("bad file_max_bytes: %q", bsCreateResp1.FileMaxBytes)
	}

	if bsCreateResp2.CompressionCodec != "" {
		t.Errorf("bad compression_codec: %q", bsCreateResp2.CompressionCodec)
	}
	if bsCreateResp2.GzipLevel != 8 {
		t.Errorf("bad gzip_level: %q", bsCreateResp2.GzipLevel)
	}
	if bsCreateResp2.FileMaxBytes != 10*MiB {
		t.Errorf("bad file_max_bytes: %q", bsCreateResp2.FileMaxBytes)
	}

	if bsCreateResp3.CompressionCodec != "snappy" {
		t.Errorf("bad compression_codec: %q", bsCreateResp3.CompressionCodec)
	}
	if bsCreateResp3.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", bsCreateResp3.GzipLevel)
	}
	if bsCreateResp3.FileMaxBytes != 0 {
		t.Errorf("bad file_max_bytes: %q", bsCreateResp3.FileMaxBytes)
	}

	// List
	var bsl []*BlobStorage
	record(t, "blobstorages/list", func(c *Client) {
		bsl, err = c.ListBlobStorages(&ListBlobStoragesInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(bsl) < 1 {
		t.Errorf("bad blob storages: %v", bsl)
	}

	// Get
	var bsGetResp *BlobStorage
	record(t, "blobstorages/get", func(c *Client) {
		bsGetResp, err = c.GetBlobStorage(&GetBlobStorageInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-blobstorage",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if bsCreateResp1.Name != bsGetResp.Name {
		t.Errorf("bad name: %q", bsCreateResp1.Name)
	}
	if bsCreateResp1.Path != bsGetResp.Path {
		t.Errorf("bad path: %q", bsCreateResp1.Path)
	}
	if bsCreateResp1.AccountName != bsGetResp.AccountName {
		t.Errorf("bad account_name: %q", bsCreateResp1.AccountName)
	}
	if bsCreateResp1.Container != bsGetResp.Container {
		t.Errorf("bad container: %q", bsCreateResp1.Container)
	}
	if bsCreateResp1.SASToken != bsGetResp.SASToken {
		t.Errorf("bad sas_token: %q", bsCreateResp1.SASToken)
	}
	if bsCreateResp1.Period != bsGetResp.Period {
		t.Errorf("bad period: %q", bsCreateResp1.Period)
	}
	if bsCreateResp1.TimestampFormat != bsGetResp.TimestampFormat {
		t.Errorf("bad timestamp_format: %q", bsCreateResp1.TimestampFormat)
	}
	if bsCreateResp1.CompressionCodec != bsGetResp.CompressionCodec {
		t.Errorf("bad compression_codec: %q", bsCreateResp1.CompressionCodec)
	}
	if bsCreateResp1.GzipLevel != bsGetResp.GzipLevel {
		t.Errorf("bad gzip_level: %q", bsCreateResp1.GzipLevel)
	}
	if bsCreateResp1.PublicKey != bsGetResp.PublicKey {
		t.Errorf("bad public_key: %q", bsCreateResp1.PublicKey)
	}
	if bsCreateResp1.Format != bsGetResp.Format {
		t.Errorf("bad format: %q", bsCreateResp1.Format)
	}
	if bsCreateResp1.FormatVersion != bsGetResp.FormatVersion {
		t.Errorf("bad format_version: %q", bsCreateResp1.FormatVersion)
	}
	if bsCreateResp1.MessageType != bsGetResp.MessageType {
		t.Errorf("bad message_type: %q", bsCreateResp1.MessageType)
	}
	if bsCreateResp1.Placement != bsGetResp.Placement {
		t.Errorf("bad placement: %q", bsCreateResp1.Placement)
	}

	// Update
	var bsUpdateResp1, bsUpdateResp2, bsUpdateResp3 *BlobStorage
	record(t, "blobstorages/update", func(c *Client) {
		bsUpdateResp1, err = c.UpdateBlobStorage(&UpdateBlobStorageInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-blobstorage",
			NewName:          String("new-test-blobstorage"),
			CompressionCodec: String("zstd"),
			FileMaxBytes:     Uint(5 * MiB),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Test that CompressionCodec can be set for a an endpoint where
	// GzipLevel was specified at creation time.
	record(t, "blobstorages/update2", func(c *Client) {
		bsUpdateResp2, err = c.UpdateBlobStorage(&UpdateBlobStorageInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-blobstorage-2",
			CompressionCodec: String("zstd"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Test that GzipLevel can be set for an endpoint where CompressionCodec
	// was set at creation time.
	record(t, "blobstorages/update3", func(c *Client) {
		bsUpdateResp3, err = c.UpdateBlobStorage(&UpdateBlobStorageInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-blobstorage-3",
			GzipLevel:      Uint(9),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if bsUpdateResp1.Name != "new-test-blobstorage" {
		t.Errorf("bad name: %q", bsUpdateResp1.Name)
	}
	if bsUpdateResp1.CompressionCodec != "zstd" {
		t.Errorf("bad compression_codec: %q", bsUpdateResp1.CompressionCodec)
	}
	if bsUpdateResp1.FileMaxBytes != 5*MiB {
		t.Errorf("bad file_max_bytes: %q", bsUpdateResp1.FileMaxBytes)
	}

	if bsUpdateResp2.CompressionCodec != "zstd" {
		t.Errorf("bad compression_codec: %q", bsUpdateResp2.CompressionCodec)
	}
	if bsUpdateResp2.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", bsUpdateResp2.GzipLevel)
	}

	if bsUpdateResp3.CompressionCodec != "" {
		t.Errorf("bad compression_codec: %q", bsUpdateResp3.CompressionCodec)
	}
	if bsUpdateResp3.GzipLevel != 9 {
		t.Errorf("bad gzip_level: %q", bsUpdateResp3.GzipLevel)
	}

	// Delete
	record(t, "blobstorages/delete", func(c *Client) {
		err = c.DeleteBlobStorage(&DeleteBlobStorageInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "new-test-blobstorage",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListBlobStorages_validation(t *testing.T) {
	var err error
	_, err = testClient.ListBlobStorages(&ListBlobStoragesInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListBlobStorages(&ListBlobStoragesInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateBlobStorage_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateBlobStorage(&CreateBlobStorageInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateBlobStorage(&CreateBlobStorageInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetBlobStorage_validation(t *testing.T) {
	var err error
	_, err = testClient.GetBlobStorage(&GetBlobStorageInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetBlobStorage(&GetBlobStorageInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetBlobStorage(&GetBlobStorageInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateBlobStorage_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateBlobStorage(&UpdateBlobStorageInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateBlobStorage(&UpdateBlobStorageInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateBlobStorage(&UpdateBlobStorageInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteBlobStorage_validation(t *testing.T) {
	var err error
	err = testClient.DeleteBlobStorage(&DeleteBlobStorageInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteBlobStorage(&DeleteBlobStorageInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteBlobStorage(&DeleteBlobStorageInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
