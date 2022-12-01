package languageservice

import (
	yamlparser "github.com/CircleCI-Public/circleci-yaml-language-server/pkg/parser"
	"github.com/CircleCI-Public/circleci-yaml-language-server/pkg/parser/validate"
	"github.com/CircleCI-Public/circleci-yaml-language-server/pkg/utils"
	"go.lsp.dev/protocol"
)

type DiagnosticType struct {
	diagnostics  *[]protocol.Diagnostic
	yamlDocument yamlparser.YamlDocument
}

func Diagnostic(params protocol.PublishDiagnosticsParams, cache *utils.Cache, context *utils.LsContext) protocol.PublishDiagnosticsParams {
	yamlDocument, err := yamlparser.ParseFromUriWithCache(params.URI, cache)

	if yamlDocument.Version < 2.1 || err != nil {
		// TODO: Handle error
		return protocol.PublishDiagnosticsParams{}
	}

	diag := DiagnosticType{
		diagnostics:  &[]protocol.Diagnostic{},
		yamlDocument: yamlDocument,
	}

	yamlDocument.ValidateYAML()
	diag.addDiagnostics(*yamlDocument.Diagnostics)

	validator := yamlparser.JSONSchemaValidator{}
	validator.ParseJsonSchema()

	diag.addDiagnostics(validator.ValidateWithJSONSchema(diag.yamlDocument.RootNode, diag.yamlDocument.Content))

	validateStruct := validate.Validate{
		Doc:         diag.yamlDocument,
		Diagnostics: &[]protocol.Diagnostic{},
		Cache:       cache,
		Context:     context,
	}
	validateStruct.Validate()
	diag.addDiagnostics(*validateStruct.Diagnostics)

	diagnosticParams := protocol.PublishDiagnosticsParams{
		URI:         params.URI,
		Diagnostics: *diag.diagnostics,
	}

	return diagnosticParams
}

func (diag *DiagnosticType) addDiagnostics(diagnostic []protocol.Diagnostic) {
	*diag.diagnostics = append(*diag.diagnostics, diagnostic...)
}
