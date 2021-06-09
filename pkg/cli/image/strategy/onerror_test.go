package strategy

import (
	"context"
	"reflect"
	"strings"
	"testing"

	operatorv1alpha1 "github.com/openshift/api/operator/v1alpha1"
	"github.com/openshift/library-go/pkg/image/reference"
)

func TestOnErrorStrategy(t *testing.T) {
	tests := []struct {
		name                 string
		icspList             []operatorv1alpha1.ImageContentSourcePolicy
		image                string
		imageSourcesExpected []string
	}{
		{
			name: "multiple ICSPs",
			icspList: []operatorv1alpha1.ImageContentSourcePolicy{
				{
					Spec: operatorv1alpha1.ImageContentSourcePolicySpec{
						RepositoryDigestMirrors: []operatorv1alpha1.RepositoryDigestMirrors{
							{
								Source: "quay.io/multiple/icsps",
								Mirrors: []string{
									"someregistry/somerepo/release",
								},
							},
							{
								Source: "quay.io/ocp-test/another-release",
								Mirrors: []string{
									"someregistry/repo/does-not-exist",
								},
							},
						},
					},
				},
				{
					Spec: operatorv1alpha1.ImageContentSourcePolicySpec{
						RepositoryDigestMirrors: []operatorv1alpha1.RepositoryDigestMirrors{
							{
								Source: "quay.io/multiple/icsps",
								Mirrors: []string{
									"anotherregistry/anotherrepo/release",
								},
							},
						},
					},
				},
			},
			image:                "quay.io/multiple/icsps:4.5",
			imageSourcesExpected: []string{"quay.io/multiple/icsps", "someregistry/somerepo/release", "anotherregistry/anotherrepo/release"},
		},
		{
			name: "multiple mirrors, single source match",
			icspList: []operatorv1alpha1.ImageContentSourcePolicy{
				{
					Spec: operatorv1alpha1.ImageContentSourcePolicySpec{
						RepositoryDigestMirrors: []operatorv1alpha1.RepositoryDigestMirrors{
							{
								Source: "docker.io/ocp-test/does-not-exist",
								Mirrors: []string{
									"does.not.exist/match/image",
								},
							},
							{
								Source: "quay.io/ocp-test/does-not-exist",
								Mirrors: []string{
									"exists/match/image",
								},
							},
						},
					},
				},
			},
			image:                "quay.io/ocp-test/does-not-exist:4.7",
			imageSourcesExpected: []string{"quay.io/ocp-test/does-not-exist", "exists/match/image"},
		},
		{
			name: "single mirror and match",
			icspList: []operatorv1alpha1.ImageContentSourcePolicy{
				{
					Spec: operatorv1alpha1.ImageContentSourcePolicySpec{
						RepositoryDigestMirrors: []operatorv1alpha1.RepositoryDigestMirrors{
							{
								Source: "quay.io/ocp-test/release",
								Mirrors: []string{
									"someregistry/mirrors/match",
								},
							},
						},
					},
				},
			},
			image:                "quay.io/ocp-test/release:4.5",
			imageSourcesExpected: []string{"quay.io/ocp-test/release", "someregistry/mirrors/match"},
		},
		{
			name: "no source match",
			icspList: []operatorv1alpha1.ImageContentSourcePolicy{
				{
					Spec: operatorv1alpha1.ImageContentSourcePolicySpec{
						RepositoryDigestMirrors: []operatorv1alpha1.RepositoryDigestMirrors{
							{
								Source: "docker.io/ocp-test/does-not-exist",
								Mirrors: []string{
									"does.not.exist/match/image",
								},
							},
							{
								Source: "quay.io/ocp-test/does-not-exist",
								Mirrors: []string{
									"exists/match/image",
								},
							},
						},
					},
				},
			},
			image:                "quay.io/passed/image:4.5",
			imageSourcesExpected: []string{"quay.io/passed/image"},
		},
		{
			name: "multiple mirrors for single source match",
			icspList: []operatorv1alpha1.ImageContentSourcePolicy{
				{
					Spec: operatorv1alpha1.ImageContentSourcePolicySpec{
						RepositoryDigestMirrors: []operatorv1alpha1.RepositoryDigestMirrors{
							{
								Source: "quay.io/ocp-test/release",
								Mirrors: []string{
									"someregistry/mirrors/match",
									"quay.io/another/release",
									"quay.io/andanother/release",
								},
							},
						},
					},
				},
			},
			image:                "quay.io/ocp-test/release:4.5",
			imageSourcesExpected: []string{"quay.io/ocp-test/release", "someregistry/mirrors/match", "quay.io/another/release", "quay.io/andanother/release"},
		},
		{
			name: "docker.io vs registry-1.docker.io",
			icspList: []operatorv1alpha1.ImageContentSourcePolicy{
				{
					Spec: operatorv1alpha1.ImageContentSourcePolicySpec{
						RepositoryDigestMirrors: []operatorv1alpha1.RepositoryDigestMirrors{
							{
								Source: "docker.io/ocp-test/release",
								Mirrors: []string{
									"quay.io/ocp-test/release",
								},
							},
						},
					},
				},
			},
			image:                "registry-1.docker.io/ocp-test/release:4.5",
			imageSourcesExpected: []string{"registry-1.docker.io/ocp-test/release", "quay.io/ocp-test/release"},
		},
		{
			name:                 "no ICSP",
			image:                "quay.io/ocp-test/release:4.5",
			imageSourcesExpected: []string{"quay.io/ocp-test/release"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := []reference.DockerImageReference{}
			for _, e := range tt.imageSourcesExpected {
				ref, _ := reference.Parse(e)
				expected = append(expected, ref)
			}

			alternates := NewICSPOnErrorStrategy("name")
			readCount := 0
			onErr := alternates.(*onErrorStrategy)
			onErr.readICSPsFromFileFunc = func(string) ([]operatorv1alpha1.ImageContentSourcePolicy, error) {
				readCount++
				return tt.icspList, nil
			}
			imageRef, _ := reference.Parse(tt.image)

			actual, err := alternates.FirstRequest(context.Background(), imageRef)
			if actual != nil || err != nil {
				t.Errorf("Unexpected values returned from FirstRequest\nactual: %v\nerr: %v", actual, err)
			}

			actual, err = alternates.OnFailure(context.Background(), imageRef)
			if err != nil {
				t.Errorf("Unexpected error %v", err)
				return
			}
			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("Unexpected alternates got = %v, want %v", actual, expected)
			}
			if readCount > 1 {
				t.Errorf("Unexpected number of ICSP reads, should be 1, got %d", readCount)
			}
		})
	}
}

func TestOnErrorStrategyErrors(t *testing.T) {
	ref, _ := reference.Parse("quay.io/ocp-test/release:4.5")

	alternates := NewICSPOnErrorStrategy("")
	_, err := alternates.OnFailure(context.Background(), ref)
	if err == nil || !strings.Contains(err.Error(), "no ImageContentSourceFile") {
		t.Errorf("Expected error empty ICSP file error, got %v", err)
	}
}
