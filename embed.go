package pharmafinder

import "embed"

//go:embed frontend/build/*
var ServerFS embed.FS

//go:embed db/migrations/*
var MigrationsFS embed.FS

//go:embed db/independent-pharmacies.json
var PharmacyJSON embed.FS
