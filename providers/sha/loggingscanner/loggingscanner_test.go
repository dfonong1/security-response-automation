package loggingscanner

import (
	"testing"

	"github.com/googlecloudplatform/security-response-automation/cloudfunctions/iam/enableauditlogs"
)

func TestReadFinding(t *testing.T) {
	const (
		auditLogDisabledFinding = `{
		"finding": {
			"name": "organizations/154584661726/sources/1986930501971458034/findings/1c35bd4b4f6d7145e441f2965c32f074",
			"parent": "organizations/154584661726/sources/1986930501971458034",
			"resourceName": "//cloudresourcemanager.googleapis.com/projects/108906606255",
			"state": "ACTIVE",
			"category": "AUDIT_LOGGING_DISABLED",
			"externalUri": "https://console.cloud.google.com/iam-admin/audit/allservices?project=test-project",
			"sourceProperties": {
				"ReactivationCount": 0,
				"ExceptionInstructions": "Add the security mark \"allow_audit_logging_disabled\" to the asset with a value of \"true\" to prevent this finding from being activated again.",
				"SeverityLevel": "Low",
				"Recommendation": "Go to https://console.cloud.google.com/iam-admin/audit/allservices?project=test-project and under \"LOG TYPE\" select \"Admin read\", \"Data read\", and \"Data write\", and then click \"SAVE\". Make sure there are no exempted users configured.",
				"ProjectId": "test-project",
				"AssetCreationTime": "2019-10-22T15:13:39.305Z",
				"ScannerName": "LOGGING_SCANNER",
				"ScanRunId": "2019-10-22T14:01:08.832-07:00",
				"Explanation": "You should enable Cloud Audit Logging for all services, to track all Admin activities including read and write access to user data."
			},
			"securityMarks": {
				"name": "organizations/154584661726/sources/1986930501971458034/findings/1c35bd4b4f6d7145e441f2965c32f074/securityMarks"
			},
			"eventTime": "2019-10-22T21:01:08.832Z",
			"createTime": "2019-10-22T21:01:39.098Z",
			"assetId": "organizations/154584661726/assets/11190834741917282179",
			"assetDisplayName": "test-project"
		   }
		}`
		auditLogDisabledRemediated = `{
		"finding": {
			"name": "organizations/154584661726/sources/1986930501971458034/findings/1c35bd4b4f6d7145e441f2965c32f074",
			"parent": "organizations/154584661726/sources/1986930501971458034",
			"resourceName": "//cloudresourcemanager.googleapis.com/projects/108906606255",
			"state": "ACTIVE",
			"category": "AUDIT_LOGGING_DISABLED",
			"externalUri": "https://console.cloud.google.com/iam-admin/audit/allservices?project=test-project",
			"sourceProperties": {
				"ReactivationCount": 0,
				"ExceptionInstructions": "Add the security mark \"allow_audit_logging_disabled\" to the asset with a value of \"true\" to prevent this finding from being activated again.",
				"SeverityLevel": "Low",
				"Recommendation": "Go to https://console.cloud.google.com/iam-admin/audit/allservices?project=test-project and under \"LOG TYPE\" select \"Admin read\", \"Data read\", and \"Data write\", and then click \"SAVE\". Make sure there are no exempted users configured.",
				"ProjectId": "test-project",
				"AssetCreationTime": "2019-10-22T15:13:39.305Z",
				"ScannerName": "LOGGING_SCANNER",
				"ScanRunId": "2019-10-22T14:01:08.832-07:00",
				"Explanation": "You should enable Cloud Audit Logging for all services, to track all Admin activities including read and write access to user data."
			},
			"securityMarks": {
				"name": "organizations/154584661726/sources/1986930501971458034/findings/1c35bd4b4f6d7145e441f2965c32f074/securityMarks",
				"marks": {
					"sra-remediated-event-time": "2019-10-22T21:01:08.832Z"
			  	}
			},
			"eventTime": "2019-10-22T21:01:08.832Z",
			"createTime": "2019-10-22T21:01:39.098Z",
			"assetId": "organizations/154584661726/assets/11190834741917282179",
			"assetDisplayName": "test-project"
		   }
		}`
		errorMessage = "remediation ignored! Finding already processed and remediated. Security Mark: \"sra-remediated-event-time: 2019-10-22T21:01:08.832Z\""
	)
	extractedValues := &enableauditlogs.Values{
		ProjectID: "test-project",
		Mark:      "2019-10-22T21:01:08.832Z",
		Name:      "organizations/154584661726/sources/1986930501971458034/findings/1c35bd4b4f6d7145e441f2965c32f074",
	}
	for _, tt := range []struct {
		name           string
		ruleName       string
		values         *enableauditlogs.Values
		bytes          []byte
		expectedErrMsg string
	}{
		{name: "read", ruleName: "audit_logging_disabled", values: extractedValues, bytes: []byte(auditLogDisabledFinding), expectedErrMsg: ""},
		{name: "remediated", ruleName: "", values: nil, bytes: []byte(auditLogDisabledRemediated), expectedErrMsg: errorMessage},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r, err := New(tt.bytes)
			if err != nil && err.Error() != tt.expectedErrMsg {
				t.Errorf("%s failed: got:%q want:%q", tt.name, err, tt.expectedErrMsg)
			}
			if r != nil {
				values := r.EnableAuditLogs()
				if *values != *tt.values {
					t.Errorf("%s failed: got:%v want:%v", tt.name, values, tt.values)
				}
				if name := r.Name(tt.bytes); name != tt.ruleName {
					t.Errorf("%q got:%q want:%q", tt.name, name, tt.ruleName)
				}
			}
		})
	}
}
