package dummy

import (
	"context"
	"testing"
	"time"

	"github.com/netdata/go-orchestrator/job/confgroup"
	"github.com/netdata/go-orchestrator/module"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDiscovery(t *testing.T) {
	tests := map[string]struct {
		cfg     Config
		wantErr bool
	}{
		"valid config": {
			cfg: Config{
				Registry: confgroup.Registry{"module1": confgroup.Default{}},
				Names:    []string{"module1", "module2"},
			},
		},
		"invalid config, registry not set": {
			cfg: Config{
				Names: []string{"module1", "module2"},
			},
			wantErr: true,
		},
		"invalid config, names not set": {
			cfg: Config{
				Names: []string{"module1", "module2"},
			},
			wantErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			d, err := NewDiscovery(test.cfg)

			if test.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, d)
			}
		})
	}
}

func TestDiscovery_Run(t *testing.T) {
	expected := []*confgroup.Group{
		{
			Source: "module1",
			Configs: []confgroup.Config{
				{
					"name":                "module1",
					"module":              "module1",
					"update_every":        module.UpdateEvery,
					"autodetection_retry": module.AutoDetectionRetry,
					"priority":            module.Priority,
				},
			},
		},
		{
			Source: "module2",
			Configs: []confgroup.Config{
				{
					"name":                "module2",
					"module":              "module2",
					"update_every":        module.UpdateEvery,
					"autodetection_retry": module.AutoDetectionRetry,
					"priority":            module.Priority,
				},
			},
		},
	}

	reg := confgroup.Registry{
		"module1": {},
		"module2": {},
	}
	cfg := Config{
		Registry: reg,
		Names:    []string{"module1", "module2"},
	}

	discovery, err := NewDiscovery(cfg)
	require.NoError(t, err)

	in := make(chan []*confgroup.Group)
	timeout := time.Second * 2

	go discovery.Run(context.Background(), in)

	var actual []*confgroup.Group
	select {
	case actual = <-in:
	case <-time.After(timeout):
		t.Logf("discovery timed out after %s", timeout)
	}
	assert.Equal(t, expected, actual)
}
