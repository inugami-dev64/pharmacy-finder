package pharmafinder

import "embed"

//go:embed frontend/build/*
var ServerFS embed.FS

//go:embed migrations/*
var MigrationsFS embed.FS
