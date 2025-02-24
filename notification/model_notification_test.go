package notification

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotificationUnmarshaJSON(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name   string
		input  string
		expect *Notification
		errVal string
	}{
		{
			name:   "no type value defined",
			input:  `{}`,
			expect: nil,
			errVal: "unknown notification type \"\"",
		},
		{
			name:  "AWS Event Bridge",
			input: `{"type":"AmazonEventBridge"}`,
			expect: &Notification{Type: "AmazonEventBridge", Value: &AmazonEventBrigeNotification{
				Type: "AmazonEventBridge",
			}},
			errVal: "",
		},
		{
			name:  "BigPanda",
			input: `{"type":"BigPanda"}`,
			expect: &Notification{Type: "BigPanda", Value: &BigPandaNotification{
				Type: "BigPanda",
			}},
			errVal: "",
		},
		{
			name:  "Email",
			input: `{"type":"Email"}`,
			expect: &Notification{Type: "Email", Value: &EmailNotification{
				Type: "Email",
			}},
			errVal: "",
		},
		{
			name:  "Jira",
			input: `{"type":"Jira"}`,
			expect: &Notification{Type: "Jira", Value: &JiraNotification{
				Type: "Jira",
			}},
			errVal: "",
		},
		{
			name:  "Office365",
			input: `{"type":"Office365"}`,
			expect: &Notification{Type: "Office365", Value: &Office365Notification{
				Type: "Office365",
			}},
			errVal: "",
		},
		{
			name:  "Opsgenie",
			input: `{"type":"Opsgenie"}`,
			expect: &Notification{Type: "Opsgenie", Value: &OpsgenieNotification{
				Type: "Opsgenie",
			}},
			errVal: "",
		},
		{
			name:  "PagerDuty",
			input: `{"type":"PagerDuty"}`,
			expect: &Notification{Type: "PagerDuty", Value: &PagerDutyNotification{
				Type: "PagerDuty",
			}},
			errVal: "",
		},
		{
			name:  "ServiceNow",
			input: `{"type":"ServiceNow"}`,
			expect: &Notification{Type: "ServiceNow", Value: &ServiceNowNotification{
				Type: "ServiceNow",
			}},
			errVal: "",
		},
		{
			name:  "Slack",
			input: `{"type":"Slack"}`,
			expect: &Notification{Type: "Slack", Value: &SlackNotification{
				Type: "Slack",
			}},
			errVal: "",
		},
		{
			name:  "Team",
			input: `{"type":"Team"}`,
			expect: &Notification{Type: "Team", Value: &TeamNotification{
				Type: "Team",
			}},
			errVal: "",
		},
		{
			name:  "TeamEmail",
			input: `{"type":"TeamEmail"}`,
			expect: &Notification{Type: "TeamEmail", Value: &TeamEmailNotification{
				Type: "TeamEmail",
			}},
			errVal: "",
		},
		{
			name:  "VictorOps",
			input: `{"type":"VictorOps"}`,
			expect: &Notification{Type: "VictorOps", Value: &VictorOpsNotification{
				Type: "VictorOps",
			}},
			errVal: "",
		},
		{
			name:  "Webhook",
			input: `{"type":"Webhook"}`,
			expect: &Notification{Type: "Webhook", Value: &WebhookNotification{
				Type: "Webhook",
			}},
			errVal: "",
		},
		{
			name:  "XMatters",
			input: `{"type":"XMatters"}`,
			expect: &Notification{Type: "XMatters", Value: &XMattersNotification{
				Type: "XMatters",
			}},
			errVal: "",
		},
		{
			name:  "SplunkPlatform",
			input: `{"type":"SplunkPlatform"}`,
			expect: &Notification{Type: "SplunkPlatform", Value: &SplunkPlatformNotification{
				Type: "SplunkPlatform",
			}},
			errVal: "",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var actual Notification
			err := json.NewDecoder(strings.NewReader(tc.input)).Decode(&actual)
			if tc.errVal != "" {
				assert.EqualError(t, err, tc.errVal, "Must match the expected error value")
			} else {
				assert.Equal(t, tc.expect, &actual, "Must match the expected notification value")
				assert.NoError(t, err, "Must not error")
			}
		})
	}
}
