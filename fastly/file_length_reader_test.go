package fastly

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileLengthReaderRemainingUpload(t *testing.T) {
	cases := []struct {
		name    string
		data    string
		offset  int64
		consume int
		want    string
	}{
		{"whole file", "abcdefgh", 0, 0, "abcdefgh"},
		{"offset", "abcdefgh", 2, 0, "cdefgh"},
		{"at EOF", "abcdefgh", 8, 0, ""},
		{"past EOF", "abcdefgh", 10, 0, ""},
		{"partially read wrapper", "abcdefgh", 0, 2, "cdefgh"},
		{"offset and partially read", "abcdefgh", 2, 2, "efgh"},
		{"empty file", "", 0, 0, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			filename := filepath.Join(t.TempDir(), "upload")
			require.NoError(t, os.WriteFile(filename, []byte(tc.data), 0o600))
			file, err := os.Open(filename)
			require.NoError(t, err)
			defer file.Close()
			_, err = file.Seek(tc.offset, io.SeekStart)
			require.NoError(t, err)
			reader, err := FileLengthReader(file)
			require.NoError(t, err)
			if tc.consume > 0 {
				_, err = io.ReadFull(reader, make([]byte, tc.consume))
				require.NoError(t, err)
			}
			initialLen := reader.Len()
			uploaded := make(chan string, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, readErr := io.ReadAll(r.Body)
				if readErr != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				uploaded <- string(body)
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()
			client, err := NewClientForEndpoint("", server.URL)
			require.NoError(t, err)
			err = client.InsertKVStoreKey(context.Background(), &InsertKVStoreKeyInput{
				StoreID: "store", Key: "key", Body: reader,
			})
			require.NoError(t, err)
			require.Equal(t, tc.want, <-uploaded)
			require.Equal(t, len(tc.want), initialLen)
		})
	}
}

func TestFileLengthReaderLenAfterRead(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "upload")
	require.NoError(t, os.WriteFile(filename, []byte("abcd"), 0o600))
	file, err := os.Open(filename)
	require.NoError(t, err)
	defer file.Close()
	reader, err := FileLengthReader(file)
	require.NoError(t, err)
	_, err = io.ReadFull(reader, make([]byte, 2))
	require.NoError(t, err)
	require.Equal(t, 2, reader.Len())
	data, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, "cd", string(data))
	require.Zero(t, reader.Len())
}
