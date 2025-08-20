package pharmafinder

import "embed"

//go:embed frontend/build/*
var ServerFS embed.FS

//go:embed db/migrations/*
var MigrationsFS embed.FS
