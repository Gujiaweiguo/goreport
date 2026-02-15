package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDashboard_Serialize(t *testing.T) {
	t.Run("成功序列化Config和Components", func(t *testing.T) {
		d := &Dashboard{
			Config: DashboardConfig{
				Width:           1920,
				Height:          1080,
				BackgroundColor: "#ffffff",
			},
			Components: []DashboardComponent{
				{
					ID:     "comp-1",
					Title:  "Test Component",
					Type:   "chart",
					Width:  400,
					Height: 300,
					X:      10,
					Y:      10,
				},
			},
		}

		err := d.serialize()
		assert.NoError(t, err)
		assert.NotEmpty(t, d.ConfigJSON)
		assert.NotEmpty(t, d.ComponentsJSON)

		var config DashboardConfig
		err = json.Unmarshal([]byte(d.ConfigJSON), &config)
		assert.NoError(t, err)
		assert.Equal(t, 1920, config.Width)
		assert.Equal(t, 1080, config.Height)

		var components []DashboardComponent
		err = json.Unmarshal([]byte(d.ComponentsJSON), &components)
		assert.NoError(t, err)
		assert.Len(t, components, 1)
		assert.Equal(t, "comp-1", components[0].ID)
	})

	t.Run("空Config和Components不序列化", func(t *testing.T) {
		d := &Dashboard{}

		err := d.serialize()
		assert.NoError(t, err)
		assert.Empty(t, d.ConfigJSON)
		assert.Empty(t, d.ComponentsJSON)
	})
}

func TestDashboard_Deserialize(t *testing.T) {
	t.Run("成功反序列化Config和Components", func(t *testing.T) {
		configJSON, _ := json.Marshal(DashboardConfig{
			Width:           1280,
			Height:          720,
			BackgroundColor: "#000000",
		})

		componentsJSON, _ := json.Marshal([]DashboardComponent{
			{
				ID:     "comp-2",
				Title:  "Component 2",
				Type:   "table",
				Width:  600,
				Height: 400,
			},
		})

		d := &Dashboard{
			ConfigJSON:     string(configJSON),
			ComponentsJSON: string(componentsJSON),
		}

		err := d.deserialize()
		assert.NoError(t, err)
		assert.Equal(t, 1280, d.Config.Width)
		assert.Equal(t, 720, d.Config.Height)
		assert.Equal(t, "#000000", d.Config.BackgroundColor)
		assert.Len(t, d.Components, 1)
		assert.Equal(t, "comp-2", d.Components[0].ID)
	})

	t.Run("无效的ConfigJSON使用默认值", func(t *testing.T) {
		d := &Dashboard{
			ConfigJSON: `{"invalid": "json`,
		}

		err := d.deserialize()
		assert.NoError(t, err)
		assert.Equal(t, 1920, d.Config.Width)
		assert.Equal(t, 1080, d.Config.Height)
		assert.Equal(t, "#0a0e27", d.Config.BackgroundColor)
	})

	t.Run("空JSON字符串不设置默认值", func(t *testing.T) {
		d := &Dashboard{}

		err := d.deserialize()
		assert.NoError(t, err)
		assert.Equal(t, 0, d.Config.Width)
		assert.Equal(t, 0, d.Config.Height)
		assert.Empty(t, d.Config.BackgroundColor)
	})

	t.Run("无效的ComponentsJSON使用空数组", func(t *testing.T) {
		d := &Dashboard{
			ConfigJSON:     `{"width": 1920, "height": 1080}`,
			ComponentsJSON: `{"invalid": "json"}`,
		}

		err := d.deserialize()
		assert.NoError(t, err)
		assert.Empty(t, d.Components)
	})
}

func TestDashboardConfig_DefaultValues(t *testing.T) {
	config := DashboardConfig{}
	assert.Equal(t, 0, config.Width)
	assert.Equal(t, 0, config.Height)
	assert.Empty(t, config.BackgroundColor)
}

func TestDashboardComponent_BasicProperties(t *testing.T) {
	comp := DashboardComponent{
		ID:      "test-id",
		Title:   "Test",
		Type:    "chart",
		Width:   100,
		Height:  100,
		X:       0,
		Y:       0,
		Visible: true,
		Locked:  false,
		Style:   map[string]interface{}{"color": "red"},
		Data:    map[string]interface{}{"value": 123},
	}

	assert.Equal(t, "test-id", comp.ID)
	assert.Equal(t, "Test", comp.Title)
	assert.Equal(t, "chart", comp.Type)
	assert.Equal(t, 100, comp.Width)
	assert.Equal(t, 100, comp.Height)
	assert.True(t, comp.Visible)
	assert.False(t, comp.Locked)
	assert.NotNil(t, comp.Style)
	assert.NotNil(t, comp.Data)
}

func TestDashboard_BeforeCreate(t *testing.T) {
	d := &Dashboard{
		Config: DashboardConfig{
			Width:           1920,
			Height:          1080,
			BackgroundColor: "#ffffff",
		},
		Components: []DashboardComponent{
			{ID: "comp-1", Type: "chart"},
		},
	}

	err := d.BeforeCreate(nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.ConfigJSON)
	assert.NotEmpty(t, d.ComponentsJSON)
}

func TestDashboard_BeforeUpdate(t *testing.T) {
	d := &Dashboard{
		Config: DashboardConfig{
			Width:  800,
			Height: 600,
		},
		Components: []DashboardComponent{
			{ID: "comp-2", Type: "table"},
		},
	}

	err := d.BeforeUpdate(nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.ConfigJSON)
	assert.NotEmpty(t, d.ComponentsJSON)
}

func TestDashboard_AfterFind(t *testing.T) {
	configJSON, _ := json.Marshal(DashboardConfig{
		Width:           1024,
		Height:          768,
		BackgroundColor: "#123456",
	})
	componentsJSON, _ := json.Marshal([]DashboardComponent{
		{ID: "comp-3", Type: "text"},
	})

	d := &Dashboard{
		ConfigJSON:     string(configJSON),
		ComponentsJSON: string(componentsJSON),
	}

	err := d.AfterFind(nil)
	assert.NoError(t, err)
	assert.Equal(t, 1024, d.Config.Width)
	assert.Equal(t, 768, d.Config.Height)
	assert.Len(t, d.Components, 1)
}
