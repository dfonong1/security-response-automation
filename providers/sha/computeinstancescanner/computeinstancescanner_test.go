package computeinstancescanner

import (
	"testing"

	"github.com/googlecloudplatform/security-response-automation/cloudfunctions/gce/removepublicip"
)

func TestReadFinding(t *testing.T) {
	const (
		publicIPAddressFinding = `{
			"notificationConfigName": "organizations/1055058813388/notificationConfigs/noticonf-active-001-id",
			"finding": {
				"name": "organizations/1055058813388/sources/1986930501971458034/findings/d7ef72093c8c1e4c135d4c43fa847b83",
				"parent": "organizations/1055058813388/sources/1986930501971458034",
				"resourceName": "//compute.googleapis.com/projects/sec-automation-dev/zones/us-central1-a/instances/4312755253150365851",
				"state": "ACTIVE",
				"category": "PUBLIC_IP_ADDRESS",
				"externalUri": "https://console.cloud.google.com/compute/instancesDetail/zones/us-central1-a/instances/remove-public-ip-test-vm",
				"sourceProperties": {
				  "ReactivationCount": 0,
				  "ExceptionInstructions": "Add the security mark \"allow_public_ip_address\" to the asset with a value of \"true\" to prevent this finding from being activated again.",
				  "SeverityLevel": "High",
				  "Recommendation": "If this is unintended, please go to https://console.cloud.google.com/compute/instancesDetail/zones/us-central1-a/instances/remove-public-ip-test-vm and click \"Edit\". For each interface under the \"Network interfaces\" heading, set \"External IP\" to \"None\" or \"Ephemeral\", then click \"Done\" and \"Save\".  If you would like to learn more about securing access to your infrastructure, see https://cloud.google.com/solutions/connecting-securely.",
				  "ProjectId": "sec-automation-dev",
				  "AssetCreationTime": "2019-10-04T10:50:45.017-07:00",
				  "ScannerName": "COMPUTE_INSTANCE_SCANNER",
				  "ScanRunId": "2019-10-10T00:01:51.204-07:00",
				  "Explanation": "To reduce the attack surface, avoid assigning public IP addresses to your VMs."
				},
				"securityMarks": {
				  "name": "organizations/1055058813388/sources/1986930501971458034/findings/d7ef72093c8c1e4c135d4c43fa847b83/securityMarks",
				  "marks": {
					"kieras-test": "true",
					"kieras-test2": "true"
				  }
				},
				"eventTime": "2019-10-10T07:01:51.204Z",
				"createTime": "2019-10-04T19:02:25.582Z"
			}
		}`
		publicIPAddressRemediated = `{
			"notificationConfigName": "organizations/1055058813388/notificationConfigs/noticonf-active-001-id",
			"finding": {
				"name": "organizations/1055058813388/sources/1986930501971458034/findings/d7ef72093c8c1e4c135d4c43fa847b83",
				"parent": "organizations/1055058813388/sources/1986930501971458034",
				"resourceName": "//compute.googleapis.com/projects/sec-automation-dev/zones/us-central1-a/instances/4312755253150365851",
				"state": "ACTIVE",
				"category": "PUBLIC_IP_ADDRESS",
				"externalUri": "https://console.cloud.google.com/compute/instancesDetail/zones/us-central1-a/instances/remove-public-ip-test-vm",
				"sourceProperties": {
				  "ReactivationCount": 0,
				  "ExceptionInstructions": "Add the security mark \"allow_public_ip_address\" to the asset with a value of \"true\" to prevent this finding from being activated again.",
				  "SeverityLevel": "High",
				  "Recommendation": "If this is unintended, please go to https://console.cloud.google.com/compute/instancesDetail/zones/us-central1-a/instances/remove-public-ip-test-vm and click \"Edit\". For each interface under the \"Network interfaces\" heading, set \"External IP\" to \"None\" or \"Ephemeral\", then click \"Done\" and \"Save\".  If you would like to learn more about securing access to your infrastructure, see https://cloud.google.com/solutions/connecting-securely.",
				  "ProjectId": "sec-automation-dev",
				  "AssetCreationTime": "2019-10-04T10:50:45.017-07:00",
				  "ScannerName": "COMPUTE_INSTANCE_SCANNER",
				  "ScanRunId": "2019-10-10T00:01:51.204-07:00",
				  "Explanation": "To reduce the attack surface, avoid assigning public IP addresses to your VMs."
				},
				"securityMarks": {
				  "name": "organizations/1055058813388/sources/1986930501971458034/findings/d7ef72093c8c1e4c135d4c43fa847b83/securityMarks",
				  "marks": {
					"sra-remediated-event-time": "2019-10-10T07:01:51.204Z"
				  }
				},
				"eventTime": "2019-10-10T07:01:51.204Z",
				"createTime": "2019-10-04T19:02:25.582Z"
			}
		}`
		errorMessage = "remediation ignored! Finding already processed and remediated. Security Mark: \"sra-remediated-event-time:  2019-10-10T07:01:51.204Z\""
	)
	extractedValues := &removepublicip.Values{
		ProjectID:    "sec-automation-dev",
		InstanceZone: "us-central1-a",
		InstanceID:   "4312755253150365851",
		Mark:         "2019-10-10T07:01:51.204Z",
		Name:         "organizations/1055058813388/sources/1986930501971458034/findings/d7ef72093c8c1e4c135d4c43fa847b83",
	}
	for _, tt := range []struct {
		name           string
		ruleName       string
		values         *removepublicip.Values
		bytes          []byte
		expectedErrMsg string
	}{
		{name: "read", ruleName: "public_ip_address", values: extractedValues, bytes: []byte(publicIPAddressFinding), expectedErrMsg: ""},
		{name: "remediated", ruleName: "", values: nil, bytes: []byte(publicIPAddressRemediated), expectedErrMsg: errorMessage},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r, err := New(tt.bytes)
			if err != nil && err.Error() != tt.expectedErrMsg {
				t.Errorf("%s failed: got:%q want:%q", tt.name, err, tt.expectedErrMsg)
			}
			if r != nil {
				if values := r.RemovePublicIP(); *values != *tt.values {
					t.Errorf("%s failed: got:%v want:%v", tt.name, values, tt.values)
				}
				if name := r.Name(tt.bytes); name != tt.ruleName {
					t.Errorf("%q got:%q want:%q", tt.name, name, tt.ruleName)
				}
			}
		})
	}
}
