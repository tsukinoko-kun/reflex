package biome

func Default(out string) *Configuration {
	all := struct {
		All bool `json:"all"`
	}{All: true}

	return &Configuration{
		Schema:  "./node_modules/@biomejs/biome/configuration_schema.json",
		Assists: nil,
		CSS: &CssConfiguration{
			Formatter: &CssFormatter{
				Enabled:    true,
				QuoteStyle: "double",
			},
			Linter: &CssLinter{
				Enabled: true,
			},
			Parser: &CssParser{
				AllowWrongLineComments: false,
			},
		},
		Extends: nil,
		Files: &FilesConfiguration{
			IgnoreUnknown: false,
			Ignore:        []string{"node_modules", out},
		},
		Formatter: &FormatterConfiguration{
			Enabled:     true,
			IndentStyle: "tab",
			LineWidth:   120,
			LineEnding:  "lf",
		},
		GraphQL: &GraphqlConfiguration{
			Formatter: &GraphqlFormatter{
				Enabled:    true,
				QuoteStyle: "double",
			},
			Linter: &GraphqlLinter{
				Enabled: true,
			},
		},
		JavaScript: &JavascriptConfiguration{
			Formatter: &JavascriptFormatter{
				QuoteStyle:       "single",
				JsxQuoteStyle:    "double",
				QuoteProperties:  "asNeeded",
				TrailingCommas:   "all",
				Semicolons:       "asNeeded",
				ArrowParentheses: "always",
			},
			Parser: &JavascriptParser{
				UnsafeParameterDecoratorsEnabled: false,
			},
			JsxRuntime: "transparent",
			Linter: &JavascriptLinter{
				Enabled: true,
			},
		},
		JSON: &JsonConfiguration{
			Parser: &JsonParser{
				AllowTrailingCommas: false,
			},
			Formatter: &JsonFormatter{
				Enabled:        true,
				TrailingCommas: "none",
			},
		},
		Linter: &LinterConfiguration{
			Enabled: true,
			Rules: map[string]any{
				"recommended": true,
				"suspicious": map[string]any{
					"all":                  true,
					"noReactSpecificProps": "off",
				},
				"correctness": all,
				"style": map[string]any{
					"all":             true,
					"noDefaultExport": "off",
				},
				"complexity":  all,
				"performance": all,
				"security":    all,
				"a11y":        all,
				"nursery": map[string]any{
					"useSortedClasses":               "error",
					"noCommonJs":                     "error",
					"noDuplicateElseIf":              "warn",
					"noDuplicateProperties":          "error",
					"noDuplicatedFields":             "error",
					"noDynamicNamespaceImportAccess": "warn",
				},
			},
		},
		OrganizeImports: &OrganizeImports{
			Enabled: true,
		},
		Overrides: nil,
		VCS: &VcsConfiguration{
			Enabled:       true,
			ClientKind:    "git",
			UseIgnoreFile: false,
		},
	}
}

// Main configuration type
type Configuration struct {
	Schema          string                   `json:"$schema,omitempty"`
	Assists         *AssistsConfiguration    `json:"assists,omitempty"`
	CSS             *CssConfiguration        `json:"css,omitempty"`
	Extends         []string                 `json:"extends,omitempty"`
	Files           *FilesConfiguration      `json:"files,omitempty"`
	Formatter       *FormatterConfiguration  `json:"formatter,omitempty"`
	GraphQL         *GraphqlConfiguration    `json:"graphql,omitempty"`
	JavaScript      *JavascriptConfiguration `json:"javascript,omitempty"`
	JSON            *JsonConfiguration       `json:"json,omitempty"`
	Linter          *LinterConfiguration     `json:"linter,omitempty"`
	OrganizeImports *OrganizeImports         `json:"organizeImports,omitempty"`
	Overrides       []OverridePattern        `json:"overrides,omitempty"`
	VCS             *VcsConfiguration        `json:"vcs,omitempty"`
}

type AssistsConfiguration struct {
	Actions *Actions `json:"actions,omitempty"`
	Enabled bool     `json:"enabled,omitempty"`
	Ignore  []string `json:"ignore,omitempty"`
	Include []string `json:"include,omitempty"`
}

type Actions struct {
	Source *Source `json:"source,omitempty"`
}

type Source struct {
	SortJsxProps  string `json:"sortJsxProps,omitempty"`
	UseSortedKeys string `json:"useSortedKeys,omitempty"`
}

type CssConfiguration struct {
	Assists   *CssAssists   `json:"assists,omitempty"`
	Formatter *CssFormatter `json:"formatter,omitempty"`
	Linter    *CssLinter    `json:"linter,omitempty"`
	Parser    *CssParser    `json:"parser,omitempty"`
}

type CssAssists struct {
	Enabled bool `json:"enabled,omitempty"`
}

type CssFormatter struct {
	Enabled     bool   `json:"enabled,omitempty"`
	IndentStyle string `json:"indentStyle,omitempty"`
	IndentWidth int    `json:"indentWidth,omitempty"`
	LineEnding  string `json:"lineEnding,omitempty"`
	LineWidth   int    `json:"lineWidth,omitempty"`
	QuoteStyle  string `json:"quoteStyle,omitempty"`
}

type CssLinter struct {
	Enabled bool `json:"enabled,omitempty"`
}

type CssParser struct {
	AllowWrongLineComments bool `json:"allowWrongLineComments"`
	CssModules             bool `json:"cssModules,omitempty"`
}

type FilesConfiguration struct {
	Ignore        []string `json:"ignore,omitempty"`
	IgnoreUnknown bool     `json:"ignoreUnknown,omitempty"`
	Include       []string `json:"include,omitempty"`
	MaxSize       int64    `json:"maxSize,omitempty"`
}

type FormatterConfiguration struct {
	AttributePosition string   `json:"attributePosition,omitempty"`
	BracketSpacing    bool     `json:"bracketSpacing,omitempty"`
	Enabled           bool     `json:"enabled,omitempty"`
	FormatWithErrors  bool     `json:"formatWithErrors,omitempty"`
	Ignore            []string `json:"ignore,omitempty"`
	Include           []string `json:"include,omitempty"`
	IndentSize        int      `json:"indentSize,omitempty"`
	IndentStyle       string   `json:"indentStyle,omitempty"`
	IndentWidth       int      `json:"indentWidth,omitempty"`
	LineEnding        string   `json:"lineEnding,omitempty"`
	LineWidth         int      `json:"lineWidth,omitempty"`
	UseEditorconfig   bool     `json:"useEditorconfig,omitempty"`
}

type GraphqlConfiguration struct {
	Formatter *GraphqlFormatter `json:"formatter,omitempty"`
	Linter    *GraphqlLinter    `json:"linter,omitempty"`
}

type GraphqlFormatter struct {
	BracketSpacing bool   `json:"bracketSpacing,omitempty"`
	Enabled        bool   `json:"enabled,omitempty"`
	IndentStyle    string `json:"indentStyle,omitempty"`
	IndentWidth    int    `json:"indentWidth,omitempty"`
	LineEnding     string `json:"lineEnding,omitempty"`
	LineWidth      int    `json:"lineWidth,omitempty"`
	QuoteStyle     string `json:"quoteStyle,omitempty"`
}

type GraphqlLinter struct {
	Enabled bool `json:"enabled,omitempty"`
}

type JavascriptConfiguration struct {
	Assists         *JavascriptAssists         `json:"assists,omitempty"`
	Formatter       *JavascriptFormatter       `json:"formatter,omitempty"`
	Globals         []string                   `json:"globals,omitempty"`
	JsxRuntime      string                     `json:"jsxRuntime,omitempty"`
	Linter          *JavascriptLinter          `json:"linter,omitempty"`
	OrganizeImports *JavascriptOrganizeImports `json:"organizeImports,omitempty"`
	Parser          *JavascriptParser          `json:"parser,omitempty"`
}

type JavascriptAssists struct {
	Enabled bool `json:"enabled,omitempty"`
}

type JavascriptFormatter struct {
	ArrowParentheses  string `json:"arrowParentheses,omitempty"`
	AttributePosition string `json:"attributePosition,omitempty"`
	BracketSameLine   bool   `json:"bracketSameLine,omitempty"`
	BracketSpacing    bool   `json:"bracketSpacing,omitempty"`
	Enabled           bool   `json:"enabled,omitempty"`
	IndentSize        int    `json:"indentSize,omitempty"`
	IndentStyle       string `json:"indentStyle,omitempty"`
	IndentWidth       int    `json:"indentWidth,omitempty"`
	JsxQuoteStyle     string `json:"jsxQuoteStyle,omitempty"`
	LineEnding        string `json:"lineEnding,omitempty"`
	LineWidth         int    `json:"lineWidth,omitempty"`
	QuoteProperties   string `json:"quoteProperties,omitempty"`
	QuoteStyle        string `json:"quoteStyle,omitempty"`
	Semicolons        string `json:"semicolons,omitempty"`
	TrailingCommas    string `json:"trailingCommas,omitempty"`
}

type JavascriptLinter struct {
	Enabled bool `json:"enabled,omitempty"`
}

type JavascriptOrganizeImports struct{}

type JavascriptParser struct {
	UnsafeParameterDecoratorsEnabled bool `json:"unsafeParameterDecoratorsEnabled"`
}

type JsonConfiguration struct {
	Assists   *JsonAssists   `json:"assists,omitempty"`
	Formatter *JsonFormatter `json:"formatter,omitempty"`
	Linter    *JsonLinter    `json:"linter,omitempty"`
	Parser    *JsonParser    `json:"parser,omitempty"`
}

type JsonAssists struct {
	Enabled bool `json:"enabled,omitempty"`
}

type JsonFormatter struct {
	Enabled        bool   `json:"enabled,omitempty"`
	IndentSize     int    `json:"indentSize,omitempty"`
	IndentStyle    string `json:"indentStyle,omitempty"`
	IndentWidth    int    `json:"indentWidth,omitempty"`
	LineEnding     string `json:"lineEnding,omitempty"`
	LineWidth      int    `json:"lineWidth,omitempty"`
	TrailingCommas string `json:"trailingCommas,omitempty"`
}

type JsonLinter struct {
	Enabled bool `json:"enabled,omitempty"`
}

type JsonParser struct {
	AllowComments       bool `json:"allowComments"`
	AllowTrailingCommas bool `json:"allowTrailingCommas"`
}

type LinterConfiguration struct {
	Enabled bool           `json:"enabled,omitempty"`
	Ignore  []string       `json:"ignore,omitempty"`
	Include []string       `json:"include,omitempty"`
	Rules   map[string]any `json:"rules,omitempty"`
}

type OrganizeImports struct {
	Enabled bool     `json:"enabled,omitempty"`
	Ignore  []string `json:"ignore,omitempty"`
	Include []string `json:"include,omitempty"`
}

type OverridePattern struct {
	CSS             *CssConfiguration        `json:"css,omitempty"`
	Formatter       *FormatterConfiguration  `json:"formatter,omitempty"`
	GraphQL         *GraphqlConfiguration    `json:"graphql,omitempty"`
	Ignore          []string                 `json:"ignore,omitempty"`
	Include         []string                 `json:"include,omitempty"`
	JavaScript      *JavascriptConfiguration `json:"javascript,omitempty"`
	JSON            *JsonConfiguration       `json:"json,omitempty"`
	Linter          *LinterConfiguration     `json:"linter,omitempty"`
	OrganizeImports *OrganizeImports         `json:"organizeImports,omitempty"`
}

type VcsConfiguration struct {
	ClientKind    string `json:"clientKind,omitempty"`
	DefaultBranch string `json:"defaultBranch,omitempty"`
	Enabled       bool   `json:"enabled,omitempty"`
	Root          string `json:"root,omitempty"`
	UseIgnoreFile bool   `json:"useIgnoreFile,omitempty"`
}
