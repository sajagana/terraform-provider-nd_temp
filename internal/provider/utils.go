package provider

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	// "terraform-provider-hashicups/client"
	// "terraform-provider-hashicups/container"
	// "terraform-provider-hashicups/models"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func ContainsString(strings []string, matchString string) bool {
	for _, stringValue := range strings {
		if stringValue == matchString {
			return true
		}
	}
	return false
}

func GetMOName(dn string) string {
	splittedDn := strings.Split(dn, "/")
	if len(splittedDn) > 1 {
		return strings.Join(strings.Split(splittedDn[len(splittedDn)-1], "-")[1:], "-")
	}
	return splittedDn[0]
}

// func CheckDn(ctx context.Context, client *client.Client, dn string, diags *diag.Diagnostics) {
// 	tflog.Debug(ctx, fmt.Sprintf("validate relation dn: %s", dn))
// 	_, err := client.Get(dn)
// 	if err != nil {
// 		tflog.Error(ctx, fmt.Sprintf("failed validate relation dn: %s", dn))
// 		diags.AddError(
// 			"Relation target dn validation failed",
// 			fmt.Sprintf("The relation target dn is not found: %s", dn),
// 		)
// 	}
// }

func DoRestRequestEscapeHtml(ctx context.Context, diags *diag.Diagnostics, client *client.Client, path, method string, payload *container.Container, escapeHtml bool) *container.Container {
	log.Printf("[DEBUG] ------ inside DoRestRequestEscapeHtml escapeHtml: %v, path: %s", escapeHtml, path)

	// Ensure path starts with a slash to assure signature is created correctly
	if !strings.HasPrefix("/", path) {
		path = fmt.Sprintf("/%s", path)
	}
	var restRequest *http.Request
	var err error

	if escapeHtml {
		restRequest, err = client.MakeRestRequest(method, path, payload, true)
	} else {
		restRequest, err = client.MakeRestRequestRaw(method, path, payload.EncodeJSON(), true)
	}
	if err != nil {
		diags.AddError(
			"Creation of rest request failed",
			fmt.Sprintf("err: %s. Please report this issue to the provider developers.", err),
		)
		return nil
	}

	cont, restResponse, err := client.Do(restRequest)

	if restResponse != nil && cont.Data() != nil && restResponse.StatusCode != 200 {
		errCode := models.StripQuotes(models.StripSquareBrackets(cont.Search("imdata", "error", "attributes", "code").String()))
		if errCode != "1" && errCode != "103" && errCode != "107" && errCode != "120" {
			diags.AddError(
				fmt.Sprintf("The %s rest request failed", strings.ToLower(method)),
				fmt.Sprintf("Code: %d Response: %s, err: %s. Please report this issue to the provider developers.", restResponse.StatusCode, cont.Data().(map[string]interface{})["imdata"], err),
			)
			return nil
		}
		tflog.Debug(ctx, models.StripQuotes(models.StripSquareBrackets(cont.Search("imdata", "error", "attributes", "text").String())))
	} else if err != nil {
		diags.AddError(
			fmt.Sprintf("The %s rest request failed", strings.ToLower(method)),
			fmt.Sprintf("Err: %s. Please report this issue to the provider developers.", err),
		)
		return nil
	}

	return cont
}

func DoRestRequest(ctx context.Context, diags *diag.Diagnostics, client *client.Client, path, method string, payload *container.Container) *container.Container {
	return DoRestRequestEscapeHtml(ctx, diags, client, path, method, payload, true)
}

func GetDeleteJsonPayload(ctx context.Context, diags *diag.Diagnostics, className, dn string) *container.Container {

	jsonString := fmt.Sprintf(`{"%s":{"attributes":{"dn": "%s","status": "deleted"}}}`, className, dn)
	jsonPayload, err := container.ParseJSON([]byte(jsonString))
	if err != nil {
		diags.AddError(
			"Construction of json payload failed",
			fmt.Sprintf("Err: %s. Please report this issue to the provider developers.", err),
		)
		return nil
	}
	return jsonPayload
}
