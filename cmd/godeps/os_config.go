package main

import buildify "github.com/bazelbuild/buildtools/build"

func generateOsConfigExprs(ruleDir string) []buildify.Expr {
	var exprs []buildify.Expr

	for _, platform := range SupportedPlatforms {
		ruleName := platform.String()
		if ruleDir == "" {
			ruleName = "__config_" + ruleName
		}

		rule := &buildify.CallExpr{
			X: &buildify.Ident{Name: "config_setting"},
			List: []buildify.Expr{
				&buildify.AssignExpr{
					LHS: &buildify.Ident{Name: "name"},
					Op:  "=",
					RHS: &buildify.StringExpr{Value: ruleName},
				},
				&buildify.AssignExpr{
					LHS: &buildify.Ident{Name: "values"},
					Op:  "=",
					RHS: &buildify.DictExpr{
						List: []*buildify.KeyValueExpr{
							{
								Key:   &buildify.StringExpr{Value: "os"},
								Value: &buildify.StringExpr{Value: platform.OS},
							},
							{
								Key:   &buildify.StringExpr{Value: "cpu"},
								Value: &buildify.StringExpr{Value: platform.Arch},
							},
						},
					},
				},
			},
		}

		exprs = append(exprs, rule)
	}

	return exprs
}
