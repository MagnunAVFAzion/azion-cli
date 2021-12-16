package create

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	sdk "github.com/aziontech/edgeservices-go-sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var resposeBody = `
{
    "id": 1753,
    "name": "{name}",
    "updated_at": "2021-12-16T01:10:07Z",
    "last_editor": "crazy.ape@azion.com",
    "active": false,
    "bound_nodes": 0,
    "permissions": [
        "read",
        "write"
    ]
}
`

func TestCreate(t *testing.T) {
	t.Run("invalid service", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_services/"),
			httpmock.StatusStringResponse(http.StatusUnprocessableEntity, "Invalid name"),
		)

		stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
		f := &cmdutil.Factory{
			HttpClient: func() (*http.Client, error) {
				return &http.Client{Transport: mock}, nil
			},
			IOStreams: &iostreams.IOStreams{
				Out: stdout,
				Err: stderr,
			},
		}
		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--name", ""})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})

	t.Run("without passing name", func(t *testing.T) {
		stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
		f := &cmdutil.Factory{
			IOStreams: &iostreams.IOStreams{
				Out: stdout,
				Err: stderr,
			},
		}
		cmd := NewCmd(f)

		cmd.SetArgs([]string{})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.EqualError(t, err, "required flag(s) \"name\" not set")
	})

	t.Run("create service with name", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_services/"),
			func(req *http.Request) (*http.Response, error) {
				request := &sdk.CreateServiceRequest{}
				body, _ := ioutil.ReadAll(req.Body)
				_ = json.Unmarshal(body, request)

				response := strings.ReplaceAll(resposeBody, "{name}", request.Name)

				return &http.Response{StatusCode: http.StatusCreated,
					Request: req,
					Body:    ioutil.NopCloser(strings.NewReader(response)),
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
				}, nil
			},
		)

		stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
		f := &cmdutil.Factory{
			HttpClient: func() (*http.Client, error) {
				return &http.Client{Transport: mock}, nil
			},
			IOStreams: &iostreams.IOStreams{
				Out: stdout,
				Err: stderr,
			},
		}

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--name", "BIRL"})
		cmd.SetIn(&bytes.Buffer{})
		cmd.SetOut(ioutil.Discard)
		cmd.SetErr(ioutil.Discard)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Equal(t, "ID: 1753\tName: BIRL \n", stdout.String())
	})
}