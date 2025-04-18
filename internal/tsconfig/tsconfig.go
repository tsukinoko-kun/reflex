package tsconfig

func Default(out string) map[string]any {
	return map[string]any{
		"compilerOptions": map[string]any{
			"target": "ES2020",
			"lib": []string{
				"DOM",
				"DOM.Iterable",
				"ES2020",
			},
			"allowJs":           true,
			"skipLibCheck":      true,
			"strict":            true,
			"strictNullChecks":  true,
			"noEmit":            true,
			"esModuleInterop":   true,
			"module":            "esnext",
			"moduleResolution":  "bundler",
			"resolveJsonModule": true,
			"isolatedModules":   true,
			"jsx":               "preserve",
			"jsxImportSource":   "react",
			"incremental":       true,
		},
		"include": []string{
			"**/*.ts",
			"**/*.tsx",
		},
		"exclude": []string{
			"node_modules",
			out,
		},
	}
}
