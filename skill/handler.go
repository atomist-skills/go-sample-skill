/*
 * Copyright © 2022 Atomist, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package skill

import (
	"encoding/json"
)

// Decode an incoming subscription payload into the concrete
// mapping type
func Decode[P interface{}](event map[string]json.RawMessage) P {
	jsonbody, _ := json.Marshal(event)
	var decoded P
	json.Unmarshal(jsonbody, &decoded)
	return decoded
}

func DecodedEventHandler[P interface{}](delegate func(payload P) Status) EventHandler {
	return func(ctx EventContext) Status {
		var status Status
		for _, e := range ctx.Data {
			decodedEvent := Decode[P](e[0])
			status = delegate(decodedEvent)
		}
		return status
	}
}
