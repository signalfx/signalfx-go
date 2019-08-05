package orgtoken

import (
	"encoding/json"
	"strings"
)

type Limit struct {
}

func (l *Limit) UnmarshalJSON(data []byte) error {
	contents := string(data)
	if strings.Contains(contents, "dpmQuota") {
		limit := &DpmLimit{}
		return json.Unmarshal(data, limit)
	} else {
		limit := &HostOrUsageLimit{}
		return json.Unmarshal(data, limit)
	}
}
