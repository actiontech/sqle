// Code generated from JavaParser.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // JavaParser

import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by JavaParser.
type JavaParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by JavaParser#compilationUnit.
	VisitCompilationUnit(ctx *CompilationUnitContext) interface{}

	// Visit a parse tree produced by JavaParser#packageDeclaration.
	VisitPackageDeclaration(ctx *PackageDeclarationContext) interface{}

	// Visit a parse tree produced by JavaParser#importDeclaration.
	VisitImportDeclaration(ctx *ImportDeclarationContext) interface{}

	// Visit a parse tree produced by JavaParser#typeDeclaration.
	VisitTypeDeclaration(ctx *TypeDeclarationContext) interface{}

	// Visit a parse tree produced by JavaParser#modifier.
	VisitModifier(ctx *ModifierContext) interface{}

	// Visit a parse tree produced by JavaParser#classOrInterfaceModifier.
	VisitClassOrInterfaceModifier(ctx *ClassOrInterfaceModifierContext) interface{}

	// Visit a parse tree produced by JavaParser#variableModifier.
	VisitVariableModifier(ctx *VariableModifierContext) interface{}

	// Visit a parse tree produced by JavaParser#classDeclaration.
	VisitClassDeclaration(ctx *ClassDeclarationContext) interface{}

	// Visit a parse tree produced by JavaParser#typeParameters.
	VisitTypeParameters(ctx *TypeParametersContext) interface{}

	// Visit a parse tree produced by JavaParser#typeParameter.
	VisitTypeParameter(ctx *TypeParameterContext) interface{}

	// Visit a parse tree produced by JavaParser#typeBound.
	VisitTypeBound(ctx *TypeBoundContext) interface{}

	// Visit a parse tree produced by JavaParser#enumDeclaration.
	VisitEnumDeclaration(ctx *EnumDeclarationContext) interface{}

	// Visit a parse tree produced by JavaParser#enumConstants.
	VisitEnumConstants(ctx *EnumConstantsContext) interface{}

	// Visit a parse tree produced by JavaParser#enumConstant.
	VisitEnumConstant(ctx *EnumConstantContext) interface{}

	// Visit a parse tree produced by JavaParser#enumBodyDeclarations.
	VisitEnumBodyDeclarations(ctx *EnumBodyDeclarationsContext) interface{}

	// Visit a parse tree produced by JavaParser#interfaceDeclaration.
	VisitInterfaceDeclaration(ctx *InterfaceDeclarationContext) interface{}

	// Visit a parse tree produced by JavaParser#classBody.
	VisitClassBody(ctx *ClassBodyContext) interface{}

	// Visit a parse tree produced by JavaParser#interfaceBody.
	VisitInterfaceBody(ctx *InterfaceBodyContext) interface{}

	// Visit a parse tree produced by JavaParser#classBodyDeclaration.
	VisitClassBodyDeclaration(ctx *ClassBodyDeclarationContext) interface{}

	// Visit a parse tree produced by JavaParser#memberDeclaration.
	VisitMemberDeclaration(ctx *MemberDeclarationContext) interface{}

	// Visit a parse tree produced by JavaParser#methodDeclaration.
	VisitMethodDeclaration(ctx *MethodDeclarationContext) interface{}

	// Visit a parse tree produced by JavaParser#methodBody.
	VisitMethodBody(ctx *MethodBodyContext) interface{}

	// Visit a parse tree produced by JavaParser#typeTypeOrVoid.
	VisitTypeTypeOrVoid(ctx *TypeTypeOrVoidContext) interface{}

	// Visit a parse tree produced by JavaParser#genericMethodDeclaration.
	VisitGenericMethodDeclaration(ctx *GenericMethodDeclarationContext) interface{}

	// Visit a parse tree produced by JavaParser#genericConstructorDeclaration.
	VisitGenericConstructorDeclaration(ctx *GenericConstructorDeclarationContext) interface{}

	// Visit a parse tree produced by JavaParser#constructorDeclaration.
	VisitConstructorDeclaration(ctx *ConstructorDeclarationContext) interface{}

	// Visit a parse tree produced by JavaParser#compactConstructorDeclaration.
	VisitCompactConstructorDeclaration(ctx *CompactConstructorDeclarationContext) interface{}

	// Visit a parse tree produced by JavaParser#fieldDeclaration.
	VisitFieldDeclaration(ctx *FieldDeclarationContext) interface{}

	// Visit a parse tree produced by JavaParser#interfaceBodyDeclaration.
	VisitInterfaceBodyDeclaration(ctx *InterfaceBodyDeclarationContext) interface{}

	// Visit a parse tree produced by JavaParser#interfaceMemberDeclaration.
	VisitInterfaceMemberDeclaration(ctx *InterfaceMemberDeclarationContext) interface{}

	// Visit a parse tree produced by JavaParser#constDeclaration.
	VisitConstDeclaration(ctx *ConstDeclarationContext) interface{}

	// Visit a parse tree produced by JavaParser#constantDeclarator.
	VisitConstantDeclarator(ctx *ConstantDeclaratorContext) interface{}

	// Visit a parse tree produced by JavaParser#interfaceMethodDeclaration.
	VisitInterfaceMethodDeclaration(ctx *InterfaceMethodDeclarationContext) interface{}

	// Visit a parse tree produced by JavaParser#interfaceMethodModifier.
	VisitInterfaceMethodModifier(ctx *InterfaceMethodModifierContext) interface{}

	// Visit a parse tree produced by JavaParser#genericInterfaceMethodDeclaration.
	VisitGenericInterfaceMethodDeclaration(ctx *GenericInterfaceMethodDeclarationContext) interface{}

	// Visit a parse tree produced by JavaParser#interfaceCommonBodyDeclaration.
	VisitInterfaceCommonBodyDeclaration(ctx *InterfaceCommonBodyDeclarationContext) interface{}

	// Visit a parse tree produced by JavaParser#variableDeclarators.
	VisitVariableDeclarators(ctx *VariableDeclaratorsContext) interface{}

	// Visit a parse tree produced by JavaParser#variableDeclarator.
	VisitVariableDeclarator(ctx *VariableDeclaratorContext) interface{}

	// Visit a parse tree produced by JavaParser#variableDeclaratorId.
	VisitVariableDeclaratorId(ctx *VariableDeclaratorIdContext) interface{}

	// Visit a parse tree produced by JavaParser#variableInitializer.
	VisitVariableInitializer(ctx *VariableInitializerContext) interface{}

	// Visit a parse tree produced by JavaParser#arrayInitializer.
	VisitArrayInitializer(ctx *ArrayInitializerContext) interface{}

	// Visit a parse tree produced by JavaParser#classOrInterfaceType.
	VisitClassOrInterfaceType(ctx *ClassOrInterfaceTypeContext) interface{}

	// Visit a parse tree produced by JavaParser#typeArgument.
	VisitTypeArgument(ctx *TypeArgumentContext) interface{}

	// Visit a parse tree produced by JavaParser#qualifiedNameList.
	VisitQualifiedNameList(ctx *QualifiedNameListContext) interface{}

	// Visit a parse tree produced by JavaParser#formalParameters.
	VisitFormalParameters(ctx *FormalParametersContext) interface{}

	// Visit a parse tree produced by JavaParser#receiverParameter.
	VisitReceiverParameter(ctx *ReceiverParameterContext) interface{}

	// Visit a parse tree produced by JavaParser#formalParameterList.
	VisitFormalParameterList(ctx *FormalParameterListContext) interface{}

	// Visit a parse tree produced by JavaParser#formalParameter.
	VisitFormalParameter(ctx *FormalParameterContext) interface{}

	// Visit a parse tree produced by JavaParser#lastFormalParameter.
	VisitLastFormalParameter(ctx *LastFormalParameterContext) interface{}

	// Visit a parse tree produced by JavaParser#lambdaLVTIList.
	VisitLambdaLVTIList(ctx *LambdaLVTIListContext) interface{}

	// Visit a parse tree produced by JavaParser#lambdaLVTIParameter.
	VisitLambdaLVTIParameter(ctx *LambdaLVTIParameterContext) interface{}

	// Visit a parse tree produced by JavaParser#qualifiedName.
	VisitQualifiedName(ctx *QualifiedNameContext) interface{}

	// Visit a parse tree produced by JavaParser#literal.
	VisitLiteral(ctx *LiteralContext) interface{}

	// Visit a parse tree produced by JavaParser#integerLiteral.
	VisitIntegerLiteral(ctx *IntegerLiteralContext) interface{}

	// Visit a parse tree produced by JavaParser#floatLiteral.
	VisitFloatLiteral(ctx *FloatLiteralContext) interface{}

	// Visit a parse tree produced by JavaParser#altAnnotationQualifiedName.
	VisitAltAnnotationQualifiedName(ctx *AltAnnotationQualifiedNameContext) interface{}

	// Visit a parse tree produced by JavaParser#annotation.
	VisitAnnotation(ctx *AnnotationContext) interface{}

	// Visit a parse tree produced by JavaParser#elementValuePairs.
	VisitElementValuePairs(ctx *ElementValuePairsContext) interface{}

	// Visit a parse tree produced by JavaParser#elementValuePair.
	VisitElementValuePair(ctx *ElementValuePairContext) interface{}

	// Visit a parse tree produced by JavaParser#elementValue.
	VisitElementValue(ctx *ElementValueContext) interface{}

	// Visit a parse tree produced by JavaParser#elementValueArrayInitializer.
	VisitElementValueArrayInitializer(ctx *ElementValueArrayInitializerContext) interface{}

	// Visit a parse tree produced by JavaParser#annotationTypeDeclaration.
	VisitAnnotationTypeDeclaration(ctx *AnnotationTypeDeclarationContext) interface{}

	// Visit a parse tree produced by JavaParser#annotationTypeBody.
	VisitAnnotationTypeBody(ctx *AnnotationTypeBodyContext) interface{}

	// Visit a parse tree produced by JavaParser#annotationTypeElementDeclaration.
	VisitAnnotationTypeElementDeclaration(ctx *AnnotationTypeElementDeclarationContext) interface{}

	// Visit a parse tree produced by JavaParser#annotationTypeElementRest.
	VisitAnnotationTypeElementRest(ctx *AnnotationTypeElementRestContext) interface{}

	// Visit a parse tree produced by JavaParser#annotationMethodOrConstantRest.
	VisitAnnotationMethodOrConstantRest(ctx *AnnotationMethodOrConstantRestContext) interface{}

	// Visit a parse tree produced by JavaParser#annotationMethodRest.
	VisitAnnotationMethodRest(ctx *AnnotationMethodRestContext) interface{}

	// Visit a parse tree produced by JavaParser#annotationConstantRest.
	VisitAnnotationConstantRest(ctx *AnnotationConstantRestContext) interface{}

	// Visit a parse tree produced by JavaParser#defaultValue.
	VisitDefaultValue(ctx *DefaultValueContext) interface{}

	// Visit a parse tree produced by JavaParser#moduleDeclaration.
	VisitModuleDeclaration(ctx *ModuleDeclarationContext) interface{}

	// Visit a parse tree produced by JavaParser#moduleBody.
	VisitModuleBody(ctx *ModuleBodyContext) interface{}

	// Visit a parse tree produced by JavaParser#moduleDirective.
	VisitModuleDirective(ctx *ModuleDirectiveContext) interface{}

	// Visit a parse tree produced by JavaParser#requiresModifier.
	VisitRequiresModifier(ctx *RequiresModifierContext) interface{}

	// Visit a parse tree produced by JavaParser#recordDeclaration.
	VisitRecordDeclaration(ctx *RecordDeclarationContext) interface{}

	// Visit a parse tree produced by JavaParser#recordHeader.
	VisitRecordHeader(ctx *RecordHeaderContext) interface{}

	// Visit a parse tree produced by JavaParser#recordComponentList.
	VisitRecordComponentList(ctx *RecordComponentListContext) interface{}

	// Visit a parse tree produced by JavaParser#recordComponent.
	VisitRecordComponent(ctx *RecordComponentContext) interface{}

	// Visit a parse tree produced by JavaParser#recordBody.
	VisitRecordBody(ctx *RecordBodyContext) interface{}

	// Visit a parse tree produced by JavaParser#block.
	VisitBlock(ctx *BlockContext) interface{}

	// Visit a parse tree produced by JavaParser#blockStatement.
	VisitBlockStatement(ctx *BlockStatementContext) interface{}

	// Visit a parse tree produced by JavaParser#localVariableDeclaration.
	VisitLocalVariableDeclaration(ctx *LocalVariableDeclarationContext) interface{}

	// Visit a parse tree produced by JavaParser#identifier.
	VisitIdentifier(ctx *IdentifierContext) interface{}

	// Visit a parse tree produced by JavaParser#typeIdentifier.
	VisitTypeIdentifier(ctx *TypeIdentifierContext) interface{}

	// Visit a parse tree produced by JavaParser#localTypeDeclaration.
	VisitLocalTypeDeclaration(ctx *LocalTypeDeclarationContext) interface{}

	// Visit a parse tree produced by JavaParser#statement.
	VisitStatement(ctx *StatementContext) interface{}

	// Visit a parse tree produced by JavaParser#catchClause.
	VisitCatchClause(ctx *CatchClauseContext) interface{}

	// Visit a parse tree produced by JavaParser#catchType.
	VisitCatchType(ctx *CatchTypeContext) interface{}

	// Visit a parse tree produced by JavaParser#finallyBlock.
	VisitFinallyBlock(ctx *FinallyBlockContext) interface{}

	// Visit a parse tree produced by JavaParser#resourceSpecification.
	VisitResourceSpecification(ctx *ResourceSpecificationContext) interface{}

	// Visit a parse tree produced by JavaParser#resources.
	VisitResources(ctx *ResourcesContext) interface{}

	// Visit a parse tree produced by JavaParser#resource.
	VisitResource(ctx *ResourceContext) interface{}

	// Visit a parse tree produced by JavaParser#switchBlockStatementGroup.
	VisitSwitchBlockStatementGroup(ctx *SwitchBlockStatementGroupContext) interface{}

	// Visit a parse tree produced by JavaParser#switchLabel.
	VisitSwitchLabel(ctx *SwitchLabelContext) interface{}

	// Visit a parse tree produced by JavaParser#forControl.
	VisitForControl(ctx *ForControlContext) interface{}

	// Visit a parse tree produced by JavaParser#forInit.
	VisitForInit(ctx *ForInitContext) interface{}

	// Visit a parse tree produced by JavaParser#enhancedForControl.
	VisitEnhancedForControl(ctx *EnhancedForControlContext) interface{}

	// Visit a parse tree produced by JavaParser#parExpression.
	VisitParExpression(ctx *ParExpressionContext) interface{}

	// Visit a parse tree produced by JavaParser#expressionList.
	VisitExpressionList(ctx *ExpressionListContext) interface{}

	// Visit a parse tree produced by JavaParser#methodCall.
	VisitMethodCall(ctx *MethodCallContext) interface{}

	// Visit a parse tree produced by JavaParser#expression.
	VisitExpression(ctx *ExpressionContext) interface{}

	// Visit a parse tree produced by JavaParser#pattern.
	VisitPattern(ctx *PatternContext) interface{}

	// Visit a parse tree produced by JavaParser#lambdaExpression.
	VisitLambdaExpression(ctx *LambdaExpressionContext) interface{}

	// Visit a parse tree produced by JavaParser#lambdaParameters.
	VisitLambdaParameters(ctx *LambdaParametersContext) interface{}

	// Visit a parse tree produced by JavaParser#lambdaBody.
	VisitLambdaBody(ctx *LambdaBodyContext) interface{}

	// Visit a parse tree produced by JavaParser#primary.
	VisitPrimary(ctx *PrimaryContext) interface{}

	// Visit a parse tree produced by JavaParser#switchExpression.
	VisitSwitchExpression(ctx *SwitchExpressionContext) interface{}

	// Visit a parse tree produced by JavaParser#switchLabeledRule.
	VisitSwitchLabeledRule(ctx *SwitchLabeledRuleContext) interface{}

	// Visit a parse tree produced by JavaParser#guardedPattern.
	VisitGuardedPattern(ctx *GuardedPatternContext) interface{}

	// Visit a parse tree produced by JavaParser#switchRuleOutcome.
	VisitSwitchRuleOutcome(ctx *SwitchRuleOutcomeContext) interface{}

	// Visit a parse tree produced by JavaParser#classType.
	VisitClassType(ctx *ClassTypeContext) interface{}

	// Visit a parse tree produced by JavaParser#creator.
	VisitCreator(ctx *CreatorContext) interface{}

	// Visit a parse tree produced by JavaParser#createdName.
	VisitCreatedName(ctx *CreatedNameContext) interface{}

	// Visit a parse tree produced by JavaParser#innerCreator.
	VisitInnerCreator(ctx *InnerCreatorContext) interface{}

	// Visit a parse tree produced by JavaParser#arrayCreatorRest.
	VisitArrayCreatorRest(ctx *ArrayCreatorRestContext) interface{}

	// Visit a parse tree produced by JavaParser#classCreatorRest.
	VisitClassCreatorRest(ctx *ClassCreatorRestContext) interface{}

	// Visit a parse tree produced by JavaParser#explicitGenericInvocation.
	VisitExplicitGenericInvocation(ctx *ExplicitGenericInvocationContext) interface{}

	// Visit a parse tree produced by JavaParser#typeArgumentsOrDiamond.
	VisitTypeArgumentsOrDiamond(ctx *TypeArgumentsOrDiamondContext) interface{}

	// Visit a parse tree produced by JavaParser#nonWildcardTypeArgumentsOrDiamond.
	VisitNonWildcardTypeArgumentsOrDiamond(ctx *NonWildcardTypeArgumentsOrDiamondContext) interface{}

	// Visit a parse tree produced by JavaParser#nonWildcardTypeArguments.
	VisitNonWildcardTypeArguments(ctx *NonWildcardTypeArgumentsContext) interface{}

	// Visit a parse tree produced by JavaParser#typeList.
	VisitTypeList(ctx *TypeListContext) interface{}

	// Visit a parse tree produced by JavaParser#typeType.
	VisitTypeType(ctx *TypeTypeContext) interface{}

	// Visit a parse tree produced by JavaParser#primitiveType.
	VisitPrimitiveType(ctx *PrimitiveTypeContext) interface{}

	// Visit a parse tree produced by JavaParser#typeArguments.
	VisitTypeArguments(ctx *TypeArgumentsContext) interface{}

	// Visit a parse tree produced by JavaParser#superSuffix.
	VisitSuperSuffix(ctx *SuperSuffixContext) interface{}

	// Visit a parse tree produced by JavaParser#explicitGenericInvocationSuffix.
	VisitExplicitGenericInvocationSuffix(ctx *ExplicitGenericInvocationSuffixContext) interface{}

	// Visit a parse tree produced by JavaParser#arguments.
	VisitArguments(ctx *ArgumentsContext) interface{}
}
