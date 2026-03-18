package openapi

// Builder provides a fluent API for constructing OpenAPI documents manually.
// This is an alternative to the automatic Generator for advanced customization.
type Builder struct {
	doc Document
}

// NewBuilder creates a new OpenAPI document Builder.
func NewBuilder() *Builder {
	return &Builder{
		doc: Document{
			OpenAPI: "3.1.0",
			Info: InfoObject{
				Title:   "API Documentation",
				Version: "1.0.0",
			},
			Paths: make(map[string]PathItem),
		},
	}
}

// SetInfo sets the API info.
func (b *Builder) SetInfo(title, version, description string) *Builder {
	b.doc.Info = InfoObject{
		Title:       title,
		Version:     version,
		Description: description,
	}
	return b
}

// AddServer adds a server.
func (b *Builder) AddServer(url, description string) *Builder {
	b.doc.Servers = append(b.doc.Servers, ServerObject{
		URL:         url,
		Description: description,
	})
	return b
}

// AddTag adds a tag.
func (b *Builder) AddTag(name, description string) *Builder {
	b.doc.Tags = append(b.doc.Tags, TagObject{
		Name:        name,
		Description: description,
	})
	return b
}

// AddPath adds a path item.
func (b *Builder) AddPath(path string, item PathItem) *Builder {
	b.doc.Paths[path] = item
	return b
}

// SetComponents sets the components.
func (b *Builder) SetComponents(components *Components) *Builder {
	b.doc.Components = components
	return b
}

// SetSecurity sets global security requirements.
func (b *Builder) SetSecurity(security []map[string][]string) *Builder {
	b.doc.Security = security
	return b
}

// Build returns the constructed document.
func (b *Builder) Build() Document {
	return b.doc
}
