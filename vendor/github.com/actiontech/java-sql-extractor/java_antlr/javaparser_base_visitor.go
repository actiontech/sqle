// Code generated from JavaParser.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // JavaParser

import "github.com/antlr4-go/antlr/v4"

type BaseJavaParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseJavaParserVisitor) VisitCompilationUnit(ctx *CompilationUnitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitPackageDeclaration(ctx *PackageDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitImportDeclaration(ctx *ImportDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitTypeDeclaration(ctx *TypeDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitModifier(ctx *ModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitClassOrInterfaceModifier(ctx *ClassOrInterfaceModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitVariableModifier(ctx *VariableModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitClassDeclaration(ctx *ClassDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitTypeParameters(ctx *TypeParametersContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitTypeParameter(ctx *TypeParameterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitTypeBound(ctx *TypeBoundContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitEnumDeclaration(ctx *EnumDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitEnumConstants(ctx *EnumConstantsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitEnumConstant(ctx *EnumConstantContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitEnumBodyDeclarations(ctx *EnumBodyDeclarationsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitInterfaceDeclaration(ctx *InterfaceDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitClassBody(ctx *ClassBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitInterfaceBody(ctx *InterfaceBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitClassBodyDeclaration(ctx *ClassBodyDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitMemberDeclaration(ctx *MemberDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitMethodDeclaration(ctx *MethodDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitMethodBody(ctx *MethodBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitTypeTypeOrVoid(ctx *TypeTypeOrVoidContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitGenericMethodDeclaration(ctx *GenericMethodDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitGenericConstructorDeclaration(ctx *GenericConstructorDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitConstructorDeclaration(ctx *ConstructorDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitCompactConstructorDeclaration(ctx *CompactConstructorDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitFieldDeclaration(ctx *FieldDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitInterfaceBodyDeclaration(ctx *InterfaceBodyDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitInterfaceMemberDeclaration(ctx *InterfaceMemberDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitConstDeclaration(ctx *ConstDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitConstantDeclarator(ctx *ConstantDeclaratorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitInterfaceMethodDeclaration(ctx *InterfaceMethodDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitInterfaceMethodModifier(ctx *InterfaceMethodModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitGenericInterfaceMethodDeclaration(ctx *GenericInterfaceMethodDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitInterfaceCommonBodyDeclaration(ctx *InterfaceCommonBodyDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitVariableDeclarators(ctx *VariableDeclaratorsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitVariableDeclarator(ctx *VariableDeclaratorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitVariableDeclaratorId(ctx *VariableDeclaratorIdContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitVariableInitializer(ctx *VariableInitializerContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitArrayInitializer(ctx *ArrayInitializerContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitClassOrInterfaceType(ctx *ClassOrInterfaceTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitTypeArgument(ctx *TypeArgumentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitQualifiedNameList(ctx *QualifiedNameListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitFormalParameters(ctx *FormalParametersContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitReceiverParameter(ctx *ReceiverParameterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitFormalParameterList(ctx *FormalParameterListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitFormalParameter(ctx *FormalParameterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitLastFormalParameter(ctx *LastFormalParameterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitLambdaLVTIList(ctx *LambdaLVTIListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitLambdaLVTIParameter(ctx *LambdaLVTIParameterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitQualifiedName(ctx *QualifiedNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitLiteral(ctx *LiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitIntegerLiteral(ctx *IntegerLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitFloatLiteral(ctx *FloatLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitAltAnnotationQualifiedName(ctx *AltAnnotationQualifiedNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitAnnotation(ctx *AnnotationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitElementValuePairs(ctx *ElementValuePairsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitElementValuePair(ctx *ElementValuePairContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitElementValue(ctx *ElementValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitElementValueArrayInitializer(ctx *ElementValueArrayInitializerContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitAnnotationTypeDeclaration(ctx *AnnotationTypeDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitAnnotationTypeBody(ctx *AnnotationTypeBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitAnnotationTypeElementDeclaration(ctx *AnnotationTypeElementDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitAnnotationTypeElementRest(ctx *AnnotationTypeElementRestContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitAnnotationMethodOrConstantRest(ctx *AnnotationMethodOrConstantRestContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitAnnotationMethodRest(ctx *AnnotationMethodRestContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitAnnotationConstantRest(ctx *AnnotationConstantRestContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitDefaultValue(ctx *DefaultValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitModuleDeclaration(ctx *ModuleDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitModuleBody(ctx *ModuleBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitModuleDirective(ctx *ModuleDirectiveContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitRequiresModifier(ctx *RequiresModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitRecordDeclaration(ctx *RecordDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitRecordHeader(ctx *RecordHeaderContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitRecordComponentList(ctx *RecordComponentListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitRecordComponent(ctx *RecordComponentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitRecordBody(ctx *RecordBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitBlock(ctx *BlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitBlockStatement(ctx *BlockStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitLocalVariableDeclaration(ctx *LocalVariableDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitIdentifier(ctx *IdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitTypeIdentifier(ctx *TypeIdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitLocalTypeDeclaration(ctx *LocalTypeDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitStatement(ctx *StatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitCatchClause(ctx *CatchClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitCatchType(ctx *CatchTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitFinallyBlock(ctx *FinallyBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitResourceSpecification(ctx *ResourceSpecificationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitResources(ctx *ResourcesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitResource(ctx *ResourceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitSwitchBlockStatementGroup(ctx *SwitchBlockStatementGroupContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitSwitchLabel(ctx *SwitchLabelContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitForControl(ctx *ForControlContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitForInit(ctx *ForInitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitEnhancedForControl(ctx *EnhancedForControlContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitParExpression(ctx *ParExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitExpressionList(ctx *ExpressionListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitMethodCall(ctx *MethodCallContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitExpression(ctx *ExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitPattern(ctx *PatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitLambdaExpression(ctx *LambdaExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitLambdaParameters(ctx *LambdaParametersContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitLambdaBody(ctx *LambdaBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitPrimary(ctx *PrimaryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitSwitchExpression(ctx *SwitchExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitSwitchLabeledRule(ctx *SwitchLabeledRuleContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitGuardedPattern(ctx *GuardedPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitSwitchRuleOutcome(ctx *SwitchRuleOutcomeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitClassType(ctx *ClassTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitCreator(ctx *CreatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitCreatedName(ctx *CreatedNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitInnerCreator(ctx *InnerCreatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitArrayCreatorRest(ctx *ArrayCreatorRestContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitClassCreatorRest(ctx *ClassCreatorRestContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitExplicitGenericInvocation(ctx *ExplicitGenericInvocationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitTypeArgumentsOrDiamond(ctx *TypeArgumentsOrDiamondContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitNonWildcardTypeArgumentsOrDiamond(ctx *NonWildcardTypeArgumentsOrDiamondContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitNonWildcardTypeArguments(ctx *NonWildcardTypeArgumentsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitTypeList(ctx *TypeListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitTypeType(ctx *TypeTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitPrimitiveType(ctx *PrimitiveTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitTypeArguments(ctx *TypeArgumentsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitSuperSuffix(ctx *SuperSuffixContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitExplicitGenericInvocationSuffix(ctx *ExplicitGenericInvocationSuffixContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseJavaParserVisitor) VisitArguments(ctx *ArgumentsContext) interface{} {
	return v.VisitChildren(ctx)
}
