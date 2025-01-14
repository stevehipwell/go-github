// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOrganizationsService_GetAllOrganizationRulesets(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/rulesets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{
			"id": 26110,
			"name": "test ruleset",
			"target": "branch",
			"source_type": "Organization",
			"source": "o",
			"enforcement": "active",
			"bypass_mode": "none",
			"node_id": "nid",
			"_links": {
			  "self": {
				"href": "https://api.github.com/orgs/o/rulesets/26110"
			  }
			}
		}]`)
	})

	ctx := context.Background()
	rulesets, _, err := client.Organizations.GetAllOrganizationRulesets(ctx, "o")
	if err != nil {
		t.Errorf("Organizations.GetAllOrganizationRulesets returned error: %v", err)
	}

	want := []*Ruleset{{
		ID:          Ptr(int64(26110)),
		Name:        "test ruleset",
		Target:      Ptr("branch"),
		SourceType:  Ptr("Organization"),
		Source:      "o",
		Enforcement: "active",
		NodeID:      Ptr("nid"),
		Links: &RulesetLinks{
			Self: &RulesetLink{HRef: Ptr("https://api.github.com/orgs/o/rulesets/26110")},
		},
	}}
	if !cmp.Equal(rulesets, want) {
		t.Errorf("Organizations.GetAllOrganizationRulesets returned %+v, want %+v", rulesets, want)
	}

	const methodName = "GetAllOrganizationRulesets"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetAllOrganizationRulesets(ctx, "o")
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_CreateOrganizationRuleset_RepoNames(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/rulesets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"id": 21,
			"name": "ruleset",
			"target": "branch",
			"source_type": "Organization",
			"source": "o",
			"enforcement": "active",
			"bypass_actors": [
			  {
				"actor_id": 234,
				"actor_type": "Team"
			  }
			],
			"conditions": {
			  "ref_name": {
				"include": [
				  "refs/heads/main",
				  "refs/heads/master"
				],
				"exclude": [
				  "refs/heads/dev*"
				]
			  },
			  "repository_name": {
				"include": [
				  "important_repository",
				  "another_important_repository"
				],
				"exclude": [
				  "unimportant_repository"
				],
				"protected": true
			  }
			},
			"rules": [
			  {
				"type": "creation"
			  },
			  {
				"type": "update",
				"parameters": {
				  "update_allows_fetch_and_merge": true
				}
			  },
			  {
				"type": "deletion"
			  },
			  {
				"type": "required_linear_history"
			  },
			  {
				"type": "required_deployments",
				"parameters": {
				  "required_deployment_environments": ["test"]
				}
			  },
			  {
				"type": "required_signatures"
			  },
			  {
				"type": "pull_request",
				"parameters": {
				"allowed_merge_methods": ["rebase","squash"],
				  "dismiss_stale_reviews_on_push": true,
				  "require_code_owner_review": true,
				  "require_last_push_approval": true,
				  "required_approving_review_count": 1,
				  "required_review_thread_resolution": true
				}
			  },
			  {
				"type": "required_status_checks",
				"parameters": {
				  "required_status_checks": [
					{
					  "context": "test",
					  "integration_id": 1
					}
				  ],
				  "strict_required_status_checks_policy": true
				}
			  },
			  {
				"type": "non_fast_forward"
			  },
			  {
				"type": "commit_message_pattern",
				"parameters": {
				  "name": "avoid test commits",
				  "negate": true,
				  "operator": "starts_with",
				  "pattern": "[test]"
				}
			  },
			  {
				"type": "commit_author_email_pattern",
				"parameters": {
				  "operator": "contains",
				  "pattern": "github"
				}
			  },
			  {
				"type": "committer_email_pattern",
				"parameters": {
				  "name": "avoid commit emails",
				  "negate": true,
				  "operator": "ends_with",
				  "pattern": "abc"
				}
			  },
			  {
				"type": "branch_name_pattern",
				"parameters": {
				  "name": "avoid branch names",
				  "negate": true,
				  "operator": "regex",
				  "pattern": "github$"
				}
			  },
			  {
				"type": "tag_name_pattern",
				"parameters": {
				  "name": "avoid tag names",
				  "negate": true,
				  "operator": "contains",
				  "pattern": "github"
				}
			  },
			  {
			    "type": "code_scanning",
			    "parameters": {
				  "code_scanning_tools": [
				    {
					  "tool": "CodeQL",
					  "security_alerts_threshold": "high_or_higher",
					  "alerts_threshold": "errors"
				    }
				  ]
			    }
			  }
			]
		  }`)
	})

	ctx := context.Background()
	ruleset, _, err := client.Organizations.CreateOrganizationRuleset(ctx, "o", Ruleset{
		ID:          Ptr(int64(21)),
		Name:        "ruleset",
		Target:      Ptr("branch"),
		SourceType:  Ptr("Organization"),
		Source:      "o",
		Enforcement: "active",
		BypassActors: []*BypassActor{
			{
				ActorID:   Ptr(int64(234)),
				ActorType: Ptr("Team"),
			},
		},
		Conditions: &RulesetConditions{
			RefName: &RulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
			RepositoryName: &RulesetRepositoryNamesConditionParameters{
				Include:   []string{"important_repository", "another_important_repository"},
				Exclude:   []string{"unimportant_repository"},
				Protected: Ptr(true),
			},
		},
		Rules: []*RepositoryRule{
			NewCreationRule(),
			NewUpdateRule(&UpdateAllowsFetchAndMergeRuleParameters{
				UpdateAllowsFetchAndMerge: true,
			}),
			NewDeletionRule(),
			NewRequiredLinearHistoryRule(),
			NewRequiredDeploymentsRule(&RequiredDeploymentEnvironmentsRuleParameters{
				RequiredDeploymentEnvironments: []string{"test"},
			}),
			NewRequiredSignaturesRule(),
			NewPullRequestRule(&PullRequestRuleParameters{
				AllowedMergeMethods:            []MergeMethod{"rebase", "squash"},
				DismissStaleReviewsOnPush:      true,
				RequireCodeOwnerReview:         true,
				RequireLastPushApproval:        true,
				RequiredApprovingReviewCount:   1,
				RequiredReviewThreadResolution: true,
			}),
			NewRequiredStatusChecksRule(&RequiredStatusChecksRuleParameters{
				RequiredStatusChecks: []RuleRequiredStatusChecks{
					{
						Context:       "test",
						IntegrationID: Ptr(int64(1)),
					},
				},
				StrictRequiredStatusChecksPolicy: true,
			}),
			NewNonFastForwardRule(),
			NewCommitMessagePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid test commits"),
				Negate:   Ptr(true),
				Operator: "starts_with",
				Pattern:  "[test]",
			}),
			NewCommitAuthorEmailPatternRule(&RulePatternParameters{
				Operator: "contains",
				Pattern:  "github",
			}),
			NewCommitterEmailPatternRule(&RulePatternParameters{
				Name:     Ptr("avoid commit emails"),
				Negate:   Ptr(true),
				Operator: "ends_with",
				Pattern:  "abc",
			}),
			NewBranchNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid branch names"),
				Negate:   Ptr(true),
				Operator: "regex",
				Pattern:  "github$",
			}),
			NewTagNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid tag names"),
				Negate:   Ptr(true),
				Operator: "contains",
				Pattern:  "github",
			}),
			NewRequiredCodeScanningRule(&RequiredCodeScanningRuleParameters{
				RequiredCodeScanningTools: []*RuleRequiredCodeScanningTool{
					{
						Tool:                    "CodeQL",
						SecurityAlertsThreshold: "high_or_higher",
						AlertsThreshold:         "errors",
					},
				},
			}),
		},
	})
	if err != nil {
		t.Errorf("Organizations.CreateOrganizationRuleset returned error: %v", err)
	}

	want := &Ruleset{
		ID:          Ptr(int64(21)),
		Name:        "ruleset",
		Target:      Ptr("branch"),
		SourceType:  Ptr("Organization"),
		Source:      "o",
		Enforcement: "active",
		BypassActors: []*BypassActor{
			{
				ActorID:   Ptr(int64(234)),
				ActorType: Ptr("Team"),
			},
		},
		Conditions: &RulesetConditions{
			RefName: &RulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
			RepositoryName: &RulesetRepositoryNamesConditionParameters{
				Include:   []string{"important_repository", "another_important_repository"},
				Exclude:   []string{"unimportant_repository"},
				Protected: Ptr(true),
			},
		},
		Rules: []*RepositoryRule{
			NewCreationRule(),
			NewUpdateRule(&UpdateAllowsFetchAndMergeRuleParameters{
				UpdateAllowsFetchAndMerge: true,
			}),
			NewDeletionRule(),
			NewRequiredLinearHistoryRule(),
			NewRequiredDeploymentsRule(&RequiredDeploymentEnvironmentsRuleParameters{
				RequiredDeploymentEnvironments: []string{"test"},
			}),
			NewRequiredSignaturesRule(),
			NewPullRequestRule(&PullRequestRuleParameters{
				AllowedMergeMethods:            []MergeMethod{"rebase", "squash"},
				DismissStaleReviewsOnPush:      true,
				RequireCodeOwnerReview:         true,
				RequireLastPushApproval:        true,
				RequiredApprovingReviewCount:   1,
				RequiredReviewThreadResolution: true,
			}),
			NewRequiredStatusChecksRule(&RequiredStatusChecksRuleParameters{
				RequiredStatusChecks: []RuleRequiredStatusChecks{
					{
						Context:       "test",
						IntegrationID: Ptr(int64(1)),
					},
				},
				StrictRequiredStatusChecksPolicy: true,
			}),
			NewNonFastForwardRule(),
			NewCommitMessagePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid test commits"),
				Negate:   Ptr(true),
				Operator: "starts_with",
				Pattern:  "[test]",
			}),
			NewCommitAuthorEmailPatternRule(&RulePatternParameters{
				Operator: "contains",
				Pattern:  "github",
			}),
			NewCommitterEmailPatternRule(&RulePatternParameters{
				Name:     Ptr("avoid commit emails"),
				Negate:   Ptr(true),
				Operator: "ends_with",
				Pattern:  "abc",
			}),
			NewBranchNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid branch names"),
				Negate:   Ptr(true),
				Operator: "regex",
				Pattern:  "github$",
			}),
			NewTagNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid tag names"),
				Negate:   Ptr(true),
				Operator: "contains",
				Pattern:  "github",
			}),
			NewRequiredCodeScanningRule(&RequiredCodeScanningRuleParameters{
				RequiredCodeScanningTools: []*RuleRequiredCodeScanningTool{
					{
						Tool:                    "CodeQL",
						SecurityAlertsThreshold: "high_or_higher",
						AlertsThreshold:         "errors",
					},
				},
			}),
		},
	}
	if !cmp.Equal(ruleset, want) {
		t.Errorf("Organizations.CreateOrganizationRuleset returned %+v, want %+v", ruleset, want)
	}

	const methodName = "CreateOrganizationRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.CreateOrganizationRuleset(ctx, "o", Ruleset{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_CreateOrganizationRuleset_RepoProperty(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/rulesets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"id": 21,
			"name": "ruleset",
			"target": "branch",
			"source_type": "Organization",
			"source": "o",
			"enforcement": "active",
			"bypass_actors": [
			  {
				"actor_id": 234,
				"actor_type": "Team"
			  }
			],
			"conditions": {
			  "repository_property": {
				"include": [
					{
						"name": "testIncludeProp",
						"source": "custom",
						"property_values": [
							"true"
						]
					}
				],
				"exclude": [
					{
						"name": "testExcludeProp",
						"property_values": [
							"false"
						]
					}
				]
			  }
			},
			"rules": [
			  {
				"type": "creation"
			  },
			  {
				"type": "update",
				"parameters": {
				  "update_allows_fetch_and_merge": true
				}
			  },
			  {
				"type": "deletion"
			  },
			  {
				"type": "required_linear_history"
			  },
			  {
				"type": "required_deployments",
				"parameters": {
				  "required_deployment_environments": ["test"]
				}
			  },
			  {
				"type": "required_signatures"
			  },
			  {
				"type": "pull_request",
				"parameters": {
				"allowed_merge_methods": ["rebase","squash"],
				  "dismiss_stale_reviews_on_push": true,
				  "require_code_owner_review": true,
				  "require_last_push_approval": true,
				  "required_approving_review_count": 1,
				  "required_review_thread_resolution": true
				}
			  },
			  {
				"type": "required_status_checks",
				"parameters": {
				  "required_status_checks": [
					{
					  "context": "test",
					  "integration_id": 1
					}
				  ],
				  "strict_required_status_checks_policy": true
				}
			  },
			  {
				"type": "non_fast_forward"
			  },
			  {
				"type": "commit_message_pattern",
				"parameters": {
				  "name": "avoid test commits",
				  "negate": true,
				  "operator": "starts_with",
				  "pattern": "[test]"
				}
			  },
			  {
				"type": "commit_author_email_pattern",
				"parameters": {
				  "operator": "contains",
				  "pattern": "github"
				}
			  },
			  {
				"type": "committer_email_pattern",
				"parameters": {
				  "name": "avoid commit emails",
				  "negate": true,
				  "operator": "ends_with",
				  "pattern": "abc"
				}
			  },
			  {
				"type": "branch_name_pattern",
				"parameters": {
				  "name": "avoid branch names",
				  "negate": true,
				  "operator": "regex",
				  "pattern": "github$"
				}
			  },
			  {
				"type": "tag_name_pattern",
				"parameters": {
				  "name": "avoid tag names",
				  "negate": true,
				  "operator": "contains",
				  "pattern": "github"
				}
			  },
			  {
			    "type": "code_scanning",
			    "parameters": {
				  "code_scanning_tools": [
				    {
					  "tool": "CodeQL",
					  "security_alerts_threshold": "high_or_higher",
					  "alerts_threshold": "errors"
				    }
				  ]
			    }
			  }
			]
		  }`)
	})

	ctx := context.Background()
	ruleset, _, err := client.Organizations.CreateOrganizationRuleset(ctx, "o", Ruleset{
		ID:          Ptr(int64(21)),
		Name:        "ruleset",
		Target:      Ptr("branch"),
		SourceType:  Ptr("Organization"),
		Source:      "o",
		Enforcement: "active",
		BypassActors: []*BypassActor{
			{
				ActorID:   Ptr(int64(234)),
				ActorType: Ptr("Team"),
			},
		},
		Conditions: &RulesetConditions{
			RepositoryProperty: &RulesetRepositoryPropertyConditionParameters{
				Include: []RulesetRepositoryPropertyTargetParameters{
					{
						Name:   "testIncludeProp",
						Source: Ptr("custom"),
						Values: []string{"true"},
					},
				},
				Exclude: []RulesetRepositoryPropertyTargetParameters{
					{
						Name:   "testExcludeProp",
						Values: []string{"false"},
					},
				},
			},
		},
		Rules: []*RepositoryRule{
			NewCreationRule(),
			NewUpdateRule(&UpdateAllowsFetchAndMergeRuleParameters{
				UpdateAllowsFetchAndMerge: true,
			}),
			NewDeletionRule(),
			NewRequiredLinearHistoryRule(),
			NewRequiredDeploymentsRule(&RequiredDeploymentEnvironmentsRuleParameters{
				RequiredDeploymentEnvironments: []string{"test"},
			}),
			NewRequiredSignaturesRule(),
			NewPullRequestRule(&PullRequestRuleParameters{
				AllowedMergeMethods:            []MergeMethod{"rebase", "squash"},
				DismissStaleReviewsOnPush:      true,
				RequireCodeOwnerReview:         true,
				RequireLastPushApproval:        true,
				RequiredApprovingReviewCount:   1,
				RequiredReviewThreadResolution: true,
			}),
			NewRequiredStatusChecksRule(&RequiredStatusChecksRuleParameters{
				RequiredStatusChecks: []RuleRequiredStatusChecks{
					{
						Context:       "test",
						IntegrationID: Ptr(int64(1)),
					},
				},
				StrictRequiredStatusChecksPolicy: true,
			}),
			NewNonFastForwardRule(),
			NewCommitMessagePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid test commits"),
				Negate:   Ptr(true),
				Operator: "starts_with",
				Pattern:  "[test]",
			}),
			NewCommitAuthorEmailPatternRule(&RulePatternParameters{
				Operator: "contains",
				Pattern:  "github",
			}),
			NewCommitterEmailPatternRule(&RulePatternParameters{
				Name:     Ptr("avoid commit emails"),
				Negate:   Ptr(true),
				Operator: "ends_with",
				Pattern:  "abc",
			}),
			NewBranchNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid branch names"),
				Negate:   Ptr(true),
				Operator: "regex",
				Pattern:  "github$",
			}),
			NewTagNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid tag names"),
				Negate:   Ptr(true),
				Operator: "contains",
				Pattern:  "github",
			}),
			NewRequiredCodeScanningRule(&RequiredCodeScanningRuleParameters{
				RequiredCodeScanningTools: []*RuleRequiredCodeScanningTool{
					{
						Tool:                    "CodeQL",
						SecurityAlertsThreshold: "high_or_higher",
						AlertsThreshold:         "errors",
					},
				},
			}),
		},
	})
	if err != nil {
		t.Errorf("Organizations.CreateOrganizationRuleset returned error: %v", err)
	}

	want := &Ruleset{
		ID:          Ptr(int64(21)),
		Name:        "ruleset",
		Target:      Ptr("branch"),
		SourceType:  Ptr("Organization"),
		Source:      "o",
		Enforcement: "active",
		BypassActors: []*BypassActor{
			{
				ActorID:   Ptr(int64(234)),
				ActorType: Ptr("Team"),
			},
		},
		Conditions: &RulesetConditions{
			RepositoryProperty: &RulesetRepositoryPropertyConditionParameters{
				Include: []RulesetRepositoryPropertyTargetParameters{
					{
						Name:   "testIncludeProp",
						Source: Ptr("custom"),
						Values: []string{"true"},
					},
				},
				Exclude: []RulesetRepositoryPropertyTargetParameters{
					{
						Name:   "testExcludeProp",
						Values: []string{"false"},
					},
				},
			},
		},
		Rules: []*RepositoryRule{
			NewCreationRule(),
			NewUpdateRule(&UpdateAllowsFetchAndMergeRuleParameters{
				UpdateAllowsFetchAndMerge: true,
			}),
			NewDeletionRule(),
			NewRequiredLinearHistoryRule(),
			NewRequiredDeploymentsRule(&RequiredDeploymentEnvironmentsRuleParameters{
				RequiredDeploymentEnvironments: []string{"test"},
			}),
			NewRequiredSignaturesRule(),
			NewPullRequestRule(&PullRequestRuleParameters{
				AllowedMergeMethods:            []MergeMethod{"rebase", "squash"},
				DismissStaleReviewsOnPush:      true,
				RequireCodeOwnerReview:         true,
				RequireLastPushApproval:        true,
				RequiredApprovingReviewCount:   1,
				RequiredReviewThreadResolution: true,
			}),
			NewRequiredStatusChecksRule(&RequiredStatusChecksRuleParameters{
				RequiredStatusChecks: []RuleRequiredStatusChecks{
					{
						Context:       "test",
						IntegrationID: Ptr(int64(1)),
					},
				},
				StrictRequiredStatusChecksPolicy: true,
			}),
			NewNonFastForwardRule(),
			NewCommitMessagePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid test commits"),
				Negate:   Ptr(true),
				Operator: "starts_with",
				Pattern:  "[test]",
			}),
			NewCommitAuthorEmailPatternRule(&RulePatternParameters{
				Operator: "contains",
				Pattern:  "github",
			}),
			NewCommitterEmailPatternRule(&RulePatternParameters{
				Name:     Ptr("avoid commit emails"),
				Negate:   Ptr(true),
				Operator: "ends_with",
				Pattern:  "abc",
			}),
			NewBranchNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid branch names"),
				Negate:   Ptr(true),
				Operator: "regex",
				Pattern:  "github$",
			}),
			NewTagNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid tag names"),
				Negate:   Ptr(true),
				Operator: "contains",
				Pattern:  "github",
			}),
			NewRequiredCodeScanningRule(&RequiredCodeScanningRuleParameters{
				RequiredCodeScanningTools: []*RuleRequiredCodeScanningTool{
					{
						Tool:                    "CodeQL",
						SecurityAlertsThreshold: "high_or_higher",
						AlertsThreshold:         "errors",
					},
				},
			}),
		},
	}
	if !cmp.Equal(ruleset, want) {
		t.Errorf("Organizations.CreateOrganizationRuleset returned %+v, want %+v", ruleset, want)
	}

	const methodName = "CreateOrganizationRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.CreateOrganizationRuleset(ctx, "o", Ruleset{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_CreateOrganizationRuleset_RepoIDs(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/rulesets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"id": 21,
			"name": "ruleset",
			"target": "branch",
			"source_type": "Organization",
			"source": "o",
			"enforcement": "active",
			"bypass_actors": [
			  {
				"actor_id": 234,
				"actor_type": "Team"
			  }
			],
			"conditions": {
			  "ref_name": {
				"include": [
				  "refs/heads/main",
				  "refs/heads/master"
				],
				"exclude": [
				  "refs/heads/dev*"
				]
			  },
			  "repository_id": {
					"repository_ids": [ 123, 456 ]
				}
			},
			"rules": [
			  {
				"type": "creation"
			  },
			  {
				"type": "update",
				"parameters": {
				  "update_allows_fetch_and_merge": true
				}
			  },
			  {
				"type": "deletion"
			  },
			  {
				"type": "required_linear_history"
			  },
			  {
				"type": "required_deployments",
				"parameters": {
				  "required_deployment_environments": ["test"]
				}
			  },
			  {
				"type": "required_signatures"
			  },
			  {
				"type": "pull_request",
				"parameters": {
					"allowed_merge_methods": ["rebase","squash"],
				  "dismiss_stale_reviews_on_push": true,
				  "require_code_owner_review": true,
				  "require_last_push_approval": true,
				  "required_approving_review_count": 1,
				  "required_review_thread_resolution": true
				}
			  },
			  {
				"type": "required_status_checks",
				"parameters": {
				  "required_status_checks": [
					{
					  "context": "test",
					  "integration_id": 1
					}
				  ],
				  "strict_required_status_checks_policy": true
				}
			  },
			  {
				"type": "non_fast_forward"
			  },
			  {
				"type": "commit_message_pattern",
				"parameters": {
				  "name": "avoid test commits",
				  "negate": true,
				  "operator": "starts_with",
				  "pattern": "[test]"
				}
			  },
			  {
				"type": "commit_author_email_pattern",
				"parameters": {
				  "operator": "contains",
				  "pattern": "github"
				}
			  },
			  {
				"type": "committer_email_pattern",
				"parameters": {
				  "name": "avoid commit emails",
				  "negate": true,
				  "operator": "ends_with",
				  "pattern": "abc"
				}
			  },
			  {
				"type": "branch_name_pattern",
				"parameters": {
				  "name": "avoid branch names",
				  "negate": true,
				  "operator": "regex",
				  "pattern": "github$"
				}
			  },
			  {
				"type": "tag_name_pattern",
				"parameters": {
				  "name": "avoid tag names",
				  "negate": true,
				  "operator": "contains",
				  "pattern": "github"
				}
			  },
			  {
			    "type": "code_scanning",
			    "parameters": {
				  "code_scanning_tools": [
				    {
					  "tool": "CodeQL",
					  "security_alerts_threshold": "high_or_higher",
					  "alerts_threshold": "errors"
				    }
				  ]
			    }
			  }
			]
		  }`)
	})

	ctx := context.Background()
	ruleset, _, err := client.Organizations.CreateOrganizationRuleset(ctx, "o", Ruleset{
		ID:          Ptr(int64(21)),
		Name:        "ruleset",
		Target:      Ptr("branch"),
		SourceType:  Ptr("Organization"),
		Source:      "o",
		Enforcement: "active",
		BypassActors: []*BypassActor{
			{
				ActorID:   Ptr(int64(234)),
				ActorType: Ptr("Team"),
			},
		},
		Conditions: &RulesetConditions{
			RefName: &RulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
			RepositoryID: &RulesetRepositoryIDsConditionParameters{
				RepositoryIDs: []int64{123, 456},
			},
		},
		Rules: []*RepositoryRule{
			NewCreationRule(),
			NewUpdateRule(&UpdateAllowsFetchAndMergeRuleParameters{
				UpdateAllowsFetchAndMerge: true,
			}),
			NewDeletionRule(),
			NewRequiredLinearHistoryRule(),
			NewRequiredDeploymentsRule(&RequiredDeploymentEnvironmentsRuleParameters{
				RequiredDeploymentEnvironments: []string{"test"},
			}),
			NewRequiredSignaturesRule(),
			NewPullRequestRule(&PullRequestRuleParameters{
				AllowedMergeMethods:            []MergeMethod{"rebase", "squash"},
				DismissStaleReviewsOnPush:      true,
				RequireCodeOwnerReview:         true,
				RequireLastPushApproval:        true,
				RequiredApprovingReviewCount:   1,
				RequiredReviewThreadResolution: true,
			}),
			NewRequiredStatusChecksRule(&RequiredStatusChecksRuleParameters{
				RequiredStatusChecks: []RuleRequiredStatusChecks{
					{
						Context:       "test",
						IntegrationID: Ptr(int64(1)),
					},
				},
				StrictRequiredStatusChecksPolicy: true,
			}),
			NewNonFastForwardRule(),
			NewCommitMessagePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid test commits"),
				Negate:   Ptr(true),
				Operator: "starts_with",
				Pattern:  "[test]",
			}),
			NewCommitAuthorEmailPatternRule(&RulePatternParameters{
				Operator: "contains",
				Pattern:  "github",
			}),
			NewCommitterEmailPatternRule(&RulePatternParameters{
				Name:     Ptr("avoid commit emails"),
				Negate:   Ptr(true),
				Operator: "ends_with",
				Pattern:  "abc",
			}),
			NewBranchNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid branch names"),
				Negate:   Ptr(true),
				Operator: "regex",
				Pattern:  "github$",
			}),
			NewTagNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid tag names"),
				Negate:   Ptr(true),
				Operator: "contains",
				Pattern:  "github",
			}),
			NewRequiredCodeScanningRule(&RequiredCodeScanningRuleParameters{
				RequiredCodeScanningTools: []*RuleRequiredCodeScanningTool{
					{
						Tool:                    "CodeQL",
						SecurityAlertsThreshold: "high_or_higher",
						AlertsThreshold:         "errors",
					},
				},
			}),
		},
	})
	if err != nil {
		t.Errorf("Organizations.CreateOrganizationRuleset returned error: %v", err)
	}

	want := &Ruleset{
		ID:          Ptr(int64(21)),
		Name:        "ruleset",
		Target:      Ptr("branch"),
		SourceType:  Ptr("Organization"),
		Source:      "o",
		Enforcement: "active",
		BypassActors: []*BypassActor{
			{
				ActorID:   Ptr(int64(234)),
				ActorType: Ptr("Team"),
			},
		},
		Conditions: &RulesetConditions{
			RefName: &RulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
			RepositoryID: &RulesetRepositoryIDsConditionParameters{
				RepositoryIDs: []int64{123, 456},
			},
		},
		Rules: []*RepositoryRule{
			NewCreationRule(),
			NewUpdateRule(&UpdateAllowsFetchAndMergeRuleParameters{
				UpdateAllowsFetchAndMerge: true,
			}),
			NewDeletionRule(),
			NewRequiredLinearHistoryRule(),
			NewRequiredDeploymentsRule(&RequiredDeploymentEnvironmentsRuleParameters{
				RequiredDeploymentEnvironments: []string{"test"},
			}),
			NewRequiredSignaturesRule(),
			NewPullRequestRule(&PullRequestRuleParameters{
				AllowedMergeMethods:            []MergeMethod{"rebase", "squash"},
				DismissStaleReviewsOnPush:      true,
				RequireCodeOwnerReview:         true,
				RequireLastPushApproval:        true,
				RequiredApprovingReviewCount:   1,
				RequiredReviewThreadResolution: true,
			}),
			NewRequiredStatusChecksRule(&RequiredStatusChecksRuleParameters{
				RequiredStatusChecks: []RuleRequiredStatusChecks{
					{
						Context:       "test",
						IntegrationID: Ptr(int64(1)),
					},
				},
				StrictRequiredStatusChecksPolicy: true,
			}),
			NewNonFastForwardRule(),
			NewCommitMessagePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid test commits"),
				Negate:   Ptr(true),
				Operator: "starts_with",
				Pattern:  "[test]",
			}),
			NewCommitAuthorEmailPatternRule(&RulePatternParameters{
				Operator: "contains",
				Pattern:  "github",
			}),
			NewCommitterEmailPatternRule(&RulePatternParameters{
				Name:     Ptr("avoid commit emails"),
				Negate:   Ptr(true),
				Operator: "ends_with",
				Pattern:  "abc",
			}),
			NewBranchNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid branch names"),
				Negate:   Ptr(true),
				Operator: "regex",
				Pattern:  "github$",
			}),
			NewTagNamePatternRule(&RulePatternParameters{
				Name:     Ptr("avoid tag names"),
				Negate:   Ptr(true),
				Operator: "contains",
				Pattern:  "github",
			}),
			NewRequiredCodeScanningRule(&RequiredCodeScanningRuleParameters{
				RequiredCodeScanningTools: []*RuleRequiredCodeScanningTool{
					{
						Tool:                    "CodeQL",
						SecurityAlertsThreshold: "high_or_higher",
						AlertsThreshold:         "errors",
					},
				},
			}),
		},
	}
	if !cmp.Equal(ruleset, want) {
		t.Errorf("Organizations.CreateOrganizationRuleset returned %+v, want %+v", ruleset, want)
	}

	const methodName = "CreateOrganizationRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.CreateOrganizationRuleset(ctx, "o", Ruleset{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_GetOrganizationRuleset(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/rulesets/26110", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"id": 26110,
			"name": "test ruleset",
			"target": "branch",
			"source_type": "Organization",
			"source": "o",
			"enforcement": "active",
			"bypass_mode": "none",
			"node_id": "nid",
			"_links": {
			  "self": {
				"href": "https://api.github.com/orgs/o/rulesets/26110"
			  }
			},
			"conditions": {
				"ref_name": {
				  "include": [
					"refs/heads/main",
					"refs/heads/master"
				  ],
				  "exclude": [
					"refs/heads/dev*"
				  ]
				},
				"repository_name": {
				  "include": [
					"important_repository",
					"another_important_repository"
				  ],
				  "exclude": [
					"unimportant_repository"
				  ],
				  "protected": true
				}
			  },
			  "rules": [
				{
				  "type": "creation"
				}
			  ]
		}`)
	})

	ctx := context.Background()
	rulesets, _, err := client.Organizations.GetOrganizationRuleset(ctx, "o", 26110)
	if err != nil {
		t.Errorf("Organizations.GetOrganizationRepositoryRuleset returned error: %v", err)
	}

	want := &Ruleset{
		ID:          Ptr(int64(26110)),
		Name:        "test ruleset",
		Target:      Ptr("branch"),
		SourceType:  Ptr("Organization"),
		Source:      "o",
		Enforcement: "active",
		NodeID:      Ptr("nid"),
		Links: &RulesetLinks{
			Self: &RulesetLink{HRef: Ptr("https://api.github.com/orgs/o/rulesets/26110")},
		},
		Conditions: &RulesetConditions{
			RefName: &RulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
			RepositoryName: &RulesetRepositoryNamesConditionParameters{
				Include:   []string{"important_repository", "another_important_repository"},
				Exclude:   []string{"unimportant_repository"},
				Protected: Ptr(true),
			},
		},
		Rules: []*RepositoryRule{
			NewCreationRule(),
		},
	}
	if !cmp.Equal(rulesets, want) {
		t.Errorf("Organizations.GetOrganizationRuleset returned %+v, want %+v", rulesets, want)
	}

	const methodName = "GetOrganizationRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetOrganizationRuleset(ctx, "o", 26110)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_GetOrganizationRulesetWithRepoPropCondition(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/rulesets/26110", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"id": 26110,
			"name": "test ruleset",
			"target": "branch",
			"source_type": "Organization",
			"source": "o",
			"enforcement": "active",
			"bypass_mode": "none",
			"node_id": "nid",
			"_links": {
			  "self": {
				"href": "https://api.github.com/orgs/o/rulesets/26110"
			  }
			},
			"conditions": {
			  "repository_property": {
				"exclude": [],
				"include": [
				  {
					"name": "testIncludeProp",
					"source": "custom",
					"property_values": [
					  "true"
					]
				  }
				]
			  }
			},
			"rules": [
			{
				"type": "creation"
			}
			]
		}`)
	})

	ctx := context.Background()
	rulesets, _, err := client.Organizations.GetOrganizationRuleset(ctx, "o", 26110)
	if err != nil {
		t.Errorf("Organizations.GetOrganizationRepositoryRuleset returned error: %v", err)
	}

	want := &Ruleset{
		ID:          Ptr(int64(26110)),
		Name:        "test ruleset",
		Target:      Ptr("branch"),
		SourceType:  Ptr("Organization"),
		Source:      "o",
		Enforcement: "active",
		NodeID:      Ptr("nid"),
		Links: &RulesetLinks{
			Self: &RulesetLink{HRef: Ptr("https://api.github.com/orgs/o/rulesets/26110")},
		},
		Conditions: &RulesetConditions{
			RepositoryProperty: &RulesetRepositoryPropertyConditionParameters{
				Include: []RulesetRepositoryPropertyTargetParameters{
					{
						Name:   "testIncludeProp",
						Source: Ptr("custom"),
						Values: []string{"true"},
					},
				},
				Exclude: []RulesetRepositoryPropertyTargetParameters{},
			},
		},
		Rules: []*RepositoryRule{
			NewCreationRule(),
		},
	}
	if !cmp.Equal(rulesets, want) {
		t.Errorf("Organizations.GetOrganizationRuleset returned %+v, want %+v", rulesets, want)
	}

	const methodName = "GetOrganizationRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.GetOrganizationRuleset(ctx, "o", 26110)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_UpdateOrganizationRuleset(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/rulesets/26110", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
			"id": 26110,
			"name": "test ruleset",
			"target": "branch",
			"source_type": "Organization",
			"source": "o",
			"enforcement": "active",
			"bypass_mode": "none",
			"node_id": "nid",
			"_links": {
			  "self": {
				"href": "https://api.github.com/orgs/o/rulesets/26110"
			  }
			},
			"conditions": {
				"ref_name": {
				  "include": [
					"refs/heads/main",
					"refs/heads/master"
				  ],
				  "exclude": [
					"refs/heads/dev*"
				  ]
				},
				"repository_name": {
				  "include": [
					"important_repository",
					"another_important_repository"
				  ],
				  "exclude": [
					"unimportant_repository"
				  ],
				  "protected": true
				}
			  },
			  "rules": [
				{
				  "type": "creation"
				}
			  ]
		}`)
	})

	ctx := context.Background()
	rulesets, _, err := client.Organizations.UpdateOrganizationRuleset(ctx, "o", 26110, Ruleset{
		Name:        "test ruleset",
		Target:      Ptr("branch"),
		Enforcement: "active",
		Conditions: &RulesetConditions{
			RefName: &RulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
			RepositoryName: &RulesetRepositoryNamesConditionParameters{
				Include:   []string{"important_repository", "another_important_repository"},
				Exclude:   []string{"unimportant_repository"},
				Protected: Ptr(true),
			},
		},
		Rules: []*RepositoryRule{
			NewCreationRule(),
		},
	})
	if err != nil {
		t.Errorf("Organizations.UpdateOrganizationRuleset returned error: %v", err)
	}

	want := &Ruleset{
		ID:          Ptr(int64(26110)),
		Name:        "test ruleset",
		Target:      Ptr("branch"),
		SourceType:  Ptr("Organization"),
		Source:      "o",
		Enforcement: "active",
		NodeID:      Ptr("nid"),
		Links: &RulesetLinks{
			Self: &RulesetLink{HRef: Ptr("https://api.github.com/orgs/o/rulesets/26110")},
		},
		Conditions: &RulesetConditions{
			RefName: &RulesetRefConditionParameters{
				Include: []string{"refs/heads/main", "refs/heads/master"},
				Exclude: []string{"refs/heads/dev*"},
			},
			RepositoryName: &RulesetRepositoryNamesConditionParameters{
				Include:   []string{"important_repository", "another_important_repository"},
				Exclude:   []string{"unimportant_repository"},
				Protected: Ptr(true),
			},
		},
		Rules: []*RepositoryRule{
			NewCreationRule(),
		},
	}
	if !cmp.Equal(rulesets, want) {
		t.Errorf("Organizations.UpdateOrganizationRuleset returned %+v, want %+v", rulesets, want)
	}

	const methodName = "UpdateOrganizationRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.UpdateOrganizationRuleset(ctx, "o", 26110, Ruleset{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_UpdateOrganizationRulesetWithRepoProp(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/rulesets/26110", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
			"id": 26110,
			"name": "test ruleset",
			"target": "branch",
			"source_type": "Organization",
			"source": "o",
			"enforcement": "active",
			"bypass_mode": "none",
			"node_id": "nid",
			"_links": {
			  "self": {
				"href": "https://api.github.com/orgs/o/rulesets/26110"
			  }
			},
			"conditions": {
			  "repository_property": {
				"exclude": [],
				"include": [
				  {
					"name": "testIncludeProp",
					"source": "custom",
					"property_values": [
					  "true"
					]
				  }
				]
			  }
			  },
			  "rules": [
				{
				  "type": "creation"
				}
			  ]
		}`)
	})

	ctx := context.Background()
	rulesets, _, err := client.Organizations.UpdateOrganizationRuleset(ctx, "o", 26110, Ruleset{
		Name:        "test ruleset",
		Target:      Ptr("branch"),
		Enforcement: "active",
		Conditions: &RulesetConditions{
			RepositoryProperty: &RulesetRepositoryPropertyConditionParameters{
				Include: []RulesetRepositoryPropertyTargetParameters{
					{
						Name:   "testIncludeProp",
						Source: Ptr("custom"),
						Values: []string{"true"},
					},
				},
				Exclude: []RulesetRepositoryPropertyTargetParameters{},
			},
		},
		Rules: []*RepositoryRule{
			NewCreationRule(),
		},
	})
	if err != nil {
		t.Errorf("Organizations.UpdateOrganizationRuleset returned error: %v", err)
	}

	want := &Ruleset{
		ID:          Ptr(int64(26110)),
		Name:        "test ruleset",
		Target:      Ptr("branch"),
		SourceType:  Ptr("Organization"),
		Source:      "o",
		Enforcement: "active",
		NodeID:      Ptr("nid"),
		Links: &RulesetLinks{
			Self: &RulesetLink{HRef: Ptr("https://api.github.com/orgs/o/rulesets/26110")},
		},
		Conditions: &RulesetConditions{
			RepositoryProperty: &RulesetRepositoryPropertyConditionParameters{
				Include: []RulesetRepositoryPropertyTargetParameters{
					{
						Name:   "testIncludeProp",
						Source: Ptr("custom"),
						Values: []string{"true"},
					},
				},
				Exclude: []RulesetRepositoryPropertyTargetParameters{},
			},
		},
		Rules: []*RepositoryRule{
			NewCreationRule(),
		},
	}
	if !cmp.Equal(rulesets, want) {
		t.Errorf("Organizations.UpdateOrganizationRuleset returned %+v, want %+v", rulesets, want)
	}

	const methodName = "UpdateOrganizationRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		got, resp, err := client.Organizations.UpdateOrganizationRuleset(ctx, "o", 26110, Ruleset{})
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, err
	})
}

func TestOrganizationsService_UpdateOrganizationRulesetClearBypassActor(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/rulesets/26110", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
			"id": 26110,
			"name": "test ruleset",
			"target": "branch",
			"source_type": "Organization",
			"source": "o",
			"enforcement": "active",
			"bypass_mode": "none",
			"conditions": {
				"repository_name": {
					"include": [
						"important_repository",
						"another_important_repository"
					],
					"exclude": [
						"unimportant_repository"
					],
					"protected": true
				},
			  "ref_name": {
					"include": [
						"refs/heads/main",
						"refs/heads/master"
					],
					"exclude": [
						"refs/heads/dev*"
					]
				}
			},
			"rules": [
			  {
					"type": "creation"
			  }
			]
		}`)
	})

	ctx := context.Background()

	_, err := client.Organizations.UpdateOrganizationRulesetClearBypassActor(ctx, "o", 26110)
	if err != nil {
		t.Errorf("Organizations.UpdateOrganizationRulesetClearBypassActor returned error: %v \n", err)
	}

	const methodName = "UpdateOrganizationRulesetClearBypassActor"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Organizations.UpdateOrganizationRulesetClearBypassActor(ctx, "o", 26110)
	})
}

func TestOrganizationsService_DeleteOrganizationRuleset(t *testing.T) {
	t.Parallel()
	client, mux, _ := setup(t)

	mux.HandleFunc("/orgs/o/rulesets/26110", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	ctx := context.Background()
	_, err := client.Organizations.DeleteOrganizationRuleset(ctx, "o", 26110)
	if err != nil {
		t.Errorf("Organizations.DeleteOrganizationRuleset returned error: %v", err)
	}

	const methodName = "DeleteOrganizationRuleset"

	testNewRequestAndDoFailure(t, methodName, client, func() (*Response, error) {
		return client.Organizations.DeleteOrganizationRuleset(ctx, "0", 26110)
	})
}
