package provider

import (
	"context"
	"os"
	"strconv"
	"testing"

	// "github.com/CiscoDevNet/terraform-provider-nd/internal/provider/nd"
	"terraform-provider-nd/internal/provider/nd"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"nd": func() (tfprotov6.ProviderServer, error) {
		ctx := context.Background()

		upgradedSdkServer, err := tf5to6server.UpgradeServer(
			ctx,
			nd.Provider().GRPCProvider,
		)

		if err != nil {
			return nil, err
		}

		providers := []func() tfprotov6.ProviderServer{
			providerserver.NewProtocol6(New("test")()),
			func() tfprotov6.ProviderServer {
				return upgradedSdkServer
			},
		}

		muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)

		if err != nil {
			return nil, err
		}

		return muxServer.ProviderServer(), nil
	},
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("ND_USERNAME"); v == "" {
		t.Fatal("ND_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("ND_PASSWORD"); v == "" {
		t.Fatal("ND_PASSWORD must be set for acceptance tests")
	}
	if v := os.Getenv("ND_URL"); v == "" {
		t.Fatal("ND_URL must be set for acceptance tests")
	}
	if v := os.Getenv("ND_VAL_REL_DN"); v == "" {
		t.Fatal("ND_VAL_REL_DN must be set for acceptance tests")
		boolValue, err := strconv.ParseBool(v)
		if err != nil || boolValue == true {
			t.Fatal("ND_VAL_REL_DN must be a 'false' boolean value")
		}
	}
}

func setGlobalAnnotationEnvVariable(t *testing.T, annotation string) {
	t.Setenv("ND_ANNOTATION", annotation)
}
